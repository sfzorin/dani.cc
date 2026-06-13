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
  var voiceOn = !!payload.voice;
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
    infoTech: $(".info-tech"),
    infoTechSec: $(".info-tech-sec"),
    infoMistake: $(".info-mistake"),
    infoMistakeSec: $(".info-mistake-sec"),
    infoWarn: $(".info-warn"),
    infoWarnSec: $(".info-warn-sec"),
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
    feedback: $(".feedback-overlay"),
    fbDiff: $(".fb-diff"),
    fbDiffOut: $(".fb-diff-out"),
    fbSave: $(".fb-save"),
    fbSkip: $(".fb-skip"),
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

  function speak(text) {
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
  function shutUp() {
    if ("speechSynthesis" in window) { try { window.speechSynthesis.cancel(); } catch (e) {} }
  }

  function clearTicker() {
    if (ticker) { clearInterval(ticker); ticker = null; }
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
  var curSide = 1;         // which side of a per-side exercise is running

  function sides(it) { return it.perSide ? 2 : 1; }

  function startItem(idx) {
    i = idx;
    var it = items[i];
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
    setProgress();
    runSide(it, 1);
  }

  // runSide runs one side of an exercise (sides==1 for non-per-side).
  function runSide(it, side) {
    curSide = side;
    el.side.textContent = it.perSide ? t("side") + " " + side + "/" + sides(it) : "";
    if (it.unit === "seconds") {
      el.done.classList.add("hidden");
      el.value.classList.add("timer");
      countdown(it.value, function (r) {
        el.value.textContent = r + "″";
        if (r <= 3 && r >= 1) speak(String(r)); // spoken 3-2-1 before the end
      }, function () { nextSide(it, side); });
    } else {
      // reps / breaths: wait for the user to tap Done
      clearTicker();
      el.value.classList.remove("timer");
      el.value.textContent = it.unit === "breaths"
        ? it.value + " " + t("breaths")
        : "× " + it.value;
      el.done.classList.remove("hidden");
    }
  }

  function nextSide(it, side) {
    if (side < sides(it)) {
      speak(cues.switch_side);
      runSide(it, side + 1); // straight into the other side, no rest
    } else {
      endItem();
    }
  }

  // doRest shows the rest/get-ready overlay and counts down. The final 3 seconds
  // are spoken (3-2-1) with "go" at zero — the countdown is part of the pause,
  // not an extra delay. `onDone` runs at zero (and is what Skip jumps to).
  function doRest(seconds, upcomingName, isReady, onDone) {
    el.stage.classList.add("hidden");
    el.controls.classList.add("hidden");
    el.rest.classList.remove("hidden");
    el.restTitle.textContent = isReady ? t("ready") : t("rest");
    el.restNext.textContent = upcomingName ? t("next") + ": " + upcomingName : "";
    pendingStart = onDone;
    speak(isReady ? cues.ready : cues.rest);
    if (upcomingName) speak(upcomingName);
    countdown(seconds, function (r) {
      el.restCount.textContent = r;
      if (r <= 3 && r >= 1) speak(String(r));
    }, function () {
      speak(cues.go);
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
    doRest(3, items[target].name, true, function () { startItem(target); });
  }

  function endItem() {
    clearTicker();
    if (i >= items.length - 1) { finish(); return; }
    var nextName = items[i + 1].name;
    var rest = items[i].rest || 0;
    var next = i + 1;
    if (rest <= 0) { startItem(next); return; }
    doRest(rest, nextName, false, function () { startItem(next); });
  }

  // finish ends the workout and shows the feedback questionnaire first; the
  // celebration + server save happen on Save/Skip (submitWorkout).
  function finish() {
    if (finished) return;
    finished = true;
    clearTicker();
    shutUp();
    releaseWake();
    el.stage.classList.add("hidden");
    el.controls.classList.add("hidden");
    el.rest.classList.add("hidden");
    if (el.feedback) {
      el.feedback.classList.remove("hidden");
    } else {
      submitWorkout(null);
    }
  }

  function fbVal(key) {
    var a = el.feedback.querySelector('.fb-group[data-key="' + key + '"] .fb-chip.active');
    if (a) return a.getAttribute("data-val");
    return key === "energy" ? "normal" : "0";
  }
  function collectFeedback() {
    return {
      difficulty: parseInt(el.fbDiff.value, 10) || 6,
      knee: parseInt(fbVal("knee"), 10) || 0,
      back: parseInt(fbVal("back"), 10) || 0,
      shoulder: parseInt(fbVal("shoulder"), 10) || 0,
      energy: fbVal("energy")
    };
  }

  // submitWorkout shows the celebration and reports completion (+ optional
  // feedback) to the server.
  async function submitWorkout(fb) {
    if (el.feedback) el.feedback.classList.add("hidden");
    el.doneOverlay.classList.remove("hidden");

    var name = (payload.name || "").trim();
    el.doneTitle.textContent = name ? name + ", " + t("bravo") : t("done_title");
    var message = enc.length ? enc[Math.floor(Math.random() * enc.length)] : "";
    el.doneEnc.textContent = message;

    celebrate();
    speak(name ? name + ", " + cues.well_done : cues.workout_end);
    if (message) setTimeout(function () { speak(message); }, 1100);

    var body = { day: payload.day };
    if (fb) { body.hasFeedback = true; body.feedback = fb; }
    try {
      var streak = await $hook("complete", body);
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
  el.skip.addEventListener("click", function () { shutUp(); endItem(); });
  el.prev.addEventListener("click", prev);
  if (el.restBack) el.restBack.addEventListener("click", prev);
  // Skip jumps straight to whatever the rest/ready countdown was about to start.
  el.restSkip.addEventListener("click", function () {
    clearTicker();
    var f = pendingStart; pendingStart = null;
    if (f) f();
  });
  el.pause.addEventListener("click", function () {
    paused = !paused;
    el.pause.textContent = paused ? t("resume") : t("pause");
    if (paused) shutUp();
  });

  // Feedback questionnaire wiring.
  if (el.fbDiff) {
    el.fbDiff.addEventListener("input", function () { el.fbDiffOut.textContent = el.fbDiff.value; });
  }
  if (el.feedback) {
    el.feedback.addEventListener("click", function (e) {
      var chip = e.target && e.target.closest ? e.target.closest(".fb-chip") : null;
      if (!chip) return;
      var chips = chip.parentElement.querySelectorAll(".fb-chip");
      for (var k = 0; k < chips.length; k++) chips[k].classList.remove("active");
      chip.classList.add("active");
    });
  }
  if (el.fbSave) el.fbSave.addEventListener("click", function () { submitWorkout(collectFeedback()); });
  if (el.fbSkip) el.fbSkip.addEventListener("click", function () { submitWorkout(null); });

  // Exercise detail card (full technique / mistake / safety warning).
  var infoWasPaused = false;
  function setSec(sec, p, text) {
    if (!sec) return;
    if (text) { p.textContent = text; sec.classList.remove("hidden"); }
    else { sec.classList.add("hidden"); }
  }
  function openInfo() {
    var it = items[i];
    if (!it) return;
    infoWasPaused = paused;
    paused = true; // freeze the timer while reading
    el.infoMedia.src = it.video || it.svg;
    el.infoName.textContent = it.name;
    setSec(el.infoTechSec, el.infoTech, it.technique);
    setSec(el.infoMistakeSec, el.infoMistake, it.mistake);
    setSec(el.infoWarnSec, el.infoWarn, it.warning);
    el.infoOverlay.classList.remove("hidden");
  }
  function closeInfo() {
    el.infoOverlay.classList.add("hidden");
    paused = infoWasPaused;
  }
  if (el.infoBtn) el.infoBtn.addEventListener("click", openInfo);
  if (el.infoClose) el.infoClose.addEventListener("click", closeInfo);

  // Static control labels
  el.pause.textContent = t("pause");
  el.skip.textContent = t("skip");
  el.done.textContent = t("done");
  el.restSkip.textContent = t("skip");
  if (el.restBack) el.restBack.textContent = t("prev");

  // ---- go ---------------------------------------------------------------
  $sys.clean(function () { clearTicker(); shutUp(); releaseWake(); });

  requestWake();
  if (items.length === 0) {
    finish();
  } else {
    // A short "get ready" 3-2-1 leads into the first exercise.
    doRest(3, items[0].name, true, function () { startItem(0); });
  }
})();
