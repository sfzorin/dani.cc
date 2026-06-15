// Morning workout player — runs entirely client-side from a server-provided
// payload. It steps through exercises (reps + timed), rest periods, voice cues
// and the celebratory finish, then reports completion to the server via $hook.
//
// Loaded as a managed inline script: $data / $hook / $sys are injected by doors.

(function () {
  var payload = JSON.parse($data("payload"));
  var items = payload.items || [];
  var labels = payload.labels || {};
  var cues = payload.cues || {};
  var enc = payload.enc || [];
  var voiceMode = payload.voiceMode || "normal";   // off | min | normal | detailed
  var voiceOn = voiceMode !== "off";
  var voiceLang = payload.bcp47 || "en-US";        // e.g. "tr-TR"
  var voicePrefix = (payload.lang || "en").toLowerCase(); // e.g. "tr"

  // Voice selection. Browsers otherwise default to an English voice even when
  // utterance.lang is "tr-TR", so we explicitly pick a voice whose lang matches
  // the target language (this is the Turkish-falls-back-to-English fix).
  var voices = [];
  var chosenVoice = null;
  function loadVoices() {
    if (!("speechSynthesis" in window)) return;
    try { voices = window.speechSynthesis.getVoices() || []; } catch (e) { voices = []; }
    chosenVoice = null;
    var i;
    for (i = 0; i < voices.length; i++) { // exact region match (tr-TR)
      if (voices[i].lang && voices[i].lang.toLowerCase() === voiceLang.toLowerCase()) { chosenVoice = voices[i]; return; }
    }
    for (i = 0; i < voices.length; i++) { // any voice of the language (tr-*)
      var lng = (voices[i].lang || "").toLowerCase().replace("_", "-");
      if (lng === voicePrefix || lng.indexOf(voicePrefix + "-") === 0) { chosenVoice = voices[i]; return; }
    }
  }
  if ("speechSynthesis" in window) {
    loadVoices();
    window.speechSynthesis.onvoiceschanged = loadVoices;
  }

  var root = document.currentScript.parentElement; // .player
  var $ = function (sel) { return root.querySelector(sel); };

  var el = {
    pbar: $(".pbar-fill"),
    phase: $(".phase"),
    index: $(".ex-index"),
    svg: $(".ex-svg"),
    name: $(".ex-name"),
    side: $(".ex-side"),
    value: $(".ex-value"),
    hint: $(".ex-hint"),
    warn: $(".ex-warn"),
    stage: $(".stage"),
    infoBtn: $(".info-btn"),
    infoOverlay: $(".info-overlay"),
    infoMedia: $(".info-media"),
    infoName: $(".info-name"),
    infoDesc: $(".info-desc"),
    infoHow: $(".info-how"),
    infoHowSec: $(".info-how-sec"),
    infoCorrect: $(".info-correct"),
    infoCorrectSec: $(".info-correct-sec"),
    infoWrong: $(".info-wrong"),
    infoWrongSec: $(".info-wrong-sec"),
    infoBreathing: $(".info-breathing"),
    infoBreathingSec: $(".info-breathing-sec"),
    infoWarn: $(".info-warn"),
    infoWarnSec: $(".info-warn-sec"),
    infoReplace: $(".info-replace"),
    infoClose: $(".info-close"),
    controls: $(".controls"),
    prev: $(".ctl-prev"),
    pause: $(".ctl-pause"),
    skip: $(".ctl-skip"),
    done: $(".ctl-done"),
    rest: $(".rest-overlay"),
    restCount: $(".rest-count"),
    restNext: $(".rest-next"),
    restSkip: $(".rest-skip"),
    restBack: $(".rest-back"),
    restTitle: $(".rest-title"),
    doneOverlay: $(".done-overlay"),
    doneTitle: $(".done-title"),
    doneEnc: $(".done-enc"),
    doneStreak: $(".done-streak")
  };
  var vis = function (e) { return e && !e.classList.contains("hidden"); };

  var i = 0;          // current item index
  var paused = false;
  var ticker = null;  // active 1s interval
  var finished = false;
  var wakeLock = null;

  // ---- helpers ----------------------------------------------------------
  function t(key) { return labels[key] || key; }

  function slotLabel(slot) {
    if (slot === "warmup") return t("warmup");
    if (slot === "cooldown") return t("cooldown");
    return t("main");
  }

  // ---- speech coordinator ----------------------------------------------
  // Two priorities, so the spoken countdown always lands on its second and
  // never queues behind technique narration:
  //   sayNow(text) — time-critical (countdown digits, Go, Rest, next-exercise,
  //                  half-time): cancels whatever is speaking, then speaks.
  //   utter(text)  — plain queued speech (narration, the rest announcement,
  //                  finish messages); any sayNow cuts the current line off.
  function utter(text) {
    if (!voiceOn || !text || !("speechSynthesis" in window)) return;
    try {
      if (!voices.length || !chosenVoice) loadVoices();
      var u = new SpeechSynthesisUtterance(text);
      if (chosenVoice) {
        u.voice = chosenVoice;
        u.lang = chosenVoice.lang;
      } else {
        u.lang = voiceLang; // no matching voice installed on this device
      }
      u.rate = 1.0;
      window.speechSynthesis.speak(u);
    } catch (e) {}
  }
  var speak = utter; // queued (finish / celebration messages)
  function sayNow(text) {
    if (!voiceOn || !text) return;
    shutUp();      // drop any in-progress narration so the cue is on time
    utter(text);
  }
  function shutUp() {
    if ("speechSynthesis" in window) { try { window.speechSynthesis.cancel(); } catch (e) {} }
  }

  function clearTicker() {
    if (ticker) { clearInterval(ticker); ticker = null; }
  }

  // talkative covers the two spoken levels: "detailed" (chatty) and "normal".
  function talkative() { return voiceMode === "normal" || voiceMode === "detailed"; }

  // ---- scheduled speech: technique narration + the rest announcement -------
  // narrTimer drives the "how to / how not to" lines; announceTimer holds the
  // delayed "Rest / next exercise" line. clearNarr cancels both — they are the
  // only pre-scheduled speech, and both must die on every transition.
  var narrTimer = null, narrQ = [], narrI = 0, announceTimer = null;
  function clearNarr() {
    if (narrTimer) { clearTimeout(narrTimer); narrTimer = null; }
    if (announceTimer) { clearTimeout(announceTimer); announceTimer = null; }
    narrQ = []; narrI = 0;
  }
  // startNarr begins the narration a few seconds into the exercise (after the Go
  // cue), spaced out and low priority. runSide stops it before the closing
  // countdown; any sayNow cue cuts the current line off. detailed = all cues.
  function startNarr(it) {
    if (!voiceOn || !talkative()) return;
    var c = it.correct || [], w = it.wrong || [], k;
    var nc = voiceMode === "detailed" ? c.length : Math.min(2, c.length);
    var nw = voiceMode === "detailed" ? w.length : Math.min(1, w.length);
    narrQ = []; narrI = 0;
    for (k = 0; k < nc; k++) narrQ.push(c[k]); // how to do it (phrased as a do)
    for (k = 0; k < nw; k++) narrQ.push(w[k]); // common mistake (phrased as a don't)
    if (!narrQ.length) return;
    var gap = voiceMode === "detailed" ? 6500 : 11000;
    function step() {
      if (paused) { narrTimer = setTimeout(step, 1000); return; }
      if (narrI >= narrQ.length) { narrTimer = null; return; }
      utter(narrQ[narrI]); narrI++;
      narrTimer = setTimeout(step, gap);
    }
    narrTimer = setTimeout(step, 3000); // small pause after "Go" before instructions
  }

  // flashHalf marks the midpoint of a timed set: shows and (when talkative)
  // speaks "half the time", clear of both the start and the closing count.
  var halfTimer = null;
  function flashHalf() {
    el.side.textContent = cues.halfway || "";
    if (halfTimer) clearTimeout(halfTimer);
    halfTimer = setTimeout(function () { el.side.textContent = ""; }, 3500);
    if (talkative()) sayNow(cues.halfway);
  }

  // Run a 1-second countdown; onTick(remaining) each second, onDone() at 0.
  // Respects the paused flag.
  function countdown(seconds, onTick, onDone) {
    clearTicker();
    var remaining = seconds;
    onTick(remaining);
    ticker = setInterval(function () {
      if (paused) return;
      remaining -= 1;
      if (remaining <= 0) {
        onTick(0);
        clearTicker();
        onDone();
      } else {
        onTick(remaining);
      }
    }, 1000);
  }

  async function requestWake() {
    try {
      if ("wakeLock" in navigator) wakeLock = await navigator.wakeLock.request("screen");
    } catch (e) {}
  }
  function releaseWake() {
    try { if (wakeLock) { wakeLock.release(); wakeLock = null; } } catch (e) {}
  }

  // ---- rendering --------------------------------------------------------
  function setProgress() {
    var pct = Math.round(((i) / items.length) * 100);
    el.pbar.style.width = pct + "%";
  }

  var pendingStart = null; // start-fn the rest/ready countdown will run at 0
  var curSide = 1;         // retained so the rep "Done" handler can call nextSide

  function startItem(idx) {
    i = idx;
    var it = items[i];
    root.classList.remove("frozen"); // a fresh set is always running, not paused
    el.rest.classList.add("hidden");
    el.stage.classList.remove("hidden");
    el.controls.classList.remove("hidden");

    el.svg.src = it.svg;
    el.svg.alt = it.name;
    el.name.textContent = it.name;
    el.hint.textContent = it.hint || "";
    el.warn.textContent = it.warning ? "⚠️ " + it.warning : "";
    el.index.textContent = (i + 1) + " / " + items.length;
    el.phase.textContent = slotLabel(it.slot) + " · " + t("round") + " " + it.round;
    // Slot-coloured counter so you can see at a glance whether you're in the
    // warm-up, main set, or cool-down.
    el.value.classList.remove("slot-warmup", "slot-main", "slot-cooldown");
    el.value.classList.add("slot-" + it.slot);
    el.phase.classList.remove("slot-warmup", "slot-main", "slot-cooldown");
    el.phase.classList.add("slot-" + it.slot);
    setProgress();
    runSide(it, 1);
  }

  // runSide runs the exercise. Per-side splitting was removed — every exercise is
  // a single set in seconds or reps; the user does both sides themselves. Voice:
  // a start cue, technique narration, "half the time" at the midpoint of a timed
  // exercise, and a 5-second warning.
  function runSide(it, side) {
    curSide = side;
    clearNarr();
    el.side.textContent = "";
    sayNow(it.vStart);   // "Go" / "Старт"
    startNarr(it);       // first instruction ~3 s later

    if (it.unit === "seconds") {
      el.done.classList.add("hidden");
      el.value.classList.add("timer");
      var D = it.value;
      var halfAt = Math.round(D * 0.5);
      countdown(D, function (r) {
        el.value.textContent = r + "″";
        if (D >= 12 && r === halfAt && r > 5) flashHalf(); // mid-set, clear of the count
        if (r === 4) clearNarr();                // silence narration before the count
        if (r <= 3 && r >= 1) sayNow(String(r)); // digits land exactly on the second
      }, function () { nextSide(it, side); });
    } else {
      // reps / breaths: wait for the user to tap Done.
      clearTicker();
      el.value.classList.remove("timer");
      el.value.textContent = it.unit === "breaths"
        ? it.value + " " + t("breaths")
        : "× " + it.value;
      el.done.classList.remove("hidden");
    }
  }

  function nextSide(it, side) {
    clearNarr();
    endItem();
  }

  // unitFor returns the spoken / written unit word for an item.
  function unitWord(it) {
    if (it.unit === "seconds") return cues.seconds || t("sec");
    if (it.unit === "breaths") return t("breaths");
    return cues.reps || t("reps");
  }
  function nextDisplay(it) { // on-screen "Name · ×12" / "Name · 30″"
    var v = it.unit === "seconds" ? it.value + "″"
      : it.unit === "breaths" ? it.value + " " + t("breaths")
      : "× " + it.value;
    return it.name + " · " + v;
  }

  // doRest shows the rest/get-ready overlay and counts down. The big number
  // counts the seconds and the final 3 are spoken (3-2-1, exactly on the tick).
  // For a real rest it then speaks "Rest" + the next exercise with its count —
  // after a short beat so it never lands on top of the previous set's "one".
  // The get-ready before an exercise is just the spoken 3-2-1, leading into "Go".
  function doRest(seconds, upItem, isReady, onDone) {
    clearNarr();
    el.stage.classList.add("hidden");
    el.controls.classList.add("hidden");
    el.done.classList.add("hidden"); // Done lives on its own row, outside .controls
    el.rest.classList.remove("hidden");
    el.restTitle.textContent = isReady ? t("ready") : t("rest");
    el.restNext.textContent = upItem ? t("next") + ": " + nextDisplay(upItem) : "";
    pendingStart = onDone;
    if (upItem && !isReady && seconds >= 6) {
      announceTimer = setTimeout(function () {
        announceTimer = null;
        if (paused || el.rest.classList.contains("hidden")) return;
        shutUp();
        if (cues.rest) utter(cues.rest); // "Rest" (queued)
        utter((cues.next_ex || t("next")) + ": " + upItem.name + ", " + upItem.value + " " + unitWord(upItem));
      }, 700);
    }
    countdown(seconds, function (r) {
      el.restCount.textContent = r;
      if (r <= 3 && r >= 1) sayNow(String(r)); // digits land exactly on the second
    }, function () {
      var f = pendingStart; pendingStart = null;
      if (f) f();
    });
  }

  // prev returns to the previous exercise (or restarts the current one when
  // resting), in case the user skipped or advanced by accident.
  function prev() {
    shutUp();
    clearTicker();
    paused = false;
    el.pause.textContent = t("pause");
    var target = vis(el.rest) ? i : Math.max(0, i - 1);
    doRest(3, items[target], true, function () { startItem(target); });
  }

  function endItem() {
    clearTicker();
    clearNarr();
    if (i >= items.length - 1) { finish(); return; }
    var rest = items[i].rest || 0;
    var next = i + 1;
    if (rest <= 0) { startItem(next); return; }
    doRest(rest, items[next], false, function () { startItem(next); });
  }

  // finish ends the workout: celebrate and report completion to the server.
  async function finish() {
    if (finished) return;
    finished = true;
    clearTicker();
    clearNarr();
    shutUp();
    releaseWake();
    el.stage.classList.add("hidden");
    el.controls.classList.add("hidden");
    el.done.classList.add("hidden");
    el.rest.classList.add("hidden");
    el.doneOverlay.classList.remove("hidden");

    var name = (payload.name || "").trim();
    el.doneTitle.textContent = name ? name + ", " + t("bravo") : t("done_title");
    var message = enc.length ? enc[Math.floor(Math.random() * enc.length)] : "";
    el.doneEnc.textContent = message;

    celebrate();
    speak(name ? name + ", " + cues.well_done : cues.workout_end);
    if (message) setTimeout(function () { speak(message); }, 1100);

    try {
      var streak = await $hook("complete", { day: payload.day });
      if (typeof streak === "number") {
        el.doneStreak.textContent = "🔥 " + streak + " · " + t("streak");
      }
    } catch (e) {}
  }

  // ---- celebration (applause + confetti) --------------------------------
  function celebrate() {
    confetti();
    try {
      var Ctx = window.AudioContext || window.webkitAudioContext;
      if (!Ctx) return;
      var ac = new Ctx();
      // Short bright chime arpeggio.
      [523, 659, 784, 1046].forEach(function (f, n) {
        var o = ac.createOscillator(), g = ac.createGain();
        o.type = "triangle"; o.frequency.value = f;
        o.connect(g); g.connect(ac.destination);
        var s = ac.currentTime + n * 0.12;
        g.gain.setValueAtTime(0.0001, s);
        g.gain.exponentialRampToValueAtTime(0.25, s + 0.03);
        g.gain.exponentialRampToValueAtTime(0.0001, s + 0.4);
        o.start(s); o.stop(s + 0.42);
      });
      // Applause-ish noise burst.
      var buffer = ac.createBuffer(1, ac.sampleRate * 1.2, ac.sampleRate);
      var d = buffer.getChannelData(0);
      for (var k = 0; k < d.length; k++) {
        var env = 1 - k / d.length;
        d[k] = (Math.random() * 2 - 1) * env * env * 0.5;
      }
      var src = ac.createBufferSource(); src.buffer = buffer;
      var bp = ac.createBiquadFilter(); bp.type = "bandpass"; bp.frequency.value = 1800;
      src.connect(bp); bp.connect(ac.destination);
      src.start(ac.currentTime + 0.1);
    } catch (e) {}
  }

  function confetti() {
    var colors = ["#fbbf24", "#f472b6", "#34d399", "#60a5fa", "#f87171"];
    for (var n = 0; n < 60; n++) {
      var c = document.createElement("div");
      c.className = "confetti-piece";
      c.style.left = Math.random() * 100 + "%";
      c.style.background = colors[n % colors.length];
      c.style.animationDelay = (Math.random() * 0.6) + "s";
      c.style.transform = "rotate(" + (Math.random() * 360) + "deg)";
      el.doneOverlay.appendChild(c);
    }
  }

  // ---- controls ---------------------------------------------------------
  el.done.addEventListener("click", function () {
    var it = items[i];
    if (it && it.unit !== "seconds") { nextSide(it, curSide); }
  });
  el.skip.addEventListener("click", function () { shutUp(); clearNarr(); endItem(); });
  el.prev.addEventListener("click", prev);
  if (el.restBack) el.restBack.addEventListener("click", prev);
  // Skip jumps straight to whatever the rest/ready countdown was about to start.
  el.restSkip.addEventListener("click", function () {
    clearTicker();
    clearNarr();
    shutUp();
    var f = pendingStart; pendingStart = null;
    if (f) f();
  });
  el.pause.addEventListener("click", function () {
    paused = !paused;
    el.pause.textContent = paused ? t("resume") : t("pause");
    root.classList.toggle("frozen", paused); // freeze the screen; controls stay live
    if (paused) shutUp();
  });

  // Exercise detail card: description, how-to, correct, mistakes, breathing,
  // warning, and a manual replace button. Opening it pauses the timer (reading
  // is a pull action, not part of the workout flow).
  var infoWasPaused = false;
  function setSec(sec, p, text) {
    if (!sec) return;
    if (text) { p.textContent = text; sec.classList.remove("hidden"); }
    else { sec.classList.add("hidden"); }
  }
  function fillList(list, arr) {
    if (!list) return;
    var sec = list.parentElement;
    list.innerHTML = "";
    if (!arr || !arr.length) { sec.classList.add("hidden"); return; }
    sec.classList.remove("hidden");
    for (var k = 0; k < arr.length; k++) {
      var li = document.createElement("li");
      li.textContent = arr[k];
      list.appendChild(li);
    }
  }
  function openInfo() {
    var it = items[i];
    if (!it) return;
    infoWasPaused = paused;
    paused = true;
    shutUp();
    el.infoMedia.src = it.video || it.svg;
    el.infoName.textContent = it.name;
    el.infoDesc.textContent = it.desc || "";
    fillList(el.infoHow, it.howTo);
    fillList(el.infoCorrect, it.correct);
    fillList(el.infoWrong, it.wrong);
    setSec(el.infoBreathingSec, el.infoBreathing, it.breathing);
    setSec(el.infoWarnSec, el.infoWarn, it.warning || it.safety);
    if (el.infoReplace) {
      if (it.repl) el.infoReplace.classList.remove("hidden");
      else el.infoReplace.classList.add("hidden");
    }
    el.infoOverlay.classList.remove("hidden");
  }
  function closeInfo() {
    el.infoOverlay.classList.add("hidden");
    paused = infoWasPaused;
    root.classList.toggle("frozen", paused);
  }
  // doReplace swaps the current exercise for its replacement and restarts it.
  function doReplace() {
    var it = items[i];
    if (!it || !it.repl) return;
    var r = it.repl;
    items[i] = {
      id: r.id, unit: r.unit, value: r.value, slot: it.slot, perSide: r.perSide,
      round: it.round, rest: it.rest, name: r.name, hint: r.hint, warning: r.warning,
      svg: r.svg, video: "", desc: "", howTo: [], correct: [], wrong: [],
      breathing: "", safety: "", vStart: cues.go, vMid: [], vLast: cues.five || "",
      vFinish: cues.well_done, repl: null
    };
    el.infoOverlay.classList.add("hidden");
    paused = false;
    shutUp();
    startItem(i);
  }
  if (el.infoBtn) el.infoBtn.addEventListener("click", openInfo);
  if (el.svg) el.svg.addEventListener("click", openInfo); // tap the figure → pause + show details
  if (el.infoClose) el.infoClose.addEventListener("click", closeInfo);
  if (el.infoReplace) el.infoReplace.addEventListener("click", doReplace);

  // Static control labels
  el.pause.textContent = t("pause");
  el.skip.textContent = t("skip");
  el.done.textContent = t("done");
  el.restSkip.textContent = t("skip");
  if (el.restBack) el.restBack.textContent = t("prev");

  // ---- go ---------------------------------------------------------------
  $sys.clean(function () { clearTicker(); clearNarr(); shutUp(); releaseWake(); });

  requestWake();
  if (items.length === 0) {
    finish();
  } else {
    // A short "get ready" 3-2-1 leads into the first exercise.
    doRest(3, items[0], true, function () { startItem(0); });
  }
})();
