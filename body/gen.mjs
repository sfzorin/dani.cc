// body/gen.mjs — exercise-figure SVG generator. Run: `node body/gen.mjs`.
// Uses the shared figure core (body/core.mjs); writes flat schematic SVGs to body/svg3d/.
// No dependencies. Production render level = L2 (single smoothed silhouette).
import {mkdirSync, writeFileSync} from 'node:fs';
import {fileURLToPath} from 'node:url';
import {dirname, join} from 'node:path';
import {joints, figureSVG, skeletonSVG, figureExtent, VIEWS, project, FT, FR, D, R} from './core.mjs';

const DIR = join(dirname(fileURLToPath(import.meta.url)), 'svg3d');
mkdirSync(DIR, {recursive:true});
const r1 = x => Math.round(x*10)/10;
const DN = [0,-1,0];

// ---- poses (one 3D pose per exercise; arrows/safety/ghost optional) ----
const arc3 = (c,r,a0,a1,n=6) =>                       // arc in the Y-Z (frontal) plane for circular arrows
  Array.from({length:n},(_,i)=>a0+(a1-a0)*i/(n-1)).map(a=>[c[0], c[1]+r*Math.sin(R(a)), c[2]+r*Math.cos(R(a))]);

const NEUTRAL = {view:'front',pelvis:[0,0,0],spine:[0,1,0],front:[1,0,0],
 Luarm:[.07,-1,0],Lfarm:[.07,-1,0],Ruarm:[.07,-1,0],Rfarm:[.07,-1,0],
 Lthigh:DN,Lshin:DN,Lfoot:[1,-.04,0],Rthigh:DN,Rshin:DN,Rfoot:[1,-.04,0]};

const EX = {
 W05_arm_circles:{view:'front',pelvis:[0,0,0],spine:[0,1,0],front:[1,0,0],
   Luarm:[0,.05,1],Lfarm:[0,.05,1],Ruarm:[0,.05,-1],Rfarm:[0,.05,-1],
   Lthigh:DN,Lshin:DN,Lfoot:[1,-.04,0],Rthigh:DN,Rshin:DN,Rfoot:[1,-.04,0],
   arrows:J=>[[arc3(J.haL,34,200,-70),'end'],[arc3(J.haR,34,-20,250),'end']]},
 M01_push_up:{view:'side',pelvis:[0,0,0],spine:[-0.96,0.28,0],front:[-0.27,-0.96,0],
   Luarm:[0,-1,0],Lfarm:[0,-1,0],Ruarm:[0,-1,0],Rfarm:[0,-1,0],
   Lthigh:[0.95,-0.30,0],Lshin:[0.95,-0.30,0],Lfoot:[0.5,-0.86,0],Rthigh:[0.95,-0.30,0],Rshin:[0.95,-0.30,0],Rfoot:[0.5,-0.86,0],
   arrows:J=>[[[[J.chest[0]-70,J.chest[1]+45,0],[J.chest[0]-70,J.chest[1]-75,0]],'both']],
   safety:J=>[[J.head,J.anL]]},
 C03_flutter_kicks:{view:'side',pelvis:[0,0,0],spine:[-1,0,0],front:[0,1,0],
   Luarm:[1,0,0],Lfarm:[-0.62,0.79,0],Ruarm:[1,0,0],Rfarm:[-0.62,0.79,0],
   Lthigh:[0.45,0.89,0],Lshin:[0.45,0.89,0],Lfoot:[0.45,0.89,0],Rthigh:[0.66,0.75,0],Rshin:[0.66,0.75,0],Rfoot:[0.66,0.75,0],
   arrows:J=>[[[[J.toL[0],J.toL[1]+55,0],[J.toL[0],J.toL[1]-25,0]],'both'],[[[J.toR[0]+30,J.toR[1]+55,0],[J.toR[0]+30,J.toR[1]-25,0]],'both']],
   safety:J=>[[J.pel,[J.pel[0]-95,J.pel[1],0]]]},
 L01_mini_squat_45:{view:'side',pelvis:[0,300,0],spine:[0.32,0.95,0],front:[1,0,0],
   Luarm:[0.4,-0.85,0],Lfarm:[0.95,-0.05,0],Ruarm:[0.4,-0.85,0],Rfarm:[0.95,-0.05,0],
   Lthigh:[0.55,-0.84,0],Lshin:[-0.05,-1,0],Lfoot:[1,-0.04,0],Rthigh:[0.55,-0.84,0],Rshin:[-0.05,-1,0],Rfoot:[1,-0.04,0],
   ghost:{pelvis:[0,326,0],spine:[0,1,0],front:[1,0,0],
     Luarm:[0.12,-0.99,0],Lfarm:[0.12,-0.99,0],Ruarm:[0.12,-0.99,0],Rfarm:[0.12,-0.99,0],
     Lthigh:DN,Lshin:DN,Lfoot:[1,-0.04,0],Rthigh:DN,Rshin:DN,Rfoot:[1,-0.04,0]},
   arrows:J=>[[[[J.pel[0]-12,J.pel[1]+8,0],[J.pel[0]-75,J.pel[1]-42,0]],'end']],
   safety:J=>[[J.toL,[J.toL[0],J.toL[1]+210,0]]]},
};

// ---- render one exercise ----
function renderEx(id,title,pose,view,level=2,sex='M'){
 const cam = VIEWS[view];
 const J = joints(pose,sex);
 const ghostJ = pose.ghost ? joints(pose.ghost,sex) : null;
 const arrows = (pose.arrows?pose.arrows(J):[]);
 const safety = (pose.safety?pose.safety(J):[]);
 const extra = [...arrows.flatMap(([pts])=>pts), ...safety.flatMap(([a,b])=>[a,b])];
 // auto-fit over figure (+ ghost + arrows/safety)
 let e = figureExtent(J,cam,sex,extra);
 if(ghostJ){const g=figureExtent(ghostJ,cam,sex); e={minx:Math.min(e.minx,g.minx),maxx:Math.max(e.maxx,g.maxx),miny:Math.min(e.miny,g.miny),maxy:Math.max(e.maxy,g.maxy)};}
 const w=e.maxx-e.minx||1, h=e.maxy-e.miny||1;
 const sc=Math.min(880/w,600/h,1.25), tx=512-sc*(e.minx+e.maxx)/2, ty=700-sc*e.maxy;
 const xf={sc,tx,ty};
 const P2 = pt => {const a=project(pt,cam); return [tx+sc*a[0], ty+sc*a[1]];};
 // arrows + safety
 const arrowSVG = arrows.map(([pts3,style])=>{
   const d='M '+pts3.map(P2).map(([x,y])=>`${r1(x)},${r1(y)}`).join(' L ');
   const ms=(style==='both'?' marker-start="url(#ah)"':'')+(style==='both'||style==='end'?' marker-end="url(#ah)"':'');
   return `<path d="${d}" fill="none" stroke="#2563EB" stroke-width="6"${ms}/>`;}).join('');
 const safeSVG = safety.map(([a,b])=>{const[ax,ay]=P2(a),[bx,by]=P2(b);
   return `<line x1="${r1(ax)}" y1="${r1(ay)}" x2="${r1(bx)}" y2="${r1(by)}" stroke="#10B981" stroke-width="4" stroke-dasharray="2 14"/>`;}).join('');
 // floor (not on top view)
 let floor='';
 if(view!=='top'){const fy=r1(ty+sc*e.maxy+2), x1=r1(Math.max(60,tx+sc*e.minx-30)), x2=r1(Math.min(964,tx+sc*e.maxx+30));
   floor=`<line x1="${x1}" y1="${fy}" x2="${x2}" y2="${fy}" stroke="#1F2937" stroke-width="4" stroke-linecap="round" opacity="0.45"/>`;}
 const fig = level===0 ? skeletonSVG(J,cam,{xf}) : figureSVG(J,cam,{sex,level,xf});
 const ghost = ghostJ ? `<g opacity="0.4">${figureSVG(ghostJ,cam,{sex,level,ghost:true,xf})}</g>` : '';
 return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 768" role="img" aria-labelledby="${id}t"><title id="${id}t">${title}</title>
  <defs><path id="FT" d="${FT}"/><marker id="ah" markerUnits="userSpaceOnUse" markerWidth="26" markerHeight="26" refX="18" refY="13" orient="auto-start-reverse"><path d="M5,3 L21,13 L5,23" fill="none" stroke="#2563EB" stroke-width="5" stroke-linecap="round" stroke-linejoin="round"/></marker></defs>
  <g>${floor}</g>${ghost}${fig}${arrowSVG}${safeSVG}
</svg>
`;
}

// ---- generate ----
let n=0;
for(const [name,pose] of Object.entries(EX)){
 writeFileSync(join(DIR,`${name}.svg`), renderEx(name.slice(0,3),name,pose,pose.view,2)); n++;
}
// level demo on one pose
writeFileSync(join(DIR,'_demo_L0.svg'), renderEx('d0','L0 skeleton',EX.L01_mini_squat_45,'side',0));
writeFileSync(join(DIR,'_demo_L1.svg'), renderEx('d1','L1 capsules',EX.L01_mini_squat_45,'side',1));
writeFileSync(join(DIR,'_demo_L2.svg'), renderEx('d2','L2 outline',EX.L01_mini_squat_45,'side',2));
for(const sex of ['M','F']) writeFileSync(join(DIR,`_neutral_${sex}.svg`), renderEx('n',`neutral ${sex}`,NEUTRAL,'front',2,sex));
console.log(`generated ${n} exercises + demos into ${DIR}`);
