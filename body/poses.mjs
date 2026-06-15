// body/poses.mjs — one pose per exercise ID for the headless figure renderer (shot.html).
// Pose = joint direction vectors (our space: x=front, y=up, z=left/+L). Omitted joints fall
// back to the flat-foot neutral. `view`: front | side | iso. Feet default to flat rig-rest
// unless overridden (e.g. plank on toes). perSide exercises show one side.
// Authored from exercises.go names + details HowTo + the spec's 18-pose table.

const DN=[0,-1,0];

// ---------- reusable archetypes ----------
// straight body tilted ~22° (shoulders high on vertical arms, feet low), belly down (back up),
// palms flat forward on the floor, toes down — hands and toes both reach the ground.
const PUSHUP_TOP = {spine:[-0.93,0.37,0],front:[-0.37,-0.93,0],head:[-0.88,0.3,0],
  Luarm:[-0.05,-1,0],Lfarm:[-0.05,-1,0],Lhand:[-1,0,0],
  Ruarm:[-0.05,-1,0],Rfarm:[-0.05,-1,0],Rhand:[-1,0,0],
  Lthigh:[0.93,-0.37,0],Lshin:[0.93,-0.37,0],Lfoot:[0.4,-0.85,0],
  Rthigh:[0.93,-0.37,0],Rshin:[0.93,-0.37,0],Rfoot:[0.4,-0.85,0]};
const FOREARM_PLANK = {...PUSHUP_TOP, Lfarm:[-0.85,-0.5,0], Rfarm:[-0.85,-0.5,0]};
const SQUAT = {spine:[0.42,0.91,0],front:[1,0,0],
  Luarm:[0.6,0,0],Lfarm:[0.92,0,0],Ruarm:[0.6,0,0],Rfarm:[0.92,0,0],
  Lthigh:[0.95,-0.12,0],Lshin:[-0.6,-0.8,0],Rthigh:[0.95,-0.12,0],Rshin:[-0.6,-0.8,0]};
const HINGE = {spine:[0.75,0.66,0],front:[1,0,0],
  Luarm:[0.2,-0.97,0],Lfarm:[0.2,-0.97,0],Ruarm:[0.2,-0.97,0],Rfarm:[0.2,-0.97,0],
  Lthigh:[0.1,-0.99,0],Lshin:[0,-1,0],Rthigh:[0.1,-0.99,0],Rshin:[0,-1,0]};
// supine = lying on back (head toward -x, face up); legs/arms set per move
const SUPINE = {spine:[-1,0.05,0],front:[0,1,0],
  Luarm:[-0.9,0.2,0],Lfarm:[-0.9,0.2,0],Ruarm:[-0.9,0.2,0],Rfarm:[-0.9,0.2,0],
  Lthigh:[-1,0,0],Lshin:[-1,0,0],Lfoot:[-0.7,-0.7,0],Rthigh:[-1,0,0],Rshin:[-1,0,0],Rfoot:[-0.7,-0.7,0]};
// prone = lying face DOWN along the body axis; head toward +x, belly down (-y), legs extend
// the OPPOSITE way (toward -x) so the body is a flat line, not folded into a crouch.
const PRONE = {spine:[0.96,0.28,0],front:[0.27,-0.96,0],
  Luarm:[0.7,-0.6,0.3],Lfarm:[0.9,0,0.2],Ruarm:[0.7,-0.6,-0.3],Rfarm:[0.9,0,-0.2],
  Lthigh:[-0.96,-0.26,0],Lshin:[-0.98,-0.12,0],Lfoot:[-0.7,-0.6,0],
  Rthigh:[-0.96,-0.26,0],Rshin:[-0.98,-0.12,0],Rfoot:[-0.7,-0.6,0]};
const STAND = {view:'front'};   // plain neutral

export const POSES = {
  _neutral: STAND,

  // ---- Warm-up ----
  W01: {view:'front', pose:{Luarm:[0,-0.95,0.1],Lfarm:[0,-0.95,0.1],Ruarm:[0,-0.95,-0.1],Rfarm:[0,-0.95,-0.1]}}, // shoulder circles
  W02: {view:'front', pose:{Luarm:[0,0.1,1],Lfarm:[0,0.1,1],Ruarm:[0,0.1,-1],Rfarm:[0,0.1,-1]}},                 // arm circles (arms out)
  W03: {view:'front', pose:{front:[0.7,0,0.7],Luarm:[0.3,0,0.6],Lfarm:[0.6,0,0.4],Ruarm:[0.3,0,-0.6],Rfarm:[0.6,0,-0.4]}}, // standing twist
  W04: {view:'side', pose:HINGE},                                                                                  // hip hinge
  W05: {view:'side', pose:{...SQUAT, Lthigh:[0.4,-0.9,0],Rthigh:[0.4,-0.9,0],Lshin:[-0.02,-1,0],Rshin:[-0.02,-1,0]}}, // shallow squat
  W06: {view:'side', pose:PUSHUP_TOP},                                                                            // high plank
  W07: {view:'front', pose:{Luarm:[0,1,0.1],Lfarm:[0,1,0.1],Ruarm:[0,1,-0.1],Rfarm:[0,1,-0.1]}},                 // breathing + reach (overhead)
  W08: {view:'front', pose:{front:[0.7,0,0.7],Luarm:[0.2,0,0.7],Lfarm:[0.5,0,0.5],Ruarm:[0.2,0,-0.7],Rfarm:[0.5,0,-0.5]}}, // torso rotations
  W09: {view:'side', pose:{Luarm:[1,0.1,0],Lfarm:[1,0.1,0],Ruarm:[-0.6,-0.5,0],Rfarm:[-0.6,-0.5,0]}},            // arm swings (one fwd/one back)

  // ---- Abs / core ----
  C01: {view:'side', pose:FOREARM_PLANK},                                                                         // forearm plank
  C02: {view:'side', pose:FOREARM_PLANK},                                                                         // RKC plank
  C03: {view:'side', pose:{...SUPINE, Lthigh:[0.45,0.89,0],Lshin:[0.45,0.89,0],Lfoot:[0.45,0.89,0],Rthigh:[0.66,0.75,0],Rshin:[0.66,0.75,0],Rfoot:[0.66,0.75,0]}}, // scissors
  C04: {view:'side', pose:{...SUPINE, Lthigh:[0.35,0.93,0],Lshin:[0.35,0.93,0],Rthigh:[0.7,0.7,0],Rshin:[0.7,0.7,0]}}, // flutter kicks
  C05: {view:'side', pose:{...SUPINE, Lthigh:[0.2,0.6,0],Lshin:[-0.7,0.4,0],Rthigh:[0.2,0.6,0],Rshin:[-0.7,0.4,0]}}, // reverse crunch (knees to chest)
  C06: {view:'side', pose:{...SUPINE, Luarm:[-1,0.1,0],Lfarm:[-1,0.1,0],Ruarm:[-1,0.1,0],Rfarm:[-1,0.1,0], Lthigh:[0.5,0.5,0],Lshin:[0.5,0.5,0],Rthigh:[0.5,0.5,0],Rshin:[0.5,0.5,0]}}, // hollow hold
  C07: {view:'iso', pose:{spine:[-0.95,0.3,0],front:[-0.2,0,-0.98], Ruarm:[0,-1,0],Rfarm:[0,-1,0], Luarm:[0,1,0],Lfarm:[0,1,0], Lthigh:[0.95,-0.05,0],Lshin:[0.95,-0.05,0],Lfoot:[0.6,-0.7,0],Rthigh:[0.95,-0.05,0],Rshin:[0.95,-0.05,0],Rfoot:[0.6,-0.7,0]}}, // side plank (on side, bottom arm support, top arm up)
  C08: {view:'iso', pose:{spine:[-0.95,0.3,0],front:[-0.2,0,-0.98], Ruarm:[0,-1,0],Rfarm:[0,-1,0], Luarm:[0,1,0],Lfarm:[0,1,0], Lthigh:[0.95,-0.05,0],Lshin:[0.95,-0.05,0],Lfoot:[0.6,-0.7,0],Rthigh:[0.95,-0.05,0],Rshin:[0.95,-0.05,0],Rfoot:[0.6,-0.7,0]}}, // side plank hip lift
  C09: {view:'side', pose:{...PUSHUP_TOP, Ruarm:[0.3,-0.5,0],Rfarm:[1,0.1,0]}},                                  // plank shoulder taps (one hand up)
  C10: {view:'side', pose:PUSHUP_TOP},                                                                            // up-down plank
  C11: {view:'side', pose:{...SUPINE, Lthigh:[0.3,0.5,0],Lshin:[-0.5,0.3,0],Rthigh:[0.3,0.5,0],Rshin:[-0.5,0.3,0], Luarm:[-0.4,-0.5,0],Lfarm:[0.3,-0.7,0],Ruarm:[-0.4,-0.5,0],Rfarm:[0.3,-0.7,0]}}, // heel taps
  C12: {view:'side', pose:{...SUPINE, Luarm:[0,0.9,0],Lfarm:[0,0.9,0], Ruarm:[0,0.9,0],Rfarm:[0,0.9,0], Lthigh:[0.3,0.6,0],Lshin:[-0.6,0.4,0], Rthigh:[0.3,0.6,0],Rshin:[-0.6,0.4,0]}}, // dead bug
  C13: {view:'side', pose:{...SUPINE, Luarm:[-0.3,0.6,0],Lfarm:[0.4,0.5,0], Ruarm:[-0.3,0.6,0],Rfarm:[0.4,0.5,0], Lthigh:[0.3,0.55,0],Lshin:[-0.6,0.3,0], Rthigh:[0.7,0.4,0],Rshin:[0.85,0.2,0]}}, // bicycle crunch

  // ---- Push-ups / shoulders ----
  P01: {view:'side', pose:PUSHUP_TOP}, P02: {view:'side', pose:PUSHUP_TOP}, P03: {view:'side', pose:PUSHUP_TOP},
  P04: {view:'side', pose:PUSHUP_TOP}, P05: {view:'side', pose:{...PUSHUP_TOP, Luarm:[0.3,-0.95,0]}},
  P06: {view:'side', pose:{spine:[-0.5,0.7,0],front:[-0.7,-0.7,0], Luarm:[-0.4,-0.9,0],Lfarm:[-0.4,-0.9,0],Ruarm:[-0.4,-0.9,0],Rfarm:[-0.4,-0.9,0], Lthigh:[0.85,-0.5,0],Lshin:[0.85,-0.5,0],Lfoot:[0.5,-0.86,0],Rthigh:[0.85,-0.5,0],Rshin:[0.85,-0.5,0],Rfoot:[0.5,-0.86,0]}}, // pike push-up
  P08: {view:'side', pose:{...PRONE, Luarm:[0.2,-0.95,0],Lfarm:[0.6,-0.3,0],Ruarm:[0.2,-0.95,0],Rfarm:[0.6,-0.3,0], spine:[0.7,0.66,0],front:[0.2,-0.97,0]}}, // sphinx push-up

  // ---- Back / scapula (prone) ----
  B01: {view:'side', pose:{...PRONE, spine:[0.8,0.55,0],front:[0.4,-0.9,0], Luarm:[0.6,-0.4,0.4],Lfarm:[0.8,-0.2,0.3],Ruarm:[0.6,-0.4,-0.4],Rfarm:[0.8,-0.2,-0.3]}}, // cobra hold
  B02: {view:'iso', pose:{...PRONE, Luarm:[0.2,0.1,0.95],Lfarm:[0.2,0.1,0.95],Ruarm:[0.2,0.1,-0.95],Rfarm:[0.2,0.1,-0.95]}}, // reverse snow angels
  B03: {view:'iso', pose:{...PRONE, Luarm:[0.3,0.5,0.6],Lfarm:[0.1,0.9,0.3],Ruarm:[0.3,0.5,-0.6],Rfarm:[0.1,0.9,-0.3]}}, // W-raise
  B04: {view:'iso', pose:{...PRONE, Luarm:[0.3,0.8,0.4],Lfarm:[0.3,0.9,0.3],Ruarm:[0.3,0.8,-0.4],Rfarm:[0.3,0.9,-0.3]}}, // Y-raise
  B05: {view:'iso', pose:{...PRONE, Luarm:[0.1,0.2,1],Lfarm:[0.1,0.2,1],Ruarm:[0.1,0.2,-1],Rfarm:[0.1,0.2,-1]}}, // T-raise
  B06: {view:'iso', pose:{...PRONE, Luarm:[0.3,0.7,0.5],Lfarm:[0.2,0.4,0.7],Ruarm:[0.3,0.7,-0.5],Rfarm:[0.2,0.4,-0.7]}}, // superman pull-down
  B07: {view:'iso', pose:{...PRONE, Luarm:[0.6,0.5,0.25],Lfarm:[0.5,0.7,0.2],Ruarm:[0.85,-0.2,-0.3],Rfarm:[0.95,-0.1,-0.2], Lthigh:[-0.92,0.32,0],Lshin:[-0.94,0.3,0], Rthigh:[-0.97,-0.18,0],Rshin:[-0.98,-0.12,0]}}, // swimmers (opp arm+leg up)
  B08: {view:'side', pose:{spine:[0.55,0.83,0],front:[-0.5,0.86,0], Luarm:[-0.55,-0.83,0],Lfarm:[-0.5,-0.86,0],Ruarm:[-0.55,-0.83,0],Rfarm:[-0.5,-0.86,0], Lthigh:[0.93,-0.36,0],Lshin:[0.96,-0.28,0],Lfoot:[0.6,-0.8,0],Rthigh:[0.93,-0.36,0],Rshin:[0.96,-0.28,0],Rfoot:[0.6,-0.8,0]}}, // reverse plank: straight body line, hands behind, hips lifted
  B09: {view:'side', pose:{...PRONE, spine:[0.85,0.5,0],front:[0.3,-0.95,0], Luarm:[0.3,0.1,0.5],Lfarm:[0.5,0.1,0.4],Ruarm:[0.3,0.1,-0.5],Rfarm:[0.5,0.1,-0.4]}}, // back extension pulses
  B10: {view:'iso', pose:{...PRONE, Luarm:[0.3,0.5,0.6],Lfarm:[0.1,0.85,0.4],Ruarm:[0.3,0.5,-0.6],Rfarm:[0.1,0.85,-0.4]}}, // cobra to W
  B11: {view:'iso', pose:{...PRONE, Luarm:[0.6,0.5,0.25],Lfarm:[0.5,0.7,0.2],Ruarm:[0.85,-0.2,-0.3],Rfarm:[0.95,-0.1,-0.2], Lthigh:[-0.92,0.32,0],Lshin:[-0.94,0.3,0], Rthigh:[-0.97,-0.18,0],Rshin:[-0.98,-0.12,0]}}, // swimmers per side

  // ---- Arms ----
  A01: {view:'front', pose:{Luarm:[0.25,-0.9,0],Lfarm:[0.5,0.75,0], Ruarm:DN,Rfarm:DN}}, // self-resist curl (one arm curled up)
  A02: {view:'front', pose:{Luarm:[0.1,-0.95,0],Lfarm:[1,0,0], Ruarm:DN,Rfarm:DN}},     // biceps isometric 90
  A03: {view:'front', pose:{Luarm:[0.2,-0.6,0],Lfarm:[0.7,0.4,0.2], Ruarm:[0.2,-0.6,0],Rfarm:[0.7,0.4,-0.2]}}, // hand-hook pull (hands at chest)
  A04: {view:'front', pose:{Luarm:[0.2,-0.9,0],Lfarm:[0.85,-0.5,0], Ruarm:DN,Rfarm:DN}}, // slow negative curl (mid)

  // ---- Legs / glutes ----
  L01: {view:'side', pose:SQUAT}, L02: {view:'side', pose:SQUAT}, L03: {view:'side', pose:HINGE},
  L04: {view:'side', pose:{...SUPINE, spine:[-0.8,0.6,0], Lthigh:[0.7,0.4,0],Lshin:[-0.3,-0.9,0],Lfoot:[0.7,-0.6,0], Rthigh:[0.85,0.4,0],Rshin:[0.85,0.4,0], Luarm:[-0.8,-0.5,0],Lfarm:[-0.8,-0.5,0],Ruarm:[-0.8,-0.5,0],Rfarm:[-0.8,-0.5,0]}}, // single-leg glute bridge
  L05: {view:'side', pose:{...SUPINE, spine:[-0.8,0.6,0], Lthigh:[0.6,0.5,0],Lshin:[-0.4,-0.85,0], Rthigh:[0.85,0.45,0],Rshin:[0.85,0.45,0], Luarm:[-0.8,-0.5,0],Lfarm:[-0.8,-0.5,0],Ruarm:[-0.8,-0.5,0],Rfarm:[-0.8,-0.5,0]}}, // bridge march
  L06: {view:'side', pose:{...SUPINE, spine:[-0.8,0.6,0], Lthigh:[0.75,0.45,0],Lshin:[-0.35,-0.9,0],Lfoot:[0.8,-0.5,0], Rthigh:[0.75,0.45,0],Rshin:[-0.35,-0.9,0],Rfoot:[0.8,-0.5,0], Luarm:[-0.8,-0.5,0],Lfarm:[-0.8,-0.5,0],Ruarm:[-0.8,-0.5,0],Rfarm:[-0.8,-0.5,0]}}, // bridge hold
  L07: {view:'side', pose:{Lfoot:[0.85,-0.5,0],Rfoot:[0.85,-0.5,0], spine:[0,1,0]}},   // calf raises (on toes)
  L09: {view:'iso', pose:{spine:[-0.95,0.3,0],front:[-0.2,0,-0.98], Ruarm:[0,-1,0],Rfarm:[0,-1,0], Luarm:[0,1,0],Lfarm:[0,1,0], Rthigh:[0.95,-0.05,0],Rshin:[0.95,-0.05,0],Rfoot:[0.6,-0.7,0], Lthigh:[0.9,0.35,0],Lshin:[0.9,0.35,0],Lfoot:[0.7,-0.3,0]}}, // side plank + top leg lift
  L10: {view:'side', pose:SQUAT}, L12: {view:'side', pose:SQUAT},
  L11: {view:'side', pose:{spine:[0.05,1,0],front:[1,0,0], Luarm:[0.3,-0.9,0],Lfarm:[0.3,-0.9,0],Ruarm:[0.3,-0.9,0],Rfarm:[0.3,-0.9,0], Lthigh:[0.5,-0.86,0],Lshin:[0.4,-0.9,0],Rthigh:[-0.6,-0.8,0],Rshin:[0.2,-0.98,0],Rfoot:[0.6,-0.8,0]}}, // reverse lunge

  // ---- Cardio / plyometrics ----
  J01: {view:'front', pose:{Luarm:[0.2,-0.9,0],Lfarm:[0.2,-0.9,0],Ruarm:[0.2,-0.9,0],Rfarm:[0.2,-0.9,0], Lfoot:[0.85,-0.5,0],Rfoot:[0.85,-0.5,0]}}, // pogo (on toes, feet together)
  J02: {view:'side', pose:{Luarm:[0.55,-0.55,0],Lfarm:[0.95,0.1,0],Ruarm:[-0.5,-0.7,0],Rfarm:[-0.65,-0.55,0], Lthigh:[0.55,-0.5,0],Lshin:[-0.05,-1,0],Lfoot:[0.6,-0.5,0], Rthigh:[0,-1,0],Rshin:[0,-1,0]}}, // marching high knees (thigh fwd-down, shin hangs)
  J03: {view:'front', pose:{Luarm:[0,1,0.55],Lfarm:[0,1,0.55],Ruarm:[0,1,-0.55],Rfarm:[0,1,-0.55], Lthigh:[0,-1,0.45],Lshin:[0,-1,0.45],Rthigh:[0,-1,-0.45],Rshin:[0,-1,-0.45]}}, // jumping jacks
  J04: {view:'side', pose:{...PUSHUP_TOP, Lthigh:[0.4,-0.4,0],Lshin:[-0.6,-0.5,0],Lfoot:[0.5,-0.5,0]}}, // mountain climbers (one knee in)
  J05: {view:'side', pose:{spine:[0.2,0.97,0],front:[1,0,0], Luarm:[-0.5,-0.7,0],Lfarm:[-0.6,-0.6,0],Ruarm:[-0.5,-0.7,0],Rfarm:[-0.6,-0.6,0], Lthigh:[0.3,-0.95,0],Lshin:[-0.05,-1,0],Lfoot:[0.9,-0.3,0],Rthigh:[0.3,-0.95,0],Rshin:[-0.05,-1,0],Rfoot:[0.9,-0.3,0]}}, // squat jump (takeoff)
  J06: {view:'front', pose:{spine:[0.1,0.99,0.1],front:[1,0,0], Luarm:[0.3,0.2,0.7],Lfarm:[0.3,0.2,0.7],Ruarm:[0.2,-0.6,-0.5],Rfarm:[0.2,-0.6,-0.5], Lthigh:[0.3,-0.9,0.2],Lshin:[0,-1,0], Rthigh:[-0.2,-0.5,-0.7],Rshin:[-0.4,-0.7,-0.5]}}, // skater hop (land one leg)
  J07: {view:'side', pose:{Luarm:[0.6,-0.5,0],Lfarm:[0.95,0.2,0],Ruarm:[-0.6,-0.6,0],Rfarm:[-0.7,-0.5,0], Lthigh:[0.72,-0.32,0],Lshin:[-0.1,-0.99,0],Lfoot:[0.6,-0.5,0], Rthigh:[0,-1,0],Rshin:[0,-1,0],Rfoot:[0.85,-0.5,0]}}, // high knees (thigh up but <horizontal, shin down)
  J08: {view:'side', pose:PUSHUP_TOP},  // burpee (plank phase)
  J09: {view:'side', pose:{spine:[0.1,0.99,0],front:[1,0,0], Luarm:[0.4,-0.85,0],Lfarm:[0.4,-0.85,0],Ruarm:[-0.4,-0.85,0],Rfarm:[-0.4,-0.85,0], Lthigh:[0.55,-0.83,0],Lshin:[0.45,-0.89,0], Rthigh:[-0.6,-0.8,0],Rshin:[0.1,-0.99,0],Rfoot:[0.7,-0.7,0]}}, // jump lunge (air)

  // ---- Cool-down ----
  CD01: {view:'iso', pose:{Luarm:[-0.55,-0.2,0.45],Lfarm:[-0.7,-0.1,0.35],Ruarm:[-0.55,-0.2,-0.45],Rfarm:[-0.7,-0.1,-0.35]}}, // chest stretch (arms back, hands clasped behind)
  CD02: {view:'side', pose:{...SUPINE, Lthigh:[0.2,0.95,0],Lshin:[0.2,0.95,0],Lfoot:[0.2,0.95,0], Rthigh:[1,-0.05,0],Rshin:[1,-0.05,0],Rfoot:[0.7,-0.6,0], Luarm:[0.2,0.9,0],Lfarm:[0.3,0.85,0],Ruarm:[0.2,0.9,0],Rfarm:[0.3,0.85,0]}}, // hamstring stretch (one leg up, other extended)
  CD03: {view:'iso', pose:{...SUPINE, front:[0.4,0.55,0.7], Lthigh:[0.3,0,0.92],Lshin:[0.7,-0.2,0.65], Rthigh:[1,-0.05,0],Rshin:[1,-0.05,0], Luarm:[0,0.1,1],Lfarm:[0,0.1,1],Ruarm:[0,0.1,-1],Rfarm:[0,0.1,-1]}}, // supine twist (one knee across)
  CD04: {view:'side', pose:{...PRONE, spine:[0.75,0.66,0],front:[0.2,-0.97,0], Luarm:[0.3,-0.9,0],Lfarm:[0.6,-0.2,0],Ruarm:[0.3,-0.9,0],Rfarm:[0.6,-0.2,0]}}, // sphinx pose
  CD05: {view:'side', pose:PRONE},  // lying breathing (prone relaxed)
  CD07: {view:'side', pose:{spine:[0.93,0.33,0],front:[0.28,-0.96,0], Luarm:[0.8,-0.25,0],Lfarm:[0.96,-0.12,0],Ruarm:[0.8,-0.25,0],Rfarm:[0.96,-0.12,0], Lthigh:[0.18,-0.98,0],Lshin:[-0.92,-0.32,0],Lfoot:[-0.85,-0.2,0], Rthigh:[0.18,-0.98,0],Rshin:[-0.92,-0.32,0],Rfoot:[-0.85,-0.2,0]}}, // child's pose: kneel, fold torso forward-down, arms reach forward
};
