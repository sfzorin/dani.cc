// @ts-check
// body/core.mjs — single source of truth for the kinematic mannequin figure.
// Used by BOTH the /body editor (browser, via <script type="module">) and the
// SVG generator (Node, body/gen.mjs). No DOM, no fs — pure geometry + SVG strings.
//
// Output is in WORLD/projected units (no view scaling). The caller wraps the
// returned inner-SVG in a <g transform="translate(tx,ty) scale(s)"> to place it.
// Projection is orthographic; an arbitrary orbit camera or the 3 fixed views work.

export const H = 660;
export const PRESETS = {
  M:{shw:.245,hpw:.113,trunk:.300,uarm:.176,farm:.157,hand:.082,thigh:.232,shin:.247,foot:.055,head:.130},
  F:{shw:.200,hpw:.120,trunk:.290,uarm:.177,farm:.160,hand:.082,thigh:.249,shin:.250,foot:.055,head:.130},
};
export const FLESH = {
  M:{neck:.022,shCap:.038,uA:.034,uB:.025,elbow:.026,fA:.027,fB:.018,wrist:.017,handL:.060,handW:.038,
     hipB:.046,tA:.050,tB:.034,knee:.034,sK:.031,sC:.039,sA:.019,ankle:.019,footL:.140,footW:.058,footH:.042,
     ribU:.230,ribL:.165,waist:.150,pelW:.180,headW:.095,shDrop:.020,shOut:.006},
  F:{neck:.018,shCap:.031,uA:.028,uB:.021,elbow:.022,fA:.023,fB:.015,wrist:.014,handL:.055,handW:.032,
     hipB:.050,tA:.050,tB:.032,knee:.029,sK:.027,sC:.034,sA:.016,ankle:.016,footL:.128,footW:.050,footH:.036,
     ribU:.188,ribL:.140,waist:.128,pelW:.205,headW:.090,shDrop:.026,shOut:.006},
};
export const FT = "M 8,-6 C 34,-6 50,0 57,9 C 63,13 61,21 51,21 L -8,21 C -18,21 -18,-3 -9,-5 Z";
export const OUTLINE = 2.4, ST = 4.0;   // world units: L2 contour half-width, L1 per-part stroke
const INK = '#1F2937', BODY = '#E5E7EB', ACC = '#2563EB';

export const P  = (sex,k)=>PRESETS[sex][k]*H;
export const FR = (sex,k)=>FLESH[sex][k]*H;

// ---- vectors ----
export const nz=v=>{const m=Math.hypot(v[0],v[1],v[2])||1;return[v[0]/m,v[1]/m,v[2]/m];};
export const add=(p,d,l)=>{d=nz(d);return[p[0]+d[0]*l,p[1]+d[1]*l,p[2]+d[2]*l];};
export const addv=(p,u,l)=>[p[0]+u[0]*l,p[1]+u[1]*l,p[2]+u[2]*l];
export const cross=(a,b)=>[a[1]*b[2]-a[2]*b[1],a[2]*b[0]-a[0]*b[2],a[0]*b[1]-a[1]*b[0]];
export const dot=(a,b)=>a[0]*b[0]+a[1]*b[1]+a[2]*b[2];
export const scl=(v,s)=>[v[0]*s,v[1]*s,v[2]*s];
export const neg=v=>scl(v,-1);
export const D=x=>x*180/Math.PI, R=x=>x*Math.PI/180;
const r1=x=>Math.round(x*10)/10;

// ---- joint constraints ----
export function hinge(T,S,allowed,minf=3,maxf=145){T=nz(T);S=nz(S);
 let p=[allowed[0]-dot(allowed,T)*T[0],allowed[1]-dot(allowed,T)*T[1],allowed[2]-dot(allowed,T)*T[2]];
 if(p[0]*p[0]+p[1]*p[1]+p[2]*p[2]<1e-9) return add(T,allowed,.05); p=nz(p);
 let th=D(Math.atan2(dot(S,p),dot(S,T))); th=Math.max(minf,Math.min(maxf,th)); const r=R(th);
 return nz([Math.cos(r)*T[0]+Math.sin(r)*p[0],Math.cos(r)*T[1]+Math.sin(r)*p[1],Math.cos(r)*T[2]+Math.sin(r)*p[2]]);}
export function cone(axis,S,maxsw){const ax=nz(axis);S=nz(S);const c=Math.max(-1,Math.min(1,dot(S,ax)));
 if(D(Math.acos(c))<=maxsw) return S; let pp=[S[0]-c*ax[0],S[1]-c*ax[1],S[2]-c*ax[2]];
 if(pp[0]*pp[0]+pp[1]*pp[1]+pp[2]*pp[2]<1e-9) pp=(Math.abs(ax[1])<.9?cross(ax,[0,1,0]):cross(ax,[1,0,0])); pp=nz(pp); const r=R(maxsw);
 return nz([Math.cos(r)*ax[0]+Math.sin(r)*pp[0],Math.cos(r)*ax[1]+Math.sin(r)*pp[1],Math.cos(r)*ax[2]+Math.sin(r)*pp[2]]);}
// Directional (asymmetric) swing limit: max angle from `axis` varies by swing direction —
// wide toward `front` (maxF, e.g. flexion/overhead), tight backward (maxB, extension) and
// to the sides (maxS). Lets a shoulder reach overhead-forward yet not swing through the back.
export function swing(axis,S,front,maxF,maxB,maxS){const ax=nz(axis);S=nz(S);const fr=nz(front);
 const c=Math.max(-1,Math.min(1,dot(S,ax)));const th=D(Math.acos(c));
 if(th<0.5) return S;
 let p=[S[0]-c*ax[0],S[1]-c*ax[1],S[2]-c*ax[2]];const pm=Math.hypot(p[0],p[1],p[2]);
 if(pm<1e-6) return S; p=[p[0]/pm,p[1]/pm,p[2]/pm];
 const f=dot(p,fr);                                   // forwardness of the swing direction, -1..1
 const maxA=f>=0 ? maxS+(maxF-maxS)*f : maxS+(maxB-maxS)*(-f);
 if(th<=maxA) return S; const r=R(maxA);
 return nz([Math.cos(r)*ax[0]+Math.sin(r)*p[0],Math.cos(r)*ax[1]+Math.sin(r)*p[1],Math.cos(r)*ax[2]+Math.sin(r)*p[2]]);}

// ---- camera / projection (orthographic) ----
// cam = {right,up,fwd} orthonormal basis. project → [screenX, screenY, depth].
// fwd points camera→scene, so larger depth = farther = drawn first (occlusion).
export const VIEWS = {
  side:  {right:[1,0,0], up:[0,1,0], fwd:[0,0,1]},   // profile: (x,-y) depth z
  front: {right:[0,0,1], up:[0,1,0], fwd:[1,0,0]},   // (z,-y) depth x
  top:   {right:[0,0,1], up:[1,0,0], fwd:[0,1,0]},   // (z,-x) depth y
};
export function cameraBasis(yaw,pitch){
  const a=R(yaw), e=R(pitch);
  const fwd=nz([-Math.cos(e)*Math.sin(a), -Math.sin(e), -Math.cos(e)*Math.cos(a)]);
  let right=cross(fwd,[0,1,0]);
  if(right[0]*right[0]+right[1]*right[1]+right[2]*right[2]<1e-6) right=[1,0,0];
  right=nz(right);
  const up=nz(cross(right,fwd));
  return {right,up,fwd};
}
export const project=(p,cam)=>[dot(p,cam.right), -dot(p,cam.up), dot(p,cam.fwd)];

// ---- skeleton ----
export function joints(pose,sex='M'){
 const p=k=>P(sex,k);
 const pel=pose.pelvis, spine=cone([0,1,0],nz(pose.spine),80), front=nz(pose.front), lat=nz(cross(spine,front));
 const sh_y=.04*H+p('shin')+p('thigh')+p('trunk');
 const neck_len=(H-p('head'))-sh_y, head_off=(H-p('head')/2)-sh_y;
 const chest=add(pel,spine,p('trunk'));
 const neck=add(chest,spine,neck_len);
 const headDir=cone(spine,nz(pose.head||spine),45);
 const J={pel,chest,neck,head:add(neck,headDir,head_off-neck_len),spine,front,lat,hrx:p('head')*.34,hry:p('head')/2};
 for(const[s,sg] of [['L',1],['R',-1]]){
  J['hip'+s]=addv(pel,lat,sg*p('hpw')/2); J['sh'+s]=addv(chest,lat,sg*p('shw')/2);
  const ua=swing(neg(spine),pose[s+'uarm'],front,178,55,168), fa=hinge(ua,pose[s+'farm'],front);
  const hd=cone(fa,nz(pose[s+'hand']||fa),90);   // wrist ROM: allow a flat palm even off a vertical forearm
  J['el'+s]=add(J['sh'+s],ua,p('uarm')); J['wr'+s]=add(J['el'+s],fa,p('farm')); J['ha'+s]=add(J['wr'+s],hd,p('hand'));
  const th=swing(neg(spine),pose[s+'thigh'],front,122,28,52), sd=hinge(th,pose[s+'shin'],neg(front)), ft=cone(front,pose[s+'foot'],75);
  J['kn'+s]=add(J['hip'+s],th,p('thigh')); J['an'+s]=add(J['kn'+s],sd,p('shin')); J['to'+s]=add(J['an'+s],ft,p('foot'));
 }
 return J;
}

export function neutralPose(){return {pelvis:[0,326,0],spine:[0,1,0],front:[1,0,0],head:[0,1,0],
 Luarm:[.07,-1,0],Lfarm:[.07,-1,0],Lhand:[.07,-1,0],Ruarm:[.07,-1,0],Rfarm:[.07,-1,0],Rhand:[.07,-1,0],
 Lthigh:[0,-1,0],Lshin:[0,-1,0],Lfoot:[1,-.04,0],Rthigh:[0,-1,0],Rshin:[0,-1,0],Rfoot:[1,-.04,0]};}

// ---- rendering ----
function capPath(ax,ay,bx,by,rA,rB){const dx=bx-ax,dy=by-ay;const L=Math.hypot(dx,dy)||1;const ux=dx/L,uy=dy/L,px=-uy,py=ux;
 const ra=r1(rA),rb=r1(rB);const f=(x,y)=>r1(x)+','+r1(y);
 return `M ${f(ax+px*rA,ay+py*rA)} L ${f(bx+px*rB,by+py*rB)} A ${rb} ${rb} 0 0 1 ${f(bx-px*rB,by-py*rB)} L ${f(ax-px*rA,ay-py*rA)} A ${ra} ${ra} 0 0 1 ${f(ax+px*rA,ay+py*rA)} Z`;}

const mid3=(a,b,t)=>[a[0]+(b[0]-a[0])*t,a[1]+(b[1]-a[1])*t,a[2]+(b[2]-a[2])*t];

/**
 * Build the figure inner-SVG for a camera, with the placement transform BAKED into
 * coordinates (no <g transform> wrapper — a scaled group of heavy overlapping fills
 * triggers a Chrome rasterizer flake where gray fails to cover dark; baking avoids it).
 * opts: {sex='M', level=2, ghost=false, xf:{sc,tx,ty}}  screen = xf.tx + xf.sc*projX, ...
 * L2 = two-pass expand-stroke (dark fat-stroke pass + gray fill pass) → uniform contour,
 *      no internal seams, concave notches not filled (survives any angle).
 * L1 = each part filled + thin stroke, depth-sorted (segmented capsule view).
 */
export function figureSVG(J,cam,opts={}){
 const sex=opts.sex||'M', level=opts.level==null?2:opts.level, ghost=!!opts.ghost, L2=level===2;
 const xf=opts.xf||{sc:1,tx:0,ty:0}, sc=xf.sc;
 const fillc=ghost?'#D5D9DF':BODY, dark=ghost?'#9CA3AF':INK;
 const fr=k=>FR(sex,k)*sc, p=k=>P(sex,k);   // fr → SCALED radius (screen units)
 const S=q=>{const a=project(q,cam);return[xf.tx+sc*a[0], xf.ty+sc*a[1], a[2]];};
 const pr={}, dz={};
 const KEYS=['pel','chest','neck','head','hipL','hipR','shL','shR','elL','elR','wrL','wrR','haL','haR','knL','knR','anL','anR','toL','toR'];
 for(const k of KEYS){const a=S(J[k]); pr[k]=[a[0],a[1]]; dz[k]=a[2];}
 const DSTR=` stroke="${dark}" stroke-width="${r1(2*OUTLINE*sc)}" stroke-linejoin="round" stroke-linecap="round"`;
 const STR=` stroke="${dark}" stroke-width="${r1(ST*sc)}" stroke-linejoin="round" stroke-linecap="round"`;
 // L2 = two-pass: every part DARK (fat dark stroke = expanded outline) then every part GRAY on
 //      top. Separate elements at final baked coords (the proven approach). Dedup identical shapes
 //      (coincident L/R) to cut overdraw. L1 = depth-sorted filled+stroked parts (capsule view).
 const darkArr=[], grayArr=[], parts=[], seen=new Set();
 const emit=(mk,depth)=>{ if(L2){const d=mk(dark,DSTR),g=mk(fillc,'');
     if(!seen.has(d)){seen.add(d);darkArr.push(d);} if(!seen.has(g)){seen.add(g);grayArr.push(g);}}
   else {const q=mk(fillc,STR); if(!seen.has(q)){seen.add(q);parts.push([depth,q]);}} };
 const circleC=(x,y,r,depth)=>emit((f,st)=>`<circle cx="${r1(x)}" cy="${r1(y)}" r="${r1(Math.max(1,r))}" fill="${f}"${st}/>`,depth);
 const capC=(ax,ay,bx,by,rA,rB,depth)=>{const L=Math.hypot(bx-ax,by-ay);
   if(L<Math.max(rA,rB)*0.7){circleC((ax+bx)/2,(ay+by)/2,Math.max(rA,rB),depth);return;}  // bone seen end-on → knob
   emit((f,st)=>`<path d="${capPath(ax,ay,bx,by,rA,rB)}" fill="${f}"${st}/>`,depth);};
 const cap=(A,B,rA,rB,depth)=>{const[ax,ay]=pr[A],[bx,by]=pr[B]; capC(ax,ay,bx,by,rA,rB,depth);};
 const knob=(K,r,depth)=>{const[x,y]=pr[K]; circleC(x,y,r,depth);};
 const head=(depth)=>{const[cx,cy]=pr.head;const[nx,ny]=pr.neck;const a=r1(D(Math.atan2(cx-nx,-(cy-ny))));
   emit((f,st)=>`<ellipse cx="${r1(cx)}" cy="${r1(cy)}" rx="${r1(J.hrx*sc)}" ry="${r1(J.hry*sc)}" fill="${f}"${st} transform="rotate(${a} ${r1(cx)} ${r1(cy)})"/>`,depth);};
 // Foot = a rounded WEDGE, not an oval around the ankle. The ankle sits on top/back of
 // the foot; heel reaches behind, toe ahead, and the SOLE is a flat line. Built in screen
 // space along the toe direction; the up-normal is flipped to point toward the leg so the
 // sole always lands on the far side from the shin (flat on the floor when standing).
 const foot=(AN,TO,KN,depth)=>{const A=pr[AN],T=pr[TO];let ux=T[0]-A[0],uy=T[1]-A[1];const ln=Math.hypot(ux,uy);
   if(ln < p('foot')*sc*0.42){const r=fr('footW')/2;  // foot seen end-on → "toe toward camera" oval
     emit((f,st)=>`<ellipse cx="${r1(A[0])}" cy="${r1(A[1]+12*sc)}" rx="${r1(r)}" ry="${r1(r*.7)}" fill="${f}"${st}/>`,depth+.05);return;}
   ux/=ln;uy/=ln;let nx=-uy,ny=ux;const K=pr[KN];const lx=A[0]-K[0],ly=A[1]-K[1]; // knee→ankle (down the leg)
   if(nx*(-lx)+ny*(-ly)<0){nx=-nx;ny=-ny;}                       // n points UP the leg; -n is the sole side
   const fl=fr('footL'),fh=fr('footH'),hb=fl*0.26,tf=fl*0.74;    // heel-back / toe-forward split
   const P=(du,dn)=>[A[0]+ux*du+nx*dn,A[1]+uy*du+ny*dn], pt=a=>r1(a[0])+','+r1(a[1]);
   const heelTop=P(-hb*0.7, fh*0.18), heelBot=P(-hb,-fh*0.55),   // achilles/heel
         toeBot=P(tf,-fh), toeTop=P(tf*0.95,-fh*0.42), instep=P(tf*0.38, fh*0.06);
   const d=`M ${pt(heelTop)} Q ${pt(P(-hb,-fh*0.05))} ${pt(heelBot)} L ${pt(toeBot)}`
         +` Q ${pt(P(tf+hb*0.3,-fh))} ${pt(toeTop)} Q ${pt(instep)} ${pt(heelTop)} Z`;
   emit((f,st)=>`<path d="${d}" fill="${f}"${st}/>`,depth+.05);};
 const torso=(depth)=>{
   // Lofted trunk = three soft volumes (pelvis → waist → ribcage), NOT one trapezoid.
   // Cross-section ellipse (lateral radius w, sagittal radius h) projected to screen; extra
   // samples + curved sides give a waist pinch and rounded shoulders instead of flat facets.
   const prof=[[0,fr('pelW')/2,fr('pelW')/2*.6],[.20,fr('pelW')/2*.97,fr('pelW')/2*.6],
     [.46,fr('waist')/2,fr('waist')/2*.62],[.70,fr('ribL')/2,fr('ribL')/2*.62],
     [.88,fr('ribU')/2*.99,fr('ribU')/2*.62],[1,fr('ribU')/2*.90,fr('ribU')/2*.62]];
   const[px0,py0]=pr.pel,[qx,qy]=pr.chest; let dx=qx-px0,dy=qy-py0; const dl=Math.hypot(dx,dy)||1; const ux=dx/dl,uy=dy/dl,nx=-uy,ny=ux;
   const latS=[dot(J.lat,cam.right),-dot(J.lat,cam.up)], frS=[dot(J.front,cam.right),-dot(J.front,cam.up)];
   const latN=latS[0]*nx+latS[1]*ny, frN=frS[0]*nx+frS[1]*ny;   // ellipse axis components along screen-normal
   const sp=prof.map(([t])=>{const a=S(addv(J.pel,J.spine,t*p('trunk')));return[a[0],a[1]];});
   const rs=prof.map(([,w,h])=>Math.hypot(w*latN,h*frN));        // silhouette half-width perp to spine (w,h already scaled)
   const sm=(pts)=>{let d='';const m=pts.length-1;for(let i=1;i<m;i++){const xc=(pts[i][0]+pts[i+1][0])/2,yc=(pts[i][1]+pts[i+1][1])/2;d+=` Q ${r1(pts[i][0])},${r1(pts[i][1])} ${r1(xc)},${r1(yc)}`;}d+=` L ${r1(pts[m][0])},${r1(pts[m][1])}`;return d;};
   emit((f,st)=>{const Lp=[],Rp=[];sp.forEach(([sx,sy],i)=>{const r=rs[i];Lp.push([sx+nx*r,sy+ny*r]);Rp.push([sx-nx*r,sy-ny*r]);});
     let d="M "+r1(Lp[0][0])+","+r1(Lp[0][1]);
     d+=sm(Lp);                                                                                            // left side, curved
     d+=` Q ${r1(qx+ux*7*sc)},${r1(qy+uy*7*sc)} ${r1(Rp[Rp.length-1][0])},${r1(Rp[Rp.length-1][1])}`;      // rounded top
     d+=sm(Rp.slice().reverse());                                                                          // right side, curved
     d+=` Q ${r1(px0-ux*4*sc)},${r1(py0-uy*4*sc)} ${r1(Lp[0][0])},${r1(Lp[0][1])} Z`;                      // rounded bottom
     return `<path d="${d}" fill="${f}"${st}/>`;},depth);};
 // Hand = a rounded mitten (narrow wrist → wider palm → softly rounded tip), NOT a capsule.
 // length > width avoids the "flipper"; a small thumb hint pops out toward the front side
 // (scaled by how face-on the hand is, so it vanishes when the hand is edge-on).
 const hand=(sd,depth)=>{const W=pr['wr'+sd],E=pr['ha'+sd];let ux=E[0]-W[0],uy=E[1]-W[1];let ln=Math.hypot(ux,uy);
   if(ln<1){ux=0;uy=1;ln=fr('handL');} ux/=ln;uy/=ln;const nx=-uy,ny=ux;
   const ph=fr('handW')/2,wh=ph*0.5,th=ph*0.78;                  // palm / wrist / tip half-widths
   const P=(du,dn)=>[W[0]+ux*du+nx*dn,W[1]+uy*du+ny*dn], pt=a=>r1(a[0])+','+r1(a[1]);
   const wl=P(0,wh),wr=P(0,-wh),pl=P(ln*0.45,ph),pR=P(ln*0.45,-ph),tl=P(ln*0.9,th),tr=P(ln*0.9,-th),apex=P(ln*1.06,0),wb=P(-wh*0.7,0);
   const d=`M ${pt(wl)} L ${pt(pl)} L ${pt(tl)} Q ${pt(apex)} ${pt(tr)} L ${pt(pR)} L ${pt(wr)} Q ${pt(wb)} ${pt(wl)} Z`;
   emit((f,st)=>`<path d="${d}" fill="${f}"${st}/>`,depth);
   const frS=[dot(J.front,cam.right),-dot(J.front,cam.up)],fp=frS[0]*nx+frS[1]*ny;     // thumb on the forward side
   if(Math.abs(fp)>0.4){const s=fp>=0?1:-1,c=P(ln*0.32,s*(ph+wh*0.15));circleC(c[0],c[1],wh*0.95,depth);}};  // soft thumb bump (no spike)
 // Shoulders = a sloping trapezius yoke from the neck base down to a rounded deltoid cap that
 // sits below+outside the joint — so the top is NOT a flat bar and the joint is not a sharp
 // corner. The torso reaches the ribcage; the deltoid caps the arm root. (All per-camera via S.)
 const shoulders=(yokeD)=>{const DROP=FR(sex,'shDrop'),OUT=FR(sex,'shOut'),nH=FR(sex,'neck')*1.15;
   for(const sd of['L','R']){const sg=sd==='L'?1:-1,d=dz['sh'+sd];
     const capW=addv(addv(J['sh'+sd],J.spine,-DROP),J.lat,sg*OUT), ncW=addv(J.chest,J.lat,sg*nH);
     const C=S(capW), N=S(ncW);
     capC(N[0],N[1],C[0],C[1],fr('neck')*0.95,fr('shCap'),yokeD);                 // trapezius slope
     capC(C[0],C[1],pr['el'+sd][0],pr['el'+sd][1],fr('uA'),fr('uB'),d);           // upper arm from deltoid
     if(!ghost) circleC(C[0],C[1],fr('shCap'),d+.4);}};                           // rounded deltoid cap
 // build far→near (single calf cap = fewer overlapping shapes)
 for(const sd of['L','R']){
   cap('hip'+sd,'kn'+sd,fr('tA'),fr('tB'),dz['hip'+sd]);
   cap('kn'+sd,'an'+sd,fr('sC'),fr('sA'),dz['hip'+sd]);
   foot('an'+sd,'to'+sd,'kn'+sd,dz['hip'+sd]);}
 torso(dz.chest);
 cap('chest','neck',fr('neck'),fr('neck')*.95,dz.chest-0.2);
 shoulders(dz.chest-0.1);
 for(const sd of['L','R']){cap('el'+sd,'wr'+sd,fr('fA'),fr('fB'),dz['sh'+sd]); hand(sd,dz['sh'+sd]);}
 if(!ghost) for(const sd of['L','R']){   // joint knobs; knee/elbow slightly oversized so they read as the joint
   knob('hip'+sd,fr('hipB'),dz['hip'+sd]+.3);
   knob('el'+sd,fr('elbow')*1.18,dz['el'+sd]+.4);knob('kn'+sd,fr('knee')*1.18,dz['kn'+sd]+.4);}
 head(dz.head+.2);
 if(L2) return darkArr.join('')+grayArr.join('');
 parts.sort((a,b)=>a[0]-b[0]); return parts.map(x=>x[1]).join('');
}

/** L0 skeleton inner-SVG with the placement transform baked in (opts.xf). */
export function skeletonSVG(J,cam,opts={}){
 const xf=opts.xf||{sc:1,tx:0,ty:0}, sc=xf.sc;
 const S=q=>{const a=project(q,cam);return[xf.tx+sc*a[0], xf.ty+sc*a[1]];};
 const pr={}; for(const k in J) if(Array.isArray(J[k]))pr[k]=S(J[k]);
 const bw=r1(3*sc), dr=r1(6*sc);
 const bones=[['pel','chest'],['chest','neck'],['neck','head'],['hipL','hipR'],['shL','shR'],['chest','shL'],['chest','shR'],['pel','hipL'],['pel','hipR'],
  ['hipL','knL'],['knL','anL'],['anL','toL'],['hipR','knR'],['knR','anR'],['anR','toR'],['shL','elL'],['elL','wrL'],['wrL','haL'],['shR','elR'],['elR','wrR'],['wrR','haR']];
 let s=''; for(const[a,b] of bones)s+=`<line x1="${r1(pr[a][0])}" y1="${r1(pr[a][1])}" x2="${r1(pr[b][0])}" y2="${r1(pr[b][1])}" stroke="${INK}" stroke-width="${bw}"/>`;
 for(const k of['pel','chest','neck','hipL','hipR','shL','shR','elL','wrL','haL','elR','wrR','haR','knL','anL','knR','anR'])s+=`<circle cx="${r1(pr[k][0])}" cy="${r1(pr[k][1])}" r="${dr}" fill="${ACC}"/>`;
 const[hx,hy]=pr.head;s+=`<circle cx="${r1(hx)}" cy="${r1(hy)}" r="${r1(J.hry*sc)}" fill="none" stroke="${ACC}" stroke-width="${bw}"/>`;
 return s;
}

/** Approx projected bounding box of the figure (for auto-fit), padded per joint. */
export function figureExtent(J,cam,sex='M',extraPts=[]){
 const PAD={head:J.hry,hipL:FR(sex,'hipB'),hipR:FR(sex,'hipB'),shL:FR(sex,'shCap'),shR:FR(sex,'shCap'),
   knL:FR(sex,'tB'),knR:FR(sex,'tB'),elL:FR(sex,'fA'),elR:FR(sex,'fA'),haL:FR(sex,'handW'),haR:FR(sex,'handW'),
   toL:FR(sex,'footW'),toR:FR(sex,'footW'),anL:FR(sex,'footW'),anR:FR(sex,'footW'),chest:FR(sex,'ribU')/2};
 let minx=1e9,maxx=-1e9,miny=1e9,maxy=-1e9;
 const acc=(x,y,pad)=>{minx=Math.min(minx,x-pad);maxx=Math.max(maxx,x+pad);miny=Math.min(miny,y-pad);maxy=Math.max(maxy,y+pad);};
 for(const k in J) if(Array.isArray(J[k])){const a=project(J[k],cam);acc(a[0],a[1],PAD[k]||FR(sex,'wrist'));}
 for(const pt of extraPts){const a=project(pt,cam);acc(a[0],a[1],0);}
 return {minx,maxx,miny,maxy};
}
