import math
TK="M 0,-10 C 20,-10 25,2 25,38 C 25,72 17,86 17,108 C 17,138 23,150 23,168 C 23,182 12,186 0,186 C -12,186 -23,182 -23,168 C -23,150 -17,138 -17,108 C -17,86 -25,72 -25,38 C -25,2 -20,-10 0,-10 Z"
# canonical foot: ankle at origin, toe toward +x, sole +y, rounded toe (not pointed)
FT="M 8,-6 C 34,-6 50,0 57,9 C 63,13 61,21 51,21 L -8,21 C -18,21 -18,-3 -9,-5 Z"
def ang(sx,sy,hx,hy): return round(math.degrees(math.atan2(-(hx-sx),(hy-sy))),1)
def path(polys):
    s=[]
    for p in polys:
        s.append("M "+" L ".join(f"{x},{y}" for x,y in p))
    return " ".join(s)
def trunk(sx,sy,hx,hy,fill,stroke=None):
    a=ang(sx,sy,hx,hy)
    st=f' stroke="{stroke}" stroke-width="6"' if stroke else ''
    return f'<use href="#TK" transform="translate({sx} {sy}) rotate({a})" fill="{fill}"{st}/>'
def foot(ax,ay,deg,fill,stroke=None):
    st=f' stroke="{stroke}" stroke-width="6"' if stroke else ''
    return f'<use href="#FT" transform="translate({ax} {ay}) rotate({round(deg,1)})" fill="{fill}"{st}/>'
def head(cx,cy,deg,fill,stroke=None):
    st=f' stroke="{stroke}" stroke-width="6"' if stroke else ''
    return f'<ellipse cx="{cx}" cy="{cy}" rx="29" ry="41" fill="{fill}"{st} transform="rotate({deg} {cx} {cy})"/>'
def limbs(polys,w,color):
    return f'<path d="{path(polys)}" fill="none" stroke="{color}" stroke-width="{w}"/>'

def body(trunk_a, groups, feet, head_a):
    out=[trunk(*trunk_a,fill="#E5E7EB",stroke="#1F2937")]
    for gi,g in enumerate(groups):
        out.append(limbs(g,40,"#1F2937"))
        for (ax,ay,deg,fgi) in feet:
            if fgi==gi: out.append(foot(ax,ay,deg,"#E5E7EB","#1F2937"))
        out.append(limbs(g,28,"#E5E7EB"))
    out.append(head(*head_a,fill="#E5E7EB",stroke="#1F2937"))
    return "\n      ".join(out)

def ghost(trunk_a, polys, feet, head_c):
    out=[]
    if trunk_a: out.append(trunk(*trunk_a,fill="#9CA3AF"))
    out.append(limbs(polys,28,"#9CA3AF"))
    for (ax,ay,deg) in feet: out.append(foot(ax,ay,deg,"#9CA3AF"))
    if head_c: out.append(f'<ellipse cx="{head_c[0]}" cy="{head_c[1]}" rx="29" ry="41" fill="#9CA3AF" transform="rotate({head_c[2]} {head_c[0]} {head_c[1]})"/>')
    return "\n      ".join(out)

MARKER='<marker id="ah" markerUnits="userSpaceOnUse" markerWidth="26" markerHeight="26" refX="18" refY="13" orient="auto-start-reverse"><path d="M5,3 L21,13 L5,23" fill="none" stroke="#2563EB" stroke-width="5" stroke-linecap="round" stroke-linejoin="round"/></marker>'

def svg(idd,title,desc,mat,shift,bodys,ghosts,arrows,safety,need_marker=True):
    defs=f'{MARKER if need_marker else ""}<path id="TK" d="{TK}"/><path id="FT" d="{FT}"/>'
    arr="".join(f'<path d="{d}" fill="none" stroke="#2563EB" stroke-width="6" stroke-linecap="round" stroke-linejoin="round"{m}/>' for d,m in arrows)
    saf="".join(f'<line x1="{a}" y1="{b}" x2="{c}" y2="{e}"/>' for a,b,c,e in safety)
    g_end=f'<g id="body-end" opacity="0.4" stroke-linecap="round" stroke-linejoin="round">\n      {ghosts}\n    </g>' if ghosts else '<g id="body-end"></g>'
    inner=f'''  {g_end}

    <g id="body-start" stroke-linecap="round" stroke-linejoin="round">
      {bodys}
    </g>

    <g id="movement-arrows">{arr}</g>

    <g id="safety-cues" fill="none" stroke="#10B981" stroke-width="4" stroke-linecap="round" stroke-dasharray="2 14">{saf}</g>'''
    if shift: inner=f'  <g transform="translate({shift} 0)">\n  {inner}\n  </g>'
    return f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 768" role="img" aria-labelledby="{idd}t {idd}d">
  <title id="{idd}t">{title}</title>
  <desc id="{idd}d">{desc}</desc>
  <defs>{defs}</defs>

  <g id="mat">{mat}</g>

{inner}
</svg>
'''

import os
D=os.path.dirname(__file__)+"/svg"
def write(name,s): open(f"{D}/{name}.svg","w").write(s)

MAT_LIE='<rect x="120" y="520" width="780" height="30" rx="15" fill="#D1FAE5" stroke="#1F2937" stroke-width="3" opacity="0.7"/>'
MAT_PL ='<rect x="120" y="560" width="780" height="22" rx="11" fill="#D1FAE5" stroke="#1F2937" stroke-width="3" opacity="0.7"/>'
MAT_FL ='<line x1="180" y1="716" x2="820" y2="716" stroke="#1F2937" stroke-width="4" stroke-linecap="round" opacity="0.45"/>'

# ---- M01 push-up ----
b=body((372,350,539,404),
  [[[(372,350),(372,453),(372,556)],[(539,404),(682,450),(826,497)],[(372,350),(347,342)]]],
  [(826,497,74,0)], (311,330,-72.1))
g=ghost((372,460,546,477),
  [[(372,460),(460,536),(372,556)],[(546,477),(697,492),(826,505)],[(372,460),(346,457)]],
  [(826,505,74)], (308,454,-84.4))
write("M01_push_up", svg("m01","Отжимания классические","Схема упражнения: верхняя и нижняя позиция отжимания со стрелкой движения и линией прямого корпуса.",
  MAT_PL,-110,b,g,[("M 480 378 L 480 484",' marker-start="url(#ah)" marker-end="url(#ah)"')],[(300,305,852,512)]))

# ---- C02 plank ----
b=body((424,440,597,462),
  [[[(424,440),(420,556),(330,556)],[(597,462),(746,481),(897,500)],[(424,440),(398,437)]]],
  [(897,500,74,0)], (360,432,-82.7))
write("C02_forearm_plank", svg("c02","Планка на локтях","Схема упражнения: удержание планки на предплечьях, предплечья вперёд, тело прямой линией от плеч до пяток.",
  MAT_PL,-106,b,"",[],[(360,434,915,524)],need_marker=False))

# ---- C03 flutter ----
b=body((270,500,445,506),
  [ [[(270,500),(378,515),(458,524)],[(270,500),(244,499)]],
    [[(445,506),(548,395),(632,268)],[(445,506),(565,415),(690,330)]] ],
  [(632,268,-56.5,1),(690,330,-34.2,1)], (206,498,-88))
g=ghost(None,[[(445,506),(575,420),(700,345)],[(445,506),(545,400),(628,272)]],
  [(700,345,-31),(628,272,-57)], None)
write("C03_flutter_kicks", svg("c03","Flutter kicks","Схема упражнения: лёжа на спине, ноги приподняты и попеременно двигаются вверх-вниз, поясница прижата к коврику.",
  MAT_LIE,0,b,g,[("M 590 198 L 590 274",' marker-start="url(#ah)" marker-end="url(#ah)"'),("M 792 282 L 792 358",' marker-start="url(#ah)" marker-end="url(#ah)"')],[(360,520,445,522)]))

# ---- C04 dead bug ----
b=body((316,500,491,506),
  [[[(316,500),(208,468),(120,452)],[(316,500),(326,384),(334,296)],[(491,506),(636,468),(784,452)],[(491,506),(534,362),(684,384)],[(316,500),(290,499)]]],
  [(784,452,-6.2,0),(684,384,8.3,0)], (252,498,-88))
write("C04_dead_bug", svg("c04","Dead bug","Схема упражнения: лёжа на спине, противоположные рука и нога вытягиваются, поясница прижата к коврику.",
  MAT_LIE,34,b,"",[("M 150 442 L 96 420",' marker-end="url(#ah)"'),("M 840 452 L 892 478",' marker-end="url(#ah)"')],[(425,514,491,516)]))

# ---- L01 mini squat ----
b=body((470,255,415,420),
  [[[(470,255),(503,365),(592,356)],[(415,420),(502,545),(468,692)],[(470,255),(478,230)]]],
  [(468,692,0,0)], (490,194,18.4))
g=ghost((470,215,470,390),
  [[(470,215),(478,329),(486,417)],[(470,390),(470,541),(468,692)],[(470,215),(470,189)]],
  [(468,692,0)], (470,151,0))
write("L01_mini_squat_45", svg("l01","Мини-присед до 45°","Схема упражнения: неглубокий присед, таз назад-вниз, спина прямая, колени по линии носков.",
  MAT_FL,0,b,g,[("M 405 418 C 360 438, 345 470, 342 505",' marker-end="url(#ah)"')],[(350,425,515,425),(531,716,531,520)]))
print("generated 5")
