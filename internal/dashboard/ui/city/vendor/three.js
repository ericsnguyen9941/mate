/*
 * three.js r186 and the addons Agent City uses (RoundedBoxGeometry, SVGLoader,
 * BufferGeometryUtils.mergeGeometries, RoomEnvironment, OrbitControls), as one
 * minified ES module so the page loads it from this origin with no build step
 * at run time. Rebuilt from the npm package with:
 *
 *   export * from 'three';
 *   export { RoundedBoxGeometry } from 'three/addons/geometries/RoundedBoxGeometry.js';
 *   export { SVGLoader } from 'three/addons/loaders/SVGLoader.js';
 *   export { mergeGeometries } from 'three/addons/utils/BufferGeometryUtils.js';
 *   export { RoomEnvironment } from 'three/addons/environments/RoomEnvironment.js';
 *   export { OrbitControls } from 'three/addons/controls/OrbitControls.js';
 *
 *   bun build entry.js --minify --target=browser --format=esm --outfile=three.js
 *
 * The MIT License
 *
 * Copyright © 2010-2026 three.js authors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
var ed="186",di={LEFT:0,MIDDLE:1,RIGHT:2,ROTATE:0,DOLLY:1,PAN:2},fi={ROTATE:0,PAN:1,DOLLY_PAN:2,DOLLY_ROTATE:3},nd=0,Al=1,id=2,_m=3,xm=0,sr=1,sd=2,Ss=3,Ms=0,Ve=1,on=2,Ln=0,rr=1,wl=2,Rl=3,Cl=4,rd=5,ym=6,bs=100,ad=101,od=102,ld=103,cd=104,hd=200,ud=201,dd=202,fd=203,pd=204,md=205,gd=206,_d=207,xd=208,yd=209,vd=210,Sd=211,Md=212,bd=213,Td=214,Ed=0,Ad=1,wd=2,Il=3,Rd=4,Cd=5,Id=6,Pd=7,Ld=0,Nd=1,Dd=2,bn=0,Pl=1,Ll=2,Nl=3,Dl=4,Ul=5,Fl=6,Ol=7,vm="attached",Sm="detached",Mm=300,Ts=301,Fi=302,Ia=303,Pa=304,ar=306,La=1000,Es=1001,Na=1002,pi=1003,Ud=1004,bm=1004,or=1005,Tm=1005,Qe=1006,Da=1007,Em=1007,Oi=1008,Am=1008,Tn=1009,Fd=1010,Od=1011,lr=1012,Bl=1013,mi=1014,$n=1015,Nn=1016,zl=1017,kl=1018,As=1020,Bd=35902,zd=35899,kd=1021,Gd=1022,Dn=1023,Bi=1026,zi=1027,Hd=1028,Gl=1029,ki=1030,Hl=1031,wm=1032,Vl=1033,Ua=33776,Fa=33777,Oa=33778,Ba=33779,Wl=35840,Xl=35841,ql=35842,Yl=35843,Zl=36196,Jl=37492,$l=37496,Kl=37488,jl=37489,za=37490,Ql=37491,tc=37808,ec=37809,nc=37810,ic=37811,sc=37812,rc=37813,ac=37814,oc=37815,lc=37816,cc=37817,hc=37818,uc=37819,dc=37820,fc=37821,pc=36492,mc=36494,gc=36495,_c=36283,xc=36284,ka=36285,yc=36286,Rm=2200,Cm=2201,Im=2202,Pm=2300,Lm=2301,Nm=2302,Dm=2303,Um=2400,Fm=2401,Om=2402,Bm=2500,zm=2501,Vd=0,Wd=1,Xd=2,km=3200,Gm=3201,Hm=3202,Vm=3203,vc=0,qd=1,Gi="",Ga="srgb",Sc="srgb-linear",Mc="linear",xe="srgb",Wm="",Xm="rg",qm="ga",Ym=0,Zm=7680,Jm=7681,$m=7682,Km=7683,jm=34055,Qm=34056,tg=5386,eg=512,ng=513,ig=514,sg=515,rg=516,ag=517,og=518,lg=519,Yd=512,Zd=513,Jd=514,Ha=515,$d=516,Kd=517,Va=518,jd=519,cg=35044,hg=35048,ug=35040,dg=35045,fg=35049,pg=35041,mg=35046,gg=35050,_g=35042,xg="100",bc="300 es",Tc=2000,yg=2001,vg={COMPUTE:"compute",RENDER:"render"},Sg={PERSPECTIVE:"perspective",LINEAR:"linear",FLAT:"flat"},Mg={NORMAL:"normal",CENTROID:"centroid",SAMPLE:"sample",FIRST:"first",EITHER:"either"},bg={TEXTURE_COMPARE:"depthTextureCompare"},Tg={NONE:0,SHARED:1,FULL:2};function Eg(t){for(let e=t.length-1;e>=0;--e)if(t[e]>=65535)return!0;return!1}var Ag={Int8Array,Uint8Array,Uint8ClampedArray,Int16Array,Uint16Array,Int32Array,Uint32Array,Float32Array,Float64Array};function ms(t,e){return new Ag[t](e)}function Qd(t){return ArrayBuffer.isView(t)&&!(t instanceof DataView)}function _s(t){return document.createElementNS("http://www.w3.org/1999/xhtml",t)}function tf(){let t=_s("canvas");return t.style.display="block",t}var Wh={},li=null;function wg(t){li=t}function Rg(){return li}function tr(...t){let e="THREE."+t.shift();if(li)li("log",e,...t);else console.log(e,...t)}function ef(t){let e=t[0];if(typeof e==="string"&&e.startsWith("TSL:")){let n=t[1];if(n&&n.isStackTrace)t[0]+=" "+n.getLocation();else t[1]='Stack trace not available. Enable "THREE.Node.captureStackTrace" to capture stack traces.'}return t}function Mt(...t){t=ef(t);let e="THREE."+t.shift();if(li)li("warn",e,...t);else{let n=t[0];if(n&&n.isStackTrace)console.warn(n.getError(e));else console.warn(e,...t)}}function Wt(...t){t=ef(t);let e="THREE."+t.shift();if(li)li("error",e,...t);else{let n=t[0];if(n&&n.isStackTrace)console.error(n.getError(e));else console.error(e,...t)}}function qn(...t){let e=t.join(" ");if(e in Wh)return;Wh[e]=!0,Mt(...t)}function nf(t,e,n){return new Promise(function(i,s){function r(){switch(t.clientWaitSync(e,t.SYNC_FLUSH_COMMANDS_BIT,0)){case t.WAIT_FAILED:s();break;case t.TIMEOUT_EXPIRED:setTimeout(r,n);break;default:i()}}setTimeout(r,n)})}var sf={[0]:1,[2]:6,[4]:7,[3]:5,[1]:0,[6]:2,[7]:4,[5]:3};class pn{addEventListener(t,e){if(this._listeners===void 0)this._listeners={};let n=this._listeners;if(n[t]===void 0)n[t]=[];if(n[t].indexOf(e)===-1)n[t].push(e)}hasEventListener(t,e){let n=this._listeners;if(n===void 0)return!1;return n[t]!==void 0&&n[t].indexOf(e)!==-1}removeEventListener(t,e){let n=this._listeners;if(n===void 0)return;let i=n[t];if(i!==void 0){let s=i.indexOf(e);if(s!==-1)i.splice(s,1)}}dispatchEvent(t){let e=this._listeners;if(e===void 0)return;let n=e[t.type];if(n!==void 0){t.target=this;let i=n.slice(0);for(let s=0,r=i.length;s<r;s++)i[s].call(this,t);t.target=null}}}var ke=["00","01","02","03","04","05","06","07","08","09","0a","0b","0c","0d","0e","0f","10","11","12","13","14","15","16","17","18","19","1a","1b","1c","1d","1e","1f","20","21","22","23","24","25","26","27","28","29","2a","2b","2c","2d","2e","2f","30","31","32","33","34","35","36","37","38","39","3a","3b","3c","3d","3e","3f","40","41","42","43","44","45","46","47","48","49","4a","4b","4c","4d","4e","4f","50","51","52","53","54","55","56","57","58","59","5a","5b","5c","5d","5e","5f","60","61","62","63","64","65","66","67","68","69","6a","6b","6c","6d","6e","6f","70","71","72","73","74","75","76","77","78","79","7a","7b","7c","7d","7e","7f","80","81","82","83","84","85","86","87","88","89","8a","8b","8c","8d","8e","8f","90","91","92","93","94","95","96","97","98","99","9a","9b","9c","9d","9e","9f","a0","a1","a2","a3","a4","a5","a6","a7","a8","a9","aa","ab","ac","ad","ae","af","b0","b1","b2","b3","b4","b5","b6","b7","b8","b9","ba","bb","bc","bd","be","bf","c0","c1","c2","c3","c4","c5","c6","c7","c8","c9","ca","cb","cc","cd","ce","cf","d0","d1","d2","d3","d4","d5","d6","d7","d8","d9","da","db","dc","dd","de","df","e0","e1","e2","e3","e4","e5","e6","e7","e8","e9","ea","eb","ec","ed","ee","ef","f0","f1","f2","f3","f4","f5","f6","f7","f8","f9","fa","fb","fc","fd","fe","ff"],Xh=1234567,Ni=Math.PI/180,Di=180/Math.PI;function rn(){let t=Math.random()*4294967295|0,e=Math.random()*4294967295|0,n=Math.random()*4294967295|0,i=Math.random()*4294967295|0;return(ke[t&255]+ke[t>>8&255]+ke[t>>16&255]+ke[t>>24&255]+"-"+ke[e&255]+ke[e>>8&255]+"-"+ke[e>>16&15|64]+ke[e>>24&255]+"-"+ke[n&63|128]+ke[n>>8&255]+"-"+ke[n>>16&255]+ke[n>>24&255]+ke[i&255]+ke[i>>8&255]+ke[i>>16&255]+ke[i>>24&255]).toLowerCase()}function Zt(t,e,n){return Math.max(e,Math.min(n,t))}function Ec(t,e){return(t%e+e)%e}function Cg(t,e,n,i,s){return i+(t-e)*(s-i)/(n-e)}function Ig(t,e,n){if(t!==e)return(n-t)/(e-t);else return 0}function $s(t,e,n){return(1-n)*t+n*e}function Pg(t,e,n,i){return $s(t,e,1-Math.exp(-n*i))}function Lg(t,e=1){return e-Math.abs(Ec(t,e*2)-e)}function Ng(t,e,n){if(t<=e)return 0;if(t>=n)return 1;return t=(t-e)/(n-e),t*t*(3-2*t)}function Dg(t,e,n){if(t<=e)return 0;if(t>=n)return 1;return t=(t-e)/(n-e),t*t*t*(t*(t*6-15)+10)}function Ug(t,e){return t+Math.floor(Math.random()*(e-t+1))}function Fg(t,e){return t+Math.random()*(e-t)}function Og(t){return t*(0.5-Math.random())}function Bg(t){if(t!==void 0)Xh=t;let e=Xh+=1831565813;return e=Math.imul(e^e>>>15,e|1),e^=e+Math.imul(e^e>>>7,e|61),((e^e>>>14)>>>0)/4294967296}function zg(t){return t*Ni}function kg(t){return t*Di}function Gg(t){return t>0&&Number.isInteger(t)&&2**Math.round(Math.log2(t))===t}function Hg(t){return Math.pow(2,Math.ceil(Math.log(t)/Math.LN2))}function Vg(t){return Math.pow(2,Math.floor(Math.log(t)/Math.LN2))}function Wg(t,e,n,i,s){let{cos:r,sin:a}=Math,o=r(n/2),l=a(n/2),c=r((e+i)/2),h=a((e+i)/2),d=r((e-i)/2),u=a((e-i)/2),f=r((i-e)/2),g=a((i-e)/2);switch(s){case"XYX":t.set(o*h,l*d,l*u,o*c);break;case"YZY":t.set(l*u,o*h,l*d,o*c);break;case"ZXZ":t.set(l*d,l*u,o*h,o*c);break;case"XZX":t.set(o*h,l*g,l*f,o*c);break;case"YXY":t.set(l*f,o*h,l*g,o*c);break;case"ZYZ":t.set(l*g,l*f,o*h,o*c);break;default:Mt("MathUtils: .setQuaternionFromProperEuler() encountered an unknown order: "+s)}}function Ze(t,e){switch(e.constructor){case Float32Array:return t;case Uint32Array:return t/4294967295;case Uint16Array:return t/65535;case Uint8Array:case Uint8ClampedArray:return t/255;case Int32Array:return Math.max(t/2147483647,-1);case Int16Array:return Math.max(t/32767,-1);case Int8Array:return Math.max(t/127,-1);default:throw Error("THREE.MathUtils: Invalid component type.")}}function te(t,e){switch(e.constructor){case Float32Array:return t;case Uint32Array:return Math.round(t*4294967295);case Uint16Array:return Math.round(t*65535);case Uint8Array:case Uint8ClampedArray:return Math.round(t*255);case Int32Array:return Math.round(t*2147483647);case Int16Array:return Math.round(t*32767);case Int8Array:return Math.round(t*127);default:throw Error("THREE.MathUtils: Invalid component type.")}}var Ac={DEG2RAD:Ni,RAD2DEG:Di,generateUUID:rn,clamp:Zt,euclideanModulo:Ec,mapLinear:Cg,inverseLerp:Ig,lerp:$s,damp:Pg,pingpong:Lg,smoothstep:Ng,smootherstep:Dg,randInt:Ug,randFloat:Fg,randFloatSpread:Og,seededRandom:Bg,degToRad:zg,radToDeg:kg,isPowerOfTwo:Gg,ceilPowerOfTwo:Hg,floorPowerOfTwo:Vg,setQuaternionFromProperEuler:Wg,normalize:te,denormalize:Ze};class rt{static{rt.prototype.isVector2=!0}constructor(t=0,e=0){this.x=t,this.y=e}get width(){return this.x}set width(t){this.x=t}get height(){return this.y}set height(t){this.y=t}set(t,e){return this.x=t,this.y=e,this}setScalar(t){return this.x=t,this.y=t,this}setX(t){return this.x=t,this}setY(t){return this.y=t,this}setComponent(t,e){switch(t){case 0:this.x=e;break;case 1:this.y=e;break;default:throw Error("THREE.Vector2: index is out of range: "+t)}return this}getComponent(t){switch(t){case 0:return this.x;case 1:return this.y;default:throw Error("THREE.Vector2: index is out of range: "+t)}}clone(){return new this.constructor(this.x,this.y)}copy(t){return this.x=t.x,this.y=t.y,this}add(t){return this.x+=t.x,this.y+=t.y,this}addScalar(t){return this.x+=t,this.y+=t,this}addVectors(t,e){return this.x=t.x+e.x,this.y=t.y+e.y,this}addScaledVector(t,e){return this.x+=t.x*e,this.y+=t.y*e,this}sub(t){return this.x-=t.x,this.y-=t.y,this}subScalar(t){return this.x-=t,this.y-=t,this}subVectors(t,e){return this.x=t.x-e.x,this.y=t.y-e.y,this}multiply(t){return this.x*=t.x,this.y*=t.y,this}multiplyScalar(t){return this.x*=t,this.y*=t,this}divide(t){return this.x/=t.x,this.y/=t.y,this}divideScalar(t){return this.multiplyScalar(1/t)}applyMatrix3(t){let e=this.x,n=this.y,i=t.elements;return this.x=i[0]*e+i[3]*n+i[6],this.y=i[1]*e+i[4]*n+i[7],this}min(t){return this.x=Math.min(this.x,t.x),this.y=Math.min(this.y,t.y),this}max(t){return this.x=Math.max(this.x,t.x),this.y=Math.max(this.y,t.y),this}clamp(t,e){return this.x=Zt(this.x,t.x,e.x),this.y=Zt(this.y,t.y,e.y),this}clampScalar(t,e){return this.x=Zt(this.x,t,e),this.y=Zt(this.y,t,e),this}clampLength(t,e){let n=this.length();return this.divideScalar(n||1).multiplyScalar(Zt(n,t,e))}floor(){return this.x=Math.floor(this.x),this.y=Math.floor(this.y),this}ceil(){return this.x=Math.ceil(this.x),this.y=Math.ceil(this.y),this}round(){return this.x=Math.round(this.x),this.y=Math.round(this.y),this}roundToZero(){return this.x=Math.trunc(this.x),this.y=Math.trunc(this.y),this}negate(){return this.x=-this.x,this.y=-this.y,this}dot(t){return this.x*t.x+this.y*t.y}cross(t){return this.x*t.y-this.y*t.x}lengthSq(){return this.x*this.x+this.y*this.y}length(){return Math.sqrt(this.x*this.x+this.y*this.y)}manhattanLength(){return Math.abs(this.x)+Math.abs(this.y)}normalize(){return this.divideScalar(this.length()||1)}angle(){return Math.atan2(-this.y,-this.x)+Math.PI}angleTo(t){let e=Math.sqrt(this.lengthSq()*t.lengthSq());if(e===0)return Math.PI/2;let n=this.dot(t)/e;return Math.acos(Zt(n,-1,1))}distanceTo(t){return Math.sqrt(this.distanceToSquared(t))}distanceToSquared(t){let e=this.x-t.x,n=this.y-t.y;return e*e+n*n}manhattanDistanceTo(t){return Math.abs(this.x-t.x)+Math.abs(this.y-t.y)}setLength(t){return this.normalize().multiplyScalar(t)}lerp(t,e){return this.x+=(t.x-this.x)*e,this.y+=(t.y-this.y)*e,this}lerpVectors(t,e,n){return this.x=t.x+(e.x-t.x)*n,this.y=t.y+(e.y-t.y)*n,this}equals(t){return t.x===this.x&&t.y===this.y}fromArray(t,e=0){return this.x=t[e],this.y=t[e+1],this}toArray(t=[],e=0){return t[e]=this.x,t[e+1]=this.y,t}fromBufferAttribute(t,e){return this.x=t.getX(e),this.y=t.getY(e),this}rotateAround(t,e){let n=Math.cos(e),i=Math.sin(e),s=this.x-t.x,r=this.y-t.y;return this.x=s*n-r*i+t.x,this.y=s*i+r*n+t.y,this}random(){return this.x=Math.random(),this.y=Math.random(),this}*[Symbol.iterator](){yield this.x,yield this.y}}class Ce{constructor(t=0,e=0,n=0,i=1){this.isQuaternion=!0,this._x=t,this._y=e,this._z=n,this._w=i}static slerpFlat(t,e,n,i,s,r,a){let o=n[i+0],l=n[i+1],c=n[i+2],h=n[i+3],d=s[r+0],u=s[r+1],f=s[r+2],g=s[r+3];if(h!==g||o!==d||l!==u||c!==f){let _=o*d+l*u+c*f+h*g;if(_<0)d=-d,u=-u,f=-f,g=-g,_=-_;let m=1-a;if(_<0.9995){let p=Math.acos(_),b=Math.sin(p);m=Math.sin(m*p)/b,a=Math.sin(a*p)/b,o=o*m+d*a,l=l*m+u*a,c=c*m+f*a,h=h*m+g*a}else{o=o*m+d*a,l=l*m+u*a,c=c*m+f*a,h=h*m+g*a;let p=1/Math.sqrt(o*o+l*l+c*c+h*h);o*=p,l*=p,c*=p,h*=p}}t[e]=o,t[e+1]=l,t[e+2]=c,t[e+3]=h}static multiplyQuaternionsFlat(t,e,n,i,s,r){let a=n[i],o=n[i+1],l=n[i+2],c=n[i+3],h=s[r],d=s[r+1],u=s[r+2],f=s[r+3];return t[e]=a*f+c*h+o*u-l*d,t[e+1]=o*f+c*d+l*h-a*u,t[e+2]=l*f+c*u+a*d-o*h,t[e+3]=c*f-a*h-o*d-l*u,t}get x(){return this._x}set x(t){this._x=t,this._onChangeCallback()}get y(){return this._y}set y(t){this._y=t,this._onChangeCallback()}get z(){return this._z}set z(t){this._z=t,this._onChangeCallback()}get w(){return this._w}set w(t){this._w=t,this._onChangeCallback()}set(t,e,n,i){return this._x=t,this._y=e,this._z=n,this._w=i,this._onChangeCallback(),this}clone(){return new this.constructor(this._x,this._y,this._z,this._w)}copy(t){return this._x=t.x,this._y=t.y,this._z=t.z,this._w=t.w,this._onChangeCallback(),this}setFromEuler(t,e=!0){let{_x:n,_y:i,_z:s,_order:r}=t,{cos:a,sin:o}=Math,l=a(n/2),c=a(i/2),h=a(s/2),d=o(n/2),u=o(i/2),f=o(s/2);switch(r){case"XYZ":this._x=d*c*h+l*u*f,this._y=l*u*h-d*c*f,this._z=l*c*f+d*u*h,this._w=l*c*h-d*u*f;break;case"YXZ":this._x=d*c*h+l*u*f,this._y=l*u*h-d*c*f,this._z=l*c*f-d*u*h,this._w=l*c*h+d*u*f;break;case"ZXY":this._x=d*c*h-l*u*f,this._y=l*u*h+d*c*f,this._z=l*c*f+d*u*h,this._w=l*c*h-d*u*f;break;case"ZYX":this._x=d*c*h-l*u*f,this._y=l*u*h+d*c*f,this._z=l*c*f-d*u*h,this._w=l*c*h+d*u*f;break;case"YZX":this._x=d*c*h+l*u*f,this._y=l*u*h+d*c*f,this._z=l*c*f-d*u*h,this._w=l*c*h-d*u*f;break;case"XZY":this._x=d*c*h-l*u*f,this._y=l*u*h-d*c*f,this._z=l*c*f+d*u*h,this._w=l*c*h+d*u*f;break;default:Mt("Quaternion: .setFromEuler() encountered an unknown order: "+r)}if(e===!0)this._onChangeCallback();return this}setFromAxisAngle(t,e){let n=e/2,i=Math.sin(n);return this._x=t.x*i,this._y=t.y*i,this._z=t.z*i,this._w=Math.cos(n),this._onChangeCallback(),this}setFromRotationMatrix(t){let e=t.elements,n=e[0],i=e[4],s=e[8],r=e[1],a=e[5],o=e[9],l=e[2],c=e[6],h=e[10],d=n+a+h;if(d>0){let u=0.5/Math.sqrt(d+1);this._w=0.25/u,this._x=(c-o)*u,this._y=(s-l)*u,this._z=(r-i)*u}else if(n>a&&n>h){let u=2*Math.sqrt(1+n-a-h);this._w=(c-o)/u,this._x=0.25*u,this._y=(i+r)/u,this._z=(s+l)/u}else if(a>h){let u=2*Math.sqrt(1+a-n-h);this._w=(s-l)/u,this._x=(i+r)/u,this._y=0.25*u,this._z=(o+c)/u}else{let u=2*Math.sqrt(1+h-n-a);this._w=(r-i)/u,this._x=(s+l)/u,this._y=(o+c)/u,this._z=0.25*u}return this._onChangeCallback(),this}setFromUnitVectors(t,e){let n=t.dot(e)+1;if(n<0.00000001)if(n=0,Math.abs(t.x)>Math.abs(t.z))this._x=-t.y,this._y=t.x,this._z=0,this._w=n;else this._x=0,this._y=-t.z,this._z=t.y,this._w=n;else this._x=t.y*e.z-t.z*e.y,this._y=t.z*e.x-t.x*e.z,this._z=t.x*e.y-t.y*e.x,this._w=n;return this.normalize()}angleTo(t){return 2*Math.acos(Math.abs(Zt(this.dot(t),-1,1)))}rotateTowards(t,e){let n=this.angleTo(t);if(n===0)return this;let i=Math.min(1,e/n);return this.slerp(t,i),this}identity(){return this.set(0,0,0,1)}invert(){return this.conjugate()}conjugate(){return this._x*=-1,this._y*=-1,this._z*=-1,this._onChangeCallback(),this}dot(t){return this._x*t._x+this._y*t._y+this._z*t._z+this._w*t._w}lengthSq(){return this._x*this._x+this._y*this._y+this._z*this._z+this._w*this._w}length(){return Math.sqrt(this._x*this._x+this._y*this._y+this._z*this._z+this._w*this._w)}normalize(){let t=this.length();if(t===0)this._x=0,this._y=0,this._z=0,this._w=1;else t=1/t,this._x=this._x*t,this._y=this._y*t,this._z=this._z*t,this._w=this._w*t;return this._onChangeCallback(),this}multiply(t){return this.multiplyQuaternions(this,t)}premultiply(t){return this.multiplyQuaternions(t,this)}multiplyQuaternions(t,e){let{_x:n,_y:i,_z:s,_w:r}=t,{_x:a,_y:o,_z:l,_w:c}=e;return this._x=n*c+r*a+i*l-s*o,this._y=i*c+r*o+s*a-n*l,this._z=s*c+r*l+n*o-i*a,this._w=r*c-n*a-i*o-s*l,this._onChangeCallback(),this}slerp(t,e){let{_x:n,_y:i,_z:s,_w:r}=t,a=this.dot(t);if(a<0)n=-n,i=-i,s=-s,r=-r,a=-a;let o=1-e;if(a<0.9995){let l=Math.acos(a),c=Math.sin(l);o=Math.sin(o*l)/c,e=Math.sin(e*l)/c,this._x=this._x*o+n*e,this._y=this._y*o+i*e,this._z=this._z*o+s*e,this._w=this._w*o+r*e,this._onChangeCallback()}else this._x=this._x*o+n*e,this._y=this._y*o+i*e,this._z=this._z*o+s*e,this._w=this._w*o+r*e,this.normalize();return this}slerpQuaternions(t,e,n){return this.copy(t).slerp(e,n)}random(){let t=2*Math.PI*Math.random(),e=2*Math.PI*Math.random(),n=Math.random(),i=Math.sqrt(1-n),s=Math.sqrt(n);return this.set(i*Math.sin(t),i*Math.cos(t),s*Math.sin(e),s*Math.cos(e))}equals(t){return t._x===this._x&&t._y===this._y&&t._z===this._z&&t._w===this._w}fromArray(t,e=0){return this._x=t[e],this._y=t[e+1],this._z=t[e+2],this._w=t[e+3],this._onChangeCallback(),this}toArray(t=[],e=0){return t[e]=this._x,t[e+1]=this._y,t[e+2]=this._z,t[e+3]=this._w,t}fromBufferAttribute(t,e){return this._x=t.getX(e),this._y=t.getY(e),this._z=t.getZ(e),this._w=t.getW(e),this._onChangeCallback(),this}toJSON(){return this.toArray()}_onChange(t){return this._onChangeCallback=t,this}_onChangeCallback(){}*[Symbol.iterator](){yield this._x,yield this._y,yield this._z,yield this._w}}class L{static{L.prototype.isVector3=!0}constructor(t=0,e=0,n=0){this.x=t,this.y=e,this.z=n}set(t,e,n){if(n===void 0)n=this.z;return this.x=t,this.y=e,this.z=n,this}setScalar(t){return this.x=t,this.y=t,this.z=t,this}setX(t){return this.x=t,this}setY(t){return this.y=t,this}setZ(t){return this.z=t,this}setComponent(t,e){switch(t){case 0:this.x=e;break;case 1:this.y=e;break;case 2:this.z=e;break;default:throw Error("THREE.Vector3: index is out of range: "+t)}return this}getComponent(t){switch(t){case 0:return this.x;case 1:return this.y;case 2:return this.z;default:throw Error("THREE.Vector3: index is out of range: "+t)}}clone(){return new this.constructor(this.x,this.y,this.z)}copy(t){return this.x=t.x,this.y=t.y,this.z=t.z,this}add(t){return this.x+=t.x,this.y+=t.y,this.z+=t.z,this}addScalar(t){return this.x+=t,this.y+=t,this.z+=t,this}addVectors(t,e){return this.x=t.x+e.x,this.y=t.y+e.y,this.z=t.z+e.z,this}addScaledVector(t,e){return this.x+=t.x*e,this.y+=t.y*e,this.z+=t.z*e,this}sub(t){return this.x-=t.x,this.y-=t.y,this.z-=t.z,this}subScalar(t){return this.x-=t,this.y-=t,this.z-=t,this}subVectors(t,e){return this.x=t.x-e.x,this.y=t.y-e.y,this.z=t.z-e.z,this}multiply(t){return this.x*=t.x,this.y*=t.y,this.z*=t.z,this}multiplyScalar(t){return this.x*=t,this.y*=t,this.z*=t,this}multiplyVectors(t,e){return this.x=t.x*e.x,this.y=t.y*e.y,this.z=t.z*e.z,this}applyEuler(t){return this.applyQuaternion(qh.setFromEuler(t))}applyAxisAngle(t,e){return this.applyQuaternion(qh.setFromAxisAngle(t,e))}applyMatrix3(t){let e=this.x,n=this.y,i=this.z,s=t.elements;return this.x=s[0]*e+s[3]*n+s[6]*i,this.y=s[1]*e+s[4]*n+s[7]*i,this.z=s[2]*e+s[5]*n+s[8]*i,this}applyNormalMatrix(t){return this.applyMatrix3(t).normalize()}applyMatrix4(t){let e=this.x,n=this.y,i=this.z,s=t.elements,r=1/(s[3]*e+s[7]*n+s[11]*i+s[15]);return this.x=(s[0]*e+s[4]*n+s[8]*i+s[12])*r,this.y=(s[1]*e+s[5]*n+s[9]*i+s[13])*r,this.z=(s[2]*e+s[6]*n+s[10]*i+s[14])*r,this}applyQuaternion(t){let e=this.x,n=this.y,i=this.z,{x:s,y:r,z:a,w:o}=t,l=2*(r*i-a*n),c=2*(a*e-s*i),h=2*(s*n-r*e);return this.x=e+o*l+r*h-a*c,this.y=n+o*c+a*l-s*h,this.z=i+o*h+s*c-r*l,this}project(t){return this.applyMatrix4(t.matrixWorldInverse).applyMatrix4(t.projectionMatrix)}unproject(t){return this.applyMatrix4(t.projectionMatrixInverse).applyMatrix4(t.matrixWorld)}transformDirection(t){let e=this.x,n=this.y,i=this.z,s=t.elements;return this.x=s[0]*e+s[4]*n+s[8]*i,this.y=s[1]*e+s[5]*n+s[9]*i,this.z=s[2]*e+s[6]*n+s[10]*i,this.normalize()}divide(t){return this.x/=t.x,this.y/=t.y,this.z/=t.z,this}divideScalar(t){return this.multiplyScalar(1/t)}min(t){return this.x=Math.min(this.x,t.x),this.y=Math.min(this.y,t.y),this.z=Math.min(this.z,t.z),this}max(t){return this.x=Math.max(this.x,t.x),this.y=Math.max(this.y,t.y),this.z=Math.max(this.z,t.z),this}clamp(t,e){return this.x=Zt(this.x,t.x,e.x),this.y=Zt(this.y,t.y,e.y),this.z=Zt(this.z,t.z,e.z),this}clampScalar(t,e){return this.x=Zt(this.x,t,e),this.y=Zt(this.y,t,e),this.z=Zt(this.z,t,e),this}clampLength(t,e){let n=this.length();return this.divideScalar(n||1).multiplyScalar(Zt(n,t,e))}floor(){return this.x=Math.floor(this.x),this.y=Math.floor(this.y),this.z=Math.floor(this.z),this}ceil(){return this.x=Math.ceil(this.x),this.y=Math.ceil(this.y),this.z=Math.ceil(this.z),this}round(){return this.x=Math.round(this.x),this.y=Math.round(this.y),this.z=Math.round(this.z),this}roundToZero(){return this.x=Math.trunc(this.x),this.y=Math.trunc(this.y),this.z=Math.trunc(this.z),this}negate(){return this.x=-this.x,this.y=-this.y,this.z=-this.z,this}dot(t){return this.x*t.x+this.y*t.y+this.z*t.z}lengthSq(){return this.x*this.x+this.y*this.y+this.z*this.z}length(){return Math.sqrt(this.x*this.x+this.y*this.y+this.z*this.z)}manhattanLength(){return Math.abs(this.x)+Math.abs(this.y)+Math.abs(this.z)}normalize(){return this.divideScalar(this.length()||1)}setLength(t){return this.normalize().multiplyScalar(t)}lerp(t,e){return this.x+=(t.x-this.x)*e,this.y+=(t.y-this.y)*e,this.z+=(t.z-this.z)*e,this}lerpVectors(t,e,n){return this.x=t.x+(e.x-t.x)*n,this.y=t.y+(e.y-t.y)*n,this.z=t.z+(e.z-t.z)*n,this}cross(t){return this.crossVectors(this,t)}crossVectors(t,e){let{x:n,y:i,z:s}=t,{x:r,y:a,z:o}=e;return this.x=i*o-s*a,this.y=s*r-n*o,this.z=n*a-i*r,this}projectOnVector(t){let e=t.lengthSq();if(e===0)return this.set(0,0,0);let n=t.dot(this)/e;return this.copy(t).multiplyScalar(n)}projectOnPlane(t){return zo.copy(this).projectOnVector(t),this.sub(zo)}reflect(t){return this.sub(zo.copy(t).multiplyScalar(2*this.dot(t)))}angleTo(t){let e=Math.sqrt(this.lengthSq()*t.lengthSq());if(e===0)return Math.PI/2;let n=this.dot(t)/e;return Math.acos(Zt(n,-1,1))}distanceTo(t){return Math.sqrt(this.distanceToSquared(t))}distanceToSquared(t){let e=this.x-t.x,n=this.y-t.y,i=this.z-t.z;return e*e+n*n+i*i}manhattanDistanceTo(t){return Math.abs(this.x-t.x)+Math.abs(this.y-t.y)+Math.abs(this.z-t.z)}setFromSpherical(t){return this.setFromSphericalCoords(t.radius,t.phi,t.theta)}setFromSphericalCoords(t,e,n){let i=Math.sin(e)*t;return this.x=i*Math.sin(n),this.y=Math.cos(e)*t,this.z=i*Math.cos(n),this}setFromCylindrical(t){return this.setFromCylindricalCoords(t.radius,t.theta,t.y)}setFromCylindricalCoords(t,e,n){return this.x=t*Math.sin(e),this.y=n,this.z=t*Math.cos(e),this}setFromMatrixPosition(t){let e=t.elements;return this.x=e[12],this.y=e[13],this.z=e[14],this}setFromMatrixScale(t){let e=this.setFromMatrixColumn(t,0).length(),n=this.setFromMatrixColumn(t,1).length(),i=this.setFromMatrixColumn(t,2).length();return this.x=e,this.y=n,this.z=i,this}setFromMatrixColumn(t,e){return this.fromArray(t.elements,e*4)}setFromMatrix3Column(t,e){return this.fromArray(t.elements,e*3)}setFromEuler(t){return this.x=t._x,this.y=t._y,this.z=t._z,this}setFromColor(t){return this.x=t.r,this.y=t.g,this.z=t.b,this}equals(t){return t.x===this.x&&t.y===this.y&&t.z===this.z}fromArray(t,e=0){return this.x=t[e],this.y=t[e+1],this.z=t[e+2],this}toArray(t=[],e=0){return t[e]=this.x,t[e+1]=this.y,t[e+2]=this.z,t}fromBufferAttribute(t,e){return this.x=t.getX(e),this.y=t.getY(e),this.z=t.getZ(e),this}random(){return this.x=Math.random(),this.y=Math.random(),this.z=Math.random(),this}randomDirection(){let t=Math.random()*Math.PI*2,e=Math.random()*2-1,n=Math.sqrt(1-e*e);return this.x=n*Math.cos(t),this.y=e,this.z=n*Math.sin(t),this}*[Symbol.iterator](){yield this.x,yield this.y,yield this.z}}var zo=new L,qh=new Ce;class Xt{static{Xt.prototype.isMatrix3=!0}constructor(t,e,n,i,s,r,a,o,l){if(this.elements=[1,0,0,0,1,0,0,0,1],t!==void 0)this.set(t,e,n,i,s,r,a,o,l)}set(t,e,n,i,s,r,a,o,l){let c=this.elements;return c[0]=t,c[1]=i,c[2]=a,c[3]=e,c[4]=s,c[5]=o,c[6]=n,c[7]=r,c[8]=l,this}identity(){return this.set(1,0,0,0,1,0,0,0,1),this}copy(t){let e=this.elements,n=t.elements;return e[0]=n[0],e[1]=n[1],e[2]=n[2],e[3]=n[3],e[4]=n[4],e[5]=n[5],e[6]=n[6],e[7]=n[7],e[8]=n[8],this}extractBasis(t,e,n){return t.setFromMatrix3Column(this,0),e.setFromMatrix3Column(this,1),n.setFromMatrix3Column(this,2),this}setFromMatrix4(t){let e=t.elements;return this.set(e[0],e[4],e[8],e[1],e[5],e[9],e[2],e[6],e[10]),this}multiply(t){return this.multiplyMatrices(this,t)}premultiply(t){return this.multiplyMatrices(t,this)}multiplyMatrices(t,e){let n=t.elements,i=e.elements,s=this.elements,r=n[0],a=n[3],o=n[6],l=n[1],c=n[4],h=n[7],d=n[2],u=n[5],f=n[8],g=i[0],_=i[3],m=i[6],p=i[1],b=i[4],S=i[7],x=i[2],T=i[5],I=i[8];return s[0]=r*g+a*p+o*x,s[3]=r*_+a*b+o*T,s[6]=r*m+a*S+o*I,s[1]=l*g+c*p+h*x,s[4]=l*_+c*b+h*T,s[7]=l*m+c*S+h*I,s[2]=d*g+u*p+f*x,s[5]=d*_+u*b+f*T,s[8]=d*m+u*S+f*I,this}multiplyScalar(t){let e=this.elements;return e[0]*=t,e[3]*=t,e[6]*=t,e[1]*=t,e[4]*=t,e[7]*=t,e[2]*=t,e[5]*=t,e[8]*=t,this}determinant(){let t=this.elements,e=t[0],n=t[1],i=t[2],s=t[3],r=t[4],a=t[5],o=t[6],l=t[7],c=t[8];return e*r*c-e*a*l-n*s*c+n*a*o+i*s*l-i*r*o}invert(){let t=this.elements,e=t[0],n=t[1],i=t[2],s=t[3],r=t[4],a=t[5],o=t[6],l=t[7],c=t[8],h=c*r-a*l,d=a*o-c*s,u=l*s-r*o,f=e*h+n*d+i*u;if(f===0)return this.set(0,0,0,0,0,0,0,0,0);let g=1/f;return t[0]=h*g,t[1]=(i*l-c*n)*g,t[2]=(a*n-i*r)*g,t[3]=d*g,t[4]=(c*e-i*o)*g,t[5]=(i*s-a*e)*g,t[6]=u*g,t[7]=(n*o-l*e)*g,t[8]=(r*e-n*s)*g,this}transpose(){let t,e=this.elements;return t=e[1],e[1]=e[3],e[3]=t,t=e[2],e[2]=e[6],e[6]=t,t=e[5],e[5]=e[7],e[7]=t,this}getNormalMatrix(t){return this.setFromMatrix4(t).invert().transpose()}transposeIntoArray(t){let e=this.elements;return t[0]=e[0],t[1]=e[3],t[2]=e[6],t[3]=e[1],t[4]=e[4],t[5]=e[7],t[6]=e[2],t[7]=e[5],t[8]=e[8],this}setUvTransform(t,e,n,i,s,r,a){let o=Math.cos(s),l=Math.sin(s);return this.set(n*o,n*l,-n*(o*r+l*a)+r+t,-i*l,i*o,-i*(-l*r+o*a)+a+e,0,0,1),this}scale(t,e){return qn("Matrix3: .scale() is deprecated. Use .makeScale() instead."),this.premultiply(ko.makeScale(t,e)),this}rotate(t){return qn("Matrix3: .rotate() is deprecated. Use .makeRotation() instead."),this.premultiply(ko.makeRotation(-t)),this}translate(t,e){return qn("Matrix3: .translate() is deprecated. Use .makeTranslation() instead."),this.premultiply(ko.makeTranslation(t,e)),this}makeTranslation(t,e){if(t.isVector2)this.set(1,0,t.x,0,1,t.y,0,0,1);else this.set(1,0,t,0,1,e,0,0,1);return this}makeRotation(t){let e=Math.cos(t),n=Math.sin(t);return this.set(e,-n,0,n,e,0,0,0,1),this}makeScale(t,e){return this.set(t,0,0,0,e,0,0,0,1),this}equals(t){let e=this.elements,n=t.elements;for(let i=0;i<9;i++)if(e[i]!==n[i])return!1;return!0}fromArray(t,e=0){for(let n=0;n<9;n++)this.elements[n]=t[n+e];return this}toArray(t=[],e=0){let n=this.elements;return t[e]=n[0],t[e+1]=n[1],t[e+2]=n[2],t[e+3]=n[3],t[e+4]=n[4],t[e+5]=n[5],t[e+6]=n[6],t[e+7]=n[7],t[e+8]=n[8],t}clone(){return new this.constructor().fromArray(this.elements)}}var ko=new Xt,Yh=new Xt().set(0.4123908,0.3575843,0.1804808,0.212639,0.7151687,0.0721923,0.0193308,0.1191948,0.9505322),Zh=new Xt().set(3.2409699,-1.5373832,-0.4986108,-0.9692436,1.8759675,0.0415551,0.0556301,-0.203977,1.0569715);function Xg(){let t={enabled:!0,workingColorSpace:"srgb-linear",spaces:{},convert:function(s,r,a){if(this.enabled===!1||r===a||!r||!a)return s;if(this.spaces[r].transfer==="srgb")s.r=Yn(s.r),s.g=Yn(s.g),s.b=Yn(s.b);if(this.spaces[r].primaries!==this.spaces[a].primaries)s.applyMatrix3(this.spaces[r].toXYZ),s.applyMatrix3(this.spaces[a].fromXYZ);if(this.spaces[a].transfer==="srgb")s.r=gs(s.r),s.g=gs(s.g),s.b=gs(s.b);return s},workingToColorSpace:function(s,r){return this.convert(s,this.workingColorSpace,r)},colorSpaceToWorking:function(s,r){return this.convert(s,r,this.workingColorSpace)},getPrimaries:function(s){return this.spaces[s].primaries},getTransfer:function(s){if(s==="")return"linear";return this.spaces[s].transfer},getToneMappingMode:function(s){return this.spaces[s].outputColorSpaceConfig.toneMappingMode||"standard"},getLuminanceCoefficients:function(s,r=this.workingColorSpace){return s.fromArray(this.spaces[r].luminanceCoefficients)},define:function(s){Object.assign(this.spaces,s)},_getMatrix:function(s,r,a){return s.copy(this.spaces[r].toXYZ).multiply(this.spaces[a].fromXYZ)},_getDrawingBufferColorSpace:function(s){return this.spaces[s].outputColorSpaceConfig.drawingBufferColorSpace},_getUnpackColorSpace:function(s=this.workingColorSpace){return this.spaces[s].workingColorSpaceConfig.unpackColorSpace},fromWorkingColorSpace:function(s,r){return qn("ColorManagement: .fromWorkingColorSpace() has been renamed to .workingToColorSpace()."),t.workingToColorSpace(s,r)},toWorkingColorSpace:function(s,r){return qn("ColorManagement: .toWorkingColorSpace() has been renamed to .colorSpaceToWorking()."),t.colorSpaceToWorking(s,r)}},e=[0.64,0.33,0.3,0.6,0.15,0.06],n=[0.2126,0.7152,0.0722],i=[0.3127,0.329];return t.define({["srgb-linear"]:{primaries:e,whitePoint:i,transfer:"linear",toXYZ:Yh,fromXYZ:Zh,luminanceCoefficients:n,workingColorSpaceConfig:{unpackColorSpace:"srgb"},outputColorSpaceConfig:{drawingBufferColorSpace:"srgb"}},["srgb"]:{primaries:e,whitePoint:i,transfer:"srgb",toXYZ:Yh,fromXYZ:Zh,luminanceCoefficients:n,outputColorSpaceConfig:{drawingBufferColorSpace:"srgb"}}}),t}var re=Xg();function Yn(t){return t<0.04045?t*0.0773993808:Math.pow(t*0.9478672986+0.0521327014,2.4)}function gs(t){return t<0.0031308?t*12.92:1.055*Math.pow(t,0.41666)-0.055}var Ji;class wc{static getDataURL(t,e="image/png"){if(/^data:/i.test(t.src))return t.src;if(typeof HTMLCanvasElement>"u")return t.src;let n;if(t instanceof HTMLCanvasElement)n=t;else{if(Ji===void 0)Ji=_s("canvas");Ji.width=t.width,Ji.height=t.height;let i=Ji.getContext("2d");if(t instanceof ImageData)i.putImageData(t,0,0);else i.drawImage(t,0,0,t.width,t.height);n=Ji}return n.toDataURL(e)}static sRGBToLinear(t){if(typeof HTMLImageElement<"u"&&t instanceof HTMLImageElement||typeof HTMLCanvasElement<"u"&&t instanceof HTMLCanvasElement||typeof ImageBitmap<"u"&&t instanceof ImageBitmap){let e=_s("canvas");e.width=t.width,e.height=t.height;let n=e.getContext("2d");n.drawImage(t,0,0,t.width,t.height);let i=n.getImageData(0,0,t.width,t.height),s=i.data;for(let r=0;r<s.length;r++)s[r]=Yn(s[r]/255)*255;return n.putImageData(i,0,0),e}else if(t.data){let e=t.data.slice(0);for(let n=0;n<e.length;n++)if(e instanceof Uint8Array||e instanceof Uint8ClampedArray)e[n]=Math.floor(Yn(e[n]/255)*255);else e[n]=Yn(e[n]);return{data:e,width:t.width,height:t.height}}else return Mt("ImageUtils.sRGBToLinear(): Unsupported image type. No color space conversion applied."),t}}var qg=0;class In{constructor(t=null){this.isTextureSource=!0,Object.defineProperty(this,"id",{value:qg++}),this.uuid=rn(),this.data=t,this.dataReady=!0,this.version=0}getSize(t){let e=this.data;if(typeof HTMLVideoElement<"u"&&e instanceof HTMLVideoElement)t.set(e.videoWidth,e.videoHeight,0);else if(typeof VideoFrame<"u"&&e instanceof VideoFrame)t.set(e.displayWidth,e.displayHeight,0);else if(e!==null)t.set(e.width,e.height,e.depth||0);else t.set(0,0,0);return t}set needsUpdate(t){if(t===!0)this.version++}toJSON(t){let e=t===void 0||typeof t==="string";if(!e&&t.images[this.uuid]!==void 0)return t.images[this.uuid];let n={uuid:this.uuid,url:""},i=this.data;if(i!==null){let s;if(Array.isArray(i)){s=[];for(let r=0,a=i.length;r<a;r++)if(i[r].isDataTexture)s.push(Go(i[r].image));else s.push(Go(i[r]))}else s=Go(i);n.url=s}if(!e)t.images[this.uuid]=n;return n}}function Go(t){if(typeof HTMLImageElement<"u"&&t instanceof HTMLImageElement||typeof HTMLCanvasElement<"u"&&t instanceof HTMLCanvasElement||typeof ImageBitmap<"u"&&t instanceof ImageBitmap)return wc.getDataURL(t);else if(t.data)return{data:Array.from(t.data),width:t.width,height:t.height,type:t.data.constructor.name};else return Mt("Texture: Unable to serialize Texture."),{}}class rf extends In{constructor(t=null){qn('Source: "Source" has been renamed to "TextureSource". Please update your code to use "THREE.TextureSource" instead.');super(t);this.isSource=!0}}var Yg=0,Ho=new L;class be extends pn{constructor(t=be.DEFAULT_IMAGE,e=be.DEFAULT_MAPPING,n=1001,i=1001,s=1006,r=1008,a=1023,o=1009,l=be.DEFAULT_ANISOTROPY,c=""){super();this.isTexture=!0,Object.defineProperty(this,"id",{value:Yg++}),this.uuid=rn(),this.name="",this.source=new In(t),this.mipmaps=[],this.mapping=e,this.channel=0,this.wrapS=n,this.wrapT=i,this.magFilter=s,this.minFilter=r,this.anisotropy=l,this.format=a,this.internalFormat=null,this.type=o,this.offset=new rt(0,0),this.repeat=new rt(1,1),this.center=new rt(0,0),this.rotation=0,this.matrixAutoUpdate=!0,this.matrix=new Xt,this.generateMipmaps=!0,this.premultiplyAlpha=!1,this.flipY=!0,this.unpackAlignment=4,this.colorSpace=c,this.userData={},this.updateRanges=[],this.version=0,this.onUpdate=null,this.renderTarget=null,this.isRenderTargetTexture=!1,this.isArrayTexture=t&&t.depth&&t.depth>1?!0:!1,this.pmremVersion=0,this.normalized=!1}get width(){return this.source.getSize(Ho).x}get height(){return this.source.getSize(Ho).y}get depth(){return this.source.getSize(Ho).z}get image(){return this.source.data}set image(t){this.source.data=t}updateMatrix(){this.matrix.setUvTransform(this.offset.x,this.offset.y,this.repeat.x,this.repeat.y,this.rotation,this.center.x,this.center.y)}addUpdateRange(t,e){this.updateRanges.push({start:t,count:e})}clearUpdateRanges(){this.updateRanges.length=0}clone(){return new this.constructor().copy(this)}copy(t){return this.name=t.name,this.source=t.source,this.mipmaps=t.mipmaps.slice(0),this.mapping=t.mapping,this.channel=t.channel,this.wrapS=t.wrapS,this.wrapT=t.wrapT,this.magFilter=t.magFilter,this.minFilter=t.minFilter,this.anisotropy=t.anisotropy,this.format=t.format,this.internalFormat=t.internalFormat,this.type=t.type,this.normalized=t.normalized,this.offset.copy(t.offset),this.repeat.copy(t.repeat),this.center.copy(t.center),this.rotation=t.rotation,this.matrixAutoUpdate=t.matrixAutoUpdate,this.matrix.copy(t.matrix),this.generateMipmaps=t.generateMipmaps,this.premultiplyAlpha=t.premultiplyAlpha,this.flipY=t.flipY,this.unpackAlignment=t.unpackAlignment,this.colorSpace=t.colorSpace,this.renderTarget=t.renderTarget,this.isRenderTargetTexture=t.isRenderTargetTexture,this.isArrayTexture=t.isArrayTexture,this.userData=JSON.parse(JSON.stringify(t.userData)),this.needsUpdate=!0,this}setValues(t){for(let e in t){let n=t[e];if(n===void 0){Mt(`Texture.setValues(): parameter '${e}' has value of undefined.`);continue}let i=this[e];if(i===void 0){Mt(`Texture.setValues(): property '${e}' does not exist.`);continue}if(i&&n&&(i.isVector2&&n.isVector2))i.copy(n);else if(i&&n&&(i.isVector3&&n.isVector3))i.copy(n);else if(i&&n&&(i.isMatrix3&&n.isMatrix3))i.copy(n);else this[e]=n}}toJSON(t){let e=t===void 0||typeof t==="string";if(!e&&t.textures[this.uuid]!==void 0)return t.textures[this.uuid];let n={metadata:{version:4.7,type:"Texture",generator:"Texture.toJSON"},uuid:this.uuid,name:this.name,image:this.source.toJSON(t).uuid,mapping:this.mapping,channel:this.channel,repeat:[this.repeat.x,this.repeat.y],offset:[this.offset.x,this.offset.y],center:[this.center.x,this.center.y],rotation:this.rotation,wrap:[this.wrapS,this.wrapT],format:this.format,internalFormat:this.internalFormat,type:this.type,normalized:this.normalized,colorSpace:this.colorSpace,minFilter:this.minFilter,magFilter:this.magFilter,anisotropy:this.anisotropy,flipY:this.flipY,generateMipmaps:this.generateMipmaps,premultiplyAlpha:this.premultiplyAlpha,unpackAlignment:this.unpackAlignment};if(Object.keys(this.userData).length>0)n.userData=this.userData;if(!e)t.textures[this.uuid]=n;return n}dispose(){this.dispatchEvent({type:"dispose"})}transformUv(t){if(this.mapping!==300)return t;if(t.applyMatrix3(this.matrix),t.x<0||t.x>1)switch(this.wrapS){case 1000:t.x=t.x-Math.floor(t.x);break;case 1001:t.x=t.x<0?0:1;break;case 1002:if(Math.abs(Math.floor(t.x)%2)===1)t.x=Math.ceil(t.x)-t.x;else t.x=t.x-Math.floor(t.x);break}if(t.y<0||t.y>1)switch(this.wrapT){case 1000:t.y=t.y-Math.floor(t.y);break;case 1001:t.y=t.y<0?0:1;break;case 1002:if(Math.abs(Math.floor(t.y)%2)===1)t.y=Math.ceil(t.y)-t.y;else t.y=t.y-Math.floor(t.y);break}if(this.flipY)t.y=1-t.y;return t}set needsUpdate(t){if(t===!0)this.version++,this.source.needsUpdate=!0}set needsPMREMUpdate(t){if(t===!0)this.pmremVersion++}}be.DEFAULT_IMAGE=null;be.DEFAULT_MAPPING=300;be.DEFAULT_ANISOTROPY=1;class fe{static{fe.prototype.isVector4=!0}constructor(t=0,e=0,n=0,i=1){this.x=t,this.y=e,this.z=n,this.w=i}get width(){return this.z}set width(t){this.z=t}get height(){return this.w}set height(t){this.w=t}set(t,e,n,i){return this.x=t,this.y=e,this.z=n,this.w=i,this}setScalar(t){return this.x=t,this.y=t,this.z=t,this.w=t,this}setX(t){return this.x=t,this}setY(t){return this.y=t,this}setZ(t){return this.z=t,this}setW(t){return this.w=t,this}setComponent(t,e){switch(t){case 0:this.x=e;break;case 1:this.y=e;break;case 2:this.z=e;break;case 3:this.w=e;break;default:throw Error("THREE.Vector4: index is out of range: "+t)}return this}getComponent(t){switch(t){case 0:return this.x;case 1:return this.y;case 2:return this.z;case 3:return this.w;default:throw Error("THREE.Vector4: index is out of range: "+t)}}clone(){return new this.constructor(this.x,this.y,this.z,this.w)}copy(t){return this.x=t.x,this.y=t.y,this.z=t.z,this.w=t.w!==void 0?t.w:1,this}add(t){return this.x+=t.x,this.y+=t.y,this.z+=t.z,this.w+=t.w,this}addScalar(t){return this.x+=t,this.y+=t,this.z+=t,this.w+=t,this}addVectors(t,e){return this.x=t.x+e.x,this.y=t.y+e.y,this.z=t.z+e.z,this.w=t.w+e.w,this}addScaledVector(t,e){return this.x+=t.x*e,this.y+=t.y*e,this.z+=t.z*e,this.w+=t.w*e,this}sub(t){return this.x-=t.x,this.y-=t.y,this.z-=t.z,this.w-=t.w,this}subScalar(t){return this.x-=t,this.y-=t,this.z-=t,this.w-=t,this}subVectors(t,e){return this.x=t.x-e.x,this.y=t.y-e.y,this.z=t.z-e.z,this.w=t.w-e.w,this}multiply(t){return this.x*=t.x,this.y*=t.y,this.z*=t.z,this.w*=t.w,this}multiplyScalar(t){return this.x*=t,this.y*=t,this.z*=t,this.w*=t,this}applyMatrix4(t){let e=this.x,n=this.y,i=this.z,s=this.w,r=t.elements;return this.x=r[0]*e+r[4]*n+r[8]*i+r[12]*s,this.y=r[1]*e+r[5]*n+r[9]*i+r[13]*s,this.z=r[2]*e+r[6]*n+r[10]*i+r[14]*s,this.w=r[3]*e+r[7]*n+r[11]*i+r[15]*s,this}divide(t){return this.x/=t.x,this.y/=t.y,this.z/=t.z,this.w/=t.w,this}divideScalar(t){return this.multiplyScalar(1/t)}setAxisAngleFromQuaternion(t){this.w=2*Math.acos(t.w);let e=Math.sqrt(1-t.w*t.w);if(e<0.0001)this.x=1,this.y=0,this.z=0;else this.x=t.x/e,this.y=t.y/e,this.z=t.z/e;return this}setAxisAngleFromRotationMatrix(t){let e,n,i,s,r=0.01,a=0.1,o=t.elements,l=o[0],c=o[4],h=o[8],d=o[1],u=o[5],f=o[9],g=o[2],_=o[6],m=o[10];if(Math.abs(c-d)<0.01&&Math.abs(h-g)<0.01&&Math.abs(f-_)<0.01){if(Math.abs(c+d)<0.1&&Math.abs(h+g)<0.1&&Math.abs(f+_)<0.1&&Math.abs(l+u+m-3)<0.1)return this.set(1,0,0,0),this;e=Math.PI;let b=(l+1)/2,S=(u+1)/2,x=(m+1)/2,T=(c+d)/4,I=(h+g)/4,R=(f+_)/4;if(b>S&&b>x)if(b<0.01)n=0,i=0.707106781,s=0.707106781;else n=Math.sqrt(b),i=T/n,s=I/n;else if(S>x)if(S<0.01)n=0.707106781,i=0,s=0.707106781;else i=Math.sqrt(S),n=T/i,s=R/i;else if(x<0.01)n=0.707106781,i=0.707106781,s=0;else s=Math.sqrt(x),n=I/s,i=R/s;return this.set(n,i,s,e),this}let p=Math.sqrt((_-f)*(_-f)+(h-g)*(h-g)+(d-c)*(d-c));if(Math.abs(p)<0.001)p=1;return this.x=(_-f)/p,this.y=(h-g)/p,this.z=(d-c)/p,this.w=Math.acos((l+u+m-1)/2),this}setFromMatrixPosition(t){let e=t.elements;return this.x=e[12],this.y=e[13],this.z=e[14],this.w=e[15],this}min(t){return this.x=Math.min(this.x,t.x),this.y=Math.min(this.y,t.y),this.z=Math.min(this.z,t.z),this.w=Math.min(this.w,t.w),this}max(t){return this.x=Math.max(this.x,t.x),this.y=Math.max(this.y,t.y),this.z=Math.max(this.z,t.z),this.w=Math.max(this.w,t.w),this}clamp(t,e){return this.x=Zt(this.x,t.x,e.x),this.y=Zt(this.y,t.y,e.y),this.z=Zt(this.z,t.z,e.z),this.w=Zt(this.w,t.w,e.w),this}clampScalar(t,e){return this.x=Zt(this.x,t,e),this.y=Zt(this.y,t,e),this.z=Zt(this.z,t,e),this.w=Zt(this.w,t,e),this}clampLength(t,e){let n=this.length();return this.divideScalar(n||1).multiplyScalar(Zt(n,t,e))}floor(){return this.x=Math.floor(this.x),this.y=Math.floor(this.y),this.z=Math.floor(this.z),this.w=Math.floor(this.w),this}ceil(){return this.x=Math.ceil(this.x),this.y=Math.ceil(this.y),this.z=Math.ceil(this.z),this.w=Math.ceil(this.w),this}round(){return this.x=Math.round(this.x),this.y=Math.round(this.y),this.z=Math.round(this.z),this.w=Math.round(this.w),this}roundToZero(){return this.x=Math.trunc(this.x),this.y=Math.trunc(this.y),this.z=Math.trunc(this.z),this.w=Math.trunc(this.w),this}negate(){return this.x=-this.x,this.y=-this.y,this.z=-this.z,this.w=-this.w,this}dot(t){return this.x*t.x+this.y*t.y+this.z*t.z+this.w*t.w}lengthSq(){return this.x*this.x+this.y*this.y+this.z*this.z+this.w*this.w}length(){return Math.sqrt(this.x*this.x+this.y*this.y+this.z*this.z+this.w*this.w)}manhattanLength(){return Math.abs(this.x)+Math.abs(this.y)+Math.abs(this.z)+Math.abs(this.w)}normalize(){return this.divideScalar(this.length()||1)}setLength(t){return this.normalize().multiplyScalar(t)}lerp(t,e){return this.x+=(t.x-this.x)*e,this.y+=(t.y-this.y)*e,this.z+=(t.z-this.z)*e,this.w+=(t.w-this.w)*e,this}lerpVectors(t,e,n){return this.x=t.x+(e.x-t.x)*n,this.y=t.y+(e.y-t.y)*n,this.z=t.z+(e.z-t.z)*n,this.w=t.w+(e.w-t.w)*n,this}equals(t){return t.x===this.x&&t.y===this.y&&t.z===this.z&&t.w===this.w}fromArray(t,e=0){return this.x=t[e],this.y=t[e+1],this.z=t[e+2],this.w=t[e+3],this}toArray(t=[],e=0){return t[e]=this.x,t[e+1]=this.y,t[e+2]=this.z,t[e+3]=this.w,t}fromBufferAttribute(t,e){return this.x=t.getX(e),this.y=t.getY(e),this.z=t.getZ(e),this.w=t.getW(e),this}random(){return this.x=Math.random(),this.y=Math.random(),this.z=Math.random(),this.w=Math.random(),this}*[Symbol.iterator](){yield this.x,yield this.y,yield this.z,yield this.w}}class Wa extends pn{constructor(t=1,e=1,n={}){super();n=Object.assign({generateMipmaps:!1,internalFormat:null,minFilter:1006,depthBuffer:!0,stencilBuffer:!1,resolveColorBuffer:!0,resolveDepthBuffer:!0,resolveStencilBuffer:!0,storeMultisampledColorBuffer:!0,storeMultisampledDepthBuffer:!0,storeMultisampledStencilBuffer:!0,depthTexture:null,samples:0,count:1,depth:1,multiview:!1,useArrayDepthTexture:!1},n),this.isRenderTarget=!0,this.width=t,this.height=e,this.depth=n.depth,this.scissor=new fe(0,0,t,e),this.scissorTest=!1,this.viewport=new fe(0,0,t,e),this.textures=[];let i={width:t,height:e,depth:n.depth},s=new be(i),r=n.count;for(let a=0;a<r;a++)this.textures[a]=s.clone(),this.textures[a].isRenderTargetTexture=!0,this.textures[a].renderTarget=this;this._setTextureOptions(n),this.depthBuffer=n.depthBuffer,this.stencilBuffer=n.stencilBuffer,this.resolveColorBuffer=n.resolveColorBuffer,this.resolveDepthBuffer=n.resolveDepthBuffer,this.resolveStencilBuffer=n.resolveStencilBuffer,this.storeMultisampledColorBuffer=n.storeMultisampledColorBuffer,this.storeMultisampledDepthBuffer=n.storeMultisampledDepthBuffer,this.storeMultisampledStencilBuffer=n.storeMultisampledStencilBuffer,this._depthTexture=null,this.depthTexture=n.depthTexture,this.samples=n.samples,this.multiview=n.multiview,this.useArrayDepthTexture=n.useArrayDepthTexture}_setTextureOptions(t={}){let e={minFilter:1006,generateMipmaps:!1,flipY:!1,internalFormat:null};if(t.mapping!==void 0)e.mapping=t.mapping;if(t.wrapS!==void 0)e.wrapS=t.wrapS;if(t.wrapT!==void 0)e.wrapT=t.wrapT;if(t.wrapR!==void 0)e.wrapR=t.wrapR;if(t.magFilter!==void 0)e.magFilter=t.magFilter;if(t.minFilter!==void 0)e.minFilter=t.minFilter;if(t.format!==void 0)e.format=t.format;if(t.type!==void 0)e.type=t.type;if(t.anisotropy!==void 0)e.anisotropy=t.anisotropy;if(t.colorSpace!==void 0)e.colorSpace=t.colorSpace;if(t.flipY!==void 0)e.flipY=t.flipY;if(t.generateMipmaps!==void 0)e.generateMipmaps=t.generateMipmaps;if(t.internalFormat!==void 0)e.internalFormat=t.internalFormat;for(let n=0;n<this.textures.length;n++)this.textures[n].setValues(e)}get texture(){return this.textures[0]}set texture(t){this.textures[0]=t}set depthTexture(t){if(this._depthTexture!==null&&this._depthTexture.renderTarget===this)this._depthTexture.renderTarget=null;if(t!==null&&t.renderTarget===null)t.renderTarget=this;this._depthTexture=t}get depthTexture(){return this._depthTexture}setSize(t,e,n=1){if(this.width!==t||this.height!==e||this.depth!==n){this.width=t,this.height=e,this.depth=n;for(let i=0,s=this.textures.length;i<s;i++)if(this.textures[i].image.width=t,this.textures[i].image.height=e,this.textures[i].image.depth=n,this.textures[i].isData3DTexture!==!0)this.textures[i].isArrayTexture=this.textures[i].image.depth>1;this.dispose()}this.viewport.set(0,0,t,e),this.scissor.set(0,0,t,e)}clone(){return new this.constructor().copy(this)}copy(t){this.width=t.width,this.height=t.height,this.depth=t.depth,this.scissor.copy(t.scissor),this.scissorTest=t.scissorTest,this.viewport.copy(t.viewport),this.textures.length=0;for(let e=0,n=t.textures.length;e<n;e++){this.textures[e]=t.textures[e].clone(),this.textures[e].isRenderTargetTexture=!0,this.textures[e].renderTarget=this;let i=Object.assign({},t.textures[e].image);this.textures[e].source=new In(i)}if(this.depthBuffer=t.depthBuffer,this.stencilBuffer=t.stencilBuffer,this.resolveColorBuffer=t.resolveColorBuffer,this.resolveDepthBuffer=t.resolveDepthBuffer,this.resolveStencilBuffer=t.resolveStencilBuffer,this.storeMultisampledColorBuffer=t.storeMultisampledColorBuffer,this.storeMultisampledDepthBuffer=t.storeMultisampledDepthBuffer,this.storeMultisampledStencilBuffer=t.storeMultisampledStencilBuffer,t.depthTexture!==null)if(t.depthTexture.renderTarget===t){let e=t.depthTexture.clone();e.renderTarget=null,this.depthTexture=e}else this.depthTexture=t.depthTexture;return this.samples=t.samples,this.multiview=t.multiview,this.useArrayDepthTexture=t.useArrayDepthTexture,this}dispose(){this.dispatchEvent({type:"dispose"})}}class Je extends Wa{constructor(t=1,e=1,n={}){super(t,e,n);this.isWebGLRenderTarget=!0}}class cr extends be{constructor(t=null,e=1,n=1,i=1){super(null);this.isDataArrayTexture=!0,this.image={data:t,width:e,height:n,depth:i},this.magFilter=1003,this.minFilter=1003,this.wrapR=1001,this.generateMipmaps=!1,this.flipY=!1,this.unpackAlignment=1,this.layerUpdates=new Set}copy(t){return super.copy(t),this.wrapR=t.wrapR,this}addLayerUpdate(t){this.layerUpdates.add(t)}clearLayerUpdates(){this.layerUpdates.clear()}}class af extends Je{constructor(t=1,e=1,n=1,i={}){super(t,e,i);this.isWebGLArrayRenderTarget=!0,this.depth=n,this.texture=new cr(null,t,e,n),this._setTextureOptions(i),this.texture.isRenderTargetTexture=!0}}class hr extends be{constructor(t=null,e=1,n=1,i=1){super(null);this.isData3DTexture=!0,this.image={data:t,width:e,height:n,depth:i},this.magFilter=1003,this.minFilter=1003,this.wrapR=1001,this.generateMipmaps=!1,this.flipY=!1,this.unpackAlignment=1}copy(t){return super.copy(t),this.wrapR=t.wrapR,this}}class of extends Je{constructor(t=1,e=1,n=1,i={}){super(t,e,i);this.isWebGL3DRenderTarget=!0,this.depth=n,this.texture=new hr(null,t,e,n),this._setTextureOptions(i),this.texture.isRenderTargetTexture=!0}}class Kt{static{Kt.prototype.isMatrix4=!0}constructor(t,e,n,i,s,r,a,o,l,c,h,d,u,f,g,_){if(this.elements=[1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1],t!==void 0)this.set(t,e,n,i,s,r,a,o,l,c,h,d,u,f,g,_)}set(t,e,n,i,s,r,a,o,l,c,h,d,u,f,g,_){let m=this.elements;return m[0]=t,m[4]=e,m[8]=n,m[12]=i,m[1]=s,m[5]=r,m[9]=a,m[13]=o,m[2]=l,m[6]=c,m[10]=h,m[14]=d,m[3]=u,m[7]=f,m[11]=g,m[15]=_,this}identity(){return this.set(1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1),this}clone(){return new Kt().fromArray(this.elements)}copy(t){let e=this.elements,n=t.elements;return e[0]=n[0],e[1]=n[1],e[2]=n[2],e[3]=n[3],e[4]=n[4],e[5]=n[5],e[6]=n[6],e[7]=n[7],e[8]=n[8],e[9]=n[9],e[10]=n[10],e[11]=n[11],e[12]=n[12],e[13]=n[13],e[14]=n[14],e[15]=n[15],this}copyPosition(t){let e=this.elements,n=t.elements;return e[12]=n[12],e[13]=n[13],e[14]=n[14],this}setFromMatrix3(t){let e=t.elements;return this.set(e[0],e[3],e[6],0,e[1],e[4],e[7],0,e[2],e[5],e[8],0,0,0,0,1),this}extractBasis(t,e,n){if(this.determinantAffine()===0)return t.set(1,0,0),e.set(0,1,0),n.set(0,0,1),this;return t.setFromMatrixColumn(this,0),e.setFromMatrixColumn(this,1),n.setFromMatrixColumn(this,2),this}makeBasis(t,e,n){return this.set(t.x,e.x,n.x,0,t.y,e.y,n.y,0,t.z,e.z,n.z,0,0,0,0,1),this}extractRotation(t){if(t.determinantAffine()===0)return this.identity();let e=this.elements,n=t.elements,i=1/$i.setFromMatrixColumn(t,0).length(),s=1/$i.setFromMatrixColumn(t,1).length(),r=1/$i.setFromMatrixColumn(t,2).length();return e[0]=n[0]*i,e[1]=n[1]*i,e[2]=n[2]*i,e[3]=0,e[4]=n[4]*s,e[5]=n[5]*s,e[6]=n[6]*s,e[7]=0,e[8]=n[8]*r,e[9]=n[9]*r,e[10]=n[10]*r,e[11]=0,e[12]=0,e[13]=0,e[14]=0,e[15]=1,this}makeRotationFromEuler(t){let e=this.elements,{x:n,y:i,z:s}=t,r=Math.cos(n),a=Math.sin(n),o=Math.cos(i),l=Math.sin(i),c=Math.cos(s),h=Math.sin(s);if(t.order==="XYZ"){let d=r*c,u=r*h,f=a*c,g=a*h;e[0]=o*c,e[4]=-o*h,e[8]=l,e[1]=u+f*l,e[5]=d-g*l,e[9]=-a*o,e[2]=g-d*l,e[6]=f+u*l,e[10]=r*o}else if(t.order==="YXZ"){let d=o*c,u=o*h,f=l*c,g=l*h;e[0]=d+g*a,e[4]=f*a-u,e[8]=r*l,e[1]=r*h,e[5]=r*c,e[9]=-a,e[2]=u*a-f,e[6]=g+d*a,e[10]=r*o}else if(t.order==="ZXY"){let d=o*c,u=o*h,f=l*c,g=l*h;e[0]=d-g*a,e[4]=-r*h,e[8]=f+u*a,e[1]=u+f*a,e[5]=r*c,e[9]=g-d*a,e[2]=-r*l,e[6]=a,e[10]=r*o}else if(t.order==="ZYX"){let d=r*c,u=r*h,f=a*c,g=a*h;e[0]=o*c,e[4]=f*l-u,e[8]=d*l+g,e[1]=o*h,e[5]=g*l+d,e[9]=u*l-f,e[2]=-l,e[6]=a*o,e[10]=r*o}else if(t.order==="YZX"){let d=r*o,u=r*l,f=a*o,g=a*l;e[0]=o*c,e[4]=g-d*h,e[8]=f*h+u,e[1]=h,e[5]=r*c,e[9]=-a*c,e[2]=-l*c,e[6]=u*h+f,e[10]=d-g*h}else if(t.order==="XZY"){let d=r*o,u=r*l,f=a*o,g=a*l;e[0]=o*c,e[4]=-h,e[8]=l*c,e[1]=d*h+g,e[5]=r*c,e[9]=u*h-f,e[2]=f*h-u,e[6]=a*c,e[10]=g*h+d}return e[3]=0,e[7]=0,e[11]=0,e[12]=0,e[13]=0,e[14]=0,e[15]=1,this}makeRotationFromQuaternion(t){return this.compose(Zg,t,Jg)}lookAt(t,e,n){let i=this.elements;if(en.subVectors(t,e),en.lengthSq()===0)en.z=1;if(en.normalize(),ni.crossVectors(n,en),ni.lengthSq()===0){if(Math.abs(n.z)===1)en.x+=0.0001;else en.z+=0.0001;en.normalize(),ni.crossVectors(n,en)}return ni.normalize(),Or.crossVectors(en,ni),i[0]=ni.x,i[4]=Or.x,i[8]=en.x,i[1]=ni.y,i[5]=Or.y,i[9]=en.y,i[2]=ni.z,i[6]=Or.z,i[10]=en.z,this}multiply(t){return this.multiplyMatrices(this,t)}premultiply(t){return this.multiplyMatrices(t,this)}multiplyMatrices(t,e){let n=t.elements,i=e.elements,s=this.elements,r=n[0],a=n[4],o=n[8],l=n[12],c=n[1],h=n[5],d=n[9],u=n[13],f=n[2],g=n[6],_=n[10],m=n[14],p=n[3],b=n[7],S=n[11],x=n[15],T=i[0],I=i[4],R=i[8],y=i[12],A=i[1],Z=i[5],D=i[9],q=i[13],ot=i[2],z=i[6],et=i[10],nt=i[14],B=i[3],ft=i[7],Q=i[11],ct=i[15];return s[0]=r*T+a*A+o*ot+l*B,s[4]=r*I+a*Z+o*z+l*ft,s[8]=r*R+a*D+o*et+l*Q,s[12]=r*y+a*q+o*nt+l*ct,s[1]=c*T+h*A+d*ot+u*B,s[5]=c*I+h*Z+d*z+u*ft,s[9]=c*R+h*D+d*et+u*Q,s[13]=c*y+h*q+d*nt+u*ct,s[2]=f*T+g*A+_*ot+m*B,s[6]=f*I+g*Z+_*z+m*ft,s[10]=f*R+g*D+_*et+m*Q,s[14]=f*y+g*q+_*nt+m*ct,s[3]=p*T+b*A+S*ot+x*B,s[7]=p*I+b*Z+S*z+x*ft,s[11]=p*R+b*D+S*et+x*Q,s[15]=p*y+b*q+S*nt+x*ct,this}multiplyScalar(t){let e=this.elements;return e[0]*=t,e[4]*=t,e[8]*=t,e[12]*=t,e[1]*=t,e[5]*=t,e[9]*=t,e[13]*=t,e[2]*=t,e[6]*=t,e[10]*=t,e[14]*=t,e[3]*=t,e[7]*=t,e[11]*=t,e[15]*=t,this}determinant(){let t=this.elements,e=t[0],n=t[4],i=t[8],s=t[12],r=t[1],a=t[5],o=t[9],l=t[13],c=t[2],h=t[6],d=t[10],u=t[14],f=t[3],g=t[7],_=t[11],m=t[15],p=o*u-l*d,b=a*u-l*h,S=a*d-o*h,x=r*u-l*c,T=r*d-o*c,I=r*h-a*c;return e*(g*p-_*b+m*S)-n*(f*p-_*x+m*T)+i*(f*b-g*x+m*I)-s*(f*S-g*T+_*I)}determinantAffine(){let t=this.elements,e=t[0],n=t[4],i=t[8],s=t[1],r=t[5],a=t[9],o=t[2],l=t[6],c=t[10];return e*(r*c-a*l)-n*(s*c-a*o)+i*(s*l-r*o)}transpose(){let t=this.elements,e;return e=t[1],t[1]=t[4],t[4]=e,e=t[2],t[2]=t[8],t[8]=e,e=t[6],t[6]=t[9],t[9]=e,e=t[3],t[3]=t[12],t[12]=e,e=t[7],t[7]=t[13],t[13]=e,e=t[11],t[11]=t[14],t[14]=e,this}setPosition(t,e,n){let i=this.elements;if(t.isVector3)i[12]=t.x,i[13]=t.y,i[14]=t.z;else i[12]=t,i[13]=e,i[14]=n;return this}invert(){let t=this.elements,e=t[0],n=t[1],i=t[2],s=t[3],r=t[4],a=t[5],o=t[6],l=t[7],c=t[8],h=t[9],d=t[10],u=t[11],f=t[12],g=t[13],_=t[14],m=t[15],p=e*a-n*r,b=e*o-i*r,S=e*l-s*r,x=n*o-i*a,T=n*l-s*a,I=i*l-s*o,R=c*g-h*f,y=c*_-d*f,A=c*m-u*f,Z=h*_-d*g,D=h*m-u*g,q=d*m-u*_,ot=p*q-b*D+S*Z+x*A-T*y+I*R;if(ot===0)return this.set(0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0);let z=1/ot;return t[0]=(a*q-o*D+l*Z)*z,t[1]=(i*D-n*q-s*Z)*z,t[2]=(g*I-_*T+m*x)*z,t[3]=(d*T-h*I-u*x)*z,t[4]=(o*A-r*q-l*y)*z,t[5]=(e*q-i*A+s*y)*z,t[6]=(_*S-f*I-m*b)*z,t[7]=(c*I-d*S+u*b)*z,t[8]=(r*D-a*A+l*R)*z,t[9]=(n*A-e*D-s*R)*z,t[10]=(f*T-g*S+m*p)*z,t[11]=(h*S-c*T-u*p)*z,t[12]=(a*y-r*Z-o*R)*z,t[13]=(e*Z-n*y+i*R)*z,t[14]=(g*b-f*x-_*p)*z,t[15]=(c*x-h*b+d*p)*z,this}scale(t){let e=this.elements,{x:n,y:i,z:s}=t;return e[0]*=n,e[4]*=i,e[8]*=s,e[1]*=n,e[5]*=i,e[9]*=s,e[2]*=n,e[6]*=i,e[10]*=s,e[3]*=n,e[7]*=i,e[11]*=s,this}getMaxScaleOnAxis(){let t=this.elements,e=t[0]*t[0]+t[1]*t[1]+t[2]*t[2],n=t[4]*t[4]+t[5]*t[5]+t[6]*t[6],i=t[8]*t[8]+t[9]*t[9]+t[10]*t[10];return Math.sqrt(Math.max(e,n,i))}makeTranslation(t,e,n){if(t.isVector3)this.set(1,0,0,t.x,0,1,0,t.y,0,0,1,t.z,0,0,0,1);else this.set(1,0,0,t,0,1,0,e,0,0,1,n,0,0,0,1);return this}makeRotationX(t){let e=Math.cos(t),n=Math.sin(t);return this.set(1,0,0,0,0,e,-n,0,0,n,e,0,0,0,0,1),this}makeRotationY(t){let e=Math.cos(t),n=Math.sin(t);return this.set(e,0,n,0,0,1,0,0,-n,0,e,0,0,0,0,1),this}makeRotationZ(t){let e=Math.cos(t),n=Math.sin(t);return this.set(e,-n,0,0,n,e,0,0,0,0,1,0,0,0,0,1),this}makeRotationAxis(t,e){let n=Math.cos(e),i=Math.sin(e),s=1-n,{x:r,y:a,z:o}=t,l=s*r,c=s*a;return this.set(l*r+n,l*a-i*o,l*o+i*a,0,l*a+i*o,c*a+n,c*o-i*r,0,l*o-i*a,c*o+i*r,s*o*o+n,0,0,0,0,1),this}makeScale(t,e,n){return this.set(t,0,0,0,0,e,0,0,0,0,n,0,0,0,0,1),this}makeShear(t,e,n,i,s,r){return this.set(1,n,s,0,t,1,r,0,e,i,1,0,0,0,0,1),this}compose(t,e,n){let i=this.elements,{_x:s,_y:r,_z:a,_w:o}=e,l=s+s,c=r+r,h=a+a,d=s*l,u=s*c,f=s*h,g=r*c,_=r*h,m=a*h,p=o*l,b=o*c,S=o*h,{x,y:T,z:I}=n;return i[0]=(1-(g+m))*x,i[1]=(u+S)*x,i[2]=(f-b)*x,i[3]=0,i[4]=(u-S)*T,i[5]=(1-(d+m))*T,i[6]=(_+p)*T,i[7]=0,i[8]=(f+b)*I,i[9]=(_-p)*I,i[10]=(1-(d+g))*I,i[11]=0,i[12]=t.x,i[13]=t.y,i[14]=t.z,i[15]=1,this}decompose(t,e,n){let i=this.elements;t.x=i[12],t.y=i[13],t.z=i[14];let s=this.determinantAffine();if(s===0)return n.set(1,1,1),e.identity(),this;let r=$i.set(i[0],i[1],i[2]).length(),a=$i.set(i[4],i[5],i[6]).length(),o=$i.set(i[8],i[9],i[10]).length();if(s<0)r=-r;xn.copy(this);let l=1/r,c=1/a,h=1/o;return xn.elements[0]*=l,xn.elements[1]*=l,xn.elements[2]*=l,xn.elements[4]*=c,xn.elements[5]*=c,xn.elements[6]*=c,xn.elements[8]*=h,xn.elements[9]*=h,xn.elements[10]*=h,e.setFromRotationMatrix(xn),n.x=r,n.y=a,n.z=o,this}makePerspective(t,e,n,i,s,r,a=2000,o=!1){let l=this.elements,c=2*s/(e-t),h=2*s/(n-i),d=(e+t)/(e-t),u=(n+i)/(n-i),f,g;if(o)f=s/(r-s),g=r*s/(r-s);else if(a===2000)f=-(r+s)/(r-s),g=-2*r*s/(r-s);else if(a===2001)f=-r/(r-s),g=-r*s/(r-s);else throw Error("THREE.Matrix4.makePerspective(): Invalid coordinate system: "+a);return l[0]=c,l[4]=0,l[8]=d,l[12]=0,l[1]=0,l[5]=h,l[9]=u,l[13]=0,l[2]=0,l[6]=0,l[10]=f,l[14]=g,l[3]=0,l[7]=0,l[11]=-1,l[15]=0,this}makeOrthographic(t,e,n,i,s,r,a=2000,o=!1){let l=this.elements,c=2/(e-t),h=2/(n-i),d=-(e+t)/(e-t),u=-(n+i)/(n-i),f,g;if(o)f=1/(r-s),g=r/(r-s);else if(a===2000)f=-2/(r-s),g=-(r+s)/(r-s);else if(a===2001)f=-1/(r-s),g=-s/(r-s);else throw Error("THREE.Matrix4.makeOrthographic(): Invalid coordinate system: "+a);return l[0]=c,l[4]=0,l[8]=0,l[12]=d,l[1]=0,l[5]=h,l[9]=0,l[13]=u,l[2]=0,l[6]=0,l[10]=f,l[14]=g,l[3]=0,l[7]=0,l[11]=0,l[15]=1,this}equals(t){let e=this.elements,n=t.elements;for(let i=0;i<16;i++)if(e[i]!==n[i])return!1;return!0}fromArray(t,e=0){for(let n=0;n<16;n++)this.elements[n]=t[n+e];return this}toArray(t=[],e=0){let n=this.elements;return t[e]=n[0],t[e+1]=n[1],t[e+2]=n[2],t[e+3]=n[3],t[e+4]=n[4],t[e+5]=n[5],t[e+6]=n[6],t[e+7]=n[7],t[e+8]=n[8],t[e+9]=n[9],t[e+10]=n[10],t[e+11]=n[11],t[e+12]=n[12],t[e+13]=n[13],t[e+14]=n[14],t[e+15]=n[15],t}}var $i=new L,xn=new Kt,Zg=new L(0,0,0),Jg=new L(1,1,1),ni=new L,Or=new L,en=new L,Jh=new Kt,$h=new Ce;class Mn{constructor(t=0,e=0,n=0,i=Mn.DEFAULT_ORDER){this.isEuler=!0,this._x=t,this._y=e,this._z=n,this._order=i}get x(){return this._x}set x(t){this._x=t,this._onChangeCallback()}get y(){return this._y}set y(t){this._y=t,this._onChangeCallback()}get z(){return this._z}set z(t){this._z=t,this._onChangeCallback()}get order(){return this._order}set order(t){this._order=t,this._onChangeCallback()}set(t,e,n,i=this._order){return this._x=t,this._y=e,this._z=n,this._order=i,this._onChangeCallback(),this}clone(){return new this.constructor(this._x,this._y,this._z,this._order)}copy(t){return this._x=t._x,this._y=t._y,this._z=t._z,this._order=t._order,this._onChangeCallback(),this}setFromRotationMatrix(t,e=this._order,n=!0){let i=t.elements,s=i[0],r=i[4],a=i[8],o=i[1],l=i[5],c=i[9],h=i[2],d=i[6],u=i[10];switch(e){case"XYZ":if(this._y=Math.asin(Zt(a,-1,1)),Math.abs(a)<0.9999999)this._x=Math.atan2(-c,u),this._z=Math.atan2(-r,s);else this._x=Math.atan2(d,l),this._z=0;break;case"YXZ":if(this._x=Math.asin(-Zt(c,-1,1)),Math.abs(c)<0.9999999)this._y=Math.atan2(a,u),this._z=Math.atan2(o,l);else this._y=Math.atan2(-h,s),this._z=0;break;case"ZXY":if(this._x=Math.asin(Zt(d,-1,1)),Math.abs(d)<0.9999999)this._y=Math.atan2(-h,u),this._z=Math.atan2(-r,l);else this._y=0,this._z=Math.atan2(o,s);break;case"ZYX":if(this._y=Math.asin(-Zt(h,-1,1)),Math.abs(h)<0.9999999)this._x=Math.atan2(d,u),this._z=Math.atan2(o,s);else this._x=0,this._z=Math.atan2(-r,l);break;case"YZX":if(this._z=Math.asin(Zt(o,-1,1)),Math.abs(o)<0.9999999)this._x=Math.atan2(-c,l),this._y=Math.atan2(-h,s);else this._x=0,this._y=Math.atan2(a,u);break;case"XZY":if(this._z=Math.asin(-Zt(r,-1,1)),Math.abs(r)<0.9999999)this._x=Math.atan2(d,l),this._y=Math.atan2(a,s);else this._x=Math.atan2(-c,u),this._y=0;break;default:Mt("Euler: .setFromRotationMatrix() encountered an unknown order: "+e)}if(this._order=e,n===!0)this._onChangeCallback();return this}setFromQuaternion(t,e,n){return Jh.makeRotationFromQuaternion(t),this.setFromRotationMatrix(Jh,e,n)}setFromVector3(t,e=this._order){return this.set(t.x,t.y,t.z,e)}reorder(t){return $h.setFromEuler(this),this.setFromQuaternion($h,t)}equals(t){return t._x===this._x&&t._y===this._y&&t._z===this._z&&t._order===this._order}fromArray(t){if(this._x=t[0],this._y=t[1],this._z=t[2],t[3]!==void 0)this._order=t[3];return this._onChangeCallback(),this}toArray(t=[],e=0){return t[e]=this._x,t[e+1]=this._y,t[e+2]=this._z,t[e+3]=this._order,t}_onChange(t){return this._onChangeCallback=t,this}_onChangeCallback(){}*[Symbol.iterator](){yield this._x,yield this._y,yield this._z,yield this._order}}Mn.DEFAULT_ORDER="XYZ";class ur{constructor(){this.mask=1}set(t){this.mask=(1<<t|0)>>>0}enable(t){this.mask|=1<<t|0}enableAll(){this.mask=-1}toggle(t){this.mask^=1<<t|0}disable(t){this.mask&=~(1<<t|0)}disableAll(){this.mask=0}test(t){return(this.mask&t.mask)!==0}isEnabled(t){return(this.mask&(1<<t|0))!==0}}var $g=0,Kh=new L,Ki=new Ce,zn=new Kt,Br=new L,Fs=new L,Kg=new L,jg=new Ce,jh=new L(1,0,0),Qh=new L(0,1,0),tu=new L(0,0,1),eu={type:"added"},Qg={type:"removed"},ji={type:"childadded",child:null},Vo={type:"childremoved",child:null};class ae extends pn{constructor(){super();this.isObject3D=!0,Object.defineProperty(this,"id",{value:$g++}),this.uuid=rn(),this.name="",this.type="Object3D",this.parent=null,this.children=[],this.up=ae.DEFAULT_UP.clone();let t=new L,e=new Mn,n=new Ce,i=new L(1,1,1);function s(){n.setFromEuler(e,!1)}function r(){e.setFromQuaternion(n,void 0,!1)}e._onChange(s),n._onChange(r),Object.defineProperties(this,{position:{configurable:!0,enumerable:!0,value:t},rotation:{configurable:!0,enumerable:!0,value:e},quaternion:{configurable:!0,enumerable:!0,value:n},scale:{configurable:!0,enumerable:!0,value:i},modelViewMatrix:{value:new Kt},normalMatrix:{value:new Xt}}),this.matrix=new Kt,this.matrixWorld=new Kt,this.matrixAutoUpdate=ae.DEFAULT_MATRIX_AUTO_UPDATE,this.matrixWorldAutoUpdate=ae.DEFAULT_MATRIX_WORLD_AUTO_UPDATE,this.matrixWorldNeedsUpdate=!1,this.layers=new ur,this.visible=!0,this.castShadow=!1,this.receiveShadow=!1,this.frustumCulled=!0,this.renderOrder=0,this.animations=[],this.customDepthMaterial=void 0,this.customDistanceMaterial=void 0,this.static=!1,this.userData={},this.pivot=null}onBeforeShadow(){}onAfterShadow(){}onBeforeRender(){}onAfterRender(){}applyMatrix4(t){if(this.matrixAutoUpdate)this.updateMatrix();this.matrix.premultiply(t),this.matrix.decompose(this.position,this.quaternion,this.scale)}applyQuaternion(t){return this.quaternion.premultiply(t),this}setRotationFromAxisAngle(t,e){this.quaternion.setFromAxisAngle(t,e)}setRotationFromEuler(t){this.quaternion.setFromEuler(t,!0)}setRotationFromMatrix(t){this.quaternion.setFromRotationMatrix(t)}setRotationFromQuaternion(t){this.quaternion.copy(t)}rotateOnAxis(t,e){return Ki.setFromAxisAngle(t,e),this.quaternion.multiply(Ki),this}rotateOnWorldAxis(t,e){return Ki.setFromAxisAngle(t,e),this.quaternion.premultiply(Ki),this}rotateX(t){return this.rotateOnAxis(jh,t)}rotateY(t){return this.rotateOnAxis(Qh,t)}rotateZ(t){return this.rotateOnAxis(tu,t)}translateOnAxis(t,e){return Kh.copy(t).applyQuaternion(this.quaternion),this.position.add(Kh.multiplyScalar(e)),this}translateX(t){return this.translateOnAxis(jh,t)}translateY(t){return this.translateOnAxis(Qh,t)}translateZ(t){return this.translateOnAxis(tu,t)}localToWorld(t){return this.updateWorldMatrix(!0,!1),t.applyMatrix4(this.matrixWorld)}worldToLocal(t){return this.updateWorldMatrix(!0,!1),t.applyMatrix4(zn.copy(this.matrixWorld).invert())}lookAt(t,e,n){if(t.isVector3)Br.copy(t);else Br.set(t,e,n);let i=this.parent;if(this.updateWorldMatrix(!0,!1),Fs.setFromMatrixPosition(this.matrixWorld),this.isCamera||this.isLight)zn.lookAt(Fs,Br,this.up);else zn.lookAt(Br,Fs,this.up);if(this.quaternion.setFromRotationMatrix(zn),i)zn.extractRotation(i.matrixWorld),Ki.setFromRotationMatrix(zn),this.quaternion.premultiply(Ki.invert())}add(t){if(arguments.length>1){for(let e=0;e<arguments.length;e++)this.add(arguments[e]);return this}if(t===this)return Wt("Object3D.add: object can't be added as a child of itself.",t),this;if(t&&t.isObject3D)t.removeFromParent(),t.parent=this,this.children.push(t),t.dispatchEvent(eu),ji.child=t,this.dispatchEvent(ji),ji.child=null;else Wt("Object3D.add: object not an instance of THREE.Object3D.",t);return this}remove(t){if(arguments.length>1){for(let n=0;n<arguments.length;n++)this.remove(arguments[n]);return this}let e=this.children.indexOf(t);if(e!==-1)t.parent=null,this.children.splice(e,1),t.dispatchEvent(Qg),Vo.child=t,this.dispatchEvent(Vo),Vo.child=null;return this}removeFromParent(){let t=this.parent;if(t!==null)t.remove(this);return this}clear(){return this.remove(...this.children)}attach(t){if(this.updateWorldMatrix(!0,!1),zn.copy(this.matrixWorld).invert(),t.parent!==null)t.parent.updateWorldMatrix(!0,!1),zn.multiply(t.parent.matrixWorld);return t.applyMatrix4(zn),t.removeFromParent(),t.parent=this,this.children.push(t),t.updateWorldMatrix(!1,!0),t.dispatchEvent(eu),ji.child=t,this.dispatchEvent(ji),ji.child=null,this}getObjectById(t){return this.getObjectByProperty("id",t)}getObjectByName(t){return this.getObjectByProperty("name",t)}getObjectByProperty(t,e){if(this[t]===e)return this;for(let n=0,i=this.children.length;n<i;n++){let r=this.children[n].getObjectByProperty(t,e);if(r!==void 0)return r}return}getObjectsByProperty(t,e,n=[]){if(this[t]===e)n.push(this);let i=this.children;for(let s=0,r=i.length;s<r;s++)i[s].getObjectsByProperty(t,e,n);return n}getWorldPosition(t){return this.updateWorldMatrix(!0,!1),t.setFromMatrixPosition(this.matrixWorld)}getWorldQuaternion(t){return this.updateWorldMatrix(!0,!1),this.matrixWorld.decompose(Fs,t,Kg),t}getWorldScale(t){return this.updateWorldMatrix(!0,!1),this.matrixWorld.decompose(Fs,jg,t),t}getWorldDirection(t){this.updateWorldMatrix(!0,!1);let e=this.matrixWorld.elements;return t.set(e[8],e[9],e[10]).normalize()}raycast(){}intersectsFrustum(){}traverse(t){t(this);let e=this.children;for(let n=0,i=e.length;n<i;n++)e[n].traverse(t)}traverseVisible(t){if(this.visible===!1)return;t(this);let e=this.children;for(let n=0,i=e.length;n<i;n++)e[n].traverseVisible(t)}traverseAncestors(t){let e=this.parent;if(e!==null)t(e),e.traverseAncestors(t)}updateMatrix(){this.matrix.compose(this.position,this.quaternion,this.scale);let t=this.pivot;if(t!==null){let{x:e,y:n,z:i}=t,s=this.matrix.elements;s[12]+=e-s[0]*e-s[4]*n-s[8]*i,s[13]+=n-s[1]*e-s[5]*n-s[9]*i,s[14]+=i-s[2]*e-s[6]*n-s[10]*i}this.matrixWorldNeedsUpdate=!0}updateMatrixWorld(t){if(this.matrixAutoUpdate)this.updateMatrix();if(this.matrixWorldNeedsUpdate||t){if(this.matrixWorldAutoUpdate===!0)if(this.parent===null)this.matrixWorld.copy(this.matrix);else this.matrixWorld.multiplyMatrices(this.parent.matrixWorld,this.matrix);this.matrixWorldNeedsUpdate=!1,t=!0}let e=this.children;for(let n=0,i=e.length;n<i;n++)e[n].updateMatrixWorld(t)}updateWorldMatrix(t,e,n=!1){let i=this.parent;if(t===!0&&i!==null)i.updateWorldMatrix(!0,!1);if(this.matrixAutoUpdate)this.updateMatrix();if(this.matrixWorldNeedsUpdate||n){if(this.matrixWorldAutoUpdate===!0)if(this.parent===null)this.matrixWorld.copy(this.matrix);else this.matrixWorld.multiplyMatrices(this.parent.matrixWorld,this.matrix);this.matrixWorldNeedsUpdate=!1,n=!0}if(e===!0){let s=this.children;for(let r=0,a=s.length;r<a;r++)s[r].updateWorldMatrix(!1,!0,n)}}toJSON(t){let e=t===void 0||typeof t==="string",n={};if(e)t={geometries:{},materials:{},textures:{},images:{},shapes:{},skeletons:{},animations:{},nodes:{}},n.metadata={version:4.7,type:"Object",generator:"Object3D.toJSON"};let i={};if(i.uuid=this.uuid,i.type=this.type,i.name=this.name,i.castShadow=this.castShadow,i.receiveShadow=this.receiveShadow,i.visible=this.visible,i.frustumCulled=this.frustumCulled,i.renderOrder=this.renderOrder,i.static=this.static,i.matrixAutoUpdate=this.matrixAutoUpdate,Object.keys(this.userData).length>0)i.userData=this.userData;if(i.layers=this.layers.mask,i.matrix=this.matrix.toArray(),i.up=this.up.toArray(),this.pivot!==null)i.pivot=this.pivot.toArray();if(this.morphTargetDictionary!==void 0)i.morphTargetDictionary=Object.assign({},this.morphTargetDictionary);if(this.morphTargetInfluences!==void 0)i.morphTargetInfluences=this.morphTargetInfluences.slice();if(this.isInstancedMesh){if(i.type="InstancedMesh",i.count=this.count,i.instanceMatrix=this.instanceMatrix.toJSON(),this.instanceColor!==null)i.instanceColor=this.instanceColor.toJSON()}if(this.isBatchedMesh){if(i.type="BatchedMesh",i.perObjectFrustumCulled=this.perObjectFrustumCulled,i.sortObjects=this.sortObjects,i.drawRanges=this._drawRanges,i.reservedRanges=this._reservedRanges,i.geometryInfo=this._geometryInfo.map((a)=>({...a,boundingBox:a.boundingBox?a.boundingBox.toJSON():void 0,boundingSphere:a.boundingSphere?a.boundingSphere.toJSON():void 0})),i.instanceInfo=this._instanceInfo.map((a)=>({...a})),i.availableInstanceIds=this._availableInstanceIds.slice(),i.availableGeometryIds=this._availableGeometryIds.slice(),i.nextIndexStart=this._nextIndexStart,i.nextVertexStart=this._nextVertexStart,i.geometryCount=this._geometryCount,i.maxInstanceCount=this._maxInstanceCount,i.maxVertexCount=this._maxVertexCount,i.maxIndexCount=this._maxIndexCount,i.geometryInitialized=this._geometryInitialized,i.matricesTexture=this._matricesTexture.toJSON(t),i.indirectTexture=this._indirectTexture.toJSON(t),this._colorsTexture!==null)i.colorsTexture=this._colorsTexture.toJSON(t);if(this.boundingSphere!==null)i.boundingSphere=this.boundingSphere.toJSON();if(this.boundingBox!==null)i.boundingBox=this.boundingBox.toJSON()}function s(a,o){if(a[o.uuid]===void 0)a[o.uuid]=o.toJSON(t);return o.uuid}if(this.isScene){if(this.background){if(this.background.isColor)i.background=this.background.toJSON();else if(this.background.isTexture)i.background=this.background.toJSON(t).uuid}if(this.environment&&this.environment.isTexture&&this.environment.isRenderTargetTexture!==!0)i.environment=this.environment.toJSON(t).uuid}else if(this.isMesh||this.isLine||this.isPoints){i.geometry=s(t.geometries,this.geometry);let a=this.geometry.parameters;if(a!==void 0&&a.shapes!==void 0){let o=a.shapes;if(Array.isArray(o))for(let l=0,c=o.length;l<c;l++){let h=o[l];s(t.shapes,h)}else s(t.shapes,o)}}if(this.isSkinnedMesh){if(i.bindMode=this.bindMode,i.bindMatrix=this.bindMatrix.toArray(),this.skeleton!==void 0)s(t.skeletons,this.skeleton),i.skeleton=this.skeleton.uuid}if(this.material!==void 0)if(Array.isArray(this.material)){let a=[];for(let o=0,l=this.material.length;o<l;o++)a.push(s(t.materials,this.material[o]));i.material=a}else i.material=s(t.materials,this.material);if(this.children.length>0){i.children=[];for(let a=0;a<this.children.length;a++)i.children.push(this.children[a].toJSON(t).object)}if(this.animations.length>0){i.animations=[];for(let a=0;a<this.animations.length;a++){let o=this.animations[a];i.animations.push(s(t.animations,o))}}if(e){let a=r(t.geometries),o=r(t.materials),l=r(t.textures),c=r(t.images),h=r(t.shapes),d=r(t.skeletons),u=r(t.animations),f=r(t.nodes);if(a.length>0)n.geometries=a;if(o.length>0)n.materials=o;if(l.length>0)n.textures=l;if(c.length>0)n.images=c;if(h.length>0)n.shapes=h;if(d.length>0)n.skeletons=d;if(u.length>0)n.animations=u;if(f.length>0)n.nodes=f}return n.object=i,n;function r(a){let o=[];for(let l in a){let c=a[l];delete c.metadata,o.push(c)}return o}}clone(t){return new this.constructor().copy(this,t)}copy(t,e=!0){if(this.name=t.name,this.up.copy(t.up),this.position.copy(t.position),this.rotation.order=t.rotation.order,this.quaternion.copy(t.quaternion),this.scale.copy(t.scale),this.pivot=t.pivot!==null?t.pivot.clone():null,this.matrix.copy(t.matrix),this.matrixWorld.copy(t.matrixWorld),this.matrixAutoUpdate=t.matrixAutoUpdate,this.matrixWorldAutoUpdate=t.matrixWorldAutoUpdate,this.matrixWorldNeedsUpdate=t.matrixWorldNeedsUpdate,this.layers.mask=t.layers.mask,this.visible=t.visible,this.castShadow=t.castShadow,this.receiveShadow=t.receiveShadow,this.frustumCulled=t.frustumCulled,this.renderOrder=t.renderOrder,this.static=t.static,this.animations=t.animations.slice(),this.userData=JSON.parse(JSON.stringify(t.userData)),e===!0)for(let n=0;n<t.children.length;n++){let i=t.children[n];this.add(i.clone())}return this}dispose(){this.dispatchEvent({type:"dispose"})}}ae.DEFAULT_UP=new L(0,1,0);ae.DEFAULT_MATRIX_AUTO_UPDATE=!0;ae.DEFAULT_MATRIX_WORLD_AUTO_UPDATE=!0;class Li extends ae{constructor(){super();this.isGroup=!0,this.type="Group"}}var t0={type:"move"};class dr{constructor(){this._targetRay=null,this._grip=null,this._hand=null}getHandSpace(){if(this._hand===null)this._hand=new Li,this._hand.matrixAutoUpdate=!1,this._hand.visible=!1,this._hand.joints={},this._hand.inputState={pinching:!1};return this._hand}getTargetRaySpace(){if(this._targetRay===null)this._targetRay=new Li,this._targetRay.matrixAutoUpdate=!1,this._targetRay.visible=!1,this._targetRay.hasLinearVelocity=!1,this._targetRay.linearVelocity=new L,this._targetRay.hasAngularVelocity=!1,this._targetRay.angularVelocity=new L;return this._targetRay}getGripSpace(){if(this._grip===null)this._grip=new Li,this._grip.matrixAutoUpdate=!1,this._grip.visible=!1,this._grip.hasLinearVelocity=!1,this._grip.linearVelocity=new L,this._grip.hasAngularVelocity=!1,this._grip.angularVelocity=new L,this._grip.eventsEnabled=!1;return this._grip}dispatchEvent(t){if(this._targetRay!==null)this._targetRay.dispatchEvent(t);if(this._grip!==null)this._grip.dispatchEvent(t);if(this._hand!==null)this._hand.dispatchEvent(t);return this}connect(t){if(t&&t.hand){let e=this._hand;if(e)for(let n of t.hand.values())this._getHandJoint(e,n)}return this.dispatchEvent({type:"connected",data:t}),this}disconnect(t){if(this.dispatchEvent({type:"disconnected",data:t}),this._targetRay!==null)this._targetRay.visible=!1;if(this._grip!==null)this._grip.visible=!1;if(this._hand!==null)this._hand.visible=!1;return this}update(t,e,n){let i=null,s=null,r=null,a=this._targetRay,o=this._grip,l=this._hand;if(t&&e.session.visibilityState!=="visible-blurred"){if(l&&t.hand){r=!0;for(let g of t.hand.values()){let _=e.getJointPose(g,n),m=this._getHandJoint(l,g);if(_!==null)m.matrix.fromArray(_.transform.matrix),m.matrix.decompose(m.position,m.rotation,m.scale),m.matrixWorldNeedsUpdate=!0,m.jointRadius=_.radius;m.visible=_!==null}let c=l.joints["index-finger-tip"],h=l.joints["thumb-tip"],d=c.position.distanceTo(h.position),u=0.02,f=0.005;if(l.inputState.pinching&&d>u+f)l.inputState.pinching=!1,this.dispatchEvent({type:"pinchend",handedness:t.handedness,target:this});else if(!l.inputState.pinching&&d<=u-f)l.inputState.pinching=!0,this.dispatchEvent({type:"pinchstart",handedness:t.handedness,target:this})}else if(o!==null&&t.gripSpace){if(s=e.getPose(t.gripSpace,n),s!==null){if(o.matrix.fromArray(s.transform.matrix),o.matrix.decompose(o.position,o.rotation,o.scale),o.matrixWorldNeedsUpdate=!0,s.linearVelocity)o.hasLinearVelocity=!0,o.linearVelocity.copy(s.linearVelocity);else o.hasLinearVelocity=!1;if(s.angularVelocity)o.hasAngularVelocity=!0,o.angularVelocity.copy(s.angularVelocity);else o.hasAngularVelocity=!1;if(o.eventsEnabled)o.dispatchEvent({type:"gripUpdated",data:t,target:this})}}if(a!==null){if(i=e.getPose(t.targetRaySpace,n),i===null&&s!==null)i=s;if(i!==null){if(a.matrix.fromArray(i.transform.matrix),a.matrix.decompose(a.position,a.rotation,a.scale),a.matrixWorldNeedsUpdate=!0,i.linearVelocity)a.hasLinearVelocity=!0,a.linearVelocity.copy(i.linearVelocity);else a.hasLinearVelocity=!1;if(i.angularVelocity)a.hasAngularVelocity=!0,a.angularVelocity.copy(i.angularVelocity);else a.hasAngularVelocity=!1;this.dispatchEvent(t0)}}}if(a!==null)a.visible=i!==null;if(o!==null)o.visible=s!==null;if(l!==null)l.visible=r!==null;return this}_getHandJoint(t,e){if(t.joints[e.jointName]===void 0){let n=new Li;n.matrixAutoUpdate=!1,n.visible=!1,t.joints[e.jointName]=n,t.add(n)}return t.joints[e.jointName]}}var lf={aliceblue:15792383,antiquewhite:16444375,aqua:65535,aquamarine:8388564,azure:15794175,beige:16119260,bisque:16770244,black:0,blanchedalmond:16772045,blue:255,blueviolet:9055202,brown:10824234,burlywood:14596231,cadetblue:6266528,chartreuse:8388352,chocolate:13789470,coral:16744272,cornflowerblue:6591981,cornsilk:16775388,crimson:14423100,cyan:65535,darkblue:139,darkcyan:35723,darkgoldenrod:12092939,darkgray:11119017,darkgreen:25600,darkgrey:11119017,darkkhaki:12433259,darkmagenta:9109643,darkolivegreen:5597999,darkorange:16747520,darkorchid:10040012,darkred:9109504,darksalmon:15308410,darkseagreen:9419919,darkslateblue:4734347,darkslategray:3100495,darkslategrey:3100495,darkturquoise:52945,darkviolet:9699539,deeppink:16716947,deepskyblue:49151,dimgray:6908265,dimgrey:6908265,dodgerblue:2003199,firebrick:11674146,floralwhite:16775920,forestgreen:2263842,fuchsia:16711935,gainsboro:14474460,ghostwhite:16316671,gold:16766720,goldenrod:14329120,gray:8421504,green:32768,greenyellow:11403055,grey:8421504,honeydew:15794160,hotpink:16738740,indianred:13458524,indigo:4915330,ivory:16777200,khaki:15787660,lavender:15132410,lavenderblush:16773365,lawngreen:8190976,lemonchiffon:16775885,lightblue:11393254,lightcoral:15761536,lightcyan:14745599,lightgoldenrodyellow:16448210,lightgray:13882323,lightgreen:9498256,lightgrey:13882323,lightpink:16758465,lightsalmon:16752762,lightseagreen:2142890,lightskyblue:8900346,lightslategray:7833753,lightslategrey:7833753,lightsteelblue:11584734,lightyellow:16777184,lime:65280,limegreen:3329330,linen:16445670,magenta:16711935,maroon:8388608,mediumaquamarine:6737322,mediumblue:205,mediumorchid:12211667,mediumpurple:9662683,mediumseagreen:3978097,mediumslateblue:8087790,mediumspringgreen:64154,mediumturquoise:4772300,mediumvioletred:13047173,midnightblue:1644912,mintcream:16121850,mistyrose:16770273,moccasin:16770229,navajowhite:16768685,navy:128,oldlace:16643558,olive:8421376,olivedrab:7048739,orange:16753920,orangered:16729344,orchid:14315734,palegoldenrod:15657130,palegreen:10025880,paleturquoise:11529966,palevioletred:14381203,papayawhip:16773077,peachpuff:16767673,peru:13468991,pink:16761035,plum:14524637,powderblue:11591910,purple:8388736,rebeccapurple:6697881,red:16711680,rosybrown:12357519,royalblue:4286945,saddlebrown:9127187,salmon:16416882,sandybrown:16032864,seagreen:3050327,seashell:16774638,sienna:10506797,silver:12632256,skyblue:8900331,slateblue:6970061,slategray:7372944,slategrey:7372944,snow:16775930,springgreen:65407,steelblue:4620980,tan:13808780,teal:32896,thistle:14204888,tomato:16737095,turquoise:4251856,violet:15631086,wheat:16113331,white:16777215,whitesmoke:16119285,yellow:16776960,yellowgreen:10145074},ii={h:0,s:0,l:0},zr={h:0,s:0,l:0};function Wo(t,e,n){if(n<0)n+=1;if(n>1)n-=1;if(n<0.16666666666666666)return t+(e-t)*6*n;if(n<0.5)return e;if(n<0.6666666666666666)return t+(e-t)*6*(0.6666666666666666-n);return t}class Pt{constructor(t,e,n){return this.isColor=!0,this.r=1,this.g=1,this.b=1,this.set(t,e,n)}set(t,e,n){if(e===void 0&&n===void 0){let i=t;if(i&&i.isColor)this.copy(i);else if(typeof i==="number")this.setHex(i);else if(typeof i==="string")this.setStyle(i)}else this.setRGB(t,e,n);return this}setScalar(t){return this.r=t,this.g=t,this.b=t,this}setHex(t,e="srgb"){return t=Math.floor(t),this.r=(t>>16&255)/255,this.g=(t>>8&255)/255,this.b=(t&255)/255,re.colorSpaceToWorking(this,e),this}setRGB(t,e,n,i=re.workingColorSpace){return this.r=t,this.g=e,this.b=n,re.colorSpaceToWorking(this,i),this}setHSL(t,e,n,i=re.workingColorSpace){if(t=Ec(t,1),e=Zt(e,0,1),n=Zt(n,0,1),e===0)this.r=this.g=this.b=n;else{let s=n<=0.5?n*(1+e):n+e-n*e,r=2*n-s;this.r=Wo(r,s,t+0.3333333333333333),this.g=Wo(r,s,t),this.b=Wo(r,s,t-0.3333333333333333)}return re.colorSpaceToWorking(this,i),this}setStyle(t,e="srgb"){function n(s){if(s===void 0)return;if(parseFloat(s)<1)Mt("Color: Alpha component of "+t+" will be ignored.")}let i;if(i=/^(\w+)\(([^\)]*)\)/.exec(t)){let s,r=i[1],a=i[2];switch(r){case"rgb":case"rgba":if(s=/^\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:,\s*(\d*\.?\d+)\s*)?$/.exec(a))return n(s[4]),this.setRGB(Math.min(255,parseInt(s[1],10))/255,Math.min(255,parseInt(s[2],10))/255,Math.min(255,parseInt(s[3],10))/255,e);if(s=/^\s*(\d+)\%\s*,\s*(\d+)\%\s*,\s*(\d+)\%\s*(?:,\s*(\d*\.?\d+)\s*)?$/.exec(a))return n(s[4]),this.setRGB(Math.min(100,parseInt(s[1],10))/100,Math.min(100,parseInt(s[2],10))/100,Math.min(100,parseInt(s[3],10))/100,e);break;case"hsl":case"hsla":if(s=/^\s*(\d*\.?\d+)\s*,\s*(\d*\.?\d+)\%\s*,\s*(\d*\.?\d+)\%\s*(?:,\s*(\d*\.?\d+)\s*)?$/.exec(a))return n(s[4]),this.setHSL(parseFloat(s[1])/360,parseFloat(s[2])/100,parseFloat(s[3])/100,e);break;default:Mt("Color: Unknown color model "+t)}}else if(i=/^\#([A-Fa-f\d]+)$/.exec(t)){let s=i[1],r=s.length;if(r===3)return this.setRGB(parseInt(s.charAt(0),16)/15,parseInt(s.charAt(1),16)/15,parseInt(s.charAt(2),16)/15,e);else if(r===6)return this.setHex(parseInt(s,16),e);else Mt("Color: Invalid hex color "+t)}else if(t&&t.length>0)return this.setColorName(t,e);return this}setColorName(t,e="srgb"){let n=lf[t.toLowerCase()];if(n!==void 0)this.setHex(n,e);else Mt("Color: Unknown color "+t);return this}clone(){return new this.constructor(this.r,this.g,this.b)}copy(t){return this.r=t.r,this.g=t.g,this.b=t.b,this}copySRGBToLinear(t){return this.r=Yn(t.r),this.g=Yn(t.g),this.b=Yn(t.b),this}copyLinearToSRGB(t){return this.r=gs(t.r),this.g=gs(t.g),this.b=gs(t.b),this}convertSRGBToLinear(){return this.copySRGBToLinear(this),this}convertLinearToSRGB(){return this.copyLinearToSRGB(this),this}getHex(t="srgb"){return re.workingToColorSpace(Ge.copy(this),t),Math.round(Zt(Ge.r*255,0,255))*65536+Math.round(Zt(Ge.g*255,0,255))*256+Math.round(Zt(Ge.b*255,0,255))}getHexString(t="srgb"){return("000000"+this.getHex(t).toString(16)).slice(-6)}getHSL(t,e=re.workingColorSpace){re.workingToColorSpace(Ge.copy(this),e);let{r:n,g:i,b:s}=Ge,r=Math.max(n,i,s),a=Math.min(n,i,s),o,l,c=(a+r)/2;if(a===r)o=0,l=0;else{let h=r-a;switch(l=c<=0.5?h/(r+a):h/(2-r-a),r){case n:o=(i-s)/h+(i<s?6:0);break;case i:o=(s-n)/h+2;break;case s:o=(n-i)/h+4;break}o/=6}return t.h=o,t.s=l,t.l=c,t}getRGB(t,e=re.workingColorSpace){return re.workingToColorSpace(Ge.copy(this),e),t.r=Ge.r,t.g=Ge.g,t.b=Ge.b,t}getStyle(t="srgb"){re.workingToColorSpace(Ge.copy(this),t);let{r:e,g:n,b:i}=Ge;if(t!=="srgb")return`color(${t} ${e.toFixed(3)} ${n.toFixed(3)} ${i.toFixed(3)})`;return`rgb(${Math.round(e*255)},${Math.round(n*255)},${Math.round(i*255)})`}offsetHSL(t,e,n){return this.getHSL(ii),this.setHSL(ii.h+t,ii.s+e,ii.l+n)}add(t){return this.r+=t.r,this.g+=t.g,this.b+=t.b,this}addColors(t,e){return this.r=t.r+e.r,this.g=t.g+e.g,this.b=t.b+e.b,this}addScalar(t){return this.r+=t,this.g+=t,this.b+=t,this}sub(t){return this.r=Math.max(0,this.r-t.r),this.g=Math.max(0,this.g-t.g),this.b=Math.max(0,this.b-t.b),this}multiply(t){return this.r*=t.r,this.g*=t.g,this.b*=t.b,this}multiplyScalar(t){return this.r*=t,this.g*=t,this.b*=t,this}lerp(t,e){return this.r+=(t.r-this.r)*e,this.g+=(t.g-this.g)*e,this.b+=(t.b-this.b)*e,this}lerpColors(t,e,n){return this.r=t.r+(e.r-t.r)*n,this.g=t.g+(e.g-t.g)*n,this.b=t.b+(e.b-t.b)*n,this}lerpHSL(t,e){this.getHSL(ii),t.getHSL(zr);let n=$s(ii.h,zr.h,e),i=$s(ii.s,zr.s,e),s=$s(ii.l,zr.l,e);return this.setHSL(n,i,s),this}setFromVector3(t){return this.r=t.x,this.g=t.y,this.b=t.z,this}applyMatrix3(t){let e=this.r,n=this.g,i=this.b,s=t.elements;return this.r=s[0]*e+s[3]*n+s[6]*i,this.g=s[1]*e+s[4]*n+s[7]*i,this.b=s[2]*e+s[5]*n+s[8]*i,this}equals(t){return t.r===this.r&&t.g===this.g&&t.b===this.b}fromArray(t,e=0){return this.r=t[e],this.g=t[e+1],this.b=t[e+2],this}toArray(t=[],e=0){return t[e]=this.r,t[e+1]=this.g,t[e+2]=this.b,t}fromBufferAttribute(t,e){return this.r=t.getX(e),this.g=t.getY(e),this.b=t.getZ(e),this}toJSON(){return this.getHex()}*[Symbol.iterator](){yield this.r,yield this.g,yield this.b}}var Ge=new Pt;Pt.NAMES=lf;class Xa{constructor(t,e=0.00025){this.isFogExp2=!0,this.name="",this.color=new Pt(t),this.density=e}clone(){return new Xa(this.color,this.density)}toJSON(){return{type:"FogExp2",name:this.name,color:this.color.getHex(),density:this.density}}}class qa{constructor(t,e=1,n=1000){this.isFog=!0,this.name="",this.color=new Pt(t),this.near=e,this.far=n}clone(){return new qa(this.color,this.near,this.far)}toJSON(){return{type:"Fog",name:this.name,color:this.color.getHex(),near:this.near,far:this.far}}}class fr extends ae{constructor(){super();if(this.isScene=!0,this.type="Scene",this.background=null,this.environment=null,this.fog=null,this.backgroundBlurriness=0,this.backgroundIntensity=1,this.backgroundRotation=new Mn,this.environmentIntensity=1,this.environmentRotation=new Mn,this.overrideMaterial=null,typeof __THREE_DEVTOOLS__<"u")__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("observe",{detail:this}))}copy(t,e){if(super.copy(t,e),t.background!==null)this.background=t.background.clone();if(t.environment!==null)this.environment=t.environment.clone();if(t.fog!==null)this.fog=t.fog.clone();if(this.backgroundBlurriness=t.backgroundBlurriness,this.backgroundIntensity=t.backgroundIntensity,this.backgroundRotation.copy(t.backgroundRotation),this.environmentIntensity=t.environmentIntensity,this.environmentRotation.copy(t.environmentRotation),t.overrideMaterial!==null)this.overrideMaterial=t.overrideMaterial.clone();return this.matrixAutoUpdate=t.matrixAutoUpdate,this}toJSON(t){let e=super.toJSON(t);if(this.fog!==null)e.object.fog=this.fog.toJSON();return e.object.backgroundBlurriness=this.backgroundBlurriness,e.object.backgroundIntensity=this.backgroundIntensity,e.object.backgroundRotation=this.backgroundRotation.toArray(),e.object.environmentIntensity=this.environmentIntensity,e.object.environmentRotation=this.environmentRotation.toArray(),e}}var yn=new L,kn=new L,Xo=new L,Gn=new L,Qi=new L,ts=new L,nu=new L,qo=new L,Yo=new L,Zo=new L,Jo=new fe,$o=new fe,Ko=new fe;class je{constructor(t=new L,e=new L,n=new L){this.a=t,this.b=e,this.c=n}static getNormal(t,e,n,i){i.subVectors(n,e),yn.subVectors(t,e),i.cross(yn);let s=i.lengthSq();if(s>0)return i.multiplyScalar(1/Math.sqrt(s));return i.set(0,0,0)}static getBarycoord(t,e,n,i,s){yn.subVectors(i,e),kn.subVectors(n,e),Xo.subVectors(t,e);let r=yn.dot(yn),a=yn.dot(kn),o=yn.dot(Xo),l=kn.dot(kn),c=kn.dot(Xo),h=r*l-a*a;if(h===0)return s.set(0,0,0),null;let d=1/h,u=(l*o-a*c)*d,f=(r*c-a*o)*d;return s.set(1-u-f,f,u)}static containsPoint(t,e,n,i){if(this.getBarycoord(t,e,n,i,Gn)===null)return!1;return Gn.x>=0&&Gn.y>=0&&Gn.x+Gn.y<=1}static getInterpolation(t,e,n,i,s,r,a,o){if(this.getBarycoord(t,e,n,i,Gn)===null){if(o.x=0,o.y=0,"z"in o)o.z=0;if("w"in o)o.w=0;return null}return o.setScalar(0),o.addScaledVector(s,Gn.x),o.addScaledVector(r,Gn.y),o.addScaledVector(a,Gn.z),o}static getInterpolatedAttribute(t,e,n,i,s,r){return Jo.setScalar(0),$o.setScalar(0),Ko.setScalar(0),Jo.fromBufferAttribute(t,e),$o.fromBufferAttribute(t,n),Ko.fromBufferAttribute(t,i),r.setScalar(0),r.addScaledVector(Jo,s.x),r.addScaledVector($o,s.y),r.addScaledVector(Ko,s.z),r}static isFrontFacing(t,e,n,i){return yn.subVectors(n,e),kn.subVectors(t,e),yn.cross(kn).dot(i)<0}set(t,e,n){return this.a.copy(t),this.b.copy(e),this.c.copy(n),this}setFromPointsAndIndices(t,e,n,i){return this.a.copy(t[e]),this.b.copy(t[n]),this.c.copy(t[i]),this}setFromAttributeAndIndices(t,e,n,i){return this.a.fromBufferAttribute(t,e),this.b.fromBufferAttribute(t,n),this.c.fromBufferAttribute(t,i),this}clone(){return new this.constructor().copy(this)}copy(t){return this.a.copy(t.a),this.b.copy(t.b),this.c.copy(t.c),this}getArea(){return yn.subVectors(this.c,this.b),kn.subVectors(this.a,this.b),yn.cross(kn).length()*0.5}getMidpoint(t){return t.addVectors(this.a,this.b).add(this.c).multiplyScalar(0.3333333333333333)}getNormal(t){return je.getNormal(this.a,this.b,this.c,t)}getPlane(t){return t.setFromCoplanarPoints(this.a,this.b,this.c)}getBarycoord(t,e){return je.getBarycoord(t,this.a,this.b,this.c,e)}getInterpolation(t,e,n,i,s){return je.getInterpolation(t,this.a,this.b,this.c,e,n,i,s)}containsPoint(t){return je.containsPoint(t,this.a,this.b,this.c)}isFrontFacing(t){return je.isFrontFacing(this.a,this.b,this.c,t)}intersectsBox(t){return t.intersectsTriangle(this)}closestPointToPoint(t,e){let n=this.a,i=this.b,s=this.c,r,a;Qi.subVectors(i,n),ts.subVectors(s,n),qo.subVectors(t,n);let o=Qi.dot(qo),l=ts.dot(qo);if(o<=0&&l<=0)return e.copy(n);Yo.subVectors(t,i);let c=Qi.dot(Yo),h=ts.dot(Yo);if(c>=0&&h<=c)return e.copy(i);let d=o*h-c*l;if(d<=0&&o>=0&&c<=0)return r=o/(o-c),e.copy(n).addScaledVector(Qi,r);Zo.subVectors(t,s);let u=Qi.dot(Zo),f=ts.dot(Zo);if(f>=0&&u<=f)return e.copy(s);let g=u*l-o*f;if(g<=0&&l>=0&&f<=0)return a=l/(l-f),e.copy(n).addScaledVector(ts,a);let _=c*f-u*h;if(_<=0&&h-c>=0&&u-f>=0)return nu.subVectors(s,i),a=(h-c)/(h-c+(u-f)),e.copy(i).addScaledVector(nu,a);let m=1/(_+g+d);return r=g*m,a=d*m,e.copy(n).addScaledVector(Qi,r).addScaledVector(ts,a)}equals(t){return t.a.equals(this.a)&&t.b.equals(this.b)&&t.c.equals(this.c)}}class Oe{constructor(t=new L(1/0,1/0,1/0),e=new L(-1/0,-1/0,-1/0)){this.isBox3=!0,this.min=t,this.max=e}set(t,e){return this.min.copy(t),this.max.copy(e),this}setFromArray(t){this.makeEmpty();for(let e=0,n=t.length;e<n;e+=3)this.expandByPoint(vn.fromArray(t,e));return this}setFromBufferAttribute(t){this.makeEmpty();for(let e=0,n=t.count;e<n;e++)this.expandByPoint(vn.fromBufferAttribute(t,e));return this}setFromPoints(t){this.makeEmpty();for(let e=0,n=t.length;e<n;e++)this.expandByPoint(t[e]);return this}setFromCenterAndSize(t,e){let n=vn.copy(e).multiplyScalar(0.5);return this.min.copy(t).sub(n),this.max.copy(t).add(n),this}setFromObject(t,e=!1){return this.makeEmpty(),this.expandByObject(t,e)}clone(){return new this.constructor().copy(this)}copy(t){return this.min.copy(t.min),this.max.copy(t.max),this}makeEmpty(){return this.min.x=this.min.y=this.min.z=1/0,this.max.x=this.max.y=this.max.z=-1/0,this}isEmpty(){return this.max.x<this.min.x||this.max.y<this.min.y||this.max.z<this.min.z}getCenter(t){return this.isEmpty()?t.set(0,0,0):t.addVectors(this.min,this.max).multiplyScalar(0.5)}getSize(t){return this.isEmpty()?t.set(0,0,0):t.subVectors(this.max,this.min)}expandByPoint(t){return this.min.min(t),this.max.max(t),this}expandByVector(t){return this.min.sub(t),this.max.add(t),this}expandByScalar(t){return this.min.addScalar(-t),this.max.addScalar(t),this}expandByObject(t,e=!1){t.updateWorldMatrix(!1,!1);let n=t.geometry;if(n!==void 0){let s=n.getAttribute("position");if(e===!0&&s!==void 0&&t.isInstancedMesh!==!0)for(let r=0,a=s.count;r<a;r++){if(t.isMesh===!0)t.getVertexPosition(r,vn);else vn.fromBufferAttribute(s,r);vn.applyMatrix4(t.matrixWorld),this.expandByPoint(vn)}else{if(t.boundingBox!==void 0){if(t.boundingBox===null)t.computeBoundingBox();kr.copy(t.boundingBox)}else{if(n.boundingBox===null)n.computeBoundingBox();kr.copy(n.boundingBox)}kr.applyMatrix4(t.matrixWorld),this.union(kr)}}let i=t.children;for(let s=0,r=i.length;s<r;s++)this.expandByObject(i[s],e);return this}containsPoint(t){return t.x>=this.min.x&&t.x<=this.max.x&&t.y>=this.min.y&&t.y<=this.max.y&&t.z>=this.min.z&&t.z<=this.max.z}containsBox(t){return this.min.x<=t.min.x&&t.max.x<=this.max.x&&this.min.y<=t.min.y&&t.max.y<=this.max.y&&this.min.z<=t.min.z&&t.max.z<=this.max.z}getParameter(t,e){return e.set((t.x-this.min.x)/(this.max.x-this.min.x),(t.y-this.min.y)/(this.max.y-this.min.y),(t.z-this.min.z)/(this.max.z-this.min.z))}intersectsBox(t){return t.max.x>=this.min.x&&t.min.x<=this.max.x&&t.max.y>=this.min.y&&t.min.y<=this.max.y&&t.max.z>=this.min.z&&t.min.z<=this.max.z}intersectsSphere(t){return this.clampPoint(t.center,vn),vn.distanceToSquared(t.center)<=t.radius*t.radius}intersectsPlane(t){let e,n;if(t.normal.x>0)e=t.normal.x*this.min.x,n=t.normal.x*this.max.x;else e=t.normal.x*this.max.x,n=t.normal.x*this.min.x;if(t.normal.y>0)e+=t.normal.y*this.min.y,n+=t.normal.y*this.max.y;else e+=t.normal.y*this.max.y,n+=t.normal.y*this.min.y;if(t.normal.z>0)e+=t.normal.z*this.min.z,n+=t.normal.z*this.max.z;else e+=t.normal.z*this.max.z,n+=t.normal.z*this.min.z;return e<=-t.constant&&n>=-t.constant}intersectsTriangle(t){if(this.isEmpty())return!1;this.getCenter(Os),Gr.subVectors(this.max,Os),es.subVectors(t.a,Os),ns.subVectors(t.b,Os),is.subVectors(t.c,Os),si.subVectors(ns,es),ri.subVectors(is,ns),vi.subVectors(es,is);let e=[0,-si.z,si.y,0,-ri.z,ri.y,0,-vi.z,vi.y,si.z,0,-si.x,ri.z,0,-ri.x,vi.z,0,-vi.x,-si.y,si.x,0,-ri.y,ri.x,0,-vi.y,vi.x,0];if(!jo(e,es,ns,is,Gr))return!1;if(e=[1,0,0,0,1,0,0,0,1],!jo(e,es,ns,is,Gr))return!1;return Hr.crossVectors(si,ri),e=[Hr.x,Hr.y,Hr.z],jo(e,es,ns,is,Gr)}clampPoint(t,e){return e.copy(t).clamp(this.min,this.max)}distanceToPoint(t){return this.clampPoint(t,vn).distanceTo(t)}getBoundingSphere(t){if(this.isEmpty())t.makeEmpty();else this.getCenter(t.center),t.radius=this.getSize(vn).length()*0.5;return t}intersect(t){if(this.min.max(t.min),this.max.min(t.max),this.isEmpty())this.makeEmpty();return this}union(t){return this.min.min(t.min),this.max.max(t.max),this}applyMatrix4(t){if(this.isEmpty())return this;return Hn[0].set(this.min.x,this.min.y,this.min.z).applyMatrix4(t),Hn[1].set(this.min.x,this.min.y,this.max.z).applyMatrix4(t),Hn[2].set(this.min.x,this.max.y,this.min.z).applyMatrix4(t),Hn[3].set(this.min.x,this.max.y,this.max.z).applyMatrix4(t),Hn[4].set(this.max.x,this.min.y,this.min.z).applyMatrix4(t),Hn[5].set(this.max.x,this.min.y,this.max.z).applyMatrix4(t),Hn[6].set(this.max.x,this.max.y,this.min.z).applyMatrix4(t),Hn[7].set(this.max.x,this.max.y,this.max.z).applyMatrix4(t),this.setFromPoints(Hn),this}translate(t){return this.min.add(t),this.max.add(t),this}equals(t){return t.min.equals(this.min)&&t.max.equals(this.max)}toJSON(){return{min:this.min.toArray(),max:this.max.toArray()}}fromJSON(t){return this.min.fromArray(t.min),this.max.fromArray(t.max),this}}var Hn=[new L,new L,new L,new L,new L,new L,new L,new L],vn=new L,kr=new Oe,es=new L,ns=new L,is=new L,si=new L,ri=new L,vi=new L,Os=new L,Gr=new L,Hr=new L,Si=new L;function jo(t,e,n,i,s){for(let r=0,a=t.length-3;r<=a;r+=3){Si.fromArray(t,r);let o=s.x*Math.abs(Si.x)+s.y*Math.abs(Si.y)+s.z*Math.abs(Si.z),l=e.dot(Si),c=n.dot(Si),h=i.dot(Si);if(Math.max(-Math.max(l,c,h),Math.min(l,c,h))>o)return!1}return!0}var Xn=e0();function e0(){let t=new ArrayBuffer(4),e=new Float32Array(t),n=new Uint32Array(t),i=new Uint32Array(512),s=new Uint32Array(512);for(let l=0;l<256;++l){let c=l-127;if(c<-27)i[l]=0,i[l|256]=32768,s[l]=24,s[l|256]=24;else if(c<-14)i[l]=1024>>-c-14,i[l|256]=1024>>-c-14|32768,s[l]=-c-1,s[l|256]=-c-1;else if(c<=15)i[l]=c+15<<10,i[l|256]=c+15<<10|32768,s[l]=13,s[l|256]=13;else if(c<128)i[l]=31744,i[l|256]=64512,s[l]=24,s[l|256]=24;else i[l]=31744,i[l|256]=64512,s[l]=13,s[l|256]=13}let r=new Uint32Array(2048),a=new Uint32Array(64),o=new Uint32Array(64);for(let l=1;l<1024;++l){let c=l<<13,h=0;while((c&8388608)===0)c<<=1,h-=8388608;c&=-8388609,h+=947912704,r[l]=c|h}for(let l=1024;l<2048;++l)r[l]=939524096+(l-1024<<13);for(let l=1;l<31;++l)a[l]=l<<23;a[31]=1199570944,a[32]=2147483648;for(let l=33;l<63;++l)a[l]=2147483648+(l-32<<23);a[63]=3347054592;for(let l=1;l<64;++l)if(l!==32)o[l]=1024;return{floatView:e,uint32View:n,baseTable:i,shiftTable:s,mantissaTable:r,exponentTable:a,offsetTable:o}}function Ke(t){if(Math.abs(t)>65504)Mt("DataUtils.toHalfFloat(): Value out of range.");t=Zt(t,-65504,65504),Xn.floatView[0]=t;let e=Xn.uint32View[0],n=e>>23&511;return Xn.baseTable[n]+((e&8388607)>>Xn.shiftTable[n])}function Zs(t){let e=t>>10;return Xn.uint32View[0]=Xn.mantissaTable[Xn.offsetTable[e]+(t&1023)]+Xn.exponentTable[e],Xn.floatView[0]}class cf{static toHalfFloat(t){return Ke(t)}static fromHalfFloat(t){return Zs(t)}}var Re=new L,Vr=new rt,n0=0;class de extends pn{constructor(t,e,n=!1){super();if(Array.isArray(t))throw TypeError("THREE.BufferAttribute: array should be a Typed Array.");this.isBufferAttribute=!0,Object.defineProperty(this,"id",{value:n0++}),this.name="",this.array=t,this.itemSize=e,this.count=t!==void 0?t.length/e:0,this.normalized=n,this.usage=35044,this.updateRanges=[],this.gpuType=1015,this.version=0}onUploadCallback(){}set needsUpdate(t){if(t===!0)this.version++}setUsage(t){return this.usage=t,this}addUpdateRange(t,e){this.updateRanges.push({start:t,count:e})}clearUpdateRanges(){this.updateRanges.length=0}copy(t){return this.name=t.name,this.array=new t.array.constructor(t.array),this.itemSize=t.itemSize,this.count=t.count,this.normalized=t.normalized,this.usage=t.usage,this.gpuType=t.gpuType,this}copyAt(t,e,n){t*=this.itemSize,n*=e.itemSize;for(let i=0,s=this.itemSize;i<s;i++)this.array[t+i]=e.array[n+i];return this}copyArray(t){return this.array.set(t),this}applyMatrix3(t){if(this.itemSize===2)for(let e=0,n=this.count;e<n;e++)Vr.fromBufferAttribute(this,e),Vr.applyMatrix3(t),this.setXY(e,Vr.x,Vr.y);else if(this.itemSize===3)for(let e=0,n=this.count;e<n;e++)Re.fromBufferAttribute(this,e),Re.applyMatrix3(t),this.setXYZ(e,Re.x,Re.y,Re.z);return this}applyMatrix4(t){for(let e=0,n=this.count;e<n;e++)Re.fromBufferAttribute(this,e),Re.applyMatrix4(t),this.setXYZ(e,Re.x,Re.y,Re.z);return this}applyNormalMatrix(t){for(let e=0,n=this.count;e<n;e++)Re.fromBufferAttribute(this,e),Re.applyNormalMatrix(t),this.setXYZ(e,Re.x,Re.y,Re.z);return this}transformDirection(t){for(let e=0,n=this.count;e<n;e++)Re.fromBufferAttribute(this,e),Re.transformDirection(t),this.setXYZ(e,Re.x,Re.y,Re.z);return this}set(t,e=0){return this.array.set(t,e),this}getComponent(t,e){let n=this.array[t*this.itemSize+e];if(this.normalized)n=Ze(n,this.array);return n}setComponent(t,e,n){if(this.normalized)n=te(n,this.array);return this.array[t*this.itemSize+e]=n,this}getX(t){let e=this.array[t*this.itemSize];if(this.normalized)e=Ze(e,this.array);return e}setX(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize]=e,this}getY(t){let e=this.array[t*this.itemSize+1];if(this.normalized)e=Ze(e,this.array);return e}setY(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize+1]=e,this}getZ(t){let e=this.array[t*this.itemSize+2];if(this.normalized)e=Ze(e,this.array);return e}setZ(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize+2]=e,this}getW(t){let e=this.array[t*this.itemSize+3];if(this.normalized)e=Ze(e,this.array);return e}setW(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize+3]=e,this}setXY(t,e,n){if(t*=this.itemSize,this.normalized)e=te(e,this.array),n=te(n,this.array);return this.array[t+0]=e,this.array[t+1]=n,this}setXYZ(t,e,n,i){if(t*=this.itemSize,this.normalized)e=te(e,this.array),n=te(n,this.array),i=te(i,this.array);return this.array[t+0]=e,this.array[t+1]=n,this.array[t+2]=i,this}setXYZW(t,e,n,i,s){if(t*=this.itemSize,this.normalized)e=te(e,this.array),n=te(n,this.array),i=te(i,this.array),s=te(s,this.array);return this.array[t+0]=e,this.array[t+1]=n,this.array[t+2]=i,this.array[t+3]=s,this}onUpload(t){return this.onUploadCallback=t,this}clone(){return new this.constructor(this.array,this.itemSize).copy(this)}toJSON(){let t={itemSize:this.itemSize,type:this.array.constructor.name,array:Array.from(this.array),normalized:this.normalized};return t.name=this.name,t.usage=this.usage,t.gpuType=this.gpuType,t}dispose(){this.dispatchEvent({type:"dispose"})}}class hf extends de{constructor(t,e,n){super(new Int8Array(t),e,n)}}class uf extends de{constructor(t,e,n){super(new Uint8Array(t),e,n)}}class df extends de{constructor(t,e,n){super(new Uint8ClampedArray(t),e,n)}}class ff extends de{constructor(t,e,n){super(new Int16Array(t),e,n)}}class Ya extends de{constructor(t,e,n){super(new Uint16Array(t),e,n)}}class pf extends de{constructor(t,e,n){super(new Int32Array(t),e,n)}}class Za extends de{constructor(t,e,n){super(new Uint32Array(t),e,n)}}class mf extends de{constructor(t,e,n){super(new Uint16Array(t),e,n);this.isFloat16BufferAttribute=!0}getX(t){let e=Zs(this.array[t*this.itemSize]);if(this.normalized)e=Ze(e,this.array);return e}setX(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize]=Ke(e),this}getY(t){let e=Zs(this.array[t*this.itemSize+1]);if(this.normalized)e=Ze(e,this.array);return e}setY(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize+1]=Ke(e),this}getZ(t){let e=Zs(this.array[t*this.itemSize+2]);if(this.normalized)e=Ze(e,this.array);return e}setZ(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize+2]=Ke(e),this}getW(t){let e=Zs(this.array[t*this.itemSize+3]);if(this.normalized)e=Ze(e,this.array);return e}setW(t,e){if(this.normalized)e=te(e,this.array);return this.array[t*this.itemSize+3]=Ke(e),this}setXY(t,e,n){if(t*=this.itemSize,this.normalized)e=te(e,this.array),n=te(n,this.array);return this.array[t+0]=Ke(e),this.array[t+1]=Ke(n),this}setXYZ(t,e,n,i){if(t*=this.itemSize,this.normalized)e=te(e,this.array),n=te(n,this.array),i=te(i,this.array);return this.array[t+0]=Ke(e),this.array[t+1]=Ke(n),this.array[t+2]=Ke(i),this}setXYZW(t,e,n,i,s){if(t*=this.itemSize,this.normalized)e=te(e,this.array),n=te(n,this.array),i=te(i,this.array),s=te(s,this.array);return this.array[t+0]=Ke(e),this.array[t+1]=Ke(n),this.array[t+2]=Ke(i),this.array[t+3]=Ke(s),this}}class Nt extends de{constructor(t,e,n){super(new Float32Array(t),e,n)}}var i0=new Oe,Bs=new L,Qo=new L;class Ue{constructor(t=new L,e=-1){this.isSphere=!0,this.center=t,this.radius=e}set(t,e){return this.center.copy(t),this.radius=e,this}setFromPoints(t,e){let n=this.center;if(e!==void 0)n.copy(e);else i0.setFromPoints(t).getCenter(n);let i=0;for(let s=0,r=t.length;s<r;s++)i=Math.max(i,n.distanceToSquared(t[s]));return this.radius=Math.sqrt(i),this}copy(t){return this.center.copy(t.center),this.radius=t.radius,this}isEmpty(){return this.radius<0}makeEmpty(){return this.center.set(0,0,0),this.radius=-1,this}containsPoint(t){return t.distanceToSquared(this.center)<=this.radius*this.radius}distanceToPoint(t){return t.distanceTo(this.center)-this.radius}intersectsSphere(t){let e=this.radius+t.radius;return t.center.distanceToSquared(this.center)<=e*e}intersectsBox(t){return t.intersectsSphere(this)}intersectsPlane(t){return Math.abs(t.distanceToPoint(this.center))<=this.radius}clampPoint(t,e){let n=this.center.distanceToSquared(t);if(e.copy(t),n>this.radius*this.radius)e.sub(this.center).normalize(),e.multiplyScalar(this.radius).add(this.center);return e}getBoundingBox(t){if(this.isEmpty())return t.makeEmpty(),t;return t.set(this.center,this.center),t.expandByScalar(this.radius),t}applyMatrix4(t){return this.center.applyMatrix4(t),this.radius=this.radius*t.getMaxScaleOnAxis(),this}translate(t){return this.center.add(t),this}expandByPoint(t){if(this.isEmpty())return this.center.copy(t),this.radius=0,this;Bs.subVectors(t,this.center);let e=Bs.lengthSq();if(e>this.radius*this.radius){let n=Math.sqrt(e),i=(n-this.radius)*0.5;this.center.addScaledVector(Bs,i/n),this.radius+=i}return this}union(t){if(t.isEmpty())return this;if(this.isEmpty())return this.copy(t),this;if(this.center.equals(t.center)===!0)this.radius=Math.max(this.radius,t.radius);else Qo.subVectors(t.center,this.center).setLength(t.radius),this.expandByPoint(Bs.copy(t.center).add(Qo)),this.expandByPoint(Bs.copy(t.center).sub(Qo));return this}equals(t){return t.center.equals(this.center)&&t.radius===this.radius}clone(){return new this.constructor().copy(this)}toJSON(){return{radius:this.radius,center:this.center.toArray()}}fromJSON(t){return this.radius=t.radius,this.center.fromArray(t.center),this}}var s0=0,un=new Kt,tl=new ae,ss=new L,nn=new Oe,zs=new Oe,Ne=new L;class jt extends pn{constructor(){super();this.isBufferGeometry=!0,Object.defineProperty(this,"id",{value:s0++}),this.uuid=rn(),this.name="",this.type="BufferGeometry",this.index=null,this.indirect=null,this.indirectOffset=0,this.attributes={},this.morphAttributes={},this.morphTargetsRelative=!1,this.groups=[],this.boundingBox=null,this.boundingSphere=null,this.drawRange={start:0,count:1/0},this.userData={},this._transformed=!1}getIndex(){return this.index}setIndex(t){if(Array.isArray(t))this.index=new((Eg(t))?Za:Ya)(t,1);else this.index=t;return this}setIndirect(t,e=0){return this.indirect=t,this.indirectOffset=e,this}getIndirect(){return this.indirect}getAttribute(t){return this.attributes[t]}setAttribute(t,e){return this.attributes[t]=e,this}deleteAttribute(t){return delete this.attributes[t],this}hasAttribute(t){return this.attributes[t]!==void 0}addGroup(t,e,n=0){this.groups.push({start:t,count:e,materialIndex:n})}clearGroups(){this.groups=[]}setDrawRange(t,e){this.drawRange.start=t,this.drawRange.count=e}applyMatrix4(t){let e=this.attributes.position;if(e!==void 0)e.applyMatrix4(t),e.needsUpdate=!0;let n=this.attributes.normal;if(n!==void 0){let s=new Xt().getNormalMatrix(t);n.applyNormalMatrix(s),n.needsUpdate=!0}let i=this.attributes.tangent;if(i!==void 0)i.transformDirection(t),i.needsUpdate=!0;if(this.boundingBox!==null)this.computeBoundingBox();if(this.boundingSphere!==null)this.computeBoundingSphere();return this._transformed=!0,this}applyQuaternion(t){return un.makeRotationFromQuaternion(t),this.applyMatrix4(un),this}rotateX(t){return un.makeRotationX(t),this.applyMatrix4(un),this}rotateY(t){return un.makeRotationY(t),this.applyMatrix4(un),this}rotateZ(t){return un.makeRotationZ(t),this.applyMatrix4(un),this}translate(t,e,n){return un.makeTranslation(t,e,n),this.applyMatrix4(un),this}scale(t,e,n){return un.makeScale(t,e,n),this.applyMatrix4(un),this}lookAt(t){return tl.lookAt(t),tl.updateMatrix(),this.applyMatrix4(tl.matrix),this}center(){return this.computeBoundingBox(),this.boundingBox.getCenter(ss).negate(),this.translate(ss.x,ss.y,ss.z),this}setFromPoints(t){let e=this.getAttribute("position");if(e===void 0){let n=[];for(let i=0,s=t.length;i<s;i++){let r=t[i];n.push(r.x,r.y,r.z||0)}this.setAttribute("position",new Nt(n,3))}else{let n=Math.min(t.length,e.count);for(let i=0;i<n;i++){let s=t[i];e.setXYZ(i,s.x,s.y,s.z||0)}if(t.length>e.count)Mt("BufferGeometry: Buffer size too small for points data. Use .dispose() and create a new geometry.");e.needsUpdate=!0}return this}computeBoundingBox(){if(this.boundingBox===null)this.boundingBox=new Oe;let t=this.attributes.position,e=this.morphAttributes.position;if(t&&t.isGLBufferAttribute){Wt("BufferGeometry.computeBoundingBox(): GLBufferAttribute requires a manual bounding box.",this),this.boundingBox.set(new L(-1/0,-1/0,-1/0),new L(1/0,1/0,1/0));return}if(t!==void 0){if(this.boundingBox.setFromBufferAttribute(t),e)for(let n=0,i=e.length;n<i;n++){let s=e[n];if(nn.setFromBufferAttribute(s),this.morphTargetsRelative)Ne.addVectors(this.boundingBox.min,nn.min),this.boundingBox.expandByPoint(Ne),Ne.addVectors(this.boundingBox.max,nn.max),this.boundingBox.expandByPoint(Ne);else this.boundingBox.expandByPoint(nn.min),this.boundingBox.expandByPoint(nn.max)}}else this.boundingBox.makeEmpty();if(isNaN(this.boundingBox.min.x)||isNaN(this.boundingBox.min.y)||isNaN(this.boundingBox.min.z))Wt('BufferGeometry.computeBoundingBox(): Computed min/max have NaN values. The "position" attribute is likely to have NaN values.',this)}computeBoundingSphere(){if(this.boundingSphere===null)this.boundingSphere=new Ue;let t=this.attributes.position,e=this.morphAttributes.position;if(t&&t.isGLBufferAttribute){Wt("BufferGeometry.computeBoundingSphere(): GLBufferAttribute requires a manual bounding sphere.",this),this.boundingSphere.set(new L,1/0);return}if(t){let n=this.boundingSphere.center;if(nn.setFromBufferAttribute(t),e)for(let s=0,r=e.length;s<r;s++){let a=e[s];if(zs.setFromBufferAttribute(a),this.morphTargetsRelative)Ne.addVectors(nn.min,zs.min),nn.expandByPoint(Ne),Ne.addVectors(nn.max,zs.max),nn.expandByPoint(Ne);else nn.expandByPoint(zs.min),nn.expandByPoint(zs.max)}nn.getCenter(n);let i=0;for(let s=0,r=t.count;s<r;s++)Ne.fromBufferAttribute(t,s),i=Math.max(i,n.distanceToSquared(Ne));if(e)for(let s=0,r=e.length;s<r;s++){let a=e[s],o=this.morphTargetsRelative;for(let l=0,c=a.count;l<c;l++){if(Ne.fromBufferAttribute(a,l),o)ss.fromBufferAttribute(t,l),Ne.add(ss);i=Math.max(i,n.distanceToSquared(Ne))}}if(this.boundingSphere.radius=Math.sqrt(i),isNaN(this.boundingSphere.radius))Wt('BufferGeometry.computeBoundingSphere(): Computed radius is NaN. The "position" attribute is likely to have NaN values.',this)}}computeTangents(){let t=this.index,e=this.attributes;if(t===null||e.position===void 0||e.normal===void 0||e.uv===void 0){Wt("BufferGeometry: .computeTangents() failed. Missing required attributes (index, position, normal or uv)");return}let{position:n,normal:i,uv:s}=e,r=this.getAttribute("tangent");if(r===void 0||r.count!==n.count)r=new de(new Float32Array(4*n.count),4),this.setAttribute("tangent",r);let a=[],o=[];for(let R=0;R<n.count;R++)a[R]=new L,o[R]=new L;let l=new L,c=new L,h=new L,d=new rt,u=new rt,f=new rt,g=new L,_=new L;function m(R,y,A){l.fromBufferAttribute(n,R),c.fromBufferAttribute(n,y),h.fromBufferAttribute(n,A),d.fromBufferAttribute(s,R),u.fromBufferAttribute(s,y),f.fromBufferAttribute(s,A),c.sub(l),h.sub(l),u.sub(d),f.sub(d);let Z=1/(u.x*f.y-f.x*u.y);if(!isFinite(Z))return;g.copy(c).multiplyScalar(f.y).addScaledVector(h,-u.y).multiplyScalar(Z),_.copy(h).multiplyScalar(u.x).addScaledVector(c,-f.x).multiplyScalar(Z),a[R].add(g),a[y].add(g),a[A].add(g),o[R].add(_),o[y].add(_),o[A].add(_)}let p=this.groups;if(p.length===0)p=[{start:0,count:t.count}];for(let R=0,y=p.length;R<y;++R){let A=p[R],{start:Z,count:D}=A;for(let q=Z,ot=Z+D;q<ot;q+=3)m(t.getX(q+0),t.getX(q+1),t.getX(q+2))}let b=new L,S=new L,x=new L,T=new L;function I(R){x.fromBufferAttribute(i,R),T.copy(x);let y=a[R];b.copy(y),b.sub(x.multiplyScalar(x.dot(y))).normalize(),S.crossVectors(T,y);let Z=S.dot(o[R])<0?-1:1;r.setXYZW(R,b.x,b.y,b.z,Z)}for(let R=0,y=p.length;R<y;++R){let A=p[R],{start:Z,count:D}=A;for(let q=Z,ot=Z+D;q<ot;q+=3)I(t.getX(q+0)),I(t.getX(q+1)),I(t.getX(q+2))}this._transformed=!0}computeVertexNormals(){let t=this.index,e=this.getAttribute("position");if(e!==void 0){let n=this.getAttribute("normal");if(n===void 0||n.count!==e.count)n=new de(new Float32Array(e.count*3),3),this.setAttribute("normal",n);else for(let d=0,u=n.count;d<u;d++)n.setXYZ(d,0,0,0);let i=new L,s=new L,r=new L,a=new L,o=new L,l=new L,c=new L,h=new L;if(t)for(let d=0,u=t.count;d<u;d+=3){let f=t.getX(d+0),g=t.getX(d+1),_=t.getX(d+2);i.fromBufferAttribute(e,f),s.fromBufferAttribute(e,g),r.fromBufferAttribute(e,_),c.subVectors(r,s),h.subVectors(i,s),c.cross(h),a.fromBufferAttribute(n,f),o.fromBufferAttribute(n,g),l.fromBufferAttribute(n,_),a.add(c),o.add(c),l.add(c),n.setXYZ(f,a.x,a.y,a.z),n.setXYZ(g,o.x,o.y,o.z),n.setXYZ(_,l.x,l.y,l.z)}else for(let d=0,u=e.count;d<u;d+=3)i.fromBufferAttribute(e,d+0),s.fromBufferAttribute(e,d+1),r.fromBufferAttribute(e,d+2),c.subVectors(r,s),h.subVectors(i,s),c.cross(h),n.setXYZ(d+0,c.x,c.y,c.z),n.setXYZ(d+1,c.x,c.y,c.z),n.setXYZ(d+2,c.x,c.y,c.z);this.normalizeNormals(),n.needsUpdate=!0}}normalizeNormals(){let t=this.attributes.normal;for(let e=0,n=t.count;e<n;e++)Ne.fromBufferAttribute(t,e),Ne.normalize(),t.setXYZ(e,Ne.x,Ne.y,Ne.z)}toNonIndexed(){function t(a,o){let{array:l,itemSize:c,normalized:h}=a,d=new l.constructor(o.length*c),u=0,f=0;for(let g=0,_=o.length;g<_;g++){if(a.isInterleavedBufferAttribute)u=o[g]*a.data.stride+a.offset;else u=o[g]*c;for(let m=0;m<c;m++)d[f++]=l[u++]}return new de(d,c,h)}if(this.index===null)return Mt("BufferGeometry.toNonIndexed(): BufferGeometry is already non-indexed."),this;let e=new jt,n=this.index.array,i=this.attributes;for(let a in i){let o=i[a],l=t(o,n);e.setAttribute(a,l)}let s=this.morphAttributes;for(let a in s){let o=[],l=s[a];for(let c=0,h=l.length;c<h;c++){let d=l[c],u=t(d,n);o.push(u)}e.morphAttributes[a]=o}e.morphTargetsRelative=this.morphTargetsRelative;let r=this.groups;for(let a=0,o=r.length;a<o;a++){let l=r[a];e.addGroup(l.start,l.count,l.materialIndex)}return e}toJSON(){let t={metadata:{version:4.7,type:"BufferGeometry",generator:"BufferGeometry.toJSON"}};if(t.uuid=this.uuid,t.type=this.parameters!==void 0&&this._transformed===!0?"BufferGeometry":this.type,t.name=this.name,Object.keys(this.userData).length>0)t.userData=this.userData;if(this.parameters!==void 0&&this._transformed!==!0){let o=this.parameters;for(let l in o)if(o[l]!==void 0)t[l]=o[l];return t}t.data={attributes:{}};let e=this.index;if(e!==null)t.data.index={type:e.array.constructor.name,array:Array.prototype.slice.call(e.array)};let n=this.attributes;for(let o in n){let l=n[o];t.data.attributes[o]=l.toJSON(t.data)}let i={},s=!1;for(let o in this.morphAttributes){let l=this.morphAttributes[o],c=[];for(let h=0,d=l.length;h<d;h++){let u=l[h];c.push(u.toJSON(t.data))}if(c.length>0)i[o]=c,s=!0}if(s)t.data.morphAttributes=i,t.data.morphTargetsRelative=this.morphTargetsRelative;let r=this.groups;if(r.length>0)t.data.groups=JSON.parse(JSON.stringify(r));let a=this.boundingSphere;if(a!==null)t.data.boundingSphere=a.toJSON();return t}clone(){return new this.constructor().copy(this)}copy(t){this.index=null,this.attributes={},this.morphAttributes={},this.groups=[],this.boundingBox=null,this.boundingSphere=null;let e={};this.name=t.name;let n=t.index;if(n!==null)this.setIndex(n.clone());let i=t.attributes;for(let l in i){let c=i[l];this.setAttribute(l,c.clone(e))}let s=t.morphAttributes;for(let l in s){let c=[],h=s[l];for(let d=0,u=h.length;d<u;d++)c.push(h[d].clone(e));this.morphAttributes[l]=c}this.morphTargetsRelative=t.morphTargetsRelative;let r=t.groups;for(let l=0,c=r.length;l<c;l++){let h=r[l];this.addGroup(h.start,h.count,h.materialIndex)}let a=t.boundingBox;if(a!==null)this.boundingBox=a.clone();let o=t.boundingSphere;if(o!==null)this.boundingSphere=o.clone();return this.drawRange.start=t.drawRange.start,this.drawRange.count=t.drawRange.count,this.userData=t.userData,this._transformed=t._transformed,this}dispose(){this.dispatchEvent({type:"dispose"})}}class ws{constructor(t,e){this.isInterleavedBuffer=!0,this.array=t,this.stride=e,this.count=t!==void 0?t.length/e:0,this.usage=35044,this.updateRanges=[],this.version=0,this.uuid=rn()}onUploadCallback(){}set needsUpdate(t){if(t===!0)this.version++}setUsage(t){return this.usage=t,this}addUpdateRange(t,e){this.updateRanges.push({start:t,count:e})}clearUpdateRanges(){this.updateRanges.length=0}copy(t){return this.array=new t.array.constructor(t.array),this.count=t.count,this.stride=t.stride,this.usage=t.usage,this}copyAt(t,e,n){t*=this.stride,n*=e.stride;for(let i=0,s=this.stride;i<s;i++)this.array[t+i]=e.array[n+i];return this}set(t,e=0){return this.array.set(t,e),this}clone(t){if(t.arrayBuffers===void 0)t.arrayBuffers={};if(this.array.buffer._uuid===void 0)this.array.buffer._uuid=rn();if(t.arrayBuffers[this.array.buffer._uuid]===void 0)t.arrayBuffers[this.array.buffer._uuid]=this.array.slice(0).buffer;let e=new this.array.constructor(t.arrayBuffers[this.array.buffer._uuid]),n=new this.constructor(e,this.stride);return n.setUsage(this.usage),n}onUpload(t){return this.onUploadCallback=t,this}toJSON(t){if(t.arrayBuffers===void 0)t.arrayBuffers={};if(this.array.buffer._uuid===void 0)this.array.buffer._uuid=rn();if(t.arrayBuffers[this.array.buffer._uuid]===void 0)t.arrayBuffers[this.array.buffer._uuid]=Array.from(new Uint32Array(this.array.buffer));let e={uuid:this.uuid,buffer:this.array.buffer._uuid,type:this.array.constructor.name,stride:this.stride};return e.usage=this.usage,e}}var Ye=new L;class ci{constructor(t,e,n,i=!1){this.isInterleavedBufferAttribute=!0,this.name="",this.data=t,this.itemSize=e,this.offset=n,this.normalized=i}get count(){return this.data.count}get array(){return this.data.array}set needsUpdate(t){this.data.needsUpdate=t}applyMatrix4(t){for(let e=0,n=this.data.count;e<n;e++)Ye.fromBufferAttribute(this,e),Ye.applyMatrix4(t),this.setXYZ(e,Ye.x,Ye.y,Ye.z);return this}applyNormalMatrix(t){for(let e=0,n=this.count;e<n;e++)Ye.fromBufferAttribute(this,e),Ye.applyNormalMatrix(t),this.setXYZ(e,Ye.x,Ye.y,Ye.z);return this}transformDirection(t){for(let e=0,n=this.count;e<n;e++)Ye.fromBufferAttribute(this,e),Ye.transformDirection(t),this.setXYZ(e,Ye.x,Ye.y,Ye.z);return this}getComponent(t,e){let n=this.array[t*this.data.stride+this.offset+e];if(this.normalized)n=Ze(n,this.array);return n}setComponent(t,e,n){if(this.normalized)n=te(n,this.array);return this.data.array[t*this.data.stride+this.offset+e]=n,this}setX(t,e){if(this.normalized)e=te(e,this.array);return this.data.array[t*this.data.stride+this.offset]=e,this}setY(t,e){if(this.normalized)e=te(e,this.array);return this.data.array[t*this.data.stride+this.offset+1]=e,this}setZ(t,e){if(this.normalized)e=te(e,this.array);return this.data.array[t*this.data.stride+this.offset+2]=e,this}setW(t,e){if(this.normalized)e=te(e,this.array);return this.data.array[t*this.data.stride+this.offset+3]=e,this}getX(t){let e=this.data.array[t*this.data.stride+this.offset];if(this.normalized)e=Ze(e,this.array);return e}getY(t){let e=this.data.array[t*this.data.stride+this.offset+1];if(this.normalized)e=Ze(e,this.array);return e}getZ(t){let e=this.data.array[t*this.data.stride+this.offset+2];if(this.normalized)e=Ze(e,this.array);return e}getW(t){let e=this.data.array[t*this.data.stride+this.offset+3];if(this.normalized)e=Ze(e,this.array);return e}setXY(t,e,n){if(t=t*this.data.stride+this.offset,this.normalized)e=te(e,this.array),n=te(n,this.array);return this.data.array[t+0]=e,this.data.array[t+1]=n,this}setXYZ(t,e,n,i){if(t=t*this.data.stride+this.offset,this.normalized)e=te(e,this.array),n=te(n,this.array),i=te(i,this.array);return this.data.array[t+0]=e,this.data.array[t+1]=n,this.data.array[t+2]=i,this}setXYZW(t,e,n,i,s){if(t=t*this.data.stride+this.offset,this.normalized)e=te(e,this.array),n=te(n,this.array),i=te(i,this.array),s=te(s,this.array);return this.data.array[t+0]=e,this.data.array[t+1]=n,this.data.array[t+2]=i,this.data.array[t+3]=s,this}clone(t){if(t===void 0){tr("InterleavedBufferAttribute.clone(): Cloning an interleaved buffer attribute will de-interleave buffer data.");let e=[];for(let n=0;n<this.count;n++){let i=n*this.data.stride+this.offset;for(let s=0;s<this.itemSize;s++)e.push(this.data.array[i+s])}return new de(new this.array.constructor(e),this.itemSize,this.normalized)}else{if(t.interleavedBuffers===void 0)t.interleavedBuffers={};if(t.interleavedBuffers[this.data.uuid]===void 0)t.interleavedBuffers[this.data.uuid]=this.data.clone(t);return new ci(t.interleavedBuffers[this.data.uuid],this.itemSize,this.offset,this.normalized)}}toJSON(t){if(t===void 0){tr("InterleavedBufferAttribute.toJSON(): Serializing an interleaved buffer attribute will de-interleave buffer data.");let e=[];for(let n=0;n<this.count;n++){let i=n*this.data.stride+this.offset;for(let s=0;s<this.itemSize;s++)e.push(this.data.array[i+s])}return{itemSize:this.itemSize,type:this.array.constructor.name,array:e,normalized:this.normalized}}else{if(t.interleavedBuffers===void 0)t.interleavedBuffers={};if(t.interleavedBuffers[this.data.uuid]===void 0)t.interleavedBuffers[this.data.uuid]=this.data.toJSON(t);return{isInterleavedBufferAttribute:!0,itemSize:this.itemSize,data:this.data.uuid,offset:this.offset,normalized:this.normalized}}}}var el=new L,r0=new L,a0=new Xt;class dn{constructor(t=new L(1,0,0),e=0){this.isPlane=!0,this.normal=t,this.constant=e}set(t,e){return this.normal.copy(t),this.constant=e,this}setComponents(t,e,n,i){return this.normal.set(t,e,n),this.constant=i,this}setFromNormalAndCoplanarPoint(t,e){return this.normal.copy(t),this.constant=-e.dot(this.normal),this}setFromCoplanarPoints(t,e,n){let i=el.subVectors(n,e).cross(r0.subVectors(t,e)).normalize();return this.setFromNormalAndCoplanarPoint(i,t),this}copy(t){return this.normal.copy(t.normal),this.constant=t.constant,this}normalize(){let t=1/this.normal.length();return this.normal.multiplyScalar(t),this.constant*=t,this}negate(){return this.constant*=-1,this.normal.negate(),this}distanceToPoint(t){return this.normal.dot(t)+this.constant}distanceToSphere(t){return this.distanceToPoint(t.center)-t.radius}projectPoint(t,e){return e.copy(t).addScaledVector(this.normal,-this.distanceToPoint(t))}intersectLine(t,e,n=!0){let i=t.delta(el),s=this.normal.dot(i);if(s===0){if(this.distanceToPoint(t.start)===0)return e.copy(t.start);return null}let r=-(t.start.dot(this.normal)+this.constant)/s;if(n===!0&&(r<0||r>1))return null;return e.copy(t.start).addScaledVector(i,r)}intersectsLine(t){let e=this.distanceToPoint(t.start),n=this.distanceToPoint(t.end);return e<0&&n>0||n<0&&e>0}intersectsBox(t){return t.intersectsPlane(this)}intersectsSphere(t){return t.intersectsPlane(this)}coplanarPoint(t){return t.copy(this.normal).multiplyScalar(-this.constant)}applyMatrix4(t,e){let n=e||a0.getNormalMatrix(t),i=this.coplanarPoint(el).applyMatrix4(t),s=this.normal.applyMatrix3(n).normalize();return this.constant=-i.dot(s),this}translate(t){return this.constant-=t.dot(this.normal),this}equals(t){return t.normal.equals(this.normal)&&t.constant===this.constant}clone(){return new this.constructor().copy(this)}toJSON(){return{normal:this.normal.toArray(),constant:this.constant}}fromJSON(t){return this.normal.fromArray(t.normal),this.constant=t.constant,this}}var o0=0;class Fe extends pn{constructor(){super();this.isMaterial=!0,Object.defineProperty(this,"id",{value:o0++}),this.uuid=rn(),this.name="",this.type="Material",this.blending=1,this.side=0,this.vertexColors=!1,this.opacity=1,this.transparent=!1,this.alphaHash=!1,this.blendSrc=204,this.blendDst=205,this.blendEquation=100,this.blendSrcAlpha=null,this.blendDstAlpha=null,this.blendEquationAlpha=null,this.blendColor=new Pt(0,0,0),this.blendAlpha=0,this.depthFunc=3,this.depthTest=!0,this.depthWrite=!0,this.stencilWriteMask=255,this.stencilFunc=519,this.stencilRef=0,this.stencilFuncMask=255,this.stencilFail=7680,this.stencilZFail=7680,this.stencilZPass=7680,this.stencilWrite=!1,this.clippingPlanes=null,this.clipIntersection=!1,this.clipShadows=!1,this.shadowSide=null,this.colorWrite=!0,this.precision=null,this.polygonOffset=!1,this.polygonOffsetFactor=0,this.polygonOffsetUnits=0,this.dithering=!1,this.alphaToCoverage=!1,this.premultipliedAlpha=!1,this.forceSinglePass=!1,this.allowOverride=!0,this.visible=!0,this.toneMapped=!0,this.userData={},this.version=0,this._alphaTest=0}get alphaTest(){return this._alphaTest}set alphaTest(t){if(this._alphaTest>0!==t>0)this.version++;this._alphaTest=t}onBeforeRender(){}onBeforeCompile(){}customProgramCacheKey(){return this.onBeforeCompile.toString()}setValues(t){if(t===void 0)return;for(let e in t){let n=t[e];if(n===void 0){Mt(`Material: parameter '${e}' has value of undefined.`);continue}let i=this[e];if(i===void 0){Mt(`Material: '${e}' is not a property of THREE.${this.type}.`);continue}if(i&&i.isColor)i.set(n);else if(i&&i.isVector2&&(n&&n.isVector2)||i&&i.isEuler&&(n&&n.isEuler)||i&&i.isVector3&&(n&&n.isVector3))i.copy(n);else this[e]=n}}toJSON(t){let e=t===void 0||typeof t==="string";if(e)t={textures:{},images:{}};let n={metadata:{version:4.7,type:"Material",generator:"Material.toJSON"}};if(n.uuid=this.uuid,n.type=this.type,n.blending=this.blending,n.side=this.side,n.shadowSide=this.shadowSide,n.vertexColors=this.vertexColors,n.opacity=this.opacity,n.transparent=this.transparent,n.blendSrc=this.blendSrc,n.blendDst=this.blendDst,n.blendEquation=this.blendEquation,n.blendSrcAlpha=this.blendSrcAlpha,n.blendDstAlpha=this.blendDstAlpha,n.blendEquationAlpha=this.blendEquationAlpha,n.blendColor=this.blendColor.getHex(),n.blendAlpha=this.blendAlpha,n.depthFunc=this.depthFunc,n.depthTest=this.depthTest,n.depthWrite=this.depthWrite,n.colorWrite=this.colorWrite,n.clipIntersection=this.clipIntersection,n.clipShadows=this.clipShadows,n.stencilWriteMask=this.stencilWriteMask,n.stencilFunc=this.stencilFunc,n.stencilRef=this.stencilRef,n.stencilFuncMask=this.stencilFuncMask,n.stencilFail=this.stencilFail,n.stencilZFail=this.stencilZFail,n.stencilZPass=this.stencilZPass,n.stencilWrite=this.stencilWrite,n.polygonOffset=this.polygonOffset,n.polygonOffsetFactor=this.polygonOffsetFactor,n.polygonOffsetUnits=this.polygonOffsetUnits,n.dithering=this.dithering,n.alphaTest=this.alphaTest,n.alphaHash=this.alphaHash,n.alphaToCoverage=this.alphaToCoverage,n.premultipliedAlpha=this.premultipliedAlpha,n.forceSinglePass=this.forceSinglePass,n.allowOverride=this.allowOverride,n.visible=this.visible,n.toneMapped=this.toneMapped,n.name=this.name,this.color&&this.color.isColor)n.color=this.color.getHex();if(this.roughness!==void 0)n.roughness=this.roughness;if(this.metalness!==void 0)n.metalness=this.metalness;if(this.sheen!==void 0)n.sheen=this.sheen;if(this.sheenColor&&this.sheenColor.isColor)n.sheenColor=this.sheenColor.getHex();if(this.sheenRoughness!==void 0)n.sheenRoughness=this.sheenRoughness;if(this.emissive&&this.emissive.isColor)n.emissive=this.emissive.getHex();if(this.emissiveIntensity!==void 0)n.emissiveIntensity=this.emissiveIntensity;if(this.specular&&this.specular.isColor)n.specular=this.specular.getHex();if(this.specularIntensity!==void 0)n.specularIntensity=this.specularIntensity;if(this.specularColor&&this.specularColor.isColor)n.specularColor=this.specularColor.getHex();if(this.shininess!==void 0)n.shininess=this.shininess;if(this.clearcoat!==void 0)n.clearcoat=this.clearcoat;if(this.clearcoatRoughness!==void 0)n.clearcoatRoughness=this.clearcoatRoughness;if(this.clearcoatMap&&this.clearcoatMap.isTexture)n.clearcoatMap=this.clearcoatMap.toJSON(t).uuid;if(this.clearcoatRoughnessMap&&this.clearcoatRoughnessMap.isTexture)n.clearcoatRoughnessMap=this.clearcoatRoughnessMap.toJSON(t).uuid;if(this.clearcoatNormalMap&&this.clearcoatNormalMap.isTexture)n.clearcoatNormalMap=this.clearcoatNormalMap.toJSON(t).uuid,n.clearcoatNormalScale=this.clearcoatNormalScale.toArray();if(this.sheenColorMap&&this.sheenColorMap.isTexture)n.sheenColorMap=this.sheenColorMap.toJSON(t).uuid;if(this.sheenRoughnessMap&&this.sheenRoughnessMap.isTexture)n.sheenRoughnessMap=this.sheenRoughnessMap.toJSON(t).uuid;if(this.dispersion!==void 0)n.dispersion=this.dispersion;if(this.retroreflectivity!==void 0)n.retroreflectivity=this.retroreflectivity;if(this.iridescence!==void 0)n.iridescence=this.iridescence;if(this.iridescenceIOR!==void 0)n.iridescenceIOR=this.iridescenceIOR;if(this.iridescenceThicknessRange!==void 0)n.iridescenceThicknessRange=this.iridescenceThicknessRange;if(this.iridescenceMap&&this.iridescenceMap.isTexture)n.iridescenceMap=this.iridescenceMap.toJSON(t).uuid;if(this.iridescenceThicknessMap&&this.iridescenceThicknessMap.isTexture)n.iridescenceThicknessMap=this.iridescenceThicknessMap.toJSON(t).uuid;if(this.anisotropy!==void 0)n.anisotropy=this.anisotropy;if(this.anisotropyRotation!==void 0)n.anisotropyRotation=this.anisotropyRotation;if(this.anisotropyMap&&this.anisotropyMap.isTexture)n.anisotropyMap=this.anisotropyMap.toJSON(t).uuid;if(this.map&&this.map.isTexture)n.map=this.map.toJSON(t).uuid;if(this.matcap&&this.matcap.isTexture)n.matcap=this.matcap.toJSON(t).uuid;if(this.alphaMap&&this.alphaMap.isTexture)n.alphaMap=this.alphaMap.toJSON(t).uuid;if(this.lightMap&&this.lightMap.isTexture)n.lightMap=this.lightMap.toJSON(t).uuid,n.lightMapIntensity=this.lightMapIntensity;if(this.aoMap&&this.aoMap.isTexture)n.aoMap=this.aoMap.toJSON(t).uuid,n.aoMapIntensity=this.aoMapIntensity;if(this.bumpMap&&this.bumpMap.isTexture)n.bumpMap=this.bumpMap.toJSON(t).uuid,n.bumpScale=this.bumpScale;if(this.normalMap&&this.normalMap.isTexture)n.normalMap=this.normalMap.toJSON(t).uuid,n.normalMapType=this.normalMapType,n.normalScale=this.normalScale.toArray();if(this.displacementMap&&this.displacementMap.isTexture)n.displacementMap=this.displacementMap.toJSON(t).uuid,n.displacementScale=this.displacementScale,n.displacementBias=this.displacementBias;if(this.roughnessMap&&this.roughnessMap.isTexture)n.roughnessMap=this.roughnessMap.toJSON(t).uuid;if(this.metalnessMap&&this.metalnessMap.isTexture)n.metalnessMap=this.metalnessMap.toJSON(t).uuid;if(this.emissiveMap&&this.emissiveMap.isTexture)n.emissiveMap=this.emissiveMap.toJSON(t).uuid;if(this.specularMap&&this.specularMap.isTexture)n.specularMap=this.specularMap.toJSON(t).uuid;if(this.specularIntensityMap&&this.specularIntensityMap.isTexture)n.specularIntensityMap=this.specularIntensityMap.toJSON(t).uuid;if(this.specularColorMap&&this.specularColorMap.isTexture)n.specularColorMap=this.specularColorMap.toJSON(t).uuid;if(this.envMap&&this.envMap.isTexture){if(n.envMap=this.envMap.toJSON(t).uuid,this.combine!==void 0)n.combine=this.combine}if(this.envMapRotation!==void 0)n.envMapRotation=this.envMapRotation.toArray();if(this.envMapIntensity!==void 0)n.envMapIntensity=this.envMapIntensity;if(this.reflectivity!==void 0)n.reflectivity=this.reflectivity;if(this.refractionRatio!==void 0)n.refractionRatio=this.refractionRatio;if(this.gradientMap&&this.gradientMap.isTexture)n.gradientMap=this.gradientMap.toJSON(t).uuid;if(this.transmission!==void 0)n.transmission=this.transmission;if(this.transmissionMap&&this.transmissionMap.isTexture)n.transmissionMap=this.transmissionMap.toJSON(t).uuid;if(this.thickness!==void 0)n.thickness=this.thickness;if(this.thicknessMap&&this.thicknessMap.isTexture)n.thicknessMap=this.thicknessMap.toJSON(t).uuid;if(this.attenuationDistance!==void 0)n.attenuationDistance=this.attenuationDistance;if(this.attenuationColor!==void 0)n.attenuationColor=this.attenuationColor.getHex();if(this.size!==void 0)n.size=this.size;if(this.sizeAttenuation!==void 0)n.sizeAttenuation=this.sizeAttenuation;if(Array.isArray(this.clippingPlanes)&&this.clippingPlanes.length>0)n.clippingPlanes=this.clippingPlanes.map((s)=>s.toJSON());if(this.rotation!==void 0)n.rotation=this.rotation;if(this.depthPacking!==void 0)n.depthPacking=this.depthPacking;if(this.linewidth!==void 0)n.linewidth=this.linewidth;if(this.linecap!==void 0)n.linecap=this.linecap;if(this.linejoin!==void 0)n.linejoin=this.linejoin;if(this.dashSize!==void 0)n.dashSize=this.dashSize;if(this.gapSize!==void 0)n.gapSize=this.gapSize;if(this.scale!==void 0)n.scale=this.scale;if(this.wireframe!==void 0)n.wireframe=this.wireframe;if(this.wireframeLinewidth!==void 0)n.wireframeLinewidth=this.wireframeLinewidth;if(this.wireframeLinecap!==void 0)n.wireframeLinecap=this.wireframeLinecap;if(this.wireframeLinejoin!==void 0)n.wireframeLinejoin=this.wireframeLinejoin;if(this.flatShading!==void 0)n.flatShading=this.flatShading;if(this.fog!==void 0)n.fog=this.fog;if(Object.keys(this.userData).length>0)n.userData=this.userData;function i(s){let r=[];for(let a in s){let o=s[a];delete o.metadata,r.push(o)}return r}if(e){let s=i(t.textures),r=i(t.images);if(s.length>0)n.textures=s;if(r.length>0)n.images=r}return n}fromJSON(t,e){if(t.uuid!==void 0)this.uuid=t.uuid;if(t.name!==void 0)this.name=t.name;if(t.color!==void 0&&this.color!==void 0)this.color.setHex(t.color);if(t.roughness!==void 0)this.roughness=t.roughness;if(t.metalness!==void 0)this.metalness=t.metalness;if(t.sheen!==void 0)this.sheen=t.sheen;if(t.sheenColor!==void 0)this.sheenColor=new Pt().setHex(t.sheenColor);if(t.sheenRoughness!==void 0)this.sheenRoughness=t.sheenRoughness;if(t.emissive!==void 0&&this.emissive!==void 0)this.emissive.setHex(t.emissive);if(t.specular!==void 0&&this.specular!==void 0)this.specular.setHex(t.specular);if(t.specularIntensity!==void 0)this.specularIntensity=t.specularIntensity;if(t.specularColor!==void 0&&this.specularColor!==void 0)this.specularColor.setHex(t.specularColor);if(t.shininess!==void 0)this.shininess=t.shininess;if(t.clearcoat!==void 0)this.clearcoat=t.clearcoat;if(t.clearcoatRoughness!==void 0)this.clearcoatRoughness=t.clearcoatRoughness;if(t.dispersion!==void 0)this.dispersion=t.dispersion;if(t.retroreflectivity!==void 0)this.retroreflectivity=t.retroreflectivity;if(t.iridescence!==void 0)this.iridescence=t.iridescence;if(t.iridescenceIOR!==void 0)this.iridescenceIOR=t.iridescenceIOR;if(t.iridescenceThicknessRange!==void 0)this.iridescenceThicknessRange=t.iridescenceThicknessRange;if(t.transmission!==void 0)this.transmission=t.transmission;if(t.thickness!==void 0)this.thickness=t.thickness;if(t.attenuationDistance!==void 0)this.attenuationDistance=t.attenuationDistance;if(t.attenuationColor!==void 0&&this.attenuationColor!==void 0)this.attenuationColor.setHex(t.attenuationColor);if(t.anisotropy!==void 0)this.anisotropy=t.anisotropy;if(t.anisotropyRotation!==void 0)this.anisotropyRotation=t.anisotropyRotation;if(t.fog!==void 0)this.fog=t.fog;if(t.flatShading!==void 0)this.flatShading=t.flatShading;if(t.blending!==void 0)this.blending=t.blending;if(t.combine!==void 0)this.combine=t.combine;if(t.side!==void 0)this.side=t.side;if(t.shadowSide!==void 0)this.shadowSide=t.shadowSide;if(t.opacity!==void 0)this.opacity=t.opacity;if(t.transparent!==void 0)this.transparent=t.transparent;if(t.alphaTest!==void 0)this.alphaTest=t.alphaTest;if(t.alphaHash!==void 0)this.alphaHash=t.alphaHash;if(t.depthFunc!==void 0)this.depthFunc=t.depthFunc;if(t.depthTest!==void 0)this.depthTest=t.depthTest;if(t.depthWrite!==void 0)this.depthWrite=t.depthWrite;if(t.colorWrite!==void 0)this.colorWrite=t.colorWrite;if(t.clippingPlanes!==void 0)this.clippingPlanes=t.clippingPlanes.map((n)=>new dn().fromJSON(n));if(t.clipIntersection!==void 0)this.clipIntersection=t.clipIntersection;if(t.clipShadows!==void 0)this.clipShadows=t.clipShadows;if(t.depthPacking!==void 0)this.depthPacking=t.depthPacking;if(t.blendSrc!==void 0)this.blendSrc=t.blendSrc;if(t.blendDst!==void 0)this.blendDst=t.blendDst;if(t.blendEquation!==void 0)this.blendEquation=t.blendEquation;if(t.blendSrcAlpha!==void 0)this.blendSrcAlpha=t.blendSrcAlpha;if(t.blendDstAlpha!==void 0)this.blendDstAlpha=t.blendDstAlpha;if(t.blendEquationAlpha!==void 0)this.blendEquationAlpha=t.blendEquationAlpha;if(t.blendColor!==void 0&&this.blendColor!==void 0)this.blendColor.setHex(t.blendColor);if(t.blendAlpha!==void 0)this.blendAlpha=t.blendAlpha;if(t.stencilWriteMask!==void 0)this.stencilWriteMask=t.stencilWriteMask;if(t.stencilFunc!==void 0)this.stencilFunc=t.stencilFunc;if(t.stencilRef!==void 0)this.stencilRef=t.stencilRef;if(t.stencilFuncMask!==void 0)this.stencilFuncMask=t.stencilFuncMask;if(t.stencilFail!==void 0)this.stencilFail=t.stencilFail;if(t.stencilZFail!==void 0)this.stencilZFail=t.stencilZFail;if(t.stencilZPass!==void 0)this.stencilZPass=t.stencilZPass;if(t.stencilWrite!==void 0)this.stencilWrite=t.stencilWrite;if(t.wireframe!==void 0)this.wireframe=t.wireframe;if(t.wireframeLinewidth!==void 0)this.wireframeLinewidth=t.wireframeLinewidth;if(t.wireframeLinecap!==void 0)this.wireframeLinecap=t.wireframeLinecap;if(t.wireframeLinejoin!==void 0)this.wireframeLinejoin=t.wireframeLinejoin;if(t.rotation!==void 0)this.rotation=t.rotation;if(t.linewidth!==void 0)this.linewidth=t.linewidth;if(t.linecap!==void 0)this.linecap=t.linecap;if(t.linejoin!==void 0)this.linejoin=t.linejoin;if(t.dashSize!==void 0)this.dashSize=t.dashSize;if(t.gapSize!==void 0)this.gapSize=t.gapSize;if(t.scale!==void 0)this.scale=t.scale;if(t.polygonOffset!==void 0)this.polygonOffset=t.polygonOffset;if(t.polygonOffsetFactor!==void 0)this.polygonOffsetFactor=t.polygonOffsetFactor;if(t.polygonOffsetUnits!==void 0)this.polygonOffsetUnits=t.polygonOffsetUnits;if(t.dithering!==void 0)this.dithering=t.dithering;if(t.alphaToCoverage!==void 0)this.alphaToCoverage=t.alphaToCoverage;if(t.premultipliedAlpha!==void 0)this.premultipliedAlpha=t.premultipliedAlpha;if(t.forceSinglePass!==void 0)this.forceSinglePass=t.forceSinglePass;if(t.allowOverride!==void 0)this.allowOverride=t.allowOverride;if(t.visible!==void 0)this.visible=t.visible;if(t.toneMapped!==void 0)this.toneMapped=t.toneMapped;if(t.userData!==void 0)this.userData=t.userData;if(t.vertexColors!==void 0)if(typeof t.vertexColors==="number")this.vertexColors=t.vertexColors>0;else this.vertexColors=t.vertexColors;if(t.size!==void 0)this.size=t.size;if(t.sizeAttenuation!==void 0)this.sizeAttenuation=t.sizeAttenuation;if(t.map!==void 0)this.map=e[t.map]||null;if(t.matcap!==void 0)this.matcap=e[t.matcap]||null;if(t.alphaMap!==void 0)this.alphaMap=e[t.alphaMap]||null;if(t.bumpMap!==void 0)this.bumpMap=e[t.bumpMap]||null;if(t.bumpScale!==void 0)this.bumpScale=t.bumpScale;if(t.normalMap!==void 0)this.normalMap=e[t.normalMap]||null;if(t.normalMapType!==void 0)this.normalMapType=t.normalMapType;if(t.normalScale!==void 0){let n=t.normalScale;if(Array.isArray(n)===!1)n=[n,n];this.normalScale=new rt().fromArray(n)}if(t.displacementMap!==void 0)this.displacementMap=e[t.displacementMap]||null;if(t.displacementScale!==void 0)this.displacementScale=t.displacementScale;if(t.displacementBias!==void 0)this.displacementBias=t.displacementBias;if(t.roughnessMap!==void 0)this.roughnessMap=e[t.roughnessMap]||null;if(t.metalnessMap!==void 0)this.metalnessMap=e[t.metalnessMap]||null;if(t.emissiveMap!==void 0)this.emissiveMap=e[t.emissiveMap]||null;if(t.emissiveIntensity!==void 0)this.emissiveIntensity=t.emissiveIntensity;if(t.specularMap!==void 0)this.specularMap=e[t.specularMap]||null;if(t.specularIntensityMap!==void 0)this.specularIntensityMap=e[t.specularIntensityMap]||null;if(t.specularColorMap!==void 0)this.specularColorMap=e[t.specularColorMap]||null;if(t.envMap!==void 0)this.envMap=e[t.envMap]||null;if(t.envMapRotation!==void 0)this.envMapRotation.fromArray(t.envMapRotation);if(t.envMapIntensity!==void 0)this.envMapIntensity=t.envMapIntensity;if(t.reflectivity!==void 0)this.reflectivity=t.reflectivity;if(t.refractionRatio!==void 0)this.refractionRatio=t.refractionRatio;if(t.lightMap!==void 0)this.lightMap=e[t.lightMap]||null;if(t.lightMapIntensity!==void 0)this.lightMapIntensity=t.lightMapIntensity;if(t.aoMap!==void 0)this.aoMap=e[t.aoMap]||null;if(t.aoMapIntensity!==void 0)this.aoMapIntensity=t.aoMapIntensity;if(t.gradientMap!==void 0)this.gradientMap=e[t.gradientMap]||null;if(t.clearcoatMap!==void 0)this.clearcoatMap=e[t.clearcoatMap]||null;if(t.clearcoatRoughnessMap!==void 0)this.clearcoatRoughnessMap=e[t.clearcoatRoughnessMap]||null;if(t.clearcoatNormalMap!==void 0)this.clearcoatNormalMap=e[t.clearcoatNormalMap]||null;if(t.clearcoatNormalScale!==void 0)this.clearcoatNormalScale=new rt().fromArray(t.clearcoatNormalScale);if(t.iridescenceMap!==void 0)this.iridescenceMap=e[t.iridescenceMap]||null;if(t.iridescenceThicknessMap!==void 0)this.iridescenceThicknessMap=e[t.iridescenceThicknessMap]||null;if(t.transmissionMap!==void 0)this.transmissionMap=e[t.transmissionMap]||null;if(t.thicknessMap!==void 0)this.thicknessMap=e[t.thicknessMap]||null;if(t.anisotropyMap!==void 0)this.anisotropyMap=e[t.anisotropyMap]||null;if(t.sheenColorMap!==void 0)this.sheenColorMap=e[t.sheenColorMap]||null;if(t.sheenRoughnessMap!==void 0)this.sheenRoughnessMap=e[t.sheenRoughnessMap]||null;return this}clone(){return new this.constructor().copy(this)}copy(t){this.name=t.name,this.blending=t.blending,this.side=t.side,this.vertexColors=t.vertexColors,this.opacity=t.opacity,this.transparent=t.transparent,this.blendSrc=t.blendSrc,this.blendDst=t.blendDst,this.blendEquation=t.blendEquation,this.blendSrcAlpha=t.blendSrcAlpha,this.blendDstAlpha=t.blendDstAlpha,this.blendEquationAlpha=t.blendEquationAlpha,this.blendColor.copy(t.blendColor),this.blendAlpha=t.blendAlpha,this.depthFunc=t.depthFunc,this.depthTest=t.depthTest,this.depthWrite=t.depthWrite,this.stencilWriteMask=t.stencilWriteMask,this.stencilFunc=t.stencilFunc,this.stencilRef=t.stencilRef,this.stencilFuncMask=t.stencilFuncMask,this.stencilFail=t.stencilFail,this.stencilZFail=t.stencilZFail,this.stencilZPass=t.stencilZPass,this.stencilWrite=t.stencilWrite;let e=t.clippingPlanes,n=null;if(e!==null){let i=e.length;n=Array(i);for(let s=0;s!==i;++s)n[s]=e[s].clone()}return this.clippingPlanes=n,this.clipIntersection=t.clipIntersection,this.clipShadows=t.clipShadows,this.shadowSide=t.shadowSide,this.colorWrite=t.colorWrite,this.precision=t.precision,this.polygonOffset=t.polygonOffset,this.polygonOffsetFactor=t.polygonOffsetFactor,this.polygonOffsetUnits=t.polygonOffsetUnits,this.dithering=t.dithering,this.alphaTest=t.alphaTest,this.alphaHash=t.alphaHash,this.alphaToCoverage=t.alphaToCoverage,this.premultipliedAlpha=t.premultipliedAlpha,this.forceSinglePass=t.forceSinglePass,this.allowOverride=t.allowOverride,this.visible=t.visible,this.toneMapped=t.toneMapped,this.userData=JSON.parse(JSON.stringify(t.userData)),this}dispose(){this.dispatchEvent({type:"dispose"})}set needsUpdate(t){if(t===!0)this.version++}}class Ja extends Fe{constructor(t){super();this.isSpriteMaterial=!0,this.type="SpriteMaterial",this.color=new Pt(16777215),this.map=null,this.alphaMap=null,this.rotation=0,this.sizeAttenuation=!0,this.transparent=!0,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.alphaMap=t.alphaMap,this.rotation=t.rotation,this.sizeAttenuation=t.sizeAttenuation,this.fog=t.fog,this}}var rs,ks=new L,as=new L,os=new L,ls=new rt,Gs=new rt,gf=new Kt,Wr=new L,Hs=new L,Xr=new L,iu=new rt,nl=new rt,su=new rt;class Rc extends ae{constructor(t=new Ja){super();if(this.isSprite=!0,this.type="Sprite",rs===void 0){rs=new jt;let e=new Float32Array([-0.5,-0.5,0,0,0,0.5,-0.5,0,1,0,0.5,0.5,0,1,1,-0.5,0.5,0,0,1]),n=new ws(e,5);rs.setIndex([0,1,2,0,2,3]),rs.setAttribute("position",new ci(n,3,0,!1)),rs.setAttribute("uv",new ci(n,2,3,!1))}this.geometry=rs,this.material=t,this.center=new rt(0.5,0.5),this.count=1}intersectsFrustum(t){return t.intersectsSprite(this)}raycast(t,e){if(t.camera===null)Wt('Sprite: "Raycaster.camera" needs to be set in order to raycast against sprites.');if(as.setFromMatrixScale(this.matrixWorld),gf.copy(t.camera.matrixWorld),this.modelViewMatrix.multiplyMatrices(t.camera.matrixWorldInverse,this.matrixWorld),os.setFromMatrixPosition(this.modelViewMatrix),t.camera.isPerspectiveCamera&&this.material.sizeAttenuation===!1)as.multiplyScalar(-os.z);let n=this.material.rotation,i,s;if(n!==0)s=Math.cos(n),i=Math.sin(n);let r=this.center;qr(Wr.set(-0.5,-0.5,0),os,r,as,i,s),qr(Hs.set(0.5,-0.5,0),os,r,as,i,s),qr(Xr.set(0.5,0.5,0),os,r,as,i,s),iu.set(0,0),nl.set(1,0),su.set(1,1);let a=t.ray.intersectTriangle(Wr,Hs,Xr,!1,ks);if(a===null){if(qr(Hs.set(-0.5,0.5,0),os,r,as,i,s),nl.set(0,1),a=t.ray.intersectTriangle(Wr,Xr,Hs,!1,ks),a===null)return}let o=t.ray.origin.distanceTo(ks);if(o<t.near||o>t.far)return;e.push({distance:o,point:ks.clone(),uv:je.getInterpolation(ks,Wr,Hs,Xr,iu,nl,su,new rt),face:null,object:this})}copy(t,e){if(super.copy(t,e),t.center!==void 0)this.center.copy(t.center);return this.material=t.material,this}}function qr(t,e,n,i,s,r){if(ls.subVectors(t,n).addScalar(0.5).multiply(i),s!==void 0)Gs.x=r*ls.x-s*ls.y,Gs.y=s*ls.x+r*ls.y;else Gs.copy(ls);t.copy(e),t.x+=Gs.x,t.y+=Gs.y,t.applyMatrix4(gf)}var Yr=new L,ru=new L;class Cc extends ae{constructor(){super();this.isLOD=!0,this._currentLevel=0,this.type="LOD",Object.defineProperties(this,{levels:{enumerable:!0,value:[]}}),this.autoUpdate=!0}copy(t){super.copy(t,!1);let e=t.levels;for(let n=0,i=e.length;n<i;n++){let s=e[n];this.addLevel(s.object.clone(),s.distance,s.hysteresis)}return this.autoUpdate=t.autoUpdate,this}addLevel(t,e=0,n=0){e=Math.abs(e);let i=this.levels,s;for(s=0;s<i.length;s++)if(e<i[s].distance)break;return i.splice(s,0,{distance:e,hysteresis:n,object:t}),this.add(t),this}removeLevel(t){let e=this.levels;for(let n=0;n<e.length;n++)if(e[n].distance===t){let i=e.splice(n,1);return this.remove(i[0].object),!0}return!1}getCurrentLevel(){return this._currentLevel}getObjectForDistance(t){let e=this.levels;if(e.length>0){let n,i;for(n=1,i=e.length;n<i;n++){let s=e[n].distance;if(e[n].object.visible)s-=s*e[n].hysteresis;if(t<s)break}return e[n-1].object}return null}raycast(t,e){if(this.levels.length>0){Yr.setFromMatrixPosition(this.matrixWorld);let i=t.ray.origin.distanceTo(Yr);this.getObjectForDistance(i).raycast(t,e)}}update(t){let e=this.levels;if(e.length>1){Yr.setFromMatrixPosition(t.matrixWorld),ru.setFromMatrixPosition(this.matrixWorld);let n=Yr.distanceTo(ru)/t.zoom;e[0].object.visible=!0;let i,s;for(i=1,s=e.length;i<s;i++){let r=e[i].distance;if(e[i].object.visible)r-=r*e[i].hysteresis;if(n>=r)e[i-1].object.visible=!1,e[i].object.visible=!0;else break}this._currentLevel=i-1;for(;i<s;i++)e[i].object.visible=!1}}toJSON(t){let e=super.toJSON(t);e.object.autoUpdate=this.autoUpdate,e.object.levels=[];let n=this.levels;for(let i=0,s=n.length;i<s;i++){let r=n[i];e.object.levels.push({object:r.object.uuid,distance:r.distance,hysteresis:r.hysteresis})}return e}}var Vn=new L,il=new L,Zr=new L,Jr=new L;class Kn{constructor(t=new L,e=new L(0,0,-1)){this.origin=t,this.direction=e}set(t,e){return this.origin.copy(t),this.direction.copy(e),this}copy(t){return this.origin.copy(t.origin),this.direction.copy(t.direction),this}at(t,e){return e.copy(this.origin).addScaledVector(this.direction,t)}lookAt(t){return this.direction.copy(t).sub(this.origin).normalize(),this}recast(t){return this.origin.copy(this.at(t,Vn)),this}closestPointToPoint(t,e){e.subVectors(t,this.origin);let n=e.dot(this.direction);if(n<0)return e.copy(this.origin);return e.copy(this.origin).addScaledVector(this.direction,n)}distanceToPoint(t){return Math.sqrt(this.distanceSqToPoint(t))}distanceSqToPoint(t){let e=Vn.subVectors(t,this.origin).dot(this.direction);if(e<0)return this.origin.distanceToSquared(t);return Vn.copy(this.origin).addScaledVector(this.direction,e),Vn.distanceToSquared(t)}distanceSqToSegment(t,e,n,i){il.copy(t).add(e).multiplyScalar(0.5),Zr.copy(e).sub(t).normalize(),Jr.copy(this.origin).sub(il);let s=t.distanceTo(e)*0.5,r=-this.direction.dot(Zr),a=Jr.dot(this.direction),o=-Jr.dot(Zr),l=Jr.lengthSq(),c=Math.abs(1-r*r),h,d,u,f;if(c>0)if(h=r*o-a,d=r*a-o,f=s*c,h>=0)if(d>=-f)if(d<=f){let g=1/c;h*=g,d*=g,u=h*(h+r*d+2*a)+d*(r*h+d+2*o)+l}else d=s,h=Math.max(0,-(r*d+a)),u=-h*h+d*(d+2*o)+l;else d=-s,h=Math.max(0,-(r*d+a)),u=-h*h+d*(d+2*o)+l;else if(d<=-f)h=Math.max(0,-(-r*s+a)),d=h>0?-s:Math.min(Math.max(-s,-o),s),u=-h*h+d*(d+2*o)+l;else if(d<=f)h=0,d=Math.min(Math.max(-s,-o),s),u=d*(d+2*o)+l;else h=Math.max(0,-(r*s+a)),d=h>0?s:Math.min(Math.max(-s,-o),s),u=-h*h+d*(d+2*o)+l;else d=r>0?-s:s,h=Math.max(0,-(r*d+a)),u=-h*h+d*(d+2*o)+l;if(n)n.copy(this.origin).addScaledVector(this.direction,h);if(i)i.copy(il).addScaledVector(Zr,d);return u}intersectSphere(t,e){if(t.radius<0)return null;Vn.subVectors(t.center,this.origin);let n=Vn.dot(this.direction),i=Vn.dot(Vn)-n*n,s=t.radius*t.radius;if(i>s)return null;let r=Math.sqrt(s-i),a=n-r,o=n+r;if(o<0)return null;if(a<0)return this.at(o,e);return this.at(a,e)}intersectsSphere(t){if(t.radius<0)return!1;return this.distanceSqToPoint(t.center)<=t.radius*t.radius}distanceToPlane(t){let e=t.normal.dot(this.direction);if(e===0){if(t.distanceToPoint(this.origin)===0)return 0;return null}let n=-(this.origin.dot(t.normal)+t.constant)/e;return n>=0?n:null}intersectPlane(t,e){let n=this.distanceToPlane(t);if(n===null)return null;return this.at(n,e)}intersectsPlane(t){let e=t.distanceToPoint(this.origin);if(e===0)return!0;if(t.normal.dot(this.direction)*e<0)return!0;return!1}intersectBox(t,e){let n,i,s,r,a,o,l=1/this.direction.x,c=1/this.direction.y,h=1/this.direction.z,d=this.origin;if(l>=0)n=(t.min.x-d.x)*l,i=(t.max.x-d.x)*l;else n=(t.max.x-d.x)*l,i=(t.min.x-d.x)*l;if(c>=0)s=(t.min.y-d.y)*c,r=(t.max.y-d.y)*c;else s=(t.max.y-d.y)*c,r=(t.min.y-d.y)*c;if(n>r||s>i)return null;if(s>n||isNaN(n))n=s;if(r<i||isNaN(i))i=r;if(h>=0)a=(t.min.z-d.z)*h,o=(t.max.z-d.z)*h;else a=(t.max.z-d.z)*h,o=(t.min.z-d.z)*h;if(n>o||a>i)return null;if(a>n||n!==n)n=a;if(o<i||i!==i)i=o;if(i<0)return null;return this.at(n>=0?n:i,e)}intersectsBox(t){return this.intersectBox(t,Vn)!==null}intersectTriangle(t,e,n,i,s){let r=this.origin,a=this.direction,{x:o,y:l,z:c}=a,h=t.x-r.x,d=t.y-r.y,u=t.z-r.z,f=e.x-r.x,g=e.y-r.y,_=e.z-r.z,m=n.x-r.x,p=n.y-r.y,b=n.z-r.z,S=Math.abs(o),x=Math.abs(l),T=Math.abs(c),I,R,y,A,Z,D,q,ot,z,et,nt,B;if(S>=x&&S>=T)if(y=o,D=h,z=f,B=m,o>=0)I=l,R=c,A=d,Z=u,q=g,ot=_,et=p,nt=b;else I=c,R=l,A=u,Z=d,q=_,ot=g,et=b,nt=p;else if(x>=T)if(y=l,D=d,z=g,B=p,l>=0)I=c,R=o,A=u,Z=h,q=_,ot=f,et=b,nt=m;else I=o,R=c,A=h,Z=u,q=f,ot=_,et=m,nt=b;else if(y=c,D=u,z=_,B=b,c>=0)I=o,R=l,A=h,Z=d,q=f,ot=g,et=m,nt=p;else I=l,R=o,A=d,Z=h,q=g,ot=f,et=p,nt=m;if(y===0)return null;let ft=I/y,Q=R/y,ct=1/y,dt=A-ft*D,Bt=Z-Q*D,Dt=q-ft*z,ee=ot-Q*z,Jt=et-ft*B,F=nt-Q*B,O=Jt*ee-F*Dt,N=dt*F-Bt*Jt,P=Dt*Bt-ee*dt;if(i){if(O<0||N<0||P<0)return null}else if((O<0||N<0||P<0)&&(O>0||N>0||P>0))return null;let G=O+N+P;if(G===0)return null;let J=ct*(O*D+N*z+P*B);if(G>0?J<0:J>0)return null;return this.at(J/G,s)}applyMatrix4(t){return this.origin.applyMatrix4(t),this.direction.transformDirection(t),this}equals(t){return t.origin.equals(this.origin)&&t.direction.equals(this.direction)}clone(){return new this.constructor().copy(this)}}class mn extends Fe{constructor(t){super();this.isMeshBasicMaterial=!0,this.type="MeshBasicMaterial",this.color=new Pt(16777215),this.map=null,this.lightMap=null,this.lightMapIntensity=1,this.aoMap=null,this.aoMapIntensity=1,this.specularMap=null,this.alphaMap=null,this.envMap=null,this.envMapRotation=new Mn,this.combine=0,this.reflectivity=1,this.refractionRatio=0.98,this.wireframe=!1,this.wireframeLinewidth=1,this.wireframeLinecap="round",this.wireframeLinejoin="round",this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.lightMap=t.lightMap,this.lightMapIntensity=t.lightMapIntensity,this.aoMap=t.aoMap,this.aoMapIntensity=t.aoMapIntensity,this.specularMap=t.specularMap,this.alphaMap=t.alphaMap,this.envMap=t.envMap,this.envMapRotation.copy(t.envMapRotation),this.combine=t.combine,this.reflectivity=t.reflectivity,this.refractionRatio=t.refractionRatio,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.wireframeLinecap=t.wireframeLinecap,this.wireframeLinejoin=t.wireframeLinejoin,this.fog=t.fog,this}}var au=new Kt,Mi=new Kn,$r=new Ue,ou=new L,Kr=new L,jr=new L,Qr=new L,sl=new L,ta=new L,lu=new L,ea=new L;class pe extends ae{constructor(t=new jt,e=new mn){super();this.isMesh=!0,this.type="Mesh",this.geometry=t,this.material=e,this.morphTargetDictionary=void 0,this.morphTargetInfluences=void 0,this.count=1,this.updateMorphTargets()}copy(t,e){if(super.copy(t,e),t.morphTargetInfluences!==void 0)this.morphTargetInfluences=t.morphTargetInfluences.slice();if(t.morphTargetDictionary!==void 0)this.morphTargetDictionary=Object.assign({},t.morphTargetDictionary);return this.material=Array.isArray(t.material)?t.material.slice():t.material,this.geometry=t.geometry,this}updateMorphTargets(){let e=this.geometry.morphAttributes,n=Object.keys(e);if(n.length>0){let i=e[n[0]];if(i!==void 0){this.morphTargetInfluences=[],this.morphTargetDictionary={};for(let s=0,r=i.length;s<r;s++){let a=i[s].name||String(s);this.morphTargetInfluences.push(0),this.morphTargetDictionary[a]=s}}}}getVertexPosition(t,e){let n=this.geometry,i=n.attributes.position,s=n.morphAttributes.position,r=n.morphTargetsRelative;e.fromBufferAttribute(i,t);let a=this.morphTargetInfluences;if(s&&a){ta.set(0,0,0);for(let o=0,l=s.length;o<l;o++){let c=a[o],h=s[o];if(c===0)continue;if(sl.fromBufferAttribute(h,t),r)ta.addScaledVector(sl,c);else ta.addScaledVector(sl.sub(e),c)}e.add(ta)}return e}intersectsFrustum(t){return t.intersectsObject(this)}raycast(t,e){let n=this.geometry,i=this.material,s=this.matrixWorld;if(i===void 0)return;if(n.boundingSphere===null)n.computeBoundingSphere();if($r.copy(n.boundingSphere),$r.applyMatrix4(s),Mi.copy(t.ray).recast(t.near),$r.containsPoint(Mi.origin)===!1){if(Mi.intersectSphere($r,ou)===null)return;if(Mi.origin.distanceToSquared(ou)>(t.far-t.near)**2)return}if(au.copy(s).invert(),Mi.copy(t.ray).applyMatrix4(au),n.boundingBox!==null){if(Mi.intersectsBox(n.boundingBox)===!1)return}this._computeIntersections(t,e,Mi)}_computeIntersections(t,e,n){let i,s=this.geometry,r=this.material,a=s.index,o=s.attributes.position,l=s.attributes.uv,c=s.attributes.uv1,h=s.attributes.normal,{groups:d,drawRange:u}=s;if(a!==null)if(Array.isArray(r))for(let f=0,g=d.length;f<g;f++){let _=d[f],m=r[_.materialIndex],p=Math.max(_.start,u.start),b=Math.min(a.count,Math.min(_.start+_.count,u.start+u.count));for(let S=p,x=b;S<x;S+=3){let T=a.getX(S),I=a.getX(S+1),R=a.getX(S+2);if(i=na(this,m,t,n,l,c,h,T,I,R),i)i.faceIndex=Math.floor(S/3),i.face.materialIndex=_.materialIndex,e.push(i)}}else{let f=Math.max(0,u.start),g=Math.min(a.count,u.start+u.count);for(let _=f,m=g;_<m;_+=3){let p=a.getX(_),b=a.getX(_+1),S=a.getX(_+2);if(i=na(this,r,t,n,l,c,h,p,b,S),i)i.faceIndex=Math.floor(_/3),e.push(i)}}else if(o!==void 0)if(Array.isArray(r))for(let f=0,g=d.length;f<g;f++){let _=d[f],m=r[_.materialIndex],p=Math.max(_.start,u.start),b=Math.min(o.count,Math.min(_.start+_.count,u.start+u.count));for(let S=p,x=b;S<x;S+=3){let T=S,I=S+1,R=S+2;if(i=na(this,m,t,n,l,c,h,T,I,R),i)i.faceIndex=Math.floor(S/3),i.face.materialIndex=_.materialIndex,e.push(i)}}else{let f=Math.max(0,u.start),g=Math.min(o.count,u.start+u.count);for(let _=f,m=g;_<m;_+=3){let p=_,b=_+1,S=_+2;if(i=na(this,r,t,n,l,c,h,p,b,S),i)i.faceIndex=Math.floor(_/3),e.push(i)}}}}function l0(t,e,n,i,s,r,a,o){let l;if(e.side===1)l=i.intersectTriangle(a,r,s,!0,o);else l=i.intersectTriangle(s,r,a,e.side===0,o);if(l===null)return null;ea.copy(o),ea.applyMatrix4(t.matrixWorld);let c=n.ray.origin.distanceTo(ea);if(c<n.near||c>n.far)return null;return{distance:c,point:ea.clone(),object:t}}function na(t,e,n,i,s,r,a,o,l,c){t.getVertexPosition(o,Kr),t.getVertexPosition(l,jr),t.getVertexPosition(c,Qr);let h=l0(t,e,n,i,Kr,jr,Qr,lu);if(h){let d=new L;if(je.getBarycoord(lu,Kr,jr,Qr,d),s)h.uv=je.getInterpolatedAttribute(s,o,l,c,d,new rt);if(r)h.uv1=je.getInterpolatedAttribute(r,o,l,c,d,new rt);if(a){if(h.normal=je.getInterpolatedAttribute(a,o,l,c,d,new L),h.normal.dot(i.direction)>0)h.normal.multiplyScalar(-1)}let u={a:o,b:l,c,normal:new L,materialIndex:0};je.getNormal(Kr,jr,Qr,u.normal),h.face=u,h.barycoord=d}return h}var Vs=new fe,cu=new fe,hu=new fe,c0=new fe,uu=new Kt,ia=new L,rl=new Ue,du=new Kt,al=new Kn;class Ic extends pe{constructor(t,e){super(t,e);this.isSkinnedMesh=!0,this.type="SkinnedMesh",this.bindMode="attached",this.bindMatrix=new Kt,this.bindMatrixInverse=new Kt,this.boundingBox=null,this.boundingSphere=null}computeBoundingBox(){let t=this.geometry;if(this.boundingBox===null)this.boundingBox=new Oe;this.boundingBox.makeEmpty();let e=t.getAttribute("position");for(let n=0;n<e.count;n++)this.getVertexPosition(n,ia),this.boundingBox.expandByPoint(ia)}computeBoundingSphere(){let t=this.geometry;if(this.boundingSphere===null)this.boundingSphere=new Ue;this.boundingSphere.makeEmpty();let e=t.getAttribute("position");for(let n=0;n<e.count;n++)this.getVertexPosition(n,ia),this.boundingSphere.expandByPoint(ia)}copy(t,e){if(super.copy(t,e),this.bindMode=t.bindMode,this.bindMatrix.copy(t.bindMatrix),this.bindMatrixInverse.copy(t.bindMatrixInverse),this.skeleton=t.skeleton,t.boundingBox!==null)this.boundingBox=t.boundingBox.clone();if(t.boundingSphere!==null)this.boundingSphere=t.boundingSphere.clone();return this}raycast(t,e){let n=this.material,i=this.matrixWorld;if(n===void 0)return;if(this.boundingSphere===null)this.computeBoundingSphere();if(rl.copy(this.boundingSphere),rl.applyMatrix4(i),t.ray.intersectsSphere(rl)===!1)return;if(du.copy(i).invert(),al.copy(t.ray).applyMatrix4(du),this.boundingBox!==null){if(al.intersectsBox(this.boundingBox)===!1)return}this._computeIntersections(t,e,al)}getVertexPosition(t,e){return super.getVertexPosition(t,e),this.applyBoneTransform(t,e),e}bind(t,e){if(this.skeleton=t,e===void 0)this.updateMatrixWorld(!0),this.skeleton.calculateInverses(),e=this.matrixWorld;this.bindMatrix.copy(e),this.bindMatrixInverse.copy(e).invert()}pose(){this.skeleton.pose()}normalizeSkinWeights(){let t=new fe,e=this.geometry.attributes.skinWeight;for(let n=0,i=e.count;n<i;n++){t.fromBufferAttribute(e,n);let s=1/t.manhattanLength();if(s!==1/0)t.multiplyScalar(s);else t.set(1,0,0,0);e.setXYZW(n,t.x,t.y,t.z,t.w)}}updateMatrixWorld(t){if(super.updateMatrixWorld(t),this.bindMode==="attached")this.bindMatrixInverse.copy(this.matrixWorld).invert();else if(this.bindMode==="detached")this.bindMatrixInverse.copy(this.bindMatrix).invert();else Mt("SkinnedMesh: Unrecognized bindMode: "+this.bindMode)}applyBoneTransform(t,e){let n=this.skeleton,i=this.geometry;if(cu.fromBufferAttribute(i.attributes.skinIndex,t),hu.fromBufferAttribute(i.attributes.skinWeight,t),e.isVector4)Vs.copy(e),e.set(0,0,0,0);else Vs.set(...e,1),e.set(0,0,0);Vs.applyMatrix4(this.bindMatrix);for(let s=0;s<4;s++){let r=hu.getComponent(s);if(r!==0){let a=cu.getComponent(s);uu.multiplyMatrices(n.bones[a].matrixWorld,n.boneInverses[a]),e.addScaledVector(c0.copy(Vs).applyMatrix4(uu),r)}}if(e.isVector4)e.w=Vs.w;return e.applyMatrix4(this.bindMatrixInverse)}}class $a extends ae{constructor(){super();this.isBone=!0,this.type="Bone"}}class an extends be{constructor(t=null,e=1,n=1,i,s,r,a,o,l=1003,c=1003,h,d){super(null,r,a,o,l,c,i,s,h,d);this.isDataTexture=!0,this.image={data:t,width:e,height:n},this.generateMipmaps=!1,this.flipY=!1,this.unpackAlignment=1}}var fu=new Kt,h0=new Kt;class Ka{constructor(t=[],e=[]){this.uuid=rn(),this.bones=t.slice(0),this.boneInverses=e,this.boneMatrices=null,this.boneTexture=null,this.init()}init(){let t=this.bones,e=this.boneInverses;if(this.boneMatrices=new Float32Array(t.length*16),e.length===0)this.calculateInverses();else if(t.length!==e.length){Mt("Skeleton: Number of inverse bone matrices does not match amount of bones."),this.boneInverses=[];for(let n=0,i=this.bones.length;n<i;n++)this.boneInverses.push(new Kt)}}calculateInverses(){this.boneInverses.length=0;for(let t=0,e=this.bones.length;t<e;t++){let n=new Kt;if(this.bones[t])n.copy(this.bones[t].matrixWorld).invert();this.boneInverses.push(n)}}pose(){for(let t=0,e=this.bones.length;t<e;t++){let n=this.bones[t];if(n)n.matrixWorld.copy(this.boneInverses[t]).invert()}for(let t=0,e=this.bones.length;t<e;t++){let n=this.bones[t];if(n){if(n.parent&&n.parent.isBone)n.matrix.copy(n.parent.matrixWorld).invert(),n.matrix.multiply(n.matrixWorld);else n.matrix.copy(n.matrixWorld);n.matrix.decompose(n.position,n.quaternion,n.scale)}}}update(){let t=this.bones,e=this.boneInverses,n=this.boneMatrices,i=this.boneTexture;for(let s=0,r=t.length;s<r;s++){let a=t[s]?t[s].matrixWorld:h0;fu.multiplyMatrices(a,e[s]),fu.toArray(n,s*16)}if(i!==null)i.needsUpdate=!0}clone(){return new Ka(this.bones,this.boneInverses)}computeBoneTexture(){let t=Math.sqrt(this.bones.length*4);t=Math.ceil(t/4)*4,t=Math.max(t,4);let e=new Float32Array(t*t*4);e.set(this.boneMatrices);let n=new an(e,t,t,1023,1015);return n.needsUpdate=!0,this.boneMatrices=e,this.boneTexture=n,this}getBoneByName(t){for(let e=0,n=this.bones.length;e<n;e++){let i=this.bones[e];if(i.name===t)return i}return}dispose(){if(this.boneTexture!==null)this.boneTexture.dispose(),this.boneTexture=null}fromJSON(t,e){this.uuid=t.uuid;for(let n=0,i=t.bones.length;n<i;n++){let s=t.bones[n],r=e[s];if(r===void 0)Mt("Skeleton: No bone found with UUID:",s),r=new $a;this.bones.push(r),this.boneInverses.push(new Kt().fromArray(t.boneInverses[n]))}return this.init(),this}toJSON(){let t={metadata:{version:4.7,type:"Skeleton",generator:"Skeleton.toJSON"},bones:[],boneInverses:[]};t.uuid=this.uuid;let e=this.bones,n=this.boneInverses;for(let i=0,s=e.length;i<s;i++){let r=e[i];t.bones.push(r.uuid);let a=n[i];t.boneInverses.push(a.toArray())}return t}}class hi extends de{constructor(t,e,n,i=1){super(t,e,n);this.isInstancedBufferAttribute=!0,this.meshPerAttribute=i}copy(t){return super.copy(t),this.meshPerAttribute=t.meshPerAttribute,this}toJSON(){let t=super.toJSON();return t.meshPerAttribute=this.meshPerAttribute,t.isInstancedBufferAttribute=!0,t}}var cs=new Kt,pu=new Kt,sa=[],mu=new Oe,u0=new Kt,Ws=new pe,Xs=new Ue;class pr extends pe{constructor(t,e,n){super(t,e);this.isInstancedMesh=!0,this.instanceMatrix=new hi(new Float32Array(n*16),16),this.instanceColor=null,this.morphTexture=null,this.count=n,this.boundingBox=null,this.boundingSphere=null;for(let i=0;i<n;i++)this.setMatrixAt(i,u0)}computeBoundingBox(){let t=this.geometry,e=this.count;if(this.boundingBox===null)this.boundingBox=new Oe;if(t.boundingBox===null)t.computeBoundingBox();this.boundingBox.makeEmpty();for(let n=0;n<e;n++)this.getMatrixAt(n,cs),mu.copy(t.boundingBox).applyMatrix4(cs),this.boundingBox.union(mu)}computeBoundingSphere(){let t=this.geometry,e=this.count;if(this.boundingSphere===null)this.boundingSphere=new Ue;if(t.boundingSphere===null)t.computeBoundingSphere();this.boundingSphere.makeEmpty();for(let n=0;n<e;n++)this.getMatrixAt(n,cs),Xs.copy(t.boundingSphere).applyMatrix4(cs),this.boundingSphere.union(Xs)}copy(t,e){if(super.copy(t,e),this.instanceMatrix.copy(t.instanceMatrix),t.morphTexture!==null)this.morphTexture=t.morphTexture.clone();if(t.instanceColor!==null)this.instanceColor=t.instanceColor.clone();if(this.count=t.count,t.boundingBox!==null)this.boundingBox=t.boundingBox.clone();if(t.boundingSphere!==null)this.boundingSphere=t.boundingSphere.clone();return this}getColorAt(t,e){if(this.instanceColor===null)return e.setRGB(1,1,1);else return e.fromArray(this.instanceColor.array,t*3)}getMatrixAt(t,e){return e.fromArray(this.instanceMatrix.array,t*16)}getMorphAt(t,e){let n=e.morphTargetInfluences,i=this.morphTexture.source.data.data,s=n.length+1,r=t*s+1;for(let a=0;a<n.length;a++)n[a]=i[r+a]}raycast(t,e){let n=this.matrixWorld,i=this.count;if(Ws.geometry=this.geometry,Ws.material=this.material,Ws.material===void 0)return;if(this.boundingSphere===null)this.computeBoundingSphere();if(Xs.copy(this.boundingSphere),Xs.applyMatrix4(n),t.ray.intersectsSphere(Xs)===!1)return;for(let s=0;s<i;s++){this.getMatrixAt(s,cs),pu.multiplyMatrices(n,cs),Ws.matrixWorld=pu,Ws.raycast(t,sa);for(let r=0,a=sa.length;r<a;r++){let o=sa[r];o.instanceId=s,o.object=this,e.push(o)}sa.length=0}}setColorAt(t,e){if(this.instanceColor===null)this.instanceColor=new hi(new Float32Array(this.instanceMatrix.count*3).fill(1),3);return e.toArray(this.instanceColor.array,t*3),this}setMatrixAt(t,e){return e.toArray(this.instanceMatrix.array,t*16),this}setMorphAt(t,e){let n=e.morphTargetInfluences,i=n.length+1;if(this.morphTexture===null)this.morphTexture=new an(new Float32Array(i*this.count),i,this.count,1028,1015);let s=this.morphTexture.source.data.data,r=0;for(let l=0;l<n.length;l++)r+=n[l];let a=this.geometry.morphTargetsRelative?1:1-r,o=i*t;return s[o]=a,s.set(n,o+1),this}updateMorphTargets(){}dispose(){if(super.dispose(),this.morphTexture!==null)this.morphTexture.dispose(),this.morphTexture=null}}var bi=new Ue,d0=new rt(0.5,0.5),ra=new L;class ui{constructor(t=new dn,e=new dn,n=new dn,i=new dn,s=new dn,r=new dn){this.planes=[t,e,n,i,s,r]}set(t,e,n,i,s,r){let a=this.planes;return a[0].copy(t),a[1].copy(e),a[2].copy(n),a[3].copy(i),a[4].copy(s),a[5].copy(r),this}copy(t){let e=this.planes;for(let n=0;n<6;n++)e[n].copy(t.planes[n]);return this}setFromProjectionMatrix(t,e=2000,n=!1){let i=this.planes,s=t.elements,r=s[0],a=s[1],o=s[2],l=s[3],c=s[4],h=s[5],d=s[6],u=s[7],f=s[8],g=s[9],_=s[10],m=s[11],p=s[12],b=s[13],S=s[14],x=s[15];if(i[0].setComponents(l-r,u-c,m-f,x-p).normalize(),i[1].setComponents(l+r,u+c,m+f,x+p).normalize(),i[2].setComponents(l+a,u+h,m+g,x+b).normalize(),i[3].setComponents(l-a,u-h,m-g,x-b).normalize(),n)i[4].setComponents(o,d,_,S).normalize(),i[5].setComponents(l-o,u-d,m-_,x-S).normalize();else if(i[4].setComponents(l-o,u-d,m-_,x-S).normalize(),e===2000)i[5].setComponents(l+o,u+d,m+_,x+S).normalize();else if(e===2001)i[5].setComponents(o,d,_,S).normalize();else throw Error("THREE.Frustum.setFromProjectionMatrix(): Invalid coordinate system: "+e);return this}intersectsObject(t){if(t.boundingSphere!==void 0){if(t.boundingSphere===null)t.computeBoundingSphere();bi.copy(t.boundingSphere).applyMatrix4(t.matrixWorld)}else{let e=t.geometry;if(e.boundingSphere===null)e.computeBoundingSphere();bi.copy(e.boundingSphere).applyMatrix4(t.matrixWorld)}return this.intersectsSphere(bi)}intersectsSprite(t){bi.center.set(0,0,0);let e=d0.distanceTo(t.center);return bi.radius=0.7071067811865476+e,bi.applyMatrix4(t.matrixWorld),this.intersectsSphere(bi)}intersectsSphere(t){let e=this.planes,n=t.center,i=-t.radius;for(let s=0;s<6;s++)if(e[s].distanceToPoint(n)<i)return!1;return!0}intersectsBox(t){let e=this.planes;for(let n=0;n<6;n++){let i=e[n];if(ra.x=i.normal.x>0?t.max.x:t.min.x,ra.y=i.normal.y>0?t.max.y:t.min.y,ra.z=i.normal.z>0?t.max.z:t.min.z,i.distanceToPoint(ra)<0)return!1}return!0}containsPoint(t){let e=this.planes;for(let n=0;n<6;n++)if(e[n].distanceToPoint(t)<0)return!1;return!0}clone(){return new this.constructor().copy(this)}}var gu=new Kt;class ja{constructor(){this.coordinateSystem=2000,this._frustums=[],this._count=0}setFromArrayCamera(t){let e=t.cameras,n=this._frustums;for(let i=0;i<e.length;i++){let s=e[i];if(gu.multiplyMatrices(s.projectionMatrix,s.matrixWorldInverse),n[i]===void 0)n[i]=new ui;n[i].setFromProjectionMatrix(gu,s.coordinateSystem,s.reversedDepth)}return this._count=e.length,this}intersectsObject(t){let e=this._frustums;for(let n=0;n<this._count;n++)if(e[n].intersectsObject(t))return!0;return!1}intersectsSprite(t){let e=this._frustums;for(let n=0;n<this._count;n++)if(e[n].intersectsSprite(t))return!0;return!1}intersectsSphere(t){let e=this._frustums;for(let n=0;n<this._count;n++)if(e[n].intersectsSphere(t))return!0;return!1}intersectsBox(t){let e=this._frustums;for(let n=0;n<this._count;n++)if(e[n].intersectsBox(t))return!0;return!1}containsPoint(t){let e=this._frustums;for(let n=0;n<this._count;n++)if(e[n].containsPoint(t))return!0;return!1}copy(t){this.coordinateSystem=t.coordinateSystem;let e=this._frustums,n=t._frustums;for(let i=0;i<t._count;i++){if(e[i]===void 0)e[i]=new ui;e[i].copy(n[i])}return this._count=t._count,this}clone(){return new ja().copy(this)}}function ol(t,e){return t-e}function f0(t,e){return t.z-e.z}function p0(t,e){return e.z-t.z}class _f{constructor(){this.index=0,this.pool=[],this.list=[]}push(t,e,n,i){let s=this.pool,r=this.list;if(this.index>=s.length)s.push({start:-1,count:-1,z:-1,index:-1});let a=s[this.index];r.push(a),this.index++,a.start=t,a.count=e,a.z=n,a.index=i}reset(){this.list.length=0,this.index=0}}var $e=new Kt,m0=new Pt(1,1,1),g0=new ui,_0=new ja,aa=new Oe,Ti=new Ue,qs=new L,_u=new L,x0=new L,ll=new _f,He=new pe,oa=[];function y0(t,e,n=0){let i=e.itemSize;if(t.isInterleavedBufferAttribute||t.array.constructor!==e.array.constructor){let s=t.count;for(let r=0;r<s;r++)for(let a=0;a<i;a++)e.setComponent(r+n,a,t.getComponent(r,a))}else e.array.set(t.array,n*i);e.needsUpdate=!0}function Ei(t,e){if(t.constructor!==e.constructor){let n=Math.min(t.length,e.length);for(let i=0;i<n;i++)e[i]=t[i]}else{let n=Math.min(t.length,e.length);e.set(new t.constructor(t.buffer,0,n))}}class Pc extends pe{constructor(t,e,n=e*2,i){super(new jt,i);this.isBatchedMesh=!0,this.perObjectFrustumCulled=!0,this.sortObjects=!0,this.boundingBox=null,this.boundingSphere=null,this.customSort=null,this._instanceInfo=[],this._geometryInfo=[],this._availableInstanceIds=[],this._availableGeometryIds=[],this._nextIndexStart=0,this._nextVertexStart=0,this._geometryCount=0,this._visibilityChanged=!0,this._geometryInitialized=!1,this._maxInstanceCount=t,this._maxVertexCount=e,this._maxIndexCount=n,this._multiDrawCounts=new Int32Array(t),this._multiDrawStarts=new Int32Array(t),this._multiDrawCount=0,this._multiDrawBytesPerElement=1,this._matricesTexture=null,this._indirectTexture=null,this._colorsTexture=null,this._initMatricesTexture(),this._initIndirectTexture()}get maxInstanceCount(){return this._maxInstanceCount}get instanceCount(){return this._instanceInfo.length-this._availableInstanceIds.length}get unusedVertexCount(){return this._maxVertexCount-this._nextVertexStart}get unusedIndexCount(){return this._maxIndexCount-this._nextIndexStart}_initMatricesTexture(){let t=Math.sqrt(this._maxInstanceCount*4);t=Math.ceil(t/4)*4,t=Math.max(t,4);let e=new Float32Array(t*t*4),n=new an(e,t,t,1023,1015);this._matricesTexture=n}_initIndirectTexture(){let t=Math.sqrt(this._maxInstanceCount);t=Math.ceil(t);let e=new Uint32Array(t*t),n=new an(e,t,t,1029,1014);this._indirectTexture=n}_initColorsTexture(){let t=Math.sqrt(this._maxInstanceCount);t=Math.ceil(t);let e=new Float32Array(t*t*4).fill(1),n=new an(e,t,t,1023,1015);n.colorSpace=re.workingColorSpace,this._colorsTexture=n}_initializeGeometry(t){let e=this.geometry,n=this._maxVertexCount,i=this._maxIndexCount;if(this._geometryInitialized===!1){for(let s in t.attributes){let r=t.getAttribute(s),{array:a,itemSize:o,normalized:l}=r,c=new a.constructor(n*o),h=new de(c,o,l);e.setAttribute(s,h)}if(t.getIndex()!==null){let s=n>65535?new Uint32Array(i):new Uint16Array(i);e.setIndex(new de(s,1))}this._geometryInitialized=!0}}_validateGeometry(t){let e=this.geometry;if(Boolean(t.getIndex())!==Boolean(e.getIndex()))throw Error('THREE.BatchedMesh: All geometries must consistently have "index".');for(let n in e.attributes){if(!t.hasAttribute(n))throw Error(`THREE.BatchedMesh: Added geometry missing "${n}". All geometries must have consistent attributes.`);let i=t.getAttribute(n),s=e.getAttribute(n);if(i.itemSize!==s.itemSize||i.normalized!==s.normalized)throw Error("THREE.BatchedMesh: All attributes must have a consistent itemSize and normalized value.")}}validateInstanceId(t){let e=this._instanceInfo;if(t<0||t>=e.length||e[t].active===!1)throw Error(`THREE.BatchedMesh: Invalid instanceId ${t}. Instance is either out of range or has been deleted.`)}validateGeometryId(t){let e=this._geometryInfo;if(t<0||t>=e.length||e[t].active===!1)throw Error(`THREE.BatchedMesh: Invalid geometryId ${t}. Geometry is either out of range or has been deleted.`)}setCustomSort(t){return this.customSort=t,this}computeBoundingBox(){if(this.boundingBox===null)this.boundingBox=new Oe;let t=this.boundingBox,e=this._instanceInfo;t.makeEmpty();for(let n=0,i=e.length;n<i;n++){if(e[n].active===!1)continue;let s=e[n].geometryIndex;this.getMatrixAt(n,$e),this.getBoundingBoxAt(s,aa).applyMatrix4($e),t.union(aa)}}computeBoundingSphere(){if(this.boundingSphere===null)this.boundingSphere=new Ue;let t=this.boundingSphere,e=this._instanceInfo;t.makeEmpty();for(let n=0,i=e.length;n<i;n++){if(e[n].active===!1)continue;let s=e[n].geometryIndex;this.getMatrixAt(n,$e),this.getBoundingSphereAt(s,Ti).applyMatrix4($e),t.union(Ti)}}addInstance(t){if(this._instanceInfo.length>=this.maxInstanceCount&&this._availableInstanceIds.length===0)throw Error("THREE.BatchedMesh: Maximum item count reached.");let n={visible:!0,active:!0,geometryIndex:t},i=null;if(this._availableInstanceIds.length>0)this._availableInstanceIds.sort(ol),i=this._availableInstanceIds.shift(),this._instanceInfo[i]=n;else i=this._instanceInfo.length,this._instanceInfo.push(n);let s=this._matricesTexture;$e.identity().toArray(s.image.data,i*16),s.needsUpdate=!0;let r=this._colorsTexture;if(r)m0.toArray(r.image.data,i*4),r.needsUpdate=!0;return this._visibilityChanged=!0,i}addGeometry(t,e=-1,n=-1){this._initializeGeometry(t),this._validateGeometry(t);let i={vertexStart:-1,vertexCount:-1,reservedVertexCount:-1,indexStart:-1,indexCount:-1,reservedIndexCount:-1,start:-1,count:-1,boundingBox:null,boundingSphere:null,active:!0},s=this._geometryInfo;i.vertexStart=this._nextVertexStart,i.reservedVertexCount=e===-1?t.getAttribute("position").count:e;let r=t.getIndex();if(r!==null)i.indexStart=this._nextIndexStart,i.reservedIndexCount=n===-1?r.count:n;if(i.indexStart!==-1&&i.indexStart+i.reservedIndexCount>this._maxIndexCount||i.vertexStart+i.reservedVertexCount>this._maxVertexCount)throw Error("THREE.BatchedMesh: Reserved space request exceeds the maximum buffer size.");let o;if(this._availableGeometryIds.length>0)this._availableGeometryIds.sort(ol),o=this._availableGeometryIds.shift(),s[o]=i;else o=this._geometryCount,this._geometryCount++,s.push(i);return this.setGeometryAt(o,t),this._nextIndexStart=i.indexStart+i.reservedIndexCount,this._nextVertexStart=i.vertexStart+i.reservedVertexCount,o}setGeometryAt(t,e){if(t>=this._geometryCount)throw Error("THREE.BatchedMesh: Maximum geometry count reached.");this._validateGeometry(e);let n=this.geometry,i=n.getIndex()!==null,s=n.getIndex(),r=e.getIndex(),a=this._geometryInfo[t];if(i&&r.count>a.reservedIndexCount||e.attributes.position.count>a.reservedVertexCount)throw Error("THREE.BatchedMesh: Reserved space not large enough for provided geometry.");let{vertexStart:o,reservedVertexCount:l}=a;a.vertexCount=e.getAttribute("position").count;for(let c in n.attributes){let h=e.getAttribute(c),d=n.getAttribute(c);y0(h,d,o);let u=h.itemSize;for(let f=h.count,g=l;f<g;f++){let _=o+f;for(let m=0;m<u;m++)d.setComponent(_,m,0)}d.needsUpdate=!0,d.addUpdateRange(o*u,l*u)}if(i){let{indexStart:c,reservedIndexCount:h}=a;a.indexCount=e.getIndex().count;for(let d=0;d<r.count;d++)s.setX(c+d,o+r.getX(d));for(let d=r.count,u=h;d<u;d++)s.setX(c+d,o);s.needsUpdate=!0,s.addUpdateRange(c,a.reservedIndexCount)}if(a.start=i?a.indexStart:a.vertexStart,a.count=i?a.indexCount:a.vertexCount,a.boundingBox=null,e.boundingBox!==null)a.boundingBox=e.boundingBox.clone();if(a.boundingSphere=null,e.boundingSphere!==null)a.boundingSphere=e.boundingSphere.clone();return this._visibilityChanged=!0,t}deleteGeometry(t){let e=this._geometryInfo;if(t>=e.length||e[t].active===!1)return this;let n=this._instanceInfo;for(let i=0,s=n.length;i<s;i++)if(n[i].active&&n[i].geometryIndex===t)this.deleteInstance(i);return e[t].active=!1,this._availableGeometryIds.push(t),this._visibilityChanged=!0,this}deleteInstance(t){return this.validateInstanceId(t),this._instanceInfo[t].active=!1,this._availableInstanceIds.push(t),this._visibilityChanged=!0,this}optimize(){let t=0,e=0,n=this._geometryInfo,i=n.map((r,a)=>a).sort((r,a)=>n[r].vertexStart-n[a].vertexStart),s=this.geometry;for(let r=0,a=n.length;r<a;r++){let o=i[r],l=n[o];if(l.active===!1)continue;if(s.index!==null){if(l.indexStart!==e){let{indexStart:c,vertexStart:h,reservedIndexCount:d}=l,u=s.index,f=u.array,g=t-h;for(let _=c;_<c+d;_++)f[_]=f[_]+g;u.array.copyWithin(e,c,c+d),u.addUpdateRange(e,d),u.needsUpdate=!0,l.indexStart=e}e+=l.reservedIndexCount}if(l.vertexStart!==t){let{vertexStart:c,reservedVertexCount:h}=l,d=s.attributes;for(let u in d){let f=d[u],{array:g,itemSize:_}=f;g.copyWithin(t*_,c*_,(c+h)*_),f.addUpdateRange(t*_,h*_),f.needsUpdate=!0}l.vertexStart=t}t+=l.reservedVertexCount,l.start=s.index?l.indexStart:l.vertexStart}return this._nextIndexStart=e,this._nextVertexStart=t,this._visibilityChanged=!0,this}getBoundingBoxAt(t,e){if(t>=this._geometryCount)return null;let n=this.geometry,i=this._geometryInfo[t];if(i.boundingBox===null){let s=new Oe,r=n.index,a=n.attributes.position;for(let o=i.start,l=i.start+i.count;o<l;o++){let c=o;if(r)c=r.getX(c);s.expandByPoint(qs.fromBufferAttribute(a,c))}i.boundingBox=s}return e.copy(i.boundingBox),e}getBoundingSphereAt(t,e){if(t>=this._geometryCount)return null;let n=this.geometry,i=this._geometryInfo[t];if(i.boundingSphere===null){let s=new Ue;this.getBoundingBoxAt(t,aa),aa.getCenter(s.center);let r=n.index,a=n.attributes.position,o=0;for(let l=i.start,c=i.start+i.count;l<c;l++){let h=l;if(r)h=r.getX(h);qs.fromBufferAttribute(a,h),o=Math.max(o,s.center.distanceToSquared(qs))}s.radius=Math.sqrt(o),i.boundingSphere=s}return e.copy(i.boundingSphere),e}setMatrixAt(t,e){this.validateInstanceId(t);let n=this._matricesTexture,i=this._matricesTexture.image.data;return e.toArray(i,t*16),n.needsUpdate=!0,this}getMatrixAt(t,e){return this.validateInstanceId(t),e.fromArray(this._matricesTexture.image.data,t*16)}setColorAt(t,e){if(this.validateInstanceId(t),this._colorsTexture===null)this._initColorsTexture();return e.toArray(this._colorsTexture.image.data,t*4),this._colorsTexture.needsUpdate=!0,this}getColorAt(t,e){if(this.validateInstanceId(t),this._colorsTexture===null)if(e.isVector4)return e.set(1,1,1,1);else return e.setRGB(1,1,1);else return e.fromArray(this._colorsTexture.image.data,t*4)}setVisibleAt(t,e){if(this.validateInstanceId(t),this._instanceInfo[t].visible===e)return this;return this._instanceInfo[t].visible=e,this._visibilityChanged=!0,this}getVisibleAt(t){return this.validateInstanceId(t),this._instanceInfo[t].visible}setGeometryIdAt(t,e){return this.validateInstanceId(t),this.validateGeometryId(e),this._instanceInfo[t].geometryIndex=e,this._visibilityChanged=!0,this}getGeometryIdAt(t){return this.validateInstanceId(t),this._instanceInfo[t].geometryIndex}getGeometryRangeAt(t,e={}){this.validateGeometryId(t);let n=this._geometryInfo[t];return e.vertexStart=n.vertexStart,e.vertexCount=n.vertexCount,e.reservedVertexCount=n.reservedVertexCount,e.indexStart=n.indexStart,e.indexCount=n.indexCount,e.reservedIndexCount=n.reservedIndexCount,e.start=n.start,e.count=n.count,e}setInstanceCount(t){let e=this._availableInstanceIds,n=this._instanceInfo;e.sort(ol);while(e[e.length-1]===n.length-1)n.pop(),e.pop();if(t<n.length)throw Error(`THREE.BatchedMesh: Instance ids outside the range ${t} are being used. Cannot shrink instance count.`);let i=new Int32Array(t),s=new Int32Array(t);Ei(this._multiDrawCounts,i),Ei(this._multiDrawStarts,s),this._multiDrawCounts=i,this._multiDrawStarts=s,this._maxInstanceCount=t;let r=this._indirectTexture,a=this._matricesTexture,o=this._colorsTexture;if(r.dispose(),this._initIndirectTexture(),Ei(r.image.data,this._indirectTexture.image.data),a.dispose(),this._initMatricesTexture(),Ei(a.image.data,this._matricesTexture.image.data),o)o.dispose(),this._initColorsTexture(),Ei(o.image.data,this._colorsTexture.image.data)}setGeometrySize(t,e){let n=[...this._geometryInfo].filter((a)=>a.active);if(Math.max(...n.map((a)=>a.vertexStart+a.reservedVertexCount))>t)throw Error(`THREE.BatchedMesh: Geometry vertex values are being used outside the range ${e}. Cannot shrink further.`);if(this.geometry.index){if(Math.max(...n.map((o)=>o.indexStart+o.reservedIndexCount))>e)throw Error(`THREE.BatchedMesh: Geometry index values are being used outside the range ${e}. Cannot shrink further.`)}let s=this.geometry;if(s.dispose(),this._maxVertexCount=t,this._maxIndexCount=e,this._geometryInitialized)this._geometryInitialized=!1,this.geometry=new jt,this._initializeGeometry(s);let r=this.geometry;if(s.index)Ei(s.index.array,r.index.array);for(let a in s.attributes)Ei(s.attributes[a].array,r.attributes[a].array)}raycast(t,e){let n=this._instanceInfo,i=this._geometryInfo,s=this.matrixWorld,r=this.geometry;if(He.material=this.material,He.geometry.index=r.index,He.geometry.attributes=r.attributes,He.geometry.boundingBox===null)He.geometry.boundingBox=new Oe;if(He.geometry.boundingSphere===null)He.geometry.boundingSphere=new Ue;for(let a=0,o=n.length;a<o;a++){if(!n[a].visible||!n[a].active)continue;let l=n[a].geometryIndex,c=i[l];He.geometry.setDrawRange(c.start,c.count),this.getMatrixAt(a,He.matrixWorld).premultiply(s),this.getBoundingBoxAt(l,He.geometry.boundingBox),this.getBoundingSphereAt(l,He.geometry.boundingSphere),He.raycast(t,oa);for(let h=0,d=oa.length;h<d;h++){let u=oa[h];u.object=this,u.batchId=a,e.push(u)}oa.length=0}He.material=null,He.geometry.index=null,He.geometry.attributes={},He.geometry.setDrawRange(0,1/0)}copy(t){if(super.copy(t),this.geometry=t.geometry.clone(),this.perObjectFrustumCulled=t.perObjectFrustumCulled,this.sortObjects=t.sortObjects,this.boundingBox=t.boundingBox!==null?t.boundingBox.clone():null,this.boundingSphere=t.boundingSphere!==null?t.boundingSphere.clone():null,this._geometryInfo=t._geometryInfo.map((e)=>({...e,boundingBox:e.boundingBox!==null?e.boundingBox.clone():null,boundingSphere:e.boundingSphere!==null?e.boundingSphere.clone():null})),this._instanceInfo=t._instanceInfo.map((e)=>({...e})),this._availableInstanceIds=t._availableInstanceIds.slice(),this._availableGeometryIds=t._availableGeometryIds.slice(),this._nextIndexStart=t._nextIndexStart,this._nextVertexStart=t._nextVertexStart,this._geometryCount=t._geometryCount,this._maxInstanceCount=t._maxInstanceCount,this._maxVertexCount=t._maxVertexCount,this._maxIndexCount=t._maxIndexCount,this._geometryInitialized=t._geometryInitialized,this._multiDrawCounts=t._multiDrawCounts.slice(),this._multiDrawStarts=t._multiDrawStarts.slice(),this._multiDrawBytesPerElement=t._multiDrawBytesPerElement,this._indirectTexture=t._indirectTexture.clone(),this._indirectTexture.image.data=this._indirectTexture.image.data.slice(),this._matricesTexture=t._matricesTexture.clone(),this._matricesTexture.image.data=this._matricesTexture.image.data.slice(),this._colorsTexture!==null)this._colorsTexture=t._colorsTexture.clone(),this._colorsTexture.image.data=this._colorsTexture.image.data.slice();return this}dispose(){if(super.dispose(),this.geometry.dispose(),this._matricesTexture.dispose(),this._matricesTexture=null,this._indirectTexture.dispose(),this._indirectTexture=null,this._colorsTexture!==null)this._colorsTexture.dispose(),this._colorsTexture=null}onBeforeRender(t,e,n,i,s){if(!this._visibilityChanged&&!this.perObjectFrustumCulled&&!this.sortObjects)return;let r=i.getIndex(),a=r===null?1:r.array.BYTES_PER_ELEMENT,o=1;if(s.wireframe)o=2,a=i.attributes.position.count>65535?4:2;let l=this._instanceInfo,c=this._multiDrawStarts,h=this._multiDrawCounts,d=this._geometryInfo,u=this.perObjectFrustumCulled,f=this._indirectTexture,g=f.image.data,_=n.isArrayCamera?_0:g0;if(u)if(n.isArrayCamera)_.setFromArrayCamera(n);else $e.multiplyMatrices(n.projectionMatrix,n.matrixWorldInverse).multiply(this.matrixWorld),_.setFromProjectionMatrix($e,n.coordinateSystem,n.reversedDepth);let m=0;if(this.sortObjects){$e.copy(this.matrixWorld).invert(),qs.setFromMatrixPosition(n.matrixWorld).applyMatrix4($e),_u.set(0,0,-1).transformDirection(n.matrixWorld).transformDirection($e);for(let S=0,x=l.length;S<x;S++)if(l[S].visible&&l[S].active){let T=l[S].geometryIndex;this.getMatrixAt(S,$e),this.getBoundingSphereAt(T,Ti).applyMatrix4($e);let I=!1;if(u)I=!_.intersectsSphere(Ti);if(!I){let R=d[T],y=x0.subVectors(Ti.center,qs).dot(_u);ll.push(R.start,R.count,y,S)}}let p=ll.list,b=this.customSort;if(b===null)p.sort(s.transparent?p0:f0);else b.call(this,p,n);for(let S=0,x=p.length;S<x;S++){let T=p[S];c[m]=T.start*a*o,h[m]=T.count*o,g[m]=T.index,m++}ll.reset()}else for(let p=0,b=l.length;p<b;p++)if(l[p].visible&&l[p].active){let S=l[p].geometryIndex,x=!1;if(u)this.getMatrixAt(p,$e),this.getBoundingSphereAt(S,Ti).applyMatrix4($e),x=!_.intersectsSphere(Ti);if(!x){let T=d[S];c[m]=T.start*a*o,h[m]=T.count*o,g[m]=p,m++}}f.needsUpdate=!0,this._multiDrawCount=m,this._multiDrawBytesPerElement=a,this._visibilityChanged=!1}onBeforeShadow(t,e,n,i,s,r){this.onBeforeRender(t,null,i,s,r)}}class We extends Fe{constructor(t){super();this.isLineBasicMaterial=!0,this.type="LineBasicMaterial",this.color=new Pt(16777215),this.map=null,this.linewidth=1,this.linecap="round",this.linejoin="round",this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.linewidth=t.linewidth,this.linecap=t.linecap,this.linejoin=t.linejoin,this.fog=t.fog,this}}var Aa=new L,wa=new L,xu=new Kt,Ys=new Kn,la=new Ue,cl=new L,yu=new L;class Zn extends ae{constructor(t=new jt,e=new We){super();this.isLine=!0,this.type="Line",this.geometry=t,this.material=e,this.morphTargetDictionary=void 0,this.morphTargetInfluences=void 0,this.updateMorphTargets()}copy(t,e){return super.copy(t,e),this.material=Array.isArray(t.material)?t.material.slice():t.material,this.geometry=t.geometry,this}computeLineDistances(){let t=this.geometry;if(t.index===null){let e=t.attributes.position,n=[0];for(let i=1,s=e.count;i<s;i++)Aa.fromBufferAttribute(e,i-1),wa.fromBufferAttribute(e,i),n[i]=n[i-1],n[i]+=Aa.distanceTo(wa);t.setAttribute("lineDistance",new Nt(n,1))}else Mt("Line.computeLineDistances(): Computation only possible with non-indexed BufferGeometry.");return this}intersectsFrustum(t){return t.intersectsObject(this)}raycast(t,e){let n=this.geometry,i=this.matrixWorld,s=t.params.Line.threshold,r=n.drawRange;if(n.boundingSphere===null)n.computeBoundingSphere();if(la.copy(n.boundingSphere),la.applyMatrix4(i),la.radius+=s,t.ray.intersectsSphere(la)===!1)return;xu.copy(i).invert(),Ys.copy(t.ray).applyMatrix4(xu);let a=s/((this.scale.x+this.scale.y+this.scale.z)/3),o=a*a,l=this.isLineSegments?2:1,c=n.index,d=n.attributes.position;if(c!==null){let u=Math.max(0,r.start),f=Math.min(c.count,r.start+r.count);for(let g=u,_=f-1;g<_;g+=l){let m=c.getX(g),p=c.getX(g+1),b=ca(this,t,Ys,o,m,p,g);if(b)e.push(b)}if(this.isLineLoop){let g=c.getX(f-1),_=c.getX(u),m=ca(this,t,Ys,o,g,_,f-1);if(m)e.push(m)}}else{let u=Math.max(0,r.start),f=Math.min(d.count,r.start+r.count);for(let g=u,_=f-1;g<_;g+=l){let m=ca(this,t,Ys,o,g,g+1,g);if(m)e.push(m)}if(this.isLineLoop){let g=ca(this,t,Ys,o,f-1,u,f-1);if(g)e.push(g)}}}updateMorphTargets(){let e=this.geometry.morphAttributes,n=Object.keys(e);if(n.length>0){let i=e[n[0]];if(i!==void 0){this.morphTargetInfluences=[],this.morphTargetDictionary={};for(let s=0,r=i.length;s<r;s++){let a=i[s].name||String(s);this.morphTargetInfluences.push(0),this.morphTargetDictionary[a]=s}}}}}function ca(t,e,n,i,s,r,a){let o=t.geometry.attributes.position;if(Aa.fromBufferAttribute(o,s),wa.fromBufferAttribute(o,r),n.distanceSqToSegment(Aa,wa,cl,yu)>i)return;cl.applyMatrix4(t.matrixWorld);let c=e.ray.origin.distanceTo(cl);if(c<e.near||c>e.far)return;return{distance:c,point:yu.clone().applyMatrix4(t.matrixWorld),index:a,face:null,faceIndex:null,barycoord:null,object:t}}var vu=new L,Su=new L;class En extends Zn{constructor(t,e){super(t,e);this.isLineSegments=!0,this.type="LineSegments"}computeLineDistances(){let t=this.geometry;if(t.index===null){let e=t.attributes.position,n=[];for(let i=0,s=e.count;i<s;i+=2)vu.fromBufferAttribute(e,i),Su.fromBufferAttribute(e,i+1),n[i]=i===0?0:n[i-1],n[i+1]=n[i]+vu.distanceTo(Su);t.setAttribute("lineDistance",new Nt(n,1))}else Mt("LineSegments.computeLineDistances(): Computation only possible with non-indexed BufferGeometry.");return this}}class Lc extends Zn{constructor(t,e){super(t,e);this.isLineLoop=!0,this.type="LineLoop"}}class Qa extends Fe{constructor(t){super();this.isPointsMaterial=!0,this.type="PointsMaterial",this.color=new Pt(16777215),this.map=null,this.alphaMap=null,this.size=1,this.sizeAttenuation=!0,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.alphaMap=t.alphaMap,this.size=t.size,this.sizeAttenuation=t.sizeAttenuation,this.fog=t.fog,this}}var Mu=new Kt,Sl=new Kn,ha=new Ue,ua=new L;class Nc extends ae{constructor(t=new jt,e=new Qa){super();this.isPoints=!0,this.type="Points",this.geometry=t,this.material=e,this.morphTargetDictionary=void 0,this.morphTargetInfluences=void 0,this.updateMorphTargets()}copy(t,e){return super.copy(t,e),this.material=Array.isArray(t.material)?t.material.slice():t.material,this.geometry=t.geometry,this}intersectsFrustum(t){return t.intersectsObject(this)}raycast(t,e){let n=this.geometry,i=this.matrixWorld,s=t.params.Points.threshold,r=n.drawRange;if(n.boundingSphere===null)n.computeBoundingSphere();if(ha.copy(n.boundingSphere),ha.applyMatrix4(i),ha.radius+=s,t.ray.intersectsSphere(ha)===!1)return;Mu.copy(i).invert(),Sl.copy(t.ray).applyMatrix4(Mu);let a=s/((this.scale.x+this.scale.y+this.scale.z)/3),o=a*a,l=n.index,h=n.attributes.position;if(l!==null){let d=Math.max(0,r.start),u=Math.min(l.count,r.start+r.count);for(let f=d,g=u;f<g;f++){let _=l.getX(f);ua.fromBufferAttribute(h,_),bu(ua,_,o,i,t,e,this)}}else{let d=Math.max(0,r.start),u=Math.min(h.count,r.start+r.count);for(let f=d,g=u;f<g;f++)ua.fromBufferAttribute(h,f),bu(ua,f,o,i,t,e,this)}}updateMorphTargets(){let e=this.geometry.morphAttributes,n=Object.keys(e);if(n.length>0){let i=e[n[0]];if(i!==void 0){this.morphTargetInfluences=[],this.morphTargetDictionary={};for(let s=0,r=i.length;s<r;s++){let a=i[s].name||String(s);this.morphTargetInfluences.push(0),this.morphTargetDictionary[a]=s}}}}}function bu(t,e,n,i,s,r,a){let o=Sl.distanceSqToPoint(t);if(o<n){let l=new L;Sl.closestPointToPoint(t,l),l.applyMatrix4(i);let c=s.ray.origin.distanceTo(l);if(c<s.near||c>s.far)return;r.push({distance:c,distanceToRay:Math.sqrt(o),point:l,index:e,face:null,faceIndex:null,barycoord:null,object:a})}}class Dc extends be{constructor(t,e,n,i,s=1006,r=1006,a,o,l){super(t,e,n,i,s,r,a,o,l);this.isVideoTexture=!0,this.generateMipmaps=!1,this._requestVideoFrameCallbackId=0;let c=this;function h(){c.needsUpdate=!0,c._requestVideoFrameCallbackId=t.requestVideoFrameCallback(h)}if("requestVideoFrameCallback"in t)this._requestVideoFrameCallbackId=t.requestVideoFrameCallback(h)}clone(){return new this.constructor(this.image).copy(this)}update(){let t=this.image;if("requestVideoFrameCallback"in t===!1&&t.readyState>=t.HAVE_CURRENT_DATA)this.needsUpdate=!0}dispose(){if(this._requestVideoFrameCallbackId!==0)this.source.data.cancelVideoFrameCallback(this._requestVideoFrameCallbackId),this._requestVideoFrameCallbackId=0;super.dispose()}}class xf extends Dc{constructor(t,e,n,i,s,r,a,o){super({},t,e,n,i,s,r,a,o);this.isVideoFrameTexture=!0}update(){}clone(){return new this.constructor().copy(this)}setFrame(t){this.image=t,this.needsUpdate=!0}}class yf extends be{constructor(t,e){super({width:t,height:e});this.isFramebufferTexture=!0,this.magFilter=1003,this.minFilter=1003,this.generateMipmaps=!1,this.needsUpdate=!0}}class mr extends be{constructor(t,e,n,i,s,r,a,o,l,c,h,d){super(null,r,a,o,l,c,i,s,h,d);this.isCompressedTexture=!0,this.image={width:e,height:n},this.mipmaps=t,this.flipY=!1,this.generateMipmaps=!1}}class vf extends mr{constructor(t,e,n,i,s,r){super(t,e,n,s,r);this.isCompressedArrayTexture=!0,this.image.depth=i,this.wrapR=1001,this.layerUpdates=new Set}copy(t){return super.copy(t),this.wrapR=t.wrapR,this}addLayerUpdate(t){this.layerUpdates.add(t)}clearLayerUpdates(){this.layerUpdates.clear()}}class Sf extends mr{constructor(t,e,n){super(void 0,t[0].width,t[0].height,e,n,301);this.isCompressedCubeTexture=!0,this.isCubeTexture=!0,this.image=t}}class Rs extends be{constructor(t=[],e=301,n,i,s,r,a,o,l,c){super(t,e,n,i,s,r,a,o,l,c);this.isCubeTexture=!0,this.flipY=!1}get images(){return this.image}set images(t){this.image=t}}class to extends be{constructor(t,e,n,i,s,r,a,o,l){super(t,e,n,i,s,r,a,o,l);this.isCanvasTexture=!0,this.needsUpdate=!0}}class Mf extends be{constructor(t,e,n,i,s,r,a,o,l){super(t,e,n,i,s,r,a,o,l);this.isHTMLTexture=!0,this.generateMipmaps=!1,this.needsUpdate=!0;let c=t?t.parentNode:null;if(c!==null&&"requestPaint"in c)c.onpaint=()=>{this.needsUpdate=!0},c.requestPaint()}dispose(){let t=this.image?this.image.parentNode:null;if(t!==null&&"onpaint"in t)t.onpaint=null;super.dispose()}}class Hi extends be{constructor(t,e,n=1014,i,s,r,a=1003,o=1003,l,c=1026,h=1){if(c!==1026&&c!==1027)throw Error("THREE.DepthTexture: format must be either THREE.DepthFormat or THREE.DepthStencilFormat");let d={width:t,height:e,depth:h};super(d,i,s,r,a,o,c,n,l);this.isDepthTexture=!0,this.flipY=!1,this.generateMipmaps=!1,this.compareFunction=null}copy(t){return super.copy(t),this.source=new In(Object.assign({},t.image)),this.compareFunction=t.compareFunction,this}toJSON(t){let e=super.toJSON(t);return e.compareFunction=this.compareFunction,e}}class Uc extends Hi{constructor(t,e=1014,n=301,i,s,r=1003,a=1003,o,l=1026){let c={width:t,height:t,depth:1},h=[c,c,c,c,c,c];super(t,t,e,n,i,s,r,a,o,l);this.image=h,this.isCubeDepthTexture=!0,this.isCubeTexture=!0}get images(){return this.image}set images(t){this.image=t}}class eo extends be{constructor(t=null){super();this.sourceTexture=t,this.isExternalTexture=!0}copy(t){return super.copy(t),this.sourceTexture=t.sourceTexture,this}}class An extends jt{constructor(t=1,e=1,n=1,i=1,s=1,r=1){super();this.type="BoxGeometry",this.parameters={width:t,height:e,depth:n,widthSegments:i,heightSegments:s,depthSegments:r};let a=this;i=Math.floor(i),s=Math.floor(s),r=Math.floor(r);let o=[],l=[],c=[],h=[],d=0,u=0;f("z","y","x",-1,-1,n,e,t,r,s,0),f("z","y","x",1,-1,n,e,-t,r,s,1),f("x","z","y",1,1,t,n,e,i,r,2),f("x","z","y",1,-1,t,n,-e,i,r,3),f("x","y","z",1,-1,t,e,n,i,s,4),f("x","y","z",-1,-1,t,e,-n,i,s,5),this.setIndex(o),this.setAttribute("position",new Nt(l,3)),this.setAttribute("normal",new Nt(c,3)),this.setAttribute("uv",new Nt(h,2));function f(g,_,m,p,b,S,x,T,I,R,y){let A=S/I,Z=x/R,D=S/2,q=x/2,ot=T/2,z=I+1,et=R+1,nt=0,B=0,ft=new L;for(let Q=0;Q<et;Q++){let ct=Q*Z-q;for(let dt=0;dt<z;dt++){let Bt=dt*A-D;ft[g]=Bt*p,ft[_]=ct*b,ft[m]=ot,l.push(ft.x,ft.y,ft.z),ft[g]=0,ft[_]=0,ft[m]=T>0?1:-1,c.push(ft.x,ft.y,ft.z),h.push(dt/I),h.push(1-Q/R),nt+=1}}for(let Q=0;Q<R;Q++)for(let ct=0;ct<I;ct++){let dt=d+ct+z*Q,Bt=d+ct+z*(Q+1),Dt=d+(ct+1)+z*(Q+1),ee=d+(ct+1)+z*Q;o.push(dt,Bt,ee),o.push(Bt,Dt,ee),B+=6}a.addGroup(u,B,y),u+=B,d+=nt}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new An(t.width,t.height,t.depth,t.widthSegments,t.heightSegments,t.depthSegments)}}class no extends jt{constructor(t=1,e=1,n=4,i=8,s=1){super();this.type="CapsuleGeometry",this.parameters={radius:t,height:e,capSegments:n,radialSegments:i,heightSegments:s},e=Math.max(0,e),n=Math.max(1,Math.floor(n)),i=Math.max(3,Math.floor(i)),s=Math.max(1,Math.floor(s));let r=[],a=[],o=[],l=[],c=e/2,h=Math.PI/2*t,d=e,u=2*h+d,f=n*2+s,g=i+1,_=new L,m=new L;for(let p=0;p<=f;p++){let b=0,S=0,x=0,T=0;if(p<=n){let y=p/n,A=y*Math.PI/2;S=-c-t*Math.cos(A),x=t*Math.sin(A),T=-t*Math.cos(A),b=y*h}else if(p<=n+s){let y=(p-n)/s;S=-c+y*e,x=t,T=0,b=h+y*d}else{let y=(p-n-s)/n,A=y*Math.PI/2;S=c+t*Math.sin(A),x=t*Math.cos(A),T=t*Math.sin(A),b=h+d+y*h}let I=Math.max(0,Math.min(1,b/u)),R=0;if(p===0)R=0.5/i;else if(p===f)R=-0.5/i;for(let y=0;y<=i;y++){let A=y/i,Z=A*Math.PI*2,D=Math.sin(Z),q=Math.cos(Z);m.x=-x*q,m.y=S,m.z=x*D,a.push(m.x,m.y,m.z),_.set(-x*q,T,x*D),_.normalize(),o.push(_.x,_.y,_.z),l.push(A+R,I)}if(p>0){let y=(p-1)*g;for(let A=0;A<i;A++){let Z=y+A,D=y+A+1,q=p*g+A,ot=p*g+A+1;r.push(Z,D,q),r.push(D,ot,q)}}}this.setIndex(r),this.setAttribute("position",new Nt(a,3)),this.setAttribute("normal",new Nt(o,3)),this.setAttribute("uv",new Nt(l,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new no(t.radius,t.height,t.capSegments,t.radialSegments,t.heightSegments)}}class io extends jt{constructor(t=1,e=32,n=0,i=Math.PI*2){super();this.type="CircleGeometry",this.parameters={radius:t,segments:e,thetaStart:n,thetaLength:i},e=Math.max(3,e);let s=[],r=[],a=[],o=[],l=new L,c=new rt;r.push(0,0,0),a.push(0,0,1),o.push(0.5,0.5);for(let h=0,d=3;h<=e;h++,d+=3){let u=n+h/e*i;l.x=t*Math.cos(u),l.y=t*Math.sin(u),r.push(l.x,l.y,l.z),a.push(0,0,1),c.x=(r[d]/t+1)/2,c.y=(r[d+1]/t+1)/2,o.push(c.x,c.y)}for(let h=1;h<=e;h++)s.push(h,h+1,0);this.setIndex(s),this.setAttribute("position",new Nt(r,3)),this.setAttribute("normal",new Nt(a,3)),this.setAttribute("uv",new Nt(o,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new io(t.radius,t.segments,t.thetaStart,t.thetaLength)}}class gr extends jt{constructor(t=1,e=1,n=1,i=32,s=1,r=!1,a=0,o=Math.PI*2){super();this.type="CylinderGeometry",this.parameters={radiusTop:t,radiusBottom:e,height:n,radialSegments:i,heightSegments:s,openEnded:r,thetaStart:a,thetaLength:o};let l=this;i=Math.floor(i),s=Math.floor(s);let c=[],h=[],d=[],u=[],f=0,g=[],_=n/2,m=0;if(p(),r===!1){if(t>0)b(!0);if(e>0)b(!1)}this.setIndex(c),this.setAttribute("position",new Nt(h,3)),this.setAttribute("normal",new Nt(d,3)),this.setAttribute("uv",new Nt(u,2));function p(){let S=new L,x=new L,T=0,I=(e-t)/n;for(let R=0;R<=s;R++){let y=[],A=R/s,Z=A*(e-t)+t;for(let D=0;D<=i;D++){let q=D/i,ot=q*o+a,z=Math.sin(ot),et=Math.cos(ot);x.x=Z*z,x.y=-A*n+_,x.z=Z*et,h.push(x.x,x.y,x.z),S.set(z,I,et).normalize(),d.push(S.x,S.y,S.z),u.push(q,1-A),y.push(f++)}g.push(y)}for(let R=0;R<i;R++)for(let y=0;y<s;y++){let A=g[y][R],Z=g[y+1][R],D=g[y+1][R+1],q=g[y][R+1];if(t>0||y!==0)c.push(A,Z,q),T+=3;if(e>0||y!==s-1)c.push(Z,D,q),T+=3}l.addGroup(m,T,0),m+=T}function b(S){let x=f,T=new rt,I=new L,R=0,y=S===!0?t:e,A=S===!0?1:-1;for(let D=1;D<=i;D++)h.push(0,_*A,0),d.push(0,A,0),u.push(0.5,0.5),f++;let Z=f;for(let D=0;D<=i;D++){let ot=D/i*o+a,z=Math.cos(ot),et=Math.sin(ot);I.x=y*et,I.y=_*A,I.z=y*z,h.push(I.x,I.y,I.z),d.push(0,A,0),T.x=z*0.5+0.5,T.y=et*0.5*A+0.5,u.push(T.x,T.y),f++}for(let D=0;D<i;D++){let q=x+D,ot=Z+D;if(S===!0)c.push(ot,ot+1,q);else c.push(ot+1,ot,q);R+=3}l.addGroup(m,R,S===!0?1:2),m+=R}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new gr(t.radiusTop,t.radiusBottom,t.height,t.radialSegments,t.heightSegments,t.openEnded,t.thetaStart,t.thetaLength)}}class _r extends gr{constructor(t=1,e=1,n=32,i=1,s=!1,r=0,a=Math.PI*2){super(0,t,e,n,i,s,r,a);this.type="ConeGeometry",this.parameters={radius:t,height:e,radialSegments:n,heightSegments:i,openEnded:s,thetaStart:r,thetaLength:a}}static fromJSON(t){return new _r(t.radius,t.height,t.radialSegments,t.heightSegments,t.openEnded,t.thetaStart,t.thetaLength)}}class gi extends jt{constructor(t=[],e=[],n=1,i=0){super();this.type="PolyhedronGeometry",this.parameters={vertices:t,indices:e,radius:n,detail:i};let s=[],r=[];if(a(i),l(n),c(),this.setAttribute("position",new Nt(s,3)),this.setAttribute("normal",new Nt(s.slice(),3)),this.setAttribute("uv",new Nt(r,2)),i===0)this.computeVertexNormals();else this.normalizeNormals();function a(p){let b=new L,S=new L,x=new L;for(let T=0;T<e.length;T+=3)u(e[T+0],b),u(e[T+1],S),u(e[T+2],x),o(b,S,x,p)}function o(p,b,S,x){let T=x+1,I=[];for(let R=0;R<=T;R++){I[R]=[];let y=p.clone().lerp(S,R/T),A=b.clone().lerp(S,R/T),Z=T-R;for(let D=0;D<=Z;D++)if(D===0&&R===T)I[R][D]=y;else I[R][D]=y.clone().lerp(A,D/Z)}for(let R=0;R<T;R++)for(let y=0;y<2*(T-R)-1;y++){let A=Math.floor(y/2);if(y%2===0)d(I[R][A+1]),d(I[R+1][A]),d(I[R][A]);else d(I[R][A+1]),d(I[R+1][A+1]),d(I[R+1][A])}}function l(p){let b=new L;for(let S=0;S<s.length;S+=3)b.x=s[S+0],b.y=s[S+1],b.z=s[S+2],b.normalize().multiplyScalar(p),s[S+0]=b.x,s[S+1]=b.y,s[S+2]=b.z}function c(){let p=new L;for(let b=0;b<s.length;b+=3){p.x=s[b+0],p.y=s[b+1],p.z=s[b+2];let S=_(p)/2/Math.PI+0.5,x=m(p)/Math.PI+0.5;r.push(S,1-x)}f(),h()}function h(){for(let p=0;p<r.length;p+=6){let b=r[p+0],S=r[p+2],x=r[p+4],T=Math.max(b,S,x),I=Math.min(b,S,x);if(T>0.9&&I<0.1){if(b<0.2)r[p+0]+=1;if(S<0.2)r[p+2]+=1;if(x<0.2)r[p+4]+=1}}}function d(p){s.push(p.x,p.y,p.z)}function u(p,b){let S=p*3;b.x=t[S+0],b.y=t[S+1],b.z=t[S+2]}function f(){let p=new L,b=new L,S=new L,x=new L,T=new rt,I=new rt,R=new rt;for(let y=0,A=0;y<s.length;y+=9,A+=6){p.set(s[y+0],s[y+1],s[y+2]),b.set(s[y+3],s[y+4],s[y+5]),S.set(s[y+6],s[y+7],s[y+8]),T.set(r[A+0],r[A+1]),I.set(r[A+2],r[A+3]),R.set(r[A+4],r[A+5]),x.copy(p).add(b).add(S).divideScalar(3);let Z=_(x);g(T,A+0,p,Z),g(I,A+2,b,Z),g(R,A+4,S,Z)}}function g(p,b,S,x){if(x<0&&p.x===1)r[b]=p.x-1;if(S.x===0&&S.z===0)r[b]=x/2/Math.PI+0.5}function _(p){return Math.atan2(p.z,-p.x)}function m(p){return Math.atan2(-p.y,Math.sqrt(p.x*p.x+p.z*p.z))}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new gi(t.vertices,t.indices,t.radius,t.detail)}}class so extends gi{constructor(t=1,e=0){let n=(1+Math.sqrt(5))/2,i=1/n,s=[-1,-1,-1,-1,-1,1,-1,1,-1,-1,1,1,1,-1,-1,1,-1,1,1,1,-1,1,1,1,0,-i,-n,0,-i,n,0,i,-n,0,i,n,-i,-n,0,-i,n,0,i,-n,0,i,n,0,-n,0,-i,n,0,-i,-n,0,i,n,0,i],r=[3,11,7,3,7,15,3,15,13,7,19,17,7,17,6,7,6,15,17,4,8,17,8,10,17,10,6,8,0,16,8,16,2,8,2,10,0,12,1,0,1,18,0,18,16,6,10,2,6,2,13,6,13,15,2,16,18,2,18,3,2,3,13,18,1,9,18,9,11,18,11,3,4,14,12,4,12,0,4,0,8,11,9,5,11,5,19,11,19,7,19,5,14,19,14,4,19,4,17,1,12,14,1,14,5,1,5,9];super(s,r,t,e);this.type="DodecahedronGeometry",this.parameters={radius:t,detail:e}}static fromJSON(t){return new so(t.radius,t.detail)}}var da=new L,fa=new L,hl=new L,pa=new je;class Fc extends jt{constructor(t=null,e=1){super();if(this.type="EdgesGeometry",this.parameters={geometry:t,thresholdAngle:e},t!==null){let i=Math.pow(10,4),s=Math.cos(Ni*e),r=t.getIndex(),a=t.getAttribute("position"),o=r?r.count:a.count,l=[0,0,0],c=["a","b","c"],h=[,,,],d={},u=[];for(let f=0;f<o;f+=3){if(r)l[0]=r.getX(f),l[1]=r.getX(f+1),l[2]=r.getX(f+2);else l[0]=f,l[1]=f+1,l[2]=f+2;let{a:g,b:_,c:m}=pa;if(g.fromBufferAttribute(a,l[0]),_.fromBufferAttribute(a,l[1]),m.fromBufferAttribute(a,l[2]),pa.getNormal(hl),h[0]=`${Math.round(g.x*i)},${Math.round(g.y*i)},${Math.round(g.z*i)}`,h[1]=`${Math.round(_.x*i)},${Math.round(_.y*i)},${Math.round(_.z*i)}`,h[2]=`${Math.round(m.x*i)},${Math.round(m.y*i)},${Math.round(m.z*i)}`,h[0]===h[1]||h[1]===h[2]||h[2]===h[0])continue;for(let p=0;p<3;p++){let b=(p+1)%3,S=h[p],x=h[b],T=pa[c[p]],I=pa[c[b]],R=`${S}_${x}`,y=`${x}_${S}`;if(y in d&&d[y]){if(hl.dot(d[y].normal)<=s)u.push(T.x,T.y,T.z),u.push(I.x,I.y,I.z);d[y]=null}else if(!(R in d))d[R]={index0:l[p],index1:l[b],normal:hl.clone()}}}for(let f in d)if(d[f]){let{index0:g,index1:_}=d[f];da.fromBufferAttribute(a,g),fa.fromBufferAttribute(a,_),u.push(da.x,da.y,da.z),u.push(fa.x,fa.y,fa.z)}this.setAttribute("position",new Nt(u,3))}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}}class gn{constructor(){this.type="Curve",this.arcLengthDivisions=200,this.needsUpdate=!1,this.cacheArcLengths=null}getPoint(){Mt("Curve: .getPoint() not implemented.")}getPointAt(t,e){let n=this.getUtoTmapping(t);return this.getPoint(n,e)}getPoints(t=5){let e=[];for(let n=0;n<=t;n++)e.push(this.getPoint(n/t));return e}getSpacedPoints(t=5){let e=[];for(let n=0;n<=t;n++)e.push(this.getPointAt(n/t));return e}getLength(){let t=this.getLengths();return t[t.length-1]}getLengths(t=this.arcLengthDivisions){if(this.cacheArcLengths&&this.cacheArcLengths.length===t+1&&!this.needsUpdate)return this.cacheArcLengths;this.needsUpdate=!1;let e=[],n,i=this.getPoint(0),s=0;e.push(0);for(let r=1;r<=t;r++)n=this.getPoint(r/t),s+=n.distanceTo(i),e.push(s),i=n;return this.cacheArcLengths=e,e}updateArcLengths(){this.needsUpdate=!0,this.getLengths()}getUtoTmapping(t,e=null){let n=this.getLengths(),i=0,s=n.length,r;if(e)r=e;else r=t*n[s-1];let a=0,o=s-1,l;while(a<=o)if(i=Math.floor(a+(o-a)/2),l=n[i]-r,l<0)a=i+1;else if(l>0)o=i-1;else{o=i;break}if(i=o,n[i]===r)return i/(s-1);let c=n[i],d=n[i+1]-c,u=(r-c)/d;return(i+u)/(s-1)}getTangent(t,e){let i=t-0.0001,s=t+0.0001;if(i<0)i=0;if(s>1)s=1;let r=this.getPoint(i),a=this.getPoint(s),o=e||(r.isVector2?new rt:new L);return o.copy(a).sub(r).normalize(),o}getTangentAt(t,e){let n=this.getUtoTmapping(t);return this.getTangent(n,e)}computeFrenetFrames(t,e=!1){let n=new L,i=[],s=[],r=[],a=new L,o=new Kt;for(let u=0;u<=t;u++){let f=u/t;i[u]=this.getTangentAt(f,new L)}s[0]=new L,r[0]=new L;let l=Number.MAX_VALUE,c=Math.abs(i[0].x),h=Math.abs(i[0].y),d=Math.abs(i[0].z);if(c<=l)l=c,n.set(1,0,0);if(h<=l)l=h,n.set(0,1,0);if(d<=l)n.set(0,0,1);a.crossVectors(i[0],n).normalize(),s[0].crossVectors(i[0],a),r[0].crossVectors(i[0],s[0]);for(let u=1;u<=t;u++){if(s[u]=s[u-1].clone(),r[u]=r[u-1].clone(),a.crossVectors(i[u-1],i[u]),a.length()>Number.EPSILON){a.normalize();let f=Math.acos(Zt(i[u-1].dot(i[u]),-1,1));s[u].applyMatrix4(o.makeRotationAxis(a,f))}r[u].crossVectors(i[u],s[u])}if(e===!0){let u=Math.acos(Zt(s[0].dot(s[t]),-1,1));if(u/=t,i[0].dot(a.crossVectors(s[0],s[t]))>0)u=-u;for(let f=1;f<=t;f++)s[f].applyMatrix4(o.makeRotationAxis(i[f],u*f)),r[f].crossVectors(i[f],s[f])}return{tangents:i,normals:s,binormals:r}}clone(){return new this.constructor().copy(this)}copy(t){return this.arcLengthDivisions=t.arcLengthDivisions,this}toJSON(){let t={metadata:{version:4.7,type:"Curve",generator:"Curve.toJSON"}};return t.arcLengthDivisions=this.arcLengthDivisions,t.type=this.type,t}fromJSON(t){return this.arcLengthDivisions=t.arcLengthDivisions,this}}class xr extends gn{constructor(t=0,e=0,n=1,i=1,s=0,r=Math.PI*2,a=!1,o=0){super();this.isEllipseCurve=!0,this.type="EllipseCurve",this.aX=t,this.aY=e,this.xRadius=n,this.yRadius=i,this.aStartAngle=s,this.aEndAngle=r,this.aClockwise=a,this.aRotation=o}getPoint(t,e=new rt){let n=e,i=Math.PI*2,s=this.aEndAngle-this.aStartAngle,r=Math.abs(s)<Number.EPSILON;while(s<0)s+=i;while(s>i)s-=i;if(s<Number.EPSILON)if(r)s=0;else s=i;if(this.aClockwise===!0&&!r)if(s===i)s=-i;else s=s-i;let a=this.aStartAngle+t*s,o=this.aX+this.xRadius*Math.cos(a),l=this.aY+this.yRadius*Math.sin(a);if(this.aRotation!==0){let c=Math.cos(this.aRotation),h=Math.sin(this.aRotation),d=o-this.aX,u=l-this.aY;o=d*c-u*h+this.aX,l=d*h+u*c+this.aY}return n.set(o,l)}copy(t){return super.copy(t),this.aX=t.aX,this.aY=t.aY,this.xRadius=t.xRadius,this.yRadius=t.yRadius,this.aStartAngle=t.aStartAngle,this.aEndAngle=t.aEndAngle,this.aClockwise=t.aClockwise,this.aRotation=t.aRotation,this}toJSON(){let t=super.toJSON();return t.aX=this.aX,t.aY=this.aY,t.xRadius=this.xRadius,t.yRadius=this.yRadius,t.aStartAngle=this.aStartAngle,t.aEndAngle=this.aEndAngle,t.aClockwise=this.aClockwise,t.aRotation=this.aRotation,t}fromJSON(t){return super.fromJSON(t),this.aX=t.aX,this.aY=t.aY,this.xRadius=t.xRadius,this.yRadius=t.yRadius,this.aStartAngle=t.aStartAngle,this.aEndAngle=t.aEndAngle,this.aClockwise=t.aClockwise,this.aRotation=t.aRotation,this}}class Oc extends xr{constructor(t,e,n,i,s,r){super(t,e,n,n,i,s,r);this.isArcCurve=!0,this.type="ArcCurve"}}function Bc(){let t=0,e=0,n=0,i=0;function s(r,a,o,l){t=r,e=o,n=-3*r+3*a-2*o-l,i=2*r-2*a+o+l}return{initCatmullRom:function(r,a,o,l,c){s(a,o,c*(o-r),c*(l-a))},initNonuniformCatmullRom:function(r,a,o,l,c,h,d){let u=(a-r)/c-(o-r)/(c+h)+(o-a)/h,f=(o-a)/h-(l-a)/(h+d)+(l-o)/d;u*=h,f*=h,s(a,o,u,f)},calc:function(r){let a=r*r,o=a*r;return t+e*r+n*a+i*o}}}var Tu=new L,Eu=new L,ul=new Bc,dl=new Bc,fl=new Bc;class zc extends gn{constructor(t=[],e=!1,n="centripetal",i=0.5){super();this.isCatmullRomCurve3=!0,this.type="CatmullRomCurve3",this.points=t,this.closed=e,this.curveType=n,this.tension=i}getPoint(t,e=new L){let n=e,i=this.points,s=i.length,r=(s-(this.closed?0:1))*t,a=Math.floor(r),o=r-a;if(this.closed)a+=a>0?0:(Math.floor(Math.abs(a)/s)+1)*s;else if(o===0&&a===s-1)a=s-2,o=1;let l,c;if(this.closed||a>0)l=i[(a-1)%s];else Eu.subVectors(i[0],i[1]).add(i[0]),l=Eu;let h=i[a%s],d=i[(a+1)%s];if(this.closed||a+2<s)c=i[(a+2)%s];else Tu.subVectors(i[s-1],i[s-2]).add(i[s-1]),c=Tu;if(this.curveType==="centripetal"||this.curveType==="chordal"){let u=this.curveType==="chordal"?0.5:0.25,f=Math.pow(l.distanceToSquared(h),u),g=Math.pow(h.distanceToSquared(d),u),_=Math.pow(d.distanceToSquared(c),u);if(g<0.0001)g=1;if(f<0.0001)f=g;if(_<0.0001)_=g;ul.initNonuniformCatmullRom(l.x,h.x,d.x,c.x,f,g,_),dl.initNonuniformCatmullRom(l.y,h.y,d.y,c.y,f,g,_),fl.initNonuniformCatmullRom(l.z,h.z,d.z,c.z,f,g,_)}else if(this.curveType==="catmullrom")ul.initCatmullRom(l.x,h.x,d.x,c.x,this.tension),dl.initCatmullRom(l.y,h.y,d.y,c.y,this.tension),fl.initCatmullRom(l.z,h.z,d.z,c.z,this.tension);return n.set(ul.calc(o),dl.calc(o),fl.calc(o)),n}copy(t){super.copy(t),this.points=[];for(let e=0,n=t.points.length;e<n;e++){let i=t.points[e];this.points.push(i.clone())}return this.closed=t.closed,this.curveType=t.curveType,this.tension=t.tension,this}toJSON(){let t=super.toJSON();t.points=[];for(let e=0,n=this.points.length;e<n;e++){let i=this.points[e];t.points.push(i.toArray())}return t.closed=this.closed,t.curveType=this.curveType,t.tension=this.tension,t}fromJSON(t){super.fromJSON(t),this.points=[];for(let e=0,n=t.points.length;e<n;e++){let i=t.points[e];this.points.push(new L().fromArray(i))}return this.closed=t.closed,this.curveType=t.curveType,this.tension=t.tension,this}}function Au(t,e,n,i,s){let r=(i-e)*0.5,a=(s-n)*0.5,o=t*t,l=t*o;return(2*n-2*i+r+a)*l+(-3*n+3*i-2*r-a)*o+r*t+n}function v0(t,e){let n=1-t;return n*n*e}function S0(t,e){return 2*(1-t)*t*e}function M0(t,e){return t*t*e}function Ks(t,e,n,i){return v0(t,e)+S0(t,n)+M0(t,i)}function b0(t,e){let n=1-t;return n*n*n*e}function T0(t,e){let n=1-t;return 3*n*n*t*e}function E0(t,e){return 3*(1-t)*t*t*e}function A0(t,e){return t*t*t*e}function js(t,e,n,i,s){return b0(t,e)+T0(t,n)+E0(t,i)+A0(t,s)}class ro extends gn{constructor(t=new rt,e=new rt,n=new rt,i=new rt){super();this.isCubicBezierCurve=!0,this.type="CubicBezierCurve",this.v0=t,this.v1=e,this.v2=n,this.v3=i}getPoint(t,e=new rt){let n=e,i=this.v0,s=this.v1,r=this.v2,a=this.v3;return n.set(js(t,i.x,s.x,r.x,a.x),js(t,i.y,s.y,r.y,a.y)),n}copy(t){return super.copy(t),this.v0.copy(t.v0),this.v1.copy(t.v1),this.v2.copy(t.v2),this.v3.copy(t.v3),this}toJSON(){let t=super.toJSON();return t.v0=this.v0.toArray(),t.v1=this.v1.toArray(),t.v2=this.v2.toArray(),t.v3=this.v3.toArray(),t}fromJSON(t){return super.fromJSON(t),this.v0.fromArray(t.v0),this.v1.fromArray(t.v1),this.v2.fromArray(t.v2),this.v3.fromArray(t.v3),this}}class kc extends gn{constructor(t=new L,e=new L,n=new L,i=new L){super();this.isCubicBezierCurve3=!0,this.type="CubicBezierCurve3",this.v0=t,this.v1=e,this.v2=n,this.v3=i}getPoint(t,e=new L){let n=e,i=this.v0,s=this.v1,r=this.v2,a=this.v3;return n.set(js(t,i.x,s.x,r.x,a.x),js(t,i.y,s.y,r.y,a.y),js(t,i.z,s.z,r.z,a.z)),n}copy(t){return super.copy(t),this.v0.copy(t.v0),this.v1.copy(t.v1),this.v2.copy(t.v2),this.v3.copy(t.v3),this}toJSON(){let t=super.toJSON();return t.v0=this.v0.toArray(),t.v1=this.v1.toArray(),t.v2=this.v2.toArray(),t.v3=this.v3.toArray(),t}fromJSON(t){return super.fromJSON(t),this.v0.fromArray(t.v0),this.v1.fromArray(t.v1),this.v2.fromArray(t.v2),this.v3.fromArray(t.v3),this}}class ao extends gn{constructor(t=new rt,e=new rt){super();this.isLineCurve=!0,this.type="LineCurve",this.v1=t,this.v2=e}getPoint(t,e=new rt){let n=e;if(t===1)n.copy(this.v2);else n.copy(this.v2).sub(this.v1),n.multiplyScalar(t).add(this.v1);return n}getPointAt(t,e){return this.getPoint(t,e)}getTangent(t,e=new rt){return e.subVectors(this.v2,this.v1).normalize()}getTangentAt(t,e){return this.getTangent(t,e)}copy(t){return super.copy(t),this.v1.copy(t.v1),this.v2.copy(t.v2),this}toJSON(){let t=super.toJSON();return t.v1=this.v1.toArray(),t.v2=this.v2.toArray(),t}fromJSON(t){return super.fromJSON(t),this.v1.fromArray(t.v1),this.v2.fromArray(t.v2),this}}class Gc extends gn{constructor(t=new L,e=new L){super();this.isLineCurve3=!0,this.type="LineCurve3",this.v1=t,this.v2=e}getPoint(t,e=new L){let n=e;if(t===1)n.copy(this.v2);else n.copy(this.v2).sub(this.v1),n.multiplyScalar(t).add(this.v1);return n}getPointAt(t,e){return this.getPoint(t,e)}getTangent(t,e=new L){return e.subVectors(this.v2,this.v1).normalize()}getTangentAt(t,e){return this.getTangent(t,e)}copy(t){return super.copy(t),this.v1.copy(t.v1),this.v2.copy(t.v2),this}toJSON(){let t=super.toJSON();return t.v1=this.v1.toArray(),t.v2=this.v2.toArray(),t}fromJSON(t){return super.fromJSON(t),this.v1.fromArray(t.v1),this.v2.fromArray(t.v2),this}}class oo extends gn{constructor(t=new rt,e=new rt,n=new rt){super();this.isQuadraticBezierCurve=!0,this.type="QuadraticBezierCurve",this.v0=t,this.v1=e,this.v2=n}getPoint(t,e=new rt){let n=e,i=this.v0,s=this.v1,r=this.v2;return n.set(Ks(t,i.x,s.x,r.x),Ks(t,i.y,s.y,r.y)),n}copy(t){return super.copy(t),this.v0.copy(t.v0),this.v1.copy(t.v1),this.v2.copy(t.v2),this}toJSON(){let t=super.toJSON();return t.v0=this.v0.toArray(),t.v1=this.v1.toArray(),t.v2=this.v2.toArray(),t}fromJSON(t){return super.fromJSON(t),this.v0.fromArray(t.v0),this.v1.fromArray(t.v1),this.v2.fromArray(t.v2),this}}class lo extends gn{constructor(t=new L,e=new L,n=new L){super();this.isQuadraticBezierCurve3=!0,this.type="QuadraticBezierCurve3",this.v0=t,this.v1=e,this.v2=n}getPoint(t,e=new L){let n=e,i=this.v0,s=this.v1,r=this.v2;return n.set(Ks(t,i.x,s.x,r.x),Ks(t,i.y,s.y,r.y),Ks(t,i.z,s.z,r.z)),n}copy(t){return super.copy(t),this.v0.copy(t.v0),this.v1.copy(t.v1),this.v2.copy(t.v2),this}toJSON(){let t=super.toJSON();return t.v0=this.v0.toArray(),t.v1=this.v1.toArray(),t.v2=this.v2.toArray(),t}fromJSON(t){return super.fromJSON(t),this.v0.fromArray(t.v0),this.v1.fromArray(t.v1),this.v2.fromArray(t.v2),this}}class co extends gn{constructor(t=[]){super();this.isSplineCurve=!0,this.type="SplineCurve",this.points=t}getPoint(t,e=new rt){let n=e,i=this.points,s=(i.length-1)*t,r=Math.floor(s),a=s-r,o=i[r===0?r:r-1],l=i[r],c=i[r>i.length-2?i.length-1:r+1],h=i[r>i.length-3?i.length-1:r+2];return n.set(Au(a,o.x,l.x,c.x,h.x),Au(a,o.y,l.y,c.y,h.y)),n}copy(t){super.copy(t),this.points=[];for(let e=0,n=t.points.length;e<n;e++){let i=t.points[e];this.points.push(i.clone())}return this}toJSON(){let t=super.toJSON();t.points=[];for(let e=0,n=this.points.length;e<n;e++){let i=this.points[e];t.points.push(i.toArray())}return t}fromJSON(t){super.fromJSON(t),this.points=[];for(let e=0,n=t.points.length;e<n;e++){let i=t.points[e];this.points.push(new rt().fromArray(i))}return this}}var Ra=Object.freeze({__proto__:null,ArcCurve:Oc,CatmullRomCurve3:zc,CubicBezierCurve:ro,CubicBezierCurve3:kc,EllipseCurve:xr,LineCurve:ao,LineCurve3:Gc,QuadraticBezierCurve:oo,QuadraticBezierCurve3:lo,SplineCurve:co});class Hc extends gn{constructor(){super();this.type="CurvePath",this.curves=[],this.autoClose=!1}add(t){this.curves.push(t)}closePath(){let t=this.curves[0].getPoint(0),e=this.curves[this.curves.length-1].getPoint(1);if(!t.equals(e)){let n=t.isVector2===!0?"LineCurve":"LineCurve3";this.curves.push(new Ra[n](e,t))}return this}getPoint(t,e){let n=t*this.getLength(),i=this.getCurveLengths(),s=0;while(s<i.length){if(i[s]>=n){let r=i[s]-n,a=this.curves[s],o=a.getLength(),l=o===0?0:1-r/o;return a.getPointAt(l,e)}s++}return null}getLength(){let t=this.getCurveLengths();return t[t.length-1]}updateArcLengths(){this.needsUpdate=!0,this.cacheLengths=null,this.getCurveLengths()}getCurveLengths(){if(this.cacheLengths&&this.cacheLengths.length===this.curves.length)return this.cacheLengths;let t=[],e=0;for(let n=0,i=this.curves.length;n<i;n++)e+=this.curves[n].getLength(),t.push(e);return this.cacheLengths=t,t}getSpacedPoints(t=40){let e=[];for(let n=0;n<=t;n++)e.push(this.getPoint(n/t));if(this.autoClose)e.push(e[0]);return e}getPoints(t=12){let e=[],n;for(let i=0,s=this.curves;i<s.length;i++){let r=s[i],a=r.isEllipseCurve?t*2:r.isLineCurve||r.isLineCurve3?1:r.isSplineCurve?t*r.points.length:t,o=r.getPoints(a);for(let l=0;l<o.length;l++){let c=o[l];if(n&&n.equals(c))continue;e.push(c),n=c}}if(this.autoClose&&e.length>1&&!e[e.length-1].equals(e[0]))e.push(e[0]);return e}copy(t){super.copy(t),this.curves=[];for(let e=0,n=t.curves.length;e<n;e++){let i=t.curves[e];this.curves.push(i.clone())}return this.autoClose=t.autoClose,this}toJSON(){let t=super.toJSON();t.autoClose=this.autoClose,t.curves=[];for(let e=0,n=this.curves.length;e<n;e++){let i=this.curves[e];t.curves.push(i.toJSON())}return t}fromJSON(t){super.fromJSON(t),this.autoClose=t.autoClose,this.curves=[];for(let e=0,n=t.curves.length;e<n;e++){let i=t.curves[e];this.curves.push(new Ra[i.type]().fromJSON(i))}return this}}class Jn extends Hc{constructor(t){super();if(this.type="Path",this.currentPoint=new rt,t)this.setFromPoints(t)}setFromPoints(t){this.moveTo(t[0].x,t[0].y);for(let e=1,n=t.length;e<n;e++)this.lineTo(t[e].x,t[e].y);return this}moveTo(t,e){return this.currentPoint.set(t,e),this}lineTo(t,e){let n=new ao(this.currentPoint.clone(),new rt(t,e));return this.curves.push(n),this.currentPoint.set(t,e),this}quadraticCurveTo(t,e,n,i){let s=new oo(this.currentPoint.clone(),new rt(t,e),new rt(n,i));return this.curves.push(s),this.currentPoint.set(n,i),this}bezierCurveTo(t,e,n,i,s,r){let a=new ro(this.currentPoint.clone(),new rt(t,e),new rt(n,i),new rt(s,r));return this.curves.push(a),this.currentPoint.set(s,r),this}splineThru(t){let e=[this.currentPoint.clone()].concat(t),n=new co(e);return this.curves.push(n),this.currentPoint.copy(t[t.length-1]),this}arc(t,e,n,i,s,r){let a=this.currentPoint.x,o=this.currentPoint.y;return this.absarc(t+a,e+o,n,i,s,r),this}absarc(t,e,n,i,s,r){return this.absellipse(t,e,n,n,i,s,r),this}ellipse(t,e,n,i,s,r,a,o){let l=this.currentPoint.x,c=this.currentPoint.y;return this.absellipse(t+l,e+c,n,i,s,r,a,o),this}absellipse(t,e,n,i,s,r,a,o){let l=new xr(t,e,n,i,s,r,a,o);if(this.curves.length>0){let h=l.getPoint(0);if(!h.equals(this.currentPoint))this.lineTo(h.x,h.y)}this.curves.push(l);let c=l.getPoint(1);return this.currentPoint.copy(c),this}copy(t){return super.copy(t),this.currentPoint.copy(t.currentPoint),this}toJSON(){let t=super.toJSON();return t.currentPoint=this.currentPoint.toArray(),t}fromJSON(t){return super.fromJSON(t),this.currentPoint.fromArray(t.currentPoint),this}}class Cs extends Jn{constructor(t){super(t);this.uuid=rn(),this.type="Shape",this.holes=[]}getPointsHoles(t){let e=[];for(let n=0,i=this.holes.length;n<i;n++)e[n]=this.holes[n].getPoints(t);return e}extractPoints(t){return{shape:this.getPoints(t),holes:this.getPointsHoles(t)}}copy(t){super.copy(t),this.holes=[];for(let e=0,n=t.holes.length;e<n;e++){let i=t.holes[e];this.holes.push(i.clone())}return this}toJSON(){let t=super.toJSON();t.uuid=this.uuid,t.holes=[];for(let e=0,n=this.holes.length;e<n;e++){let i=this.holes[e];t.holes.push(i.toJSON())}return t}fromJSON(t){super.fromJSON(t),this.uuid=t.uuid,this.holes=[];for(let e=0,n=t.holes.length;e<n;e++){let i=t.holes[e];this.holes.push(new Jn().fromJSON(i))}return this}}function w0(t,e,n=2){let i=e&&e.length,s=i?e[0]*n:t.length,r=bf(t,0,s,n,!0),a=[];if(!r||r.next===r.prev)return a;let o,l,c;if(i)r=L0(t,e,r,n);if(t.length>80*n){o=t[0],l=t[1];let h=o,d=l;for(let u=n;u<s;u+=n){let f=t[u],g=t[u+1];if(f<o)o=f;if(g<l)l=g;if(f>h)h=f;if(g>d)d=g}c=Math.max(h-o,d-l),c=c!==0?32767/c:0}return er(r,a,n,o,l,c,0),a}function bf(t,e,n,i,s){let r;if(s===V0(t,e,n,i)>0)for(let a=e;a<n;a+=i)r=wu(a/i|0,t[a],t[a+1],r);else for(let a=n-i;a>=e;a-=i)r=wu(a/i|0,t[a],t[a+1],r);if(r&&xs(r,r.next))ir(r),r=r.next;return r}function Ui(t,e){if(!t)return t;if(!e)e=t;let n=t,i;do if(i=!1,!n.steiner&&(xs(n,n.next)||Me(n.prev,n,n.next)===0)){if(ir(n),n=e=n.prev,n===n.next)break;i=!0}else n=n.next;while(i||n!==e);return e}function er(t,e,n,i,s,r,a){if(!t)return;if(!a&&r)O0(t,i,s,r);let o=t;while(t.prev!==t.next){let l=t.prev,c=t.next;if(r?C0(t,i,s,r):R0(t)){e.push(l.i,t.i,c.i),ir(t),t=c.next,o=c.next;continue}if(t=c,t===o){if(!a)er(Ui(t),e,n,i,s,r,1);else if(a===1)t=I0(Ui(t),e),er(t,e,n,i,s,r,2);else if(a===2)P0(t,e,n,i,s,r);break}}}function R0(t){let e=t.prev,n=t,i=t.next;if(Me(e,n,i)>=0)return!1;let s=e.x,r=n.x,a=i.x,o=e.y,l=n.y,c=i.y,h=Math.min(s,r,a),d=Math.min(o,l,c),u=Math.max(s,r,a),f=Math.max(o,l,c),g=i.next;while(g!==e){if(g.x>=h&&g.x<=u&&g.y>=d&&g.y<=f&&Js(s,o,r,l,a,c,g.x,g.y)&&Me(g.prev,g,g.next)>=0)return!1;g=g.next}return!0}function C0(t,e,n,i){let s=t.prev,r=t,a=t.next;if(Me(s,r,a)>=0)return!1;let o=s.x,l=r.x,c=a.x,h=s.y,d=r.y,u=a.y,f=Math.min(o,l,c),g=Math.min(h,d,u),_=Math.max(o,l,c),m=Math.max(h,d,u),p=Ml(f,g,e,n,i),b=Ml(_,m,e,n,i),{prevZ:S,nextZ:x}=t;while(S&&S.z>=p&&x&&x.z<=b){if(S.x>=f&&S.x<=_&&S.y>=g&&S.y<=m&&S!==s&&S!==a&&Js(o,h,l,d,c,u,S.x,S.y)&&Me(S.prev,S,S.next)>=0)return!1;if(S=S.prevZ,x.x>=f&&x.x<=_&&x.y>=g&&x.y<=m&&x!==s&&x!==a&&Js(o,h,l,d,c,u,x.x,x.y)&&Me(x.prev,x,x.next)>=0)return!1;x=x.nextZ}while(S&&S.z>=p){if(S.x>=f&&S.x<=_&&S.y>=g&&S.y<=m&&S!==s&&S!==a&&Js(o,h,l,d,c,u,S.x,S.y)&&Me(S.prev,S,S.next)>=0)return!1;S=S.prevZ}while(x&&x.z<=b){if(x.x>=f&&x.x<=_&&x.y>=g&&x.y<=m&&x!==s&&x!==a&&Js(o,h,l,d,c,u,x.x,x.y)&&Me(x.prev,x,x.next)>=0)return!1;x=x.nextZ}return!0}function I0(t,e){let n=t;do{let i=n.prev,s=n.next.next;if(!xs(i,s)&&Ef(i,n,n.next,s)&&nr(i,s)&&nr(s,i))e.push(i.i,n.i,s.i),ir(n),ir(n.next),n=t=s;n=n.next}while(n!==t);return Ui(n)}function P0(t,e,n,i,s,r){let a=t;do{let o=a.next.next;while(o!==a.prev){if(a.i!==o.i&&k0(a,o)){let l=Af(a,o);a=Ui(a,a.next),l=Ui(l,l.next),er(a,e,n,i,s,r,0),er(l,e,n,i,s,r,0);return}o=o.next}a=a.next}while(a!==t)}function L0(t,e,n,i){let s=[];for(let r=0,a=e.length;r<a;r++){let o=e[r]*i,l=r<a-1?e[r+1]*i:t.length,c=bf(t,o,l,i,!1);if(c===c.next)c.steiner=!0;s.push(z0(c))}s.sort(N0);for(let r=0;r<s.length;r++)n=D0(s[r],n);return n}function N0(t,e){let n=t.x-e.x;if(n===0){if(n=t.y-e.y,n===0){let i=(t.next.y-t.y)/(t.next.x-t.x),s=(e.next.y-e.y)/(e.next.x-e.x);n=i-s}}return n}function D0(t,e){let n=U0(t,e);if(!n)return e;let i=Af(n,t);return Ui(i,i.next),Ui(n,n.next)}function U0(t,e){let n=e,{x:i,y:s}=t,r=-1/0,a;if(xs(t,n))return n;do{if(xs(t,n.next))return n.next;else if(s<=n.y&&s>=n.next.y&&n.next.y!==n.y){let d=n.x+(s-n.y)*(n.next.x-n.x)/(n.next.y-n.y);if(d<=i&&d>r){if(r=d,a=n.x<n.next.x?n:n.next,d===i)return a}}n=n.next}while(n!==e);if(!a)return null;let o=a,l=a.x,c=a.y,h=1/0;n=a;do{if(i>=n.x&&n.x>=l&&i!==n.x&&Tf(s<c?i:r,s,l,c,s<c?r:i,s,n.x,n.y)){let d=Math.abs(s-n.y)/(i-n.x);if(nr(n,t)&&(d<h||d===h&&(n.x>a.x||n.x===a.x&&F0(a,n))))a=n,h=d}n=n.next}while(n!==o);return a}function F0(t,e){return Me(t.prev,t,e.prev)<0&&Me(e.next,t,t.next)<0}function O0(t,e,n,i){let s=t;do{if(s.z===0)s.z=Ml(s.x,s.y,e,n,i);s.prevZ=s.prev,s.nextZ=s.next,s=s.next}while(s!==t);s.prevZ.nextZ=null,s.prevZ=null,B0(s)}function B0(t){let e,n=1;do{let i=t,s;t=null;let r=null;e=0;while(i){e++;let a=i,o=0;for(let c=0;c<n;c++)if(o++,a=a.nextZ,!a)break;let l=n;while(o>0||l>0&&a){if(o!==0&&(l===0||!a||i.z<=a.z))s=i,i=i.nextZ,o--;else s=a,a=a.nextZ,l--;if(r)r.nextZ=s;else t=s;s.prevZ=r,r=s}i=a}r.nextZ=null,n*=2}while(e>1);return t}function Ml(t,e,n,i,s){return t=(t-n)*s|0,e=(e-i)*s|0,t=(t|t<<8)&16711935,t=(t|t<<4)&252645135,t=(t|t<<2)&858993459,t=(t|t<<1)&1431655765,e=(e|e<<8)&16711935,e=(e|e<<4)&252645135,e=(e|e<<2)&858993459,e=(e|e<<1)&1431655765,t|e<<1}function z0(t){let e=t,n=t;do{if(e.x<n.x||e.x===n.x&&e.y<n.y)n=e;e=e.next}while(e!==t);return n}function Tf(t,e,n,i,s,r,a,o){return(s-a)*(e-o)>=(t-a)*(r-o)&&(t-a)*(i-o)>=(n-a)*(e-o)&&(n-a)*(r-o)>=(s-a)*(i-o)}function Js(t,e,n,i,s,r,a,o){return!(t===a&&e===o)&&Tf(t,e,n,i,s,r,a,o)}function k0(t,e){return t.next.i!==e.i&&t.prev.i!==e.i&&!G0(t,e)&&(nr(t,e)&&nr(e,t)&&H0(t,e)&&(Me(t.prev,t,e.prev)||Me(t,e.prev,e))||xs(t,e)&&Me(t.prev,t,t.next)>0&&Me(e.prev,e,e.next)>0)}function Me(t,e,n){return(e.y-t.y)*(n.x-e.x)-(e.x-t.x)*(n.y-e.y)}function xs(t,e){return t.x===e.x&&t.y===e.y}function Ef(t,e,n,i){let s=ga(Me(t,e,n)),r=ga(Me(t,e,i)),a=ga(Me(n,i,t)),o=ga(Me(n,i,e));if(s!==r&&a!==o)return!0;if(s===0&&ma(t,n,e))return!0;if(r===0&&ma(t,i,e))return!0;if(a===0&&ma(n,t,i))return!0;if(o===0&&ma(n,e,i))return!0;return!1}function ma(t,e,n){return e.x<=Math.max(t.x,n.x)&&e.x>=Math.min(t.x,n.x)&&e.y<=Math.max(t.y,n.y)&&e.y>=Math.min(t.y,n.y)}function ga(t){return t>0?1:t<0?-1:0}function G0(t,e){let n=t;do{if(n.i!==t.i&&n.next.i!==t.i&&n.i!==e.i&&n.next.i!==e.i&&Ef(n,n.next,t,e))return!0;n=n.next}while(n!==t);return!1}function nr(t,e){return Me(t.prev,t,t.next)<0?Me(t,e,t.next)>=0&&Me(t,t.prev,e)>=0:Me(t,e,t.prev)<0||Me(t,t.next,e)<0}function H0(t,e){let n=t,i=!1,s=(t.x+e.x)/2,r=(t.y+e.y)/2;do{if(n.y>r!==n.next.y>r&&n.next.y!==n.y&&s<(n.next.x-n.x)*(r-n.y)/(n.next.y-n.y)+n.x)i=!i;n=n.next}while(n!==t);return i}function Af(t,e){let n=bl(t.i,t.x,t.y),i=bl(e.i,e.x,e.y),s=t.next,r=e.prev;return t.next=e,e.prev=t,n.next=s,s.prev=n,i.next=n,n.prev=i,r.next=i,i.prev=r,i}function wu(t,e,n,i){let s=bl(t,e,n);if(!i)s.prev=s,s.next=s;else s.next=i.next,s.prev=i,i.next.prev=s,i.next=s;return s}function ir(t){if(t.next.prev=t.prev,t.prev.next=t.next,t.prevZ)t.prevZ.nextZ=t.nextZ;if(t.nextZ)t.nextZ.prevZ=t.prevZ}function bl(t,e,n){return{i:t,x:e,y:n,prev:null,next:null,z:0,prevZ:null,nextZ:null,steiner:!1}}function V0(t,e,n,i){let s=0;for(let r=e,a=n-i;r<n;r+=i)s+=(t[a]-t[r])*(t[r+1]+t[a+1]),a=r;return s}class wf{static triangulate(t,e,n=2){return w0(t,e,n)}}class sn{static area(t){let e=t.length,n=0;for(let i=e-1,s=0;s<e;i=s++)n+=t[i].x*t[s].y-t[s].x*t[i].y;return n*0.5}static isClockWise(t){return sn.area(t)<0}static triangulateShape(t,e){let n=[],i=[],s=[];Ru(t),Cu(n,t);let r=t.length;e.forEach(Ru);for(let o=0;o<e.length;o++)i.push(r),r+=e[o].length,Cu(n,e[o]);let a=wf.triangulate(n,i);for(let o=0;o<a.length;o+=3)s.push(a.slice(o,o+3));return s}}function Ru(t){let e=t.length;if(e>2&&t[e-1].equals(t[0]))t.pop()}function Cu(t,e){for(let n=0;n<e.length;n++)t.push(e[n].x),t.push(e[n].y)}class ho extends jt{constructor(t=new Cs([new rt(0.5,0.5),new rt(-0.5,0.5),new rt(-0.5,-0.5),new rt(0.5,-0.5)]),e={}){super();this.type="ExtrudeGeometry",this.parameters={shapes:t,options:e},t=Array.isArray(t)?t:[t];let n=this,i=[],s=[];for(let a=0,o=t.length;a<o;a++){let l=t[a];r(l)}this.setAttribute("position",new Nt(i,3)),this.setAttribute("uv",new Nt(s,2)),this.computeVertexNormals();function r(a){let o=[],l=e.curveSegments!==void 0?e.curveSegments:12,c=e.steps!==void 0?e.steps:1,h=e.depth!==void 0?e.depth:1,d=e.bevelEnabled!==void 0?e.bevelEnabled:!0,u=e.bevelThickness!==void 0?e.bevelThickness:0.2,f=e.bevelSize!==void 0?e.bevelSize:u-0.1,g=e.bevelOffset!==void 0?e.bevelOffset:0,_=e.bevelSegments!==void 0?e.bevelSegments:3,m=e.extrudePath,p=e.UVGenerator!==void 0?e.UVGenerator:W0,b,S=!1,x,T,I,R;if(m){b=m.getSpacedPoints(c),S=!0,d=!1;let C=m.isCatmullRomCurve3?m.closed:!1;x=m.computeFrenetFrames(c,C),T=new L,I=new L,R=new L}if(!d)_=0,u=0,f=0,g=0;let y=a.extractPoints(l),{shape:A,holes:Z}=y;if(!sn.isClockWise(A)){A=A.reverse();for(let C=0,H=Z.length;C<H;C++){let U=Z[C];if(sn.isClockWise(U))Z[C]=U.reverse()}}function q(C){let Y=C[0];for(let tt=1;tt<=C.length;tt++){let at=tt%C.length,mt=C[at],gt=mt.x-Y.x,X=mt.y-Y.y,v=gt*gt+X*X,Lt=Math.max(Math.abs(mt.x),Math.abs(mt.y),Math.abs(Y.x),Math.abs(Y.y)),zt=0.000000000000000000010000000000000001*Lt*Lt;if(v<=zt){C.splice(at,1),tt--;continue}Y=mt}}q(A),Z.forEach(q);let ot=Z.length,z=A;for(let C=0;C<ot;C++){let H=Z[C];A=A.concat(H)}function et(C,H,U){if(!H)Wt("ExtrudeGeometry: vec does not exist");return C.clone().addScaledVector(H,U)}let nt=A.length;function B(C,H,U){let Y,tt,at,mt=C.x-H.x,gt=C.y-H.y,X=U.x-C.x,v=U.y-C.y,Lt=mt*mt+gt*gt,zt=mt*v-gt*X;if(Math.abs(zt)>Number.EPSILON){let $t=Math.sqrt(Lt),w=Math.sqrt(X*X+v*v),M=H.x-gt/$t,k=H.y+mt/$t,it=U.x-v/w,pt=U.y+X/w,yt=((it-M)*v-(pt-k)*X)/(mt*v-gt*X);Y=M+mt*yt-C.x,tt=k+gt*yt-C.y;let St=Y*Y+tt*tt;if(St<=2)return new rt(Y,tt);else at=Math.sqrt(St/2)}else{let $t=!1;if(mt>Number.EPSILON){if(X>Number.EPSILON)$t=!0}else if(mt<-Number.EPSILON){if(X<-Number.EPSILON)$t=!0}else if(Math.sign(gt)===Math.sign(v))$t=!0;if($t)Y=-gt,tt=mt,at=Math.sqrt(Lt);else Y=mt,tt=gt,at=Math.sqrt(Lt/2)}return new rt(Y/at,tt/at)}let ft=[];for(let C=0,H=z.length,U=H-1,Y=C+1;C<H;C++,U++,Y++){if(U===H)U=0;if(Y===H)Y=0;ft[C]=B(z[C],z[U],z[Y])}let Q=[],ct,dt=ft.concat();for(let C=0,H=ot;C<H;C++){let U=Z[C];ct=[];for(let Y=0,tt=U.length,at=tt-1,mt=Y+1;Y<tt;Y++,at++,mt++){if(at===tt)at=0;if(mt===tt)mt=0;ct[Y]=B(U[Y],U[at],U[mt])}Q.push(ct),dt=dt.concat(ct)}let Bt;if(_===0)Bt=sn.triangulateShape(z,Z);else{let C=[],H=[];for(let U=0;U<_;U++){let Y=U/_,tt=u*Math.cos(Y*Math.PI/2),at=f*Math.sin(Y*Math.PI/2)+g;for(let mt=0,gt=z.length;mt<gt;mt++){let X=et(z[mt],ft[mt],at);if(N(X.x,X.y,-tt),Y===0)C.push(X)}for(let mt=0,gt=ot;mt<gt;mt++){let X=Z[mt];ct=Q[mt];let v=[];for(let Lt=0,zt=X.length;Lt<zt;Lt++){let $t=et(X[Lt],ct[Lt],at);if(N($t.x,$t.y,-tt),Y===0)v.push($t)}if(Y===0)H.push(v)}}Bt=sn.triangulateShape(C,H)}let Dt=Bt.length,ee=f+g;for(let C=0;C<nt;C++){let H=d?et(A[C],dt[C],ee):A[C];if(!S)N(H.x,H.y,0);else I.copy(x.normals[0]).multiplyScalar(H.x),T.copy(x.binormals[0]).multiplyScalar(H.y),R.copy(b[0]).add(I).add(T),N(R.x,R.y,R.z)}for(let C=1;C<=c;C++)for(let H=0;H<nt;H++){let U=d?et(A[H],dt[H],ee):A[H];if(!S)N(U.x,U.y,h/c*C);else I.copy(x.normals[C]).multiplyScalar(U.x),T.copy(x.binormals[C]).multiplyScalar(U.y),R.copy(b[C]).add(I).add(T),N(R.x,R.y,R.z)}for(let C=_-1;C>=0;C--){let H=C/_,U=u*Math.cos(H*Math.PI/2),Y=f*Math.sin(H*Math.PI/2)+g;for(let tt=0,at=z.length;tt<at;tt++){let mt=et(z[tt],ft[tt],Y);N(mt.x,mt.y,h+U)}for(let tt=0,at=Z.length;tt<at;tt++){let mt=Z[tt];ct=Q[tt];for(let gt=0,X=mt.length;gt<X;gt++){let v=et(mt[gt],ct[gt],Y);if(!S)N(v.x,v.y,h+U);else N(v.x,v.y+b[c-1].y,b[c-1].x+U)}}}Jt(),F();function Jt(){let C=i.length/3;if(d){let H=0,U=nt*H;for(let Y=0;Y<Dt;Y++){let tt=Bt[Y];P(tt[2]+U,tt[1]+U,tt[0]+U)}H=c+_*2,U=nt*H;for(let Y=0;Y<Dt;Y++){let tt=Bt[Y];P(tt[0]+U,tt[1]+U,tt[2]+U)}}else{for(let H=0;H<Dt;H++){let U=Bt[H];P(U[2],U[1],U[0])}for(let H=0;H<Dt;H++){let U=Bt[H];P(U[0]+nt*c,U[1]+nt*c,U[2]+nt*c)}}n.addGroup(C,i.length/3-C,0)}function F(){let C=i.length/3,H=0;O(z,H),H+=z.length;for(let U=0,Y=Z.length;U<Y;U++){let tt=Z[U];O(tt,H),H+=tt.length}n.addGroup(C,i.length/3-C,1)}function O(C,H){let U=C.length;while(--U>=0){let Y=U,tt=U-1;if(tt<0)tt=C.length-1;for(let at=0,mt=c+_*2;at<mt;at++){let gt=nt*at,X=nt*(at+1),v=H+Y+gt,Lt=H+tt+gt,zt=H+tt+X,$t=H+Y+X;G(v,Lt,zt,$t)}}}function N(C,H,U){o.push(C),o.push(H),o.push(U)}function P(C,H,U){J(C),J(H),J(U);let Y=i.length/3,tt=p.generateTopUV(n,i,Y-3,Y-2,Y-1);_t(tt[0]),_t(tt[1]),_t(tt[2])}function G(C,H,U,Y){J(C),J(H),J(Y),J(H),J(U),J(Y);let tt=i.length/3,at=p.generateSideWallUV(n,i,tt-6,tt-3,tt-2,tt-1);_t(at[0]),_t(at[1]),_t(at[3]),_t(at[1]),_t(at[2]),_t(at[3])}function J(C){i.push(o[C*3+0]),i.push(o[C*3+1]),i.push(o[C*3+2])}function _t(C){s.push(C.x),s.push(C.y)}}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}toJSON(){let t=super.toJSON(),e=this.parameters.shapes,n=this.parameters.options;return X0(e,n,t)}static fromJSON(t,e){let n=[];for(let s=0,r=t.shapes.length;s<r;s++){let a=e[t.shapes[s]];n.push(a)}let i=t.options.extrudePath;if(i!==void 0)t.options.extrudePath=new Ra[i.type]().fromJSON(i);return new ho(n,t.options)}}var W0={generateTopUV:function(t,e,n,i,s){let r=e[n*3],a=e[n*3+1],o=e[i*3],l=e[i*3+1],c=e[s*3],h=e[s*3+1];return[new rt(r,a),new rt(o,l),new rt(c,h)]},generateSideWallUV:function(t,e,n,i,s,r){let a=e[n*3],o=e[n*3+1],l=e[n*3+2],c=e[i*3],h=e[i*3+1],d=e[i*3+2],u=e[s*3],f=e[s*3+1],g=e[s*3+2],_=e[r*3],m=e[r*3+1],p=e[r*3+2];if(Math.abs(o-h)<Math.abs(a-c))return[new rt(a,1-l),new rt(c,1-d),new rt(u,1-g),new rt(_,1-p)];else return[new rt(o,1-l),new rt(h,1-d),new rt(f,1-g),new rt(m,1-p)]}};function X0(t,e,n){if(n.shapes=[],Array.isArray(t))for(let i=0,s=t.length;i<s;i++){let r=t[i];n.shapes.push(r.uuid)}else n.shapes.push(t.uuid);if(n.options=Object.assign({},e),e.extrudePath!==void 0)n.options.extrudePath=e.extrudePath.toJSON();return n}class uo extends gi{constructor(t=1,e=0){let n=(1+Math.sqrt(5))/2,i=[-1,n,0,1,n,0,-1,-n,0,1,-n,0,0,-1,n,0,1,n,0,-1,-n,0,1,-n,n,0,-1,n,0,1,-n,0,-1,-n,0,1],s=[0,11,5,0,5,1,0,1,7,0,7,10,0,10,11,1,5,9,5,11,4,11,10,2,10,7,6,7,1,8,3,9,4,3,4,2,3,2,6,3,6,8,3,8,9,4,9,5,2,4,11,6,2,10,8,6,7,9,8,1];super(i,s,t,e);this.type="IcosahedronGeometry",this.parameters={radius:t,detail:e}}static fromJSON(t){return new uo(t.radius,t.detail)}}class fo extends jt{constructor(t=[new rt(0,-0.5),new rt(0.5,0),new rt(0,0.5)],e=12,n=0,i=Math.PI*2){super();this.type="LatheGeometry",this.parameters={points:t,segments:e,phiStart:n,phiLength:i},e=Math.floor(e),i=Zt(i,0,Math.PI*2);let s=[],r=[],a=[],o=[],l=[],c=1/e,h=new L,d=new rt,u=new L,f=new L,g=new L,_=0,m=0;for(let p=0;p<=t.length-1;p++)switch(p){case 0:_=t[p+1].x-t[p].x,m=t[p+1].y-t[p].y,u.x=m*1,u.y=-_,u.z=m*0,g.copy(u),u.normalize(),o.push(u.x,u.y,u.z);break;case t.length-1:o.push(g.x,g.y,g.z);break;default:_=t[p+1].x-t[p].x,m=t[p+1].y-t[p].y,u.x=m*1,u.y=-_,u.z=m*0,f.copy(u),u.x+=g.x,u.y+=g.y,u.z+=g.z,u.normalize(),o.push(u.x,u.y,u.z),g.copy(f)}for(let p=0;p<=e;p++){let b=n+p*c*i,S=Math.sin(b),x=Math.cos(b);for(let T=0;T<=t.length-1;T++){h.x=t[T].x*S,h.y=t[T].y,h.z=t[T].x*x,r.push(h.x,h.y,h.z),d.x=p/e,d.y=T/(t.length-1),a.push(d.x,d.y);let I=o[3*T+0]*S,R=o[3*T+1],y=o[3*T+0]*x;l.push(I,R,y)}}for(let p=0;p<e;p++)for(let b=0;b<t.length-1;b++){let S=b+p*t.length,x=S,T=S+t.length,I=S+t.length+1,R=S+1;s.push(x,T,R),s.push(I,R,T)}this.setIndex(s),this.setAttribute("position",new Nt(r,3)),this.setAttribute("uv",new Nt(a,2)),this.setAttribute("normal",new Nt(l,3))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new fo(t.points,t.segments,t.phiStart,t.phiLength)}}class yr extends gi{constructor(t=1,e=0){let n=[1,0,0,-1,0,0,0,1,0,0,-1,0,0,0,1,0,0,-1],i=[0,2,4,0,4,3,0,3,5,0,5,2,1,2,5,1,5,3,1,3,4,1,4,2];super(n,i,t,e);this.type="OctahedronGeometry",this.parameters={radius:t,detail:e}}static fromJSON(t){return new yr(t.radius,t.detail)}}class Is extends jt{constructor(t=1,e=1,n=1,i=1){super();this.type="PlaneGeometry",this.parameters={width:t,height:e,widthSegments:n,heightSegments:i};let s=t/2,r=e/2,a=Math.floor(n),o=Math.floor(i),l=a+1,c=o+1,h=t/a,d=e/o,u=[],f=[],g=[],_=[];for(let m=0;m<c;m++){let p=m*d-r;for(let b=0;b<l;b++){let S=b*h-s;f.push(S,-p,0),g.push(0,0,1),_.push(b/a),_.push(1-m/o)}}for(let m=0;m<o;m++)for(let p=0;p<a;p++){let b=p+l*m,S=p+l*(m+1),x=p+1+l*(m+1),T=p+1+l*m;u.push(b,S,T),u.push(S,x,T)}this.setIndex(u),this.setAttribute("position",new Nt(f,3)),this.setAttribute("normal",new Nt(g,3)),this.setAttribute("uv",new Nt(_,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new Is(t.width,t.height,t.widthSegments,t.heightSegments)}}class po extends jt{constructor(t=0.5,e=1,n=32,i=1,s=0,r=Math.PI*2){super();this.type="RingGeometry",this.parameters={innerRadius:t,outerRadius:e,thetaSegments:n,phiSegments:i,thetaStart:s,thetaLength:r},n=Math.max(3,n),i=Math.max(1,i);let a=[],o=[],l=[],c=[],h=t,d=(e-t)/i,u=new L,f=new rt;for(let g=0;g<=i;g++){for(let _=0;_<=n;_++){let m=s+_/n*r;u.x=h*Math.cos(m),u.y=h*Math.sin(m),o.push(u.x,u.y,u.z),l.push(0,0,1),f.x=(u.x/e+1)/2,f.y=(u.y/e+1)/2,c.push(f.x,f.y)}h+=d}for(let g=0;g<i;g++){let _=g*(n+1);for(let m=0;m<n;m++){let p=m+_,b=p,S=p+n+1,x=p+n+2,T=p+1;a.push(b,S,T),a.push(S,x,T)}}this.setIndex(a),this.setAttribute("position",new Nt(o,3)),this.setAttribute("normal",new Nt(l,3)),this.setAttribute("uv",new Nt(c,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new po(t.innerRadius,t.outerRadius,t.thetaSegments,t.phiSegments,t.thetaStart,t.thetaLength)}}class mo extends jt{constructor(t=new Cs([new rt(0,0.5),new rt(-0.5,-0.5),new rt(0.5,-0.5)]),e=12){super();this.type="ShapeGeometry",this.parameters={shapes:t,curveSegments:e};let n=[],i=[],s=[],r=[],a=0,o=0;if(Array.isArray(t)===!1)l(t);else for(let c=0;c<t.length;c++)l(t[c]),this.addGroup(a,o,c),a+=o,o=0;this.setIndex(n),this.setAttribute("position",new Nt(i,3)),this.setAttribute("normal",new Nt(s,3)),this.setAttribute("uv",new Nt(r,2));function l(c){let h=i.length/3,d=c.extractPoints(e),{shape:u,holes:f}=d;if(sn.isClockWise(u)===!1)u=u.reverse();for(let _=0,m=f.length;_<m;_++){let p=f[_];if(sn.isClockWise(p)===!0)f[_]=p.reverse()}let g=sn.triangulateShape(u,f);for(let _=0,m=f.length;_<m;_++){let p=f[_];u=u.concat(p)}for(let _=0,m=u.length;_<m;_++){let p=u[_];i.push(p.x,p.y,0),s.push(0,0,1),r.push(p.x,p.y)}for(let _=0,m=g.length;_<m;_++){let p=g[_],b=p[0]+h,S=p[1]+h,x=p[2]+h;n.push(b,S,x),o+=3}}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}toJSON(){let t=super.toJSON(),e=this.parameters.shapes;return q0(e,t)}static fromJSON(t,e){let n=[];for(let i=0,s=t.shapes.length;i<s;i++){let r=e[t.shapes[i]];n.push(r)}return new mo(n,t.curveSegments)}}function q0(t,e){if(e.shapes=[],Array.isArray(t))for(let n=0,i=t.length;n<i;n++){let s=t[n];e.shapes.push(s.uuid)}else e.shapes.push(t.uuid);return e}class vr extends jt{constructor(t=1,e=32,n=16,i=0,s=Math.PI*2,r=0,a=Math.PI){super();this.type="SphereGeometry",this.parameters={radius:t,widthSegments:e,heightSegments:n,phiStart:i,phiLength:s,thetaStart:r,thetaLength:a},e=Math.max(3,Math.floor(e)),n=Math.max(2,Math.floor(n));let o=Math.min(r+a,Math.PI),l=0,c=[],h=new L,d=new L,u=[],f=[],g=[],_=[];for(let m=0;m<=n;m++){let p=[],b=m/n,S=r+b*a,x=t*Math.cos(S),T=Math.sqrt(t*t-x*x),I=0;if(m===0&&r===0)I=0.5/e;else if(m===n&&o===Math.PI)I=-0.5/e;for(let R=0;R<=e;R++){let y=R/e,A=i+y*s;h.x=-T*Math.cos(A),h.y=x,h.z=T*Math.sin(A),f.push(h.x,h.y,h.z),d.copy(h).normalize(),g.push(d.x,d.y,d.z),_.push(y+I,1-b),p.push(l++)}c.push(p)}for(let m=0;m<n;m++)for(let p=0;p<e;p++){let b=c[m][p+1],S=c[m][p],x=c[m+1][p],T=c[m+1][p+1];if(m!==0||r>0)u.push(b,S,T);if(m!==n-1||o<Math.PI)u.push(S,x,T)}this.setIndex(u),this.setAttribute("position",new Nt(f,3)),this.setAttribute("normal",new Nt(g,3)),this.setAttribute("uv",new Nt(_,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new vr(t.radius,t.widthSegments,t.heightSegments,t.phiStart,t.phiLength,t.thetaStart,t.thetaLength)}}class go extends gi{constructor(t=1,e=0){let n=[1,1,1,-1,-1,1,-1,1,-1,1,-1,-1],i=[2,1,0,0,3,2,1,3,0,2,3,1];super(n,i,t,e);this.type="TetrahedronGeometry",this.parameters={radius:t,detail:e}}static fromJSON(t){return new go(t.radius,t.detail)}}class _o extends jt{constructor(t=1,e=0.4,n=12,i=48,s=Math.PI*2,r=0,a=Math.PI*2){super();this.type="TorusGeometry",this.parameters={radius:t,tube:e,radialSegments:n,tubularSegments:i,arc:s,thetaStart:r,thetaLength:a},n=Math.floor(n),i=Math.floor(i);let o=[],l=[],c=[],h=[],d=new L,u=new L,f=new L;for(let g=0;g<=n;g++){let _=r+g/n*a;for(let m=0;m<=i;m++){let p=m/i*s;u.x=(t+e*Math.cos(_))*Math.cos(p),u.y=(t+e*Math.cos(_))*Math.sin(p),u.z=e*Math.sin(_),l.push(u.x,u.y,u.z),d.x=t*Math.cos(p),d.y=t*Math.sin(p),f.subVectors(u,d).normalize(),c.push(f.x,f.y,f.z),h.push(m/i),h.push(g/n)}}for(let g=1;g<=n;g++)for(let _=1;_<=i;_++){let m=(i+1)*g+_-1,p=(i+1)*(g-1)+_-1,b=(i+1)*(g-1)+_,S=(i+1)*g+_;o.push(m,p,S),o.push(p,b,S)}this.setIndex(o),this.setAttribute("position",new Nt(l,3)),this.setAttribute("normal",new Nt(c,3)),this.setAttribute("uv",new Nt(h,2))}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new _o(t.radius,t.tube,t.radialSegments,t.tubularSegments,t.arc,t.thetaStart,t.thetaLength)}}class xo extends jt{constructor(t=1,e=0.4,n=64,i=8,s=2,r=3){super();this.type="TorusKnotGeometry",this.parameters={radius:t,tube:e,tubularSegments:n,radialSegments:i,p:s,q:r},n=Math.floor(n),i=Math.floor(i);let a=[],o=[],l=[],c=[],h=new L,d=new L,u=new L,f=new L,g=new L,_=new L,m=new L;for(let b=0;b<=n;++b){let S=b/n*s*Math.PI*2;p(S,s,r,t,u),p(S+0.01,s,r,t,f),_.subVectors(f,u),m.addVectors(f,u),g.crossVectors(_,m),m.crossVectors(g,_),g.normalize(),m.normalize();for(let x=0;x<=i;++x){let T=x/i*Math.PI*2,I=-e*Math.cos(T),R=e*Math.sin(T);h.x=u.x+(I*m.x+R*g.x),h.y=u.y+(I*m.y+R*g.y),h.z=u.z+(I*m.z+R*g.z),o.push(h.x,h.y,h.z),d.subVectors(h,u).normalize(),l.push(d.x,d.y,d.z),c.push(b/n),c.push(x/i)}}for(let b=1;b<=n;b++)for(let S=1;S<=i;S++){let x=(i+1)*(b-1)+(S-1),T=(i+1)*b+(S-1),I=(i+1)*b+S,R=(i+1)*(b-1)+S;a.push(x,T,R),a.push(T,I,R)}this.setIndex(a),this.setAttribute("position",new Nt(o,3)),this.setAttribute("normal",new Nt(l,3)),this.setAttribute("uv",new Nt(c,2));function p(b,S,x,T,I){let R=Math.cos(b),y=Math.sin(b),A=x/S*b,Z=Math.cos(A);I.x=T*(2+Z)*0.5*R,I.y=T*(2+Z)*y*0.5,I.z=T*Math.sin(A)*0.5}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}static fromJSON(t){return new xo(t.radius,t.tube,t.tubularSegments,t.radialSegments,t.p,t.q)}}class yo extends jt{constructor(t=new lo(new L(-1,-1,0),new L(-1,1,0),new L(1,1,0)),e=64,n=1,i=8,s=!1){super();this.type="TubeGeometry",this.parameters={path:t,tubularSegments:e,radius:n,radialSegments:i,closed:s};let r=t.computeFrenetFrames(e,s);this.tangents=r.tangents,this.normals=r.normals,this.binormals=r.binormals;let a=new L,o=new L,l=new rt,c=new L,h=[],d=[],u=[],f=[];g(),this.setIndex(f),this.setAttribute("position",new Nt(h,3)),this.setAttribute("normal",new Nt(d,3)),this.setAttribute("uv",new Nt(u,2));function g(){for(let b=0;b<e;b++)_(b);_(s===!1?e:0),p(),m()}function _(b){c=t.getPointAt(b/e,c);let S=r.normals[b],x=r.binormals[b];for(let T=0;T<=i;T++){let I=T/i*Math.PI*2,R=Math.sin(I),y=-Math.cos(I);o.x=y*S.x+R*x.x,o.y=y*S.y+R*x.y,o.z=y*S.z+R*x.z,o.normalize(),d.push(o.x,o.y,o.z),a.x=c.x+n*o.x,a.y=c.y+n*o.y,a.z=c.z+n*o.z,h.push(a.x,a.y,a.z)}}function m(){for(let b=1;b<=e;b++)for(let S=1;S<=i;S++){let x=(i+1)*(b-1)+(S-1),T=(i+1)*b+(S-1),I=(i+1)*b+S,R=(i+1)*(b-1)+S;f.push(x,T,R),f.push(T,I,R)}}function p(){for(let b=0;b<=e;b++)for(let S=0;S<=i;S++)l.x=b/e,l.y=S/i,u.push(l.x,l.y)}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}toJSON(){let t=super.toJSON();return t.path=this.parameters.path.toJSON(),t}static fromJSON(t){return new yo(new Ra[t.path.type]().fromJSON(t.path),t.tubularSegments,t.radius,t.radialSegments,t.closed)}}class Vc extends jt{constructor(t=null){super();if(this.type="WireframeGeometry",this.parameters={geometry:t},t!==null){let e=[],n=new Set,i=new L,s=new L;if(t.index!==null){let r=t.attributes.position,{index:a,groups:o}=t;if(o.length===0)o=[{start:0,count:a.count,materialIndex:0}];for(let l=0,c=o.length;l<c;++l){let h=o[l],{start:d,count:u}=h;for(let f=d,g=d+u;f<g;f+=3)for(let _=0;_<3;_++){let m=a.getX(f+_),p=a.getX(f+(_+1)%3);if(i.fromBufferAttribute(r,m),s.fromBufferAttribute(r,p),Iu(i,s,n)===!0)e.push(i.x,i.y,i.z),e.push(s.x,s.y,s.z)}}}else{let r=t.attributes.position;for(let a=0,o=r.count/3;a<o;a++)for(let l=0;l<3;l++){let c=3*a+l,h=3*a+(l+1)%3;if(i.fromBufferAttribute(r,c),s.fromBufferAttribute(r,h),Iu(i,s,n)===!0)e.push(i.x,i.y,i.z),e.push(s.x,s.y,s.z)}}this.setAttribute("position",new Nt(e,3))}}copy(t){return super.copy(t),this.parameters=Object.assign({},t.parameters),this}}function Iu(t,e,n){let i=`${t.x},${t.y},${t.z}-${e.x},${e.y},${e.z}`,s=`${e.x},${e.y},${e.z}-${t.x},${t.y},${t.z}`;if(n.has(i)===!0||n.has(s)===!0)return!1;else return n.add(i),n.add(s),!0}var Pu=Object.freeze({__proto__:null,BoxGeometry:An,CapsuleGeometry:no,CircleGeometry:io,ConeGeometry:_r,CylinderGeometry:gr,DodecahedronGeometry:so,EdgesGeometry:Fc,ExtrudeGeometry:ho,IcosahedronGeometry:uo,LatheGeometry:fo,OctahedronGeometry:yr,PlaneGeometry:Is,PolyhedronGeometry:gi,RingGeometry:po,ShapeGeometry:mo,SphereGeometry:vr,TetrahedronGeometry:go,TorusGeometry:_o,TorusKnotGeometry:xo,TubeGeometry:yo,WireframeGeometry:Vc});class Wc extends Fe{constructor(t){super();this.isShadowMaterial=!0,this.type="ShadowMaterial",this.color=new Pt(0),this.transparent=!0,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.fog=t.fog,this}}function Vi(t){let e={};for(let n in t){e[n]={};for(let i in t[n]){let s=t[n][i];if(Lu(s))if(s.isRenderTargetTexture)Mt("UniformsUtils: Textures of render targets cannot be cloned via cloneUniforms() or mergeUniforms()."),e[n][i]=null;else e[n][i]=s.clone();else if(Array.isArray(s))if(Lu(s[0])){let r=[];for(let a=0,o=s.length;a<o;a++)r[a]=s[a].clone();e[n][i]=r}else e[n][i]=s.slice();else e[n][i]=s}}return e}function Xe(t){let e={};for(let n=0;n<t.length;n++){let i=Vi(t[n]);for(let s in i)e[s]=i[s]}return e}function Lu(t){return t&&(t.isColor||t.isMatrix3||t.isMatrix4||t.isVector2||t.isVector3||t.isVector4||t.isTexture||t.isQuaternion)}function Y0(t){let e=[];for(let n=0;n<t.length;n++)e.push(t[n].clone());return e}function Xc(t){let e=t.getRenderTarget();if(e===null)return t.outputColorSpace;if(e.isXRRenderTarget===!0)return e.texture.colorSpace;return re.workingColorSpace}var Rf={clone:Vi,merge:Xe},Z0=`void main() {
	gl_Position = projectionMatrix * modelViewMatrix * vec4( position, 1.0 );
}`,J0=`void main() {
	gl_FragColor = vec4( 1.0, 0.0, 0.0, 1.0 );
}`;class ln extends Fe{constructor(t){super();if(this.isShaderMaterial=!0,this.type="ShaderMaterial",this.defines={},this.uniforms={},this.uniformsGroups=[],this.vertexShader=Z0,this.fragmentShader=J0,this.linewidth=1,this.wireframe=!1,this.wireframeLinewidth=1,this.fog=!1,this.lights=!1,this.clipping=!1,this.forceSinglePass=!0,this.extensions={clipCullDistance:!1,multiDraw:!1},this.defaultAttributeValues={color:[1,1,1],uv:[0,0],uv1:[0,0]},this.index0AttributeName=void 0,this.uniformsNeedUpdate=!1,this.glslVersion=null,t!==void 0)this.setValues(t)}copy(t){return super.copy(t),this.fragmentShader=t.fragmentShader,this.vertexShader=t.vertexShader,this.uniforms=Vi(t.uniforms),this.uniformsGroups=Y0(t.uniformsGroups),this.defines=Object.assign({},t.defines),this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.fog=t.fog,this.lights=t.lights,this.clipping=t.clipping,this.extensions=Object.assign({},t.extensions),this.glslVersion=t.glslVersion,this.defaultAttributeValues=Object.assign({},t.defaultAttributeValues),this.index0AttributeName=t.index0AttributeName,this.uniformsNeedUpdate=t.uniformsNeedUpdate,this}toJSON(t){let e=super.toJSON(t);e.glslVersion=this.glslVersion,e.uniforms={};for(let i in this.uniforms){let r=this.uniforms[i].value;if(r&&r.isTexture)e.uniforms[i]={type:"t",value:r.toJSON(t).uuid};else if(r&&r.isColor)e.uniforms[i]={type:"c",value:r.getHex()};else if(r&&r.isVector2)e.uniforms[i]={type:"v2",value:r.toArray()};else if(r&&r.isVector3)e.uniforms[i]={type:"v3",value:r.toArray()};else if(r&&r.isVector4)e.uniforms[i]={type:"v4",value:r.toArray()};else if(r&&r.isMatrix3)e.uniforms[i]={type:"m3",value:r.toArray()};else if(r&&r.isMatrix4)e.uniforms[i]={type:"m4",value:r.toArray()};else e.uniforms[i]={value:r}}if(Object.keys(this.defines).length>0)e.defines=this.defines;e.vertexShader=this.vertexShader,e.fragmentShader=this.fragmentShader,e.lights=this.lights,e.clipping=this.clipping;let n={};for(let i in this.extensions)if(this.extensions[i]===!0)n[i]=!0;if(Object.keys(n).length>0)e.extensions=n;return e}fromJSON(t,e){if(super.fromJSON(t,e),t.uniforms!==void 0)for(let n in t.uniforms){let i=t.uniforms[n];switch(this.uniforms[n]={},i.type){case"t":this.uniforms[n].value=e[i.value]||null;break;case"c":this.uniforms[n].value=new Pt().setHex(i.value);break;case"v2":this.uniforms[n].value=new rt().fromArray(i.value);break;case"v3":this.uniforms[n].value=new L().fromArray(i.value);break;case"v4":this.uniforms[n].value=new fe().fromArray(i.value);break;case"m3":this.uniforms[n].value=new Xt().fromArray(i.value);break;case"m4":this.uniforms[n].value=new Kt().fromArray(i.value);break;default:this.uniforms[n].value=i.value}}if(t.defines!==void 0)this.defines=t.defines;if(t.vertexShader!==void 0)this.vertexShader=t.vertexShader;if(t.fragmentShader!==void 0)this.fragmentShader=t.fragmentShader;if(t.glslVersion!==void 0)this.glslVersion=t.glslVersion;if(t.extensions!==void 0)for(let n in t.extensions)this.extensions[n]=t.extensions[n];if(t.lights!==void 0)this.lights=t.lights;if(t.clipping!==void 0)this.clipping=t.clipping;return this}}class vo extends ln{constructor(t){super(t);this.isRawShaderMaterial=!0,this.type="RawShaderMaterial"}}class Wi extends Fe{constructor(t){super();this.isMeshStandardMaterial=!0,this.type="MeshStandardMaterial",this.defines={STANDARD:""},this.color=new Pt(16777215),this.roughness=1,this.metalness=0,this.map=null,this.lightMap=null,this.lightMapIntensity=1,this.aoMap=null,this.aoMapIntensity=1,this.emissive=new Pt(0),this.emissiveIntensity=1,this.emissiveMap=null,this.bumpMap=null,this.bumpScale=1,this.normalMap=null,this.normalMapType=0,this.normalScale=new rt(1,1),this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.roughnessMap=null,this.metalnessMap=null,this.alphaMap=null,this.envMap=null,this.envMapRotation=new Mn,this.envMapIntensity=1,this.wireframe=!1,this.wireframeLinewidth=1,this.wireframeLinecap="round",this.wireframeLinejoin="round",this.flatShading=!1,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.defines={STANDARD:""},this.color.copy(t.color),this.roughness=t.roughness,this.metalness=t.metalness,this.map=t.map,this.lightMap=t.lightMap,this.lightMapIntensity=t.lightMapIntensity,this.aoMap=t.aoMap,this.aoMapIntensity=t.aoMapIntensity,this.emissive.copy(t.emissive),this.emissiveMap=t.emissiveMap,this.emissiveIntensity=t.emissiveIntensity,this.bumpMap=t.bumpMap,this.bumpScale=t.bumpScale,this.normalMap=t.normalMap,this.normalMapType=t.normalMapType,this.normalScale.copy(t.normalScale),this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.roughnessMap=t.roughnessMap,this.metalnessMap=t.metalnessMap,this.alphaMap=t.alphaMap,this.envMap=t.envMap,this.envMapRotation.copy(t.envMapRotation),this.envMapIntensity=t.envMapIntensity,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.wireframeLinecap=t.wireframeLinecap,this.wireframeLinejoin=t.wireframeLinejoin,this.flatShading=t.flatShading,this.fog=t.fog,this}}class qc extends Wi{constructor(t){super();this.isMeshPhysicalMaterial=!0,this.defines={STANDARD:"",PHYSICAL:""},this.type="MeshPhysicalMaterial",this.anisotropyRotation=0,this.anisotropyMap=null,this.clearcoatMap=null,this.clearcoatRoughness=0,this.clearcoatRoughnessMap=null,this.clearcoatNormalScale=new rt(1,1),this.clearcoatNormalMap=null,this.ior=1.5,Object.defineProperty(this,"reflectivity",{get:function(){return Zt(2.5*(this.ior-1)/(this.ior+1),0,1)},set:function(e){this.ior=(1+0.4*e)/(1-0.4*e)}}),this.iridescenceMap=null,this.iridescenceIOR=1.3,this.iridescenceThicknessRange=[100,400],this.iridescenceThicknessMap=null,this.sheenColor=new Pt(0),this.sheenColorMap=null,this.sheenRoughness=1,this.sheenRoughnessMap=null,this.transmissionMap=null,this.thickness=0,this.thicknessMap=null,this.attenuationDistance=1/0,this.attenuationColor=new Pt(1,1,1),this.specularIntensity=1,this.specularIntensityMap=null,this.specularColor=new Pt(1,1,1),this.specularColorMap=null,this._anisotropy=0,this._clearcoat=0,this._dispersion=0,this._iridescence=0,this._retroreflectivity=0,this._sheen=0,this._transmission=0,this.setValues(t)}get anisotropy(){return this._anisotropy}set anisotropy(t){if(this._anisotropy>0!==t>0)this.version++;this._anisotropy=t}get clearcoat(){return this._clearcoat}set clearcoat(t){if(this._clearcoat>0!==t>0)this.version++;this._clearcoat=t}get iridescence(){return this._iridescence}set iridescence(t){if(this._iridescence>0!==t>0)this.version++;this._iridescence=t}get dispersion(){return this._dispersion}set dispersion(t){if(this._dispersion>0!==t>0)this.version++;this._dispersion=t}get retroreflectivity(){return this._retroreflectivity}set retroreflectivity(t){if(this._retroreflectivity>0!==t>0)this.version++;this._retroreflectivity=t}get sheen(){return this._sheen}set sheen(t){if(this._sheen>0!==t>0)this.version++;this._sheen=t}get transmission(){return this._transmission}set transmission(t){if(this._transmission>0!==t>0)this.version++;this._transmission=t}copy(t){return super.copy(t),this.defines={STANDARD:"",PHYSICAL:""},this.anisotropy=t.anisotropy,this.anisotropyRotation=t.anisotropyRotation,this.anisotropyMap=t.anisotropyMap,this.clearcoat=t.clearcoat,this.clearcoatMap=t.clearcoatMap,this.clearcoatRoughness=t.clearcoatRoughness,this.clearcoatRoughnessMap=t.clearcoatRoughnessMap,this.clearcoatNormalMap=t.clearcoatNormalMap,this.clearcoatNormalScale.copy(t.clearcoatNormalScale),this.dispersion=t.dispersion,this.ior=t.ior,this.iridescence=t.iridescence,this.iridescenceMap=t.iridescenceMap,this.iridescenceIOR=t.iridescenceIOR,this.iridescenceThicknessRange=[...t.iridescenceThicknessRange],this.iridescenceThicknessMap=t.iridescenceThicknessMap,this.retroreflectivity=t.retroreflectivity,this.sheen=t.sheen,this.sheenColor.copy(t.sheenColor),this.sheenColorMap=t.sheenColorMap,this.sheenRoughness=t.sheenRoughness,this.sheenRoughnessMap=t.sheenRoughnessMap,this.transmission=t.transmission,this.transmissionMap=t.transmissionMap,this.thickness=t.thickness,this.thicknessMap=t.thicknessMap,this.attenuationDistance=t.attenuationDistance,this.attenuationColor.copy(t.attenuationColor),this.specularIntensity=t.specularIntensity,this.specularIntensityMap=t.specularIntensityMap,this.specularColor.copy(t.specularColor),this.specularColorMap=t.specularColorMap,this}}class Yc extends Fe{constructor(t){super();this.isMeshPhongMaterial=!0,this.type="MeshPhongMaterial",this.color=new Pt(16777215),this.specular=new Pt(1118481),this.shininess=30,this.map=null,this.lightMap=null,this.lightMapIntensity=1,this.aoMap=null,this.aoMapIntensity=1,this.emissive=new Pt(0),this.emissiveIntensity=1,this.emissiveMap=null,this.bumpMap=null,this.bumpScale=1,this.normalMap=null,this.normalMapType=0,this.normalScale=new rt(1,1),this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.specularMap=null,this.alphaMap=null,this.envMap=null,this.envMapRotation=new Mn,this.combine=0,this.reflectivity=1,this.envMapIntensity=1,this.refractionRatio=0.98,this.wireframe=!1,this.wireframeLinewidth=1,this.wireframeLinecap="round",this.wireframeLinejoin="round",this.flatShading=!1,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.specular.copy(t.specular),this.shininess=t.shininess,this.map=t.map,this.lightMap=t.lightMap,this.lightMapIntensity=t.lightMapIntensity,this.aoMap=t.aoMap,this.aoMapIntensity=t.aoMapIntensity,this.emissive.copy(t.emissive),this.emissiveMap=t.emissiveMap,this.emissiveIntensity=t.emissiveIntensity,this.bumpMap=t.bumpMap,this.bumpScale=t.bumpScale,this.normalMap=t.normalMap,this.normalMapType=t.normalMapType,this.normalScale.copy(t.normalScale),this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.specularMap=t.specularMap,this.alphaMap=t.alphaMap,this.envMap=t.envMap,this.envMapRotation.copy(t.envMapRotation),this.combine=t.combine,this.reflectivity=t.reflectivity,this.envMapIntensity=t.envMapIntensity,this.refractionRatio=t.refractionRatio,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.wireframeLinecap=t.wireframeLinecap,this.wireframeLinejoin=t.wireframeLinejoin,this.flatShading=t.flatShading,this.fog=t.fog,this}}class Zc extends Fe{constructor(t){super();this.isMeshToonMaterial=!0,this.defines={TOON:""},this.type="MeshToonMaterial",this.color=new Pt(16777215),this.map=null,this.gradientMap=null,this.lightMap=null,this.lightMapIntensity=1,this.aoMap=null,this.aoMapIntensity=1,this.emissive=new Pt(0),this.emissiveIntensity=1,this.emissiveMap=null,this.bumpMap=null,this.bumpScale=1,this.normalMap=null,this.normalMapType=0,this.normalScale=new rt(1,1),this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.alphaMap=null,this.wireframe=!1,this.wireframeLinewidth=1,this.wireframeLinecap="round",this.wireframeLinejoin="round",this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.gradientMap=t.gradientMap,this.lightMap=t.lightMap,this.lightMapIntensity=t.lightMapIntensity,this.aoMap=t.aoMap,this.aoMapIntensity=t.aoMapIntensity,this.emissive.copy(t.emissive),this.emissiveMap=t.emissiveMap,this.emissiveIntensity=t.emissiveIntensity,this.bumpMap=t.bumpMap,this.bumpScale=t.bumpScale,this.normalMap=t.normalMap,this.normalMapType=t.normalMapType,this.normalScale.copy(t.normalScale),this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.alphaMap=t.alphaMap,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.wireframeLinecap=t.wireframeLinecap,this.wireframeLinejoin=t.wireframeLinejoin,this.fog=t.fog,this}}class Jc extends Fe{constructor(t){super();this.isMeshNormalMaterial=!0,this.type="MeshNormalMaterial",this.bumpMap=null,this.bumpScale=1,this.normalMap=null,this.normalMapType=0,this.normalScale=new rt(1,1),this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.wireframe=!1,this.wireframeLinewidth=1,this.flatShading=!1,this.setValues(t)}copy(t){return super.copy(t),this.bumpMap=t.bumpMap,this.bumpScale=t.bumpScale,this.normalMap=t.normalMap,this.normalMapType=t.normalMapType,this.normalScale.copy(t.normalScale),this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.flatShading=t.flatShading,this}}class Sr extends Fe{constructor(t){super();this.isMeshLambertMaterial=!0,this.type="MeshLambertMaterial",this.color=new Pt(16777215),this.map=null,this.lightMap=null,this.lightMapIntensity=1,this.aoMap=null,this.aoMapIntensity=1,this.emissive=new Pt(0),this.emissiveIntensity=1,this.emissiveMap=null,this.bumpMap=null,this.bumpScale=1,this.normalMap=null,this.normalMapType=0,this.normalScale=new rt(1,1),this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.specularMap=null,this.alphaMap=null,this.envMap=null,this.envMapRotation=new Mn,this.combine=0,this.reflectivity=1,this.envMapIntensity=1,this.refractionRatio=0.98,this.wireframe=!1,this.wireframeLinewidth=1,this.wireframeLinecap="round",this.wireframeLinejoin="round",this.flatShading=!1,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.color.copy(t.color),this.map=t.map,this.lightMap=t.lightMap,this.lightMapIntensity=t.lightMapIntensity,this.aoMap=t.aoMap,this.aoMapIntensity=t.aoMapIntensity,this.emissive.copy(t.emissive),this.emissiveMap=t.emissiveMap,this.emissiveIntensity=t.emissiveIntensity,this.bumpMap=t.bumpMap,this.bumpScale=t.bumpScale,this.normalMap=t.normalMap,this.normalMapType=t.normalMapType,this.normalScale.copy(t.normalScale),this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.specularMap=t.specularMap,this.alphaMap=t.alphaMap,this.envMap=t.envMap,this.envMapRotation.copy(t.envMapRotation),this.combine=t.combine,this.reflectivity=t.reflectivity,this.envMapIntensity=t.envMapIntensity,this.refractionRatio=t.refractionRatio,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.wireframeLinecap=t.wireframeLinecap,this.wireframeLinejoin=t.wireframeLinejoin,this.flatShading=t.flatShading,this.fog=t.fog,this}}class So extends Fe{constructor(t){super();this.isMeshDepthMaterial=!0,this.type="MeshDepthMaterial",this.depthPacking=3200,this.map=null,this.alphaMap=null,this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.wireframe=!1,this.wireframeLinewidth=1,this.setValues(t)}copy(t){return super.copy(t),this.depthPacking=t.depthPacking,this.map=t.map,this.alphaMap=t.alphaMap,this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this}}class Mo extends Fe{constructor(t){super();this.isMeshDistanceMaterial=!0,this.type="MeshDistanceMaterial",this.map=null,this.alphaMap=null,this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.setValues(t)}copy(t){return super.copy(t),this.map=t.map,this.alphaMap=t.alphaMap,this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this}}class $c extends Fe{constructor(t){super();this.isMeshMatcapMaterial=!0,this.defines={MATCAP:""},this.type="MeshMatcapMaterial",this.color=new Pt(16777215),this.matcap=null,this.map=null,this.bumpMap=null,this.bumpScale=1,this.normalMap=null,this.normalMapType=0,this.normalScale=new rt(1,1),this.displacementMap=null,this.displacementScale=1,this.displacementBias=0,this.alphaMap=null,this.wireframe=!1,this.wireframeLinewidth=1,this.flatShading=!1,this.fog=!0,this.setValues(t)}copy(t){return super.copy(t),this.defines={MATCAP:""},this.color.copy(t.color),this.matcap=t.matcap,this.map=t.map,this.bumpMap=t.bumpMap,this.bumpScale=t.bumpScale,this.normalMap=t.normalMap,this.normalMapType=t.normalMapType,this.normalScale.copy(t.normalScale),this.displacementMap=t.displacementMap,this.displacementScale=t.displacementScale,this.displacementBias=t.displacementBias,this.alphaMap=t.alphaMap,this.wireframe=t.wireframe,this.wireframeLinewidth=t.wireframeLinewidth,this.flatShading=t.flatShading,this.fog=t.fog,this}}class Kc extends We{constructor(t){super();this.isLineDashedMaterial=!0,this.type="LineDashedMaterial",this.scale=1,this.dashSize=3,this.gapSize=1,this.setValues(t)}copy(t){return super.copy(t),this.scale=t.scale,this.dashSize=t.dashSize,this.gapSize=t.gapSize,this}}function Sn(t,e){if(!t||t.constructor===e)return t;if(typeof e.BYTES_PER_ELEMENT==="number")return new e(t);return Array.prototype.slice.call(t)}function Qs(t){return t!==void 0&&t.inTangents!==void 0&&t.outTangents!==void 0}function Cf(t){function e(s,r){return t[s]-t[r]}let n=t.length,i=Array(n);for(let s=0;s!==n;++s)i[s]=s;return i.sort(e),i}function Tl(t,e,n){let i=t.length,s=new t.constructor(i);for(let r=0,a=0;a!==i;++r){let o=n[r]*e;for(let l=0;l!==e;++l)s[a++]=t[o+l]}return s}function If(t,e,n,i){let s=1,r=t[0];while(r!==void 0&&r[i]===void 0)r=t[s++];if(r===void 0)return;let a=r[i];if(a===void 0)return;if(Array.isArray(a))do{if(a=r[i],a!==void 0)e.push(r.time),n.push(...a);r=t[s++]}while(r!==void 0);else if(a.toArray!==void 0)do{if(a=r[i],a!==void 0)e.push(r.time),a.toArray(n,n.length);r=t[s++]}while(r!==void 0);else do{if(a=r[i],a!==void 0)e.push(r.time),n.push(a);r=t[s++]}while(r!==void 0)}function $0(t,e,n,i,s=30){let r=t.clone();r.name=e;let a=[];for(let l=0;l<r.tracks.length;++l){let c=r.tracks[l],h=c.getValueSize(),d=[],u=[];for(let f=0;f<c.times.length;++f){let g=c.times[f]*s;if(g<n||g>=i)continue;d.push(c.times[f]);for(let _=0;_<h;++_)u.push(c.values[f*h+_])}if(d.length===0)continue;c.times=Sn(d,c.times.constructor),c.values=Sn(u,c.values.constructor),a.push(c)}r.tracks=a;let o=1/0;for(let l=0;l<r.tracks.length;++l)if(o>r.tracks[l].times[0])o=r.tracks[l].times[0];for(let l=0;l<r.tracks.length;++l)r.tracks[l].shift(-1*o);return r.resetDuration(),r}function K0(t,e=0,n=t,i=30){if(i<=0)i=30;let s=n.tracks.length,r=e/i;for(let a=0;a<s;++a){let o=n.tracks[a],l=o.ValueTypeName;if(l==="bool"||l==="string")continue;let c=t.tracks.find(function(p){return p.name===o.name&&p.ValueTypeName===l});if(c===void 0)continue;let h=0,d=o.getValueSize();if(o.createInterpolant.isInterpolantFactoryMethodGLTFCubicSpline)h=d/3;let u=0,f=c.getValueSize();if(c.createInterpolant.isInterpolantFactoryMethodGLTFCubicSpline)u=f/3;let g=o.times.length-1,_;if(r<=o.times[0]){let p=h,b=d-h;_=o.values.slice(p,b)}else if(r>=o.times[g]){let p=g*d+h,b=p+d-h;_=o.values.slice(p,b)}else{let p=o.createInterpolant(),b=h,S=d-h;p.evaluate(r),_=p.resultBuffer.slice(b,S)}if(l==="quaternion")new Ce().fromArray(_).normalize().conjugate().toArray(_);let m=c.times.length;for(let p=0;p<m;++p){let b=p*f+u;if(l==="quaternion")Ce.multiplyQuaternionsFlat(c.values,b,_,0,c.values,b);else{let S=f-u*2;for(let x=0;x<S;++x)c.values[b+x]-=_[x]}}}return t.blendMode=2501,t}class Pf{static convertArray(t,e){return Sn(t,e)}static isTypedArray(t){return Qd(t)}static hasTangents(t){return Qs(t)}static getKeyframeOrder(t){return Cf(t)}static sortedArray(t,e,n){return Tl(t,e,n)}static flattenJSON(t,e,n,i){If(t,e,n,i)}static subclip(t,e,n,i,s=30){return $0(t,e,n,i,s)}static makeClipAdditive(t,e=0,n=t,i=30){return K0(t,e,n,i)}}class Xi{constructor(t,e,n,i){this.parameterPositions=t,this._cachedIndex=0,this.resultBuffer=i!==void 0?i:new e.constructor(n),this.sampleValues=e,this.valueSize=n,this.settings=null,this.DefaultSettings_={}}evaluate(t){let e=this.parameterPositions,n=this._cachedIndex,i=e[n],s=e[n-1];t:{e:{let r;n:{i:if(!(t<i)){for(let a=n+2;;){if(i===void 0){if(t<s)break i;return n=e.length,this._cachedIndex=n,this.copySampleValue_(n-1)}if(n===a)break;if(s=i,i=e[++n],t<i)break e}r=e.length;break n}if(!(t>=s)){let a=e[1];if(t<a)n=2,s=a;for(let o=n-2;;){if(s===void 0)return this._cachedIndex=0,this.copySampleValue_(0);if(n===o)break;if(i=s,s=e[--n-1],t>=s)break e}r=n,n=0;break n}break t}while(n<r){let a=n+r>>>1;if(t<e[a])r=a;else n=a+1}if(i=e[n],s=e[n-1],s===void 0)return this._cachedIndex=0,this.copySampleValue_(0);if(i===void 0)return n=e.length,this._cachedIndex=n,this.copySampleValue_(n-1)}this._cachedIndex=n,this.intervalChanged_(n,s,i)}return this.interpolate_(n,s,t,i)}getSettings_(){return this.settings||this.DefaultSettings_}copySampleValue_(t){let e=this.resultBuffer,n=this.sampleValues,i=this.valueSize,s=t*i;for(let r=0;r!==i;++r)e[r]=n[s+r];return e}interpolate_(){throw Error("THREE.Interpolant: Call to abstract method.")}intervalChanged_(){}}class jc extends Xi{constructor(t,e,n,i){super(t,e,n,i);this._weightPrev=-0,this._offsetPrev=-0,this._weightNext=-0,this._offsetNext=-0,this.DefaultSettings_={endingStart:2400,endingEnd:2400}}intervalChanged_(t,e,n){let i=this.parameterPositions,s=t-2,r=t+1,a=i[s],o=i[r];if(a===void 0)switch(this.getSettings_().endingStart){case 2401:s=t,a=2*e-n;break;case 2402:s=i.length-2,a=e+i[s]-i[s+1];break;default:s=t,a=n}if(o===void 0)switch(this.getSettings_().endingEnd){case 2401:r=t,o=2*n-e;break;case 2402:r=1,o=n+i[1]-i[0];break;default:r=t-1,o=e}let l=(n-e)*0.5,c=this.valueSize;this._weightPrev=l/(e-a),this._weightNext=l/(o-n),this._offsetPrev=s*c,this._offsetNext=r*c}interpolate_(t,e,n,i){let s=this.resultBuffer,r=this.sampleValues,a=this.valueSize,o=t*a,l=o-a,c=this._offsetPrev,h=this._offsetNext,d=this._weightPrev,u=this._weightNext,f=(n-e)/(i-e),g=f*f,_=g*f,m=-d*_+2*d*g-d*f,p=(1+d)*_+(-1.5-2*d)*g+(-0.5+d)*f+1,b=(-1-u)*_+(1.5+u)*g+0.5*f,S=u*_-u*g;for(let x=0;x!==a;++x)s[x]=m*r[c+x]+p*r[l+x]+b*r[o+x]+S*r[h+x];return s}}class bo extends Xi{constructor(t,e,n,i){super(t,e,n,i)}interpolate_(t,e,n,i){let s=this.resultBuffer,r=this.sampleValues,a=this.valueSize,o=t*a,l=o-a,c=(n-e)/(i-e),h=1-c;for(let d=0;d!==a;++d)s[d]=r[l+d]*h+r[o+d]*c;return s}}class Qc extends Xi{constructor(t,e,n,i){super(t,e,n,i)}interpolate_(t){return this.copySampleValue_(t-1)}}class th extends Xi{interpolate_(t,e,n,i){let s=this.resultBuffer,r=this.sampleValues,a=this.valueSize,o=t*a,l=o-a,c=this.inTangents,h=this.outTangents;if(!c||!h){let f=(n-e)/(i-e),g=1-f;for(let _=0;_!==a;++_)s[_]=r[l+_]*g+r[o+_]*f;return s}let d=a*2,u=t-1;for(let f=0;f!==a;++f){let g=r[l+f],_=r[o+f],m=u*d+f*2,p=h[m],b=h[m+1],S=t*d+f*2,x=c[S],T=c[S+1],I=Q0(n,e,p,x,i);s[f]=Lf(I,g,b,T,_)}return s}}function Lf(t,e,n,i,s){let r=1-t;return r*r*r*e+3*r*r*t*n+3*r*t*t*i+t*t*t*s}function j0(t,e,n,i,s){let r=1-t;return 3*r*r*(n-e)+6*r*t*(i-n)+3*t*t*(s-i)}function Q0(t,e,n,i,s){let r=(t-e)/(s-e);for(let a=0;a<8;a++){let o=Lf(r,e,n,i,s)-t;if(Math.abs(o)<0.0000000001)break;let l=j0(r,e,n,i,s);if(Math.abs(l)<0.0000000001)break;r=Math.max(0,Math.min(1,r-o/l))}return r}class cn{constructor(t,e,n,i){if(t===void 0)throw Error("THREE.KeyframeTrack: track name is undefined");if(e===void 0||e.length===0)throw Error("THREE.KeyframeTrack: no keyframes in track named "+t);this.name=t,this.times=Sn(e,this.TimeBufferType),this.values=Sn(n,this.ValueBufferType),this.setInterpolation(i||this.DefaultInterpolation)}static toJSON(t){let e=t.constructor,n;if(e.toJSON!==this.toJSON)n=e.toJSON(t);else{n={name:t.name,times:Sn(t.times,Array),values:Sn(t.values,Array)};let i=t.getInterpolation();if(i!==t.DefaultInterpolation)n.interpolation=i;if(Qs(t.settings))n.settings={inTangents:Sn(t.settings.inTangents,Array),outTangents:Sn(t.settings.outTangents,Array)}}return n.type=t.ValueTypeName,n}InterpolantFactoryMethodDiscrete(t){return new Qc(this.times,this.values,this.getValueSize(),t)}InterpolantFactoryMethodLinear(t){return new bo(this.times,this.values,this.getValueSize(),t)}InterpolantFactoryMethodSmooth(t){return new jc(this.times,this.values,this.getValueSize(),t)}InterpolantFactoryMethodBezier(t){let e=new th(this.times,this.values,this.getValueSize(),t);if(this.settings)e.inTangents=this.settings.inTangents,e.outTangents=this.settings.outTangents;return e}setInterpolation(t){let e;switch(t){case 2300:e=this.InterpolantFactoryMethodDiscrete;break;case 2301:e=this.InterpolantFactoryMethodLinear;break;case 2302:e=this.InterpolantFactoryMethodSmooth;break;case 2303:e=this.InterpolantFactoryMethodBezier;break}if(e===void 0){let n="unsupported interpolation for "+this.ValueTypeName+" keyframe track named "+this.name;if(this.createInterpolant===void 0)if(t!==this.DefaultInterpolation)this.setInterpolation(this.DefaultInterpolation);else throw Error(n);return Mt("KeyframeTrack:",n),this}return this.createInterpolant=e,this}getInterpolation(){switch(this.createInterpolant){case this.InterpolantFactoryMethodDiscrete:return 2300;case this.InterpolantFactoryMethodLinear:return 2301;case this.InterpolantFactoryMethodSmooth:return 2302;case this.InterpolantFactoryMethodBezier:return 2303}}getValueSize(){return this.values.length/this.times.length}shift(t){if(t!==0){let e=this.times;for(let n=0,i=e.length;n!==i;++n)e[n]+=t}return this}scale(t){if(t!==1){let e=this.times;for(let n=0,i=e.length;n!==i;++n)e[n]*=t;if(Qs(this.settings))Nu(this.settings.inTangents,t),Nu(this.settings.outTangents,t)}return this}trim(t,e){let n=this.times,i=n.length,s=0,r=i-1;while(s!==i&&n[s]<t)++s;while(r!==-1&&n[r]>e)--r;if(++r,s!==0||r!==i){if(s>=r)r=Math.max(r,1),s=r-1;let a=this.getValueSize();this.times=n.slice(s,r),this.values=this.values.slice(s*a,r*a)}return this}validate(){let t=!0,e=this.getValueSize();if(e-Math.floor(e)!==0)Wt("KeyframeTrack: Invalid value size in track.",this),t=!1;let n=this.times,i=this.values,s=n.length;if(s===0)Wt("KeyframeTrack: Track is empty.",this),t=!1;let r=null;for(let a=0;a!==s;a++){let o=n[a];if(typeof o==="number"&&isNaN(o)){Wt("KeyframeTrack: Time is not a valid number.",this,a,o),t=!1;break}if(r!==null&&r>o){Wt("KeyframeTrack: Out of order keys.",this,a,o,r),t=!1;break}r=o}if(i!==void 0){if(Qd(i))for(let a=0,o=i.length;a!==o;++a){let l=i[a];if(isNaN(l)){Wt("KeyframeTrack: Value is not a valid number.",this,a,l),t=!1;break}}}return t}optimize(){let t=this.times.slice(),e=this.values.slice(),n=this.getValueSize(),i=this.getInterpolation()===2302,s=t.length-1,r=1;for(let a=1;a<s;++a){let o=!1,l=t[a],c=t[a+1];if(l!==c&&(a!==1||l!==t[0]))if(!i){let h=a*n,d=h-n,u=h+n;for(let f=0;f!==n;++f){let g=e[h+f];if(g!==e[d+f]||g!==e[u+f]){o=!0;break}}}else o=!0;if(o){if(a!==r){t[r]=t[a];let h=a*n,d=r*n;for(let u=0;u!==n;++u)e[d+u]=e[h+u]}++r}}if(s>0){t[r]=t[s];for(let a=s*n,o=r*n,l=0;l!==n;++l)e[o+l]=e[a+l];++r}if(r!==t.length)this.times=t.slice(0,r),this.values=e.slice(0,r*n);else this.times=t,this.values=e;return this}clone(){let t=this.times.slice(),e=this.values.slice(),i=new this.constructor(this.name,t,e);if(i.createInterpolant=this.createInterpolant,Qs(this.settings))i.settings={inTangents:this.settings.inTangents.slice(),outTangents:this.settings.outTangents.slice()};return i}}function Nu(t,e){for(let n=0,i=t.length;n!==i;n+=2)t[n]*=e}cn.prototype.ValueTypeName="";cn.prototype.TimeBufferType=Float32Array;cn.prototype.ValueBufferType=Float32Array;cn.prototype.DefaultInterpolation=2301;class _i extends cn{constructor(t,e,n){super(t,e,n)}}_i.prototype.ValueTypeName="bool";_i.prototype.ValueBufferType=Array;_i.prototype.DefaultInterpolation=2300;_i.prototype.InterpolantFactoryMethodLinear=void 0;_i.prototype.InterpolantFactoryMethodSmooth=void 0;class To extends cn{constructor(t,e,n,i){super(t,e,n,i)}}To.prototype.ValueTypeName="color";class Mr extends cn{constructor(t,e,n,i){super(t,e,n,i)}}Mr.prototype.ValueTypeName="number";class eh extends Xi{constructor(t,e,n,i){super(t,e,n,i)}interpolate_(t,e,n,i){let s=this.resultBuffer,r=this.sampleValues,a=this.valueSize,o=(n-e)/(i-e),l=t*a;for(let c=l+a;l!==c;l+=4)Ce.slerpFlat(s,0,r,l-a,r,l,o);return s}}class br extends cn{constructor(t,e,n,i){super(t,e,n,i)}InterpolantFactoryMethodLinear(t){return new eh(this.times,this.values,this.getValueSize(),t)}}br.prototype.ValueTypeName="quaternion";br.prototype.InterpolantFactoryMethodSmooth=void 0;class xi extends cn{constructor(t,e,n){super(t,e,n)}}xi.prototype.ValueTypeName="string";xi.prototype.ValueBufferType=Array;xi.prototype.DefaultInterpolation=2300;xi.prototype.InterpolantFactoryMethodLinear=void 0;xi.prototype.InterpolantFactoryMethodSmooth=void 0;class Eo extends cn{constructor(t,e,n,i){super(t,e,n,i)}}Eo.prototype.ValueTypeName="vector";class ys{constructor(t="",e=-1,n=[],i=2500){if(this.name=t,this.tracks=n,this.duration=e,this.blendMode=i,this.uuid=rn(),this.userData={},this.duration<0)this.resetDuration()}static parse(t){let e=[],n=t.tracks,i=1/(t.fps||1);for(let r=0,a=n.length;r!==a;++r)e.push(e_(n[r]).scale(i));let s=new this(t.name,t.duration,e,t.blendMode);return s.uuid=t.uuid,s.userData=JSON.parse(t.userData||"{}"),s}static toJSON(t){let e=[],n=t.tracks,i={name:t.name,duration:t.duration,tracks:e,uuid:t.uuid,blendMode:t.blendMode,userData:JSON.stringify(t.userData)};for(let s=0,r=n.length;s!==r;++s)e.push(cn.toJSON(n[s]));return i}static CreateFromMorphTargetSequence(t,e,n,i){let s=e.length,r=[];for(let a=0;a<s;a++){let o=[],l=[];o.push((a+s-1)%s,a,(a+1)%s),l.push(0,1,0);let c=Cf(o);if(o=Tl(o,1,c),l=Tl(l,1,c),!i&&o[0]===0)o.push(s),l.push(l[0]);r.push(new Mr(".morphTargetInfluences["+e[a].name+"]",o,l).scale(1/n))}return new this(t,-1,r)}static findByName(t,e){let n=t;if(!Array.isArray(t)){let i=t;n=i.geometry&&i.geometry.animations||i.animations}for(let i=0;i<n.length;i++)if(n[i].name===e)return n[i];return null}static CreateClipsFromMorphTargetSequences(t,e,n){let i={},s=/^([\w-]*?)([\d]+)$/;for(let a=0,o=t.length;a<o;a++){let l=t[a],c=l.name.match(s);if(c&&c.length>1){let h=c[1],d=i[h];if(!d)i[h]=d=[];d.push(l)}}let r=[];for(let a in i)r.push(this.CreateFromMorphTargetSequence(a,i[a],e,n));return r}resetDuration(){let t=this.tracks,e=0;for(let n=0,i=t.length;n!==i;++n){let s=this.tracks[n];e=Math.max(e,s.times[s.times.length-1])}return this.duration=e,this}trim(){for(let t=0;t<this.tracks.length;t++)this.tracks[t].trim(0,this.duration);return this}validate(){let t=!0;for(let e=0;e<this.tracks.length;e++)t=t&&this.tracks[e].validate();return t}optimize(){for(let t=0;t<this.tracks.length;t++)this.tracks[t].optimize();return this}clone(){let t=[];for(let n=0;n<this.tracks.length;n++)t.push(this.tracks[n].clone());let e=new this.constructor(this.name,this.duration,t,this.blendMode);return e.userData=JSON.parse(JSON.stringify(this.userData)),e}toJSON(){return this.constructor.toJSON(this)}}function t_(t){switch(t.toLowerCase()){case"scalar":case"double":case"float":case"number":case"integer":return Mr;case"vector":case"vector2":case"vector3":case"vector4":return Eo;case"color":return To;case"quaternion":return br;case"bool":case"boolean":return _i;case"string":return xi}throw Error("THREE.KeyframeTrack: Unsupported typeName: "+t)}function e_(t){if(t.type===void 0)throw Error("THREE.KeyframeTrack: track type undefined, can not parse");let e=t_(t.type);if(t.times===void 0){let i=[],s=[];If(t.keys,i,s,"value"),t.times=i,t.values=s}let n;if(e.parse!==void 0)n=e.parse(t);else n=new e(t.name,t.times,t.values,t.interpolation);if(Qs(t.settings))n.settings={inTangents:Sn(t.settings.inTangents,Float32Array),outTangents:Sn(t.settings.outTangents,Float32Array)};return n}var Pn={enabled:!1,files:{},add:function(t,e){if(this.enabled===!1)return;if(Du(t))return;this.files[t]=e},get:function(t){if(this.enabled===!1)return;if(Du(t))return;return this.files[t]},remove:function(t){delete this.files[t]},clear:function(){this.files={}}};function Du(t){try{let e=t.slice(t.indexOf(":")+1);return new URL(e).protocol==="blob:"}catch(e){return!1}}class Ao{constructor(t,e,n){let i=this,s=!1,r=0,a=0,o=void 0,l=[];this.onStart=void 0,this.onLoad=t,this.onProgress=e,this.onError=n,this._abortController=null,this.itemStart=function(c){if(a++,s===!1){if(i.onStart!==void 0)i.onStart(c,r,a)}s=!0},this.itemEnd=function(c){if(r++,i.onProgress!==void 0)i.onProgress(c,r,a);if(r===a){if(s=!1,i.onLoad!==void 0)i.onLoad()}},this.itemError=function(c){if(i.onError!==void 0)i.onError(c)},this.resolveURL=function(c){if(c=c.normalize("NFC"),o)return o(c);return c},this.setURLModifier=function(c){return o=c,this},this.addHandler=function(c,h){return l.push(c,h),this},this.removeHandler=function(c){let h=l.indexOf(c);if(h!==-1)l.splice(h,2);return this},this.getHandler=function(c){for(let h=0,d=l.length;h<d;h+=2){let u=l[h],f=l[h+1];if(u.global)u.lastIndex=0;if(u.test(c))return f}return null},this.abort=function(){return this.abortController.abort(),this._abortController=null,this}}get abortController(){if(!this._abortController)this._abortController=new AbortController;return this._abortController}}var Nf=new Ao;class Be{constructor(t){if(this.manager=t!==void 0?t:Nf,this.crossOrigin="anonymous",this.withCredentials=!1,this.path="",this.resourcePath="",this.requestHeader={},typeof __THREE_DEVTOOLS__<"u")__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("observe",{detail:this}))}load(){}loadAsync(t,e){let n=this;return new Promise(function(i,s){n.load(t,i,e,s)})}parse(){}setCrossOrigin(t){return this.crossOrigin=t,this}setWithCredentials(t){return this.withCredentials=t,this}setPath(t){return this.path=t,this}setResourcePath(t){return this.resourcePath=t,this}setRequestHeader(t){return this.requestHeader=t,this}abort(){return this}}Be.DEFAULT_MATERIAL_NAME="__DEFAULT";var Wn={};class Df extends Error{constructor(t,e){super(t);this.response=e}}class fn extends Be{constructor(t){super(t);this.mimeType="",this.responseType="",this._abortController=new AbortController}load(t,e,n,i){if(t===void 0)t="";if(this.path!==void 0)t=this.path+t;t=this.manager.resolveURL(t);let s=Pn.get(`file:${t}`);if(s!==void 0){this.manager.itemStart(t),setTimeout(()=>{if(e)e(s);this.manager.itemEnd(t)},0);return}if(Wn[t]!==void 0){Wn[t].push({onLoad:e,onProgress:n,onError:i});return}Wn[t]=[],Wn[t].push({onLoad:e,onProgress:n,onError:i});let r=new Request(t,{headers:new Headers(this.requestHeader),credentials:this.withCredentials?"include":"same-origin",signal:typeof AbortSignal.any==="function"?AbortSignal.any([this._abortController.signal,this.manager.abortController.signal]):this._abortController.signal}),a=this.mimeType,o=this.responseType;fetch(r).then((l)=>{if(l.status===200||l.status===0){if(l.status===0)Mt("FileLoader: HTTP Status 0 received.");if(typeof ReadableStream>"u"||l.body===void 0||l.body.getReader===void 0)return l;let c=Wn[t],h=l.body.getReader(),d=l.headers.get("X-File-Size")||l.headers.get("Content-Length"),u=d?parseInt(d):0,f=u!==0,g=0,_=new ReadableStream({start(m){p();function p(){h.read().then(({done:b,value:S})=>{if(b)m.close();else{g+=S.byteLength;let x=new ProgressEvent("progress",{lengthComputable:f,loaded:g,total:u});for(let T=0,I=c.length;T<I;T++){let R=c[T];if(R.onProgress)R.onProgress(x)}m.enqueue(S),p()}},(b)=>{m.error(b)})}}});return new Response(_)}else throw new Df(`fetch for "${l.url}" responded with ${l.status}: ${l.statusText}`,l)}).then((l)=>{switch(o){case"arraybuffer":return l.arrayBuffer();case"blob":return l.blob();case"document":return l.text().then((c)=>new DOMParser().parseFromString(c,a));case"json":return l.json();default:if(a==="")return l.text();else{let h=/charset="?([^;"\s]*)"?/i.exec(a),d=h&&h[1]?h[1].toLowerCase():void 0,u=new TextDecoder(d);return l.arrayBuffer().then((f)=>u.decode(f))}}}).then((l)=>{Pn.add(`file:${t}`,l);let c=Wn[t];delete Wn[t];for(let h=0,d=c.length;h<d;h++){let u=c[h];if(u.onLoad)u.onLoad(l)}}).catch((l)=>{let c=Wn[t];if(c===void 0)throw this.manager.itemError(t),l;delete Wn[t];for(let h=0,d=c.length;h<d;h++){let u=c[h];if(u.onError)u.onError(l)}this.manager.itemError(t)}).finally(()=>{this.manager.itemEnd(t)}),this.manager.itemStart(t)}setResponseType(t){return this.responseType=t,this}setMimeType(t){return this.mimeType=t,this}abort(){return this._abortController.abort(),this._abortController=new AbortController,this}}class Uf extends Be{constructor(t){super(t)}load(t,e,n,i){let s=this,r=new fn(this.manager);r.setPath(this.path),r.setRequestHeader(this.requestHeader),r.setWithCredentials(this.withCredentials),r.load(t,function(a){try{e(s.parse(JSON.parse(a)))}catch(o){if(i)i(o);else Wt(o);s.manager.itemError(t)}},n,i)}parse(t){let e=[];for(let n=0;n<t.length;n++){let i=ys.parse(t[n]);e.push(i)}return e}}class Ff extends Be{constructor(t){super(t)}load(t,e,n,i){let s=this,r=[],a=new mr,o=new fn(this.manager);o.setPath(this.path),o.setResponseType("arraybuffer"),o.setRequestHeader(this.requestHeader),o.setWithCredentials(s.withCredentials);let l=0;function c(h){o.load(t[h],function(d){let u=s.parse(d,!0);if(r[h]={width:u.width,height:u.height,format:u.format,mipmaps:u.mipmaps},l+=1,l===6){if(u.mipmapCount===1)a.minFilter=1006;if(a.image=r,a.format=u.format,a.needsUpdate=!0,e)e(a)}},n,i)}if(Array.isArray(t))for(let h=0,d=t.length;h<d;++h)c(h);else o.load(t,function(h){let d=s.parse(h,!0);if(d.isCubemap){let u=d.mipmaps.length/d.mipmapCount;for(let f=0;f<u;f++){r[f]={mipmaps:[]};for(let g=0;g<d.mipmapCount;g++)r[f].mipmaps.push(d.mipmaps[f*d.mipmapCount+g]),r[f].format=d.format,r[f].width=d.width,r[f].height=d.height}a.image=r}else a.image.width=d.width,a.image.height=d.height,a.mipmaps=d.mipmaps;if(d.mipmapCount===1)a.minFilter=1006;if(a.format=d.format,a.needsUpdate=!0,e)e(a)},n,i);return a}}var hs=new WeakMap;class vs extends Be{constructor(t){super(t)}load(t,e,n,i){if(this.path!==void 0)t=this.path+t;t=this.manager.resolveURL(t);let s=this,r=Pn.get(`image:${t}`);if(r!==void 0){if(r.complete===!0)s.manager.itemStart(t),setTimeout(function(){if(e)e(r);s.manager.itemEnd(t)},0);else{let h=hs.get(r);if(h===void 0)h=[],hs.set(r,h);h.push({onLoad:e,onError:i})}return r}let a=_s("img");function o(){if(c(),e)e(this);let h=hs.get(this)||[];for(let d=0;d<h.length;d++){let u=h[d];if(u.onLoad)u.onLoad(this)}hs.delete(this),s.manager.itemEnd(t)}function l(h){if(c(),i)i(h);Pn.remove(`image:${t}`);let d=hs.get(this)||[];for(let u=0;u<d.length;u++){let f=d[u];if(f.onError)f.onError(h)}hs.delete(this),s.manager.itemError(t),s.manager.itemEnd(t)}function c(){a.removeEventListener("load",o,!1),a.removeEventListener("error",l,!1)}if(a.addEventListener("load",o,!1),a.addEventListener("error",l,!1),t.slice(0,5)!=="data:"){if(this.crossOrigin!==void 0)a.crossOrigin=this.crossOrigin}return Pn.add(`image:${t}`,a),s.manager.itemStart(t),a.src=t,a}}class Of extends Be{constructor(t){super(t)}load(t,e,n,i){let s=new Rs;s.colorSpace="srgb";let r=new vs(this.manager);r.setCrossOrigin(this.crossOrigin),r.setPath(this.path);let a=0;function o(l){r.load(t[l],function(c){if(s.images[l]=c,a++,a===6){if(s.needsUpdate=!0,e)e(s)}},void 0,i)}for(let l=0;l<t.length;++l)o(l);return s}}class Bf extends Be{constructor(t){super(t)}load(t,e,n,i){let s=this,r=new an,a=new fn(this.manager);return a.setResponseType("arraybuffer"),a.setRequestHeader(this.requestHeader),a.setPath(this.path),a.setWithCredentials(s.withCredentials),a.load(t,function(o){let l;try{l=s.parse(o)}catch(c){if(i!==void 0)i(c);else Wt(c);return}if(s._applyTexData(r,l),e)e(r,l)},n,i),r}createDataTexture(t){let e=new an;return this._applyTexData(e,this.parse(t)),e}_applyTexData(t,e){if(e.image!==void 0)t.image=e.image;else if(e.data!==void 0)t.image.width=e.width,t.image.height=e.height,t.image.data=e.data;if(t.wrapS=e.wrapS!==void 0?e.wrapS:1001,t.wrapT=e.wrapT!==void 0?e.wrapT:1001,t.magFilter=e.magFilter!==void 0?e.magFilter:1006,t.minFilter=e.minFilter!==void 0?e.minFilter:1006,t.anisotropy=e.anisotropy!==void 0?e.anisotropy:1,e.colorSpace!==void 0)t.colorSpace=e.colorSpace;if(e.flipY!==void 0)t.flipY=e.flipY;if(e.format!==void 0)t.format=e.format;if(e.type!==void 0)t.type=e.type;if(e.mipmaps!==void 0)t.mipmaps=e.mipmaps,t.minFilter=1008;if(e.mipmapCount===1)t.minFilter=1006;if(e.generateMipmaps!==void 0)t.generateMipmaps=e.generateMipmaps;t.needsUpdate=!0}}class zf extends Be{constructor(t){super(t)}load(t,e,n,i){let s=new be,r=new vs(this.manager);return r.setCrossOrigin(this.crossOrigin),r.setPath(this.path),r.load(t,function(a){if(s.image=a,s.needsUpdate=!0,e!==void 0)e(s)},n,i),s}}class jn extends ae{constructor(t,e=1){super();this.isLight=!0,this.type="Light",this.color=new Pt(t),this.intensity=e}copy(t,e){return super.copy(t,e),this.color.copy(t.color),this.intensity=t.intensity,this}toJSON(t){let e=super.toJSON(t);return e.object.color=this.color.getHex(),e.object.intensity=this.intensity,e}}class nh extends jn{constructor(t,e,n){super(t,n);this.isHemisphereLight=!0,this.type="HemisphereLight",this.position.copy(ae.DEFAULT_UP),this.updateMatrix(),this.groundColor=new Pt(e)}copy(t,e){return super.copy(t,e),this.groundColor.copy(t.groundColor),this}toJSON(t){let e=super.toJSON(t);return e.object.groundColor=this.groundColor.getHex(),e}}var pl=new Kt,Uu=new L,Fu=new L;class Tr{constructor(t){this.camera=t,this.intensity=1,this.bias=0,this.biasNode=null,this.normalBias=0,this.radius=1,this.blurSamples=8,this.mapSize=new rt(512,512),this.mapType=1009,this.map=null,this.mapPass=null,this.matrix=new Kt,this.autoUpdate=!0,this.needsUpdate=!1,this._frustum=new ui,this._frameExtents=new rt(1,1),this._viewportCount=1,this._viewports=[new fe(0,0,1,1)]}getViewportCount(){return this._viewportCount}getCamera(){return this.camera}getFrustum(){return this._frustum}updateMatrices(t){let e=this.camera;Uu.setFromMatrixPosition(t.matrixWorld),e.position.copy(Uu),Fu.setFromMatrixPosition(t.target.matrixWorld),e.lookAt(Fu),e.updateMatrixWorld(),this._updateMatrix(e,this.matrix,this._frustum)}_updateMatrix(t,e,n,i){pl.multiplyMatrices(t.projectionMatrix,t.matrixWorldInverse),n.setFromProjectionMatrix(pl,t.coordinateSystem,t.reversedDepth);let s=this._frameExtents,r=i?i.z/s.x:1,a=i?i.w/s.y:1,o=i?i.x/s.x:0,l=i?i.y/s.y:0;if(t.coordinateSystem===2001||t.reversedDepth)e.set(0.5*r,0,0,0.5*r+o,0,0.5*a,0,0.5*a+l,0,0,1,0,0,0,0,1);else e.set(0.5*r,0,0,0.5*r+o,0,0.5*a,0,0.5*a+l,0,0,0.5,0.5,0,0,0,1);e.multiply(pl)}getViewport(t){return this._viewports[t]}getFrameExtents(){return this._frameExtents}dispose(){if(this.map)this.map.dispose();if(this.mapPass)this.mapPass.dispose()}copy(t){return this.camera=t.camera.clone(),this.intensity=t.intensity,this.bias=t.bias,this.radius=t.radius,this.autoUpdate=t.autoUpdate,this.needsUpdate=t.needsUpdate,this.normalBias=t.normalBias,this.blurSamples=t.blurSamples,this.mapSize.copy(t.mapSize),this.biasNode=t.biasNode,this}clone(){return new this.constructor().copy(this)}toJSON(){let t={};return t.intensity=this.intensity,t.bias=this.bias,t.normalBias=this.normalBias,t.radius=this.radius,t.blurSamples=this.blurSamples,t.mapSize=this.mapSize.toArray(),t.camera=this.camera.toJSON(!1).object,delete t.camera.matrix,t}}var _a=new L,xa=new Ce,Cn=new L;class Er extends ae{constructor(){super();this.isCamera=!0,this.type="Camera",this.matrixWorldInverse=new Kt,this.projectionMatrix=new Kt,this.projectionMatrixInverse=new Kt,this.coordinateSystem=2000,this._reversedDepth=!1}get reversedDepth(){return this._reversedDepth}copy(t,e){return super.copy(t,e),this.matrixWorldInverse.copy(t.matrixWorldInverse),this.projectionMatrix.copy(t.projectionMatrix),this.projectionMatrixInverse.copy(t.projectionMatrixInverse),this.coordinateSystem=t.coordinateSystem,this}getWorldDirection(t){return super.getWorldDirection(t).negate()}updateMatrixWorld(t){if(super.updateMatrixWorld(t),this.matrixWorld.decompose(_a,xa,Cn),Cn.x===1&&Cn.y===1&&Cn.z===1)this.matrixWorldInverse.copy(this.matrixWorld).invert();else this.matrixWorldInverse.compose(_a,xa,Cn.set(1,1,1)).invert()}updateWorldMatrix(t,e,n=!1){if(super.updateWorldMatrix(t,e,n),this.matrixWorld.decompose(_a,xa,Cn),Cn.x===1&&Cn.y===1&&Cn.z===1)this.matrixWorldInverse.copy(this.matrixWorld).invert();else this.matrixWorldInverse.compose(_a,xa,Cn.set(1,1,1)).invert()}clone(){return new this.constructor().copy(this)}}var ai=new L,Ou=new rt,Bu=new rt;class De extends Er{constructor(t=50,e=1,n=0.1,i=2000){super();this.isPerspectiveCamera=!0,this.type="PerspectiveCamera",this.fov=t,this.zoom=1,this.near=n,this.far=i,this.focus=10,this.aspect=e,this.view=null,this.filmGauge=35,this.filmOffset=0,this.updateProjectionMatrix()}copy(t,e){return super.copy(t,e),this.fov=t.fov,this.zoom=t.zoom,this.near=t.near,this.far=t.far,this.focus=t.focus,this.aspect=t.aspect,this.view=t.view===null?null:Object.assign({},t.view),this.filmGauge=t.filmGauge,this.filmOffset=t.filmOffset,this}setFocalLength(t){let e=0.5*this.getFilmHeight()/t;this.fov=Di*2*Math.atan(e),this.updateProjectionMatrix()}getFocalLength(){let t=Math.tan(Ni*0.5*this.fov);return 0.5*this.getFilmHeight()/t}getEffectiveFOV(){return Di*2*Math.atan(Math.tan(Ni*0.5*this.fov)/this.zoom)}getFilmWidth(){return this.filmGauge*Math.min(this.aspect,1)}getFilmHeight(){return this.filmGauge/Math.max(this.aspect,1)}getViewBounds(t,e,n){ai.set(-1,-1,0.5).applyMatrix4(this.projectionMatrixInverse),e.set(ai.x,ai.y).multiplyScalar(-t/ai.z),ai.set(1,1,0.5).applyMatrix4(this.projectionMatrixInverse),n.set(ai.x,ai.y).multiplyScalar(-t/ai.z)}getViewSize(t,e){return this.getViewBounds(t,Ou,Bu),e.subVectors(Bu,Ou)}setViewOffset(t,e,n,i,s,r){if(this.aspect=t/e,this.view===null)this.view={enabled:!0,fullWidth:1,fullHeight:1,offsetX:0,offsetY:0,width:1,height:1};this.view.enabled=!0,this.view.fullWidth=t,this.view.fullHeight=e,this.view.offsetX=n,this.view.offsetY=i,this.view.width=s,this.view.height=r,this.updateProjectionMatrix()}clearViewOffset(){if(this.view!==null)this.view.enabled=!1;this.updateProjectionMatrix()}updateProjectionMatrix(){let t=this.near,e=t*Math.tan(Ni*0.5*this.fov)/this.zoom,n=2*e,i=this.aspect*n,s=-0.5*i,r=this.view;if(this.view!==null&&this.view.enabled){let{fullWidth:o,fullHeight:l}=r;s+=r.offsetX*i/o,e-=r.offsetY*n/l,i*=r.width/o,n*=r.height/l}let a=this.filmOffset;if(a!==0)s+=t*a/this.getFilmWidth();this.projectionMatrix.makePerspective(s,s+i,e,e-n,t,this.far,this.coordinateSystem,this.reversedDepth),this.projectionMatrixInverse.copy(this.projectionMatrix).invert()}toJSON(t){let e=super.toJSON(t);if(e.object.fov=this.fov,e.object.zoom=this.zoom,e.object.near=this.near,e.object.far=this.far,e.object.focus=this.focus,e.object.aspect=this.aspect,this.view!==null)e.object.view=Object.assign({},this.view);return e.object.filmGauge=this.filmGauge,e.object.filmOffset=this.filmOffset,e}}class kf extends Tr{constructor(){super(new De(50,1,0.5,500));this.isSpotLightShadow=!0,this.focus=1,this.aspect=1}updateMatrices(t){let e=this.camera,n=Di*2*t.angle*this.focus,i=this.mapSize.width/this.mapSize.height*this.aspect,s=t.distance||e.far;if(n!==e.fov||i!==e.aspect||s!==e.far)e.fov=n,e.aspect=i,e.far=s,e.updateProjectionMatrix();super.updateMatrices(t)}copy(t){return super.copy(t),this.focus=t.focus,this.aspect=t.aspect,this}toJSON(){let t=super.toJSON();return t.focus=this.focus,t.aspect=this.aspect,t}}class ih extends jn{constructor(t,e,n=0,i=Math.PI/3,s=0,r=2){super(t,e);this.isSpotLight=!0,this.type="SpotLight",this.position.copy(ae.DEFAULT_UP),this.updateMatrix(),this.target=new ae,this.distance=n,this.angle=i,this.penumbra=s,this.decay=r,this.map=null,this.shadow=new kf}get power(){return this.intensity*Math.PI}set power(t){this.intensity=t/Math.PI}dispose(){super.dispose(),this.shadow.dispose()}copy(t,e){return super.copy(t,e),this.distance=t.distance,this.angle=t.angle,this.penumbra=t.penumbra,this.decay=t.decay,this.target=t.target.clone(),this.map=t.map,this.shadow=t.shadow.clone(),this}toJSON(t){let e=super.toJSON(t);if(e.object.distance=this.distance,e.object.angle=this.angle,e.object.decay=this.decay,e.object.penumbra=this.penumbra,e.object.target=this.target.uuid,this.map&&this.map.isTexture)e.object.map=this.map.toJSON(t).uuid;return e.object.shadow=this.shadow.toJSON(),e}}class Gf extends Tr{constructor(){super(new De(90,1,0.5,500));this.isPointLightShadow=!0}}class Ar extends jn{constructor(t,e,n=0,i=2){super(t,e);this.isPointLight=!0,this.type="PointLight",this.distance=n,this.decay=i,this.shadow=new Gf}get power(){return this.intensity*4*Math.PI}set power(t){this.intensity=t/(4*Math.PI)}dispose(){super.dispose(),this.shadow.dispose()}copy(t,e){return super.copy(t,e),this.distance=t.distance,this.decay=t.decay,this.shadow=t.shadow.clone(),this}toJSON(t){let e=super.toJSON(t);return e.object.distance=this.distance,e.object.decay=this.decay,e.object.shadow=this.shadow.toJSON(),e}}class Ps extends Er{constructor(t=-1,e=1,n=1,i=-1,s=0.1,r=2000){super();this.isOrthographicCamera=!0,this.type="OrthographicCamera",this.zoom=1,this.view=null,this.left=t,this.right=e,this.top=n,this.bottom=i,this.near=s,this.far=r,this.updateProjectionMatrix()}copy(t,e){return super.copy(t,e),this.left=t.left,this.right=t.right,this.top=t.top,this.bottom=t.bottom,this.near=t.near,this.far=t.far,this.zoom=t.zoom,this.view=t.view===null?null:Object.assign({},t.view),this}setViewOffset(t,e,n,i,s,r){if(this.view===null)this.view={enabled:!0,fullWidth:1,fullHeight:1,offsetX:0,offsetY:0,width:1,height:1};this.view.enabled=!0,this.view.fullWidth=t,this.view.fullHeight=e,this.view.offsetX=n,this.view.offsetY=i,this.view.width=s,this.view.height=r,this.updateProjectionMatrix()}clearViewOffset(){if(this.view!==null)this.view.enabled=!1;this.updateProjectionMatrix()}updateProjectionMatrix(){let t=(this.right-this.left)/(2*this.zoom),e=(this.top-this.bottom)/(2*this.zoom),n=(this.right+this.left)/2,i=(this.top+this.bottom)/2,s=n-t,r=n+t,a=i+e,o=i-e;if(this.view!==null&&this.view.enabled){let l=(this.right-this.left)/this.view.fullWidth/this.zoom,c=(this.top-this.bottom)/this.view.fullHeight/this.zoom;s+=l*this.view.offsetX,r=s+l*this.view.width,a-=c*this.view.offsetY,o=a-c*this.view.height}this.projectionMatrix.makeOrthographic(s,r,a,o,this.near,this.far,this.coordinateSystem,this.reversedDepth),this.projectionMatrixInverse.copy(this.projectionMatrix).invert()}toJSON(t){let e=super.toJSON(t);if(e.object.zoom=this.zoom,e.object.left=this.left,e.object.right=this.right,e.object.top=this.top,e.object.bottom=this.bottom,e.object.near=this.near,e.object.far=this.far,this.view!==null)e.object.view=Object.assign({},this.view);return e}}class Hf extends Tr{constructor(){super(new Ps(-5,5,5,-5,0.5,500));this.isDirectionalLightShadow=!0}}class sh extends jn{constructor(t,e){super(t,e);this.isDirectionalLight=!0,this.type="DirectionalLight",this.position.copy(ae.DEFAULT_UP),this.updateMatrix(),this.target=new ae,this.shadow=new Hf}dispose(){super.dispose(),this.shadow.dispose()}copy(t){return super.copy(t),this.target=t.target.clone(),this.shadow=t.shadow.clone(),this}toJSON(t){let e=super.toJSON(t);return e.object.shadow=this.shadow.toJSON(),e.object.target=this.target.uuid,e}}class rh extends jn{constructor(t,e){super(t,e);this.isAmbientLight=!0,this.type="AmbientLight"}}class ah extends jn{constructor(t,e,n=10,i=10){super(t,e);this.isRectAreaLight=!0,this.type="RectAreaLight",this.width=n,this.height=i}get power(){return this.intensity*this.width*this.height*Math.PI}set power(t){this.intensity=t/(this.width*this.height*Math.PI)}copy(t){return super.copy(t),this.width=t.width,this.height=t.height,this}toJSON(t){let e=super.toJSON(t);return e.object.width=this.width,e.object.height=this.height,e}}class wo{constructor(){this.isSphericalHarmonics3=!0,this.coefficients=[];for(let t=0;t<9;t++)this.coefficients.push(new L)}set(t){for(let e=0;e<9;e++)this.coefficients[e].copy(t[e]);return this}zero(){for(let t=0;t<9;t++)this.coefficients[t].set(0,0,0);return this}getAt(t,e){let{x:n,y:i,z:s}=t,r=this.coefficients;return e.copy(r[0]).multiplyScalar(0.282095),e.addScaledVector(r[1],0.488603*i),e.addScaledVector(r[2],0.488603*s),e.addScaledVector(r[3],0.488603*n),e.addScaledVector(r[4],1.092548*(n*i)),e.addScaledVector(r[5],1.092548*(i*s)),e.addScaledVector(r[6],0.315392*(3*s*s-1)),e.addScaledVector(r[7],1.092548*(n*s)),e.addScaledVector(r[8],0.546274*(n*n-i*i)),e}getIrradianceAt(t,e){let{x:n,y:i,z:s}=t,r=this.coefficients;return e.copy(r[0]).multiplyScalar(0.886227),e.addScaledVector(r[1],1.023328*i),e.addScaledVector(r[2],1.023328*s),e.addScaledVector(r[3],1.023328*n),e.addScaledVector(r[4],0.858086*n*i),e.addScaledVector(r[5],0.858086*i*s),e.addScaledVector(r[6],0.743125*s*s-0.247708),e.addScaledVector(r[7],0.858086*n*s),e.addScaledVector(r[8],0.429043*(n*n-i*i)),e}add(t){for(let e=0;e<9;e++)this.coefficients[e].add(t.coefficients[e]);return this}addScaledSH(t,e){for(let n=0;n<9;n++)this.coefficients[n].addScaledVector(t.coefficients[n],e);return this}scale(t){for(let e=0;e<9;e++)this.coefficients[e].multiplyScalar(t);return this}lerp(t,e){for(let n=0;n<9;n++)this.coefficients[n].lerp(t.coefficients[n],e);return this}equals(t){for(let e=0;e<9;e++)if(!this.coefficients[e].equals(t.coefficients[e]))return!1;return!0}copy(t){return this.set(t.coefficients)}clone(){return new this.constructor().copy(this)}fromArray(t,e=0){let n=this.coefficients;for(let i=0;i<9;i++)n[i].fromArray(t,e+i*3);return this}toArray(t=[],e=0){let n=this.coefficients;for(let i=0;i<9;i++)n[i].toArray(t,e+i*3);return t}static getBasisAt(t,e){let{x:n,y:i,z:s}=t;e[0]=0.282095,e[1]=0.488603*i,e[2]=0.488603*s,e[3]=0.488603*n,e[4]=1.092548*n*i,e[5]=1.092548*i*s,e[6]=0.315392*(3*s*s-1),e[7]=1.092548*n*s,e[8]=0.546274*(n*n-i*i)}}class oh extends jn{constructor(t=new wo,e=1){super(void 0,e);this.isLightProbe=!0,this.sh=t}copy(t){return super.copy(t),this.sh.copy(t.sh),this}toJSON(t){let e=super.toJSON(t);return e.object.sh=this.sh.toArray(),e}}var zu={};class Ro extends Be{constructor(t){super(t);this.textures={}}load(t,e,n,i){let s=this,r=new fn(s.manager);r.setPath(s.path),r.setRequestHeader(s.requestHeader),r.setWithCredentials(s.withCredentials),r.load(t,function(a){try{e(s.parse(JSON.parse(a)))}catch(o){if(i)i(o);else Wt(o);s.manager.itemError(t)}},n,i)}parse(t){let e=this.createMaterialFromType(t.type);return e.fromJSON(t,this.textures),e}setTextures(t){return this.textures=t,this}createMaterialFromType(t){return Ro.createMaterialFromType(t)}static createMaterialFromType(t){let n={ShadowMaterial:Wc,SpriteMaterial:Ja,RawShaderMaterial:vo,ShaderMaterial:ln,PointsMaterial:Qa,MeshPhysicalMaterial:qc,MeshStandardMaterial:Wi,MeshPhongMaterial:Yc,MeshToonMaterial:Zc,MeshNormalMaterial:Jc,MeshLambertMaterial:Sr,MeshDepthMaterial:So,MeshDistanceMaterial:Mo,MeshBasicMaterial:mn,MeshMatcapMaterial:$c,LineDashedMaterial:Kc,LineBasicMaterial:We,Material:Fe,...zu}[t],i;if(n===void 0)qn(`MaterialLoader: Unknown material type "${t}". Use .registerMaterial() before starting the deserialization process.`),i=new Fe;else i=new n;return i}static registerMaterial(t,e){zu[t]=e}}class Ca{static extractUrlBase(t){let e=t.lastIndexOf("/");if(e===-1)return"./";return t.slice(0,e+1)}static resolveURL(t,e){if(typeof t!=="string"||t==="")return"";if(/^https?:\/\//i.test(e)&&/^\//.test(t))e=e.replace(/(^https?:\/\/[^\/]+).*/i,"$1");if(/^(https?:)?\/\//i.test(t))return t;if(/^data:.*,.*$/i.test(t))return t;if(/^blob:.*$/i.test(t))return t;return e+t}}class lh extends jt{constructor(){super();this.isInstancedBufferGeometry=!0,this.type="InstancedBufferGeometry",this.instanceCount=1/0}copy(t){return super.copy(t),this.instanceCount=t.instanceCount,this}toJSON(){let t=super.toJSON();return t.instanceCount=this.instanceCount,t.isInstancedBufferGeometry=!0,t}}class ch extends Be{constructor(t){super(t)}load(t,e,n,i){let s=this,r=new fn(s.manager);r.setPath(s.path),r.setRequestHeader(s.requestHeader),r.setWithCredentials(s.withCredentials),r.load(t,function(a){try{e(s.parse(JSON.parse(a)))}catch(o){if(i)i(o);else Wt(o);s.manager.itemError(t)}},n,i)}parse(t){let e={},n={};function i(u,f){if(e[f]!==void 0)return e[f];let _=u.interleavedBuffers[f],m=s(u,_.buffer),p=ms(_.type,m),b=new ws(p,_.stride);if(b.uuid=_.uuid,_.usage!==void 0)b.setUsage(_.usage);return e[f]=b,b}function s(u,f){if(n[f]!==void 0)return n[f];let _=u.arrayBuffers[f],m=new Uint32Array(_).buffer;return n[f]=m,m}let r=t.isInstancedBufferGeometry?new lh:new jt,a=t.data.index;if(a!==void 0){let u=ms(a.type,a.array);r.setIndex(new de(u,1))}let o=t.data.attributes;for(let u in o){let f=o[u],g;if(f.isInterleavedBufferAttribute){let _=i(t.data,f.data);g=new ci(_,f.itemSize,f.offset,f.normalized)}else{let _=ms(f.type,f.array);g=new(f.isInstancedBufferAttribute?hi:de)(_,f.itemSize,f.normalized)}if(f.name!==void 0)g.name=f.name;if(f.usage!==void 0)g.setUsage(f.usage);if(f.gpuType!==void 0)g.gpuType=f.gpuType;r.setAttribute(u,g)}let l=t.data.morphAttributes;if(l)for(let u in l){let f=l[u],g=[];for(let _=0,m=f.length;_<m;_++){let p=f[_],b;if(p.isInterleavedBufferAttribute){let S=i(t.data,p.data);b=new ci(S,p.itemSize,p.offset,p.normalized)}else{let S=ms(p.type,p.array);b=new de(S,p.itemSize,p.normalized)}if(p.name!==void 0)b.name=p.name;if(p.usage!==void 0)b.setUsage(p.usage);if(p.gpuType!==void 0)b.gpuType=p.gpuType;g.push(b)}r.morphAttributes[u]=g}if(t.data.morphTargetsRelative)r.morphTargetsRelative=!0;let h=t.data.groups||t.data.drawcalls||t.data.offsets;if(h!==void 0)for(let u=0,f=h.length;u!==f;++u){let g=h[u];r.addGroup(g.start,g.count,g.materialIndex)}let d=t.data.boundingSphere;if(d!==void 0)r.boundingSphere=new Ue().fromJSON(d);if(t.name)r.name=t.name;if(t.userData)r.userData=t.userData;return r}}var ml={};class Vf extends Be{constructor(t){super(t)}load(t,e,n,i){let s=this,r=this.path===""?Ca.extractUrlBase(t):this.path;this.resourcePath=this.resourcePath||r;let a=new fn(this.manager);a.setPath(this.path),a.setRequestHeader(this.requestHeader),a.setWithCredentials(this.withCredentials),a.load(t,function(o){let l=null;try{l=JSON.parse(o)}catch(h){if(i!==void 0)i(h);Wt("ObjectLoader: Can't parse "+t+".",h.message);return}let c=l.metadata;if(c===void 0||c.type===void 0||c.type.toLowerCase()==="geometry"){if(i!==void 0)i(Error("THREE.ObjectLoader: Can't load "+t));Wt("ObjectLoader: Can't load "+t);return}s.parse(l,e)},n,i)}async loadAsync(t,e){let n=this,i=this.path===""?Ca.extractUrlBase(t):this.path;this.resourcePath=this.resourcePath||i;let s=new fn(this.manager);s.setPath(this.path),s.setRequestHeader(this.requestHeader),s.setWithCredentials(this.withCredentials);let r=await s.loadAsync(t,e),a;try{a=JSON.parse(r)}catch(l){throw Error("THREE.ObjectLoader: Can't parse "+t+". "+l.message)}let o=a.metadata;if(o===void 0||o.type===void 0||o.type.toLowerCase()==="geometry")throw Error("THREE.ObjectLoader: Can't load "+t);return await n.parseAsync(a)}parse(t,e){let n=this.parseAnimations(t.animations),i=this.parseShapes(t.shapes),s=this.parseGeometries(t.geometries,i),r=this.parseImages(t.images,function(){if(e!==void 0)e(l)}),a=this.parseTextures(t.textures,r),o=this.parseMaterials(t.materials,a),l=this.parseObject(t.object,s,o,a,n),c=this.parseSkeletons(t.skeletons,l);if(this.bindSkeletons(l,c),this.bindLightTargets(l),e!==void 0){let h=!1;for(let d in r)if(r[d].data instanceof HTMLImageElement){h=!0;break}if(h===!1)e(l)}return l}async parseAsync(t){let e=this.parseAnimations(t.animations),n=this.parseShapes(t.shapes),i=this.parseGeometries(t.geometries,n),s=await this.parseImagesAsync(t.images),r=this.parseTextures(t.textures,s),a=this.parseMaterials(t.materials,r),o=this.parseObject(t.object,i,a,r,e),l=this.parseSkeletons(t.skeletons,o);return this.bindSkeletons(o,l),this.bindLightTargets(o),o}static registerGeometry(t,e){ml[t]=e}parseShapes(t){let e={};if(t!==void 0)for(let n=0,i=t.length;n<i;n++){let s=new Cs().fromJSON(t[n]);e[s.uuid]=s}return e}parseSkeletons(t,e){let n={},i={};if(e.traverse(function(s){if(s.isBone)i[s.uuid]=s}),t!==void 0)for(let s=0,r=t.length;s<r;s++){let a=new Ka().fromJSON(t[s],i);n[a.uuid]=a}return n}parseGeometries(t,e){let n={};if(t!==void 0){let i=new ch;for(let s=0,r=t.length;s<r;s++){let a,o=t[s];switch(o.type){case"BufferGeometry":case"InstancedBufferGeometry":a=i.parse(o);break;default:if(o.type in Pu)a=Pu[o.type].fromJSON(o,e);else if(o.type in ml)a=ml[o.type].fromJSON(o,e);else Mt(`ObjectLoader: Unknown geometry type "${o.type}". Use .registerGeometry() before starting the deserialization process.`)}if(a.uuid=o.uuid,o.name!==void 0)a.name=o.name;if(o.userData!==void 0)a.userData=o.userData;n[o.uuid]=a}}return n}parseMaterials(t,e){let n={},i={};if(t!==void 0){let s=new Ro;s.setTextures(e);for(let r=0,a=t.length;r<a;r++){let o=t[r];if(n[o.uuid]===void 0)n[o.uuid]=s.parse(o);i[o.uuid]=n[o.uuid]}}return i}parseAnimations(t){let e={};if(t!==void 0)for(let n=0;n<t.length;n++){let i=t[n],s=ys.parse(i);e[s.uuid]=s}return e}parseImages(t,e){let n=this,i={},s;function r(o){return o=n.manager.resolveURL(o),n.manager.itemStart(o),s.load(o,function(){n.manager.itemEnd(o)},void 0,function(){n.manager.itemError(o),n.manager.itemEnd(o)})}function a(o){if(typeof o==="string"){let l=o,c=/^(\/\/)|([a-z]+:(\/\/)?)/i.test(l)?l:n.resourcePath+l;return r(c)}else if(o.data)return{data:ms(o.type,o.data),width:o.width,height:o.height};else return null}if(t!==void 0&&t.length>0){let o=new Ao(e);s=new vs(o),s.setCrossOrigin(this.crossOrigin);for(let l=0,c=t.length;l<c;l++){let h=t[l],d=h.url;if(Array.isArray(d)){let u=[];for(let f=0,g=d.length;f<g;f++){let _=d[f],m=a(_);if(m!==null)if(m instanceof HTMLImageElement)u.push(m);else u.push(new an(m.data,m.width,m.height))}i[h.uuid]=new In(u)}else{let u=a(h.url);i[h.uuid]=new In(u)}}}return i}async parseImagesAsync(t){let e=this,n={},i;async function s(r){if(typeof r==="string"){let a=r,o=/^(\/\/)|([a-z]+:(\/\/)?)/i.test(a)?a:e.resourcePath+a;return await i.loadAsync(o)}else if(r.data)return{data:ms(r.type,r.data),width:r.width,height:r.height};else return null}if(t!==void 0&&t.length>0){i=new vs(this.manager),i.setCrossOrigin(this.crossOrigin);for(let r=0,a=t.length;r<a;r++){let o=t[r],l=o.url;if(Array.isArray(l)){let c=[];for(let h=0,d=l.length;h<d;h++){let u=l[h],f=await s(u);if(f!==null)if(f instanceof HTMLImageElement)c.push(f);else c.push(new an(f.data,f.width,f.height))}n[o.uuid]=new In(c)}else{let c=await s(o.url);n[o.uuid]=new In(c)}}}return n}parseTextures(t,e){function n(s,r){if(typeof s==="number")return s;return Mt("ObjectLoader.parseTexture: Constant should be in numeric form.",s),r[s]}let i={};if(t!==void 0)for(let s=0,r=t.length;s<r;s++){let a=t[s];if(a.image===void 0)Mt('ObjectLoader: No "image" specified for',a.uuid);if(e[a.image]===void 0)Mt("ObjectLoader: Undefined image",a.image);let o=e[a.image],l=o.data,c;if(Array.isArray(l)){if(c=new Rs,l.length===6)c.needsUpdate=!0}else{if(l&&l.data)c=new an;else c=new be;if(l)c.needsUpdate=!0}if(c.source=o,c.uuid=a.uuid,a.name!==void 0)c.name=a.name;if(a.mapping!==void 0)c.mapping=n(a.mapping,n_);if(a.channel!==void 0)c.channel=a.channel;if(a.offset!==void 0)c.offset.fromArray(a.offset);if(a.repeat!==void 0)c.repeat.fromArray(a.repeat);if(a.center!==void 0)c.center.fromArray(a.center);if(a.rotation!==void 0)c.rotation=a.rotation;if(a.wrap!==void 0)c.wrapS=n(a.wrap[0],ku),c.wrapT=n(a.wrap[1],ku);if(a.format!==void 0)c.format=a.format;if(a.internalFormat!==void 0)c.internalFormat=a.internalFormat;if(a.type!==void 0)c.type=a.type;if(a.colorSpace!==void 0)c.colorSpace=a.colorSpace;if(a.minFilter!==void 0)c.minFilter=n(a.minFilter,Gu);if(a.magFilter!==void 0)c.magFilter=n(a.magFilter,Gu);if(a.anisotropy!==void 0)c.anisotropy=a.anisotropy;if(a.flipY!==void 0)c.flipY=a.flipY;if(a.generateMipmaps!==void 0)c.generateMipmaps=a.generateMipmaps;if(a.premultiplyAlpha!==void 0)c.premultiplyAlpha=a.premultiplyAlpha;if(a.unpackAlignment!==void 0)c.unpackAlignment=a.unpackAlignment;if(a.compareFunction!==void 0)c.compareFunction=a.compareFunction;if(a.normalized!==void 0)c.normalized=a.normalized;if(a.userData!==void 0)c.userData=a.userData;i[a.uuid]=c}return i}parseObject(t,e,n,i,s){let r;function a(d){if(e[d]===void 0)Mt("ObjectLoader: Undefined geometry",d);return e[d]}function o(d){if(d===void 0)return;if(Array.isArray(d)){let u=[];for(let f=0,g=d.length;f<g;f++){let _=d[f];if(n[_]===void 0)Mt("ObjectLoader: Undefined material",_);u.push(n[_])}return u}if(n[d]===void 0)Mt("ObjectLoader: Undefined material",d);return n[d]}function l(d){if(i[d]===void 0)Mt("ObjectLoader: Undefined texture",d);return i[d]}let c,h;switch(t.type){case"Scene":if(r=new fr,t.background!==void 0)if(Number.isInteger(t.background))r.background=new Pt(t.background);else r.background=l(t.background);if(t.environment!==void 0)r.environment=l(t.environment);if(t.fog!==void 0){if(t.fog.type==="Fog")r.fog=new qa(t.fog.color,t.fog.near,t.fog.far);else if(t.fog.type==="FogExp2")r.fog=new Xa(t.fog.color,t.fog.density);if(t.fog.name!=="")r.fog.name=t.fog.name}if(t.backgroundBlurriness!==void 0)r.backgroundBlurriness=t.backgroundBlurriness;if(t.backgroundIntensity!==void 0)r.backgroundIntensity=t.backgroundIntensity;if(t.backgroundRotation!==void 0)r.backgroundRotation.fromArray(t.backgroundRotation);if(t.environmentIntensity!==void 0)r.environmentIntensity=t.environmentIntensity;if(t.environmentRotation!==void 0)r.environmentRotation.fromArray(t.environmentRotation);break;case"PerspectiveCamera":if(r=new De(t.fov,t.aspect,t.near,t.far),t.focus!==void 0)r.focus=t.focus;if(t.zoom!==void 0)r.zoom=t.zoom;if(t.filmGauge!==void 0)r.filmGauge=t.filmGauge;if(t.filmOffset!==void 0)r.filmOffset=t.filmOffset;if(t.view!==void 0)r.view=Object.assign({},t.view);break;case"OrthographicCamera":if(r=new Ps(t.left,t.right,t.top,t.bottom,t.near,t.far),t.zoom!==void 0)r.zoom=t.zoom;if(t.view!==void 0)r.view=Object.assign({},t.view);break;case"AmbientLight":r=new rh(t.color,t.intensity);break;case"DirectionalLight":r=new sh(t.color,t.intensity),r.target=t.target||"";break;case"PointLight":r=new Ar(t.color,t.intensity,t.distance,t.decay);break;case"RectAreaLight":r=new ah(t.color,t.intensity,t.width,t.height);break;case"SpotLight":r=new ih(t.color,t.intensity,t.distance,t.angle,t.penumbra,t.decay),r.target=t.target||"";break;case"HemisphereLight":r=new nh(t.color,t.groundColor,t.intensity);break;case"LightProbe":let d=new wo().fromArray(t.sh);r=new oh(d,t.intensity);break;case"SkinnedMesh":if(c=a(t.geometry),h=o(t.material),r=new Ic(c,h),t.bindMode!==void 0)r.bindMode=t.bindMode;if(t.bindMatrix!==void 0)r.bindMatrix.fromArray(t.bindMatrix);if(t.skeleton!==void 0)r.skeleton=t.skeleton;break;case"Mesh":c=a(t.geometry),h=o(t.material),r=new pe(c,h);break;case"InstancedMesh":c=a(t.geometry),h=o(t.material);let{count:u,instanceMatrix:f,instanceColor:g}=t;if(r=new pr(c,h,u),r.instanceMatrix=new hi(new Float32Array(f.array),16),g!==void 0)r.instanceColor=new hi(new Float32Array(g.array),g.itemSize);break;case"BatchedMesh":if(c=a(t.geometry),h=o(t.material),r=new Pc(t.maxInstanceCount,t.maxVertexCount,t.maxIndexCount,h),r.geometry=c,r.perObjectFrustumCulled=t.perObjectFrustumCulled,r.sortObjects=t.sortObjects,r._drawRanges=t.drawRanges,r._reservedRanges=t.reservedRanges,r._geometryInfo=t.geometryInfo.map((_)=>{let m=null,p=null;if(_.boundingBox!==void 0)m=new Oe().fromJSON(_.boundingBox);if(_.boundingSphere!==void 0)p=new Ue().fromJSON(_.boundingSphere);return{..._,boundingBox:m,boundingSphere:p}}),r._instanceInfo=t.instanceInfo,r._availableInstanceIds=t._availableInstanceIds,r._availableGeometryIds=t._availableGeometryIds,r._nextIndexStart=t.nextIndexStart,r._nextVertexStart=t.nextVertexStart,r._geometryCount=t.geometryCount,r._maxInstanceCount=t.maxInstanceCount,r._maxVertexCount=t.maxVertexCount,r._maxIndexCount=t.maxIndexCount,r._geometryInitialized=t.geometryInitialized,r._matricesTexture=l(t.matricesTexture.uuid),r._indirectTexture=l(t.indirectTexture.uuid),t.colorsTexture!==void 0)r._colorsTexture=l(t.colorsTexture.uuid);if(t.boundingSphere!==void 0)r.boundingSphere=new Ue().fromJSON(t.boundingSphere);if(t.boundingBox!==void 0)r.boundingBox=new Oe().fromJSON(t.boundingBox);break;case"LOD":r=new Cc;break;case"Line":r=new Zn(a(t.geometry),o(t.material));break;case"LineLoop":r=new Lc(a(t.geometry),o(t.material));break;case"LineSegments":r=new En(a(t.geometry),o(t.material));break;case"PointCloud":case"Points":r=new Nc(a(t.geometry),o(t.material));break;case"Sprite":r=new Rc(o(t.material));break;case"Group":r=new Li;break;case"Bone":r=new $a;break;default:r=new ae}if(r.uuid=t.uuid,t.name!==void 0)r.name=t.name;if(t.matrix!==void 0){if(r.matrix.fromArray(t.matrix),t.matrixAutoUpdate!==void 0)r.matrixAutoUpdate=t.matrixAutoUpdate;if(r.matrixAutoUpdate)r.matrix.decompose(r.position,r.quaternion,r.scale)}else{if(t.position!==void 0)r.position.fromArray(t.position);if(t.rotation!==void 0)r.rotation.fromArray(t.rotation);if(t.quaternion!==void 0)r.quaternion.fromArray(t.quaternion);if(t.scale!==void 0)r.scale.fromArray(t.scale)}if(t.up!==void 0)r.up.fromArray(t.up);if(t.pivot!==void 0)r.pivot=new L().fromArray(t.pivot);if(t.morphTargetDictionary!==void 0)r.morphTargetDictionary=Object.assign({},t.morphTargetDictionary);if(t.morphTargetInfluences!==void 0)r.morphTargetInfluences=t.morphTargetInfluences.slice();if(t.castShadow!==void 0)r.castShadow=t.castShadow;if(t.receiveShadow!==void 0)r.receiveShadow=t.receiveShadow;if(t.shadow){if(t.shadow.intensity!==void 0)r.shadow.intensity=t.shadow.intensity;if(t.shadow.bias!==void 0)r.shadow.bias=t.shadow.bias;if(t.shadow.normalBias!==void 0)r.shadow.normalBias=t.shadow.normalBias;if(t.shadow.radius!==void 0)r.shadow.radius=t.shadow.radius;if(t.shadow.blurSamples!==void 0)r.shadow.blurSamples=t.shadow.blurSamples;if(t.shadow.focus!==void 0)r.shadow.focus=t.shadow.focus;if(t.shadow.aspect!==void 0)r.shadow.aspect=t.shadow.aspect;if(t.shadow.mapSize!==void 0)r.shadow.mapSize.fromArray(t.shadow.mapSize);if(t.shadow.camera!==void 0)r.shadow.camera=this.parseObject(t.shadow.camera)}if(t.visible!==void 0)r.visible=t.visible;if(t.frustumCulled!==void 0)r.frustumCulled=t.frustumCulled;if(t.renderOrder!==void 0)r.renderOrder=t.renderOrder;if(t.static!==void 0)r.static=t.static;if(t.userData!==void 0)r.userData=t.userData;if(t.layers!==void 0)r.layers.mask=t.layers;if(t.children!==void 0){let d=t.children;for(let u=0;u<d.length;u++)r.add(this.parseObject(d[u],e,n,i,s))}if(t.animations!==void 0){let d=t.animations;for(let u=0;u<d.length;u++){let f=d[u];r.animations.push(s[f])}}if(t.type==="LOD"){if(t.autoUpdate!==void 0)r.autoUpdate=t.autoUpdate;let d=t.levels;for(let u=0;u<d.length;u++){let f=d[u],g=r.getObjectByProperty("uuid",f.object);if(g!==void 0)r.addLevel(g,f.distance,f.hysteresis)}}return r}bindSkeletons(t,e){if(Object.keys(e).length===0)return;t.traverse(function(n){if(n.isSkinnedMesh===!0&&n.skeleton!==void 0){let i=e[n.skeleton];if(i===void 0)Mt("ObjectLoader: No skeleton found with UUID:",n.skeleton);else n.bind(i,n.bindMatrix)}})}bindLightTargets(t){t.traverse(function(e){if(e.isDirectionalLight||e.isSpotLight){let n=e.target,i=t.getObjectByProperty("uuid",n);if(i!==void 0)e.target=i;else e.target=new ae}})}}var n_={UVMapping:300,CubeReflectionMapping:301,CubeRefractionMapping:302,EquirectangularReflectionMapping:303,EquirectangularRefractionMapping:304,CubeUVReflectionMapping:306},ku={RepeatWrapping:1000,ClampToEdgeWrapping:1001,MirroredRepeatWrapping:1002},Gu={NearestFilter:1003,NearestMipmapNearestFilter:1004,NearestMipmapLinearFilter:1005,LinearFilter:1006,LinearMipmapNearestFilter:1007,LinearMipmapLinearFilter:1008},gl=new WeakMap;class Wf extends Be{constructor(t){super(t);if(this.isImageBitmapLoader=!0,typeof createImageBitmap>"u")Mt("ImageBitmapLoader: createImageBitmap() not supported.");if(typeof fetch>"u")Mt("ImageBitmapLoader: fetch() not supported.");this.options={premultiplyAlpha:"none"},this._abortController=new AbortController}setOptions(t){return this.options=t,this}load(t,e,n,i){if(t===void 0)t="";if(this.path!==void 0)t=this.path+t;t=this.manager.resolveURL(t);let s=this,r=Pn.get(`image-bitmap:${t}`);if(r!==void 0){if(s.manager.itemStart(t),r.then){r.then((l)=>{if(gl.has(r)===!0){if(i)i(gl.get(r));s.manager.itemError(t),s.manager.itemEnd(t)}else{if(e)e(l);s.manager.itemEnd(t)}});return}setTimeout(function(){if(e)e(r);s.manager.itemEnd(t)},0);return}let a={};a.credentials=this.crossOrigin==="anonymous"?"same-origin":"include",a.headers=this.requestHeader,a.signal=typeof AbortSignal.any==="function"?AbortSignal.any([this._abortController.signal,this.manager.abortController.signal]):this._abortController.signal;let o=fetch(t,a).then(function(l){return l.blob()}).then(function(l){return createImageBitmap(l,Object.assign({},s.options,{colorSpaceConversion:"none"}))}).then(function(l){if(Pn.add(`image-bitmap:${t}`,l),e)e(l);return s.manager.itemEnd(t),l}).catch(function(l){if(i)i(l);gl.set(o,l),Pn.remove(`image-bitmap:${t}`),s.manager.itemError(t),s.manager.itemEnd(t)});Pn.add(`image-bitmap:${t}`,o),s.manager.itemStart(t)}abort(){return this._abortController.abort(),this._abortController=new AbortController,this}}var ya;class Co{static getContext(){if(ya===void 0)ya=new(window.AudioContext||window.webkitAudioContext);return ya}static setContext(t){ya=t}}class Xf extends Be{constructor(t){super(t)}load(t,e,n,i){let s=this,r=new fn(this.manager);r.setResponseType("arraybuffer"),r.setPath(this.path),r.setRequestHeader(this.requestHeader),r.setWithCredentials(this.withCredentials),r.load(t,function(o){try{let l=o.slice(0),c=Co.getContext(),h=t+"#decode";s.manager.itemStart(h),c.decodeAudioData(l,function(d){e(d),s.manager.itemEnd(h)}).catch(function(d){a(d),s.manager.itemEnd(h)})}catch(l){a(l)}},n,i);function a(o){if(i)i(o);else Wt(o);s.manager.itemError(t)}}}var Hu=new Kt,Vu=new Kt,Ai=new Kt;class qf{constructor(){this.type="StereoCamera",this.aspect=1,this.eyeSep=0.064,this.cameraL=new De,this.cameraL.layers.enable(1),this.cameraL.matrixAutoUpdate=!1,this.cameraR=new De,this.cameraR.layers.enable(2),this.cameraR.matrixAutoUpdate=!1,this._cache={focus:null,fov:null,aspect:null,near:null,far:null,zoom:null,eyeSep:null}}update(t){let e=this._cache;if(e.focus!==t.focus||e.fov!==t.fov||e.aspect!==t.aspect*this.aspect||e.near!==t.near||e.far!==t.far||e.zoom!==t.zoom||e.eyeSep!==this.eyeSep){e.focus=t.focus,e.fov=t.fov,e.aspect=t.aspect*this.aspect,e.near=t.near,e.far=t.far,e.zoom=t.zoom,e.eyeSep=this.eyeSep,Ai.copy(t.projectionMatrix);let i=e.eyeSep/2,s=i*e.near/e.focus,r=e.near*Math.tan(Ni*e.fov*0.5)/e.zoom,a,o;Vu.elements[12]=-i,Hu.elements[12]=i,a=-r*e.aspect+s,o=r*e.aspect+s,Ai.elements[0]=2*e.near/(o-a),Ai.elements[8]=(o+a)/(o-a),this.cameraL.projectionMatrix.copy(Ai),a=-r*e.aspect-s,o=r*e.aspect-s,Ai.elements[0]=2*e.near/(o-a),Ai.elements[8]=(o+a)/(o-a),this.cameraR.projectionMatrix.copy(Ai)}this.cameraL.matrix.copy(t.matrixWorld).multiply(Vu),this.cameraL.matrixWorldNeedsUpdate=!0,this.cameraR.matrix.copy(t.matrixWorld).multiply(Hu),this.cameraR.matrixWorldNeedsUpdate=!0}}var us=-90,ds=1;class hh extends ae{constructor(t,e,n){super();this.type="CubeCamera",this.renderTarget=n,this.coordinateSystem=null,this.activeMipmapLevel=0;let i=new De(us,ds,t,e);i.layers=this.layers,this.add(i);let s=new De(us,ds,t,e);s.layers=this.layers,this.add(s);let r=new De(us,ds,t,e);r.layers=this.layers,this.add(r);let a=new De(us,ds,t,e);a.layers=this.layers,this.add(a);let o=new De(us,ds,t,e);o.layers=this.layers,this.add(o);let l=new De(us,ds,t,e);l.layers=this.layers,this.add(l)}updateCoordinateSystem(){let t=this.coordinateSystem,e=this.children.concat(),[n,i,s,r,a,o]=e;for(let l of e)this.remove(l);if(t===2000)n.up.set(0,1,0),n.lookAt(1,0,0),i.up.set(0,1,0),i.lookAt(-1,0,0),s.up.set(0,0,-1),s.lookAt(0,1,0),r.up.set(0,0,1),r.lookAt(0,-1,0),a.up.set(0,1,0),a.lookAt(0,0,1),o.up.set(0,1,0),o.lookAt(0,0,-1);else if(t===2001)n.up.set(0,-1,0),n.lookAt(-1,0,0),i.up.set(0,-1,0),i.lookAt(1,0,0),s.up.set(0,0,1),s.lookAt(0,1,0),r.up.set(0,0,-1),r.lookAt(0,-1,0),a.up.set(0,-1,0),a.lookAt(0,0,1),o.up.set(0,-1,0),o.lookAt(0,0,-1);else throw Error("THREE.CubeCamera.updateCoordinateSystem(): Invalid coordinate system: "+t);for(let l of e)this.add(l),l.updateMatrixWorld()}update(t,e){if(this.parent===null)this.updateMatrixWorld();let{renderTarget:n,activeMipmapLevel:i}=this;if(this.coordinateSystem!==t.coordinateSystem)this.coordinateSystem=t.coordinateSystem,this.updateCoordinateSystem();let[s,r,a,o,l,c]=this.children,h=t.getRenderTarget(),d=t.getActiveCubeFace(),u=t.getActiveMipmapLevel(),f=t.xr.enabled;t.xr.enabled=!1;let g=n.texture.generateMipmaps;n.texture.generateMipmaps=!1;let _=!1;if(t.isWebGLRenderer===!0)_=t.state.buffers.depth.getReversed();else _=t.reversedDepthBuffer;if(t.setRenderTarget(n,0,i),_&&t.autoClear===!1)t.clearDepth();if(t.render(e,s),t.setRenderTarget(n,1,i),_&&t.autoClear===!1)t.clearDepth();if(t.render(e,r),t.setRenderTarget(n,2,i),_&&t.autoClear===!1)t.clearDepth();if(t.render(e,a),t.setRenderTarget(n,3,i),_&&t.autoClear===!1)t.clearDepth();if(t.render(e,o),t.setRenderTarget(n,4,i),_&&t.autoClear===!1)t.clearDepth();if(t.render(e,l),n.texture.generateMipmaps=g,t.setRenderTarget(n,5,i),_&&t.autoClear===!1)t.clearDepth();t.render(e,c),t.setRenderTarget(h,d,u),t.xr.enabled=f,n.texture.needsPMREMUpdate=!0}}class uh extends De{constructor(t=[]){super();this.isArrayCamera=!0,this.isMultiViewCamera=!1,this.cameras=t}}class dh{constructor(){this._previousTime=0,this._currentTime=0,this._startTime=performance.now(),this._delta=0,this._elapsed=0,this._timescale=1,this._document=null,this._pageVisibilityHandler=null}connect(t){if(this._document=t,t.hidden!==void 0)this._pageVisibilityHandler=i_.bind(this),t.addEventListener("visibilitychange",this._pageVisibilityHandler,!1)}disconnect(){if(this._pageVisibilityHandler!==null)this._document.removeEventListener("visibilitychange",this._pageVisibilityHandler),this._pageVisibilityHandler=null;this._document=null}getDelta(){return this._delta/1000}getElapsed(){return this._elapsed/1000}getTimescale(){return this._timescale}setTimescale(t){return this._timescale=t,this}reset(){return this._currentTime=performance.now()-this._startTime,this}dispose(){this.disconnect()}update(t){if(this._pageVisibilityHandler!==null&&this._document.hidden===!0)this._delta=0;else this._previousTime=this._currentTime,this._currentTime=(t!==void 0?t:performance.now())-this._startTime,this._delta=(this._currentTime-this._previousTime)*this._timescale,this._elapsed+=this._delta;return this}}function i_(){if(this._document.hidden===!1)this.reset()}var wi=new L,_l=new Ce,s_=new L,Ri=new L,Ci=new L;class Yf extends ae{constructor(){super();this.type="AudioListener",this.context=Co.getContext(),this.gain=this.context.createGain(),this.gain.connect(this.context.destination),this.filter=null,this.timeDelta=0,this._timer=new dh}getInput(){return this.gain}removeFilter(){if(this.filter!==null)this.gain.disconnect(this.filter),this.filter.disconnect(this.context.destination),this.gain.connect(this.context.destination),this.filter=null;return this}getFilter(){return this.filter}setFilter(t){if(this.filter!==null)this.gain.disconnect(this.filter),this.filter.disconnect(this.context.destination);else this.gain.disconnect(this.context.destination);return this.filter=t,this.gain.connect(this.filter),this.filter.connect(this.context.destination),this}getMasterVolume(){return this.gain.gain.value}setMasterVolume(t){return this.gain.gain.setTargetAtTime(t,this.context.currentTime,0.01),this}updateMatrixWorld(t){super.updateMatrixWorld(t),this._timer.update();let e=this.context.listener;if(this.timeDelta=this._timer.getDelta(),this.matrixWorld.decompose(wi,_l,s_),Ri.set(0,0,-1).applyQuaternion(_l),Ci.set(0,1,0).applyQuaternion(_l),e.positionX){let n=this.context.currentTime+this.timeDelta;e.positionX.linearRampToValueAtTime(wi.x,n),e.positionY.linearRampToValueAtTime(wi.y,n),e.positionZ.linearRampToValueAtTime(wi.z,n),e.forwardX.linearRampToValueAtTime(Ri.x,n),e.forwardY.linearRampToValueAtTime(Ri.y,n),e.forwardZ.linearRampToValueAtTime(Ri.z,n),e.upX.linearRampToValueAtTime(Ci.x,n),e.upY.linearRampToValueAtTime(Ci.y,n),e.upZ.linearRampToValueAtTime(Ci.z,n)}else e.setPosition(wi.x,wi.y,wi.z),e.setOrientation(Ri.x,Ri.y,Ri.z,Ci.x,Ci.y,Ci.z)}}class fh extends ae{constructor(t){super();this.type="Audio",this.listener=t,this.context=t.context,this.gain=this.context.createGain(),this.gain.connect(t.getInput()),this.autoplay=!1,this.buffer=null,this.detune=0,this.loop=!1,this.loopStart=0,this.loopEnd=0,this.offset=0,this.duration=void 0,this.playbackRate=1,this.isPlaying=!1,this.hasPlaybackControl=!0,this.source=null,this.sourceType="empty",this._startedAt=0,this._progress=0,this._connected=!1,this.filters=[]}getOutput(){return this.gain}setNodeSource(t){return this.hasPlaybackControl=!1,this.sourceType="audioNode",this.source=t,this.connect(),this}setMediaElementSource(t){return this.hasPlaybackControl=!1,this.sourceType="mediaNode",this.source=this.context.createMediaElementSource(t),this.connect(),this}setMediaStreamSource(t){return this.hasPlaybackControl=!1,this.sourceType="mediaStreamNode",this.source=this.context.createMediaStreamSource(t),this.connect(),this}setBuffer(t){if(this.buffer=t,this.sourceType="buffer",this.autoplay)this.play();return this}play(t=0){if(this.isPlaying===!0){Mt("Audio: Audio is already playing.");return}if(this.hasPlaybackControl===!1){Mt("Audio: this Audio has no playback control.");return}this._startedAt=this.context.currentTime+t;let e=this.context.createBufferSource();return e.buffer=this.buffer,e.loop=this.loop,e.loopStart=this.loopStart,e.loopEnd=this.loopEnd,e.onended=this.onEnded.bind(this),e.start(this._startedAt,this._progress+this.offset,this.duration),this.isPlaying=!0,this.source=e,this.setDetune(this.detune),this.setPlaybackRate(this.playbackRate),this.connect()}pause(){if(this.hasPlaybackControl===!1){Mt("Audio: this Audio has no playback control.");return}if(this.isPlaying===!0){if(this._progress+=Math.max(this.context.currentTime-this._startedAt,0)*this.playbackRate,this.loop===!0)this._progress=this._progress%(this.duration||this.buffer.duration);this.source.stop(),this.source.onended=null,this.isPlaying=!1}return this}stop(t=0){if(this.hasPlaybackControl===!1){Mt("Audio: this Audio has no playback control.");return}if(this._progress=0,this.source!==null)this.source.stop(this.context.currentTime+t),this.source.onended=null;return this.isPlaying=!1,this}connect(){if(this.filters.length>0){this.source.connect(this.filters[0]);for(let t=1,e=this.filters.length;t<e;t++)this.filters[t-1].connect(this.filters[t]);this.filters[this.filters.length-1].connect(this.getOutput())}else this.source.connect(this.getOutput());return this._connected=!0,this}disconnect(){if(this._connected===!1)return;if(this.filters.length>0){this.source.disconnect(this.filters[0]);for(let t=1,e=this.filters.length;t<e;t++)this.filters[t-1].disconnect(this.filters[t]);this.filters[this.filters.length-1].disconnect(this.getOutput())}else this.source.disconnect(this.getOutput());return this._connected=!1,this}getFilters(){return this.filters}setFilters(t){if(!t)t=[];if(this._connected===!0)this.disconnect(),this.filters=t.slice(),this.connect();else this.filters=t.slice();return this}setDetune(t){if(this.detune=t,this.isPlaying===!0&&this.source.detune!==void 0)this.source.detune.setTargetAtTime(this.detune,this.context.currentTime,0.01);return this}getDetune(){return this.detune}getFilter(){return this.getFilters()[0]}setFilter(t){return this.setFilters(t?[t]:[])}setPlaybackRate(t){if(this.hasPlaybackControl===!1){Mt("Audio: this Audio has no playback control.");return}if(this.playbackRate=t,this.isPlaying===!0)this.source.playbackRate.setTargetAtTime(this.playbackRate,this.context.currentTime,0.01);return this}getPlaybackRate(){return this.playbackRate}onEnded(){this.isPlaying=!1,this._progress=0}getLoop(){if(this.hasPlaybackControl===!1)return Mt("Audio: this Audio has no playback control."),!1;return this.loop}setLoop(t){if(this.hasPlaybackControl===!1){Mt("Audio: this Audio has no playback control.");return}if(this.loop=t,this.isPlaying===!0)this.source.loop=this.loop;return this}setLoopStart(t){return this.loopStart=t,this}setLoopEnd(t){return this.loopEnd=t,this}getVolume(){return this.gain.gain.value}setVolume(t){return this.gain.gain.setTargetAtTime(t,this.context.currentTime,0.01),this}copy(t,e){if(super.copy(t,e),t.sourceType!=="buffer")return Mt("Audio: Audio source type cannot be copied."),this;return this.autoplay=t.autoplay,this.buffer=t.buffer,this.detune=t.detune,this.loop=t.loop,this.loopStart=t.loopStart,this.loopEnd=t.loopEnd,this.offset=t.offset,this.duration=t.duration,this.playbackRate=t.playbackRate,this.hasPlaybackControl=t.hasPlaybackControl,this.sourceType=t.sourceType,this.filters=t.filters.slice(),this}clone(t){return new this.constructor(this.listener).copy(this,t)}}var Ii=new L,Wu=new Ce,r_=new L,Pi=new L;class Zf extends fh{constructor(t){super(t);this.panner=this.context.createPanner(),this.panner.panningModel="HRTF",this.panner.connect(this.gain)}connect(){return super.connect(),this.panner.connect(this.gain),this}disconnect(){return super.disconnect(),this.panner.disconnect(this.gain),this}getOutput(){return this.panner}getRefDistance(){return this.panner.refDistance}setRefDistance(t){return this.panner.refDistance=t,this}getRolloffFactor(){return this.panner.rolloffFactor}setRolloffFactor(t){return this.panner.rolloffFactor=t,this}getDistanceModel(){return this.panner.distanceModel}setDistanceModel(t){return this.panner.distanceModel=t,this}getMaxDistance(){return this.panner.maxDistance}setMaxDistance(t){return this.panner.maxDistance=t,this}setDirectionalCone(t,e,n){return this.panner.coneInnerAngle=t,this.panner.coneOuterAngle=e,this.panner.coneOuterGain=n,this}updateMatrixWorld(t){if(super.updateMatrixWorld(t),this.hasPlaybackControl===!0&&this.isPlaying===!1)return;this.matrixWorld.decompose(Ii,Wu,r_),Pi.set(0,0,1).applyQuaternion(Wu);let e=this.panner;if(e.positionX){let n=this.context.currentTime+this.listener.timeDelta;e.positionX.linearRampToValueAtTime(Ii.x,n),e.positionY.linearRampToValueAtTime(Ii.y,n),e.positionZ.linearRampToValueAtTime(Ii.z,n),e.orientationX.linearRampToValueAtTime(Pi.x,n),e.orientationY.linearRampToValueAtTime(Pi.y,n),e.orientationZ.linearRampToValueAtTime(Pi.z,n)}else e.setPosition(Ii.x,Ii.y,Ii.z),e.setOrientation(Pi.x,Pi.y,Pi.z)}}class Jf{constructor(t,e=2048){this.analyser=t.context.createAnalyser(),this.analyser.fftSize=e,this.data=new Uint8Array(this.analyser.frequencyBinCount),t.getOutput().connect(this.analyser)}getFrequencyData(){return this.analyser.getByteFrequencyData(this.data),this.data}getAverageFrequency(){let t=0,e=this.getFrequencyData();for(let n=0;n<e.length;n++)t+=e[n];return t/e.length}}class ph{constructor(t,e,n){this.binding=t,this.valueSize=n;let i,s,r;switch(e){case"quaternion":i=this._slerp,s=this._slerpAdditive,r=this._setAdditiveIdentityQuaternion,this.buffer=new Float64Array(n*6),this._workIndex=5;break;case"string":case"bool":i=this._select,s=this._select,r=this._setAdditiveIdentityOther,this.buffer=Array(n*5);break;default:i=this._lerp,s=this._lerpAdditive,r=this._setAdditiveIdentityNumeric,this.buffer=new Float64Array(n*5)}this._mixBufferRegion=i,this._mixBufferRegionAdditive=s,this._setIdentity=r,this._origIndex=3,this._addIndex=4,this.cumulativeWeight=0,this.cumulativeWeightAdditive=0,this.useCount=0,this.referenceCount=0}accumulate(t,e){let n=this.buffer,i=this.valueSize,s=t*i+i,r=this.cumulativeWeight;if(r===0){for(let a=0;a!==i;++a)n[s+a]=n[a];r=e}else{r+=e;let a=e/r;this._mixBufferRegion(n,s,0,a,i)}this.cumulativeWeight=r}accumulateAdditive(t){let e=this.buffer,n=this.valueSize,i=n*this._addIndex;if(this.cumulativeWeightAdditive===0)this._setIdentity();this._mixBufferRegionAdditive(e,i,0,t,n),this.cumulativeWeightAdditive+=t}apply(t){let e=this.valueSize,n=this.buffer,i=t*e+e,s=this.cumulativeWeight,r=this.cumulativeWeightAdditive,a=this.binding;if(this.cumulativeWeight=0,this.cumulativeWeightAdditive=0,s<1){let o=e*this._origIndex;this._mixBufferRegion(n,i,o,1-s,e)}if(r>0)this._mixBufferRegionAdditive(n,i,this._addIndex*e,1,e);for(let o=e,l=e+e;o!==l;++o)if(n[o]!==n[o+e]){a.setValue(n,i);break}}saveOriginalState(){let t=this.binding,e=this.buffer,n=this.valueSize,i=n*this._origIndex;t.getValue(e,i);for(let s=n,r=i;s!==r;++s)e[s]=e[i+s%n];this._setIdentity(),this.cumulativeWeight=0,this.cumulativeWeightAdditive=0}restoreOriginalState(){let t=this.valueSize*3;this.binding.setValue(this.buffer,t)}_setAdditiveIdentityNumeric(){let t=this._addIndex*this.valueSize,e=t+this.valueSize;for(let n=t;n<e;n++)this.buffer[n]=0}_setAdditiveIdentityQuaternion(){this._setAdditiveIdentityNumeric(),this.buffer[this._addIndex*this.valueSize+3]=1}_setAdditiveIdentityOther(){let t=this._origIndex*this.valueSize,e=this._addIndex*this.valueSize;for(let n=0;n<this.valueSize;n++)this.buffer[e+n]=this.buffer[t+n]}_select(t,e,n,i,s){if(i>=0.5)for(let r=0;r!==s;++r)t[e+r]=t[n+r]}_slerp(t,e,n,i){Ce.slerpFlat(t,e,t,e,t,n,i)}_slerpAdditive(t,e,n,i,s){let r=this._workIndex*s;Ce.multiplyQuaternionsFlat(t,r,t,e,t,n),Ce.slerpFlat(t,e,t,e,t,r,i)}_lerp(t,e,n,i,s){let r=1-i;for(let a=0;a!==s;++a){let o=e+a;t[o]=t[o]*r+t[n+a]*i}}_lerpAdditive(t,e,n,i,s){for(let r=0;r!==s;++r){let a=e+r;t[a]=t[a]+t[n+r]*i}}}var mh="\\[\\]\\.:\\/",a_=new RegExp("["+mh+"]","g"),gh="[^"+mh+"]",o_="[^"+mh.replace("\\.","")+"]",l_=/((?:WC+[\/:])*)/.source.replace("WC",gh),c_=/(WCOD+)?/.source.replace("WCOD",o_),h_=/(?:\.(WC+)(?:\[(.+)\])?)?/.source.replace("WC",gh),u_=/\.(WC+)(?:\[(.+)\])?/.source.replace("WC",gh),d_=new RegExp("^"+l_+c_+h_+u_+"$"),f_=["material","materials","bones","map"];class $f{constructor(t,e,n){let i=n||le.parseTrackName(e);this._targetGroup=t,this._bindings=t.subscribe_(e,i)}getValue(t,e){this.bind();let n=this._targetGroup.nCachedObjects_,i=this._bindings[n];if(i!==void 0)i.getValue(t,e)}setValue(t,e){let n=this._bindings;for(let i=this._targetGroup.nCachedObjects_,s=n.length;i!==s;++i)n[i].setValue(t,e)}bind(){let t=this._bindings;for(let e=this._targetGroup.nCachedObjects_,n=t.length;e!==n;++e)t[e].bind()}unbind(){let t=this._bindings;for(let e=this._targetGroup.nCachedObjects_,n=t.length;e!==n;++e)t[e].unbind()}}class le{constructor(t,e,n){this.path=e,this.parsedPath=n||le.parseTrackName(e),this.node=le.findNode(t,this.parsedPath.nodeName),this.rootNode=t,this.getValue=this._getValue_unbound,this.setValue=this._setValue_unbound}static create(t,e,n){if(!(t&&t.isAnimationObjectGroup))return new le(t,e,n);else return new le.Composite(t,e,n)}static sanitizeNodeName(t){return t.replace(/\s/g,"_").replace(a_,"")}static parseTrackName(t){let e=d_.exec(t);if(e===null)throw Error("THREE.PropertyBinding: Cannot parse trackName: "+t);let n={nodeName:e[2],objectName:e[3],objectIndex:e[4],propertyName:e[5],propertyIndex:e[6]},i=n.nodeName&&n.nodeName.lastIndexOf(".");if(i!==void 0&&i!==-1){let s=n.nodeName.substring(i+1);if(f_.indexOf(s)!==-1)n.nodeName=n.nodeName.substring(0,i),n.objectName=s}if(n.propertyName===null||n.propertyName.length===0)throw Error("THREE.PropertyBinding: can not parse propertyName from trackName: "+t);return n}static findNode(t,e){if(e===void 0||e===""||e==="."||e===-1||e===t.name||e===t.uuid)return t;if(t.skeleton){let n=t.skeleton.getBoneByName(e);if(n!==void 0)return n}if(t.children){let n=function(s){for(let r=0;r<s.length;r++){let a=s[r];if(a.name===e||a.uuid===e)return a;let o=n(a.children);if(o)return o}return null},i=n(t.children);if(i)return i}return null}_getValue_unavailable(){}_setValue_unavailable(){}_getValue_direct(t,e){t[e]=this.targetObject[this.propertyName]}_getValue_array(t,e){let n=this.resolvedProperty;for(let i=0,s=n.length;i!==s;++i)t[e++]=n[i]}_getValue_arrayElement(t,e){t[e]=this.resolvedProperty[this.propertyIndex]}_getValue_toArray(t,e){this.resolvedProperty.toArray(t,e)}_setValue_direct(t,e){this.targetObject[this.propertyName]=t[e]}_setValue_direct_setNeedsUpdate(t,e){this.targetObject[this.propertyName]=t[e],this.targetObject.needsUpdate=!0}_setValue_direct_setMatrixWorldNeedsUpdate(t,e){this.targetObject[this.propertyName]=t[e],this.targetObject.matrixWorldNeedsUpdate=!0}_setValue_array(t,e){let n=this.resolvedProperty;for(let i=0,s=n.length;i!==s;++i)n[i]=t[e++]}_setValue_array_setNeedsUpdate(t,e){let n=this.resolvedProperty;for(let i=0,s=n.length;i!==s;++i)n[i]=t[e++];this.targetObject.needsUpdate=!0}_setValue_array_setMatrixWorldNeedsUpdate(t,e){let n=this.resolvedProperty;for(let i=0,s=n.length;i!==s;++i)n[i]=t[e++];this.targetObject.matrixWorldNeedsUpdate=!0}_setValue_arrayElement(t,e){this.resolvedProperty[this.propertyIndex]=t[e]}_setValue_arrayElement_setNeedsUpdate(t,e){this.resolvedProperty[this.propertyIndex]=t[e],this.targetObject.needsUpdate=!0}_setValue_arrayElement_setMatrixWorldNeedsUpdate(t,e){this.resolvedProperty[this.propertyIndex]=t[e],this.targetObject.matrixWorldNeedsUpdate=!0}_setValue_fromArray(t,e){this.resolvedProperty.fromArray(t,e)}_setValue_fromArray_setNeedsUpdate(t,e){this.resolvedProperty.fromArray(t,e),this.targetObject.needsUpdate=!0}_setValue_fromArray_setMatrixWorldNeedsUpdate(t,e){this.resolvedProperty.fromArray(t,e),this.targetObject.matrixWorldNeedsUpdate=!0}_getValue_unbound(t,e){this.bind(),this.getValue(t,e)}_setValue_unbound(t,e){this.bind(),this.setValue(t,e)}bind(){let t=this.node,e=this.parsedPath,{objectName:n,propertyName:i,propertyIndex:s}=e;if(!t)t=le.findNode(this.rootNode,e.nodeName),this.node=t;if(this.getValue=this._getValue_unavailable,this.setValue=this._setValue_unavailable,!t){Mt("PropertyBinding: No target node found for track: "+this.path+".");return}if(n){let l=e.objectIndex;switch(n){case"materials":if(!t.material){Wt("PropertyBinding: Can not bind to material as node does not have a material.",this);return}if(!t.material.materials){Wt("PropertyBinding: Can not bind to material.materials as node.material does not have a materials array.",this);return}t=t.material.materials;break;case"bones":if(!t.skeleton){Wt("PropertyBinding: Can not bind to bones as node does not have a skeleton.",this);return}t=t.skeleton.bones;for(let c=0;c<t.length;c++)if(t[c].name===l){l=c;break}break;case"map":if("map"in t){t=t.map;break}if(!t.material){Wt("PropertyBinding: Can not bind to material as node does not have a material.",this);return}if(!t.material.map){Wt("PropertyBinding: Can not bind to material.map as node.material does not have a map.",this);return}t=t.material.map;break;default:if(t[n]===void 0){Wt("PropertyBinding: Can not bind to objectName of node undefined.",this);return}t=t[n]}if(l!==void 0){if(t[l]===void 0){Wt("PropertyBinding: Trying to bind to objectIndex of objectName, but is undefined.",this,t);return}t=t[l]}}let r=t[i];if(r===void 0){let l=e.nodeName;Wt("PropertyBinding: Trying to update property for track: "+l+"."+i+" but it wasn't found.",t);return}let a=this.Versioning.None;if(this.targetObject=t,t.isMaterial===!0)a=this.Versioning.NeedsUpdate;else if(t.isObject3D===!0)a=this.Versioning.MatrixWorldNeedsUpdate;let o=this.BindingType.Direct;if(s!==void 0){if(i==="morphTargetInfluences"){if(!t.geometry){Wt("PropertyBinding: Can not bind to morphTargetInfluences because node does not have a geometry.",this);return}if(!t.geometry.morphAttributes){Wt("PropertyBinding: Can not bind to morphTargetInfluences because node does not have a geometry.morphAttributes.",this);return}if(t.morphTargetDictionary[s]!==void 0)s=t.morphTargetDictionary[s]}o=this.BindingType.ArrayElement,this.resolvedProperty=r,this.propertyIndex=s}else if(r.fromArray!==void 0&&r.toArray!==void 0)o=this.BindingType.HasFromToArray,this.resolvedProperty=r;else if(Array.isArray(r))o=this.BindingType.EntireArray,this.resolvedProperty=r;else this.propertyName=i;this.getValue=this.GetterByBindingType[o],this.setValue=this.SetterByBindingTypeAndVersioning[o][a]}unbind(){this.node=null,this.getValue=this._getValue_unbound,this.setValue=this._setValue_unbound}}le.Composite=$f;le.prototype.BindingType={Direct:0,EntireArray:1,ArrayElement:2,HasFromToArray:3};le.prototype.Versioning={None:0,NeedsUpdate:1,MatrixWorldNeedsUpdate:2};le.prototype.GetterByBindingType=[le.prototype._getValue_direct,le.prototype._getValue_array,le.prototype._getValue_arrayElement,le.prototype._getValue_toArray];le.prototype.SetterByBindingTypeAndVersioning=[[le.prototype._setValue_direct,le.prototype._setValue_direct_setNeedsUpdate,le.prototype._setValue_direct_setMatrixWorldNeedsUpdate],[le.prototype._setValue_array,le.prototype._setValue_array_setNeedsUpdate,le.prototype._setValue_array_setMatrixWorldNeedsUpdate],[le.prototype._setValue_arrayElement,le.prototype._setValue_arrayElement_setNeedsUpdate,le.prototype._setValue_arrayElement_setMatrixWorldNeedsUpdate],[le.prototype._setValue_fromArray,le.prototype._setValue_fromArray_setNeedsUpdate,le.prototype._setValue_fromArray_setMatrixWorldNeedsUpdate]];class Kf{constructor(){this.isAnimationObjectGroup=!0,this.uuid=rn(),this._objects=Array.prototype.slice.call(arguments),this.nCachedObjects_=0;let t={};this._indicesByUUID=t;for(let n=0,i=arguments.length;n!==i;++n)t[arguments[n].uuid]=n;this._paths=[],this._parsedPaths=[],this._bindings=[],this._bindingsIndicesByPath={};let e=this;this.stats={objects:{get total(){return e._objects.length},get inUse(){return this.total-e.nCachedObjects_}},get bindingsPerObject(){return e._bindings.length}}}add(){let t=this._objects,e=this._indicesByUUID,n=this._paths,i=this._parsedPaths,s=this._bindings,r=s.length,a=void 0,o=t.length,l=this.nCachedObjects_;for(let c=0,h=arguments.length;c!==h;++c){let d=arguments[c],u=d.uuid,f=e[u];if(f===void 0){f=o++,e[u]=f,t.push(d);for(let g=0,_=r;g!==_;++g)s[g].push(new le(d,n[g],i[g]))}else if(f<l){a=t[f];let g=--l,_=t[g];e[_.uuid]=f,t[f]=_,e[u]=g,t[g]=d;for(let m=0,p=r;m!==p;++m){let b=s[m],S=b[g],x=b[f];if(b[f]=S,x===void 0)x=new le(d,n[m],i[m]);b[g]=x}}else if(t[f]!==a)Wt("AnimationObjectGroup: Different objects with the same UUID detected. Clean the caches or recreate your infrastructure when reloading scenes.")}this.nCachedObjects_=l}remove(){let t=this._objects,e=this._indicesByUUID,n=this._bindings,i=n.length,s=this.nCachedObjects_;for(let r=0,a=arguments.length;r!==a;++r){let o=arguments[r],l=o.uuid,c=e[l];if(c!==void 0&&c>=s){let h=s++,d=t[h];e[d.uuid]=c,t[c]=d,e[l]=h,t[h]=o;for(let u=0,f=i;u!==f;++u){let g=n[u],_=g[h],m=g[c];g[c]=_,g[h]=m}}}this.nCachedObjects_=s}uncache(){let t=this._objects,e=this._indicesByUUID,n=this._bindings,i=n.length,s=this.nCachedObjects_,r=t.length;for(let a=0,o=arguments.length;a!==o;++a){let l=arguments[a],c=l.uuid,h=e[c];if(h!==void 0)if(delete e[c],h<s){let d=--s,u=t[d],f=--r,g=t[f];if(h!==d)e[u.uuid]=h;if(t[h]=u,d!==f)e[g.uuid]=d;t[d]=g,t.pop();for(let _=0,m=i;_!==m;++_){let p=n[_],b=p[d],S=p[f];p[h]=b,p[d]=S,p.pop()}}else{let d=--r,u=t[d];if(h!==d)e[u.uuid]=h;t[h]=u,t.pop();for(let f=0,g=i;f!==g;++f){let _=n[f];_[h]=_[d],_.pop()}}}this.nCachedObjects_=s}subscribe_(t,e){let n=this._bindingsIndicesByPath,i=n[t],s=this._bindings;if(i!==void 0)return s[i];let r=this._paths,a=this._parsedPaths,o=this._objects,l=o.length,c=this.nCachedObjects_,h=Array(l);i=s.length,n[t]=i,r.push(t),a.push(e),s.push(h);for(let d=c,u=o.length;d!==u;++d){let f=o[d];h[d]=new le(f,t,e)}return h}unsubscribe_(t){let e=this._bindingsIndicesByPath,n=e[t];if(n!==void 0){let i=this._paths,s=this._parsedPaths,r=this._bindings,a=r.length-1,o=r[a],l=i[a];e[l]=n,r[n]=o,r.pop(),s[n]=s[a],s.pop(),i[n]=i[a],i.pop()}}}class _h{constructor(t,e,n=null,i=e.blendMode){this._mixer=t,this._clip=e,this._localRoot=n,this.blendMode=i;let s=e.tracks,r=s.length,a=Array(r),o={endingStart:2400,endingEnd:2400};for(let l=0;l!==r;++l){let c=s[l].createInterpolant(null);a[l]=c,c.settings=o}this._interpolantSettings=o,this._interpolants=a,this._propertyBindings=Array(r),this._cacheIndex=null,this._byClipCacheIndex=null,this._timeScaleInterpolant=null,this._restoreTimeScale=null,this._weightInterpolant=null,this.loop=2201,this._loopCount=-1,this._startTime=null,this.time=0,this.timeScale=1,this._effectiveTimeScale=1,this.weight=1,this._effectiveWeight=1,this.repetitions=1/0,this.paused=!1,this.enabled=!0,this.clampWhenFinished=!1,this.zeroSlopeAtStart=!0,this.zeroSlopeAtEnd=!0}play(){return this._mixer._activateAction(this),this}stop(){return this._mixer._deactivateAction(this),this.reset()}reset(){return this.paused=!1,this.enabled=!0,this.time=0,this._loopCount=-1,this._startTime=null,this.stopFading().stopWarping()}isRunning(){return this.enabled&&!this.paused&&this.timeScale!==0&&this._startTime===null&&this._mixer._isActiveAction(this)}isScheduled(){return this._mixer._isActiveAction(this)}startAt(t){return this._startTime=t,this}setLoop(t,e){return this.loop=t,this.repetitions=e,this}setEffectiveWeight(t){return this.weight=t,this._effectiveWeight=this.enabled?t:0,this.stopFading()}getEffectiveWeight(){return this._effectiveWeight}fadeIn(t){return this._scheduleFading(t,0,1)}fadeOut(t){return this._scheduleFading(t,1,0)}crossFadeFrom(t,e,n=!1){if(t.fadeOut(e),this.fadeIn(e),n===!0){let i=this._clip.duration,s=t._clip.duration,r=s/i,a=i/s;t._restoreTimeScale=t.timeScale,this._restoreTimeScale=this.timeScale,t.warp(1,r,e),this.warp(a,1,e)}return this}crossFadeTo(t,e,n=!1){return t.crossFadeFrom(this,e,n)}stopFading(){let t=this._weightInterpolant;if(t!==null)this._weightInterpolant=null,this._mixer._takeBackControlInterpolant(t);return this}setEffectiveTimeScale(t){return this.timeScale=t,this._effectiveTimeScale=this.paused?0:t,this.stopWarping()}getEffectiveTimeScale(){return this._effectiveTimeScale}setDuration(t){return this.timeScale=this._clip.duration/t,this.stopWarping()}syncWith(t){return this.time=t.time,this.timeScale=t.timeScale,this.stopWarping()}halt(t){return this.warp(this._effectiveTimeScale,0,t)}warp(t,e,n){let i=this._mixer,s=i.time,r=this.timeScale,a=this._timeScaleInterpolant;if(a===null)a=i._lendControlInterpolant(),this._timeScaleInterpolant=a;let o=a.parameterPositions,l=a.sampleValues;return o[0]=s,o[1]=s+n,l[0]=t/r,l[1]=e/r,this}stopWarping(){let t=this._timeScaleInterpolant;if(t!==null)this._timeScaleInterpolant=null,this._mixer._takeBackControlInterpolant(t);return this._restoreTimeScale=null,this}getMixer(){return this._mixer}getClip(){return this._clip}getRoot(){return this._localRoot||this._mixer._root}_update(t,e,n,i){if(!this.enabled){this._updateWeight(t);return}let s=this._startTime;if(s!==null){let o=(t-s)*n;if(o<0||n===0)e=0;else this._startTime=null,e=n*o}e*=this._updateTimeScale(t);let r=this._updateTime(e),a=this._updateWeight(t);if(a>0){let o=this._interpolants,l=this._propertyBindings;switch(this.blendMode){case 2501:for(let c=0,h=o.length;c!==h;++c)o[c].evaluate(r),l[c].accumulateAdditive(a);break;case 2500:default:for(let c=0,h=o.length;c!==h;++c)o[c].evaluate(r),l[c].accumulate(i,a)}}}_updateWeight(t){let e=0;if(this.enabled){e=this.weight;let n=this._weightInterpolant;if(n!==null){let i=n.evaluate(t)[0];if(e*=i,t>n.parameterPositions[1]){if(this.stopFading(),i===0)this.enabled=!1}}}return this._effectiveWeight=e,e}_updateTimeScale(t){let e=0;if(!this.paused){e=this.timeScale;let n=this._timeScaleInterpolant;if(n!==null){let i=n.evaluate(t)[0];if(e*=i,t>n.parameterPositions[1]){if(e===0)this.paused=!0;else{if(this._restoreTimeScale!==null)e=this._restoreTimeScale;this.timeScale=e}this.stopWarping()}}}return this._effectiveTimeScale=e,e}_updateTime(t){let e=this._clip.duration,n=this.loop,i=this.time+t,s=this._loopCount,r=n===2202;if(t===0){if(s===-1)return i;return r&&(s&1)===1?e-i:i}if(n===2200){if(s===-1)this._loopCount=0,this._setEndings(!0,!0,!1);t:{if(i>=e)i=e;else if(i<0)i=0;else{this.time=i;break t}if(this.clampWhenFinished)this.paused=!0;else this.enabled=!1;this.time=i,this._mixer.dispatchEvent({type:"finished",action:this,direction:t<0?-1:1})}}else{if(s===-1)if(t>=0)s=0,this._setEndings(!0,this.repetitions===0,r);else this._setEndings(this.repetitions===0,!0,r);if(i>=e||i<0){let a=Math.floor(i/e);i-=e*a,s+=Math.abs(a);let o=this.repetitions-s;if(o<=0){if(this.clampWhenFinished)this.paused=!0;else this.enabled=!1;i=t>0?e:0,this.time=i,this._mixer.dispatchEvent({type:"finished",action:this,direction:t>0?1:-1})}else{if(o===1){let l=t<0;this._setEndings(l,!l,r)}else this._setEndings(!1,!1,r);this._loopCount=s,this.time=i,this._mixer.dispatchEvent({type:"loop",action:this,loopDelta:a})}}else this._loopCount=s,this.time=i;if(r&&(s&1)===1)return e-i}return i}_setEndings(t,e,n){let i=this._interpolantSettings;if(n)i.endingStart=2401,i.endingEnd=2401;else{if(t)i.endingStart=this.zeroSlopeAtStart?2401:2400;else i.endingStart=2402;if(e)i.endingEnd=this.zeroSlopeAtEnd?2401:2400;else i.endingEnd=2402}}_scheduleFading(t,e,n){let i=this._mixer,s=i.time,r=this._weightInterpolant;if(r===null)r=i._lendControlInterpolant(),this._weightInterpolant=r;let a=r.parameterPositions,o=r.sampleValues;return a[0]=s,o[0]=e,a[1]=s+t,o[1]=n,this}}var p_=new Float32Array(1);class jf extends pn{constructor(t){super();if(this._root=t,this._initMemoryManager(),this._accuIndex=0,this.time=0,this.timeScale=1,typeof __THREE_DEVTOOLS__<"u")__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("observe",{detail:this}))}_bindAction(t,e){let n=t._localRoot||this._root,i=t._clip.tracks,s=i.length,{_propertyBindings:r,_interpolants:a}=t,o=n.uuid,l=this._bindingsByRootAndName,c=l[o];if(c===void 0)c={},l[o]=c;for(let h=0;h!==s;++h){let d=i[h],u=d.name,f=c[u];if(f!==void 0)++f.referenceCount,r[h]=f;else{if(f=r[h],f!==void 0){if(f._cacheIndex===null)++f.referenceCount,this._addInactiveBinding(f,o,u);continue}let g=e&&e._propertyBindings[h].binding.parsedPath;f=new ph(le.create(n,u,g),d.ValueTypeName,d.getValueSize()),++f.referenceCount,this._addInactiveBinding(f,o,u),r[h]=f}a[h].resultBuffer=f.buffer}}_activateAction(t){if(!this._isActiveAction(t)){if(t._cacheIndex===null){let n=(t._localRoot||this._root).uuid,i=t._clip.uuid,s=this._actionsByClip[i];this._bindAction(t,s&&s.knownActions[0]),this._addInactiveAction(t,i,n)}let e=t._propertyBindings;for(let n=0,i=e.length;n!==i;++n){let s=e[n];if(s.useCount++===0)this._lendBinding(s),s.saveOriginalState()}this._lendAction(t)}}_deactivateAction(t){if(this._isActiveAction(t)){let e=t._propertyBindings;for(let n=0,i=e.length;n!==i;++n){let s=e[n];if(--s.useCount===0)s.restoreOriginalState(),this._takeBackBinding(s)}this._takeBackAction(t)}}_initMemoryManager(){this._actions=[],this._nActiveActions=0,this._actionsByClip={},this._bindings=[],this._nActiveBindings=0,this._bindingsByRootAndName={},this._controlInterpolants=[],this._nActiveControlInterpolants=0;let t=this;this.stats={actions:{get total(){return t._actions.length},get inUse(){return t._nActiveActions}},bindings:{get total(){return t._bindings.length},get inUse(){return t._nActiveBindings}},controlInterpolants:{get total(){return t._controlInterpolants.length},get inUse(){return t._nActiveControlInterpolants}}}}_isActiveAction(t){let e=t._cacheIndex;return e!==null&&e<this._nActiveActions}_addInactiveAction(t,e,n){let i=this._actions,s=this._actionsByClip,r=s[e];if(r===void 0)r={knownActions:[t],actionByRoot:{}},t._byClipCacheIndex=0,s[e]=r;else{let a=r.knownActions;t._byClipCacheIndex=a.length,a.push(t)}t._cacheIndex=i.length,i.push(t),r.actionByRoot[n]=t}_removeInactiveAction(t){let e=this._actions,n=e[e.length-1],i=t._cacheIndex;n._cacheIndex=i,e[i]=n,e.pop(),t._cacheIndex=null;let s=t._clip.uuid,r=this._actionsByClip,a=r[s],o=a.knownActions,l=o[o.length-1],c=t._byClipCacheIndex;l._byClipCacheIndex=c,o[c]=l,o.pop(),t._byClipCacheIndex=null;let h=a.actionByRoot,d=(t._localRoot||this._root).uuid;if(delete h[d],o.length===0)delete r[s];this._removeInactiveBindingsForAction(t)}_removeInactiveBindingsForAction(t){let e=t._propertyBindings;for(let n=0,i=e.length;n!==i;++n){let s=e[n];if(--s.referenceCount===0)this._removeInactiveBinding(s)}}_lendAction(t){let e=this._actions,n=t._cacheIndex,i=this._nActiveActions++,s=e[i];t._cacheIndex=i,e[i]=t,s._cacheIndex=n,e[n]=s}_takeBackAction(t){let e=this._actions,n=t._cacheIndex,i=--this._nActiveActions,s=e[i];t._cacheIndex=i,e[i]=t,s._cacheIndex=n,e[n]=s}_addInactiveBinding(t,e,n){let i=this._bindingsByRootAndName,s=this._bindings,r=i[e];if(r===void 0)r={},i[e]=r;r[n]=t,t._cacheIndex=s.length,s.push(t)}_removeInactiveBinding(t){let e=this._bindings,n=t.binding,i=n.rootNode.uuid,s=n.path,r=this._bindingsByRootAndName,a=r[i],o=e[e.length-1],l=t._cacheIndex;if(o._cacheIndex=l,e[l]=o,e.pop(),delete a[s],Object.keys(a).length===0)delete r[i]}_lendBinding(t){let e=this._bindings,n=t._cacheIndex,i=this._nActiveBindings++,s=e[i];t._cacheIndex=i,e[i]=t,s._cacheIndex=n,e[n]=s}_takeBackBinding(t){let e=this._bindings,n=t._cacheIndex,i=--this._nActiveBindings,s=e[i];t._cacheIndex=i,e[i]=t,s._cacheIndex=n,e[n]=s}_lendControlInterpolant(){let t=this._controlInterpolants,e=this._nActiveControlInterpolants++,n=t[e];if(n===void 0)n=new bo(new Float32Array(2),new Float32Array(2),1,p_),n.__cacheIndex=e,t[e]=n;return n}_takeBackControlInterpolant(t){let e=this._controlInterpolants,n=t.__cacheIndex,i=--this._nActiveControlInterpolants,s=e[i];t.__cacheIndex=i,e[i]=t,s.__cacheIndex=n,e[n]=s}clipAction(t,e,n){let i=e||this._root,s=i.uuid,r=typeof t==="string"?ys.findByName(i,t):t,a=r!==null?r.uuid:t,o=this._actionsByClip[a],l=null;if(n===void 0)if(r!==null)n=r.blendMode;else n=2500;if(o!==void 0){let h=o.actionByRoot[s];if(h!==void 0&&h.blendMode===n)return h;if(l=o.knownActions[0],r===null)r=l._clip}if(r===null)return null;let c=new _h(this,r,e,n);return this._bindAction(c,l),this._addInactiveAction(c,a,s),c}existingAction(t,e){let n=e||this._root,i=n.uuid,s=typeof t==="string"?ys.findByName(n,t):t,r=s?s.uuid:t,a=this._actionsByClip[r];if(a!==void 0)return a.actionByRoot[i]||null;return null}stopAllAction(){let t=this._actions,e=this._nActiveActions;for(let n=e-1;n>=0;--n)t[n].stop();return this}update(t){t*=this.timeScale;let e=this._actions,n=this._nActiveActions,i=this.time+=t,s=Math.sign(t),r=this._accuIndex^=1;for(let l=0;l!==n;++l)e[l]._update(i,t,s,r);let a=this._bindings,o=this._nActiveBindings;for(let l=0;l!==o;++l)a[l].apply(r);return this}setTime(t){this.time=0;for(let e=0;e<this._actions.length;e++)this._actions[e].time=0;return this.update(t)}getRoot(){return this._root}uncacheClip(t){let e=this._actions,n=t.uuid,i=this._actionsByClip,s=i[n];if(s!==void 0){let r=s.knownActions;for(let a=0,o=r.length;a!==o;++a){let l=r[a];this._deactivateAction(l);let c=l._cacheIndex,h=e[e.length-1];l._cacheIndex=null,l._byClipCacheIndex=null,h._cacheIndex=c,e[c]=h,e.pop(),this._removeInactiveBindingsForAction(l)}delete i[n]}}uncacheRoot(t){let e=t.uuid,n=this._actionsByClip;for(let r in n){let a=n[r].actionByRoot,o=a[e];if(o!==void 0)this._deactivateAction(o),this._removeInactiveAction(o)}let i=this._bindingsByRootAndName,s=i[e];if(s!==void 0)for(let r in s){let a=s[r];a.restoreOriginalState(),this._removeInactiveBinding(a)}}uncacheAction(t,e){let n=this.existingAction(t,e);if(n!==null)this._deactivateAction(n),this._removeInactiveAction(n)}}class Qf extends Wa{constructor(t=1,e=1,n=1,i={}){super(t,e,i);this.isRenderTarget3D=!0,this.depth=n;for(let s=0;s<this.textures.length;s++){let r=new hr(null,t,e,n);r.isRenderTargetTexture=!0,r.renderTarget=this,this.textures[s]=r}this._setTextureOptions(i)}}class xh{constructor(t){this.value=t}clone(){return new xh(this.value.clone===void 0?this.value:this.value.clone())}}var m_=0;class tp extends pn{constructor(){super();this.isUniformsGroup=!0,Object.defineProperty(this,"id",{value:m_++}),this.name="",this.usage=35044,this.uniforms=[]}add(t){return this.uniforms.push(t),this}remove(t){let e=this.uniforms.indexOf(t);if(e!==-1)this.uniforms.splice(e,1);return this}setName(t){return this.name=t,this}setUsage(t){return this.usage=t,this}dispose(){this.dispatchEvent({type:"dispose"})}copy(t){this.name=t.name,this.usage=t.usage;let e=t.uniforms;this.uniforms.length=0;for(let n=0,i=e.length;n<i;n++){let s=Array.isArray(e[n])?e[n]:[e[n]];for(let r=0;r<s.length;r++)this.uniforms.push(s[r].clone())}return this}clone(){return new this.constructor().copy(this)}}class ep extends ws{constructor(t,e,n=1){super(t,e);this.isInstancedInterleavedBuffer=!0,this.meshPerAttribute=n}copy(t){return super.copy(t),this.meshPerAttribute=t.meshPerAttribute,this}clone(t){let e=super.clone(t);return e.meshPerAttribute=this.meshPerAttribute,e}toJSON(t){let e=super.toJSON(t);return e.isInstancedInterleavedBuffer=!0,e.meshPerAttribute=this.meshPerAttribute,e}}class np{constructor(t,e,n,i,s,r=!1){this.isGLBufferAttribute=!0,this.name="",this.buffer=t,this.type=e,this.itemSize=n,this.elementSize=i,this.count=s,this.normalized=r,this.version=0}set needsUpdate(t){if(t===!0)this.version++}setBuffer(t){return this.buffer=t,this}setType(t,e){return this.type=t,this.elementSize=e,this}setItemSize(t){return this.itemSize=t,this}setCount(t){return this.count=t,this}}var Xu=new Kt;class ip{constructor(t,e,n=0,i=1/0){this.ray=new Kn(t,e),this.near=n,this.far=i,this.camera=null,this.layers=new ur,this.params={Mesh:{},Line:{threshold:1},LOD:{},Points:{threshold:1},Sprite:{}}}set(t,e){this.ray.set(t,e)}setFromCamera(t,e){if(e.isPerspectiveCamera)this.ray.origin.setFromMatrixPosition(e.matrixWorld),this.ray.direction.set(t.x,t.y,0.5).unproject(e).sub(this.ray.origin).normalize(),this.camera=e;else if(e.isOrthographicCamera)this.ray.origin.set(t.x,t.y,e.projectionMatrix.elements[14]).unproject(e),this.ray.direction.set(0,0,-1).transformDirection(e.matrixWorld),this.camera=e;else Wt("Raycaster: Unsupported camera type: "+e.type)}setFromXRController(t){return Xu.identity().extractRotation(t.matrixWorld),this.ray.origin.setFromMatrixPosition(t.matrixWorld),this.ray.direction.set(0,0,-1).applyMatrix4(Xu),this}intersectObject(t,e=!0,n=[]){return El(t,this,n,e),n.sort(qu),n}intersectObjects(t,e=!0,n=[]){for(let i=0,s=t.length;i<s;i++)El(t[i],this,n,e);return n.sort(qu),n}}function qu(t,e){return t.distance-e.distance}function El(t,e,n,i){let s=!0;if(t.layers.test(e.layers)){if(t.raycast(e,n)===!1)s=!1}if(s===!0&&i===!0){let r=t.children;for(let a=0,o=r.length;a<o;a++)El(r[a],e,n,!0)}}class sp{constructor(t=!0){this.autoStart=t,this.startTime=0,this.oldTime=0,this.elapsedTime=0,this.running=!1,Mt("Clock: This module has been deprecated. Please use THREE.Timer instead.")}start(){this.startTime=performance.now(),this.oldTime=this.startTime,this.elapsedTime=0,this.running=!0}stop(){this.getElapsedTime(),this.running=!1,this.autoStart=!1}getElapsedTime(){return this.getDelta(),this.elapsedTime}getDelta(){let t=0;if(this.autoStart&&!this.running)return this.start(),0;if(this.running){let e=performance.now();t=(e-this.oldTime)/1000,this.oldTime=e,this.elapsedTime+=t}return t}}class wr{constructor(t=1,e=0,n=0){this.radius=t,this.phi=e,this.theta=n}set(t,e,n){return this.radius=t,this.phi=e,this.theta=n,this}copy(t){return this.radius=t.radius,this.phi=t.phi,this.theta=t.theta,this}makeSafe(){return this.phi=Zt(this.phi,0.000001,Math.PI-0.000001),this}setFromVector3(t){return this.setFromCartesianCoords(t.x,t.y,t.z)}setFromCartesianCoords(t,e,n){if(this.radius=Math.sqrt(t*t+e*e+n*n),this.radius===0)this.theta=0,this.phi=0;else this.theta=Math.atan2(t,n),this.phi=Math.acos(Zt(e/this.radius,-1,1));return this}clone(){return new this.constructor().copy(this)}}class rp{constructor(t=1,e=0,n=0){this.radius=t,this.theta=e,this.y=n}set(t,e,n){return this.radius=t,this.theta=e,this.y=n,this}copy(t){return this.radius=t.radius,this.theta=t.theta,this.y=t.y,this}setFromVector3(t){return this.setFromCartesianCoords(t.x,t.y,t.z)}setFromCartesianCoords(t,e,n){return this.radius=Math.sqrt(t*t+n*n),this.theta=Math.atan2(t,n),this.y=e,this}clone(){return new this.constructor().copy(this)}}class yh{static{yh.prototype.isMatrix2=!0}constructor(t,e,n,i){if(this.elements=[1,0,0,1],t!==void 0)this.set(t,e,n,i)}identity(){return this.set(1,0,0,1),this}fromArray(t,e=0){for(let n=0;n<4;n++)this.elements[n]=t[n+e];return this}set(t,e,n,i){let s=this.elements;return s[0]=t,s[2]=e,s[1]=n,s[3]=i,this}}var Yu=new rt;class Rr{constructor(t=new rt(1/0,1/0),e=new rt(-1/0,-1/0)){this.isBox2=!0,this.min=t,this.max=e}set(t,e){return this.min.copy(t),this.max.copy(e),this}setFromPoints(t){this.makeEmpty();for(let e=0,n=t.length;e<n;e++)this.expandByPoint(t[e]);return this}setFromCenterAndSize(t,e){let n=Yu.copy(e).multiplyScalar(0.5);return this.min.copy(t).sub(n),this.max.copy(t).add(n),this}clone(){return new this.constructor().copy(this)}copy(t){return this.min.copy(t.min),this.max.copy(t.max),this}makeEmpty(){return this.min.x=this.min.y=1/0,this.max.x=this.max.y=-1/0,this}isEmpty(){return this.max.x<this.min.x||this.max.y<this.min.y}getCenter(t){return this.isEmpty()?t.set(0,0):t.addVectors(this.min,this.max).multiplyScalar(0.5)}getSize(t){return this.isEmpty()?t.set(0,0):t.subVectors(this.max,this.min)}expandByPoint(t){return this.min.min(t),this.max.max(t),this}expandByVector(t){return this.min.sub(t),this.max.add(t),this}expandByScalar(t){return this.min.addScalar(-t),this.max.addScalar(t),this}containsPoint(t){return t.x>=this.min.x&&t.x<=this.max.x&&t.y>=this.min.y&&t.y<=this.max.y}containsBox(t){return this.min.x<=t.min.x&&t.max.x<=this.max.x&&this.min.y<=t.min.y&&t.max.y<=this.max.y}getParameter(t,e){return e.set((t.x-this.min.x)/(this.max.x-this.min.x),(t.y-this.min.y)/(this.max.y-this.min.y))}intersectsBox(t){return t.max.x>=this.min.x&&t.min.x<=this.max.x&&t.max.y>=this.min.y&&t.min.y<=this.max.y}clampPoint(t,e){return e.copy(t).clamp(this.min,this.max)}distanceToPoint(t){return this.clampPoint(t,Yu).distanceTo(t)}intersect(t){if(this.min.max(t.min),this.max.min(t.max),this.isEmpty())this.makeEmpty();return this}union(t){return this.min.min(t.min),this.max.max(t.max),this}translate(t){return this.min.add(t),this.max.add(t),this}equals(t){return t.min.equals(this.min)&&t.max.equals(this.max)}}var Zu=new L,va=new L,fs=new L,ps=new L,xl=new L,g_=new L,__=new L;class ap{constructor(t=new L,e=new L){this.start=t,this.end=e}set(t,e){return this.start.copy(t),this.end.copy(e),this}copy(t){return this.start.copy(t.start),this.end.copy(t.end),this}getCenter(t){return t.addVectors(this.start,this.end).multiplyScalar(0.5)}delta(t){return t.subVectors(this.end,this.start)}distanceSq(){return this.start.distanceToSquared(this.end)}distance(){return this.start.distanceTo(this.end)}at(t,e){return this.delta(e).multiplyScalar(t).add(this.start)}closestPointToPointParameter(t,e){Zu.subVectors(t,this.start),va.subVectors(this.end,this.start);let n=va.dot(va);if(n===0)return 0;let s=va.dot(Zu)/n;if(e)s=Zt(s,0,1);return s}closestPointToPoint(t,e,n){let i=this.closestPointToPointParameter(t,e);return this.delta(n).multiplyScalar(i).add(this.start)}distanceSqToLine3(t,e=g_,n=__){let s,r,a=this.start,o=t.start,l=this.end,c=t.end;fs.subVectors(l,a),ps.subVectors(c,o),xl.subVectors(a,o);let h=fs.dot(fs),d=ps.dot(ps),u=ps.dot(xl);if(h<=0.00000000000000010000000000000001&&d<=0.00000000000000010000000000000001)return e.copy(a),n.copy(o),e.sub(n),e.dot(e);if(h<=0.00000000000000010000000000000001)s=0,r=u/d,r=Zt(r,0,1);else{let f=fs.dot(xl);if(d<=0.00000000000000010000000000000001)r=0,s=Zt(-f/h,0,1);else{let g=fs.dot(ps),_=h*d-g*g;if(_!==0)s=Zt((g*u-f*d)/_,0,1);else s=0;if(r=(g*s+u)/d,r<0)r=0,s=Zt(-f/h,0,1);else if(r>1)r=1,s=Zt((g-f)/h,0,1)}}return e.copy(a).addScaledVector(fs,s),n.copy(o).addScaledVector(ps,r),e.distanceToSquared(n)}applyMatrix4(t){return this.start.applyMatrix4(t),this.end.applyMatrix4(t),this}equals(t){return t.start.equals(this.start)&&t.end.equals(this.end)}clone(){return new this.constructor().copy(this)}}var Ju=new L;class op extends ae{constructor(t,e){super();this.light=t,this.matrixAutoUpdate=!1,this.color=e,this.type="SpotLightHelper";let n=new jt,i=[0,0,0,0,0,1,0,0,0,1,0,1,0,0,0,-1,0,1,0,0,0,0,1,1,0,0,0,0,-1,1];for(let r=0,a=1,o=32;r<o;r++,a++){let l=r/o*Math.PI*2,c=a/o*Math.PI*2;i.push(Math.cos(l),Math.sin(l),1,Math.cos(c),Math.sin(c),1)}n.setAttribute("position",new Nt(i,3));let s=new We({fog:!1,toneMapped:!1});this.cone=new En(n,s),this.add(this.cone),this.update()}dispose(){super.dispose(),this.cone.geometry.dispose(),this.cone.material.dispose()}update(){if(this.light.updateWorldMatrix(!0,!1),this.light.target.updateWorldMatrix(!0,!1),this.parent)this.parent.updateWorldMatrix(!0),this.matrix.copy(this.parent.matrixWorld).invert().multiply(this.light.matrixWorld);else this.matrix.copy(this.light.matrixWorld);this.matrixWorldNeedsUpdate=!0;let t=this.light.distance?this.light.distance:1000,e=t*Math.tan(this.light.angle);if(this.cone.scale.set(e,e,t),Ju.setFromMatrixPosition(this.light.target.matrixWorld),this.cone.lookAt(Ju),this.color!==void 0)this.cone.material.color.set(this.color);else this.cone.material.color.copy(this.light.color)}}var oi=new L,Sa=new Kt,yl=new Kt;class lp extends En{constructor(t){let e=cp(t),n=new jt,i=[],s=[];for(let l=0;l<e.length;l++){let c=e[l];if(c.parent&&c.parent.isBone)i.push(0,0,0),i.push(0,0,0),s.push(0,0,0),s.push(0,0,0)}n.setAttribute("position",new Nt(i,3)),n.setAttribute("color",new Nt(s,3));let r=new We({vertexColors:!0,depthTest:!1,depthWrite:!1,toneMapped:!1,transparent:!0});super(n,r);this.isSkeletonHelper=!0,this.type="SkeletonHelper",this.root=t,this.bones=e,this.matrix=t.matrixWorld,this.matrixAutoUpdate=!1;let a=new Pt(255),o=new Pt(65280);this.setColors(a,o)}updateMatrixWorld(t){let e=this.bones,n=this.geometry,i=n.getAttribute("position");yl.copy(this.root.matrixWorld).invert();for(let s=0,r=0;s<e.length;s++){let a=e[s];if(a.parent&&a.parent.isBone)Sa.multiplyMatrices(yl,a.matrixWorld),oi.setFromMatrixPosition(Sa),i.setXYZ(r,oi.x,oi.y,oi.z),Sa.multiplyMatrices(yl,a.parent.matrixWorld),oi.setFromMatrixPosition(Sa),i.setXYZ(r+1,oi.x,oi.y,oi.z),r+=2}n.getAttribute("position").needsUpdate=!0,super.updateMatrixWorld(t)}setColors(t,e){let i=this.geometry.getAttribute("color");for(let s=0;s<i.count;s+=2)i.setXYZ(s,t.r,t.g,t.b),i.setXYZ(s+1,e.r,e.g,e.b);return i.needsUpdate=!0,this}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}}function cp(t){let e=[];if(t.isBone===!0)e.push(t);for(let n=0;n<t.children.length;n++)e.push(...cp(t.children[n]));return e}class hp extends pe{constructor(t,e,n){let i=new vr(e,4,2),s=new mn({wireframe:!0,fog:!1,toneMapped:!1});super(i,s);this.light=t,this.color=n,this.type="PointLightHelper",this.matrix=this.light.matrixWorld,this.matrixAutoUpdate=!1,this.update()}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}update(){if(this.matrixWorldNeedsUpdate=!0,this.light.updateWorldMatrix(!0,!1),this.color!==void 0)this.material.color.set(this.color);else this.material.color.copy(this.light.color)}}var x_=new L,$u=new Pt,Ku=new Pt;class up extends ae{constructor(t,e,n){super();this.light=t,this.matrix=t.matrixWorld,this.matrixAutoUpdate=!1,this.color=n,this.type="HemisphereLightHelper";let i=new yr(e);if(i.rotateY(Math.PI*0.5),this.material=new mn({wireframe:!0,fog:!1,toneMapped:!1}),this.color===void 0)this.material.vertexColors=!0;let s=i.getAttribute("position"),r=new Float32Array(s.count*3);i.setAttribute("color",new de(r,3)),this.add(new pe(i,this.material)),this.update()}dispose(){super.dispose(),this.children[0].geometry.dispose(),this.children[0].material.dispose()}update(){let t=this.children[0];if(this.color!==void 0)this.material.color.set(this.color);else{let e=t.geometry.getAttribute("color");$u.copy(this.light.color),Ku.copy(this.light.groundColor);for(let n=0,i=e.count;n<i;n++){let s=n<i/2?$u:Ku;e.setXYZ(n,s.r,s.g,s.b)}e.needsUpdate=!0}this.matrixWorldNeedsUpdate=!0,this.light.updateWorldMatrix(!0,!1),t.lookAt(x_.setFromMatrixPosition(this.light.matrixWorld).negate())}}class dp extends En{constructor(t=10,e=10,n=4473924,i=8947848){n=new Pt(n),i=new Pt(i);let s=e/2,r=t/e,a=t/2,o=[],l=[];for(let d=0,u=0,f=-a;d<=e;d++,f+=r){o.push(-a,0,f,a,0,f),o.push(f,0,-a,f,0,a);let g=d===s?n:i;g.toArray(l,u),u+=3,g.toArray(l,u),u+=3,g.toArray(l,u),u+=3,g.toArray(l,u),u+=3}let c=new jt;c.setAttribute("position",new Nt(o,3)),c.setAttribute("color",new Nt(l,3));let h=new We({vertexColors:!0,toneMapped:!1});super(c,h);this.type="GridHelper"}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}}class fp extends En{constructor(t=10,e=16,n=8,i=64,s=4473924,r=8947848){s=new Pt(s),r=new Pt(r);let a=[],o=[];if(e>1)for(let h=0;h<e;h++){let d=h/e*(Math.PI*2),u=Math.sin(d)*t,f=Math.cos(d)*t;a.push(0,0,0),a.push(u,0,f);let g=h&1?s:r;o.push(g.r,g.g,g.b),o.push(g.r,g.g,g.b)}for(let h=0;h<n;h++){let d=h&1?s:r,u=t-t/n*h;for(let f=0;f<i;f++){let g=f/i*(Math.PI*2),_=Math.sin(g)*u,m=Math.cos(g)*u;a.push(_,0,m),o.push(d.r,d.g,d.b),g=(f+1)/i*(Math.PI*2),_=Math.sin(g)*u,m=Math.cos(g)*u,a.push(_,0,m),o.push(d.r,d.g,d.b)}}let l=new jt;l.setAttribute("position",new Nt(a,3)),l.setAttribute("color",new Nt(o,3));let c=new We({vertexColors:!0,toneMapped:!1});super(l,c);this.type="PolarGridHelper"}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}}var ju=new L,Ma=new L,Qu=new L;class pp extends ae{constructor(t,e,n){super();if(this.light=t,this.matrix=t.matrixWorld,this.matrixAutoUpdate=!1,this.color=n,this.type="DirectionalLightHelper",e===void 0)e=1;let i=new jt;i.setAttribute("position",new Nt([-e,e,0,e,e,0,e,-e,0,-e,-e,0,-e,e,0],3));let s=new We({fog:!1,toneMapped:!1});this.lightPlane=new Zn(i,s),this.add(this.lightPlane),i=new jt,i.setAttribute("position",new Nt([0,0,0,0,0,1],3)),this.targetLine=new Zn(i,s),this.add(this.targetLine),this.update()}dispose(){super.dispose(),this.lightPlane.geometry.dispose(),this.lightPlane.material.dispose(),this.targetLine.geometry.dispose(),this.targetLine.material.dispose()}update(){if(this.matrixWorldNeedsUpdate=!0,this.light.updateWorldMatrix(!0,!1),this.light.target.updateWorldMatrix(!0,!1),ju.setFromMatrixPosition(this.light.matrixWorld),Ma.setFromMatrixPosition(this.light.target.matrixWorld),Qu.subVectors(Ma,ju),this.lightPlane.lookAt(Ma),this.color!==void 0)this.lightPlane.material.color.set(this.color),this.targetLine.material.color.set(this.color);else this.lightPlane.material.color.copy(this.light.color),this.targetLine.material.color.copy(this.light.color);this.targetLine.lookAt(Ma),this.targetLine.scale.z=Qu.length()}}var ba=new L,Te=new Er;class mp extends En{constructor(t){let e=new jt,n=new We({color:16777215,vertexColors:!0,toneMapped:!1}),i=[],s=[],r={};a("n1","n2"),a("n2","n4"),a("n4","n3"),a("n3","n1"),a("f1","f2"),a("f2","f4"),a("f4","f3"),a("f3","f1"),a("n1","f1"),a("n2","f2"),a("n3","f3"),a("n4","f4"),a("p","n1"),a("p","n2"),a("p","n3"),a("p","n4"),a("u1","u2"),a("u2","u3"),a("u3","u1"),a("c","t"),a("p","c"),a("cn1","cn2"),a("cn3","cn4"),a("cf1","cf2"),a("cf3","cf4");function a(f,g){o(f),o(g)}function o(f){if(i.push(0,0,0),s.push(0,0,0),r[f]===void 0)r[f]=[];r[f].push(i.length/3-1)}e.setAttribute("position",new Nt(i,3)),e.setAttribute("color",new Nt(s,3));super(e,n);if(this.type="CameraHelper",this.camera=t,this.camera.updateProjectionMatrix)this.camera.updateProjectionMatrix();this.matrix=t.matrixWorld,this.matrixAutoUpdate=!1,this.pointMap=r,this.update();let l=new Pt(16755200),c=new Pt(16711680),h=new Pt(43775),d=new Pt(16777215),u=new Pt(3355443);this.setColors(l,c,h,d,u)}setColors(t,e,n,i,s){let a=this.geometry.getAttribute("color");return a.setXYZ(0,t.r,t.g,t.b),a.setXYZ(1,t.r,t.g,t.b),a.setXYZ(2,t.r,t.g,t.b),a.setXYZ(3,t.r,t.g,t.b),a.setXYZ(4,t.r,t.g,t.b),a.setXYZ(5,t.r,t.g,t.b),a.setXYZ(6,t.r,t.g,t.b),a.setXYZ(7,t.r,t.g,t.b),a.setXYZ(8,t.r,t.g,t.b),a.setXYZ(9,t.r,t.g,t.b),a.setXYZ(10,t.r,t.g,t.b),a.setXYZ(11,t.r,t.g,t.b),a.setXYZ(12,t.r,t.g,t.b),a.setXYZ(13,t.r,t.g,t.b),a.setXYZ(14,t.r,t.g,t.b),a.setXYZ(15,t.r,t.g,t.b),a.setXYZ(16,t.r,t.g,t.b),a.setXYZ(17,t.r,t.g,t.b),a.setXYZ(18,t.r,t.g,t.b),a.setXYZ(19,t.r,t.g,t.b),a.setXYZ(20,t.r,t.g,t.b),a.setXYZ(21,t.r,t.g,t.b),a.setXYZ(22,t.r,t.g,t.b),a.setXYZ(23,t.r,t.g,t.b),a.setXYZ(24,e.r,e.g,e.b),a.setXYZ(25,e.r,e.g,e.b),a.setXYZ(26,e.r,e.g,e.b),a.setXYZ(27,e.r,e.g,e.b),a.setXYZ(28,e.r,e.g,e.b),a.setXYZ(29,e.r,e.g,e.b),a.setXYZ(30,e.r,e.g,e.b),a.setXYZ(31,e.r,e.g,e.b),a.setXYZ(32,n.r,n.g,n.b),a.setXYZ(33,n.r,n.g,n.b),a.setXYZ(34,n.r,n.g,n.b),a.setXYZ(35,n.r,n.g,n.b),a.setXYZ(36,n.r,n.g,n.b),a.setXYZ(37,n.r,n.g,n.b),a.setXYZ(38,i.r,i.g,i.b),a.setXYZ(39,i.r,i.g,i.b),a.setXYZ(40,s.r,s.g,s.b),a.setXYZ(41,s.r,s.g,s.b),a.setXYZ(42,s.r,s.g,s.b),a.setXYZ(43,s.r,s.g,s.b),a.setXYZ(44,s.r,s.g,s.b),a.setXYZ(45,s.r,s.g,s.b),a.setXYZ(46,s.r,s.g,s.b),a.setXYZ(47,s.r,s.g,s.b),a.setXYZ(48,s.r,s.g,s.b),a.setXYZ(49,s.r,s.g,s.b),a.needsUpdate=!0,this}update(){let t=this.geometry,e=this.pointMap,n=1,i=1,s,r;if(Te.projectionMatrixInverse.copy(this.camera.projectionMatrixInverse),this.camera.reversedDepth===!0)s=1,r=0;else if(this.camera.coordinateSystem===2000)s=-1,r=1;else if(this.camera.coordinateSystem===2001)s=0,r=1;else throw Error("THREE.CameraHelper.update(): Invalid coordinate system: "+this.camera.coordinateSystem);Ae("c",e,t,Te,0,0,s),Ae("t",e,t,Te,0,0,r),Ae("n1",e,t,Te,-1,-1,s),Ae("n2",e,t,Te,1,-1,s),Ae("n3",e,t,Te,-1,1,s),Ae("n4",e,t,Te,1,1,s),Ae("f1",e,t,Te,-1,-1,r),Ae("f2",e,t,Te,1,-1,r),Ae("f3",e,t,Te,-1,1,r),Ae("f4",e,t,Te,1,1,r),Ae("u1",e,t,Te,0.7,1.1,s),Ae("u2",e,t,Te,-0.7,1.1,s),Ae("u3",e,t,Te,0,2,s),Ae("cf1",e,t,Te,-1,0,r),Ae("cf2",e,t,Te,1,0,r),Ae("cf3",e,t,Te,0,-1,r),Ae("cf4",e,t,Te,0,1,r),Ae("cn1",e,t,Te,-1,0,s),Ae("cn2",e,t,Te,1,0,s),Ae("cn3",e,t,Te,0,-1,s),Ae("cn4",e,t,Te,0,1,s),t.getAttribute("position").needsUpdate=!0}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}}function Ae(t,e,n,i,s,r,a){ba.set(s,r,a).unproject(i);let o=e[t];if(o!==void 0){let l=n.getAttribute("position");for(let c=0,h=o.length;c<h;c++)l.setXYZ(o[c],ba.x,ba.y,ba.z)}}var Ta=new Oe;class gp extends En{constructor(t,e=16776960){let n=new Uint16Array([0,1,1,2,2,3,3,0,4,5,5,6,6,7,7,4,0,4,1,5,2,6,3,7]),i=new Float32Array(24),s=new jt;s.setIndex(new de(n,1)),s.setAttribute("position",new de(i,3));super(s,new We({color:e,toneMapped:!1}));this.object=t,this.type="BoxHelper",this.matrixAutoUpdate=!1,this.update()}update(){if(this.object!==void 0)Ta.setFromObject(this.object);if(Ta.isEmpty())return;let{min:t,max:e}=Ta,n=this.geometry.attributes.position,i=n.array;i[0]=e.x,i[1]=e.y,i[2]=e.z,i[3]=t.x,i[4]=e.y,i[5]=e.z,i[6]=t.x,i[7]=t.y,i[8]=e.z,i[9]=e.x,i[10]=t.y,i[11]=e.z,i[12]=e.x,i[13]=e.y,i[14]=t.z,i[15]=t.x,i[16]=e.y,i[17]=t.z,i[18]=t.x,i[19]=t.y,i[20]=t.z,i[21]=e.x,i[22]=t.y,i[23]=t.z,n.needsUpdate=!0,this.geometry.computeBoundingSphere()}setFromObject(t){return this.object=t,this.update(),this}copy(t,e){return super.copy(t,e),this.object=t.object,this}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}}class _p extends En{constructor(t,e=16776960){let n=new Uint16Array([0,1,1,2,2,3,3,0,4,5,5,6,6,7,7,4,0,4,1,5,2,6,3,7]),i=[1,1,1,-1,1,1,-1,-1,1,1,-1,1,1,1,-1,-1,1,-1,-1,-1,-1,1,-1,-1],s=new jt;s.setIndex(new de(n,1)),s.setAttribute("position",new Nt(i,3));super(s,new We({color:e,toneMapped:!1}));this.box=t,this.type="Box3Helper",this.geometry.computeBoundingSphere()}updateMatrixWorld(t){let e=this.box;if(e.isEmpty())return;e.getCenter(this.position),e.getSize(this.scale),this.scale.multiplyScalar(0.5),super.updateMatrixWorld(t)}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}}class xp extends Zn{constructor(t,e=1,n=16776960){let i=n,s=[1,-1,0,-1,1,0,-1,-1,0,1,1,0,-1,1,0,-1,-1,0,1,-1,0,1,1,0],r=new jt;r.setAttribute("position",new Nt(s,3)),r.computeBoundingSphere();super(r,new We({color:i,toneMapped:!1}));this.type="PlaneHelper",this.plane=t,this.size=e;let a=[1,1,0,-1,1,0,-1,-1,0,1,1,0,-1,-1,0,1,-1,0],o=new jt;o.setAttribute("position",new Nt(a,3)),o.computeBoundingSphere(),this.add(new pe(o,new mn({color:i,opacity:0.2,transparent:!0,depthWrite:!1,toneMapped:!1})))}updateMatrixWorld(t){this.position.set(0,0,0),this.scale.set(0.5*this.size,0.5*this.size,1),this.lookAt(this.plane.normal),this.translateZ(-this.plane.constant),super.updateMatrixWorld(t)}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose(),this.children[0].geometry.dispose(),this.children[0].material.dispose()}}var td=new L,Ea,vl;class yp extends ae{constructor(t=new L(0,0,1),e=new L(0,0,0),n=1,i=16776960,s=n*0.2,r=s*0.2){super();if(this.type="ArrowHelper",Ea===void 0)Ea=new jt,Ea.setAttribute("position",new Nt([0,0,0,0,1,0],3)),vl=new _r(0.5,1,5,1),vl.translate(0,-0.5,0);this.position.copy(e),this.line=new Zn(Ea,new We({color:i,toneMapped:!1})),this.line.matrixAutoUpdate=!1,this.add(this.line),this.cone=new pe(vl,new mn({color:i,toneMapped:!1})),this.cone.matrixAutoUpdate=!1,this.add(this.cone),this.setDirection(t),this.setLength(n,s,r)}setDirection(t){if(t.y>0.99999)this.quaternion.set(0,0,0,1);else if(t.y<-0.99999)this.quaternion.set(1,0,0,0);else{td.set(t.z,0,-t.x).normalize();let e=Math.acos(t.y);this.quaternion.setFromAxisAngle(td,e)}}setLength(t,e=t*0.2,n=e*0.2){this.line.scale.set(1,Math.max(0.0001,t-e),1),this.line.updateMatrix(),this.cone.scale.set(n,e,n),this.cone.position.y=t,this.cone.updateMatrix()}setColor(t){this.line.material.color.set(t),this.cone.material.color.set(t)}copy(t){return super.copy(t,!1),this.line.copy(t.line),this.cone.copy(t.cone),this}dispose(){super.dispose(),this.line.geometry.dispose(),this.line.material.dispose(),this.cone.geometry.dispose(),this.cone.material.dispose()}}class vp extends En{constructor(t=1){let e=[0,0,0,t,0,0,0,0,0,0,t,0,0,0,0,0,0,t],n=[1,0,0,1,0.6,0,0,1,0,0.6,1,0,0,0,1,0,0.6,1],i=new jt;i.setAttribute("position",new Nt(e,3)),i.setAttribute("color",new Nt(n,3));let s=new We({vertexColors:!0,toneMapped:!1});super(i,s);this.type="AxesHelper"}setColors(t,e,n){let i=new Pt,s=this.geometry.attributes.color.array;return i.set(t),i.toArray(s,0),i.toArray(s,3),i.set(e),i.toArray(s,6),i.toArray(s,9),i.set(n),i.toArray(s,12),i.toArray(s,15),this.geometry.attributes.color.needsUpdate=!0,this}dispose(){super.dispose(),this.geometry.dispose(),this.material.dispose()}}class Un{constructor(){this.type="ShapePath",this.color=new Pt,this.subPaths=[],this.currentPath=null,this.userData={}}moveTo(t,e){return this.currentPath=new Jn,this.subPaths.push(this.currentPath),this.currentPath.moveTo(t,e),this}lineTo(t,e){return this.currentPath.lineTo(t,e),this}quadraticCurveTo(t,e,n,i){return this.currentPath.quadraticCurveTo(t,e,n,i),this}bezierCurveTo(t,e,n,i,s,r){return this.currentPath.bezierCurveTo(t,e,n,i,s,r),this}splineThru(t){return this.currentPath.splineThru(t),this}toShapes(){function t(o,l){let c=!1,h=l.length;for(let d=0,u=h-1;d<h;u=d++){let f=l[d],g=l[u];if(f.y>o.y!==g.y>o.y&&o.x<(g.x-f.x)*(o.y-f.y)/(g.y-f.y)+f.x)c=!c}return c}function e(o,l){let c=l.getCenter(new rt);if(t(c,o))return c;let h=c.y,d=[],u=o.length;for(let f=0;f<u;f++){let g=o[f],_=o[(f+1)%u];if(g.y>h!==_.y>h){let m=g.x+(h-g.y)*(_.x-g.x)/(_.y-g.y);d.push(m)}}if(d.length>1)d.sort((f,g)=>f-g),c.x=(d[0]+d[1])/2;return c}let n=this.userData.style&&this.userData.style.fillRule||"nonzero";if(n!=="nonzero"&&n!=="evenodd")Mt('Fill-rule "'+n+'" is not supported, falling back to "nonzero".'),n="nonzero";let i=n==="nonzero"?(o)=>o!==0:(o)=>(o&1)!==0,s=[];for(let o of this.subPaths){let l=o.getPoints();if(l.length<3)continue;let c=sn.area(l);if(c===0)continue;let h=new Rr;for(let d=0;d<l.length;d++)h.expandByPoint(l[d]);s.push({subPath:o,points:l,boundingBox:h,interiorPoint:e(l,h),absArea:Math.abs(c),winding:c<0?-1:1,container:null,exclude:!1,role:null})}s.sort((o,l)=>l.absArea-o.absArea);for(let o=0;o<s.length;o++){let l=s[o],c=0;for(let h=o-1;h>=0;h--){let d=s[h];if(!d.boundingBox.containsBox(l.boundingBox))continue;if(!t(l.interiorPoint,d.points))continue;l.container=d.exclude?d.container:d,c=d.winding,l.winding+=c;break}if(i(l.winding)===i(c))l.exclude=!0}for(let o of s){if(o.exclude)continue;o.role=o.container===null||o.container.role==="hole"?"outer":"hole"}let r=[],a=new Map;for(let o of s){if(o.exclude||o.role!=="outer")continue;let l=new Cs;l.curves=o.subPath.curves,r.push(l),a.set(o,l)}for(let o of s){if(o.exclude||o.role!=="hole")continue;let l=a.get(o.container);if(!l)continue;let c=new Jn;c.curves=o.subPath.curves,l.holes.push(c)}return r}}class Io extends pn{constructor(t,e=null){super();this.object=t,this.domElement=e,this.enabled=!0,this.state=-1,this.keys={},this.mouseButtons={LEFT:null,MIDDLE:null,RIGHT:null},this.touches={ONE:null,TWO:null}}connect(t){if(this.domElement!==null)this.disconnect();this.domElement=t}disconnect(){}dispose(){}update(){}}function y_(t,e){let n=t.image&&t.image.width?t.image.width/t.image.height:1;if(n>e)t.repeat.x=1,t.repeat.y=n/e,t.offset.x=0,t.offset.y=(1-t.repeat.y)/2;else t.repeat.x=e/n,t.repeat.y=1,t.offset.x=(1-t.repeat.x)/2,t.offset.y=0;return t}function v_(t,e){let n=t.image&&t.image.width?t.image.width/t.image.height:1;if(n>e)t.repeat.x=e/n,t.repeat.y=1,t.offset.x=(1-t.repeat.x)/2,t.offset.y=0;else t.repeat.x=1,t.repeat.y=n/e,t.offset.x=0,t.offset.y=(1-t.repeat.y)/2;return t}function S_(t){return t.repeat.x=1,t.repeat.y=1,t.offset.x=0,t.offset.y=0,t}function Po(t,e,n,i){let s=M_(i);switch(n){case 1021:return t*e;case 1028:return t*e/s.components*s.byteLength;case 1029:return t*e/s.components*s.byteLength;case 1030:return t*e*2/s.components*s.byteLength;case 1031:return t*e*2/s.components*s.byteLength;case 1022:return t*e*3/s.components*s.byteLength;case 1023:return t*e*4/s.components*s.byteLength;case 1033:return t*e*4/s.components*s.byteLength;case 33776:case 33777:return Math.floor((t+3)/4)*Math.floor((e+3)/4)*8;case 33778:case 33779:return Math.floor((t+3)/4)*Math.floor((e+3)/4)*16;case 35841:case 35843:return Math.max(t,16)*Math.max(e,8)/4;case 35840:case 35842:return Math.max(t,8)*Math.max(e,8)/2;case 36196:case 37492:case 37488:case 37489:return Math.floor((t+3)/4)*Math.floor((e+3)/4)*8;case 37496:case 37490:case 37491:return Math.floor((t+3)/4)*Math.floor((e+3)/4)*16;case 37808:return Math.floor((t+3)/4)*Math.floor((e+3)/4)*16;case 37809:return Math.floor((t+4)/5)*Math.floor((e+3)/4)*16;case 37810:return Math.floor((t+4)/5)*Math.floor((e+4)/5)*16;case 37811:return Math.floor((t+5)/6)*Math.floor((e+4)/5)*16;case 37812:return Math.floor((t+5)/6)*Math.floor((e+5)/6)*16;case 37813:return Math.floor((t+7)/8)*Math.floor((e+4)/5)*16;case 37814:return Math.floor((t+7)/8)*Math.floor((e+5)/6)*16;case 37815:return Math.floor((t+7)/8)*Math.floor((e+7)/8)*16;case 37816:return Math.floor((t+9)/10)*Math.floor((e+4)/5)*16;case 37817:return Math.floor((t+9)/10)*Math.floor((e+5)/6)*16;case 37818:return Math.floor((t+9)/10)*Math.floor((e+7)/8)*16;case 37819:return Math.floor((t+9)/10)*Math.floor((e+9)/10)*16;case 37820:return Math.floor((t+11)/12)*Math.floor((e+9)/10)*16;case 37821:return Math.floor((t+11)/12)*Math.floor((e+11)/12)*16;case 36492:case 36494:case 36495:return Math.ceil(t/4)*Math.ceil(e/4)*16;case 36283:case 36284:return Math.ceil(t/4)*Math.ceil(e/4)*8;case 36285:case 36286:return Math.ceil(t/4)*Math.ceil(e/4)*16}throw Error(`Unable to determine texture byte length for ${n} format.`)}function M_(t){switch(t){case 1009:case 1010:return{byteLength:1,components:1};case 1012:case 1011:case 1016:return{byteLength:2,components:1};case 1017:case 1018:return{byteLength:2,components:4};case 1014:case 1013:case 1015:return{byteLength:4,components:1};case 35902:case 35899:return{byteLength:4,components:3}}throw Error(`THREE.TextureUtils: Unknown texture type ${t}.`)}class Sp{static contain(t,e){return y_(t,e)}static cover(t,e){return v_(t,e)}static fill(t){return S_(t)}static getByteLength(t,e,n,i){return Po(t,e,n,i)}}if(typeof __THREE_DEVTOOLS__<"u")__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("register",{detail:{revision:"186"}}));if(typeof window<"u")if(window.__THREE__)Mt("WARNING: Multiple instances of Three.js being imported.");else window.__THREE__="186";function Vp(){let t=null,e=!1,n=null,i=null;function s(r,a){i=t.requestAnimationFrame(s),n(r,a)}return{start:function(){if(e===!0)return;if(n===null)return;if(t===null)return;i=t.requestAnimationFrame(s),e=!0},stop:function(){if(t!==null)t.cancelAnimationFrame(i);e=!1},setAnimationLoop:function(r){n=r},setContext:function(r){t=r}}}function b_(t){let e=new WeakMap;function n(o,l){let{array:c,usage:h}=o,d=c.byteLength,u=t.createBuffer();t.bindBuffer(l,u),t.bufferData(l,c,h),o.onUploadCallback();let f;if(c instanceof Float32Array)f=t.FLOAT;else if(typeof Float16Array<"u"&&c instanceof Float16Array)f=t.HALF_FLOAT;else if(c instanceof Uint16Array)if(o.isFloat16BufferAttribute)f=t.HALF_FLOAT;else f=t.UNSIGNED_SHORT;else if(c instanceof Int16Array)f=t.SHORT;else if(c instanceof Uint32Array)f=t.UNSIGNED_INT;else if(c instanceof Int32Array)f=t.INT;else if(c instanceof Int8Array)f=t.BYTE;else if(c instanceof Uint8Array)f=t.UNSIGNED_BYTE;else if(c instanceof Uint8ClampedArray)f=t.UNSIGNED_BYTE;else throw Error("THREE.WebGLAttributes: Unsupported buffer data format: "+c);return{buffer:u,type:f,bytesPerElement:c.BYTES_PER_ELEMENT,version:o.version,size:d}}function i(o,l,c){let{array:h,updateRanges:d}=l;if(t.bindBuffer(c,o),d.length===0)t.bufferSubData(c,0,h);else{d.sort((f,g)=>f.start-g.start);let u=0;for(let f=1;f<d.length;f++){let g=d[u],_=d[f];if(_.start<=g.start+g.count+1)g.count=Math.max(g.count,_.start+_.count-g.start);else++u,d[u]=_}d.length=u+1;for(let f=0,g=d.length;f<g;f++){let _=d[f];t.bufferSubData(c,_.start*h.BYTES_PER_ELEMENT,h,_.start,_.count)}l.clearUpdateRanges()}l.onUploadCallback()}function s(o){if(o.isInterleavedBufferAttribute)o=o.data;return e.get(o)}function r(o){if(o.isInterleavedBufferAttribute)o=o.data;let l=e.get(o);if(l)t.deleteBuffer(l.buffer),e.delete(o)}function a(o,l){if(o.isInterleavedBufferAttribute)o=o.data;if(o.isGLBufferAttribute){let h=e.get(o);if(!h||h.version<o.version)e.set(o,{buffer:o.buffer,type:o.type,bytesPerElement:o.elementSize,version:o.version});return}let c=e.get(o);if(c===void 0)e.set(o,n(o,l));else if(c.version<o.version){if(c.size!==o.array.byteLength)throw Error("THREE.WebGLAttributes: The size of the buffer attribute's array buffer does not match the original size. Resizing buffer attributes is not supported.");i(c.buffer,o,l),c.version=o.version}}return{get:s,remove:r,update:a}}var T_=`#ifdef USE_ALPHAHASH
	if ( diffuseColor.a < getAlphaHashThreshold( vPosition ) ) discard;
#endif`,E_=`#ifdef USE_ALPHAHASH
	const float ALPHA_HASH_SCALE = 0.05;
	float hash2D( vec2 value ) {
		return fract( 1.0e4 * sin( 17.0 * value.x + 0.1 * value.y ) * ( 0.1 + abs( sin( 13.0 * value.y + value.x ) ) ) );
	}
	float hash3D( vec3 value ) {
		return hash2D( vec2( hash2D( value.xy ), value.z ) );
	}
	float getAlphaHashThreshold( vec3 position ) {
		float maxDeriv = max(
			length( dFdx( position.xyz ) ),
			length( dFdy( position.xyz ) )
		);
		float pixScale = 1.0 / ( ALPHA_HASH_SCALE * maxDeriv );
		vec2 pixScales = vec2(
			exp2( floor( log2( pixScale ) ) ),
			exp2( ceil( log2( pixScale ) ) )
		);
		vec2 alpha = vec2(
			hash3D( floor( pixScales.x * position.xyz ) ),
			hash3D( floor( pixScales.y * position.xyz ) )
		);
		float lerpFactor = fract( log2( pixScale ) );
		float x = ( 1.0 - lerpFactor ) * alpha.x + lerpFactor * alpha.y;
		float a = min( lerpFactor, 1.0 - lerpFactor );
		vec3 cases = vec3(
			x * x / ( 2.0 * a * ( 1.0 - a ) ),
			( x - 0.5 * a ) / ( 1.0 - a ),
			1.0 - ( ( 1.0 - x ) * ( 1.0 - x ) / ( 2.0 * a * ( 1.0 - a ) ) )
		);
		float threshold = ( x < ( 1.0 - a ) )
			? ( ( x < a ) ? cases.x : cases.y )
			: cases.z;
		return clamp( threshold , 1.0e-6, 1.0 );
	}
#endif`,A_=`#ifdef USE_ALPHAMAP
	diffuseColor.a *= texture2D( alphaMap, vAlphaMapUv ).g;
#endif`,w_=`#ifdef USE_ALPHAMAP
	uniform sampler2D alphaMap;
#endif`,R_=`#ifdef USE_ALPHATEST
	#ifdef ALPHA_TO_COVERAGE
	diffuseColor.a = smoothstep( alphaTest, alphaTest + fwidth( diffuseColor.a ), diffuseColor.a );
	if ( diffuseColor.a == 0.0 ) discard;
	#else
	if ( diffuseColor.a < alphaTest ) discard;
	#endif
#endif`,C_=`#ifdef USE_ALPHATEST
	uniform float alphaTest;
#endif`,I_=`#ifdef USE_AOMAP
	float ambientOcclusion = ( texture2D( aoMap, vAoMapUv ).r - 1.0 ) * aoMapIntensity + 1.0;
	reflectedLight.indirectDiffuse *= ambientOcclusion;
	#if defined( USE_CLEARCOAT ) 
		clearcoatSpecularIndirect *= ambientOcclusion;
	#endif
	#if defined( USE_SHEEN ) 
		sheenSpecularIndirect *= ambientOcclusion;
	#endif
	#if defined( USE_ENVMAP ) && defined( STANDARD )
		float dotNV = saturate( dot( geometryNormal, geometryViewDir ) );
		reflectedLight.indirectSpecular *= computeSpecularOcclusion( dotNV, ambientOcclusion, material.roughness );
	#endif
#endif`,P_=`#ifdef USE_AOMAP
	uniform sampler2D aoMap;
	uniform float aoMapIntensity;
#endif`,L_=`#ifdef USE_BATCHING
	#if ! defined( GL_ANGLE_multi_draw )
	#define gl_DrawID _gl_DrawID
	uniform int _gl_DrawID;
	#endif
	uniform highp sampler2D batchingTexture;
	uniform highp usampler2D batchingIdTexture;
	mat4 getBatchingMatrix( const in float i ) {
		int size = textureSize( batchingTexture, 0 ).x;
		int j = int( i ) * 4;
		int x = j % size;
		int y = j / size;
		vec4 v1 = texelFetch( batchingTexture, ivec2( x, y ), 0 );
		vec4 v2 = texelFetch( batchingTexture, ivec2( x + 1, y ), 0 );
		vec4 v3 = texelFetch( batchingTexture, ivec2( x + 2, y ), 0 );
		vec4 v4 = texelFetch( batchingTexture, ivec2( x + 3, y ), 0 );
		return mat4( v1, v2, v3, v4 );
	}
	float getIndirectIndex( const in int i ) {
		int size = textureSize( batchingIdTexture, 0 ).x;
		int x = i % size;
		int y = i / size;
		return float( texelFetch( batchingIdTexture, ivec2( x, y ), 0 ).r );
	}
#endif
#ifdef USE_BATCHING_COLOR
	uniform sampler2D batchingColorTexture;
	vec4 getBatchingColor( const in float i ) {
		int size = textureSize( batchingColorTexture, 0 ).x;
		int j = int( i );
		int x = j % size;
		int y = j / size;
		return texelFetch( batchingColorTexture, ivec2( x, y ), 0 );
	}
#endif`,N_=`#ifdef USE_BATCHING
	mat4 batchingMatrix = getBatchingMatrix( getIndirectIndex( gl_DrawID ) );
#endif`,D_=`vec3 transformed = vec3( position );
#ifdef USE_ALPHAHASH
	vPosition = vec3( position );
#endif`,U_=`vec3 objectNormal = vec3( normal );
#ifdef USE_TANGENT
	vec3 objectTangent = vec3( tangent.xyz );
#endif`,F_=`float G_BlinnPhong_Implicit( ) {
	return 0.25;
}
float D_BlinnPhong( const in float shininess, const in float dotNH ) {
	return RECIPROCAL_PI * ( shininess * 0.5 + 1.0 ) * pow( dotNH, shininess );
}
vec3 BRDF_BlinnPhong( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, const in vec3 specularColor, const in float shininess ) {
	vec3 halfDir = normalize( lightDir + viewDir );
	float dotNH = saturate( dot( normal, halfDir ) );
	float dotVH = saturate( dot( viewDir, halfDir ) );
	vec3 F = F_Schlick( specularColor, 1.0, dotVH );
	float G = G_BlinnPhong_Implicit( );
	float D = D_BlinnPhong( shininess, dotNH );
	return F * ( G * D );
} // validated`,O_=`#ifdef USE_IRIDESCENCE
	const mat3 XYZ_TO_REC709 = mat3(
		 3.2404542, -0.9692660,  0.0556434,
		-1.5371385,  1.8760108, -0.2040259,
		-0.4985314,  0.0415560,  1.0572252
	);
	vec3 Fresnel0ToIor( vec3 fresnel0 ) {
		vec3 sqrtF0 = sqrt( fresnel0 );
		return ( vec3( 1.0 ) + sqrtF0 ) / ( vec3( 1.0 ) - sqrtF0 );
	}
	vec3 IorToFresnel0( vec3 transmittedIor, float incidentIor ) {
		return pow2( ( transmittedIor - vec3( incidentIor ) ) / ( transmittedIor + vec3( incidentIor ) ) );
	}
	float IorToFresnel0( float transmittedIor, float incidentIor ) {
		return pow2( ( transmittedIor - incidentIor ) / ( transmittedIor + incidentIor ));
	}
	vec3 evalSensitivity( float OPD, vec3 shift ) {
		float phase = 2.0 * PI * OPD * 1.0e-9;
		vec3 val = vec3( 5.4856e-13, 4.4201e-13, 5.2481e-13 );
		vec3 pos = vec3( 1.6810e+06, 1.7953e+06, 2.2084e+06 );
		vec3 var = vec3( 4.3278e+09, 9.3046e+09, 6.6121e+09 );
		vec3 xyz = val * sqrt( 2.0 * PI * var ) * cos( pos * phase + shift ) * exp( - pow2( phase ) * var );
		xyz.x += 9.7470e-14 * sqrt( 2.0 * PI * 4.5282e+09 ) * cos( 2.2399e+06 * phase + shift[ 0 ] ) * exp( - 4.5282e+09 * pow2( phase ) );
		xyz /= 1.0685e-7;
		vec3 rgb = XYZ_TO_REC709 * xyz;
		return rgb;
	}
	vec3 evalIridescence( float outsideIOR, float eta2, float cosTheta1, float thinFilmThickness, vec3 baseF0 ) {
		vec3 I;
		float iridescenceIOR = mix( outsideIOR, eta2, smoothstep( 0.0, 0.03, thinFilmThickness ) );
		float sinTheta2Sq = pow2( outsideIOR / iridescenceIOR ) * ( 1.0 - pow2( cosTheta1 ) );
		float cosTheta2Sq = 1.0 - sinTheta2Sq;
		if ( cosTheta2Sq < 0.0 ) {
			return vec3( 1.0 );
		}
		float cosTheta2 = sqrt( cosTheta2Sq );
		float R0 = IorToFresnel0( iridescenceIOR, outsideIOR );
		float R12 = F_Schlick( R0, 1.0, cosTheta1 );
		float T121 = 1.0 - R12;
		float phi12 = 0.0;
		if ( iridescenceIOR < outsideIOR ) phi12 = PI;
		float phi21 = PI - phi12;
		vec3 baseIOR = Fresnel0ToIor( clamp( baseF0, 0.0, 0.9999 ) );		vec3 R1 = IorToFresnel0( baseIOR, iridescenceIOR );
		vec3 R23 = F_Schlick( R1, 1.0, cosTheta2 );
		vec3 phi23 = vec3( 0.0 );
		if ( baseIOR[ 0 ] < iridescenceIOR ) phi23[ 0 ] = PI;
		if ( baseIOR[ 1 ] < iridescenceIOR ) phi23[ 1 ] = PI;
		if ( baseIOR[ 2 ] < iridescenceIOR ) phi23[ 2 ] = PI;
		float OPD = 2.0 * iridescenceIOR * thinFilmThickness * cosTheta2;
		vec3 phi = vec3( phi21 ) + phi23;
		vec3 R123 = clamp( R12 * R23, 1e-5, 0.9999 );
		vec3 r123 = sqrt( R123 );
		vec3 Rs = pow2( T121 ) * R23 / ( vec3( 1.0 ) - R123 );
		vec3 C0 = R12 + Rs;
		I = C0;
		vec3 Cm = Rs - T121;
		for ( int m = 1; m <= 2; ++ m ) {
			Cm *= r123;
			vec3 Sm = 2.0 * evalSensitivity( float( m ) * OPD, float( m ) * phi );
			I += Cm * Sm;
		}
		return max( I, vec3( 0.0 ) );
	}
#endif`,B_=`#ifdef USE_BUMPMAP
	uniform sampler2D bumpMap;
	uniform float bumpScale;
	vec2 dHdxy_fwd() {
		vec2 dSTdx = dFdx( vBumpMapUv );
		vec2 dSTdy = dFdy( vBumpMapUv );
		float Hll = bumpScale * texture2D( bumpMap, vBumpMapUv ).x;
		float dBx = bumpScale * texture2D( bumpMap, vBumpMapUv + dSTdx ).x - Hll;
		float dBy = bumpScale * texture2D( bumpMap, vBumpMapUv + dSTdy ).x - Hll;
		return vec2( dBx, dBy );
	}
	vec3 perturbNormalArb( vec3 surf_pos, vec3 surf_norm, vec2 dHdxy, float faceDirection ) {
		vec3 vSigmaX = normalize( dFdx( surf_pos.xyz ) );
		vec3 vSigmaY = normalize( dFdy( surf_pos.xyz ) );
		vec3 vN = surf_norm;
		vec3 R1 = cross( vSigmaY, vN );
		vec3 R2 = cross( vN, vSigmaX );
		float fDet = dot( vSigmaX, R1 ) * faceDirection;
		vec3 vGrad = sign( fDet ) * ( dHdxy.x * R1 + dHdxy.y * R2 );
		return normalize( abs( fDet ) * surf_norm - vGrad );
	}
#endif`,z_=`#if NUM_CLIPPING_PLANES > 0
	vec4 plane;
	#ifdef ALPHA_TO_COVERAGE
		float distanceToPlane, distanceGradient;
		float clipOpacity = 1.0;
		#pragma unroll_loop_start
		for ( int i = 0; i < UNION_CLIPPING_PLANES; i ++ ) {
			plane = clippingPlanes[ i ];
			distanceToPlane = - dot( vClipPosition, plane.xyz ) + plane.w;
			distanceGradient = fwidth( distanceToPlane ) / 2.0;
			clipOpacity *= smoothstep( - distanceGradient, distanceGradient, distanceToPlane );
			if ( clipOpacity == 0.0 ) discard;
		}
		#pragma unroll_loop_end
		#if UNION_CLIPPING_PLANES < NUM_CLIPPING_PLANES
			float unionClipOpacity = 1.0;
			#pragma unroll_loop_start
			for ( int i = UNION_CLIPPING_PLANES; i < NUM_CLIPPING_PLANES; i ++ ) {
				plane = clippingPlanes[ i ];
				distanceToPlane = - dot( vClipPosition, plane.xyz ) + plane.w;
				distanceGradient = fwidth( distanceToPlane ) / 2.0;
				unionClipOpacity *= 1.0 - smoothstep( - distanceGradient, distanceGradient, distanceToPlane );
			}
			#pragma unroll_loop_end
			clipOpacity *= 1.0 - unionClipOpacity;
		#endif
		diffuseColor.a *= clipOpacity;
		if ( diffuseColor.a == 0.0 ) discard;
	#else
		#pragma unroll_loop_start
		for ( int i = 0; i < UNION_CLIPPING_PLANES; i ++ ) {
			plane = clippingPlanes[ i ];
			if ( dot( vClipPosition, plane.xyz ) > plane.w ) discard;
		}
		#pragma unroll_loop_end
		#if UNION_CLIPPING_PLANES < NUM_CLIPPING_PLANES
			bool clipped = true;
			#pragma unroll_loop_start
			for ( int i = UNION_CLIPPING_PLANES; i < NUM_CLIPPING_PLANES; i ++ ) {
				plane = clippingPlanes[ i ];
				clipped = ( dot( vClipPosition, plane.xyz ) > plane.w ) && clipped;
			}
			#pragma unroll_loop_end
			if ( clipped ) discard;
		#endif
	#endif
#endif`,k_=`#if NUM_CLIPPING_PLANES > 0
	varying vec3 vClipPosition;
	uniform vec4 clippingPlanes[ NUM_CLIPPING_PLANES ];
#endif`,G_=`#if NUM_CLIPPING_PLANES > 0
	varying vec3 vClipPosition;
#endif`,H_=`#if NUM_CLIPPING_PLANES > 0
	vClipPosition = - mvPosition.xyz;
#endif`,V_=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA )
	diffuseColor *= vColor;
#endif`,W_=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA )
	varying vec4 vColor;
#endif`,X_=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA ) || defined( USE_INSTANCING_COLOR ) || defined( USE_BATCHING_COLOR )
	varying vec4 vColor;
#endif`,q_=`#if defined( USE_COLOR ) || defined( USE_COLOR_ALPHA ) || defined( USE_INSTANCING_COLOR ) || defined( USE_BATCHING_COLOR )
	vColor = vec4( 1.0 );
#endif
#ifdef USE_COLOR_ALPHA
	vColor *= color;
#elif defined( USE_COLOR )
	vColor.rgb *= color;
#endif
#ifdef USE_INSTANCING_COLOR
	vColor.rgb *= instanceColor.rgb;
#endif
#ifdef USE_BATCHING_COLOR
	vColor *= getBatchingColor( getIndirectIndex( gl_DrawID ) );
#endif`,Y_=`#define PI 3.141592653589793
#define PI2 6.283185307179586
#define PI_HALF 1.5707963267948966
#define RECIPROCAL_PI 0.3183098861837907
#define RECIPROCAL_PI2 0.15915494309189535
#define EPSILON 1e-6
#ifndef saturate
#define saturate( a ) clamp( a, 0.0, 1.0 )
#endif
#define whiteComplement( a ) ( 1.0 - saturate( a ) )
float pow2( const in float x ) { return x*x; }
vec3 pow2( const in vec3 x ) { return x*x; }
float pow3( const in float x ) { return x*x*x; }
float pow4( const in float x ) { float x2 = x*x; return x2*x2; }
float max3( const in vec3 v ) { return max( max( v.x, v.y ), v.z ); }
float average( const in vec3 v ) { return dot( v, vec3( 0.3333333 ) ); }
highp float rand( const in vec2 uv ) {
	const highp float a = 12.9898, b = 78.233, c = 43758.5453;
	highp float dt = dot( uv.xy, vec2( a,b ) ), sn = mod( dt, PI );
	return fract( sin( sn ) * c );
}
#ifdef HIGH_PRECISION
	float precisionSafeLength( vec3 v ) { return length( v ); }
#else
	float precisionSafeLength( vec3 v ) {
		float maxComponent = max3( abs( v ) );
		return length( v / maxComponent ) * maxComponent;
	}
#endif
struct IncidentLight {
	vec3 color;
	vec3 direction;
	bool visible;
};
struct ReflectedLight {
	vec3 directDiffuse;
	vec3 directSpecular;
	vec3 indirectDiffuse;
	vec3 indirectSpecular;
};
#ifdef USE_ALPHAHASH
	varying vec3 vPosition;
#endif
vec3 transformDirection( in vec3 dir, in mat4 matrix ) {
	return normalize( ( matrix * vec4( dir, 0.0 ) ).xyz );
}
#define inverseTransformDirection transformDirectionByInverseViewMatrix
vec3 transformNormalByInverseViewMatrix( in vec3 normal, in mat4 viewMatrix ) {
	return normalize( ( vec4( normal, 0.0 ) * viewMatrix ).xyz );
}
vec3 transformDirectionByInverseViewMatrix( in vec3 dir, in mat4 viewMatrix ) {
	return normalize( ( vec4( dir, 0.0 ) * viewMatrix ).xyz );
}
bool isPerspectiveMatrix( mat4 m ) {
	return m[ 2 ][ 3 ] == - 1.0;
}
vec2 equirectUv( in vec3 dir ) {
	float u = atan( dir.z, dir.x ) * RECIPROCAL_PI2 + 0.5;
	float v = asin( clamp( dir.y, - 1.0, 1.0 ) ) * RECIPROCAL_PI + 0.5;
	return vec2( u, v );
}
vec3 BRDF_Lambert( const in vec3 diffuseColor ) {
	return RECIPROCAL_PI * diffuseColor;
}
vec3 F_Schlick( const in vec3 f0, const in float f90, const in float dotVH ) {
	float fresnel = exp2( ( - 5.55473 * dotVH - 6.98316 ) * dotVH );
	return f0 * ( 1.0 - fresnel ) + ( f90 * fresnel );
}
float F_Schlick( const in float f0, const in float f90, const in float dotVH ) {
	float fresnel = exp2( ( - 5.55473 * dotVH - 6.98316 ) * dotVH );
	return f0 * ( 1.0 - fresnel ) + ( f90 * fresnel );
} // validated`,Z_=`#ifdef ENVMAP_TYPE_CUBE_UV
	#define cubeUV_minMipLevel 4.0
	#define cubeUV_minTileSize 16.0
	float getFace( vec3 direction ) {
		vec3 absDirection = abs( direction );
		float face = - 1.0;
		if ( absDirection.x > absDirection.z ) {
			if ( absDirection.x > absDirection.y )
				face = direction.x > 0.0 ? 0.0 : 3.0;
			else
				face = direction.y > 0.0 ? 1.0 : 4.0;
		} else {
			if ( absDirection.z > absDirection.y )
				face = direction.z > 0.0 ? 2.0 : 5.0;
			else
				face = direction.y > 0.0 ? 1.0 : 4.0;
		}
		return face;
	}
	vec2 getUV( vec3 direction, float face ) {
		vec2 uv;
		if ( face == 0.0 ) {
			uv = vec2( direction.z, direction.y ) / abs( direction.x );
		} else if ( face == 1.0 ) {
			uv = vec2( - direction.x, - direction.z ) / abs( direction.y );
		} else if ( face == 2.0 ) {
			uv = vec2( - direction.x, direction.y ) / abs( direction.z );
		} else if ( face == 3.0 ) {
			uv = vec2( - direction.z, direction.y ) / abs( direction.x );
		} else if ( face == 4.0 ) {
			uv = vec2( - direction.x, direction.z ) / abs( direction.y );
		} else {
			uv = vec2( direction.x, direction.y ) / abs( direction.z );
		}
		return 0.5 * ( uv + 1.0 );
	}
	vec3 bilinearCubeUV( sampler2D envMap, vec3 direction, float mipInt ) {
		float face = getFace( direction );
		float filterInt = max( cubeUV_minMipLevel - mipInt, 0.0 );
		mipInt = max( mipInt, cubeUV_minMipLevel );
		float faceSize = exp2( mipInt );
		highp vec2 uv = getUV( direction, face ) * ( faceSize - 2.0 ) + 1.0;
		if ( face > 2.0 ) {
			uv.y += faceSize;
			face -= 3.0;
		}
		uv.x += face * faceSize;
		uv.x += filterInt * 3.0 * cubeUV_minTileSize;
		uv.y += 4.0 * ( exp2( CUBEUV_MAX_MIP ) - faceSize );
		uv.x *= CUBEUV_TEXEL_WIDTH;
		uv.y *= CUBEUV_TEXEL_HEIGHT;
		#ifdef texture2DGradEXT
			return texture2DGradEXT( envMap, uv, vec2( 0.0 ), vec2( 0.0 ) ).rgb;
		#else
			return texture2D( envMap, uv ).rgb;
		#endif
	}
	#define cubeUV_r0 1.0
	#define cubeUV_m0 - 2.0
	#define cubeUV_r1 0.8
	#define cubeUV_m1 - 1.0
	#define cubeUV_r4 0.4
	#define cubeUV_m4 2.0
	#define cubeUV_r5 0.305
	#define cubeUV_m5 3.0
	#define cubeUV_r6 0.21
	#define cubeUV_m6 4.0
	float roughnessToMip( float roughness ) {
		float mip = 0.0;
		if ( roughness >= cubeUV_r1 ) {
			mip = ( cubeUV_r0 - roughness ) * ( cubeUV_m1 - cubeUV_m0 ) / ( cubeUV_r0 - cubeUV_r1 ) + cubeUV_m0;
		} else if ( roughness >= cubeUV_r4 ) {
			mip = ( cubeUV_r1 - roughness ) * ( cubeUV_m4 - cubeUV_m1 ) / ( cubeUV_r1 - cubeUV_r4 ) + cubeUV_m1;
		} else if ( roughness >= cubeUV_r5 ) {
			mip = ( cubeUV_r4 - roughness ) * ( cubeUV_m5 - cubeUV_m4 ) / ( cubeUV_r4 - cubeUV_r5 ) + cubeUV_m4;
		} else if ( roughness >= cubeUV_r6 ) {
			mip = ( cubeUV_r5 - roughness ) * ( cubeUV_m6 - cubeUV_m5 ) / ( cubeUV_r5 - cubeUV_r6 ) + cubeUV_m5;
		} else {
			mip = - 2.0 * log2( 1.16 * roughness );		}
		return mip;
	}
	vec4 textureCubeUV( sampler2D envMap, vec3 sampleDir, float roughness ) {
		float mip = clamp( roughnessToMip( roughness ), cubeUV_m0, CUBEUV_MAX_MIP );
		float mipF = fract( mip );
		float mipInt = floor( mip );
		vec3 color0 = bilinearCubeUV( envMap, sampleDir, mipInt );
		if ( mipF == 0.0 ) {
			return vec4( color0, 1.0 );
		} else {
			vec3 color1 = bilinearCubeUV( envMap, sampleDir, mipInt + 1.0 );
			return vec4( mix( color0, color1, mipF ), 1.0 );
		}
	}
#endif`,J_=`vec3 transformedNormal = objectNormal;
#ifdef USE_TANGENT
	vec3 transformedTangent = objectTangent;
#endif
#ifdef USE_BATCHING
	mat3 bm = mat3( batchingMatrix );
	transformedNormal /= vec3( dot( bm[ 0 ], bm[ 0 ] ), dot( bm[ 1 ], bm[ 1 ] ), dot( bm[ 2 ], bm[ 2 ] ) );
	transformedNormal = bm * transformedNormal;
	#ifdef USE_TANGENT
		transformedTangent = bm * transformedTangent;
	#endif
#endif
#ifdef USE_INSTANCING
	mat3 im = mat3( instanceMatrix );
	transformedNormal /= vec3( dot( im[ 0 ], im[ 0 ] ), dot( im[ 1 ], im[ 1 ] ), dot( im[ 2 ], im[ 2 ] ) );
	transformedNormal = im * transformedNormal;
	#ifdef USE_TANGENT
		transformedTangent = im * transformedTangent;
	#endif
#endif
transformedNormal = normalMatrix * transformedNormal;
#ifdef FLIP_SIDED
	transformedNormal = - transformedNormal;
#endif
#ifdef USE_TANGENT
	transformedTangent = ( modelViewMatrix * vec4( transformedTangent, 0.0 ) ).xyz;
#endif`,$_=`#ifdef USE_DISPLACEMENTMAP
	uniform sampler2D displacementMap;
	uniform float displacementScale;
	uniform float displacementBias;
#endif`,K_=`#ifdef USE_DISPLACEMENTMAP
	transformed += normalize( objectNormal ) * ( texture2D( displacementMap, vDisplacementMapUv ).x * displacementScale + displacementBias );
#endif`,j_=`#ifdef USE_EMISSIVEMAP
	vec4 emissiveColor = texture2D( emissiveMap, vEmissiveMapUv );
	#ifdef DECODE_VIDEO_TEXTURE_EMISSIVE
		emissiveColor = sRGBTransferEOTF( emissiveColor );
	#endif
	totalEmissiveRadiance *= emissiveColor.rgb;
#endif`,Q_=`#ifdef USE_EMISSIVEMAP
	uniform sampler2D emissiveMap;
#endif`,tx="gl_FragColor = linearToOutputTexel( gl_FragColor );",ex=`vec4 LinearTransferOETF( in vec4 value ) {
	return value;
}
vec4 sRGBTransferEOTF( in vec4 value ) {
	return vec4( mix( pow( value.rgb * 0.9478672986 + vec3( 0.0521327014 ), vec3( 2.4 ) ), value.rgb * 0.0773993808, vec3( lessThanEqual( value.rgb, vec3( 0.04045 ) ) ) ), value.a );
}
vec4 sRGBTransferOETF( in vec4 value ) {
	return vec4( mix( pow( value.rgb, vec3( 0.41666 ) ) * 1.055 - vec3( 0.055 ), value.rgb * 12.92, vec3( lessThanEqual( value.rgb, vec3( 0.0031308 ) ) ) ), value.a );
}`,nx=`#ifdef USE_ENVMAP
	#ifdef ENV_WORLDPOS
		vec3 cameraToFrag;
		if ( isOrthographic ) {
			cameraToFrag = normalize( vec3( - viewMatrix[ 0 ][ 2 ], - viewMatrix[ 1 ][ 2 ], - viewMatrix[ 2 ][ 2 ] ) );
		} else {
			cameraToFrag = normalize( vWorldPosition - cameraPosition );
		}
		vec3 worldNormal = transformNormalByInverseViewMatrix( normal, viewMatrix );
		#ifdef ENVMAP_MODE_REFLECTION
			vec3 reflectVec = reflect( cameraToFrag, worldNormal );
		#else
			vec3 reflectVec = refract( cameraToFrag, worldNormal, refractionRatio );
		#endif
	#else
		vec3 reflectVec = vReflect;
	#endif
	#ifdef ENVMAP_TYPE_CUBE
		vec4 envColor = textureCube( envMap, envMapRotation * reflectVec );
		#ifdef ENVMAP_BLENDING_MULTIPLY
			outgoingLight = mix( outgoingLight, outgoingLight * envColor.xyz, specularStrength * reflectivity );
		#elif defined( ENVMAP_BLENDING_MIX )
			outgoingLight = mix( outgoingLight, envColor.xyz, specularStrength * reflectivity );
		#elif defined( ENVMAP_BLENDING_ADD )
			outgoingLight += envColor.xyz * specularStrength * reflectivity;
		#endif
	#endif
#endif`,ix=`#ifdef USE_ENVMAP
	uniform float envMapIntensity;
	uniform mat3 envMapRotation;
	#ifdef ENVMAP_TYPE_CUBE
		uniform samplerCube envMap;
	#else
		uniform sampler2D envMap;
	#endif
#endif`,sx=`#ifdef USE_ENVMAP
	uniform float reflectivity;
	#if defined( USE_BUMPMAP ) || defined( USE_NORMALMAP ) || defined( PHONG ) || defined( LAMBERT )
		#define ENV_WORLDPOS
	#endif
	#ifdef ENV_WORLDPOS
		varying vec3 vWorldPosition;
		uniform float refractionRatio;
	#else
		varying vec3 vReflect;
	#endif
#endif`,rx=`#ifdef USE_ENVMAP
	#if defined( USE_BUMPMAP ) || defined( USE_NORMALMAP ) || defined( PHONG ) || defined( LAMBERT )
		#define ENV_WORLDPOS
	#endif
	#ifdef ENV_WORLDPOS
		
		varying vec3 vWorldPosition;
	#else
		varying vec3 vReflect;
		uniform float refractionRatio;
	#endif
#endif`,ax=`#ifdef USE_ENVMAP
	#ifdef ENV_WORLDPOS
		vWorldPosition = worldPosition.xyz;
	#else
		vec3 cameraToVertex;
		if ( isOrthographic ) {
			cameraToVertex = normalize( vec3( - viewMatrix[ 0 ][ 2 ], - viewMatrix[ 1 ][ 2 ], - viewMatrix[ 2 ][ 2 ] ) );
		} else {
			cameraToVertex = normalize( worldPosition.xyz - cameraPosition );
		}
		vec3 worldNormal = transformNormalByInverseViewMatrix( transformedNormal, viewMatrix );
		#ifdef ENVMAP_MODE_REFLECTION
			vReflect = reflect( cameraToVertex, worldNormal );
		#else
			vReflect = refract( cameraToVertex, worldNormal, refractionRatio );
		#endif
	#endif
#endif`,ox=`#ifdef USE_FOG
	vFogDepth = - mvPosition.z;
#endif`,lx=`#ifdef USE_FOG
	varying float vFogDepth;
#endif`,cx=`#ifdef USE_FOG
	#ifdef FOG_EXP2
		float fogFactor = 1.0 - exp( - fogDensity * fogDensity * vFogDepth * vFogDepth );
	#else
		float fogFactor = smoothstep( fogNear, fogFar, vFogDepth );
	#endif
	gl_FragColor.rgb = mix( gl_FragColor.rgb, fogColor, fogFactor );
#endif`,hx=`#ifdef USE_FOG
	uniform vec3 fogColor;
	varying float vFogDepth;
	#ifdef FOG_EXP2
		uniform float fogDensity;
	#else
		uniform float fogNear;
		uniform float fogFar;
	#endif
#endif`,ux=`#ifdef USE_GRADIENTMAP
	uniform sampler2D gradientMap;
#endif
vec3 getGradientIrradiance( vec3 normal, vec3 lightDirection ) {
	float dotNL = dot( normal, lightDirection );
	vec2 coord = vec2( dotNL * 0.5 + 0.5, 0.0 );
	#ifdef USE_GRADIENTMAP
		return vec3( texture2D( gradientMap, coord ).r );
	#else
		vec2 fw = fwidth( coord ) * 0.5;
		return mix( vec3( 0.7 ), vec3( 1.0 ), smoothstep( 0.7 - fw.x, 0.7 + fw.x, coord.x ) );
	#endif
}`,dx=`#ifdef USE_LIGHTMAP
	uniform sampler2D lightMap;
	uniform float lightMapIntensity;
#endif`,fx=`LambertMaterial material;
material.diffuseColor = diffuseColor.rgb;
material.specularStrength = specularStrength;`,px=`varying vec3 vViewPosition;
struct LambertMaterial {
	vec3 diffuseColor;
	float specularStrength;
};
void RE_Direct_Lambert( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in LambertMaterial material, inout ReflectedLight reflectedLight ) {
	float dotNL = saturate( dot( geometryNormal, directLight.direction ) );
	vec3 irradiance = dotNL * directLight.color;
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
void RE_IndirectDiffuse_Lambert( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in LambertMaterial material, inout ReflectedLight reflectedLight ) {
	reflectedLight.indirectDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
#define RE_Direct				RE_Direct_Lambert
#define RE_IndirectDiffuse		RE_IndirectDiffuse_Lambert`,mx=`uniform bool receiveShadow;
uniform vec3 ambientLightColor;
#if defined( USE_LIGHT_PROBES )
	uniform vec3 lightProbe[ 9 ];
#endif
vec3 shGetIrradianceAt( in vec3 normal, in vec3 shCoefficients[ 9 ] ) {
	float x = normal.x, y = normal.y, z = normal.z;
	vec3 result = shCoefficients[ 0 ] * 0.886227;
	result += shCoefficients[ 1 ] * 2.0 * 0.511664 * y;
	result += shCoefficients[ 2 ] * 2.0 * 0.511664 * z;
	result += shCoefficients[ 3 ] * 2.0 * 0.511664 * x;
	result += shCoefficients[ 4 ] * 2.0 * 0.429043 * x * y;
	result += shCoefficients[ 5 ] * 2.0 * 0.429043 * y * z;
	result += shCoefficients[ 6 ] * ( 0.743125 * z * z - 0.247708 );
	result += shCoefficients[ 7 ] * 2.0 * 0.429043 * x * z;
	result += shCoefficients[ 8 ] * 0.429043 * ( x * x - y * y );
	return result;
}
vec3 getLightProbeIrradiance( const in vec3 lightProbe[ 9 ], const in vec3 normal ) {
	vec3 worldNormal = transformNormalByInverseViewMatrix( normal, viewMatrix );
	vec3 irradiance = shGetIrradianceAt( worldNormal, lightProbe );
	return irradiance;
}
vec3 getAmbientLightIrradiance( const in vec3 ambientLightColor ) {
	vec3 irradiance = ambientLightColor;
	return irradiance;
}
float getDistanceAttenuation( const in float lightDistance, const in float cutoffDistance, const in float decayExponent ) {
	float distanceFalloff = 1.0 / max( pow( lightDistance, decayExponent ), 0.01 );
	if ( cutoffDistance > 0.0 ) {
		distanceFalloff *= pow2( saturate( 1.0 - pow4( lightDistance / cutoffDistance ) ) );
	}
	return distanceFalloff;
}
float getSpotAttenuation( const in float coneCosine, const in float penumbraCosine, const in float angleCosine ) {
	return smoothstep( coneCosine, penumbraCosine, angleCosine );
}
#if NUM_SUN_LIGHTS > 0
	struct SunLight {
		vec3 direction;
		vec3 color;
	};
	uniform SunLight sunLights[ NUM_SUN_LIGHTS ];
	void getSunLightInfo( const in SunLight sunLight, out IncidentLight light ) {
		light.color = sunLight.color;
		light.direction = sunLight.direction;
		light.visible = true;
	}
#endif
#if NUM_DIR_LIGHTS > 0
	struct DirectionalLight {
		vec3 direction;
		vec3 color;
	};
	uniform DirectionalLight directionalLights[ NUM_DIR_LIGHTS ];
	void getDirectionalLightInfo( const in DirectionalLight directionalLight, out IncidentLight light ) {
		light.color = directionalLight.color;
		light.direction = directionalLight.direction;
		light.visible = true;
	}
#endif
#if NUM_POINT_LIGHTS > 0
	struct PointLight {
		vec3 position;
		vec3 color;
		float distance;
		float decay;
	};
	uniform PointLight pointLights[ NUM_POINT_LIGHTS ];
	void getPointLightInfo( const in PointLight pointLight, const in vec3 geometryPosition, out IncidentLight light ) {
		vec3 lVector = pointLight.position - geometryPosition;
		light.direction = normalize( lVector );
		float lightDistance = length( lVector );
		light.color = pointLight.color;
		light.color *= getDistanceAttenuation( lightDistance, pointLight.distance, pointLight.decay );
		light.visible = ( light.color != vec3( 0.0 ) );
	}
#endif
#if NUM_SPOT_LIGHTS > 0
	struct SpotLight {
		vec3 position;
		vec3 direction;
		vec3 color;
		float distance;
		float decay;
		float coneCos;
		float penumbraCos;
	};
	uniform SpotLight spotLights[ NUM_SPOT_LIGHTS ];
	void getSpotLightInfo( const in SpotLight spotLight, const in vec3 geometryPosition, out IncidentLight light ) {
		vec3 lVector = spotLight.position - geometryPosition;
		light.direction = normalize( lVector );
		float angleCos = dot( light.direction, spotLight.direction );
		float spotAttenuation = getSpotAttenuation( spotLight.coneCos, spotLight.penumbraCos, angleCos );
		if ( spotAttenuation > 0.0 ) {
			float lightDistance = length( lVector );
			light.color = spotLight.color * spotAttenuation;
			light.color *= getDistanceAttenuation( lightDistance, spotLight.distance, spotLight.decay );
			light.visible = ( light.color != vec3( 0.0 ) );
		} else {
			light.color = vec3( 0.0 );
			light.visible = false;
		}
	}
#endif
#if NUM_RECT_AREA_LIGHTS > 0
	struct RectAreaLight {
		vec3 color;
		vec3 position;
		vec3 halfWidth;
		vec3 halfHeight;
	};
	uniform sampler2D ltc_1;	uniform sampler2D ltc_2;
	uniform RectAreaLight rectAreaLights[ NUM_RECT_AREA_LIGHTS ];
#endif
#if NUM_HEMI_LIGHTS > 0
	struct HemisphereLight {
		vec3 direction;
		vec3 skyColor;
		vec3 groundColor;
	};
	uniform HemisphereLight hemisphereLights[ NUM_HEMI_LIGHTS ];
	vec3 getHemisphereLightIrradiance( const in HemisphereLight hemiLight, const in vec3 normal ) {
		float dotNL = dot( normal, hemiLight.direction );
		float hemiDiffuseWeight = 0.5 * dotNL + 0.5;
		vec3 irradiance = mix( hemiLight.groundColor, hemiLight.skyColor, hemiDiffuseWeight );
		return irradiance;
	}
#endif
#include <lightprobes_pars_fragment>`,gx=`#ifdef USE_ENVMAP
	vec3 getIBLIrradiance( const in vec3 normal ) {
		#ifdef ENVMAP_TYPE_CUBE_UV
			vec3 worldNormal = transformNormalByInverseViewMatrix( normal, viewMatrix );
			vec4 envMapColor = textureCubeUV( envMap, envMapRotation * worldNormal, 1.0 );
			return PI * envMapColor.rgb * envMapIntensity;
		#else
			return vec3( 0.0 );
		#endif
	}
	vec3 getIBLRadiance( const in vec3 viewDir, const in vec3 normal, const in float roughness ) {
		#ifdef ENVMAP_TYPE_CUBE_UV
			vec3 reflectVec = reflect( - viewDir, normal );
			reflectVec = normalize( mix( reflectVec, normal, pow4( roughness ) ) );
			reflectVec = transformDirectionByInverseViewMatrix( reflectVec, viewMatrix );
			vec4 envMapColor = textureCubeUV( envMap, envMapRotation * reflectVec, roughness );
			return envMapColor.rgb * envMapIntensity;
		#else
			return vec3( 0.0 );
		#endif
	}
	#ifdef USE_RETROREFLECTION
		vec3 getIBLRetroRadiance( const in vec3 viewDir, const in vec3 normal, const in float roughness ) {
			#ifdef ENVMAP_TYPE_CUBE_UV
				vec3 retroVec = normalize( mix( viewDir, normal, pow4( roughness ) ) );
				retroVec = transformDirectionByInverseViewMatrix( retroVec, viewMatrix );
				vec4 envMapColor = textureCubeUV( envMap, envMapRotation * retroVec, roughness );
				return envMapColor.rgb * envMapIntensity;
			#else
				return vec3( 0.0 );
			#endif
		}
	#endif
	#ifdef USE_ANISOTROPY
		vec3 getIBLAnisotropyRadiance( const in vec3 viewDir, const in vec3 normal, const in float roughness, const in vec3 bitangent, const in float anisotropy ) {
			#ifdef ENVMAP_TYPE_CUBE_UV
				vec3 bentNormal = cross( bitangent, viewDir );
				bentNormal = normalize( cross( bentNormal, bitangent ) );
				bentNormal = normalize( mix( bentNormal, normal, pow2( pow2( 1.0 - anisotropy * ( 1.0 - roughness ) ) ) ) );
				return getIBLRadiance( viewDir, bentNormal, roughness );
			#else
				return vec3( 0.0 );
			#endif
		}
		#ifdef USE_RETROREFLECTION
			vec3 getIBLAnisotropyRetroRadiance( const in vec3 viewDir, const in vec3 normal, const in float roughness, const in vec3 bitangent, const in float anisotropy ) {
				#ifdef ENVMAP_TYPE_CUBE_UV
					vec3 bentNormal = cross( bitangent, viewDir );
					bentNormal = normalize( cross( bentNormal, bitangent ) );
					bentNormal = normalize( mix( bentNormal, normal, pow2( pow2( 1.0 - anisotropy * ( 1.0 - roughness ) ) ) ) );
					return getIBLRetroRadiance( viewDir, bentNormal, roughness );
				#else
					return vec3( 0.0 );
				#endif
			}
		#endif
	#endif
#endif`,_x=`ToonMaterial material;
material.diffuseColor = diffuseColor.rgb;`,xx=`varying vec3 vViewPosition;
struct ToonMaterial {
	vec3 diffuseColor;
};
void RE_Direct_Toon( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in ToonMaterial material, inout ReflectedLight reflectedLight ) {
	vec3 irradiance = getGradientIrradiance( geometryNormal, directLight.direction ) * directLight.color;
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
void RE_IndirectDiffuse_Toon( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in ToonMaterial material, inout ReflectedLight reflectedLight ) {
	reflectedLight.indirectDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
#define RE_Direct				RE_Direct_Toon
#define RE_IndirectDiffuse		RE_IndirectDiffuse_Toon`,yx=`BlinnPhongMaterial material;
material.diffuseColor = diffuseColor.rgb;
material.specularColor = specular;
material.specularShininess = shininess;
material.specularStrength = specularStrength;`,vx=`varying vec3 vViewPosition;
struct BlinnPhongMaterial {
	vec3 diffuseColor;
	vec3 specularColor;
	float specularShininess;
	float specularStrength;
};
void RE_Direct_BlinnPhong( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in BlinnPhongMaterial material, inout ReflectedLight reflectedLight ) {
	float dotNL = saturate( dot( geometryNormal, directLight.direction ) );
	vec3 irradiance = dotNL * directLight.color;
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
	reflectedLight.directSpecular += irradiance * BRDF_BlinnPhong( directLight.direction, geometryViewDir, geometryNormal, material.specularColor, material.specularShininess ) * material.specularStrength;
}
void RE_IndirectDiffuse_BlinnPhong( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in BlinnPhongMaterial material, inout ReflectedLight reflectedLight ) {
	reflectedLight.indirectDiffuse += irradiance * BRDF_Lambert( material.diffuseColor );
}
#define RE_Direct				RE_Direct_BlinnPhong
#define RE_IndirectDiffuse		RE_IndirectDiffuse_BlinnPhong`,Sx=`PhysicalMaterial material;
material.diffuseColor = diffuseColor.rgb;
material.diffuseContribution = diffuseColor.rgb * ( 1.0 - metalnessFactor );
material.metalness = metalnessFactor;
vec3 dxy = max( abs( dFdx( nonPerturbedNormal ) ), abs( dFdy( nonPerturbedNormal ) ) );
float geometryRoughness = max( max( dxy.x, dxy.y ), dxy.z );
material.roughness = max( roughnessFactor, 0.0525 );material.roughness += geometryRoughness;
material.roughness = min( material.roughness, 1.0 );
#ifdef IOR
	material.ior = ior;
	#ifdef USE_SPECULAR
		float specularIntensityFactor = specularIntensity;
		vec3 specularColorFactor = specularColor;
		#ifdef USE_SPECULAR_COLORMAP
			specularColorFactor *= texture2D( specularColorMap, vSpecularColorMapUv ).rgb;
		#endif
		#ifdef USE_SPECULAR_INTENSITYMAP
			specularIntensityFactor *= texture2D( specularIntensityMap, vSpecularIntensityMapUv ).a;
		#endif
		material.specularF90 = mix( specularIntensityFactor, 1.0, metalnessFactor );
	#else
		float specularIntensityFactor = 1.0;
		vec3 specularColorFactor = vec3( 1.0 );
		material.specularF90 = 1.0;
	#endif
	material.specularColor = min( pow2( ( material.ior - 1.0 ) / ( material.ior + 1.0 ) ) * specularColorFactor, vec3( 1.0 ) ) * specularIntensityFactor;
	material.specularColorBlended = mix( material.specularColor, diffuseColor.rgb, metalnessFactor );
#else
	material.specularColor = vec3( 0.04 );
	material.specularColorBlended = mix( material.specularColor, diffuseColor.rgb, metalnessFactor );
	material.specularF90 = 1.0;
#endif
#ifdef USE_CLEARCOAT
	material.clearcoat = clearcoat;
	material.clearcoatRoughness = clearcoatRoughness;
	material.clearcoatF0 = vec3( 0.04 );
	material.clearcoatF90 = 1.0;
	#ifdef USE_CLEARCOATMAP
		material.clearcoat *= texture2D( clearcoatMap, vClearcoatMapUv ).x;
	#endif
	#ifdef USE_CLEARCOAT_ROUGHNESSMAP
		material.clearcoatRoughness *= texture2D( clearcoatRoughnessMap, vClearcoatRoughnessMapUv ).y;
	#endif
	material.clearcoat = saturate( material.clearcoat );	material.clearcoatRoughness = max( material.clearcoatRoughness, 0.0525 );
	material.clearcoatRoughness += geometryRoughness;
	material.clearcoatRoughness = min( material.clearcoatRoughness, 1.0 );
#endif
#ifdef USE_DISPERSION
	material.dispersion = dispersion;
#endif
#ifdef USE_RETROREFLECTION
	material.retroreflectivity = retroreflectivity;
#endif
#ifdef USE_IRIDESCENCE
	material.iridescence = iridescence;
	material.iridescenceIOR = iridescenceIOR;
	#ifdef USE_IRIDESCENCEMAP
		material.iridescence *= texture2D( iridescenceMap, vIridescenceMapUv ).r;
	#endif
	#ifdef USE_IRIDESCENCE_THICKNESSMAP
		material.iridescenceThickness = (iridescenceThicknessMaximum - iridescenceThicknessMinimum) * texture2D( iridescenceThicknessMap, vIridescenceThicknessMapUv ).g + iridescenceThicknessMinimum;
	#else
		material.iridescenceThickness = iridescenceThicknessMaximum;
	#endif
#endif
#ifdef USE_SHEEN
	material.sheenColor = sheenColor;
	#ifdef USE_SHEEN_COLORMAP
		material.sheenColor *= texture2D( sheenColorMap, vSheenColorMapUv ).rgb;
	#endif
	material.sheenRoughness = clamp( sheenRoughness, 0.0001, 1.0 );
	#ifdef USE_SHEEN_ROUGHNESSMAP
		material.sheenRoughness *= texture2D( sheenRoughnessMap, vSheenRoughnessMapUv ).a;
	#endif
#endif
#ifdef USE_ANISOTROPY
	#ifdef USE_ANISOTROPYMAP
		mat2 anisotropyMat = mat2( anisotropyVector.x, anisotropyVector.y, - anisotropyVector.y, anisotropyVector.x );
		vec3 anisotropyPolar = texture2D( anisotropyMap, vAnisotropyMapUv ).rgb;
		vec2 anisotropyV = anisotropyMat * normalize( 2.0 * anisotropyPolar.rg - vec2( 1.0 ) ) * anisotropyPolar.b;
	#else
		vec2 anisotropyV = anisotropyVector;
	#endif
	material.anisotropy = length( anisotropyV );
	if( material.anisotropy == 0.0 ) {
		anisotropyV = vec2( 1.0, 0.0 );
	} else {
		anisotropyV /= material.anisotropy;
		material.anisotropy = saturate( material.anisotropy );
	}
	material.alphaT = mix( pow2( material.roughness ), 1.0, pow2( material.anisotropy ) );
	material.anisotropyT = tbn[ 0 ] * anisotropyV.x + tbn[ 1 ] * anisotropyV.y;
	material.anisotropyB = tbn[ 1 ] * anisotropyV.x - tbn[ 0 ] * anisotropyV.y;
#endif`,Mx=`uniform sampler2D dfgLUT;
struct PhysicalMaterial {
	vec3 diffuseColor;
	vec3 diffuseContribution;
	vec3 specularColor;
	vec3 specularColorBlended;
	float roughness;
	float metalness;
	float specularF90;
	float dispersion;
	vec2 dfg;
	vec3 multiScatteringCompensation;
	#ifdef USE_RETROREFLECTION
		float retroreflectivity;
	#endif
	#ifdef USE_CLEARCOAT
		float clearcoat;
		float clearcoatRoughness;
		vec3 clearcoatF0;
		float clearcoatF90;
	#endif
	#ifdef USE_IRIDESCENCE
		float iridescence;
		float iridescenceIOR;
		float iridescenceThickness;
		vec3 iridescenceFresnel;
		vec3 iridescenceF0Dielectric;
		vec3 iridescenceF0Metallic;
	#endif
	#ifdef USE_SHEEN
		vec3 sheenColor;
		float sheenRoughness;
	#endif
	#ifdef IOR
		float ior;
	#endif
	#ifdef USE_TRANSMISSION
		float transmission;
		float transmissionAlpha;
		float thickness;
		float attenuationDistance;
		vec3 attenuationColor;
	#endif
	#ifdef USE_ANISOTROPY
		float anisotropy;
		float alphaT;
		vec3 anisotropyT;
		vec3 anisotropyB;
	#endif
};
vec3 clearcoatSpecularDirect = vec3( 0.0 );
vec3 clearcoatSpecularIndirect = vec3( 0.0 );
vec3 sheenSpecularDirect = vec3( 0.0 );
vec3 sheenSpecularIndirect = vec3(0.0 );
vec3 Schlick_to_F0( const in vec3 f, const in float f90, const in float dotVH ) {
    float x = clamp( 1.0 - dotVH, 0.0, 1.0 );
    float x2 = x * x;
    float x5 = clamp( x * x2 * x2, 0.0, 0.9999 );
    return ( f - vec3( f90 ) * x5 ) / ( 1.0 - x5 );
}
float V_GGX_SmithCorrelated( const in float alpha, const in float dotNL, const in float dotNV ) {
	float a2 = pow2( alpha );
	float gv = dotNL * sqrt( a2 + ( 1.0 - a2 ) * pow2( dotNV ) );
	float gl = dotNV * sqrt( a2 + ( 1.0 - a2 ) * pow2( dotNL ) );
	return 0.5 / max( gv + gl, EPSILON );
}
float D_GGX( const in float alpha, const in float dotNH ) {
	float a2 = pow2( alpha );
	float denom = pow2( dotNH ) * ( a2 - 1.0 ) + 1.0;
	return RECIPROCAL_PI * a2 / pow2( denom );
}
#ifdef USE_ANISOTROPY
	float V_GGX_SmithCorrelated_Anisotropic( const in float alphaT, const in float alphaB, const in float dotTV, const in float dotBV, const in float dotTL, const in float dotBL, const in float dotNV, const in float dotNL ) {
		float gv = dotNL * length( vec3( alphaT * dotTV, alphaB * dotBV, dotNV ) );
		float gl = dotNV * length( vec3( alphaT * dotTL, alphaB * dotBL, dotNL ) );
		return 0.5 / max( gv + gl, EPSILON );
	}
	float D_GGX_Anisotropic( const in float alphaT, const in float alphaB, const in float dotNH, const in float dotTH, const in float dotBH ) {
		float a2 = alphaT * alphaB;
		highp vec3 v = vec3( alphaB * dotTH, alphaT * dotBH, a2 * dotNH );
		highp float v2 = dot( v, v );
		float w2 = a2 / v2;
		return RECIPROCAL_PI * a2 * pow2 ( w2 );
	}
#endif
#ifdef USE_CLEARCOAT
	vec3 BRDF_GGX_Clearcoat( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, const in PhysicalMaterial material) {
		vec3 f0 = material.clearcoatF0;
		float f90 = material.clearcoatF90;
		float roughness = material.clearcoatRoughness;
		float alpha = pow2( roughness );
		vec3 halfDir = normalize( lightDir + viewDir );
		float dotNL = saturate( dot( normal, lightDir ) );
		float dotNV = saturate( dot( normal, viewDir ) );
		float dotNH = saturate( dot( normal, halfDir ) );
		float dotVH = saturate( dot( viewDir, halfDir ) );
		vec3 F = F_Schlick( f0, f90, dotVH );
		float V = V_GGX_SmithCorrelated( alpha, dotNL, dotNV );
		float D = D_GGX( alpha, dotNH );
		return F * ( V * D );
	}
#endif
vec3 BRDF_GGX( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, const in PhysicalMaterial material ) {
	vec3 f0 = material.specularColorBlended;
	float f90 = material.specularF90;
	float roughness = material.roughness;
	float alpha = pow2( roughness );
	vec3 halfDir = normalize( lightDir + viewDir );
	float dotNL = saturate( dot( normal, lightDir ) );
	float dotNV = saturate( dot( normal, viewDir ) );
	float dotNH = saturate( dot( normal, halfDir ) );
	float dotVH = saturate( dot( viewDir, halfDir ) );
	vec3 F = F_Schlick( f0, f90, dotVH );
	#ifdef USE_IRIDESCENCE
		F = mix( F, material.iridescenceFresnel, material.iridescence );
	#endif
	#ifdef USE_ANISOTROPY
		float dotTL = dot( material.anisotropyT, lightDir );
		float dotTV = dot( material.anisotropyT, viewDir );
		float dotTH = dot( material.anisotropyT, halfDir );
		float dotBL = dot( material.anisotropyB, lightDir );
		float dotBV = dot( material.anisotropyB, viewDir );
		float dotBH = dot( material.anisotropyB, halfDir );
		float V = V_GGX_SmithCorrelated_Anisotropic( material.alphaT, alpha, dotTV, dotBV, dotTL, dotBL, dotNV, dotNL );
		float D = D_GGX_Anisotropic( material.alphaT, alpha, dotNH, dotTH, dotBH );
	#else
		float V = V_GGX_SmithCorrelated( alpha, dotNL, dotNV );
		float D = D_GGX( alpha, dotNH );
	#endif
	return F * ( V * D );
}
vec2 LTC_Uv( const in vec3 N, const in vec3 V, const in float roughness ) {
	const float LUT_SIZE = 64.0;
	const float LUT_SCALE = ( LUT_SIZE - 1.0 ) / LUT_SIZE;
	const float LUT_BIAS = 0.5 / LUT_SIZE;
	float dotNV = saturate( dot( N, V ) );
	vec2 uv = vec2( roughness, sqrt( 1.0 - dotNV ) );
	uv = uv * LUT_SCALE + LUT_BIAS;
	return uv;
}
float LTC_ClippedSphereFormFactor( const in vec3 f ) {
	float l = length( f );
	return max( ( l * l + f.z ) / ( l + 1.0 ), 0.0 );
}
vec3 LTC_EdgeVectorFormFactor( const in vec3 v1, const in vec3 v2 ) {
	float x = dot( v1, v2 );
	float y = abs( x );
	float a = 0.8543985 + ( 0.4965155 + 0.0145206 * y ) * y;
	float b = 3.4175940 + ( 4.1616724 + y ) * y;
	float v = a / b;
	float theta_sintheta = ( x > 0.0 ) ? v : 0.5 * inversesqrt( max( 1.0 - x * x, 1e-7 ) ) - v;
	return cross( v1, v2 ) * theta_sintheta;
}
vec3 LTC_Evaluate( const in vec3 N, const in vec3 V, const in vec3 P, const in mat3 mInv, const in vec3 rectCoords[ 4 ] ) {
	vec3 v1 = rectCoords[ 1 ] - rectCoords[ 0 ];
	vec3 v2 = rectCoords[ 3 ] - rectCoords[ 0 ];
	vec3 lightNormal = cross( v1, v2 );
	if( dot( lightNormal, P - rectCoords[ 0 ] ) < 0.0 ) return vec3( 0.0 );
	vec3 T1, T2;
	T1 = normalize( V - N * dot( V, N ) );
	T2 = - cross( N, T1 );
	mat3 mat = mInv * transpose( mat3( T1, T2, N ) );
	vec3 coords[ 4 ];
	coords[ 0 ] = mat * ( rectCoords[ 0 ] - P );
	coords[ 1 ] = mat * ( rectCoords[ 1 ] - P );
	coords[ 2 ] = mat * ( rectCoords[ 2 ] - P );
	coords[ 3 ] = mat * ( rectCoords[ 3 ] - P );
	coords[ 0 ] = normalize( coords[ 0 ] );
	coords[ 1 ] = normalize( coords[ 1 ] );
	coords[ 2 ] = normalize( coords[ 2 ] );
	coords[ 3 ] = normalize( coords[ 3 ] );
	vec3 vectorFormFactor = vec3( 0.0 );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 0 ], coords[ 1 ] );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 1 ], coords[ 2 ] );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 2 ], coords[ 3 ] );
	vectorFormFactor += LTC_EdgeVectorFormFactor( coords[ 3 ], coords[ 0 ] );
	float result = LTC_ClippedSphereFormFactor( vectorFormFactor );
	return vec3( result );
}
#if defined( USE_SHEEN )
float D_Charlie( float roughness, float dotNH ) {
	float alpha = pow2( roughness );
	float invAlpha = 1.0 / alpha;
	float cos2h = dotNH * dotNH;
	float sin2h = max( 1.0 - cos2h, 0.0078125 );
	return ( 2.0 + invAlpha ) * pow( sin2h, invAlpha * 0.5 ) / ( 2.0 * PI );
}
float V_Neubelt( float dotNV, float dotNL ) {
	return saturate( 1.0 / ( 4.0 * ( dotNL + dotNV - dotNL * dotNV ) ) );
}
vec3 BRDF_Sheen( const in vec3 lightDir, const in vec3 viewDir, const in vec3 normal, vec3 sheenColor, const in float sheenRoughness ) {
	vec3 halfDir = normalize( lightDir + viewDir );
	float dotNL = saturate( dot( normal, lightDir ) );
	float dotNV = saturate( dot( normal, viewDir ) );
	float dotNH = saturate( dot( normal, halfDir ) );
	float D = D_Charlie( sheenRoughness, dotNH );
	float V = V_Neubelt( dotNV, dotNL );
	return sheenColor * ( D * V );
}
#endif
float IBLSheenBRDF( const in vec3 normal, const in vec3 viewDir, const in float roughness ) {
	float dotNV = saturate( dot( normal, viewDir ) );
	float r2 = roughness * roughness;
	float rInv = 1.0 / ( roughness + 0.1 );
	float a = -1.9362 + 1.0678 * roughness + 0.4573 * r2 - 0.8469 * rInv;
	float b = -0.6014 + 0.5538 * roughness - 0.4670 * r2 - 0.1255 * rInv;
	float DG = exp( a * dotNV + b );
	return saturate( DG );
}
vec3 EnvironmentBRDF( const in vec3 normal, const in vec3 viewDir, const in vec3 specularColor, const in float specularF90, const in float roughness ) {
	float dotNV = saturate( dot( normal, viewDir ) );
	vec2 fab = texture2D( dfgLUT, vec2( roughness, dotNV ) ).rg;
	return specularColor * fab.x + specularF90 * fab.y;
}
#ifdef USE_IRIDESCENCE
void computeMultiscatteringIridescence( const in vec2 fab, const in vec3 specularColor, const in float specularF90, const in float iridescence, const in vec3 iridescenceF0, inout vec3 singleScatter, inout vec3 multiScatter ) {
#else
void computeMultiscattering( const in vec2 fab, const in vec3 specularColor, const in float specularF90, inout vec3 singleScatter, inout vec3 multiScatter ) {
#endif
	#ifdef USE_IRIDESCENCE
		vec3 Fr = mix( specularColor, iridescenceF0, iridescence );
	#else
		vec3 Fr = specularColor;
	#endif
	vec3 FssEss = Fr * fab.x + specularF90 * fab.y;
	float Ess = fab.x + fab.y;
	float Ems = 1.0 - Ess;
	vec3 Favg = Fr + ( 1.0 - Fr ) * 0.047619;	vec3 Fms = FssEss * Favg / ( 1.0 - Ems * Favg );
	singleScatter += FssEss;
	multiScatter += Fms * Ems;
}
#if NUM_RECT_AREA_LIGHTS > 0
	void RE_Direct_RectArea_Physical( const in RectAreaLight rectAreaLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight ) {
		vec3 normal = geometryNormal;
		vec3 viewDir = geometryViewDir;
		vec3 position = geometryPosition;
		vec3 lightPos = rectAreaLight.position;
		vec3 halfWidth = rectAreaLight.halfWidth;
		vec3 halfHeight = rectAreaLight.halfHeight;
		vec3 lightColor = rectAreaLight.color;
		float roughness = material.roughness;
		vec3 rectCoords[ 4 ];
		rectCoords[ 0 ] = lightPos + halfWidth - halfHeight;		rectCoords[ 1 ] = lightPos - halfWidth - halfHeight;
		rectCoords[ 2 ] = lightPos - halfWidth + halfHeight;
		rectCoords[ 3 ] = lightPos + halfWidth + halfHeight;
		vec2 uv = LTC_Uv( normal, viewDir, roughness );
		vec4 t1 = texture2D( ltc_1, uv );
		vec4 t2 = texture2D( ltc_2, uv );
		mat3 mInv = mat3(
			vec3( t1.x, 0, t1.y ),
			vec3(    0, 1,    0 ),
			vec3( t1.z, 0, t1.w )
		);
		vec3 fresnel = ( material.specularColorBlended * t2.x + ( material.specularF90 - material.specularColorBlended ) * t2.y );
		reflectedLight.directSpecular += lightColor * fresnel * LTC_Evaluate( normal, viewDir, position, mInv, rectCoords );
		reflectedLight.directDiffuse += lightColor * material.diffuseContribution * LTC_Evaluate( normal, viewDir, position, mat3( 1.0 ), rectCoords );
		#ifdef USE_CLEARCOAT
			vec3 Ncc = geometryClearcoatNormal;
			vec2 uvClearcoat = LTC_Uv( Ncc, viewDir, material.clearcoatRoughness );
			vec4 t1Clearcoat = texture2D( ltc_1, uvClearcoat );
			vec4 t2Clearcoat = texture2D( ltc_2, uvClearcoat );
			mat3 mInvClearcoat = mat3(
				vec3( t1Clearcoat.x, 0, t1Clearcoat.y ),
				vec3(             0, 1,             0 ),
				vec3( t1Clearcoat.z, 0, t1Clearcoat.w )
			);
			vec3 fresnelClearcoat = material.clearcoatF0 * t2Clearcoat.x + ( material.clearcoatF90 - material.clearcoatF0 ) * t2Clearcoat.y;
			clearcoatSpecularDirect += lightColor * fresnelClearcoat * LTC_Evaluate( Ncc, viewDir, position, mInvClearcoat, rectCoords );
		#endif
	}
#endif
void RE_Direct_Physical( const in IncidentLight directLight, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight ) {
	float dotNL = saturate( dot( geometryNormal, directLight.direction ) );
	vec3 irradiance = dotNL * directLight.color;
	#ifdef USE_CLEARCOAT
		float dotNLcc = saturate( dot( geometryClearcoatNormal, directLight.direction ) );
		vec3 ccIrradiance = dotNLcc * directLight.color;
		clearcoatSpecularDirect += ccIrradiance * BRDF_GGX_Clearcoat( directLight.direction, geometryViewDir, geometryClearcoatNormal, material );
	#endif
	#ifdef USE_SHEEN
 
 		sheenSpecularDirect += irradiance * BRDF_Sheen( directLight.direction, geometryViewDir, geometryNormal, material.sheenColor, material.sheenRoughness );
 
 		float sheenAlbedoV = IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness );
 		float sheenAlbedoL = IBLSheenBRDF( geometryNormal, directLight.direction, material.sheenRoughness );
 
 		float sheenEnergyComp = 1.0 - max3( material.sheenColor ) * max( sheenAlbedoV, sheenAlbedoL );
 
 		irradiance *= sheenEnergyComp;
 
 	#endif
	vec3 specularBRDF = BRDF_GGX( directLight.direction, geometryViewDir, geometryNormal, material );
	#ifdef USE_RETROREFLECTION
		vec3 retroViewDir = reflect( - geometryViewDir, geometryNormal );
		vec3 retroSpecularBRDF = BRDF_GGX( directLight.direction, retroViewDir, geometryNormal, material );
		specularBRDF = mix( specularBRDF, retroSpecularBRDF, saturate( material.retroreflectivity ) );
	#endif
	reflectedLight.directSpecular += irradiance * specularBRDF * material.multiScatteringCompensation;
	vec3 halfDir = normalize( directLight.direction + geometryViewDir );
	float dotVH = saturate( dot( geometryViewDir, halfDir ) );
	vec3 F = F_Schlick( material.specularColor, material.specularF90, dotVH );
	#ifdef USE_RETROREFLECTION
		vec3 retroHalfDir = normalize( directLight.direction + retroViewDir );
		float dotRetroVH = saturate( dot( retroViewDir, retroHalfDir ) );
		vec3 retroF = F_Schlick( material.specularColor, material.specularF90, dotRetroVH );
		F = mix( F, retroF, saturate( material.retroreflectivity ) );
	#endif
	reflectedLight.directDiffuse += irradiance * BRDF_Lambert( material.diffuseContribution ) * ( 1.0 - F );
}
void RE_IndirectDiffuse_Physical( const in vec3 irradiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight ) {
	vec3 singleScattering = vec3( 0.0 );
	vec3 multiScattering = vec3( 0.0 );
	#ifdef USE_IRIDESCENCE
		computeMultiscatteringIridescence( material.dfg, material.specularColor, material.specularF90, material.iridescence, material.iridescenceF0Dielectric, singleScattering, multiScattering );
	#else
		computeMultiscattering( material.dfg, material.specularColor, material.specularF90, singleScattering, multiScattering );
	#endif
	vec3 diffuse = irradiance * BRDF_Lambert( material.diffuseContribution ) * ( 1.0 - singleScattering - multiScattering );
	#ifdef USE_SHEEN
		float sheenAlbedo = IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness );
		sheenSpecularIndirect += irradiance * material.sheenColor * sheenAlbedo * RECIPROCAL_PI;
		float sheenEnergyComp = 1.0 - max3( material.sheenColor ) * sheenAlbedo;
		diffuse *= sheenEnergyComp;
	#endif
	reflectedLight.indirectDiffuse += diffuse;
}
void RE_IndirectSpecular_Physical( const in vec3 radiance, const in vec3 irradiance, const in vec3 clearcoatRadiance, const in vec3 geometryPosition, const in vec3 geometryNormal, const in vec3 geometryViewDir, const in vec3 geometryClearcoatNormal, const in PhysicalMaterial material, inout ReflectedLight reflectedLight) {
	#ifdef USE_CLEARCOAT
		clearcoatSpecularIndirect += clearcoatRadiance * EnvironmentBRDF( geometryClearcoatNormal, geometryViewDir, material.clearcoatF0, material.clearcoatF90, material.clearcoatRoughness );
	#endif
	#ifdef USE_SHEEN
		sheenSpecularIndirect += irradiance * material.sheenColor * IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness ) * RECIPROCAL_PI;
 	#endif
	vec3 singleScatteringDielectric = vec3( 0.0 );
	vec3 multiScatteringDielectric = vec3( 0.0 );
	vec3 singleScatteringMetallic = vec3( 0.0 );
	vec3 multiScatteringMetallic = vec3( 0.0 );
	#ifdef USE_IRIDESCENCE
		computeMultiscatteringIridescence( material.dfg, material.specularColor, material.specularF90, material.iridescence, material.iridescenceF0Dielectric, singleScatteringDielectric, multiScatteringDielectric );
		computeMultiscatteringIridescence( material.dfg, material.diffuseColor, material.specularF90, material.iridescence, material.iridescenceF0Metallic, singleScatteringMetallic, multiScatteringMetallic );
	#else
		computeMultiscattering( material.dfg, material.specularColor, material.specularF90, singleScatteringDielectric, multiScatteringDielectric );
		computeMultiscattering( material.dfg, material.diffuseColor, material.specularF90, singleScatteringMetallic, multiScatteringMetallic );
	#endif
	vec3 singleScattering = mix( singleScatteringDielectric, singleScatteringMetallic, material.metalness );
	vec3 multiScattering = mix( multiScatteringDielectric, multiScatteringMetallic, material.metalness );
	vec3 totalScatteringDielectric = singleScatteringDielectric + multiScatteringDielectric;
	vec3 diffuse = material.diffuseContribution * ( 1.0 - totalScatteringDielectric );
	vec3 cosineWeightedIrradiance = irradiance * RECIPROCAL_PI;
	vec3 indirectSpecular = radiance * singleScattering;
	indirectSpecular += multiScattering * cosineWeightedIrradiance;
	vec3 indirectDiffuse = diffuse * cosineWeightedIrradiance;
	#ifdef USE_SHEEN
		float sheenAlbedo = IBLSheenBRDF( geometryNormal, geometryViewDir, material.sheenRoughness );
		float sheenEnergyComp = 1.0 - max3( material.sheenColor ) * sheenAlbedo;
		indirectSpecular *= sheenEnergyComp;
		indirectDiffuse *= sheenEnergyComp;
	#endif
	reflectedLight.indirectSpecular += indirectSpecular;
	reflectedLight.indirectDiffuse += indirectDiffuse;
}
#define RE_Direct				RE_Direct_Physical
#define RE_Direct_RectArea		RE_Direct_RectArea_Physical
#define RE_IndirectDiffuse		RE_IndirectDiffuse_Physical
#define RE_IndirectSpecular		RE_IndirectSpecular_Physical
float computeSpecularOcclusion( const in float dotNV, const in float ambientOcclusion, const in float roughness ) {
	return saturate( pow( dotNV + ambientOcclusion, exp2( - 16.0 * roughness - 1.0 ) ) - 1.0 + ambientOcclusion );
}`,bx=`
vec3 geometryPosition = - vViewPosition;
vec3 geometryNormal = normal;
vec3 geometryViewDir = ( isOrthographic ) ? vec3( 0, 0, 1 ) : normalize( vViewPosition );
vec3 geometryClearcoatNormal = vec3( 0.0 );
#ifdef USE_CLEARCOAT
	geometryClearcoatNormal = clearcoatNormal;
#endif
#ifdef USE_IRIDESCENCE
	float dotNVi = saturate( dot( normal, geometryViewDir ) );
	if ( material.iridescenceThickness == 0.0 ) {
		material.iridescence = 0.0;
	} else {
		material.iridescence = saturate( material.iridescence );
	}
	if ( material.iridescence > 0.0 ) {
		vec3 iridescenceFresnelDielectric = evalIridescence( 1.0, material.iridescenceIOR, dotNVi, material.iridescenceThickness, material.specularColor );
		vec3 iridescenceFresnelMetallic = evalIridescence( 1.0, material.iridescenceIOR, dotNVi, material.iridescenceThickness, material.diffuseColor );
		material.iridescenceFresnel = mix( iridescenceFresnelDielectric, iridescenceFresnelMetallic, material.metalness );
		material.iridescenceF0Dielectric = Schlick_to_F0( iridescenceFresnelDielectric, 1.0, dotNVi );
		material.iridescenceF0Metallic = Schlick_to_F0( iridescenceFresnelMetallic, 1.0, dotNVi );
	}
#endif
#ifdef STANDARD
	float dotNVms = saturate( dot( geometryNormal, geometryViewDir ) );
	material.dfg = texture2D( dfgLUT, vec2( material.roughness, dotNVms ) ).rg;
	#if ( NUM_SUN_LIGHTS > 0 || NUM_DIR_LIGHTS > 0 || NUM_POINT_LIGHTS > 0 || NUM_SPOT_LIGHTS > 0 )
		float EssMs = material.dfg.x + material.dfg.y;
		material.multiScatteringCompensation = 1.0 + material.specularColorBlended * ( 1.0 / EssMs - 1.0 );
	#endif
#endif
IncidentLight directLight;
#if ( NUM_POINT_LIGHTS > 0 ) && defined( RE_Direct )
	PointLight pointLight;
	#if defined( USE_SHADOWMAP ) && NUM_POINT_LIGHT_SHADOWS > 0
	PointLightShadow pointLightShadow;
	#endif
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_POINT_LIGHTS; i ++ ) {
		pointLight = pointLights[ i ];
		getPointLightInfo( pointLight, geometryPosition, directLight );
		#if defined( USE_SHADOWMAP ) && ( UNROLLED_LOOP_INDEX < NUM_POINT_LIGHT_SHADOWS ) && ( defined( SHADOWMAP_TYPE_PCF ) || defined( SHADOWMAP_TYPE_BASIC ) )
		pointLightShadow = pointLightShadows[ i ];
		directLight.color *= ( directLight.visible && receiveShadow ) ? getPointShadow( pointShadowMap[ i ], pointLightShadow.shadowMapSize, pointLightShadow.shadowIntensity, pointLightShadow.shadowBias, pointLightShadow.shadowRadius, vPointShadowCoord[ i ], pointLightShadow.shadowCameraNear, pointLightShadow.shadowCameraFar ) : 1.0;
		#endif
		RE_Direct( directLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if ( NUM_SPOT_LIGHTS > 0 ) && defined( RE_Direct )
	SpotLight spotLight;
	vec4 spotColor;
	vec3 spotLightCoord;
	bool inSpotLightMap;
	#if defined( USE_SHADOWMAP ) && NUM_SPOT_LIGHT_SHADOWS > 0
	SpotLightShadow spotLightShadow;
	#endif
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SPOT_LIGHTS; i ++ ) {
		spotLight = spotLights[ i ];
		getSpotLightInfo( spotLight, geometryPosition, directLight );
		#if ( UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS_WITH_MAPS )
		#define SPOT_LIGHT_MAP_INDEX UNROLLED_LOOP_INDEX
		#elif ( UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS )
		#define SPOT_LIGHT_MAP_INDEX NUM_SPOT_LIGHT_MAPS
		#else
		#define SPOT_LIGHT_MAP_INDEX ( UNROLLED_LOOP_INDEX - NUM_SPOT_LIGHT_SHADOWS + NUM_SPOT_LIGHT_SHADOWS_WITH_MAPS )
		#endif
		#if ( SPOT_LIGHT_MAP_INDEX < NUM_SPOT_LIGHT_MAPS )
			spotLightCoord = vSpotLightCoord[ i ].xyz / vSpotLightCoord[ i ].w;
			inSpotLightMap = all( lessThan( abs( spotLightCoord * 2. - 1. ), vec3( 1.0 ) ) );
			spotColor = texture2D( spotLightMap[ SPOT_LIGHT_MAP_INDEX ], spotLightCoord.xy );
			directLight.color = inSpotLightMap ? directLight.color * spotColor.rgb : directLight.color;
		#endif
		#undef SPOT_LIGHT_MAP_INDEX
		#if defined( USE_SHADOWMAP ) && ( UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS )
		spotLightShadow = spotLightShadows[ i ];
		directLight.color *= ( directLight.visible && receiveShadow ) ? getShadow( spotShadowMap[ i ], spotLightShadow.shadowMapSize, spotLightShadow.shadowIntensity, spotLightShadow.shadowBias, spotLightShadow.shadowRadius, vSpotLightCoord[ i ] ) : 1.0;
		#endif
		RE_Direct( directLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if ( NUM_SUN_LIGHTS > 0 ) && defined( RE_Direct )
	SunLight sunLight;
	#if defined( USE_SHADOWMAP ) && NUM_SUN_LIGHT_SHADOWS > 0
	SunLightShadow sunLightShadow;
	#endif
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SUN_LIGHTS; i ++ ) {
		sunLight = sunLights[ i ];
		getSunLightInfo( sunLight, directLight );
		#if defined( USE_SHADOWMAP ) && ( UNROLLED_LOOP_INDEX < NUM_SUN_LIGHT_SHADOWS )
		sunLightShadow = sunLightShadows[ i ];
		directLight.color *= ( directLight.visible && receiveShadow ) ? getSunShadow( sunShadowMap[ i ], sunLightShadow, UNROLLED_LOOP_INDEX ) : 1.0;
		#endif
		RE_Direct( directLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if ( NUM_DIR_LIGHTS > 0 ) && defined( RE_Direct )
	DirectionalLight directionalLight;
	#if defined( USE_SHADOWMAP ) && NUM_DIR_LIGHT_SHADOWS > 0
	DirectionalLightShadow directionalLightShadow;
	#endif
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_DIR_LIGHTS; i ++ ) {
		directionalLight = directionalLights[ i ];
		getDirectionalLightInfo( directionalLight, directLight );
		#if defined( USE_SHADOWMAP ) && ( UNROLLED_LOOP_INDEX < NUM_DIR_LIGHT_SHADOWS )
		directionalLightShadow = directionalLightShadows[ i ];
		directLight.color *= ( directLight.visible && receiveShadow ) ? getShadow( directionalShadowMap[ i ], directionalLightShadow.shadowMapSize, directionalLightShadow.shadowIntensity, directionalLightShadow.shadowBias, directionalLightShadow.shadowRadius, vDirectionalShadowCoord[ i ] ) : 1.0;
		#endif
		RE_Direct( directLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if ( NUM_RECT_AREA_LIGHTS > 0 ) && defined( RE_Direct_RectArea )
	RectAreaLight rectAreaLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_RECT_AREA_LIGHTS; i ++ ) {
		rectAreaLight = rectAreaLights[ i ];
		RE_Direct_RectArea( rectAreaLight, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
	}
	#pragma unroll_loop_end
#endif
#if defined( RE_IndirectDiffuse )
	vec3 iblIrradiance = vec3( 0.0 );
	vec3 irradiance = getAmbientLightIrradiance( ambientLightColor );
	#if defined( USE_LIGHT_PROBES )
		irradiance += getLightProbeIrradiance( lightProbe, geometryNormal );
	#endif
	#if ( NUM_HEMI_LIGHTS > 0 )
		#pragma unroll_loop_start
		for ( int i = 0; i < NUM_HEMI_LIGHTS; i ++ ) {
			irradiance += getHemisphereLightIrradiance( hemisphereLights[ i ], geometryNormal );
		}
		#pragma unroll_loop_end
	#endif
	#ifdef USE_LIGHT_PROBES_GRID
		vec3 probeWorldPos = ( ( vec4( geometryPosition, 1.0 ) - viewMatrix[ 3 ] ) * viewMatrix ).xyz;
		vec3 probeWorldNormal = transformNormalByInverseViewMatrix( geometryNormal, viewMatrix );
		irradiance += getLightProbeGridIrradiance( probeWorldPos, probeWorldNormal );
	#endif
#endif
#if defined( RE_IndirectSpecular )
	vec3 radiance = vec3( 0.0 );
	vec3 clearcoatRadiance = vec3( 0.0 );
#endif`,Tx=`#if defined( RE_IndirectDiffuse )
	#ifdef USE_LIGHTMAP
		vec4 lightMapTexel = texture2D( lightMap, vLightMapUv );
		vec3 lightMapIrradiance = lightMapTexel.rgb * lightMapIntensity;
		irradiance += lightMapIrradiance;
	#endif
	#if defined( USE_ENVMAP ) && defined( ENVMAP_TYPE_CUBE_UV )
		#if defined( STANDARD ) || defined( LAMBERT ) || defined( PHONG )
			iblIrradiance += getIBLIrradiance( geometryNormal );
		#endif
	#endif
#endif
#if defined( USE_ENVMAP ) && defined( RE_IndirectSpecular )
	#ifdef USE_ANISOTROPY
		vec3 iblRadiance = getIBLAnisotropyRadiance( geometryViewDir, geometryNormal, material.roughness, material.anisotropyB, material.anisotropy );
	#else
		vec3 iblRadiance = getIBLRadiance( geometryViewDir, geometryNormal, material.roughness );
	#endif
	#ifdef USE_RETROREFLECTION
		#ifdef USE_ANISOTROPY
			vec3 retroIBLRadiance = getIBLAnisotropyRetroRadiance( geometryViewDir, geometryNormal, material.roughness, material.anisotropyB, material.anisotropy );
		#else
			vec3 retroIBLRadiance = getIBLRetroRadiance( geometryViewDir, geometryNormal, material.roughness );
		#endif
		iblRadiance = mix( iblRadiance, retroIBLRadiance, saturate( material.retroreflectivity ) );
	#endif
	radiance += iblRadiance;
	#ifdef USE_CLEARCOAT
		clearcoatRadiance += getIBLRadiance( geometryViewDir, geometryClearcoatNormal, material.clearcoatRoughness );
	#endif
#endif`,Ex=`#if defined( RE_IndirectDiffuse )
	#if defined( LAMBERT ) || defined( PHONG )
		irradiance += iblIrradiance;
	#endif
	RE_IndirectDiffuse( irradiance, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
#endif
#if defined( RE_IndirectSpecular )
	RE_IndirectSpecular( radiance, iblIrradiance, clearcoatRadiance, geometryPosition, geometryNormal, geometryViewDir, geometryClearcoatNormal, material, reflectedLight );
#endif`,Ax=`#ifdef USE_LIGHT_PROBES_GRID
uniform highp sampler3D probesSH;
uniform vec3 probesMin;
uniform vec3 probesMax;
uniform vec3 probesResolution;
vec3 getLightProbeGridIrradiance( vec3 worldPos, vec3 worldNormal ) {
	vec3 res = probesResolution;
	vec3 gridRange = probesMax - probesMin;
	vec3 resMinusOne = res - 1.0;
	vec3 probeSpacing = gridRange / resMinusOne;
	vec3 samplePos = worldPos + worldNormal * probeSpacing * 0.5;
	vec3 uvw = clamp( ( samplePos - probesMin ) / gridRange, 0.0, 1.0 );
	uvw = uvw * resMinusOne / res + 0.5 / res;
	float nz          = res.z;
	float paddedSlices = nz + 2.0;
	float atlasDepth  = 7.0 * paddedSlices;
	float uvZBase     = uvw.z * nz + 1.0;
	vec4 s0 = texture( probesSH, vec3( uvw.xy, ( uvZBase                       ) / atlasDepth ) );
	vec4 s1 = texture( probesSH, vec3( uvw.xy, ( uvZBase +       paddedSlices   ) / atlasDepth ) );
	vec4 s2 = texture( probesSH, vec3( uvw.xy, ( uvZBase + 2.0 * paddedSlices   ) / atlasDepth ) );
	vec4 s3 = texture( probesSH, vec3( uvw.xy, ( uvZBase + 3.0 * paddedSlices   ) / atlasDepth ) );
	vec4 s4 = texture( probesSH, vec3( uvw.xy, ( uvZBase + 4.0 * paddedSlices   ) / atlasDepth ) );
	vec4 s5 = texture( probesSH, vec3( uvw.xy, ( uvZBase + 5.0 * paddedSlices   ) / atlasDepth ) );
	vec4 s6 = texture( probesSH, vec3( uvw.xy, ( uvZBase + 6.0 * paddedSlices   ) / atlasDepth ) );
	vec3 c0 = s0.xyz;
	vec3 c1 = vec3( s0.w, s1.xy );
	vec3 c2 = vec3( s1.zw, s2.x );
	vec3 c3 = s2.yzw;
	vec3 c4 = s3.xyz;
	vec3 c5 = vec3( s3.w, s4.xy );
	vec3 c6 = vec3( s4.zw, s5.x );
	vec3 c7 = s5.yzw;
	vec3 c8 = s6.xyz;
	float x = worldNormal.x, y = worldNormal.y, z = worldNormal.z;
	vec3 result = c0 * 0.886227;
	result += c1 * 2.0 * 0.511664 * y;
	result += c2 * 2.0 * 0.511664 * z;
	result += c3 * 2.0 * 0.511664 * x;
	result += c4 * 2.0 * 0.429043 * x * y;
	result += c5 * 2.0 * 0.429043 * y * z;
	result += c6 * ( 0.743125 * z * z - 0.247708 );
	result += c7 * 2.0 * 0.429043 * x * z;
	result += c8 * 0.429043 * ( x * x - y * y );
	return max( result, vec3( 0.0 ) );
}
#endif`,wx=`#if defined( USE_LOGARITHMIC_DEPTH_BUFFER )
	gl_FragDepth = vIsPerspective == 0.0 ? gl_FragCoord.z : log2( vFragDepth ) * logDepthBufFC * 0.5;
#endif`,Rx=`#if defined( USE_LOGARITHMIC_DEPTH_BUFFER )
	uniform float logDepthBufFC;
	varying float vFragDepth;
	varying float vIsPerspective;
#endif`,Cx=`#ifdef USE_LOGARITHMIC_DEPTH_BUFFER
	varying float vFragDepth;
	varying float vIsPerspective;
#endif`,Ix=`#ifdef USE_LOGARITHMIC_DEPTH_BUFFER
	vFragDepth = 1.0 + gl_Position.w;
	vIsPerspective = float( isPerspectiveMatrix( projectionMatrix ) );
#endif`,Px=`#ifdef USE_MAP
	vec4 sampledDiffuseColor = texture2D( map, vMapUv );
	#ifdef DECODE_VIDEO_TEXTURE
		sampledDiffuseColor = sRGBTransferEOTF( sampledDiffuseColor );
	#endif
	diffuseColor *= sampledDiffuseColor;
#endif`,Lx=`#ifdef USE_MAP
	uniform sampler2D map;
#endif`,Nx=`#if defined( USE_MAP ) || defined( USE_ALPHAMAP )
	#if defined( USE_POINTS_UV )
		vec2 uv = vUv;
	#else
		vec2 uv = ( uvTransform * vec3( gl_PointCoord.x, 1.0 - gl_PointCoord.y, 1 ) ).xy;
	#endif
#endif
#ifdef USE_MAP
	diffuseColor *= texture2D( map, uv );
#endif
#ifdef USE_ALPHAMAP
	diffuseColor.a *= texture2D( alphaMap, uv ).g;
#endif`,Dx=`#if defined( USE_POINTS_UV )
	varying vec2 vUv;
#else
	#if defined( USE_MAP ) || defined( USE_ALPHAMAP )
		uniform mat3 uvTransform;
	#endif
#endif
#ifdef USE_MAP
	uniform sampler2D map;
#endif
#ifdef USE_ALPHAMAP
	uniform sampler2D alphaMap;
#endif`,Ux=`float metalnessFactor = metalness;
#ifdef USE_METALNESSMAP
	vec4 texelMetalness = texture2D( metalnessMap, vMetalnessMapUv );
	metalnessFactor *= texelMetalness.b;
#endif`,Fx=`#ifdef USE_METALNESSMAP
	uniform sampler2D metalnessMap;
#endif`,Ox=`#ifdef USE_INSTANCING_MORPH
	float morphTargetInfluences[ MORPHTARGETS_COUNT ];
	float morphTargetBaseInfluence = texelFetch( morphTexture, ivec2( 0, gl_InstanceID ), 0 ).r;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		morphTargetInfluences[i] =  texelFetch( morphTexture, ivec2( i + 1, gl_InstanceID ), 0 ).r;
	}
#endif`,Bx=`#if defined( USE_MORPHCOLORS )
	vColor *= morphTargetBaseInfluence;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		#if defined( USE_COLOR_ALPHA )
			if ( morphTargetInfluences[ i ] != 0.0 ) vColor += getMorph( gl_VertexID, i, 2 ) * morphTargetInfluences[ i ];
		#elif defined( USE_COLOR )
			if ( morphTargetInfluences[ i ] != 0.0 ) vColor += getMorph( gl_VertexID, i, 2 ).rgb * morphTargetInfluences[ i ];
		#endif
	}
#endif`,zx=`#ifdef USE_MORPHNORMALS
	objectNormal *= morphTargetBaseInfluence;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		if ( morphTargetInfluences[ i ] != 0.0 ) objectNormal += getMorph( gl_VertexID, i, 1 ).xyz * morphTargetInfluences[ i ];
	}
#endif`,kx=`#ifdef USE_MORPHTARGETS
	#ifndef USE_INSTANCING_MORPH
		uniform float morphTargetBaseInfluence;
		uniform float morphTargetInfluences[ MORPHTARGETS_COUNT ];
	#endif
	uniform sampler2DArray morphTargetsTexture;
	uniform ivec2 morphTargetsTextureSize;
	vec4 getMorph( const in int vertexIndex, const in int morphTargetIndex, const in int offset ) {
		int texelIndex = vertexIndex * MORPHTARGETS_TEXTURE_STRIDE + offset;
		int y = texelIndex / morphTargetsTextureSize.x;
		int x = texelIndex - y * morphTargetsTextureSize.x;
		ivec3 morphUV = ivec3( x, y, morphTargetIndex );
		return texelFetch( morphTargetsTexture, morphUV, 0 );
	}
#endif`,Gx=`#ifdef USE_MORPHTARGETS
	transformed *= morphTargetBaseInfluence;
	for ( int i = 0; i < MORPHTARGETS_COUNT; i ++ ) {
		if ( morphTargetInfluences[ i ] != 0.0 ) transformed += getMorph( gl_VertexID, i, 0 ).xyz * morphTargetInfluences[ i ];
	}
#endif`,Hx=`float faceDirection = gl_FrontFacing ? 1.0 : - 1.0;
#ifdef FLAT_SHADED
	vec3 fdx = dFdx( vViewPosition );
	vec3 fdy = dFdy( vViewPosition );
	vec3 normal = normalize( cross( fdx, fdy ) );
#else
	vec3 normal = normalize( vNormal );
	#ifdef DOUBLE_SIDED
		normal *= faceDirection;
	#endif
#endif
#if defined( USE_NORMALMAP_TANGENTSPACE ) || defined( USE_CLEARCOAT_NORMALMAP ) || defined( USE_ANISOTROPY )
	#ifdef USE_TANGENT
		mat3 tbn = mat3( normalize( vTangent ), normalize( vBitangent ), normal );
	#else
		mat3 tbn = getTangentFrame( - vViewPosition, normal,
		#if defined( USE_NORMALMAP )
			vNormalMapUv
		#elif defined( USE_CLEARCOAT_NORMALMAP )
			vClearcoatNormalMapUv
		#else
			vUv
		#endif
		);
	#endif
	#ifdef DOUBLE_SIDED
		tbn[0] *= faceDirection;
		tbn[1] *= faceDirection;
	#endif
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	#ifdef USE_TANGENT
		mat3 tbn2 = mat3( normalize( vTangent ), normalize( vBitangent ), normal );
	#else
		mat3 tbn2 = getTangentFrame( - vViewPosition, normal, vClearcoatNormalMapUv );
	#endif
	#ifdef DOUBLE_SIDED
		tbn2[0] *= faceDirection;
		tbn2[1] *= faceDirection;
	#endif
#endif
vec3 nonPerturbedNormal = normal;`,Vx=`#ifdef USE_NORMALMAP_OBJECTSPACE
	normal = texture2D( normalMap, vNormalMapUv ).xyz * 2.0 - 1.0;
	#ifdef FLIP_SIDED
		normal = - normal;
	#endif
	#ifdef DOUBLE_SIDED
		normal = normal * faceDirection;
	#endif
	normal = normalize( normalMatrix * normal );
#elif defined( USE_NORMALMAP_TANGENTSPACE )
	vec3 mapN = texture2D( normalMap, vNormalMapUv ).xyz * 2.0 - 1.0;
	#if defined( USE_PACKED_NORMALMAP )
		mapN = vec3( mapN.xy, sqrt( saturate( 1.0 - dot( mapN.xy, mapN.xy ) ) ) );
	#endif
	mapN.xy *= normalScale;
	normal = normalize( tbn * mapN );
#elif defined( USE_BUMPMAP )
	normal = perturbNormalArb( - vViewPosition, normal, dHdxy_fwd(), faceDirection );
#endif`,Wx=`#ifndef FLAT_SHADED
	varying vec3 vNormal;
	#ifdef USE_TANGENT
		varying vec3 vTangent;
		varying vec3 vBitangent;
	#endif
#endif`,Xx=`#ifndef FLAT_SHADED
	varying vec3 vNormal;
	#ifdef USE_TANGENT
		varying vec3 vTangent;
		varying vec3 vBitangent;
	#endif
#endif`,qx=`#ifndef FLAT_SHADED
	vNormal = normalize( transformedNormal );
	#ifdef USE_TANGENT
		vTangent = normalize( transformedTangent );
		vBitangent = normalize( cross( vNormal, vTangent ) * tangent.w );
		#ifdef FLIP_SIDED
			vBitangent = - vBitangent;
		#endif
	#endif
#endif`,Yx=`#ifdef USE_NORMALMAP
	uniform sampler2D normalMap;
	uniform vec2 normalScale;
#endif
#ifdef USE_NORMALMAP_OBJECTSPACE
	uniform mat3 normalMatrix;
#endif
#if ! defined ( USE_TANGENT ) && ( defined ( USE_NORMALMAP_TANGENTSPACE ) || defined ( USE_CLEARCOAT_NORMALMAP ) || defined( USE_ANISOTROPY ) )
	mat3 getTangentFrame( vec3 eye_pos, vec3 surf_norm, vec2 uv ) {
		vec3 q0 = dFdx( eye_pos.xyz );
		vec3 q1 = dFdy( eye_pos.xyz );
		vec2 st0 = dFdx( uv.st );
		vec2 st1 = dFdy( uv.st );
		vec3 N = surf_norm;
		vec3 q1perp = cross( q1, N );
		vec3 q0perp = cross( N, q0 );
		vec3 T = q1perp * st0.x + q0perp * st1.x;
		vec3 B = q1perp * st0.y + q0perp * st1.y;
		float det = max( dot( T, T ), dot( B, B ) );
		float scale = ( det == 0.0 ) ? 0.0 : inversesqrt( det );
		return mat3( T * scale, B * scale, N );
	}
#endif`,Zx=`#ifdef USE_CLEARCOAT
	vec3 clearcoatNormal = nonPerturbedNormal;
#endif`,Jx=`#ifdef USE_CLEARCOAT_NORMALMAP
	vec3 clearcoatMapN = texture2D( clearcoatNormalMap, vClearcoatNormalMapUv ).xyz * 2.0 - 1.0;
	clearcoatMapN.xy *= clearcoatNormalScale;
	clearcoatNormal = normalize( tbn2 * clearcoatMapN );
#endif`,$x=`#ifdef USE_CLEARCOATMAP
	uniform sampler2D clearcoatMap;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	uniform sampler2D clearcoatNormalMap;
	uniform vec2 clearcoatNormalScale;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	uniform sampler2D clearcoatRoughnessMap;
#endif`,Kx=`#ifdef USE_IRIDESCENCEMAP
	uniform sampler2D iridescenceMap;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	uniform sampler2D iridescenceThicknessMap;
#endif`,jx=`#ifdef OPAQUE
diffuseColor.a = 1.0;
#endif
#ifdef USE_TRANSMISSION
diffuseColor.a *= material.transmissionAlpha;
#endif
gl_FragColor = vec4( outgoingLight, diffuseColor.a );`,Qx=`vec3 packNormalToRGB( const in vec3 normal ) {
	return normalize( normal ) * 0.5 + 0.5;
}
vec3 unpackRGBToNormal( const in vec3 rgb ) {
	return 2.0 * rgb.xyz - 1.0;
}
const float PackUpscale = 256. / 255.;const float UnpackDownscale = 255. / 256.;const float ShiftRight8 = 1. / 256.;
const float Inv255 = 1. / 255.;
const vec4 PackFactors = vec4( 1.0, 256.0, 256.0 * 256.0, 256.0 * 256.0 * 256.0 );
const vec2 UnpackFactors2 = vec2( UnpackDownscale, 1.0 / PackFactors.g );
const vec3 UnpackFactors3 = vec3( UnpackDownscale / PackFactors.rg, 1.0 / PackFactors.b );
const vec4 UnpackFactors4 = vec4( UnpackDownscale / PackFactors.rgb, 1.0 / PackFactors.a );
vec4 packDepthToRGBA( const in float v ) {
	if( v <= 0.0 )
		return vec4( 0., 0., 0., 0. );
	if( v >= 1.0 )
		return vec4( 1., 1., 1., 1. );
	float vuf;
	float af = modf( v * PackFactors.a, vuf );
	float bf = modf( vuf * ShiftRight8, vuf );
	float gf = modf( vuf * ShiftRight8, vuf );
	return vec4( vuf * Inv255, gf * PackUpscale, bf * PackUpscale, af );
}
vec3 packDepthToRGB( const in float v ) {
	if( v <= 0.0 )
		return vec3( 0., 0., 0. );
	if( v >= 1.0 )
		return vec3( 1., 1., 1. );
	float vuf;
	float bf = modf( v * PackFactors.b, vuf );
	float gf = modf( vuf * ShiftRight8, vuf );
	return vec3( vuf * Inv255, gf * PackUpscale, bf );
}
vec2 packDepthToRG( const in float v ) {
	if( v <= 0.0 )
		return vec2( 0., 0. );
	if( v >= 1.0 )
		return vec2( 1., 1. );
	float vuf;
	float gf = modf( v * 256., vuf );
	return vec2( vuf * Inv255, gf );
}
float unpackRGBAToDepth( const in vec4 v ) {
	return dot( v, UnpackFactors4 );
}
float unpackRGBToDepth( const in vec3 v ) {
	return dot( v, UnpackFactors3 );
}
float unpackRGToDepth( const in vec2 v ) {
	return v.r * UnpackFactors2.r + v.g * UnpackFactors2.g;
}
vec4 pack2HalfToRGBA( const in vec2 v ) {
	vec4 r = vec4( v.x, fract( v.x * 255.0 ), v.y, fract( v.y * 255.0 ) );
	return vec4( r.x - r.y / 255.0, r.y, r.z - r.w / 255.0, r.w );
}
vec2 unpackRGBATo2Half( const in vec4 v ) {
	return vec2( v.x + ( v.y / 255.0 ), v.z + ( v.w / 255.0 ) );
}
float viewZToOrthographicDepth( const in float viewZ, const in float near, const in float far ) {
	return ( viewZ + near ) / ( near - far );
}
float orthographicDepthToViewZ( const in float depth, const in float near, const in float far ) {
	#ifdef USE_REVERSED_DEPTH_BUFFER
	
		return depth * ( far - near ) - far;
	#else
		return depth * ( near - far ) - near;
	#endif
}
float viewZToPerspectiveDepth( const in float viewZ, const in float near, const in float far ) {
	return ( ( near + viewZ ) * far ) / ( ( far - near ) * viewZ );
}
float perspectiveDepthToViewZ( const in float depth, const in float near, const in float far ) {
	
	#ifdef USE_REVERSED_DEPTH_BUFFER
		return ( near * far ) / ( ( near - far ) * depth - near );
	#else
		return ( near * far ) / ( ( far - near ) * depth - far );
	#endif
}`,ty=`#ifdef PREMULTIPLIED_ALPHA
	gl_FragColor.rgb *= gl_FragColor.a;
#endif`,ey=`vec4 mvPosition = vec4( transformed, 1.0 );
#ifdef USE_BATCHING
	mvPosition = batchingMatrix * mvPosition;
#endif
#ifdef USE_INSTANCING
	mvPosition = instanceMatrix * mvPosition;
#endif
mvPosition = modelViewMatrix * mvPosition;
gl_Position = projectionMatrix * mvPosition;`,ny=`#ifdef DITHERING
	gl_FragColor.rgb = dithering( gl_FragColor.rgb );
#endif`,iy=`#ifdef DITHERING
	vec3 dithering( vec3 color ) {
		float grid_position = rand( gl_FragCoord.xy );
		vec3 dither_shift_RGB = vec3( 0.25 / 255.0, -0.25 / 255.0, 0.25 / 255.0 );
		dither_shift_RGB = mix( 2.0 * dither_shift_RGB, -2.0 * dither_shift_RGB, grid_position );
		return color + dither_shift_RGB;
	}
#endif`,sy=`float roughnessFactor = roughness;
#ifdef USE_ROUGHNESSMAP
	vec4 texelRoughness = texture2D( roughnessMap, vRoughnessMapUv );
	roughnessFactor *= texelRoughness.g;
#endif`,ry=`#ifdef USE_ROUGHNESSMAP
	uniform sampler2D roughnessMap;
#endif`,ay=`#if NUM_SPOT_LIGHT_COORDS > 0
	varying vec4 vSpotLightCoord[ NUM_SPOT_LIGHT_COORDS ];
#endif
#if NUM_SPOT_LIGHT_MAPS > 0
	uniform sampler2D spotLightMap[ NUM_SPOT_LIGHT_MAPS ];
#endif
#ifdef USE_SHADOWMAP
	#if NUM_SUN_LIGHT_SHADOWS > 0
		#define SUN_LIGHT_CASCADES 2
		#if defined( SHADOWMAP_TYPE_PCF )
			uniform sampler2DShadow sunShadowMap[ NUM_SUN_LIGHT_SHADOWS ];
		#else
			uniform sampler2D sunShadowMap[ NUM_SUN_LIGHT_SHADOWS ];
		#endif
		uniform mat4 sunShadowMatrix[ NUM_SUN_LIGHT_SHADOWS * SUN_LIGHT_CASCADES ];
		uniform vec4 sunShadowCascade[ NUM_SUN_LIGHT_SHADOWS * SUN_LIGHT_CASCADES ];
		varying vec4 vSunShadowWorldPosition;
		varying vec3 vSunShadowWorldNormal;
		struct SunLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform SunLightShadow sunLightShadows[ NUM_SUN_LIGHT_SHADOWS ];
	#endif
	#if NUM_DIR_LIGHT_SHADOWS > 0
		#if defined( SHADOWMAP_TYPE_PCF )
			uniform sampler2DShadow directionalShadowMap[ NUM_DIR_LIGHT_SHADOWS ];
		#else
			uniform sampler2D directionalShadowMap[ NUM_DIR_LIGHT_SHADOWS ];
		#endif
		varying vec4 vDirectionalShadowCoord[ NUM_DIR_LIGHT_SHADOWS ];
		struct DirectionalLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform DirectionalLightShadow directionalLightShadows[ NUM_DIR_LIGHT_SHADOWS ];
	#endif
	#if NUM_SPOT_LIGHT_SHADOWS > 0
		#if defined( SHADOWMAP_TYPE_PCF )
			uniform sampler2DShadow spotShadowMap[ NUM_SPOT_LIGHT_SHADOWS ];
		#else
			uniform sampler2D spotShadowMap[ NUM_SPOT_LIGHT_SHADOWS ];
		#endif
		struct SpotLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform SpotLightShadow spotLightShadows[ NUM_SPOT_LIGHT_SHADOWS ];
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
		#if defined( SHADOWMAP_TYPE_PCF )
			uniform samplerCubeShadow pointShadowMap[ NUM_POINT_LIGHT_SHADOWS ];
		#elif defined( SHADOWMAP_TYPE_BASIC )
			uniform samplerCube pointShadowMap[ NUM_POINT_LIGHT_SHADOWS ];
		#endif
		varying vec4 vPointShadowCoord[ NUM_POINT_LIGHT_SHADOWS ];
		struct PointLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
			float shadowCameraNear;
			float shadowCameraFar;
		};
		uniform PointLightShadow pointLightShadows[ NUM_POINT_LIGHT_SHADOWS ];
	#endif
	#if defined( SHADOWMAP_TYPE_PCF )
		float interleavedGradientNoise( vec2 position ) {
			return fract( 52.9829189 * fract( dot( position, vec2( 0.06711056, 0.00583715 ) ) ) );
		}
		vec2 vogelDiskSample( int sampleIndex, int samplesCount, float phi ) {
			const float goldenAngle = 2.399963229728653;
			float r = sqrt( ( float( sampleIndex ) + 0.5 ) / float( samplesCount ) );
			float theta = float( sampleIndex ) * goldenAngle + phi;
			return vec2( cos( theta ), sin( theta ) ) * r;
		}
	#endif
	#if defined( SHADOWMAP_TYPE_PCF )
		float getShadow( sampler2DShadow shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord ) {
			float shadow = 1.0;
			shadowCoord.xyz /= shadowCoord.w;
			shadowCoord.z += shadowBias;
			bool inFrustum = shadowCoord.x >= 0.0 && shadowCoord.x <= 1.0 && shadowCoord.y >= 0.0 && shadowCoord.y <= 1.0;
			bool frustumTest = inFrustum && shadowCoord.z <= 1.0;
			if ( frustumTest ) {
				vec2 texelSize = vec2( 1.0 ) / shadowMapSize;
				float radius = shadowRadius * texelSize.x;
				float phi = interleavedGradientNoise( gl_FragCoord.xy ) * PI2;
				shadow = (
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 0, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 1, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 2, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 3, 5, phi ) * radius, shadowCoord.z ) ) +
					texture( shadowMap, vec3( shadowCoord.xy + vogelDiskSample( 4, 5, phi ) * radius, shadowCoord.z ) )
				) * 0.2;
			}
			return mix( 1.0, shadow, shadowIntensity );
		}
	#elif defined( SHADOWMAP_TYPE_VSM )
		float getShadow( sampler2D shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord ) {
			float shadow = 1.0;
			shadowCoord.xyz /= shadowCoord.w;
			#ifdef USE_REVERSED_DEPTH_BUFFER
				shadowCoord.z -= shadowBias;
			#else
				shadowCoord.z += shadowBias;
			#endif
			bool inFrustum = shadowCoord.x >= 0.0 && shadowCoord.x <= 1.0 && shadowCoord.y >= 0.0 && shadowCoord.y <= 1.0;
			bool frustumTest = inFrustum && shadowCoord.z <= 1.0;
			if ( frustumTest ) {
				vec2 distribution = texture2D( shadowMap, shadowCoord.xy ).rg;
				float mean = distribution.x;
				float variance = distribution.y * distribution.y;
				#ifdef USE_REVERSED_DEPTH_BUFFER
					float hard_shadow = step( mean, shadowCoord.z );
				#else
					float hard_shadow = step( shadowCoord.z, mean );
				#endif
				
				if ( hard_shadow == 1.0 ) {
					shadow = 1.0;
				} else {
					variance = max( variance, 0.0000001 );
					float d = shadowCoord.z - mean;
					float p_max = variance / ( variance + d * d );
					p_max = clamp( ( p_max - 0.3 ) / 0.65, 0.0, 1.0 );
					shadow = max( hard_shadow, p_max );
				}
			}
			return mix( 1.0, shadow, shadowIntensity );
		}
	#else
		float getShadow( sampler2D shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord ) {
			float shadow = 1.0;
			shadowCoord.xyz /= shadowCoord.w;
			#ifdef USE_REVERSED_DEPTH_BUFFER
				shadowCoord.z -= shadowBias;
			#else
				shadowCoord.z += shadowBias;
			#endif
			bool inFrustum = shadowCoord.x >= 0.0 && shadowCoord.x <= 1.0 && shadowCoord.y >= 0.0 && shadowCoord.y <= 1.0;
			bool frustumTest = inFrustum && shadowCoord.z <= 1.0;
			if ( frustumTest ) {
				float depth = texture2D( shadowMap, shadowCoord.xy ).r;
				#ifdef USE_REVERSED_DEPTH_BUFFER
					shadow = step( depth, shadowCoord.z );
				#else
					shadow = step( shadowCoord.z, depth );
				#endif
			}
			return mix( 1.0, shadow, shadowIntensity );
		}
	#endif
	#if NUM_SUN_LIGHT_SHADOWS > 0
		float getSunShadow(
			#if defined( SHADOWMAP_TYPE_PCF )
				sampler2DShadow shadowMap,
			#else
				sampler2D shadowMap,
			#endif
			SunLightShadow sunLightShadow,
			int shadowIndex
		) {
			vec4 shadowWorldPosition = vec4( vSunShadowWorldPosition.xyz + vSunShadowWorldNormal * sunLightShadow.shadowNormalBias, 1.0 );
			float viewDepth = vSunShadowWorldPosition.w;
			int cascadeOffset = shadowIndex * SUN_LIGHT_CASCADES;
			float shadow = 1.0;
			for ( int i = SUN_LIGHT_CASCADES - 1; i >= 0; i -- ) {
				vec4 cascade = sunShadowCascade[ cascadeOffset + i ];
				if ( viewDepth >= cascade.x && viewDepth < cascade.y ) {
					float cascadeShadow = getShadow(
						shadowMap,
						sunLightShadow.shadowMapSize,
						sunLightShadow.shadowIntensity,
						sunLightShadow.shadowBias,
						sunLightShadow.shadowRadius,
						sunShadowMatrix[ cascadeOffset + i ] * shadowWorldPosition
					);
					shadow = mix( cascadeShadow, shadow, smoothstep( cascade.z, cascade.y, viewDepth ) );
				}
			}
			return shadow;
		}
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
	#if defined( SHADOWMAP_TYPE_PCF )
	float getPointShadow( samplerCubeShadow shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord, float shadowCameraNear, float shadowCameraFar ) {
		float shadow = 1.0;
		vec3 lightToPosition = shadowCoord.xyz;
		vec3 bd3D = normalize( lightToPosition );
		vec3 absVec = abs( lightToPosition );
		float viewSpaceZ = max( max( absVec.x, absVec.y ), absVec.z );
		if ( viewSpaceZ - shadowCameraFar <= 0.0 && viewSpaceZ - shadowCameraNear >= 0.0 ) {
			#ifdef USE_REVERSED_DEPTH_BUFFER
				float dp = ( shadowCameraNear * ( shadowCameraFar - viewSpaceZ ) ) / ( viewSpaceZ * ( shadowCameraFar - shadowCameraNear ) );
				dp -= shadowBias;
			#else
				float dp = ( shadowCameraFar * ( viewSpaceZ - shadowCameraNear ) ) / ( viewSpaceZ * ( shadowCameraFar - shadowCameraNear ) );
				dp += shadowBias;
			#endif
			float texelSize = shadowRadius / shadowMapSize.x;
			vec3 absDir = abs( bd3D );
			vec3 tangent = absDir.x > absDir.z ? vec3( 0.0, 1.0, 0.0 ) : vec3( 1.0, 0.0, 0.0 );
			tangent = normalize( cross( bd3D, tangent ) );
			vec3 bitangent = cross( bd3D, tangent );
			float phi = interleavedGradientNoise( gl_FragCoord.xy ) * PI2;
			vec2 sample0 = vogelDiskSample( 0, 5, phi );
			vec2 sample1 = vogelDiskSample( 1, 5, phi );
			vec2 sample2 = vogelDiskSample( 2, 5, phi );
			vec2 sample3 = vogelDiskSample( 3, 5, phi );
			vec2 sample4 = vogelDiskSample( 4, 5, phi );
			shadow = (
				texture( shadowMap, vec4( bd3D + ( tangent * sample0.x + bitangent * sample0.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample1.x + bitangent * sample1.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample2.x + bitangent * sample2.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample3.x + bitangent * sample3.y ) * texelSize, dp ) ) +
				texture( shadowMap, vec4( bd3D + ( tangent * sample4.x + bitangent * sample4.y ) * texelSize, dp ) )
			) * 0.2;
		}
		return mix( 1.0, shadow, shadowIntensity );
	}
	#elif defined( SHADOWMAP_TYPE_BASIC )
	float getPointShadow( samplerCube shadowMap, vec2 shadowMapSize, float shadowIntensity, float shadowBias, float shadowRadius, vec4 shadowCoord, float shadowCameraNear, float shadowCameraFar ) {
		float shadow = 1.0;
		vec3 lightToPosition = shadowCoord.xyz;
		vec3 absVec = abs( lightToPosition );
		float viewSpaceZ = max( max( absVec.x, absVec.y ), absVec.z );
		if ( viewSpaceZ - shadowCameraFar <= 0.0 && viewSpaceZ - shadowCameraNear >= 0.0 ) {
			float dp = ( shadowCameraFar * ( viewSpaceZ - shadowCameraNear ) ) / ( viewSpaceZ * ( shadowCameraFar - shadowCameraNear ) );
			dp += shadowBias;
			vec3 bd3D = normalize( lightToPosition );
			float depth = textureCube( shadowMap, bd3D ).r;
			#ifdef USE_REVERSED_DEPTH_BUFFER
				depth = 1.0 - depth;
			#endif
			shadow = step( dp, depth );
		}
		return mix( 1.0, shadow, shadowIntensity );
	}
	#endif
	#endif
#endif`,oy=`#if NUM_SPOT_LIGHT_COORDS > 0
	uniform mat4 spotLightMatrix[ NUM_SPOT_LIGHT_COORDS ];
	varying vec4 vSpotLightCoord[ NUM_SPOT_LIGHT_COORDS ];
#endif
#ifdef USE_SHADOWMAP
	#if NUM_SUN_LIGHT_SHADOWS > 0
		varying vec4 vSunShadowWorldPosition;
		varying vec3 vSunShadowWorldNormal;
	#endif
	#if NUM_DIR_LIGHT_SHADOWS > 0
		uniform mat4 directionalShadowMatrix[ NUM_DIR_LIGHT_SHADOWS ];
		varying vec4 vDirectionalShadowCoord[ NUM_DIR_LIGHT_SHADOWS ];
		struct DirectionalLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform DirectionalLightShadow directionalLightShadows[ NUM_DIR_LIGHT_SHADOWS ];
	#endif
	#if NUM_SPOT_LIGHT_SHADOWS > 0
		struct SpotLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
		};
		uniform SpotLightShadow spotLightShadows[ NUM_SPOT_LIGHT_SHADOWS ];
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
		uniform mat4 pointShadowMatrix[ NUM_POINT_LIGHT_SHADOWS ];
		varying vec4 vPointShadowCoord[ NUM_POINT_LIGHT_SHADOWS ];
		struct PointLightShadow {
			float shadowIntensity;
			float shadowBias;
			float shadowNormalBias;
			float shadowRadius;
			vec2 shadowMapSize;
			float shadowCameraNear;
			float shadowCameraFar;
		};
		uniform PointLightShadow pointLightShadows[ NUM_POINT_LIGHT_SHADOWS ];
	#endif
#endif`,ly=`#if ( defined( USE_SHADOWMAP ) && ( NUM_DIR_LIGHT_SHADOWS > 0 || NUM_SUN_LIGHT_SHADOWS > 0 || NUM_POINT_LIGHT_SHADOWS > 0 ) ) || ( NUM_SPOT_LIGHT_COORDS > 0 )
	#ifdef HAS_NORMAL
		vec3 shadowWorldNormal = transformNormalByInverseViewMatrix( transformedNormal, viewMatrix );
	#else
		vec3 shadowWorldNormal = vec3( 0.0 );
	#endif
	vec4 shadowWorldPosition;
#endif
#if defined( USE_SHADOWMAP )
	#if NUM_SUN_LIGHT_SHADOWS > 0
		vSunShadowWorldPosition = vec4( worldPosition.xyz, - mvPosition.z );
		vSunShadowWorldNormal = shadowWorldNormal;
	#endif
	#if NUM_DIR_LIGHT_SHADOWS > 0
		#pragma unroll_loop_start
		for ( int i = 0; i < NUM_DIR_LIGHT_SHADOWS; i ++ ) {
			shadowWorldPosition = worldPosition + vec4( shadowWorldNormal * directionalLightShadows[ i ].shadowNormalBias, 0 );
			vDirectionalShadowCoord[ i ] = directionalShadowMatrix[ i ] * shadowWorldPosition;
		}
		#pragma unroll_loop_end
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0
		#pragma unroll_loop_start
		for ( int i = 0; i < NUM_POINT_LIGHT_SHADOWS; i ++ ) {
			shadowWorldPosition = worldPosition + vec4( shadowWorldNormal * pointLightShadows[ i ].shadowNormalBias, 0 );
			vPointShadowCoord[ i ] = pointShadowMatrix[ i ] * shadowWorldPosition;
		}
		#pragma unroll_loop_end
	#endif
#endif
#if NUM_SPOT_LIGHT_COORDS > 0
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SPOT_LIGHT_COORDS; i ++ ) {
		shadowWorldPosition = worldPosition;
		#if ( defined( USE_SHADOWMAP ) && UNROLLED_LOOP_INDEX < NUM_SPOT_LIGHT_SHADOWS )
			shadowWorldPosition.xyz += shadowWorldNormal * spotLightShadows[ i ].shadowNormalBias;
		#endif
		vSpotLightCoord[ i ] = spotLightMatrix[ i ] * shadowWorldPosition;
	}
	#pragma unroll_loop_end
#endif`,cy=`float getShadowMask() {
	float shadow = 1.0;
	#ifdef USE_SHADOWMAP
	#if NUM_SUN_LIGHT_SHADOWS > 0
	SunLightShadow sunLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SUN_LIGHT_SHADOWS; i ++ ) {
		sunLight = sunLightShadows[ i ];
		shadow *= receiveShadow ? getSunShadow( sunShadowMap[ i ], sunLight, UNROLLED_LOOP_INDEX ) : 1.0;
	}
	#pragma unroll_loop_end
	#endif
	#if NUM_DIR_LIGHT_SHADOWS > 0
	DirectionalLightShadow directionalLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_DIR_LIGHT_SHADOWS; i ++ ) {
		directionalLight = directionalLightShadows[ i ];
		shadow *= receiveShadow ? getShadow( directionalShadowMap[ i ], directionalLight.shadowMapSize, directionalLight.shadowIntensity, directionalLight.shadowBias, directionalLight.shadowRadius, vDirectionalShadowCoord[ i ] ) : 1.0;
	}
	#pragma unroll_loop_end
	#endif
	#if NUM_SPOT_LIGHT_SHADOWS > 0
	SpotLightShadow spotLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_SPOT_LIGHT_SHADOWS; i ++ ) {
		spotLight = spotLightShadows[ i ];
		shadow *= receiveShadow ? getShadow( spotShadowMap[ i ], spotLight.shadowMapSize, spotLight.shadowIntensity, spotLight.shadowBias, spotLight.shadowRadius, vSpotLightCoord[ i ] ) : 1.0;
	}
	#pragma unroll_loop_end
	#endif
	#if NUM_POINT_LIGHT_SHADOWS > 0 && ( defined( SHADOWMAP_TYPE_PCF ) || defined( SHADOWMAP_TYPE_BASIC ) )
	PointLightShadow pointLight;
	#pragma unroll_loop_start
	for ( int i = 0; i < NUM_POINT_LIGHT_SHADOWS; i ++ ) {
		pointLight = pointLightShadows[ i ];
		shadow *= receiveShadow ? getPointShadow( pointShadowMap[ i ], pointLight.shadowMapSize, pointLight.shadowIntensity, pointLight.shadowBias, pointLight.shadowRadius, vPointShadowCoord[ i ], pointLight.shadowCameraNear, pointLight.shadowCameraFar ) : 1.0;
	}
	#pragma unroll_loop_end
	#endif
	#endif
	return shadow;
}`,hy=`#ifdef USE_SKINNING
	mat4 boneMatX = getBoneMatrix( skinIndex.x );
	mat4 boneMatY = getBoneMatrix( skinIndex.y );
	mat4 boneMatZ = getBoneMatrix( skinIndex.z );
	mat4 boneMatW = getBoneMatrix( skinIndex.w );
#endif`,uy=`#ifdef USE_SKINNING
	uniform mat4 bindMatrix;
	uniform mat4 bindMatrixInverse;
	uniform highp sampler2D boneTexture;
	mat4 getBoneMatrix( const in float i ) {
		int size = textureSize( boneTexture, 0 ).x;
		int j = int( i ) * 4;
		int x = j % size;
		int y = j / size;
		vec4 v1 = texelFetch( boneTexture, ivec2( x, y ), 0 );
		vec4 v2 = texelFetch( boneTexture, ivec2( x + 1, y ), 0 );
		vec4 v3 = texelFetch( boneTexture, ivec2( x + 2, y ), 0 );
		vec4 v4 = texelFetch( boneTexture, ivec2( x + 3, y ), 0 );
		return mat4( v1, v2, v3, v4 );
	}
#endif`,dy=`#ifdef USE_SKINNING
	vec4 skinVertex = bindMatrix * vec4( transformed, 1.0 );
	vec4 skinned = vec4( 0.0 );
	skinned += boneMatX * skinVertex * skinWeight.x;
	skinned += boneMatY * skinVertex * skinWeight.y;
	skinned += boneMatZ * skinVertex * skinWeight.z;
	skinned += boneMatW * skinVertex * skinWeight.w;
	transformed = ( bindMatrixInverse * skinned ).xyz;
#endif`,fy=`#ifdef USE_SKINNING
	mat4 skinMatrix = mat4( 0.0 );
	skinMatrix += skinWeight.x * boneMatX;
	skinMatrix += skinWeight.y * boneMatY;
	skinMatrix += skinWeight.z * boneMatZ;
	skinMatrix += skinWeight.w * boneMatW;
	skinMatrix = bindMatrixInverse * skinMatrix * bindMatrix;
	objectNormal = vec4( skinMatrix * vec4( objectNormal, 0.0 ) ).xyz;
	#ifdef USE_TANGENT
		objectTangent = vec4( skinMatrix * vec4( objectTangent, 0.0 ) ).xyz;
	#endif
#endif`,py=`float specularStrength;
#ifdef USE_SPECULARMAP
	vec4 texelSpecular = texture2D( specularMap, vSpecularMapUv );
	specularStrength = texelSpecular.r;
#else
	specularStrength = 1.0;
#endif`,my=`#ifdef USE_SPECULARMAP
	uniform sampler2D specularMap;
#endif`,gy=`#if defined( TONE_MAPPING )
	gl_FragColor.rgb = toneMapping( gl_FragColor.rgb );
#endif`,_y=`#ifndef saturate
#define saturate( a ) clamp( a, 0.0, 1.0 )
#endif
uniform float toneMappingExposure;
vec3 LinearToneMapping( vec3 color ) {
	return saturate( toneMappingExposure * color );
}
vec3 ReinhardToneMapping( vec3 color ) {
	color *= toneMappingExposure;
	return saturate( color / ( vec3( 1.0 ) + color ) );
}
vec3 CineonToneMapping( vec3 color ) {
	color *= toneMappingExposure;
	color = max( vec3( 0.0 ), color - 0.004 );
	return pow( ( color * ( 6.2 * color + 0.5 ) ) / ( color * ( 6.2 * color + 1.7 ) + 0.06 ), vec3( 2.2 ) );
}
vec3 RRTAndODTFit( vec3 v ) {
	vec3 a = v * ( v + 0.0245786 ) - 0.000090537;
	vec3 b = v * ( 0.983729 * v + 0.4329510 ) + 0.238081;
	return a / b;
}
vec3 ACESFilmicToneMapping( vec3 color ) {
	const mat3 ACESInputMat = mat3(
		vec3( 0.59719, 0.07600, 0.02840 ),		vec3( 0.35458, 0.90834, 0.13383 ),
		vec3( 0.04823, 0.01566, 0.83777 )
	);
	const mat3 ACESOutputMat = mat3(
		vec3(  1.60475, -0.10208, -0.00327 ),		vec3( -0.53108,  1.10813, -0.07276 ),
		vec3( -0.07367, -0.00605,  1.07602 )
	);
	color *= toneMappingExposure / 0.6;
	color = ACESInputMat * color;
	color = RRTAndODTFit( color );
	color = ACESOutputMat * color;
	return saturate( color );
}
const mat3 LINEAR_REC2020_TO_LINEAR_SRGB = mat3(
	vec3( 1.6605, - 0.1246, - 0.0182 ),
	vec3( - 0.5876, 1.1329, - 0.1006 ),
	vec3( - 0.0728, - 0.0083, 1.1187 )
);
const mat3 LINEAR_SRGB_TO_LINEAR_REC2020 = mat3(
	vec3( 0.6274, 0.0691, 0.0164 ),
	vec3( 0.3293, 0.9195, 0.0880 ),
	vec3( 0.0433, 0.0113, 0.8956 )
);
vec3 agxDefaultContrastApprox( vec3 x ) {
	vec3 x2 = x * x;
	vec3 x4 = x2 * x2;
	return + 15.5 * x4 * x2
		- 40.14 * x4 * x
		+ 31.96 * x4
		- 6.868 * x2 * x
		+ 0.4298 * x2
		+ 0.1191 * x
		- 0.00232;
}
vec3 AgXToneMapping( vec3 color ) {
	const mat3 AgXInsetMatrix = mat3(
		vec3( 0.856627153315983, 0.137318972929847, 0.11189821299995 ),
		vec3( 0.0951212405381588, 0.761241990602591, 0.0767994186031903 ),
		vec3( 0.0482516061458583, 0.101439036467562, 0.811302368396859 )
	);
	const mat3 AgXOutsetMatrix = mat3(
		vec3( 1.1271005818144368, - 0.1413297634984383, - 0.14132976349843826 ),
		vec3( - 0.11060664309660323, 1.157823702216272, - 0.11060664309660294 ),
		vec3( - 0.016493938717834573, - 0.016493938717834257, 1.2519364065950405 )
	);
	const float AgxMinEv = - 12.47393;	const float AgxMaxEv = 4.026069;
	color *= toneMappingExposure;
	color = LINEAR_SRGB_TO_LINEAR_REC2020 * color;
	color = AgXInsetMatrix * color;
	color = max( color, 1e-10 );	color = log2( color );
	color = ( color - AgxMinEv ) / ( AgxMaxEv - AgxMinEv );
	color = clamp( color, 0.0, 1.0 );
	color = agxDefaultContrastApprox( color );
	color = AgXOutsetMatrix * color;
	color = pow( max( vec3( 0.0 ), color ), vec3( 2.2 ) );
	color = LINEAR_REC2020_TO_LINEAR_SRGB * color;
	color = clamp( color, 0.0, 1.0 );
	return color;
}
vec3 NeutralToneMapping( vec3 color ) {
	const float StartCompression = 0.8 - 0.04;
	const float Desaturation = 0.15;
	color *= toneMappingExposure;
	float x = min( color.r, min( color.g, color.b ) );
	float offset = x < 0.08 ? x - 6.25 * x * x : 0.04;
	color -= offset;
	float peak = max( color.r, max( color.g, color.b ) );
	if ( peak < StartCompression ) return color;
	float d = 1. - StartCompression;
	float newPeak = 1. - d * d / ( peak + d - StartCompression );
	color *= newPeak / peak;
	float g = 1. - 1. / ( Desaturation * ( peak - newPeak ) + 1. );
	return mix( color, vec3( newPeak ), g );
}
vec3 CustomToneMapping( vec3 color ) { return color; }`,xy=`#ifdef USE_TRANSMISSION
	material.transmission = transmission;
	material.transmissionAlpha = 1.0;
	material.thickness = thickness;
	material.attenuationDistance = attenuationDistance;
	material.attenuationColor = attenuationColor;
	#ifdef USE_TRANSMISSIONMAP
		material.transmission *= texture2D( transmissionMap, vTransmissionMapUv ).r;
	#endif
	#ifdef USE_THICKNESSMAP
		material.thickness *= texture2D( thicknessMap, vThicknessMapUv ).g;
	#endif
	vec3 pos = vWorldPosition;
	vec3 v = normalize( cameraPosition - pos );
	vec3 n = transformNormalByInverseViewMatrix( normal, viewMatrix );
	vec4 transmitted = getIBLVolumeRefraction(
		n, v, material.roughness, material.diffuseContribution, material.specularColorBlended, material.specularF90,
		pos, modelMatrix, viewMatrix, projectionMatrix, material.dispersion, material.ior, material.thickness,
		material.attenuationColor, material.attenuationDistance );
	material.transmissionAlpha = mix( material.transmissionAlpha, transmitted.a, material.transmission );
	totalDiffuse = mix( totalDiffuse, transmitted.rgb, material.transmission );
#endif`,yy=`#ifdef USE_TRANSMISSION
	uniform float transmission;
	uniform float thickness;
	uniform float attenuationDistance;
	uniform vec3 attenuationColor;
	#ifdef USE_TRANSMISSIONMAP
		uniform sampler2D transmissionMap;
	#endif
	#ifdef USE_THICKNESSMAP
		uniform sampler2D thicknessMap;
	#endif
	uniform vec2 transmissionSamplerSize;
	uniform sampler2D transmissionSamplerMap;
	uniform mat4 modelMatrix;
	uniform mat4 projectionMatrix;
	varying vec3 vWorldPosition;
	float w0( float a ) {
		return ( 1.0 / 6.0 ) * ( a * ( a * ( - a + 3.0 ) - 3.0 ) + 1.0 );
	}
	float w1( float a ) {
		return ( 1.0 / 6.0 ) * ( a *  a * ( 3.0 * a - 6.0 ) + 4.0 );
	}
	float w2( float a ){
		return ( 1.0 / 6.0 ) * ( a * ( a * ( - 3.0 * a + 3.0 ) + 3.0 ) + 1.0 );
	}
	float w3( float a ) {
		return ( 1.0 / 6.0 ) * ( a * a * a );
	}
	float g0( float a ) {
		return w0( a ) + w1( a );
	}
	float g1( float a ) {
		return w2( a ) + w3( a );
	}
	float h0( float a ) {
		return - 1.0 + w1( a ) / ( w0( a ) + w1( a ) );
	}
	float h1( float a ) {
		return 1.0 + w3( a ) / ( w2( a ) + w3( a ) );
	}
	vec4 bicubic( sampler2D tex, vec2 uv, vec4 texelSize, float lod ) {
		uv = uv * texelSize.zw + 0.5;
		vec2 iuv = floor( uv );
		vec2 fuv = fract( uv );
		float g0x = g0( fuv.x );
		float g1x = g1( fuv.x );
		float h0x = h0( fuv.x );
		float h1x = h1( fuv.x );
		float h0y = h0( fuv.y );
		float h1y = h1( fuv.y );
		vec2 p0 = ( vec2( iuv.x + h0x, iuv.y + h0y ) - 0.5 ) * texelSize.xy;
		vec2 p1 = ( vec2( iuv.x + h1x, iuv.y + h0y ) - 0.5 ) * texelSize.xy;
		vec2 p2 = ( vec2( iuv.x + h0x, iuv.y + h1y ) - 0.5 ) * texelSize.xy;
		vec2 p3 = ( vec2( iuv.x + h1x, iuv.y + h1y ) - 0.5 ) * texelSize.xy;
		return g0( fuv.y ) * ( g0x * textureLod( tex, p0, lod ) + g1x * textureLod( tex, p1, lod ) ) +
			g1( fuv.y ) * ( g0x * textureLod( tex, p2, lod ) + g1x * textureLod( tex, p3, lod ) );
	}
	vec4 textureBicubic( sampler2D sampler, vec2 uv, float lod ) {
		vec2 fLodSize = vec2( textureSize( sampler, int( lod ) ) );
		vec2 cLodSize = vec2( textureSize( sampler, int( lod + 1.0 ) ) );
		vec2 fLodSizeInv = 1.0 / fLodSize;
		vec2 cLodSizeInv = 1.0 / cLodSize;
		vec4 fSample = bicubic( sampler, uv, vec4( fLodSizeInv, fLodSize ), floor( lod ) );
		vec4 cSample = bicubic( sampler, uv, vec4( cLodSizeInv, cLodSize ), ceil( lod ) );
		return mix( fSample, cSample, fract( lod ) );
	}
	vec3 getVolumeTransmissionRay( const in vec3 n, const in vec3 v, const in float thickness, const in float ior, const in mat4 modelMatrix ) {
		vec3 refractionVector = refract( - v, normalize( n ), 1.0 / ior );
		vec3 modelScale;
		modelScale.x = length( vec3( modelMatrix[ 0 ].xyz ) );
		modelScale.y = length( vec3( modelMatrix[ 1 ].xyz ) );
		modelScale.z = length( vec3( modelMatrix[ 2 ].xyz ) );
		return normalize( refractionVector ) * thickness * modelScale;
	}
	float applyIorToRoughness( const in float roughness, const in float ior ) {
		return roughness * clamp( ior * 2.0 - 2.0, 0.0, 1.0 );
	}
	vec4 getTransmissionSample( const in vec2 fragCoord, const in float roughness, const in float ior ) {
		float lod = log2( transmissionSamplerSize.x ) * applyIorToRoughness( roughness, ior );
		return textureBicubic( transmissionSamplerMap, fragCoord.xy, lod );
	}
	vec3 volumeAttenuation( const in float transmissionDistance, const in vec3 attenuationColor, const in float attenuationDistance ) {
		if ( isinf( attenuationDistance ) ) {
			return vec3( 1.0 );
		} else {
			vec3 attenuationCoefficient = -log( attenuationColor ) / attenuationDistance;
			vec3 transmittance = exp( - attenuationCoefficient * transmissionDistance );			return transmittance;
		}
	}
	vec4 getIBLVolumeRefraction( const in vec3 n, const in vec3 v, const in float roughness, const in vec3 diffuseColor,
		const in vec3 specularColor, const in float specularF90, const in vec3 position, const in mat4 modelMatrix,
		const in mat4 viewMatrix, const in mat4 projMatrix, const in float dispersion, const in float ior, const in float thickness,
		const in vec3 attenuationColor, const in float attenuationDistance ) {
		vec4 transmittedLight;
		vec3 transmittance;
		#ifdef USE_DISPERSION
			float halfSpread = ( ior - 1.0 ) * 0.025 * dispersion;
			vec3 iors = vec3( ior - halfSpread, ior, ior + halfSpread );
			for ( int i = 0; i < 3; i ++ ) {
				vec3 transmissionRay = getVolumeTransmissionRay( n, v, thickness, iors[ i ], modelMatrix );
				vec3 refractedRayExit = position + transmissionRay;
				vec4 ndcPos = projMatrix * viewMatrix * vec4( refractedRayExit, 1.0 );
				vec2 refractionCoords = ndcPos.xy / ndcPos.w;
				refractionCoords += 1.0;
				refractionCoords /= 2.0;
				vec4 transmissionSample = getTransmissionSample( refractionCoords, roughness, iors[ i ] );
				transmittedLight[ i ] = transmissionSample[ i ];
				transmittedLight.a += transmissionSample.a;
				transmittance[ i ] = diffuseColor[ i ] * volumeAttenuation( length( transmissionRay ), attenuationColor, attenuationDistance )[ i ];
			}
			transmittedLight.a /= 3.0;
		#else
			vec3 transmissionRay = getVolumeTransmissionRay( n, v, thickness, ior, modelMatrix );
			vec3 refractedRayExit = position + transmissionRay;
			vec4 ndcPos = projMatrix * viewMatrix * vec4( refractedRayExit, 1.0 );
			vec2 refractionCoords = ndcPos.xy / ndcPos.w;
			refractionCoords += 1.0;
			refractionCoords /= 2.0;
			transmittedLight = getTransmissionSample( refractionCoords, roughness, ior );
			transmittance = diffuseColor * volumeAttenuation( length( transmissionRay ), attenuationColor, attenuationDistance );
		#endif
		vec3 attenuatedColor = transmittance * transmittedLight.rgb;
		vec3 F = EnvironmentBRDF( n, v, specularColor, specularF90, roughness );
		float transmittanceFactor = ( transmittance.r + transmittance.g + transmittance.b ) / 3.0;
		return vec4( ( 1.0 - F ) * attenuatedColor, 1.0 - ( 1.0 - transmittedLight.a ) * transmittanceFactor );
	}
#endif`,vy=`#if defined( USE_UV ) || defined( USE_ANISOTROPY )
	varying vec2 vUv;
#endif
#ifdef USE_MAP
	varying vec2 vMapUv;
#endif
#ifdef USE_ALPHAMAP
	varying vec2 vAlphaMapUv;
#endif
#ifdef USE_LIGHTMAP
	varying vec2 vLightMapUv;
#endif
#ifdef USE_AOMAP
	varying vec2 vAoMapUv;
#endif
#ifdef USE_BUMPMAP
	varying vec2 vBumpMapUv;
#endif
#ifdef USE_NORMALMAP
	varying vec2 vNormalMapUv;
#endif
#ifdef USE_EMISSIVEMAP
	varying vec2 vEmissiveMapUv;
#endif
#ifdef USE_METALNESSMAP
	varying vec2 vMetalnessMapUv;
#endif
#ifdef USE_ROUGHNESSMAP
	varying vec2 vRoughnessMapUv;
#endif
#ifdef USE_ANISOTROPYMAP
	varying vec2 vAnisotropyMapUv;
#endif
#ifdef USE_CLEARCOATMAP
	varying vec2 vClearcoatMapUv;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	varying vec2 vClearcoatNormalMapUv;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	varying vec2 vClearcoatRoughnessMapUv;
#endif
#ifdef USE_IRIDESCENCEMAP
	varying vec2 vIridescenceMapUv;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	varying vec2 vIridescenceThicknessMapUv;
#endif
#ifdef USE_SHEEN_COLORMAP
	varying vec2 vSheenColorMapUv;
#endif
#ifdef USE_SHEEN_ROUGHNESSMAP
	varying vec2 vSheenRoughnessMapUv;
#endif
#ifdef USE_SPECULARMAP
	varying vec2 vSpecularMapUv;
#endif
#ifdef USE_SPECULAR_COLORMAP
	varying vec2 vSpecularColorMapUv;
#endif
#ifdef USE_SPECULAR_INTENSITYMAP
	varying vec2 vSpecularIntensityMapUv;
#endif
#ifdef USE_TRANSMISSIONMAP
	uniform mat3 transmissionMapTransform;
	varying vec2 vTransmissionMapUv;
#endif
#ifdef USE_THICKNESSMAP
	uniform mat3 thicknessMapTransform;
	varying vec2 vThicknessMapUv;
#endif`,Sy=`#if defined( USE_UV ) || defined( USE_ANISOTROPY )
	varying vec2 vUv;
#endif
#ifdef USE_MAP
	uniform mat3 mapTransform;
	varying vec2 vMapUv;
#endif
#ifdef USE_ALPHAMAP
	uniform mat3 alphaMapTransform;
	varying vec2 vAlphaMapUv;
#endif
#ifdef USE_LIGHTMAP
	uniform mat3 lightMapTransform;
	varying vec2 vLightMapUv;
#endif
#ifdef USE_AOMAP
	uniform mat3 aoMapTransform;
	varying vec2 vAoMapUv;
#endif
#ifdef USE_BUMPMAP
	uniform mat3 bumpMapTransform;
	varying vec2 vBumpMapUv;
#endif
#ifdef USE_NORMALMAP
	uniform mat3 normalMapTransform;
	varying vec2 vNormalMapUv;
#endif
#ifdef USE_DISPLACEMENTMAP
	uniform mat3 displacementMapTransform;
	varying vec2 vDisplacementMapUv;
#endif
#ifdef USE_EMISSIVEMAP
	uniform mat3 emissiveMapTransform;
	varying vec2 vEmissiveMapUv;
#endif
#ifdef USE_METALNESSMAP
	uniform mat3 metalnessMapTransform;
	varying vec2 vMetalnessMapUv;
#endif
#ifdef USE_ROUGHNESSMAP
	uniform mat3 roughnessMapTransform;
	varying vec2 vRoughnessMapUv;
#endif
#ifdef USE_ANISOTROPYMAP
	uniform mat3 anisotropyMapTransform;
	varying vec2 vAnisotropyMapUv;
#endif
#ifdef USE_CLEARCOATMAP
	uniform mat3 clearcoatMapTransform;
	varying vec2 vClearcoatMapUv;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	uniform mat3 clearcoatNormalMapTransform;
	varying vec2 vClearcoatNormalMapUv;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	uniform mat3 clearcoatRoughnessMapTransform;
	varying vec2 vClearcoatRoughnessMapUv;
#endif
#ifdef USE_SHEEN_COLORMAP
	uniform mat3 sheenColorMapTransform;
	varying vec2 vSheenColorMapUv;
#endif
#ifdef USE_SHEEN_ROUGHNESSMAP
	uniform mat3 sheenRoughnessMapTransform;
	varying vec2 vSheenRoughnessMapUv;
#endif
#ifdef USE_IRIDESCENCEMAP
	uniform mat3 iridescenceMapTransform;
	varying vec2 vIridescenceMapUv;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	uniform mat3 iridescenceThicknessMapTransform;
	varying vec2 vIridescenceThicknessMapUv;
#endif
#ifdef USE_SPECULARMAP
	uniform mat3 specularMapTransform;
	varying vec2 vSpecularMapUv;
#endif
#ifdef USE_SPECULAR_COLORMAP
	uniform mat3 specularColorMapTransform;
	varying vec2 vSpecularColorMapUv;
#endif
#ifdef USE_SPECULAR_INTENSITYMAP
	uniform mat3 specularIntensityMapTransform;
	varying vec2 vSpecularIntensityMapUv;
#endif
#ifdef USE_TRANSMISSIONMAP
	uniform mat3 transmissionMapTransform;
	varying vec2 vTransmissionMapUv;
#endif
#ifdef USE_THICKNESSMAP
	uniform mat3 thicknessMapTransform;
	varying vec2 vThicknessMapUv;
#endif`,My=`#if defined( USE_UV ) || defined( USE_ANISOTROPY )
	vUv = vec3( uv, 1 ).xy;
#endif
#ifdef USE_MAP
	vMapUv = ( mapTransform * vec3( MAP_UV, 1 ) ).xy;
#endif
#ifdef USE_ALPHAMAP
	vAlphaMapUv = ( alphaMapTransform * vec3( ALPHAMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_LIGHTMAP
	vLightMapUv = ( lightMapTransform * vec3( LIGHTMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_AOMAP
	vAoMapUv = ( aoMapTransform * vec3( AOMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_BUMPMAP
	vBumpMapUv = ( bumpMapTransform * vec3( BUMPMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_NORMALMAP
	vNormalMapUv = ( normalMapTransform * vec3( NORMALMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_DISPLACEMENTMAP
	vDisplacementMapUv = ( displacementMapTransform * vec3( DISPLACEMENTMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_EMISSIVEMAP
	vEmissiveMapUv = ( emissiveMapTransform * vec3( EMISSIVEMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_METALNESSMAP
	vMetalnessMapUv = ( metalnessMapTransform * vec3( METALNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_ROUGHNESSMAP
	vRoughnessMapUv = ( roughnessMapTransform * vec3( ROUGHNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_ANISOTROPYMAP
	vAnisotropyMapUv = ( anisotropyMapTransform * vec3( ANISOTROPYMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_CLEARCOATMAP
	vClearcoatMapUv = ( clearcoatMapTransform * vec3( CLEARCOATMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_CLEARCOAT_NORMALMAP
	vClearcoatNormalMapUv = ( clearcoatNormalMapTransform * vec3( CLEARCOAT_NORMALMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_CLEARCOAT_ROUGHNESSMAP
	vClearcoatRoughnessMapUv = ( clearcoatRoughnessMapTransform * vec3( CLEARCOAT_ROUGHNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_IRIDESCENCEMAP
	vIridescenceMapUv = ( iridescenceMapTransform * vec3( IRIDESCENCEMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_IRIDESCENCE_THICKNESSMAP
	vIridescenceThicknessMapUv = ( iridescenceThicknessMapTransform * vec3( IRIDESCENCE_THICKNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SHEEN_COLORMAP
	vSheenColorMapUv = ( sheenColorMapTransform * vec3( SHEEN_COLORMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SHEEN_ROUGHNESSMAP
	vSheenRoughnessMapUv = ( sheenRoughnessMapTransform * vec3( SHEEN_ROUGHNESSMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SPECULARMAP
	vSpecularMapUv = ( specularMapTransform * vec3( SPECULARMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SPECULAR_COLORMAP
	vSpecularColorMapUv = ( specularColorMapTransform * vec3( SPECULAR_COLORMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_SPECULAR_INTENSITYMAP
	vSpecularIntensityMapUv = ( specularIntensityMapTransform * vec3( SPECULAR_INTENSITYMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_TRANSMISSIONMAP
	vTransmissionMapUv = ( transmissionMapTransform * vec3( TRANSMISSIONMAP_UV, 1 ) ).xy;
#endif
#ifdef USE_THICKNESSMAP
	vThicknessMapUv = ( thicknessMapTransform * vec3( THICKNESSMAP_UV, 1 ) ).xy;
#endif`;var by=`#if defined( USE_ENVMAP ) || defined( DISTANCE ) || defined ( USE_SHADOWMAP ) || defined ( USE_TRANSMISSION ) || NUM_SPOT_LIGHT_COORDS > 0
	vec4 worldPosition = vec4( transformed, 1.0 );
	#ifdef USE_BATCHING
		worldPosition = batchingMatrix * worldPosition;
	#endif
	#ifdef USE_INSTANCING
		worldPosition = instanceMatrix * worldPosition;
	#endif
	worldPosition = modelMatrix * worldPosition;
#endif`,Ty=`varying vec2 vUv;
uniform mat3 uvTransform;
void main() {
	vUv = ( uvTransform * vec3( uv, 1 ) ).xy;
	gl_Position = vec4( position.xy, 1.0, 1.0 );
}`,Ey=`uniform sampler2D t2D;
uniform float backgroundIntensity;
varying vec2 vUv;
void main() {
	vec4 texColor = texture2D( t2D, vUv );
	#ifdef DECODE_VIDEO_TEXTURE
		texColor = vec4( mix( pow( texColor.rgb * 0.9478672986 + vec3( 0.0521327014 ), vec3( 2.4 ) ), texColor.rgb * 0.0773993808, vec3( lessThanEqual( texColor.rgb, vec3( 0.04045 ) ) ) ), texColor.w );
	#endif
	texColor.rgb *= backgroundIntensity;
	gl_FragColor = texColor;
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,Ay=`varying vec3 vWorldDirection;
#include <common>
void main() {
	vWorldDirection = transformDirection( position, modelMatrix );
	#include <begin_vertex>
	#include <project_vertex>
	gl_Position.z = gl_Position.w;
}`,wy=`#ifdef ENVMAP_TYPE_CUBE
	uniform samplerCube envMap;
#elif defined( ENVMAP_TYPE_CUBE_UV )
	uniform sampler2D envMap;
#endif
uniform float backgroundBlurriness;
uniform float backgroundIntensity;
uniform mat3 backgroundRotation;
varying vec3 vWorldDirection;
#include <cube_uv_reflection_fragment>
void main() {
	#ifdef ENVMAP_TYPE_CUBE
		vec4 texColor = textureCube( envMap, backgroundRotation * vWorldDirection );
	#elif defined( ENVMAP_TYPE_CUBE_UV )
		vec4 texColor = textureCubeUV( envMap, backgroundRotation * vWorldDirection, backgroundBlurriness );
	#else
		vec4 texColor = vec4( 0.0, 0.0, 0.0, 1.0 );
	#endif
	texColor.rgb *= backgroundIntensity;
	gl_FragColor = texColor;
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,Ry=`varying vec3 vWorldDirection;
#include <common>
void main() {
	vWorldDirection = transformDirection( position, modelMatrix );
	#include <begin_vertex>
	#include <project_vertex>
	gl_Position.z = gl_Position.w;
}`,Cy=`uniform samplerCube tCube;
uniform float tFlip;
uniform float opacity;
varying vec3 vWorldDirection;
void main() {
	vec4 texColor = textureCube( tCube, vec3( tFlip * vWorldDirection.x, vWorldDirection.yz ) );
	gl_FragColor = texColor;
	gl_FragColor.a *= opacity;
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,Iy=`#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
varying vec2 vHighPrecisionZW;
void main() {
	#include <uv_vertex>
	#include <batching_vertex>
	#include <skinbase_vertex>
	#include <morphinstance_vertex>
	#ifdef USE_DISPLACEMENTMAP
		#include <beginnormal_vertex>
		#include <morphnormal_vertex>
		#include <skinnormal_vertex>
	#endif
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vHighPrecisionZW = gl_Position.zw;
}`,Py=`#if DEPTH_PACKING == 3200
	uniform float opacity;
#endif
#include <common>
#include <packing>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
varying vec2 vHighPrecisionZW;
void main() {
	vec4 diffuseColor = vec4( 1.0 );
	#include <clipping_planes_fragment>
	#if DEPTH_PACKING == 3200
		diffuseColor.a = opacity;
	#endif
	#include <map_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <logdepthbuf_fragment>
	#ifdef USE_REVERSED_DEPTH_BUFFER
		float fragCoordZ = vHighPrecisionZW[ 0 ] / vHighPrecisionZW[ 1 ];
	#else
		float fragCoordZ = 0.5 * vHighPrecisionZW[ 0 ] / vHighPrecisionZW[ 1 ] + 0.5;
	#endif
	#if DEPTH_PACKING == 3200
		gl_FragColor = vec4( vec3( 1.0 - fragCoordZ ), opacity );
	#elif DEPTH_PACKING == 3201
		gl_FragColor = packDepthToRGBA( fragCoordZ );
	#elif DEPTH_PACKING == 3202
		gl_FragColor = vec4( packDepthToRGB( fragCoordZ ), 1.0 );
	#elif DEPTH_PACKING == 3203
		gl_FragColor = vec4( packDepthToRG( fragCoordZ ), 0.0, 1.0 );
	#endif
}`,Ly=`#define DISTANCE
varying vec3 vWorldPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <batching_vertex>
	#include <skinbase_vertex>
	#include <morphinstance_vertex>
	#ifdef USE_DISPLACEMENTMAP
		#include <beginnormal_vertex>
		#include <morphnormal_vertex>
		#include <skinnormal_vertex>
	#endif
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <worldpos_vertex>
	#include <clipping_planes_vertex>
	vWorldPosition = worldPosition.xyz;
}`,Ny=`#define DISTANCE
uniform vec3 referencePosition;
uniform float nearDistance;
uniform float farDistance;
varying vec3 vWorldPosition;
#include <common>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( 1.0 );
	#include <clipping_planes_fragment>
	#include <map_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	float dist = length( vWorldPosition - referencePosition );
	dist = ( dist - nearDistance ) / ( farDistance - nearDistance );
	dist = saturate( dist );
	gl_FragColor = vec4( dist, 0.0, 0.0, 1.0 );
}`,Dy=`varying vec3 vWorldDirection;
#include <common>
void main() {
	vWorldDirection = transformDirection( position, modelMatrix );
	#include <begin_vertex>
	#include <project_vertex>
}`,Uy=`uniform sampler2D tEquirect;
varying vec3 vWorldDirection;
#include <common>
void main() {
	vec3 direction = normalize( vWorldDirection );
	vec2 sampleUV = equirectUv( direction );
	gl_FragColor = texture2D( tEquirect, sampleUV );
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
}`,Fy=`uniform float scale;
attribute float lineDistance;
varying float vLineDistance;
#include <common>
#include <uv_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	vLineDistance = scale * lineDistance;
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <fog_vertex>
}`,Oy=`uniform vec3 diffuse;
uniform float opacity;
uniform float dashSize;
uniform float totalSize;
varying float vLineDistance;
#include <common>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <fog_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	if ( mod( vLineDistance, totalSize ) > dashSize ) {
		discard;
	}
	vec3 outgoingLight = vec3( 0.0 );
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	outgoingLight = diffuseColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
}`,By=`#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <envmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#if defined ( USE_ENVMAP ) || defined ( USE_SKINNING )
		#include <beginnormal_vertex>
		#include <morphnormal_vertex>
		#include <skinbase_vertex>
		#include <skinnormal_vertex>
		#include <defaultnormal_vertex>
	#endif
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <worldpos_vertex>
	#include <envmap_vertex>
	#include <fog_vertex>
}`,zy=`uniform vec3 diffuse;
uniform float opacity;
#ifndef FLAT_SHADED
	varying vec3 vNormal;
#endif
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_pars_fragment>
#include <fog_pars_fragment>
#include <specularmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <specularmap_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	#ifdef USE_LIGHTMAP
		vec4 lightMapTexel = texture2D( lightMap, vLightMapUv );
		reflectedLight.indirectDiffuse += lightMapTexel.rgb * lightMapIntensity * RECIPROCAL_PI;
	#else
		reflectedLight.indirectDiffuse += vec3( 1.0 );
	#endif
	#include <aomap_fragment>
	reflectedLight.indirectDiffuse *= diffuseColor.rgb;
	vec3 outgoingLight = reflectedLight.indirectDiffuse;
	#include <envmap_fragment>
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,ky=`#define LAMBERT
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <envmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <envmap_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,Gy=`#define LAMBERT
uniform vec3 diffuse;
uniform vec3 emissive;
uniform float opacity;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <cube_uv_reflection_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_pars_fragment>
#include <envmap_physical_pars_fragment>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_lambert_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <specularmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <specularmap_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_lambert_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 outgoingLight = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse + totalEmissiveRadiance;
	#include <envmap_fragment>
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,Hy=`#define MATCAP
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <color_pars_vertex>
#include <displacementmap_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <fog_vertex>
	vViewPosition = - mvPosition.xyz;
}`,Vy=`#define MATCAP
uniform vec3 diffuse;
uniform float opacity;
uniform sampler2D matcap;
varying vec3 vViewPosition;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <fog_pars_fragment>
#include <normal_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	vec3 viewDir = normalize( vViewPosition );
	vec3 x = normalize( vec3( viewDir.z, 0.0, - viewDir.x ) );
	vec3 y = cross( viewDir, x );
	vec2 uv = vec2( dot( x, normal ), dot( y, normal ) ) * 0.495 + 0.5;
	#ifdef USE_MATCAP
		vec4 matcapColor = texture2D( matcap, uv );
	#else
		vec4 matcapColor = vec4( vec3( mix( 0.2, 0.8, uv.y ) ), 1.0 );
	#endif
	vec3 outgoingLight = diffuseColor.rgb * matcapColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,Wy=`#define NORMAL
#if defined( FLAT_SHADED ) || defined( USE_BUMPMAP ) || defined( USE_NORMALMAP_TANGENTSPACE )
	varying vec3 vViewPosition;
#endif
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphinstance_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
#if defined( FLAT_SHADED ) || defined( USE_BUMPMAP ) || defined( USE_NORMALMAP_TANGENTSPACE )
	vViewPosition = - mvPosition.xyz;
#endif
}`,Xy=`#define NORMAL
uniform float opacity;
#if defined( FLAT_SHADED ) || defined( USE_BUMPMAP ) || defined( USE_NORMALMAP_TANGENTSPACE )
	varying vec3 vViewPosition;
#endif
#include <uv_pars_fragment>
#include <normal_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( 0.0, 0.0, 0.0, opacity );
	#include <clipping_planes_fragment>
	#include <logdepthbuf_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	gl_FragColor = vec4( normalize( normal ) * 0.5 + 0.5, diffuseColor.a );
	#ifdef OPAQUE
		gl_FragColor.a = 1.0;
	#endif
}`,qy=`#define PHONG
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <envmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphinstance_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <envmap_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,Yy=`#define PHONG
uniform vec3 diffuse;
uniform vec3 emissive;
uniform vec3 specular;
uniform float shininess;
uniform float opacity;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <cube_uv_reflection_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_pars_fragment>
#include <envmap_physical_pars_fragment>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_phong_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <specularmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <specularmap_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_phong_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 outgoingLight = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse + reflectedLight.directSpecular + reflectedLight.indirectSpecular + totalEmissiveRadiance;
	#include <envmap_fragment>
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,Zy=`#define STANDARD
varying vec3 vViewPosition;
#ifdef USE_TRANSMISSION
	varying vec3 vWorldPosition;
#endif
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
#ifdef USE_TRANSMISSION
	vWorldPosition = worldPosition.xyz;
#endif
}`,Jy=`#define STANDARD
#ifdef PHYSICAL
	#define IOR
	#define USE_SPECULAR
#endif
uniform vec3 diffuse;
uniform vec3 emissive;
uniform float roughness;
uniform float metalness;
uniform float opacity;
#ifdef IOR
	uniform float ior;
#endif
#ifdef USE_SPECULAR
	uniform float specularIntensity;
	uniform vec3 specularColor;
	#ifdef USE_SPECULAR_COLORMAP
		uniform sampler2D specularColorMap;
	#endif
	#ifdef USE_SPECULAR_INTENSITYMAP
		uniform sampler2D specularIntensityMap;
	#endif
#endif
#ifdef USE_CLEARCOAT
	uniform float clearcoat;
	uniform float clearcoatRoughness;
#endif
#ifdef USE_DISPERSION
	uniform float dispersion;
#endif
#ifdef USE_RETROREFLECTION
	uniform float retroreflectivity;
#endif
#ifdef USE_IRIDESCENCE
	uniform float iridescence;
	uniform float iridescenceIOR;
	uniform float iridescenceThicknessMinimum;
	uniform float iridescenceThicknessMaximum;
#endif
#ifdef USE_SHEEN
	uniform vec3 sheenColor;
	uniform float sheenRoughness;
	#ifdef USE_SHEEN_COLORMAP
		uniform sampler2D sheenColorMap;
	#endif
	#ifdef USE_SHEEN_ROUGHNESSMAP
		uniform sampler2D sheenRoughnessMap;
	#endif
#endif
#ifdef USE_ANISOTROPY
	uniform vec2 anisotropyVector;
	#ifdef USE_ANISOTROPYMAP
		uniform sampler2D anisotropyMap;
	#endif
#endif
varying vec3 vViewPosition;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <iridescence_fragment>
#include <cube_uv_reflection_fragment>
#include <envmap_common_pars_fragment>
#include <envmap_physical_pars_fragment>
#include <fog_pars_fragment>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_physical_pars_fragment>
#include <transmission_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <clearcoat_pars_fragment>
#include <iridescence_pars_fragment>
#include <roughnessmap_pars_fragment>
#include <metalnessmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <roughnessmap_fragment>
	#include <metalnessmap_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <clearcoat_normal_fragment_begin>
	#include <clearcoat_normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_physical_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 totalDiffuse = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse;
	vec3 totalSpecular = reflectedLight.directSpecular + reflectedLight.indirectSpecular;
	#include <transmission_fragment>
	vec3 outgoingLight = totalDiffuse + totalSpecular + totalEmissiveRadiance;
	#ifdef USE_SHEEN
 
		outgoingLight = outgoingLight + sheenSpecularDirect + sheenSpecularIndirect;
 
 	#endif
	#ifdef USE_CLEARCOAT
		float dotNVcc = saturate( dot( geometryClearcoatNormal, geometryViewDir ) );
		vec3 Fcc = F_Schlick( material.clearcoatF0, material.clearcoatF90, dotNVcc );
		outgoingLight = outgoingLight * ( 1.0 - material.clearcoat * Fcc ) + ( clearcoatSpecularDirect + clearcoatSpecularIndirect ) * material.clearcoat;
	#endif
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,$y=`#define TOON
varying vec3 vViewPosition;
#include <common>
#include <batching_pars_vertex>
#include <uv_pars_vertex>
#include <displacementmap_pars_vertex>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <normal_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <shadowmap_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <normal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <displacementmap_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	vViewPosition = - mvPosition.xyz;
	#include <worldpos_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,Ky=`#define TOON
uniform vec3 diffuse;
uniform vec3 emissive;
uniform float opacity;
#include <common>
#include <dithering_pars_fragment>
#include <color_pars_fragment>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <aomap_pars_fragment>
#include <lightmap_pars_fragment>
#include <emissivemap_pars_fragment>
#include <gradientmap_pars_fragment>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <normal_pars_fragment>
#include <lights_toon_pars_fragment>
#include <shadowmap_pars_fragment>
#include <bumpmap_pars_fragment>
#include <normalmap_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	ReflectedLight reflectedLight = ReflectedLight( vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ), vec3( 0.0 ) );
	vec3 totalEmissiveRadiance = emissive;
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <color_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	#include <normal_fragment_begin>
	#include <normal_fragment_maps>
	#include <emissivemap_fragment>
	#include <lights_toon_fragment>
	#include <lights_fragment_begin>
	#include <lights_fragment_maps>
	#include <lights_fragment_end>
	#include <aomap_fragment>
	vec3 outgoingLight = reflectedLight.directDiffuse + reflectedLight.indirectDiffuse + totalEmissiveRadiance;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
	#include <dithering_fragment>
}`,jy=`uniform float size;
uniform float scale;
#include <common>
#include <color_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
#ifdef USE_POINTS_UV
	varying vec2 vUv;
	uniform mat3 uvTransform;
#endif
void main() {
	#ifdef USE_POINTS_UV
		vUv = ( uvTransform * vec3( uv, 1 ) ).xy;
	#endif
	#include <color_vertex>
	#include <morphinstance_vertex>
	#include <morphcolor_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <project_vertex>
	gl_PointSize = size;
	#ifdef USE_SIZEATTENUATION
		bool isPerspective = isPerspectiveMatrix( projectionMatrix );
		if ( isPerspective ) gl_PointSize *= ( scale / - mvPosition.z );
	#endif
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <worldpos_vertex>
	#include <fog_vertex>
}`,Qy=`uniform vec3 diffuse;
uniform float opacity;
#include <common>
#include <color_pars_fragment>
#include <map_particle_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <fog_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	vec3 outgoingLight = vec3( 0.0 );
	#include <logdepthbuf_fragment>
	#include <map_particle_fragment>
	#include <color_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	outgoingLight = diffuseColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
}`,tv=`#include <common>
#include <batching_pars_vertex>
#include <fog_pars_vertex>
#include <morphtarget_pars_vertex>
#include <skinning_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <shadowmap_pars_vertex>
void main() {
	#include <batching_vertex>
	#include <beginnormal_vertex>
	#include <morphinstance_vertex>
	#include <morphnormal_vertex>
	#include <skinbase_vertex>
	#include <skinnormal_vertex>
	#include <defaultnormal_vertex>
	#include <begin_vertex>
	#include <morphtarget_vertex>
	#include <skinning_vertex>
	#include <project_vertex>
	#include <logdepthbuf_vertex>
	#include <worldpos_vertex>
	#include <shadowmap_vertex>
	#include <fog_vertex>
}`,ev=`uniform vec3 color;
uniform float opacity;
#include <common>
#include <fog_pars_fragment>
#include <bsdfs>
#include <lights_pars_begin>
#include <logdepthbuf_pars_fragment>
#include <shadowmap_pars_fragment>
#include <shadowmask_pars_fragment>
void main() {
	#include <logdepthbuf_fragment>
	gl_FragColor = vec4( color, opacity * ( 1.0 - getShadowMask() ) );
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
	#include <premultiplied_alpha_fragment>
}`,nv=`uniform float rotation;
uniform vec2 center;
#include <common>
#include <uv_pars_vertex>
#include <fog_pars_vertex>
#include <logdepthbuf_pars_vertex>
#include <clipping_planes_pars_vertex>
void main() {
	#include <uv_vertex>
	vec4 mvPosition = modelViewMatrix[ 3 ];
	vec2 scale = vec2( length( modelMatrix[ 0 ].xyz ), length( modelMatrix[ 1 ].xyz ) );
	#ifndef USE_SIZEATTENUATION
		bool isPerspective = isPerspectiveMatrix( projectionMatrix );
		if ( isPerspective ) scale *= - mvPosition.z;
	#endif
	vec2 alignedPosition = ( position.xy - ( center - vec2( 0.5 ) ) ) * scale;
	vec2 rotatedPosition;
	rotatedPosition.x = cos( rotation ) * alignedPosition.x - sin( rotation ) * alignedPosition.y;
	rotatedPosition.y = sin( rotation ) * alignedPosition.x + cos( rotation ) * alignedPosition.y;
	mvPosition.xy += rotatedPosition;
	gl_Position = projectionMatrix * mvPosition;
	#include <logdepthbuf_vertex>
	#include <clipping_planes_vertex>
	#include <fog_vertex>
}`,iv=`uniform vec3 diffuse;
uniform float opacity;
#include <common>
#include <uv_pars_fragment>
#include <map_pars_fragment>
#include <alphamap_pars_fragment>
#include <alphatest_pars_fragment>
#include <alphahash_pars_fragment>
#include <fog_pars_fragment>
#include <logdepthbuf_pars_fragment>
#include <clipping_planes_pars_fragment>
void main() {
	vec4 diffuseColor = vec4( diffuse, opacity );
	#include <clipping_planes_fragment>
	vec3 outgoingLight = vec3( 0.0 );
	#include <logdepthbuf_fragment>
	#include <map_fragment>
	#include <alphamap_fragment>
	#include <alphatest_fragment>
	#include <alphahash_fragment>
	outgoingLight = diffuseColor.rgb;
	#include <opaque_fragment>
	#include <tonemapping_fragment>
	#include <colorspace_fragment>
	#include <fog_fragment>
}`,ie={alphahash_fragment:T_,alphahash_pars_fragment:E_,alphamap_fragment:A_,alphamap_pars_fragment:w_,alphatest_fragment:R_,alphatest_pars_fragment:C_,aomap_fragment:I_,aomap_pars_fragment:P_,batching_pars_vertex:L_,batching_vertex:N_,begin_vertex:D_,beginnormal_vertex:U_,bsdfs:F_,iridescence_fragment:O_,bumpmap_pars_fragment:B_,clipping_planes_fragment:z_,clipping_planes_pars_fragment:k_,clipping_planes_pars_vertex:G_,clipping_planes_vertex:H_,color_fragment:V_,color_pars_fragment:W_,color_pars_vertex:X_,color_vertex:q_,common:Y_,cube_uv_reflection_fragment:Z_,defaultnormal_vertex:J_,displacementmap_pars_vertex:$_,displacementmap_vertex:K_,emissivemap_fragment:j_,emissivemap_pars_fragment:Q_,colorspace_fragment:tx,colorspace_pars_fragment:ex,envmap_fragment:nx,envmap_common_pars_fragment:ix,envmap_pars_fragment:sx,envmap_pars_vertex:rx,envmap_physical_pars_fragment:gx,envmap_vertex:ax,fog_vertex:ox,fog_pars_vertex:lx,fog_fragment:cx,fog_pars_fragment:hx,gradientmap_pars_fragment:ux,lightmap_pars_fragment:dx,lights_lambert_fragment:fx,lights_lambert_pars_fragment:px,lights_pars_begin:mx,lights_toon_fragment:_x,lights_toon_pars_fragment:xx,lights_phong_fragment:yx,lights_phong_pars_fragment:vx,lights_physical_fragment:Sx,lights_physical_pars_fragment:Mx,lights_fragment_begin:bx,lights_fragment_maps:Tx,lights_fragment_end:Ex,lightprobes_pars_fragment:Ax,logdepthbuf_fragment:wx,logdepthbuf_pars_fragment:Rx,logdepthbuf_pars_vertex:Cx,logdepthbuf_vertex:Ix,map_fragment:Px,map_pars_fragment:Lx,map_particle_fragment:Nx,map_particle_pars_fragment:Dx,metalnessmap_fragment:Ux,metalnessmap_pars_fragment:Fx,morphinstance_vertex:Ox,morphcolor_vertex:Bx,morphnormal_vertex:zx,morphtarget_pars_vertex:kx,morphtarget_vertex:Gx,normal_fragment_begin:Hx,normal_fragment_maps:Vx,normal_pars_fragment:Wx,normal_pars_vertex:Xx,normal_vertex:qx,normalmap_pars_fragment:Yx,clearcoat_normal_fragment_begin:Zx,clearcoat_normal_fragment_maps:Jx,clearcoat_pars_fragment:$x,iridescence_pars_fragment:Kx,opaque_fragment:jx,packing:Qx,premultiplied_alpha_fragment:ty,project_vertex:ey,dithering_fragment:ny,dithering_pars_fragment:iy,roughnessmap_fragment:sy,roughnessmap_pars_fragment:ry,shadowmap_pars_fragment:ay,shadowmap_pars_vertex:oy,shadowmap_vertex:ly,shadowmask_pars_fragment:cy,skinbase_vertex:hy,skinning_pars_vertex:uy,skinning_vertex:dy,skinnormal_vertex:fy,specularmap_fragment:py,specularmap_pars_fragment:my,tonemapping_fragment:gy,tonemapping_pars_fragment:_y,transmission_fragment:xy,transmission_pars_fragment:yy,uv_pars_fragment:vy,uv_pars_vertex:Sy,uv_vertex:My,worldpos_vertex:by,background_vert:Ty,background_frag:Ey,backgroundCube_vert:Ay,backgroundCube_frag:wy,cube_vert:Ry,cube_frag:Cy,depth_vert:Iy,depth_frag:Py,distance_vert:Ly,distance_frag:Ny,equirect_vert:Dy,equirect_frag:Uy,linedashed_vert:Fy,linedashed_frag:Oy,meshbasic_vert:By,meshbasic_frag:zy,meshlambert_vert:ky,meshlambert_frag:Gy,meshmatcap_vert:Hy,meshmatcap_frag:Vy,meshnormal_vert:Wy,meshnormal_frag:Xy,meshphong_vert:qy,meshphong_frag:Yy,meshphysical_vert:Zy,meshphysical_frag:Jy,meshtoon_vert:$y,meshtoon_frag:Ky,points_vert:jy,points_frag:Qy,shadow_vert:tv,shadow_frag:ev,sprite_vert:nv,sprite_frag:iv},At={common:{diffuse:{value:new Pt(16777215)},opacity:{value:1},map:{value:null},mapTransform:{value:new Xt},alphaMap:{value:null},alphaMapTransform:{value:new Xt},alphaTest:{value:0}},specularmap:{specularMap:{value:null},specularMapTransform:{value:new Xt}},envmap:{envMap:{value:null},envMapRotation:{value:new Xt},reflectivity:{value:1},ior:{value:1.5},refractionRatio:{value:0.98},dfgLUT:{value:null}},aomap:{aoMap:{value:null},aoMapIntensity:{value:1},aoMapTransform:{value:new Xt}},lightmap:{lightMap:{value:null},lightMapIntensity:{value:1},lightMapTransform:{value:new Xt}},bumpmap:{bumpMap:{value:null},bumpMapTransform:{value:new Xt},bumpScale:{value:1}},normalmap:{normalMap:{value:null},normalMapTransform:{value:new Xt},normalScale:{value:new rt(1,1)}},displacementmap:{displacementMap:{value:null},displacementMapTransform:{value:new Xt},displacementScale:{value:1},displacementBias:{value:0}},emissivemap:{emissiveMap:{value:null},emissiveMapTransform:{value:new Xt}},metalnessmap:{metalnessMap:{value:null},metalnessMapTransform:{value:new Xt}},roughnessmap:{roughnessMap:{value:null},roughnessMapTransform:{value:new Xt}},gradientmap:{gradientMap:{value:null}},fog:{fogDensity:{value:0.00025},fogNear:{value:1},fogFar:{value:2000},fogColor:{value:new Pt(16777215)}},lights:{ambientLightColor:{value:[]},lightProbe:{value:[]},sunLights:{value:[],properties:{direction:{},color:{}}},sunLightShadows:{value:[],properties:{shadowIntensity:1,shadowBias:{},shadowNormalBias:{},shadowRadius:{},shadowMapSize:{}}},sunShadowMatrix:{value:[]},sunShadowCascade:{value:[]},directionalLights:{value:[],properties:{direction:{},color:{}}},directionalLightShadows:{value:[],properties:{shadowIntensity:1,shadowBias:{},shadowNormalBias:{},shadowRadius:{},shadowMapSize:{}}},directionalShadowMatrix:{value:[]},spotLights:{value:[],properties:{color:{},position:{},direction:{},distance:{},coneCos:{},penumbraCos:{},decay:{}}},spotLightShadows:{value:[],properties:{shadowIntensity:1,shadowBias:{},shadowNormalBias:{},shadowRadius:{},shadowMapSize:{}}},spotLightMap:{value:[]},spotLightMatrix:{value:[]},pointLights:{value:[],properties:{color:{},position:{},decay:{},distance:{}}},pointLightShadows:{value:[],properties:{shadowIntensity:1,shadowBias:{},shadowNormalBias:{},shadowRadius:{},shadowMapSize:{},shadowCameraNear:{},shadowCameraFar:{}}},pointShadowMatrix:{value:[]},hemisphereLights:{value:[],properties:{direction:{},skyColor:{},groundColor:{}}},rectAreaLights:{value:[],properties:{color:{},position:{},width:{},height:{}}},ltc_1:{value:null},ltc_2:{value:null},probesSH:{value:null},probesMin:{value:new L},probesMax:{value:new L},probesResolution:{value:new L}},points:{diffuse:{value:new Pt(16777215)},opacity:{value:1},size:{value:1},scale:{value:1},map:{value:null},alphaMap:{value:null},alphaMapTransform:{value:new Xt},alphaTest:{value:0},uvTransform:{value:new Xt}},sprite:{diffuse:{value:new Pt(16777215)},opacity:{value:1},center:{value:new rt(0.5,0.5)},rotation:{value:0},map:{value:null},mapTransform:{value:new Xt},alphaMap:{value:null},alphaMapTransform:{value:new Xt},alphaTest:{value:0}}},On={basic:{uniforms:Xe([At.common,At.specularmap,At.envmap,At.aomap,At.lightmap,At.fog]),vertexShader:ie.meshbasic_vert,fragmentShader:ie.meshbasic_frag},lambert:{uniforms:Xe([At.common,At.specularmap,At.envmap,At.aomap,At.lightmap,At.emissivemap,At.bumpmap,At.normalmap,At.displacementmap,At.fog,At.lights,{emissive:{value:new Pt(0)},envMapIntensity:{value:1}}]),vertexShader:ie.meshlambert_vert,fragmentShader:ie.meshlambert_frag},phong:{uniforms:Xe([At.common,At.specularmap,At.envmap,At.aomap,At.lightmap,At.emissivemap,At.bumpmap,At.normalmap,At.displacementmap,At.fog,At.lights,{emissive:{value:new Pt(0)},specular:{value:new Pt(1118481)},shininess:{value:30},envMapIntensity:{value:1}}]),vertexShader:ie.meshphong_vert,fragmentShader:ie.meshphong_frag},standard:{uniforms:Xe([At.common,At.envmap,At.aomap,At.lightmap,At.emissivemap,At.bumpmap,At.normalmap,At.displacementmap,At.roughnessmap,At.metalnessmap,At.fog,At.lights,{emissive:{value:new Pt(0)},roughness:{value:1},metalness:{value:0},envMapIntensity:{value:1}}]),vertexShader:ie.meshphysical_vert,fragmentShader:ie.meshphysical_frag},toon:{uniforms:Xe([At.common,At.aomap,At.lightmap,At.emissivemap,At.bumpmap,At.normalmap,At.displacementmap,At.gradientmap,At.fog,At.lights,{emissive:{value:new Pt(0)}}]),vertexShader:ie.meshtoon_vert,fragmentShader:ie.meshtoon_frag},matcap:{uniforms:Xe([At.common,At.bumpmap,At.normalmap,At.displacementmap,At.fog,{matcap:{value:null}}]),vertexShader:ie.meshmatcap_vert,fragmentShader:ie.meshmatcap_frag},points:{uniforms:Xe([At.points,At.fog]),vertexShader:ie.points_vert,fragmentShader:ie.points_frag},dashed:{uniforms:Xe([At.common,At.fog,{scale:{value:1},dashSize:{value:1},totalSize:{value:2}}]),vertexShader:ie.linedashed_vert,fragmentShader:ie.linedashed_frag},depth:{uniforms:Xe([At.common,At.displacementmap]),vertexShader:ie.depth_vert,fragmentShader:ie.depth_frag},normal:{uniforms:Xe([At.common,At.bumpmap,At.normalmap,At.displacementmap,{opacity:{value:1}}]),vertexShader:ie.meshnormal_vert,fragmentShader:ie.meshnormal_frag},sprite:{uniforms:Xe([At.sprite,At.fog]),vertexShader:ie.sprite_vert,fragmentShader:ie.sprite_frag},background:{uniforms:{uvTransform:{value:new Xt},t2D:{value:null},backgroundIntensity:{value:1}},vertexShader:ie.background_vert,fragmentShader:ie.background_frag},backgroundCube:{uniforms:{envMap:{value:null},backgroundBlurriness:{value:0},backgroundIntensity:{value:1},backgroundRotation:{value:new Xt}},vertexShader:ie.backgroundCube_vert,fragmentShader:ie.backgroundCube_frag},cube:{uniforms:{tCube:{value:null},tFlip:{value:-1},opacity:{value:1}},vertexShader:ie.cube_vert,fragmentShader:ie.cube_frag},equirect:{uniforms:{tEquirect:{value:null}},vertexShader:ie.equirect_vert,fragmentShader:ie.equirect_frag},distance:{uniforms:Xe([At.common,At.displacementmap,{referencePosition:{value:new L},nearDistance:{value:1},farDistance:{value:1000}}]),vertexShader:ie.distance_vert,fragmentShader:ie.distance_frag},shadow:{uniforms:Xe([At.lights,At.fog,{color:{value:new Pt(0)},opacity:{value:1}}]),vertexShader:ie.shadow_vert,fragmentShader:ie.shadow_frag}};On.physical={uniforms:Xe([On.standard.uniforms,{clearcoat:{value:0},clearcoatMap:{value:null},clearcoatMapTransform:{value:new Xt},clearcoatNormalMap:{value:null},clearcoatNormalMapTransform:{value:new Xt},clearcoatNormalScale:{value:new rt(1,1)},clearcoatRoughness:{value:0},clearcoatRoughnessMap:{value:null},clearcoatRoughnessMapTransform:{value:new Xt},dispersion:{value:0},retroreflectivity:{value:0},iridescence:{value:0},iridescenceMap:{value:null},iridescenceMapTransform:{value:new Xt},iridescenceIOR:{value:1.3},iridescenceThicknessMinimum:{value:100},iridescenceThicknessMaximum:{value:400},iridescenceThicknessMap:{value:null},iridescenceThicknessMapTransform:{value:new Xt},sheen:{value:0},sheenColor:{value:new Pt(0)},sheenColorMap:{value:null},sheenColorMapTransform:{value:new Xt},sheenRoughness:{value:1},sheenRoughnessMap:{value:null},sheenRoughnessMapTransform:{value:new Xt},transmission:{value:0},transmissionMap:{value:null},transmissionMapTransform:{value:new Xt},transmissionSamplerSize:{value:new rt},transmissionSamplerMap:{value:null},thickness:{value:0},thicknessMap:{value:null},thicknessMapTransform:{value:new Xt},attenuationDistance:{value:0},attenuationColor:{value:new Pt(0)},specularColor:{value:new Pt(1,1,1)},specularColorMap:{value:null},specularColorMapTransform:{value:new Xt},specularIntensity:{value:1},specularIntensityMap:{value:null},specularIntensityMapTransform:{value:new Xt},anisotropyVector:{value:new rt},anisotropyMap:{value:null},anisotropyMapTransform:{value:new Xt}}]),vertexShader:ie.meshphysical_vert,fragmentShader:ie.meshphysical_frag};var Lo={r:0,b:0,g:0},sv=new Kt,Wp=new Xt;Wp.set(-1,0,0,0,1,0,0,0,1);function rv(t,e,n,i,s,r){let a=new Pt(0),o=s===!0?0:1,l,c,h=null,d=0,u=null;function f(b){let S=b.isScene===!0?b.background:null;if(S&&S.isTexture){let x=b.backgroundBlurriness>0;S=e.get(S,x)}return S}function g(b){let S=!1,x=f(b);if(x===null)m(a,o);else if(x&&x.isColor)m(x,1),S=!0;let T=t.xr.getEnvironmentBlendMode();if(T==="additive")n.buffers.color.setClear(0,0,0,1,r);else if(T==="alpha-blend")n.buffers.color.setClear(0,0,0,0,r);if(t.autoClear||S)n.buffers.depth.setTest(!0),n.buffers.depth.setMask(!0),n.buffers.color.setMask(!0),t.clear(t.autoClearColor,t.autoClearDepth,t.autoClearStencil)}function _(b,S){let x=f(S);if(x&&(x.isCubeTexture||x.mapping===ar)){if(c===void 0)c=new pe(new An(1,1,1),new ln({name:"BackgroundCubeMaterial",uniforms:Vi(On.backgroundCube.uniforms),vertexShader:On.backgroundCube.vertexShader,fragmentShader:On.backgroundCube.fragmentShader,side:Ve,depthTest:!1,depthWrite:!1,fog:!1,allowOverride:!1})),c.geometry.deleteAttribute("normal"),c.geometry.deleteAttribute("uv"),c.onBeforeRender=function(T,I,R){this.matrixWorld.copyPosition(R.matrixWorld)},Object.defineProperty(c.material,"envMap",{get:function(){return this.uniforms.envMap.value}}),i.update(c);if(c.material.uniforms.envMap.value=x,c.material.uniforms.backgroundBlurriness.value=S.backgroundBlurriness,c.material.uniforms.backgroundIntensity.value=S.backgroundIntensity,c.material.uniforms.backgroundRotation.value.setFromMatrix4(sv.makeRotationFromEuler(S.backgroundRotation)).transpose(),x.isCubeTexture&&x.isRenderTargetTexture===!1)c.material.uniforms.backgroundRotation.value.premultiply(Wp);if(c.material.toneMapped=re.getTransfer(x.colorSpace)!==xe,h!==x||d!==x.version||u!==t.toneMapping)c.material.needsUpdate=!0,h=x,d=x.version,u=t.toneMapping;c.layers.enableAll(),b.unshift(c,c.geometry,c.material,0,0,null)}else if(x&&x.isTexture){if(l===void 0)l=new pe(new Is(2,2),new ln({name:"BackgroundMaterial",uniforms:Vi(On.background.uniforms),vertexShader:On.background.vertexShader,fragmentShader:On.background.fragmentShader,side:Ms,depthTest:!1,depthWrite:!1,fog:!1,allowOverride:!1})),l.geometry.deleteAttribute("normal"),Object.defineProperty(l.material,"map",{get:function(){return this.uniforms.t2D.value}}),i.update(l);if(l.material.uniforms.t2D.value=x,l.material.uniforms.backgroundIntensity.value=S.backgroundIntensity,l.material.toneMapped=re.getTransfer(x.colorSpace)!==xe,x.matrixAutoUpdate===!0)x.updateMatrix();if(l.material.uniforms.uvTransform.value.copy(x.matrix),h!==x||d!==x.version||u!==t.toneMapping)l.material.needsUpdate=!0,h=x,d=x.version,u=t.toneMapping;l.layers.enableAll(),b.unshift(l,l.geometry,l.material,0,0,null)}}function m(b,S){b.getRGB(Lo,Xc(t)),n.buffers.color.setClear(Lo.r,Lo.g,Lo.b,S,r)}function p(){if(c!==void 0)c.geometry.dispose(),c.material.dispose(),c=void 0;if(l!==void 0)l.geometry.dispose(),l.material.dispose(),l=void 0}return{getClearColor:function(){return a},setClearColor:function(b,S=1){a.set(b),o=S,m(a,o)},getClearAlpha:function(){return o},setClearAlpha:function(b){o=b,m(a,o)},render:g,addToRenderList:_,dispose:p}}function av(t,e){let n=t.getParameter(t.MAX_VERTEX_ATTRIBS),i={},s=u(null),r=s,a=!1;function o(D,q,ot,z,et){let nt=!1,B=d(D,z,ot,q);if(r!==B)r=B,c(r.object);if(nt=f(D,z,ot,et),nt)g(D,z,ot,et);if(et!==null)e.update(et,t.ELEMENT_ARRAY_BUFFER);if(nt||a){if(a=!1,x(D,q,ot,z),et!==null)t.bindBuffer(t.ELEMENT_ARRAY_BUFFER,e.get(et).buffer)}}function l(){return t.createVertexArray()}function c(D){return t.bindVertexArray(D)}function h(D){return t.deleteVertexArray(D)}function d(D,q,ot,z){let et=z.wireframe===!0,nt=i[q.id];if(nt===void 0)nt={},i[q.id]=nt;let B=D.isInstancedMesh===!0?D.id:0,ft=nt[B];if(ft===void 0)ft={},nt[B]=ft;let Q=ft[ot.id];if(Q===void 0)Q={},ft[ot.id]=Q;let ct=Q[et];if(ct===void 0)ct=u(l()),Q[et]=ct;return ct}function u(D){let q=[],ot=[],z=[];for(let et=0;et<n;et++)q[et]=0,ot[et]=0,z[et]=0;return{geometry:null,program:null,wireframe:!1,newAttributes:q,enabledAttributes:ot,attributeDivisors:z,object:D,attributes:{},index:null}}function f(D,q,ot,z){let et=r.attributes,nt=q.attributes,B=0,ft=ot.getAttributes();for(let Q in ft)if(ft[Q].location>=0){let dt=et[Q],Bt=nt[Q];if(Bt===void 0){if(Q==="instanceMatrix"&&D.instanceMatrix)Bt=D.instanceMatrix;if(Q==="instanceColor"&&D.instanceColor)Bt=D.instanceColor}if(dt===void 0)return!0;if(dt.attribute!==Bt)return!0;if(Bt&&dt.data!==Bt.data)return!0;B++}if(r.attributesNum!==B)return!0;if(r.index!==z)return!0;return!1}function g(D,q,ot,z){let et={},nt=q.attributes,B=0,ft=ot.getAttributes();for(let Q in ft)if(ft[Q].location>=0){let dt=nt[Q];if(dt===void 0){if(Q==="instanceMatrix"&&D.instanceMatrix)dt=D.instanceMatrix;if(Q==="instanceColor"&&D.instanceColor)dt=D.instanceColor}let Bt={};if(Bt.attribute=dt,dt&&dt.data)Bt.data=dt.data;et[Q]=Bt,B++}r.attributes=et,r.attributesNum=B,r.index=z}function _(){let D=r.newAttributes;for(let q=0,ot=D.length;q<ot;q++)D[q]=0}function m(D){p(D,0)}function p(D,q){let ot=r.newAttributes,z=r.enabledAttributes,et=r.attributeDivisors;if(ot[D]=1,z[D]===0)t.enableVertexAttribArray(D),z[D]=1;if(et[D]!==q)t.vertexAttribDivisor(D,q),et[D]=q}function b(){let D=r.newAttributes,q=r.enabledAttributes;for(let ot=0,z=q.length;ot<z;ot++)if(q[ot]!==D[ot])t.disableVertexAttribArray(ot),q[ot]=0}function S(D,q,ot,z,et,nt,B){if(B===!0)t.vertexAttribIPointer(D,q,ot,et,nt);else t.vertexAttribPointer(D,q,ot,z,et,nt)}function x(D,q,ot,z){_();let et=z.attributes,nt=ot.getAttributes(),B=q.defaultAttributeValues;for(let ft in nt){let Q=nt[ft];if(Q.location>=0){let ct=et[ft];if(ct===void 0){if(ft==="instanceMatrix"&&D.instanceMatrix)ct=D.instanceMatrix;if(ft==="instanceColor"&&D.instanceColor)ct=D.instanceColor}if(ct!==void 0){let dt=ct.normalized,Bt=ct.itemSize,Dt=e.get(ct);if(Dt===void 0)continue;let{buffer:ee,type:Jt,bytesPerElement:F}=Dt,O=Jt===t.INT||Jt===t.UNSIGNED_INT||ct.gpuType===Bl;if(ct.isInterleavedBufferAttribute){let N=ct.data,P=N.stride,G=ct.offset;if(N.isInstancedInterleavedBuffer){for(let J=0;J<Q.locationSize;J++)p(Q.location+J,N.meshPerAttribute);if(D.isInstancedMesh!==!0&&z._maxInstanceCount===void 0)z._maxInstanceCount=N.meshPerAttribute*N.count}else for(let J=0;J<Q.locationSize;J++)m(Q.location+J);t.bindBuffer(t.ARRAY_BUFFER,ee);for(let J=0;J<Q.locationSize;J++)S(Q.location+J,Bt/Q.locationSize,Jt,dt,P*F,(G+Bt/Q.locationSize*J)*F,O)}else{if(ct.isInstancedBufferAttribute){for(let N=0;N<Q.locationSize;N++)p(Q.location+N,ct.meshPerAttribute);if(D.isInstancedMesh!==!0&&z._maxInstanceCount===void 0)z._maxInstanceCount=ct.meshPerAttribute*ct.count}else for(let N=0;N<Q.locationSize;N++)m(Q.location+N);t.bindBuffer(t.ARRAY_BUFFER,ee);for(let N=0;N<Q.locationSize;N++)S(Q.location+N,Bt/Q.locationSize,Jt,dt,Bt*F,Bt/Q.locationSize*N*F,O)}}else if(B!==void 0){let dt=B[ft];if(dt!==void 0)switch(dt.length){case 2:t.vertexAttrib2fv(Q.location,dt);break;case 3:t.vertexAttrib3fv(Q.location,dt);break;case 4:t.vertexAttrib4fv(Q.location,dt);break;default:t.vertexAttrib1fv(Q.location,dt)}}}}b()}function T(){A();for(let D in i){let q=i[D];for(let ot in q){let z=q[ot];for(let et in z){let nt=z[et];for(let B in nt)h(nt[B].object),delete nt[B];delete z[et]}}delete i[D]}}function I(D){if(i[D.id]===void 0)return;let q=i[D.id];for(let ot in q){let z=q[ot];for(let et in z){let nt=z[et];for(let B in nt)h(nt[B].object),delete nt[B];delete z[et]}}delete i[D.id]}function R(D){for(let q in i){let ot=i[q];for(let z in ot){let et=ot[z];if(et[D.id]===void 0)continue;let nt=et[D.id];for(let B in nt)h(nt[B].object),delete nt[B];delete et[D.id]}}}function y(D){for(let q in i){let ot=i[q],z=D.isInstancedMesh===!0?D.id:0,et=ot[z];if(et===void 0)continue;for(let nt in et){let B=et[nt];for(let ft in B)h(B[ft].object),delete B[ft];delete et[nt]}if(delete ot[z],Object.keys(ot).length===0)delete i[q]}}function A(){if(Z(),a=!0,r===s)return;r=s,c(r.object)}function Z(){s.geometry=null,s.program=null,s.wireframe=!1}return{setup:o,reset:A,resetDefaultState:Z,dispose:T,releaseStatesOfGeometry:I,releaseStatesOfObject:y,releaseStatesOfProgram:R,initAttributes:_,enableAttribute:m,disableUnusedAttributes:b}}function ov(t,e,n){let i;function s(l){i=l}function r(l,c){t.drawArrays(i,l,c),n.update(c,i,1)}function a(l,c,h){if(h===0)return;t.drawArraysInstanced(i,l,c,h),n.update(c,i,h)}function o(l,c,h){if(h===0)return;e.get("WEBGL_multi_draw").multiDrawArraysWEBGL(i,l,0,c,0,h);let u=0;for(let f=0;f<h;f++)u+=c[f];n.update(u,i,1)}this.setMode=s,this.render=r,this.renderInstances=a,this.renderMultiDraw=o}function lv(t,e,n,i){let s;function r(){if(s!==void 0)return s;if(e.has("EXT_texture_filter_anisotropic")===!0){let R=e.get("EXT_texture_filter_anisotropic");s=t.getParameter(R.MAX_TEXTURE_MAX_ANISOTROPY_EXT)}else s=0;return s}function a(R){if(R!==Dn&&i.convert(R)!==t.getParameter(t.IMPLEMENTATION_COLOR_READ_FORMAT))return!1;return!0}function o(R){let y=R===Nn&&(e.has("EXT_color_buffer_half_float")||e.has("EXT_color_buffer_float"));if(R!==Tn&&R!==$n&&!y&&i.convert(R)!==t.getParameter(t.IMPLEMENTATION_COLOR_READ_TYPE))return!1;return!0}function l(R){if(R==="highp"){if(t.getShaderPrecisionFormat(t.VERTEX_SHADER,t.HIGH_FLOAT).precision>0&&t.getShaderPrecisionFormat(t.FRAGMENT_SHADER,t.HIGH_FLOAT).precision>0)return"highp";R="mediump"}if(R==="mediump"){if(t.getShaderPrecisionFormat(t.VERTEX_SHADER,t.MEDIUM_FLOAT).precision>0&&t.getShaderPrecisionFormat(t.FRAGMENT_SHADER,t.MEDIUM_FLOAT).precision>0)return"mediump"}return"lowp"}let c=n.precision!==void 0?n.precision:"highp",h=l(c);if(h!==c)Mt("WebGLRenderer:",c,"not supported, using",h,"instead."),c=h;let d=n.logarithmicDepthBuffer===!0,u=n.reversedDepthBuffer===!0&&e.has("EXT_clip_control");if(n.reversedDepthBuffer===!0&&u===!1)Mt("WebGLRenderer: Unable to use reversed depth buffer due to missing EXT_clip_control extension. Fallback to default depth buffer.");let f=t.getParameter(t.MAX_TEXTURE_IMAGE_UNITS),g=t.getParameter(t.MAX_VERTEX_TEXTURE_IMAGE_UNITS),_=t.getParameter(t.MAX_TEXTURE_SIZE),m=t.getParameter(t.MAX_CUBE_MAP_TEXTURE_SIZE),p=t.getParameter(t.MAX_VERTEX_ATTRIBS),b=t.getParameter(t.MAX_VERTEX_UNIFORM_VECTORS),S=t.getParameter(t.MAX_VARYING_VECTORS),x=t.getParameter(t.MAX_FRAGMENT_UNIFORM_VECTORS),T=t.getParameter(t.MAX_SAMPLES),I=t.getParameter(t.SAMPLES);return{isWebGL2:!0,getMaxAnisotropy:r,getMaxPrecision:l,textureFormatReadable:a,textureTypeReadable:o,precision:c,logarithmicDepthBuffer:d,reversedDepthBuffer:u,maxTextures:f,maxVertexTextures:g,maxTextureSize:_,maxCubemapSize:m,maxAttributes:p,maxVertexUniforms:b,maxVaryings:S,maxFragmentUniforms:x,maxSamples:T,samples:I}}function cv(t){let e=this,n=null,i=0,s=!1,r=!1,a=new dn,o=new Xt,l={value:null,needsUpdate:!1};this.uniform=l,this.numPlanes=0,this.numIntersection=0,this.init=function(d,u){let f=d.length!==0||u||i!==0||s;return s=u,i=d.length,f},this.beginShadows=function(){r=!0,h(null)},this.endShadows=function(){r=!1},this.setGlobalState=function(d,u){n=h(d,u,0)},this.setState=function(d,u,f){let{clippingPlanes:g,clipIntersection:_,clipShadows:m}=d,p=t.get(d);if(!s||g===null||g.length===0||r&&!m)if(r)h(null);else c();else{let b=r?0:i,S=b*4,x=p.clippingState||null;l.value=x,x=h(g,u,S,f);for(let T=0;T!==S;++T)x[T]=n[T];p.clippingState=x,this.numIntersection=_?this.numPlanes:0,this.numPlanes+=b}};function c(){if(l.value!==n)l.value=n,l.needsUpdate=i>0;e.numPlanes=i,e.numIntersection=0}function h(d,u,f,g){let _=d!==null?d.length:0,m=null;if(_!==0){if(m=l.value,g!==!0||m===null){let p=f+_*4,b=u.matrixWorldInverse;if(o.getNormalMatrix(b),m===null||m.length<p)m=new Float32Array(p);for(let S=0,x=f;S!==_;++S,x+=4)a.copy(d[S]).applyMatrix4(b,o),a.normal.toArray(m,x),m[x+3]=a.constant}l.value=m,l.needsUpdate=!0}return e.numPlanes=_,e.numIntersection=0,m}}var Ns=4,hv=6,uv=20,dv=256,Cr=new Ps,Mp=new Pt,vh=null,Sh=0,Mh=0,bh=!1,fv=new L,qi=new L;class Ah{constructor(t){this._renderer=t,this._pingPongRenderTarget=null,this._lodMax=0,this._cubeSize=0,this._sizeLods=[],this._lodMeshes=[],this._backgroundBox=null,this._cubemapMaterial=null,this._equirectMaterial=null,this._blurMaterial=null,this._ggxMaterial=null}fromScene(t,e=0,n=0.1,i=100,s={}){let{size:r=256,position:a=fv}=s;vh=this._renderer.getRenderTarget(),Sh=this._renderer.getActiveCubeFace(),Mh=this._renderer.getActiveMipmapLevel(),bh=this._renderer.xr.enabled,this._renderer.xr.enabled=!1,this._setSize(r);let o=this._allocateTargets();if(o.depthBuffer=!0,this._sceneToCubeUV(t,n,i,o,a),e>0)this._blur(o,0,0,e);return this._applyPMREM(o),this._cleanup(o),o}fromEquirectangular(t,e=null){return this._fromTexture(t,e)}fromCubemap(t,e=null){return this._fromTexture(t,e)}compileCubemapShader(){if(this._cubemapMaterial===null)this._cubemapMaterial=Ep(),this._compileMaterial(this._cubemapMaterial)}compileEquirectangularShader(){if(this._equirectMaterial===null)this._equirectMaterial=Tp(),this._compileMaterial(this._equirectMaterial)}dispose(){if(this._dispose(),this._cubemapMaterial!==null)this._cubemapMaterial.dispose();if(this._equirectMaterial!==null)this._equirectMaterial.dispose();if(this._backgroundBox!==null)this._backgroundBox.geometry.dispose(),this._backgroundBox.material.dispose()}_setSize(t){this._lodMax=Math.floor(Math.log2(t)),this._cubeSize=Math.pow(2,this._lodMax)}_dispose(){if(this._blurMaterial!==null)this._blurMaterial.dispose();if(this._ggxMaterial!==null)this._ggxMaterial.dispose();if(this._pingPongRenderTarget!==null)this._pingPongRenderTarget.dispose();for(let t=0;t<this._lodMeshes.length;t++)this._lodMeshes[t].geometry.dispose()}_cleanup(t){this._renderer.setRenderTarget(vh,Sh,Mh),this._renderer.xr.enabled=bh,t.scissorTest=!1,Ls(t,0,0,t.width,t.height)}_fromTexture(t,e){if(t.mapping===Ts||t.mapping===Fi)this._setSize(t.image.length===0?16:t.image[0].width||t.image[0].image.width);else this._setSize(t.image.width/4);vh=this._renderer.getRenderTarget(),Sh=this._renderer.getActiveCubeFace(),Mh=this._renderer.getActiveMipmapLevel(),bh=this._renderer.xr.enabled,this._renderer.xr.enabled=!1;let n=e||this._allocateTargets();return this._textureToCubeUV(t,n),this._applyPMREM(n),this._cleanup(n),n}_allocateTargets(){let t=3*Math.max(this._cubeSize,112),e=4*this._cubeSize,n={magFilter:Qe,minFilter:Qe,generateMipmaps:!1,type:Nn,format:Dn,colorSpace:Sc,depthBuffer:!1},i=bp(t,e,n);if(this._pingPongRenderTarget===null||this._pingPongRenderTarget.width!==t||this._pingPongRenderTarget.height!==e){if(this._pingPongRenderTarget!==null)this._dispose();this._pingPongRenderTarget=bp(t,e,n);let{_lodMax:s}=this;({lodMeshes:this._lodMeshes,sizeLods:this._sizeLods}=pv(s)),this._blurMaterial=gv(s,t,e),this._ggxMaterial=mv(s,t,e)}return i}_compileMaterial(t){let e=new pe(new jt,t);this._renderer.compile(e,Cr)}_sceneToCubeUV(t,e,n,i,s){let o=new De(90,1,e,n),l=[1,-1,1,1,1,1],c=[1,1,1,-1,-1,-1],h=this._renderer,{autoClear:d,toneMapping:u}=h;if(h.getClearColor(Mp),h.toneMapping=bn,h.autoClear=!1,h.state.buffers.depth.getReversed())h.setRenderTarget(i),h.clearDepth(),h.setRenderTarget(null);if(this._backgroundBox===null)this._backgroundBox=new pe(new An,new mn({name:"PMREM.Background",side:Ve,depthWrite:!1,depthTest:!1}));let g=this._backgroundBox,_=g.material,m=!1,p=t.background;if(p){if(p.isColor)_.color.copy(p),t.background=null,m=!0}else _.color.copy(Mp),m=!0;for(let b=0;b<6;b++){let S=b%3;if(S===0)o.up.set(0,l[b],0),o.position.set(s.x,s.y,s.z),o.lookAt(s.x+c[b],s.y,s.z);else if(S===1)o.up.set(0,0,l[b]),o.position.set(s.x,s.y,s.z),o.lookAt(s.x,s.y+c[b],s.z);else o.up.set(0,l[b],0),o.position.set(s.x,s.y,s.z),o.lookAt(s.x,s.y,s.z+c[b]);let x=this._cubeSize;if(Ls(i,S*x,b>2?x:0,x,x),h.setRenderTarget(i),m)h.render(g,o);h.render(t,o)}h.toneMapping=u,h.autoClear=d,t.background=p}_textureToCubeUV(t,e){let n=this._renderer,i=t.mapping===Ts||t.mapping===Fi;if(i){if(this._cubemapMaterial===null)this._cubemapMaterial=Ep();this._cubemapMaterial.uniforms.flipEnvMap.value=t.isRenderTargetTexture===!1?-1:1}else if(this._equirectMaterial===null)this._equirectMaterial=Tp();let s=i?this._cubemapMaterial:this._equirectMaterial,r=this._lodMeshes[0];r.material=s;let a=s.uniforms;a.envMap.value=t;let o=this._cubeSize;Ls(e,0,0,3*o,2*o),n.setRenderTarget(e),n.render(r,Cr)}_applyPMREM(t){let e=this._renderer,n=e.autoClear;e.autoClear=!1;let i=this._lodMeshes.length;for(let s=1;s<i;s++)this._applyGGXFilter(t,s-1,s);e.autoClear=n}_applyGGXFilter(t,e,n){let i=this._renderer,s=this._pingPongRenderTarget,r=this._ggxMaterial,a=this._lodMeshes[n];a.material=r;let o=r.uniforms,l=n/(this._lodMeshes.length-1),c=e/(this._lodMeshes.length-1),h=Math.sqrt(l*l-c*c),d=l*1.25,u=h*d,{_lodMax:f}=this,g=this._sizeLods[n],_=3*g*(n>f-Ns?n-f+Ns:0),m=4*(this._cubeSize-g);o.envMap.value=t.texture,o.roughness.value=u,o.mipInt.value=f-e,Ls(s,_,m,3*g,2*g),i.setRenderTarget(s),i.render(a,Cr),o.envMap.value=s.texture,o.roughness.value=0,o.mipInt.value=f-n,Ls(t,_,m,3*g,2*g),i.setRenderTarget(t),i.render(a,Cr)}_blur(t,e,n,i){let s=this._pingPongRenderTarget,r=Math.min(i,Math.PI)/Math.SQRT2;this._blurPass(t,s,e,n,r),this._blurPass(s,t,n,n,r)}_blurPass(t,e,n,i,s){let r=this._renderer,a=this._blurMaterial,o=this._lodMeshes[i];o.material=a;let l=a.uniforms;l.envMap.value=t.texture,l.sigma.value=s,l.mipInt.value=this._lodMax-n;let c=this._sizeLods[i],h=3*c*(i>this._lodMax-Ns?i-this._lodMax+Ns:0),d=4*(this._cubeSize-c);Ls(e,h,d,3*c,2*c),r.setRenderTarget(e),r.render(o,Cr)}}function pv(t){let e=[],n=[],i=t,s=t-Ns+1+hv;for(let r=0;r<s;r++){let a=Math.pow(2,i);e.push(a);let o=1/(a-2),l=-o,c=1+o,h=[l,l,c,l,c,c,l,l,c,c,l,c],d=6,u=6,f=3,g=new Float32Array(f*u*d),_=new Float32Array(f*u*d);for(let p=0;p<d;p++){let b=p%3*2/3-1,S=p>2?0:-1,x=[b,S,0,b+0.6666666666666666,S,0,b+0.6666666666666666,S+1,0,b,S,0,b+0.6666666666666666,S+1,0,b,S+1,0];g.set(x,f*u*p);for(let T=0;T<u;T++){let I=h[T*2]*2-1,R=h[T*2+1]*2-1;if(p===0)qi.set(1,R,I);else if(p===1)qi.set(-I,1,-R);else if(p===2)qi.set(-I,R,1);else if(p===3)qi.set(-1,R,-I);else if(p===4)qi.set(-I,-1,R);else qi.set(I,R,-1);qi.toArray(_,(p*u+T)*f)}}let m=new jt;if(m.setAttribute("position",new de(g,f)),m.setAttribute("outputDirection",new de(_,f)),n.push(new pe(m,null)),i>Ns)i--}return{lodMeshes:n,sizeLods:e}}function bp(t,e,n){let i=new Je(t,e,n);return i.texture.mapping=ar,i.texture.name="PMREM.cubeUv",i.scissorTest=!0,i}function Ls(t,e,n,i,s){t.viewport.set(e,n,i,s),t.scissor.set(e,n,i,s)}function mv(t,e,n){return new ln({name:"PMREMGGXConvolution",defines:{GGX_SAMPLES:dv,CUBEUV_TEXEL_WIDTH:1/e,CUBEUV_TEXEL_HEIGHT:1/n,CUBEUV_MAX_MIP:`${t}.0`},uniforms:{envMap:{value:null},roughness:{value:0},mipInt:{value:0}},vertexShader:Do(),fragmentShader:`

			precision highp float;
			precision highp int;

			varying vec3 vOutputDirection;

			uniform sampler2D envMap;
			uniform float roughness;
			uniform float mipInt;

			#define ENVMAP_TYPE_CUBE_UV
			#include <cube_uv_reflection_fragment>

			#define PI 3.14159265359

			// Van der Corput radical inverse
			float radicalInverse_VdC(uint bits) {
				bits = (bits << 16u) | (bits >> 16u);
				bits = ((bits & 0x55555555u) << 1u) | ((bits & 0xAAAAAAAAu) >> 1u);
				bits = ((bits & 0x33333333u) << 2u) | ((bits & 0xCCCCCCCCu) >> 2u);
				bits = ((bits & 0x0F0F0F0Fu) << 4u) | ((bits & 0xF0F0F0F0u) >> 4u);
				bits = ((bits & 0x00FF00FFu) << 8u) | ((bits & 0xFF00FF00u) >> 8u);
				return float(bits) * 2.3283064365386963e-10; // / 0x100000000
			}

			// Hammersley sequence
			vec2 hammersley(uint i, uint N) {
				return vec2(float(i) / float(N), radicalInverse_VdC(i));
			}

			// GGX VNDF importance sampling (Eric Heitz 2018)
			// "Sampling the GGX Distribution of Visible Normals"
			// https://jcgt.org/published/0007/04/01/
			vec3 importanceSampleGGX_VNDF(vec2 Xi, vec3 V, float roughness) {
				float alpha = roughness * roughness;

				// Section 4.1: Orthonormal basis
				vec3 T1 = vec3(1.0, 0.0, 0.0);
				vec3 T2 = cross(V, T1);

				// Section 4.2: Parameterization of projected area
				float r = sqrt(Xi.x);
				float phi = 2.0 * PI * Xi.y;
				float t1 = r * cos(phi);
				float t2 = r * sin(phi);
				float s = 0.5 * (1.0 + V.z);
				t2 = (1.0 - s) * sqrt(1.0 - t1 * t1) + s * t2;

				// Section 4.3: Reprojection onto hemisphere
				vec3 Nh = t1 * T1 + t2 * T2 + sqrt(max(0.0, 1.0 - t1 * t1 - t2 * t2)) * V;

				// Section 3.4: Transform back to ellipsoid configuration
				return normalize(vec3(alpha * Nh.x, alpha * Nh.y, max(0.0, Nh.z)));
			}

			void main() {
				vec3 N = normalize(vOutputDirection);
				vec3 V = N; // Assume view direction equals normal for pre-filtering

				vec3 prefilteredColor = vec3(0.0);
				float totalWeight = 0.0;

				// For very low roughness, just sample the environment directly
				if (roughness < 0.001) {
					gl_FragColor = vec4(bilinearCubeUV(envMap, N, mipInt), 1.0);
					return;
				}

				// Tangent space basis for VNDF sampling
				vec3 up = abs(N.z) < 0.999 ? vec3(0.0, 0.0, 1.0) : vec3(1.0, 0.0, 0.0);
				vec3 tangent = normalize(cross(up, N));
				vec3 bitangent = cross(N, tangent);

				for(uint i = 0u; i < uint(GGX_SAMPLES); i++) {
					vec2 Xi = hammersley(i, uint(GGX_SAMPLES));

					// For PMREM, V = N, so in tangent space V is always (0, 0, 1)
					vec3 H_tangent = importanceSampleGGX_VNDF(Xi, vec3(0.0, 0.0, 1.0), roughness);

					// Transform H back to world space
					vec3 H = normalize(tangent * H_tangent.x + bitangent * H_tangent.y + N * H_tangent.z);
					vec3 L = normalize(2.0 * dot(V, H) * H - V);

					float NdotL = max(dot(N, L), 0.0);

					if(NdotL > 0.0) {
						// Sample environment at fixed mip level
						// VNDF importance sampling handles the distribution filtering
						vec3 sampleColor = bilinearCubeUV(envMap, L, mipInt);

						// Weight by NdotL for the split-sum approximation
						// VNDF PDF naturally accounts for the visible microfacet distribution
						prefilteredColor += sampleColor * NdotL;
						totalWeight += NdotL;
					}
				}

				if (totalWeight > 0.0) {
					prefilteredColor = prefilteredColor / totalWeight;
				}

				gl_FragColor = vec4(prefilteredColor, 1.0);
			}
		`,blending:Ln,depthTest:!1,depthWrite:!1})}function gv(t,e,n){return new ln({name:"SphericalGaussianBlur",defines:{SAMPLES:uv,CUBEUV_TEXEL_WIDTH:1/e,CUBEUV_TEXEL_HEIGHT:1/n,CUBEUV_MAX_MIP:`${t}.0`},uniforms:{envMap:{value:null},sigma:{value:0},mipInt:{value:0}},vertexShader:Do(),fragmentShader:`

			precision highp float;
			precision highp int;

			varying vec3 vOutputDirection;

			uniform sampler2D envMap;
			uniform float sigma;
			uniform float mipInt;

			#define ENVMAP_TYPE_CUBE_UV
			#include <cube_uv_reflection_fragment>

			#define PI 3.14159265359
			#define GOLDEN_ANGLE 2.39996322973

			void main() {

				if ( sigma == 0.0 ) {

					gl_FragColor = vec4( bilinearCubeUV( envMap, vOutputDirection, mipInt ), 1.0 );
					return;

				}

				vec3 outputDirection = normalize( vOutputDirection );

				vec3 up = abs( outputDirection.z ) < 0.999 ? vec3( 0.0, 0.0, 1.0 ) : vec3( 1.0, 0.0, 0.0 );
				vec3 tangent = normalize( cross( up, outputDirection ) );
				vec3 bitangent = cross( outputDirection, tangent );

				// Truncate the kernel at three standard deviations or at the antipode.
				float thetaMax = min( 3.0 * sigma, PI );
				float truncation = 1.0 - exp( - 0.5 * thetaMax * thetaMax / ( sigma * sigma ) );

				vec3 accumColor = vec3( 0.0 );
				float accumWeight = 0.0;

				for ( int i = 0; i < SAMPLES; i ++ ) {

					// Stratified inverse-CDF sampling of the Gaussian, placed on a golden-angle spiral.
					float stratum = ( float( i ) + 0.5 ) / float( SAMPLES );
					float theta = sigma * sqrt( - 2.0 * log( 1.0 - stratum * truncation ) );
					float phi = float( i ) * GOLDEN_ANGLE;

					vec3 offset = cos( phi ) * tangent + sin( phi ) * bitangent;
					vec3 sampleDirection = cos( theta ) * outputDirection + sin( theta ) * offset;

					// Correct the planar sample density to solid angle.
					float weight = sin( theta ) / theta;

					accumColor += weight * bilinearCubeUV( envMap, sampleDirection, mipInt );
					accumWeight += weight;

				}

				gl_FragColor = vec4( accumColor / accumWeight, 1.0 );

			}
		`,blending:Ln,depthTest:!1,depthWrite:!1})}function Tp(){return new ln({name:"EquirectangularToCubeUV",uniforms:{envMap:{value:null}},vertexShader:Do(),fragmentShader:`

			precision mediump float;
			precision mediump int;

			varying vec3 vOutputDirection;

			uniform sampler2D envMap;

			#include <common>

			void main() {

				vec3 outputDirection = normalize( vOutputDirection );
				vec2 uv = equirectUv( outputDirection );

				gl_FragColor = vec4( texture2D ( envMap, uv ).rgb, 1.0 );

			}
		`,blending:Ln,depthTest:!1,depthWrite:!1})}function Ep(){return new ln({name:"CubemapToCubeUV",uniforms:{envMap:{value:null},flipEnvMap:{value:-1}},vertexShader:Do(),fragmentShader:`

			precision mediump float;
			precision mediump int;

			uniform float flipEnvMap;

			varying vec3 vOutputDirection;

			uniform samplerCube envMap;

			void main() {

				gl_FragColor = textureCube( envMap, vec3( flipEnvMap * vOutputDirection.x, vOutputDirection.yz ) );

			}
		`,blending:Ln,depthTest:!1,depthWrite:!1})}function Do(){return`

		precision mediump float;
		precision mediump int;

		attribute vec3 outputDirection;

		varying vec3 vOutputDirection;

		void main() {

			vOutputDirection = outputDirection;
			gl_Position = vec4( position, 1.0 );

		}
	`}class Ch extends Je{constructor(t=1,e={}){super(t,t,e);this.isWebGLCubeRenderTarget=!0;let n={width:t,height:t,depth:1},i=[n,n,n,n,n,n];this.texture=new Rs(i),this._setTextureOptions(e),this.texture.isRenderTargetTexture=!0}fromEquirectangularTexture(t,e){this.texture.type=e.type,this.texture.colorSpace=e.colorSpace,this.texture.generateMipmaps=e.generateMipmaps,this.texture.minFilter=e.minFilter,this.texture.magFilter=e.magFilter;let n={uniforms:{tEquirect:{value:null}},vertexShader:`

				varying vec3 vWorldDirection;

				vec3 transformDirection( in vec3 dir, in mat4 matrix ) {

					return normalize( ( matrix * vec4( dir, 0.0 ) ).xyz );

				}

				void main() {

					vWorldDirection = transformDirection( position, modelMatrix );

					#include <begin_vertex>
					#include <project_vertex>

				}
			`,fragmentShader:`

				uniform sampler2D tEquirect;

				varying vec3 vWorldDirection;

				#include <common>

				void main() {

					vec3 direction = normalize( vWorldDirection );

					vec2 sampleUV = equirectUv( direction );

					gl_FragColor = texture2D( tEquirect, sampleUV );

				}
			`},i=new An(5,5,5),s=new ln({name:"CubemapFromEquirect",uniforms:Vi(n.uniforms),vertexShader:n.vertexShader,fragmentShader:n.fragmentShader,side:Ve,blending:Ln});s.uniforms.tEquirect.value=e;let r=new pe(i,s),a=e.minFilter;if(e.minFilter===Oi)e.minFilter=Qe;return new hh(1,10,this).update(t,r),e.minFilter=a,r.geometry.dispose(),r.material.dispose(),this}clear(t,e=!0,n=!0,i=!0){let s=t.getRenderTarget();for(let r=0;r<6;r++)t.setRenderTarget(this,r),t.clear(e,n,i);t.setRenderTarget(s)}}function _v(t){let e=new WeakMap,n=new WeakMap,i=null;function s(u,f=!1){if(u===null||u===void 0)return null;if(f)return a(u);return r(u)}function r(u){if(u&&u.isTexture){let f=u.mapping;if(f===Ia||f===Pa)if(e.has(u)){let g=e.get(u).texture;return o(g,u.mapping)}else{let g=u.image;if(g&&g.height>0){let _=new Ch(g.height);return _.fromEquirectangularTexture(t,u),e.set(u,_),u.addEventListener("dispose",c),o(_.texture,u.mapping)}else return null}}return u}function a(u){if(u&&u.isTexture){let f=u.mapping,g=f===Ia||f===Pa,_=f===Ts||f===Fi;if(g||_){let m=n.get(u),p=m!==void 0?m.texture.pmremVersion:0;if(u.isRenderTargetTexture&&u.pmremVersion!==p){if(i===null)i=new Ah(t);return m=g?i.fromEquirectangular(u,m):i.fromCubemap(u,m),m.texture.pmremVersion=u.pmremVersion,n.set(u,m),m.texture}else if(m!==void 0)return m.texture;else{let b=u.image;if(g&&b&&b.height>0||_&&b&&l(b)){if(i===null)i=new Ah(t);return m=g?i.fromEquirectangular(u):i.fromCubemap(u),m.texture.pmremVersion=u.pmremVersion,n.set(u,m),u.addEventListener("dispose",h),m.texture}else return null}}}return u}function o(u,f){if(f===Ia)u.mapping=Ts;else if(f===Pa)u.mapping=Fi;return u}function l(u){let f=0,g=6;for(let _=0;_<g;_++)if(u[_]!==void 0)f++;return f===g}function c(u){let f=u.target;f.removeEventListener("dispose",c);let g=e.get(f);if(g!==void 0)e.delete(f),g.dispose()}function h(u){let f=u.target;f.removeEventListener("dispose",h);let g=n.get(f);if(g!==void 0)n.delete(f),g.dispose()}function d(){if(e=new WeakMap,n=new WeakMap,i!==null)i.dispose(),i=null}return{get:s,dispose:d}}function xv(t){let e={};function n(i){if(e[i]!==void 0)return e[i];let s=t.getExtension(i);return e[i]=s,s}return{has:function(i){return n(i)!==null},init:function(){n("EXT_color_buffer_float"),n("WEBGL_clip_cull_distance"),n("OES_texture_float_linear"),n("EXT_color_buffer_half_float"),n("WEBGL_multisampled_render_to_texture"),n("WEBGL_render_shared_exponent")},get:function(i){let s=n(i);if(s===null)qn("WebGLRenderer: "+i+" extension not supported.");return s}}}function yv(t,e,n,i){let s={},r=new WeakMap;function a(d){let u=d.target;if(u.index!==null)e.remove(u.index);for(let g in u.attributes)e.remove(u.attributes[g]);u.removeEventListener("dispose",a),delete s[u.id];let f=r.get(u);if(f)e.remove(f),r.delete(u);if(i.releaseStatesOfGeometry(u),u.isInstancedBufferGeometry===!0)delete u._maxInstanceCount;n.memory.geometries--}function o(d,u){if(s[u.id]===!0)return u;return u.addEventListener("dispose",a),s[u.id]=!0,n.memory.geometries++,u}function l(d){let u=d.attributes;for(let f in u)e.update(u[f],t.ARRAY_BUFFER)}function c(d){let u=[],f=d.index,g=d.attributes.position,_=0;if(g===void 0)return;if(f!==null){let b=f.array;_=f.version;for(let S=0,x=b.length;S<x;S+=3){let T=b[S+0],I=b[S+1],R=b[S+2];u.push(T,I,I,R,R,T)}}else{let b=g.array;_=g.version;for(let S=0,x=b.length/3-1;S<x;S+=3){let T=S+0,I=S+1,R=S+2;u.push(T,I,I,R,R,T)}}let m=new(g.count>=65535?Za:Ya)(u,1);m.version=_;let p=r.get(d);if(p)e.remove(p);r.set(d,m)}function h(d){let u=r.get(d);if(u){let f=d.index;if(f!==null){if(u.version<f.version)c(d)}}else c(d);return r.get(d)}return{get:o,update:l,getWireframeAttribute:h}}function vv(t,e,n){let i;function s(d){i=d}let r,a;function o(d){r=d.type,a=d.bytesPerElement}function l(d,u){t.drawElements(i,u,r,d*a),n.update(u,i,1)}function c(d,u,f){if(f===0)return;t.drawElementsInstanced(i,u,r,d*a,f),n.update(u,i,f)}function h(d,u,f){if(f===0)return;e.get("WEBGL_multi_draw").multiDrawElementsWEBGL(i,u,0,r,d,0,f);let _=0;for(let m=0;m<f;m++)_+=u[m];n.update(_,i,1)}this.setMode=s,this.setIndex=o,this.render=l,this.renderInstances=c,this.renderMultiDraw=h}function Sv(t){let e={geometries:0,textures:0},n={frame:0,calls:0,triangles:0,points:0,lines:0};function i(r,a,o){switch(n.calls++,a){case t.TRIANGLES:n.triangles+=o*(r/3);break;case t.LINES:n.lines+=o*(r/2);break;case t.LINE_STRIP:n.lines+=o*(r-1);break;case t.LINE_LOOP:n.lines+=o*r;break;case t.POINTS:n.points+=o*r;break;default:Wt("WebGLInfo: Unknown draw mode:",a);break}}function s(){n.calls=0,n.triangles=0,n.points=0,n.lines=0}return{memory:e,render:n,programs:null,autoReset:!0,reset:s,update:i}}function Mv(t,e,n){let i=new WeakMap,s=new fe;function r(a,o,l){let c=a.morphTargetInfluences,h=o.morphAttributes.position||o.morphAttributes.normal||o.morphAttributes.color,d=h!==void 0?h.length:0,u=i.get(o);if(u===void 0||u.count!==d){let A=function(){R.dispose(),i.delete(o),o.removeEventListener("dispose",A)};if(u!==void 0)u.texture.dispose();let f=o.morphAttributes.position!==void 0,g=o.morphAttributes.normal!==void 0,_=o.morphAttributes.color!==void 0,m=o.morphAttributes.position||[],p=o.morphAttributes.normal||[],b=o.morphAttributes.color||[],S=0;if(f===!0)S=1;if(g===!0)S=2;if(_===!0)S=3;let x=o.attributes.position.count*S,T=1;if(x>e.maxTextureSize)T=Math.ceil(x/e.maxTextureSize),x=e.maxTextureSize;let I=new Float32Array(x*T*4*d),R=new cr(I,x,T,d);R.type=$n,R.needsUpdate=!0;let y=S*4;for(let Z=0;Z<d;Z++){let D=m[Z],q=p[Z],ot=b[Z],z=x*T*4*Z;for(let et=0;et<D.count;et++){let nt=et*y;if(f===!0)s.fromBufferAttribute(D,et),I[z+nt+0]=s.x,I[z+nt+1]=s.y,I[z+nt+2]=s.z,I[z+nt+3]=0;if(g===!0)s.fromBufferAttribute(q,et),I[z+nt+4]=s.x,I[z+nt+5]=s.y,I[z+nt+6]=s.z,I[z+nt+7]=0;if(_===!0)s.fromBufferAttribute(ot,et),I[z+nt+8]=s.x,I[z+nt+9]=s.y,I[z+nt+10]=s.z,I[z+nt+11]=ot.itemSize===4?s.w:1}}u={count:d,texture:R,size:new rt(x,T)},i.set(o,u),o.addEventListener("dispose",A)}if(a.isInstancedMesh===!0&&a.morphTexture!==null)l.getUniforms().setValue(t,"morphTexture",a.morphTexture,n);else{let f=0;for(let _=0;_<c.length;_++)f+=c[_];let g=o.morphTargetsRelative?1:1-f;l.getUniforms().setValue(t,"morphTargetBaseInfluence",g),l.getUniforms().setValue(t,"morphTargetInfluences",c)}l.getUniforms().setValue(t,"morphTargetsTexture",u.texture,n),l.getUniforms().setValue(t,"morphTargetsTextureSize",u.size)}return{update:r}}function bv(t,e,n,i,s){let r=new WeakMap;function a(c){let h=s.render.frame,d=c.geometry,u=e.get(c,d);if(r.get(u)!==h)e.update(u),r.set(u,h);if(c.isInstancedMesh){if(c.hasEventListener("dispose",l)===!1)c.addEventListener("dispose",l);if(r.get(c)!==h){if(n.update(c.instanceMatrix,t.ARRAY_BUFFER),c.instanceColor!==null)n.update(c.instanceColor,t.ARRAY_BUFFER);r.set(c,h)}}if(c.isSkinnedMesh){let f=c.skeleton;if(r.get(f)!==h)f.update(),r.set(f,h)}return u}function o(){r=new WeakMap}function l(c){let h=c.target;if(h.removeEventListener("dispose",l),i.releaseStatesOfObject(h),n.remove(h.instanceMatrix),h.instanceColor!==null)n.remove(h.instanceColor)}return{update:a,dispose:o}}var Tv={[Pl]:"LINEAR_TONE_MAPPING",[Ll]:"REINHARD_TONE_MAPPING",[Nl]:"CINEON_TONE_MAPPING",[Dl]:"ACES_FILMIC_TONE_MAPPING",[Fl]:"AGX_TONE_MAPPING",[Ol]:"NEUTRAL_TONE_MAPPING",[Ul]:"CUSTOM_TONE_MAPPING"};function Ev(t,e,n,i,s,r){let a=new Je(e,n,{type:t,depthBuffer:s,stencilBuffer:r,samples:i?4:0,storeMultisampledDepthBuffer:!1,storeMultisampledStencilBuffer:!1,resolveDepthBuffer:!1,resolveStencilBuffer:!1}),o=null,l=null,c=new jt;c.setAttribute("position",new Nt([-1,3,0,-1,-1,0,3,-1,0],3)),c.setAttribute("uv",new Nt([0,2,0,0,2,0],2));let h=new vo({uniforms:{tDiffuse:{value:null}},vertexShader:`
			precision highp float;

			uniform mat4 modelViewMatrix;
			uniform mat4 projectionMatrix;

			attribute vec3 position;
			attribute vec2 uv;

			varying vec2 vUv;

			void main() {
				vUv = uv;
				gl_Position = projectionMatrix * modelViewMatrix * vec4( position, 1.0 );
			}`,fragmentShader:`
			precision highp float;

			uniform sampler2D tDiffuse;

			varying vec2 vUv;

			#include <tonemapping_pars_fragment>
			#include <colorspace_pars_fragment>

			void main() {
				gl_FragColor = texture2D( tDiffuse, vUv );

				#ifdef LINEAR_TONE_MAPPING
					gl_FragColor.rgb = LinearToneMapping( gl_FragColor.rgb );
				#elif defined( REINHARD_TONE_MAPPING )
					gl_FragColor.rgb = ReinhardToneMapping( gl_FragColor.rgb );
				#elif defined( CINEON_TONE_MAPPING )
					gl_FragColor.rgb = CineonToneMapping( gl_FragColor.rgb );
				#elif defined( ACES_FILMIC_TONE_MAPPING )
					gl_FragColor.rgb = ACESFilmicToneMapping( gl_FragColor.rgb );
				#elif defined( AGX_TONE_MAPPING )
					gl_FragColor.rgb = AgXToneMapping( gl_FragColor.rgb );
				#elif defined( NEUTRAL_TONE_MAPPING )
					gl_FragColor.rgb = NeutralToneMapping( gl_FragColor.rgb );
				#elif defined( CUSTOM_TONE_MAPPING )
					gl_FragColor.rgb = CustomToneMapping( gl_FragColor.rgb );
				#endif

				#ifdef SRGB_TRANSFER
					gl_FragColor = sRGBTransferOETF( gl_FragColor );
				#endif
			}`,depthTest:!1,depthWrite:!1}),d=new pe(c,h),u=new Ps(-1,1,1,-1,0,1),f=null,g=null,_=!1,m,p=null,b=[],S=!1;this.setSize=function(x,T){if(a.setSize(x,T),o!==null)o.setSize(x,T);if(l!==null)l.setSize(x,T);for(let I=0;I<b.length;I++){let R=b[I];if(R.setSize)R.setSize(x,T)}},this.setEffects=function(x){b=x,S=b.length>0&&b[0].isRenderPass===!0;let{width:T,height:I}=a;if(b.length>0&&o===null)o=new Je(T,I,{type:Nn,depthBuffer:!1,stencilBuffer:!1}),l=new Je(T,I,{type:Nn,depthBuffer:!1,stencilBuffer:!1});for(let R=0;R<b.length;R++){let y=b[R];if(y.setSize)y.setSize(T,I)}},this.begin=function(x,T){if(_)return!1;if(x.toneMapping===bn&&b.length===0)return!1;if(p=T,T!==null){let{width:I,height:R}=T;if(a.width!==I||a.height!==R)this.setSize(I,R)}if(S===!1)x.setRenderTarget(a);return m=x.toneMapping,x.toneMapping=bn,!0},this.hasRenderPass=function(){return S},this.end=function(x,T){x.toneMapping=m,_=!0;let I=a,R=o;for(let y=0;y<b.length;y++){let A=b[y];if(A.enabled===!1)continue;if(A.render(x,R,I,T),A.needsSwap!==!1)I=R,R=R===o?l:o}if(f!==x.outputColorSpace||g!==x.toneMapping){if(f=x.outputColorSpace,g=x.toneMapping,h.defines={},re.getTransfer(f)===xe)h.defines.SRGB_TRANSFER="";let y=Tv[g];if(y)h.defines[y]="";h.needsUpdate=!0}h.uniforms.tDiffuse.value=I.texture,x.setRenderTarget(p),x.render(d,u),p=null,_=!1},this.isCompositing=function(){return _},this.dispose=function(){if(a.dispose(),o!==null)o.dispose();if(l!==null)l.dispose();c.dispose(),h.dispose()}}var Xp=new be,wh=new Hi(1,1),qp=new cr,Yp=new hr,Zp=new Rs,Ap=[],wp=[],Rp=new Float32Array(16),Cp=new Float32Array(9),Ip=new Float32Array(4);function Ds(t,e,n){let i=t[0];if(i<=0||i>0)return t;let s=e*n,r=Ap[s];if(r===void 0)r=new Float32Array(s),Ap[s]=r;if(e!==0){i.toArray(r,0);for(let a=1,o=0;a!==e;++a)o+=n,t[a].toArray(r,o)}return r}function Ie(t,e){if(t.length!==e.length)return!1;for(let n=0,i=t.length;n<i;n++)if(t[n]!==e[n])return!1;return!0}function Pe(t,e){for(let n=0,i=e.length;n<i;n++)t[n]=e[n]}function Uo(t,e){let n=wp[e];if(n===void 0)n=new Int32Array(e),wp[e]=n;for(let i=0;i!==e;++i)n[i]=t.allocateTextureUnit();return n}function Av(t,e){let n=this.cache;if(n[0]===e)return;t.uniform1f(this.addr,e),n[0]=e}function wv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y)t.uniform2f(this.addr,e.x,e.y),n[0]=e.x,n[1]=e.y}else{if(Ie(n,e))return;t.uniform2fv(this.addr,e),Pe(n,e)}}function Rv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y||n[2]!==e.z)t.uniform3f(this.addr,e.x,e.y,e.z),n[0]=e.x,n[1]=e.y,n[2]=e.z}else if(e.r!==void 0){if(n[0]!==e.r||n[1]!==e.g||n[2]!==e.b)t.uniform3f(this.addr,e.r,e.g,e.b),n[0]=e.r,n[1]=e.g,n[2]=e.b}else{if(Ie(n,e))return;t.uniform3fv(this.addr,e),Pe(n,e)}}function Cv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y||n[2]!==e.z||n[3]!==e.w)t.uniform4f(this.addr,e.x,e.y,e.z,e.w),n[0]=e.x,n[1]=e.y,n[2]=e.z,n[3]=e.w}else{if(Ie(n,e))return;t.uniform4fv(this.addr,e),Pe(n,e)}}function Iv(t,e){let n=this.cache,i=e.elements;if(i===void 0){if(Ie(n,e))return;t.uniformMatrix2fv(this.addr,!1,e),Pe(n,e)}else{if(Ie(n,i))return;Ip.set(i),t.uniformMatrix2fv(this.addr,!1,Ip),Pe(n,i)}}function Pv(t,e){let n=this.cache,i=e.elements;if(i===void 0){if(Ie(n,e))return;t.uniformMatrix3fv(this.addr,!1,e),Pe(n,e)}else{if(Ie(n,i))return;Cp.set(i),t.uniformMatrix3fv(this.addr,!1,Cp),Pe(n,i)}}function Lv(t,e){let n=this.cache,i=e.elements;if(i===void 0){if(Ie(n,e))return;t.uniformMatrix4fv(this.addr,!1,e),Pe(n,e)}else{if(Ie(n,i))return;Rp.set(i),t.uniformMatrix4fv(this.addr,!1,Rp),Pe(n,i)}}function Nv(t,e){let n=this.cache;if(n[0]===e)return;t.uniform1i(this.addr,e),n[0]=e}function Dv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y)t.uniform2i(this.addr,e.x,e.y),n[0]=e.x,n[1]=e.y}else{if(Ie(n,e))return;t.uniform2iv(this.addr,e),Pe(n,e)}}function Uv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y||n[2]!==e.z)t.uniform3i(this.addr,e.x,e.y,e.z),n[0]=e.x,n[1]=e.y,n[2]=e.z}else{if(Ie(n,e))return;t.uniform3iv(this.addr,e),Pe(n,e)}}function Fv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y||n[2]!==e.z||n[3]!==e.w)t.uniform4i(this.addr,e.x,e.y,e.z,e.w),n[0]=e.x,n[1]=e.y,n[2]=e.z,n[3]=e.w}else{if(Ie(n,e))return;t.uniform4iv(this.addr,e),Pe(n,e)}}function Ov(t,e){let n=this.cache;if(n[0]===e)return;t.uniform1ui(this.addr,e),n[0]=e}function Bv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y)t.uniform2ui(this.addr,e.x,e.y),n[0]=e.x,n[1]=e.y}else{if(Ie(n,e))return;t.uniform2uiv(this.addr,e),Pe(n,e)}}function zv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y||n[2]!==e.z)t.uniform3ui(this.addr,e.x,e.y,e.z),n[0]=e.x,n[1]=e.y,n[2]=e.z}else{if(Ie(n,e))return;t.uniform3uiv(this.addr,e),Pe(n,e)}}function kv(t,e){let n=this.cache;if(e.x!==void 0){if(n[0]!==e.x||n[1]!==e.y||n[2]!==e.z||n[3]!==e.w)t.uniform4ui(this.addr,e.x,e.y,e.z,e.w),n[0]=e.x,n[1]=e.y,n[2]=e.z,n[3]=e.w}else{if(Ie(n,e))return;t.uniform4uiv(this.addr,e),Pe(n,e)}}function Gv(t,e,n){let i=this.cache,s=n.allocateTextureUnit();if(i[0]!==s)t.uniform1i(this.addr,s),i[0]=s;let r;if(this.type===t.SAMPLER_2D_SHADOW)wh.compareFunction=n.isReversedDepthBuffer()?Va:Ha,r=wh;else r=Xp;n.setTexture2D(e||r,s)}function Hv(t,e,n){let i=this.cache,s=n.allocateTextureUnit();if(i[0]!==s)t.uniform1i(this.addr,s),i[0]=s;n.setTexture3D(e||Yp,s)}function Vv(t,e,n){let i=this.cache,s=n.allocateTextureUnit();if(i[0]!==s)t.uniform1i(this.addr,s),i[0]=s;n.setTextureCube(e||Zp,s)}function Wv(t,e,n){let i=this.cache,s=n.allocateTextureUnit();if(i[0]!==s)t.uniform1i(this.addr,s),i[0]=s;n.setTexture2DArray(e||qp,s)}function Xv(t){switch(t){case 5126:return Av;case 35664:return wv;case 35665:return Rv;case 35666:return Cv;case 35674:return Iv;case 35675:return Pv;case 35676:return Lv;case 5124:case 35670:return Nv;case 35667:case 35671:return Dv;case 35668:case 35672:return Uv;case 35669:case 35673:return Fv;case 5125:return Ov;case 36294:return Bv;case 36295:return zv;case 36296:return kv;case 35678:case 36198:case 36298:case 36306:case 35682:return Gv;case 35679:case 36299:case 36307:return Hv;case 35680:case 36300:case 36308:case 36293:return Vv;case 36289:case 36303:case 36311:case 36292:return Wv}}function qv(t,e){t.uniform1fv(this.addr,e)}function Yv(t,e){let n=Ds(e,this.size,2);t.uniform2fv(this.addr,n)}function Zv(t,e){let n=Ds(e,this.size,3);t.uniform3fv(this.addr,n)}function Jv(t,e){let n=Ds(e,this.size,4);t.uniform4fv(this.addr,n)}function $v(t,e){let n=Ds(e,this.size,4);t.uniformMatrix2fv(this.addr,!1,n)}function Kv(t,e){let n=Ds(e,this.size,9);t.uniformMatrix3fv(this.addr,!1,n)}function jv(t,e){let n=Ds(e,this.size,16);t.uniformMatrix4fv(this.addr,!1,n)}function Qv(t,e){t.uniform1iv(this.addr,e)}function tS(t,e){t.uniform2iv(this.addr,e)}function eS(t,e){t.uniform3iv(this.addr,e)}function nS(t,e){t.uniform4iv(this.addr,e)}function iS(t,e){t.uniform1uiv(this.addr,e)}function sS(t,e){t.uniform2uiv(this.addr,e)}function rS(t,e){t.uniform3uiv(this.addr,e)}function aS(t,e){t.uniform4uiv(this.addr,e)}function oS(t,e,n){let i=this.cache,s=e.length,r=Uo(n,s);if(!Ie(i,r))t.uniform1iv(this.addr,r),Pe(i,r);let a;if(this.type===t.SAMPLER_2D_SHADOW)a=wh;else a=Xp;for(let o=0;o!==s;++o)n.setTexture2D(e[o]||a,r[o])}function lS(t,e,n){let i=this.cache,s=e.length,r=Uo(n,s);if(!Ie(i,r))t.uniform1iv(this.addr,r),Pe(i,r);for(let a=0;a!==s;++a)n.setTexture3D(e[a]||Yp,r[a])}function cS(t,e,n){let i=this.cache,s=e.length,r=Uo(n,s);if(!Ie(i,r))t.uniform1iv(this.addr,r),Pe(i,r);for(let a=0;a!==s;++a)n.setTextureCube(e[a]||Zp,r[a])}function hS(t,e,n){let i=this.cache,s=e.length,r=Uo(n,s);if(!Ie(i,r))t.uniform1iv(this.addr,r),Pe(i,r);for(let a=0;a!==s;++a)n.setTexture2DArray(e[a]||qp,r[a])}function uS(t){switch(t){case 5126:return qv;case 35664:return Yv;case 35665:return Zv;case 35666:return Jv;case 35674:return $v;case 35675:return Kv;case 35676:return jv;case 5124:case 35670:return Qv;case 35667:case 35671:return tS;case 35668:case 35672:return eS;case 35669:case 35673:return nS;case 5125:return iS;case 36294:return sS;case 36295:return rS;case 36296:return aS;case 35678:case 36198:case 36298:case 36306:case 35682:return oS;case 35679:case 36299:case 36307:return lS;case 35680:case 36300:case 36308:case 36293:return cS;case 36289:case 36303:case 36311:case 36292:return hS}}class Jp{constructor(t,e,n){this.id=t,this.addr=n,this.cache=[],this.type=e.type,this.setValue=Xv(e.type)}}class $p{constructor(t,e,n){this.id=t,this.addr=n,this.cache=[],this.type=e.type,this.size=e.size,this.setValue=uS(e.type)}}class Kp{constructor(t){this.id=t,this.seq=[],this.map={}}setValue(t,e,n){let i=this.seq;for(let s=0,r=i.length;s!==r;++s){let a=i[s];a.setValue(t,e[a.id],n)}}}var Th=/(\w+)(\])?(\[|\.)?/g;function Pp(t,e){t.seq.push(e),t.map[e.id]=e}function dS(t,e,n){let i=t.name,s=i.length;Th.lastIndex=0;while(!0){let r=Th.exec(i),a=Th.lastIndex,o=r[1],l=r[2]==="]",c=r[3];if(l)o=o|0;if(c===void 0||c==="["&&a+2===s){Pp(n,c===void 0?new Jp(o,t,e):new $p(o,t,e));break}else{let d=n.map[o];if(d===void 0)d=new Kp(o),Pp(n,d);n=d}}}class Lr{constructor(t,e){this.seq=[],this.map={};let n=t.getProgramParameter(e,t.ACTIVE_UNIFORMS);for(let r=0;r<n;++r){let a=t.getActiveUniform(e,r),o=t.getUniformLocation(e,a.name);dS(a,o,this)}let i=[],s=[];for(let r of this.seq)if(r.type===t.SAMPLER_2D_SHADOW||r.type===t.SAMPLER_CUBE_SHADOW||r.type===t.SAMPLER_2D_ARRAY_SHADOW)i.push(r);else s.push(r);if(i.length>0)this.seq=i.concat(s)}setValue(t,e,n,i){let s=this.map[e];if(s!==void 0)s.setValue(t,n,i)}setOptional(t,e,n){let i=e[n];if(i!==void 0)this.setValue(t,n,i)}static upload(t,e,n,i){for(let s=0,r=e.length;s!==r;++s){let a=e[s],o=n[a.id];if(o.needsUpdate!==!1)a.setValue(t,o.value,i)}}static seqWithValue(t,e){let n=[];for(let i=0,s=t.length;i!==s;++i){let r=t[i];if(r.id in e)n.push(r)}return n}}function Lp(t,e,n){let i=t.createShader(e);return t.shaderSource(i,n),t.compileShader(i),i}var fS=37297,pS=0;function mS(t,e){let n=t.split(`
`),i=[],s=Math.max(e-6,0),r=Math.min(e+6,n.length);for(let a=s;a<r;a++){let o=a+1;i.push(`${o===e?">":" "} ${o}: ${n[a]}`)}return i.join(`
`)}var Np=new Xt;function gS(t){re._getMatrix(Np,re.workingColorSpace,t);let e=`mat3( ${Np.elements.map((n)=>n.toFixed(4))} )`;switch(re.getTransfer(t)){case Mc:return[e,"LinearTransferOETF"];case xe:return[e,"sRGBTransferOETF"];default:return Mt("WebGLProgram: Unsupported color space: ",t),[e,"LinearTransferOETF"]}}function Dp(t,e,n){let i=t.getShaderParameter(e,t.COMPILE_STATUS),r=(t.getShaderInfoLog(e)||"").trim();if(i&&r==="")return"";let a=/ERROR: 0:(\d+)/.exec(r);if(a){let o=parseInt(a[1]);return n.toUpperCase()+`

`+r+`

`+mS(t.getShaderSource(e),o)}else return r}function _S(t,e){let n=gS(e);return[`vec4 ${t}( vec4 value ) {`,`	return ${n[1]}( vec4( value.rgb * ${n[0]}, value.a ) );`,"}"].join(`
`)}var xS={[Pl]:"Linear",[Ll]:"Reinhard",[Nl]:"Cineon",[Dl]:"ACESFilmic",[Fl]:"AgX",[Ol]:"Neutral",[Ul]:"Custom"};function yS(t,e){let n=xS[e];if(n===void 0)return Mt("WebGLProgram: Unsupported toneMapping:",e),"vec3 "+t+"( vec3 color ) { return LinearToneMapping( color ); }";return"vec3 "+t+"( vec3 color ) { return "+n+"ToneMapping( color ); }"}var No=new L;function vS(){re.getLuminanceCoefficients(No);let t=No.x.toFixed(4),e=No.y.toFixed(4),n=No.z.toFixed(4);return["float luminance( const in vec3 rgb ) {",`	const vec3 weights = vec3( ${t}, ${e}, ${n} );`,"\treturn dot( weights, rgb );","}"].join(`
`)}function SS(t){return[t.extensionClipCullDistance?"#extension GL_ANGLE_clip_cull_distance : require":"",t.extensionMultiDraw?"#extension GL_ANGLE_multi_draw : require":""].filter(Pr).join(`
`)}function MS(t){let e=[];for(let n in t){let i=t[n];if(i===!1)continue;e.push("#define "+n+" "+i)}return e.join(`
`)}function bS(t,e){let n={},i=t.getProgramParameter(e,t.ACTIVE_ATTRIBUTES);for(let s=0;s<i;s++){let r=t.getActiveAttrib(e,s),a=r.name,o=1;if(r.type===t.FLOAT_MAT2)o=2;if(r.type===t.FLOAT_MAT3)o=3;if(r.type===t.FLOAT_MAT4)o=4;n[a]={type:r.type,location:t.getAttribLocation(e,a),locationSize:o}}return n}function Pr(t){return t!==""}function Up(t,e){let n=e.numSpotLightShadows+e.numSpotLightMaps-e.numSpotLightShadowsWithMaps;return t.replace(/NUM_SUN_LIGHTS/g,e.numSunLights).replace(/NUM_DIR_LIGHTS/g,e.numDirLights).replace(/NUM_SPOT_LIGHTS/g,e.numSpotLights).replace(/NUM_SPOT_LIGHT_MAPS/g,e.numSpotLightMaps).replace(/NUM_SPOT_LIGHT_COORDS/g,n).replace(/NUM_RECT_AREA_LIGHTS/g,e.numRectAreaLights).replace(/NUM_POINT_LIGHTS/g,e.numPointLights).replace(/NUM_HEMI_LIGHTS/g,e.numHemiLights).replace(/NUM_SUN_LIGHT_SHADOWS/g,e.numSunLightShadows).replace(/NUM_DIR_LIGHT_SHADOWS/g,e.numDirLightShadows).replace(/NUM_SPOT_LIGHT_SHADOWS_WITH_MAPS/g,e.numSpotLightShadowsWithMaps).replace(/NUM_SPOT_LIGHT_SHADOWS/g,e.numSpotLightShadows).replace(/NUM_POINT_LIGHT_SHADOWS/g,e.numPointLightShadows)}function Fp(t,e){return t.replace(/NUM_CLIPPING_PLANES/g,e.numClippingPlanes).replace(/UNION_CLIPPING_PLANES/g,e.numClippingPlanes-e.numClipIntersection)}var TS=/^[ \t]*#include +<([\w\d./]+)>/gm;function Rh(t){return t.replace(TS,AS)}var ES=new Map;function AS(t,e){let n=ie[e];if(n===void 0){let i=ES.get(e);if(i!==void 0)n=ie[i],Mt('WebGLRenderer: Shader chunk "%s" has been deprecated. Use "%s" instead.',e,i);else throw Error("THREE.WebGLProgram: Can not resolve #include <"+e+">")}return Rh(n)}var wS=/#pragma unroll_loop_start\s+for\s*\(\s*int\s+i\s*=\s*(\d+)\s*;\s*i\s*<\s*(\d+)\s*;\s*i\s*\+\+\s*\)\s*{([\s\S]+?)}\s+#pragma unroll_loop_end/g;function Op(t){return t.replace(wS,RS)}function RS(t,e,n,i){let s="";for(let r=parseInt(e);r<parseInt(n);r++)s+=i.replace(/\[\s*i\s*\]/g,"[ "+r+" ]").replace(/UNROLLED_LOOP_INDEX/g,r);return s}function Bp(t){let e=`precision ${t.precision} float;
	precision ${t.precision} int;
	precision ${t.precision} sampler2D;
	precision ${t.precision} samplerCube;
	precision ${t.precision} sampler3D;
	precision ${t.precision} sampler2DArray;
	precision ${t.precision} sampler2DShadow;
	precision ${t.precision} samplerCubeShadow;
	precision ${t.precision} sampler2DArrayShadow;
	precision ${t.precision} isampler2D;
	precision ${t.precision} isampler3D;
	precision ${t.precision} isamplerCube;
	precision ${t.precision} isampler2DArray;
	precision ${t.precision} usampler2D;
	precision ${t.precision} usampler3D;
	precision ${t.precision} usamplerCube;
	precision ${t.precision} usampler2DArray;
	`;if(t.precision==="highp")e+=`
#define HIGH_PRECISION`;else if(t.precision==="mediump")e+=`
#define MEDIUM_PRECISION`;else if(t.precision==="lowp")e+=`
#define LOW_PRECISION`;return e}var CS={[sr]:"SHADOWMAP_TYPE_PCF",[Ss]:"SHADOWMAP_TYPE_VSM"};function IS(t){return CS[t.shadowMapType]||"SHADOWMAP_TYPE_BASIC"}var PS={[Ts]:"ENVMAP_TYPE_CUBE",[Fi]:"ENVMAP_TYPE_CUBE",[ar]:"ENVMAP_TYPE_CUBE_UV"};function LS(t){if(t.envMap===!1)return"ENVMAP_TYPE_CUBE";return PS[t.envMapMode]||"ENVMAP_TYPE_CUBE"}var NS={[Fi]:"ENVMAP_MODE_REFRACTION"};function DS(t){if(t.envMap===!1)return"ENVMAP_MODE_REFLECTION";return NS[t.envMapMode]||"ENVMAP_MODE_REFLECTION"}var US={[Ld]:"ENVMAP_BLENDING_MULTIPLY",[Nd]:"ENVMAP_BLENDING_MIX",[Dd]:"ENVMAP_BLENDING_ADD"};function FS(t){if(t.envMap===!1)return"ENVMAP_BLENDING_NONE";return US[t.combine]||"ENVMAP_BLENDING_NONE"}function OS(t){let e=t.envMapCubeUVHeight;if(e===null)return null;let n=Math.log2(e)-2,i=1/e;return{texelWidth:1/(3*Math.max(Math.pow(2,n),112)),texelHeight:i,maxMip:n}}function BS(t,e,n,i){let s=t.getContext(),{defines:r,vertexShader:a,fragmentShader:o}=n,l=IS(n),c=LS(n),h=DS(n),d=FS(n),u=OS(n),f=SS(n),g=MS(r),_=s.createProgram(),m,p,b=n.glslVersion?"#version "+n.glslVersion+`
`:"";if(n.isRawShaderMaterial){if(m=["#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,g].filter(Pr).join(`
`),m.length>0)m+=`
`;if(p=["#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,g].filter(Pr).join(`
`),p.length>0)p+=`
`}else m=[Bp(n),"#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,g,n.extensionClipCullDistance?"#define USE_CLIP_DISTANCE":"",n.batching?"#define USE_BATCHING":"",n.batchingColor?"#define USE_BATCHING_COLOR":"",n.instancing?"#define USE_INSTANCING":"",n.instancingColor?"#define USE_INSTANCING_COLOR":"",n.instancingMorph?"#define USE_INSTANCING_MORPH":"",n.useFog&&n.fog?"#define USE_FOG":"",n.useFog&&n.fogExp2?"#define FOG_EXP2":"",n.map?"#define USE_MAP":"",n.envMap?"#define USE_ENVMAP":"",n.envMap?"#define "+h:"",n.lightMap?"#define USE_LIGHTMAP":"",n.aoMap?"#define USE_AOMAP":"",n.bumpMap?"#define USE_BUMPMAP":"",n.normalMap?"#define USE_NORMALMAP":"",n.normalMapObjectSpace?"#define USE_NORMALMAP_OBJECTSPACE":"",n.normalMapTangentSpace?"#define USE_NORMALMAP_TANGENTSPACE":"",n.displacementMap?"#define USE_DISPLACEMENTMAP":"",n.emissiveMap?"#define USE_EMISSIVEMAP":"",n.anisotropy?"#define USE_ANISOTROPY":"",n.anisotropyMap?"#define USE_ANISOTROPYMAP":"",n.clearcoatMap?"#define USE_CLEARCOATMAP":"",n.clearcoatRoughnessMap?"#define USE_CLEARCOAT_ROUGHNESSMAP":"",n.clearcoatNormalMap?"#define USE_CLEARCOAT_NORMALMAP":"",n.iridescenceMap?"#define USE_IRIDESCENCEMAP":"",n.iridescenceThicknessMap?"#define USE_IRIDESCENCE_THICKNESSMAP":"",n.specularMap?"#define USE_SPECULARMAP":"",n.specularColorMap?"#define USE_SPECULAR_COLORMAP":"",n.specularIntensityMap?"#define USE_SPECULAR_INTENSITYMAP":"",n.roughnessMap?"#define USE_ROUGHNESSMAP":"",n.metalnessMap?"#define USE_METALNESSMAP":"",n.alphaMap?"#define USE_ALPHAMAP":"",n.alphaHash?"#define USE_ALPHAHASH":"",n.transmission?"#define USE_TRANSMISSION":"",n.transmissionMap?"#define USE_TRANSMISSIONMAP":"",n.thicknessMap?"#define USE_THICKNESSMAP":"",n.sheenColorMap?"#define USE_SHEEN_COLORMAP":"",n.sheenRoughnessMap?"#define USE_SHEEN_ROUGHNESSMAP":"",n.mapUv?"#define MAP_UV "+n.mapUv:"",n.alphaMapUv?"#define ALPHAMAP_UV "+n.alphaMapUv:"",n.lightMapUv?"#define LIGHTMAP_UV "+n.lightMapUv:"",n.aoMapUv?"#define AOMAP_UV "+n.aoMapUv:"",n.emissiveMapUv?"#define EMISSIVEMAP_UV "+n.emissiveMapUv:"",n.bumpMapUv?"#define BUMPMAP_UV "+n.bumpMapUv:"",n.normalMapUv?"#define NORMALMAP_UV "+n.normalMapUv:"",n.displacementMapUv?"#define DISPLACEMENTMAP_UV "+n.displacementMapUv:"",n.metalnessMapUv?"#define METALNESSMAP_UV "+n.metalnessMapUv:"",n.roughnessMapUv?"#define ROUGHNESSMAP_UV "+n.roughnessMapUv:"",n.anisotropyMapUv?"#define ANISOTROPYMAP_UV "+n.anisotropyMapUv:"",n.clearcoatMapUv?"#define CLEARCOATMAP_UV "+n.clearcoatMapUv:"",n.clearcoatNormalMapUv?"#define CLEARCOAT_NORMALMAP_UV "+n.clearcoatNormalMapUv:"",n.clearcoatRoughnessMapUv?"#define CLEARCOAT_ROUGHNESSMAP_UV "+n.clearcoatRoughnessMapUv:"",n.iridescenceMapUv?"#define IRIDESCENCEMAP_UV "+n.iridescenceMapUv:"",n.iridescenceThicknessMapUv?"#define IRIDESCENCE_THICKNESSMAP_UV "+n.iridescenceThicknessMapUv:"",n.sheenColorMapUv?"#define SHEEN_COLORMAP_UV "+n.sheenColorMapUv:"",n.sheenRoughnessMapUv?"#define SHEEN_ROUGHNESSMAP_UV "+n.sheenRoughnessMapUv:"",n.specularMapUv?"#define SPECULARMAP_UV "+n.specularMapUv:"",n.specularColorMapUv?"#define SPECULAR_COLORMAP_UV "+n.specularColorMapUv:"",n.specularIntensityMapUv?"#define SPECULAR_INTENSITYMAP_UV "+n.specularIntensityMapUv:"",n.transmissionMapUv?"#define TRANSMISSIONMAP_UV "+n.transmissionMapUv:"",n.thicknessMapUv?"#define THICKNESSMAP_UV "+n.thicknessMapUv:"",n.vertexTangents&&n.flatShading===!1?"#define USE_TANGENT":"",n.vertexNormals?"#define HAS_NORMAL":"",n.vertexColors?"#define USE_COLOR":"",n.vertexAlphas?"#define USE_COLOR_ALPHA":"",n.vertexUv1s?"#define USE_UV1":"",n.vertexUv2s?"#define USE_UV2":"",n.vertexUv3s?"#define USE_UV3":"",n.pointsUvs?"#define USE_POINTS_UV":"",n.flatShading?"#define FLAT_SHADED":"",n.skinning?"#define USE_SKINNING":"",n.morphTargets?"#define USE_MORPHTARGETS":"",n.morphNormals&&n.flatShading===!1?"#define USE_MORPHNORMALS":"",n.morphColors?"#define USE_MORPHCOLORS":"",n.morphTargetsCount>0?"#define MORPHTARGETS_TEXTURE_STRIDE "+n.morphTextureStride:"",n.morphTargetsCount>0?"#define MORPHTARGETS_COUNT "+n.morphTargetsCount:"",n.doubleSided?"#define DOUBLE_SIDED":"",n.flipSided?"#define FLIP_SIDED":"",n.shadowMapEnabled?"#define USE_SHADOWMAP":"",n.shadowMapEnabled?"#define "+l:"",n.sizeAttenuation?"#define USE_SIZEATTENUATION":"",n.numLightProbes>0?"#define USE_LIGHT_PROBES":"",n.logarithmicDepthBuffer?"#define USE_LOGARITHMIC_DEPTH_BUFFER":"",n.reversedDepthBuffer?"#define USE_REVERSED_DEPTH_BUFFER":"","uniform mat4 modelMatrix;","uniform mat4 modelViewMatrix;","uniform mat4 projectionMatrix;","uniform mat4 viewMatrix;","uniform mat3 normalMatrix;","uniform vec3 cameraPosition;","uniform bool isOrthographic;","#ifdef USE_INSTANCING","\tattribute mat4 instanceMatrix;","#endif","#ifdef USE_INSTANCING_COLOR","\tattribute vec3 instanceColor;","#endif","#ifdef USE_INSTANCING_MORPH","\tuniform sampler2D morphTexture;","#endif","attribute vec3 position;","attribute vec3 normal;","attribute vec2 uv;","#ifdef USE_UV1","\tattribute vec2 uv1;","#endif","#ifdef USE_UV2","\tattribute vec2 uv2;","#endif","#ifdef USE_UV3","\tattribute vec2 uv3;","#endif","#ifdef USE_TANGENT","\tattribute vec4 tangent;","#endif","#if defined( USE_COLOR_ALPHA )","\tattribute vec4 color;","#elif defined( USE_COLOR )","\tattribute vec3 color;","#endif","#ifdef USE_SKINNING","\tattribute vec4 skinIndex;","\tattribute vec4 skinWeight;","#endif",`
`].filter(Pr).join(`
`),p=[Bp(n),"#define SHADER_TYPE "+n.shaderType,"#define SHADER_NAME "+n.shaderName,g,n.useFog&&n.fog?"#define USE_FOG":"",n.useFog&&n.fogExp2?"#define FOG_EXP2":"",n.alphaToCoverage?"#define ALPHA_TO_COVERAGE":"",n.map?"#define USE_MAP":"",n.matcap?"#define USE_MATCAP":"",n.envMap?"#define USE_ENVMAP":"",n.envMap?"#define "+c:"",n.envMap?"#define "+h:"",n.envMap?"#define "+d:"",u?"#define CUBEUV_TEXEL_WIDTH "+u.texelWidth:"",u?"#define CUBEUV_TEXEL_HEIGHT "+u.texelHeight:"",u?"#define CUBEUV_MAX_MIP "+u.maxMip+".0":"",n.lightMap?"#define USE_LIGHTMAP":"",n.aoMap?"#define USE_AOMAP":"",n.bumpMap?"#define USE_BUMPMAP":"",n.normalMap?"#define USE_NORMALMAP":"",n.normalMapObjectSpace?"#define USE_NORMALMAP_OBJECTSPACE":"",n.normalMapTangentSpace?"#define USE_NORMALMAP_TANGENTSPACE":"",n.packedNormalMap?"#define USE_PACKED_NORMALMAP":"",n.emissiveMap?"#define USE_EMISSIVEMAP":"",n.anisotropy?"#define USE_ANISOTROPY":"",n.anisotropyMap?"#define USE_ANISOTROPYMAP":"",n.clearcoat?"#define USE_CLEARCOAT":"",n.clearcoatMap?"#define USE_CLEARCOATMAP":"",n.clearcoatRoughnessMap?"#define USE_CLEARCOAT_ROUGHNESSMAP":"",n.clearcoatNormalMap?"#define USE_CLEARCOAT_NORMALMAP":"",n.dispersion?"#define USE_DISPERSION":"",n.retroreflection?"#define USE_RETROREFLECTION":"",n.iridescence?"#define USE_IRIDESCENCE":"",n.iridescenceMap?"#define USE_IRIDESCENCEMAP":"",n.iridescenceThicknessMap?"#define USE_IRIDESCENCE_THICKNESSMAP":"",n.specularMap?"#define USE_SPECULARMAP":"",n.specularColorMap?"#define USE_SPECULAR_COLORMAP":"",n.specularIntensityMap?"#define USE_SPECULAR_INTENSITYMAP":"",n.roughnessMap?"#define USE_ROUGHNESSMAP":"",n.metalnessMap?"#define USE_METALNESSMAP":"",n.alphaMap?"#define USE_ALPHAMAP":"",n.alphaTest?"#define USE_ALPHATEST":"",n.alphaHash?"#define USE_ALPHAHASH":"",n.sheen?"#define USE_SHEEN":"",n.sheenColorMap?"#define USE_SHEEN_COLORMAP":"",n.sheenRoughnessMap?"#define USE_SHEEN_ROUGHNESSMAP":"",n.transmission?"#define USE_TRANSMISSION":"",n.transmissionMap?"#define USE_TRANSMISSIONMAP":"",n.thicknessMap?"#define USE_THICKNESSMAP":"",n.vertexTangents&&n.flatShading===!1?"#define USE_TANGENT":"",n.vertexColors||n.instancingColor?"#define USE_COLOR":"",n.vertexAlphas||n.batchingColor?"#define USE_COLOR_ALPHA":"",n.vertexUv1s?"#define USE_UV1":"",n.vertexUv2s?"#define USE_UV2":"",n.vertexUv3s?"#define USE_UV3":"",n.pointsUvs?"#define USE_POINTS_UV":"",n.gradientMap?"#define USE_GRADIENTMAP":"",n.flatShading?"#define FLAT_SHADED":"",n.doubleSided?"#define DOUBLE_SIDED":"",n.flipSided?"#define FLIP_SIDED":"",n.shadowMapEnabled?"#define USE_SHADOWMAP":"",n.shadowMapEnabled?"#define "+l:"",n.premultipliedAlpha?"#define PREMULTIPLIED_ALPHA":"",n.numLightProbes>0?"#define USE_LIGHT_PROBES":"",n.numLightProbeGrids>0?"#define USE_LIGHT_PROBES_GRID":"",n.decodeVideoTexture?"#define DECODE_VIDEO_TEXTURE":"",n.decodeVideoTextureEmissive?"#define DECODE_VIDEO_TEXTURE_EMISSIVE":"",n.logarithmicDepthBuffer?"#define USE_LOGARITHMIC_DEPTH_BUFFER":"",n.reversedDepthBuffer?"#define USE_REVERSED_DEPTH_BUFFER":"","uniform mat4 viewMatrix;","uniform vec3 cameraPosition;","uniform bool isOrthographic;",n.toneMapping!==bn?"#define TONE_MAPPING":"",n.toneMapping!==bn?ie.tonemapping_pars_fragment:"",n.toneMapping!==bn?yS("toneMapping",n.toneMapping):"",n.dithering?"#define DITHERING":"",n.opaque?"#define OPAQUE":"",ie.colorspace_pars_fragment,_S("linearToOutputTexel",n.outputColorSpace),vS(),n.useDepthPacking?"#define DEPTH_PACKING "+n.depthPacking:"",`
`].filter(Pr).join(`
`);if(a=Rh(a),a=Up(a,n),a=Fp(a,n),o=Rh(o),o=Up(o,n),o=Fp(o,n),a=Op(a),o=Op(o),n.isRawShaderMaterial!==!0)b=`#version 300 es
`,m=[f,"#define attribute in","#define varying out","#define texture2D texture"].join(`
`)+`
`+m,p=["#define varying in",n.glslVersion===bc?"":"layout(location = 0) out highp vec4 pc_fragColor;",n.glslVersion===bc?"":"#define gl_FragColor pc_fragColor","#define gl_FragDepthEXT gl_FragDepth","#define texture2D texture","#define textureCube texture","#define texture2DProj textureProj","#define texture2DLodEXT textureLod","#define texture2DProjLodEXT textureProjLod","#define textureCubeLodEXT textureLod","#define texture2DGradEXT textureGrad","#define texture2DProjGradEXT textureProjGrad","#define textureCubeGradEXT textureGrad"].join(`
`)+`
`+p;let S=b+m+a,x=b+p+o,T=Lp(s,s.VERTEX_SHADER,S),I=Lp(s,s.FRAGMENT_SHADER,x);if(s.attachShader(_,T),s.attachShader(_,I),n.index0AttributeName!==void 0)s.bindAttribLocation(_,0,n.index0AttributeName);else if(n.hasPositionAttribute===!0)s.bindAttribLocation(_,0,"position");s.linkProgram(_);function R(D){if(t.debug.checkShaderErrors){let q=s.getProgramInfoLog(_)||"",ot=s.getShaderInfoLog(T)||"",z=s.getShaderInfoLog(I)||"",et=q.trim(),nt=ot.trim(),B=z.trim(),ft=!0,Q=!0;if(s.getProgramParameter(_,s.LINK_STATUS)===!1)if(ft=!1,typeof t.debug.onShaderError==="function")t.debug.onShaderError(s,_,T,I);else{let ct=Dp(s,T,"vertex"),dt=Dp(s,I,"fragment");Wt("WebGLProgram: Shader Error "+s.getError()+" - VALIDATE_STATUS "+s.getProgramParameter(_,s.VALIDATE_STATUS)+`

Material Name: `+D.name+`
Material Type: `+D.type+`

Program Info Log: `+et+`
`+ct+`
`+dt)}else if(et!=="")Mt("WebGLProgram: Program Info Log:",et);else if(nt===""||B==="")Q=!1;if(Q)D.diagnostics={runnable:ft,programLog:et,vertexShader:{log:nt,prefix:m},fragmentShader:{log:B,prefix:p}}}s.deleteShader(T),s.deleteShader(I),y=new Lr(s,_),A=bS(s,_)}let y;this.getUniforms=function(){if(y===void 0)R(this);return y};let A;this.getAttributes=function(){if(A===void 0)R(this);return A};let Z=n.rendererExtensionParallelShaderCompile===!1;return this.isReady=function(){if(Z===!1)Z=s.getProgramParameter(_,fS);return Z},this.destroy=function(){i.releaseStatesOfProgram(this),s.deleteProgram(_),this.program=void 0},this.type=n.shaderType,this.name=n.shaderName,this.id=pS++,this.cacheKey=e,this.usedTimes=1,this.program=_,this.vertexShader=T,this.fragmentShader=I,this}var zS=0;class jp{constructor(){this.shaderCache=new Map,this.materialCache=new Map}update(t,e,n){let i=this._getShaderCacheForMaterial(t);if(i.has(e)===!1)i.add(e),e.usedTimes++;if(i.has(n)===!1)i.add(n),n.usedTimes++;return this}remove(t){let e=this.materialCache.get(t);for(let n of e)if(n.usedTimes--,n.usedTimes===0)this.shaderCache.delete(n.code);return this.materialCache.delete(t),this}getVertexShaderStage(t){return this._getShaderStage(t.vertexShader)}getFragmentShaderStage(t){return this._getShaderStage(t.fragmentShader)}dispose(){this.shaderCache.clear(),this.materialCache.clear()}_getShaderCacheForMaterial(t){let e=this.materialCache,n=e.get(t);if(n===void 0)n=new Set,e.set(t,n);return n}_getShaderStage(t){let e=this.shaderCache,n=e.get(t);if(n===void 0)n=new Qp(t),e.set(t,n);return n}}class Qp{constructor(t){this.id=zS++,this.code=t,this.usedTimes=0}}function kS(t){return t===ki||t===za||t===ka}function GS(t,e,n,i,s,r){let a=new ur,o=new jp,l=new Set,c=[],h=new Map,{logarithmicDepthBuffer:d,precision:u}=i,f={MeshDepthMaterial:"depth",MeshDistanceMaterial:"distance",MeshNormalMaterial:"normal",MeshBasicMaterial:"basic",MeshLambertMaterial:"lambert",MeshPhongMaterial:"phong",MeshToonMaterial:"toon",MeshStandardMaterial:"physical",MeshPhysicalMaterial:"physical",MeshMatcapMaterial:"matcap",LineBasicMaterial:"basic",LineDashedMaterial:"dashed",PointsMaterial:"points",ShadowMaterial:"shadow",SpriteMaterial:"sprite"};function g(y){if(l.add(y),y===0)return"uv";return`uv${y}`}function _(y,A,Z,D,q,ot){let z=D.fog,et=q.geometry,nt=y.isMeshStandardMaterial||y.isMeshLambertMaterial||y.isMeshPhongMaterial?D.environment:null,B=y.isMeshStandardMaterial||y.isMeshLambertMaterial&&!y.envMap||y.isMeshPhongMaterial&&!y.envMap,ft=e.get(y.envMap||nt,B),Q=!!ft&&ft.mapping===ar?ft.image.height:null,ct=f[y.type];if(y.precision!==null){if(u=i.getMaxPrecision(y.precision),u!==y.precision)Mt("WebGLProgram.getParameters:",y.precision,"not supported, using",u,"instead.")}let dt=et.morphAttributes.position||et.morphAttributes.normal||et.morphAttributes.color,Bt=dt!==void 0?dt.length:0,Dt=0;if(et.morphAttributes.position!==void 0)Dt=1;if(et.morphAttributes.normal!==void 0)Dt=2;if(et.morphAttributes.color!==void 0)Dt=3;let ee,Jt,F,O;if(ct){let ye=On[ct];ee=ye.vertexShader,Jt=ye.fragmentShader}else{ee=y.vertexShader,Jt=y.fragmentShader;let ye=o.getVertexShaderStage(y),he=o.getFragmentShaderStage(y);o.update(y,ye,he),F=ye.id,O=he.id}let N=t.getRenderTarget(),P=t.state.buffers.depth.getReversed(),G=q.isInstancedMesh===!0,J=q.isBatchedMesh===!0,_t=!!y.map,C=!!y.matcap,H=!!ft,U=!!y.aoMap,Y=!!y.lightMap,tt=!!y.bumpMap&&y.wireframe===!1,at=!!y.normalMap,mt=!!y.displacementMap,gt=!!y.emissiveMap,X=!!y.metalnessMap,v=!!y.roughnessMap,Lt=y.anisotropy>0,zt=y.clearcoat>0,$t=y.dispersion>0,w=y.retroreflectivity>0,M=y.iridescence>0,k=y.sheen>0,it=y.transmission>0,pt=Lt&&!!y.anisotropyMap,yt=zt&&!!y.clearcoatMap,St=zt&&!!y.clearcoatNormalMap,lt=zt&&!!y.clearcoatRoughnessMap,ut=M&&!!y.iridescenceMap,wt=M&&!!y.iridescenceThicknessMap,Ht=k&&!!y.sheenColorMap,bt=k&&!!y.sheenRoughnessMap,vt=!!y.specularMap,qt=!!y.specularColorMap,Yt=!!y.specularIntensityMap,ce=it&&!!y.transmissionMap,W=it&&!!y.thicknessMap,Tt=!!y.gradientMap,ht=!!y.alphaMap,Et=y.alphaTest>0,Ut=!!y.alphaHash,xt=!!y.extensions,Rt=bn;if(y.toneMapped){if(N===null||N.isXRRenderTarget===!0)Rt=t.toneMapping}let Qt={shaderID:ct,shaderType:y.type,shaderName:y.name,vertexShader:ee,fragmentShader:Jt,defines:y.defines,customVertexShaderID:F,customFragmentShaderID:O,isRawShaderMaterial:y.isRawShaderMaterial===!0,glslVersion:y.glslVersion,precision:u,batching:J,batchingColor:J&&q._colorsTexture!==null,instancing:G,instancingColor:G&&q.instanceColor!==null,instancingMorph:G&&q.morphTexture!==null,outputColorSpace:N===null?t.outputColorSpace:N.isXRRenderTarget===!0?N.texture.colorSpace:re.workingColorSpace,alphaToCoverage:!!y.alphaToCoverage,map:_t,matcap:C,envMap:H,envMapMode:H&&ft.mapping,envMapCubeUVHeight:Q,aoMap:U,lightMap:Y,bumpMap:tt,normalMap:at,displacementMap:mt,emissiveMap:gt,normalMapObjectSpace:at&&y.normalMapType===qd,normalMapTangentSpace:at&&y.normalMapType===vc,packedNormalMap:at&&y.normalMapType===vc&&kS(y.normalMap.format),metalnessMap:X,roughnessMap:v,anisotropy:Lt,anisotropyMap:pt,clearcoat:zt,clearcoatMap:yt,clearcoatNormalMap:St,clearcoatRoughnessMap:lt,dispersion:$t,retroreflection:w,iridescence:M,iridescenceMap:ut,iridescenceThicknessMap:wt,sheen:k,sheenColorMap:Ht,sheenRoughnessMap:bt,specularMap:vt,specularColorMap:qt,specularIntensityMap:Yt,transmission:it,transmissionMap:ce,thicknessMap:W,gradientMap:Tt,opaque:y.transparent===!1&&y.blending===rr&&y.alphaToCoverage===!1,alphaMap:ht,alphaTest:Et,alphaHash:Ut,combine:y.combine,mapUv:_t&&g(y.map.channel),aoMapUv:U&&g(y.aoMap.channel),lightMapUv:Y&&g(y.lightMap.channel),bumpMapUv:tt&&g(y.bumpMap.channel),normalMapUv:at&&g(y.normalMap.channel),displacementMapUv:mt&&g(y.displacementMap.channel),emissiveMapUv:gt&&g(y.emissiveMap.channel),metalnessMapUv:X&&g(y.metalnessMap.channel),roughnessMapUv:v&&g(y.roughnessMap.channel),anisotropyMapUv:pt&&g(y.anisotropyMap.channel),clearcoatMapUv:yt&&g(y.clearcoatMap.channel),clearcoatNormalMapUv:St&&g(y.clearcoatNormalMap.channel),clearcoatRoughnessMapUv:lt&&g(y.clearcoatRoughnessMap.channel),iridescenceMapUv:ut&&g(y.iridescenceMap.channel),iridescenceThicknessMapUv:wt&&g(y.iridescenceThicknessMap.channel),sheenColorMapUv:Ht&&g(y.sheenColorMap.channel),sheenRoughnessMapUv:bt&&g(y.sheenRoughnessMap.channel),specularMapUv:vt&&g(y.specularMap.channel),specularColorMapUv:qt&&g(y.specularColorMap.channel),specularIntensityMapUv:Yt&&g(y.specularIntensityMap.channel),transmissionMapUv:ce&&g(y.transmissionMap.channel),thicknessMapUv:W&&g(y.thicknessMap.channel),alphaMapUv:ht&&g(y.alphaMap.channel),vertexTangents:!!et.attributes.tangent&&(at||Lt),vertexNormals:!!et.attributes.normal,vertexColors:y.vertexColors,vertexAlphas:y.vertexColors===!0&&!!et.attributes.color&&et.attributes.color.itemSize===4,pointsUvs:q.isPoints===!0&&!!et.attributes.uv&&(_t||ht),fog:!!z,useFog:y.fog===!0,fogExp2:!!z&&z.isFogExp2,flatShading:y.wireframe===!1&&(y.flatShading===!0||et.attributes.normal===void 0&&at===!1&&(y.isMeshLambertMaterial||y.isMeshPhongMaterial||y.isMeshStandardMaterial||y.isMeshPhysicalMaterial)),sizeAttenuation:y.sizeAttenuation===!0,logarithmicDepthBuffer:d,reversedDepthBuffer:P,skinning:q.isSkinnedMesh===!0,hasPositionAttribute:et.attributes.position!==void 0,morphTargets:et.morphAttributes.position!==void 0,morphNormals:et.morphAttributes.normal!==void 0,morphColors:et.morphAttributes.color!==void 0,morphTargetsCount:Bt,morphTextureStride:Dt,numSunLights:A.sun.length,numDirLights:A.directional.length,numPointLights:A.point.length,numSpotLights:A.spot.length,numSpotLightMaps:A.spotLightMap.length,numRectAreaLights:A.rectArea.length,numHemiLights:A.hemi.length,numSunLightShadows:A.sunShadowMap.length,numDirLightShadows:A.directionalShadowMap.length,numPointLightShadows:A.pointShadowMap.length,numSpotLightShadows:A.spotShadowMap.length,numSpotLightShadowsWithMaps:A.numSpotLightShadowsWithMaps,numLightProbes:A.numLightProbes,numLightProbeGrids:ot.length,numClippingPlanes:r.numPlanes,numClipIntersection:r.numIntersection,dithering:y.dithering,shadowMapEnabled:t.shadowMap.enabled&&Z.length>0,shadowMapType:t.shadowMap.type,toneMapping:Rt,decodeVideoTexture:_t&&y.map.isVideoTexture===!0&&re.getTransfer(y.map.colorSpace)===xe,decodeVideoTextureEmissive:gt&&y.emissiveMap.isVideoTexture===!0&&re.getTransfer(y.emissiveMap.colorSpace)===xe,premultipliedAlpha:y.premultipliedAlpha,doubleSided:y.side===on,flipSided:y.side===Ve,useDepthPacking:y.depthPacking>=0,depthPacking:y.depthPacking||0,index0AttributeName:y.index0AttributeName,extensionClipCullDistance:xt&&y.extensions.clipCullDistance===!0&&n.has("WEBGL_clip_cull_distance"),extensionMultiDraw:(xt&&y.extensions.multiDraw===!0||J)&&n.has("WEBGL_multi_draw"),rendererExtensionParallelShaderCompile:n.has("KHR_parallel_shader_compile"),customProgramCacheKey:y.customProgramCacheKey()};return Qt.vertexUv1s=l.has(1),Qt.vertexUv2s=l.has(2),Qt.vertexUv3s=l.has(3),l.clear(),Qt}function m(y){let A=[];if(y.shaderID)A.push(y.shaderID);else A.push(y.customVertexShaderID),A.push(y.customFragmentShaderID);if(y.defines!==void 0)for(let Z in y.defines)A.push(Z),A.push(y.defines[Z]);if(y.isRawShaderMaterial===!1)p(A,y),b(A,y),A.push(t.outputColorSpace);return A.push(y.customProgramCacheKey),A.join()}function p(y,A){y.push(A.precision),y.push(A.outputColorSpace),y.push(A.envMapMode),y.push(A.envMapCubeUVHeight),y.push(A.mapUv),y.push(A.alphaMapUv),y.push(A.lightMapUv),y.push(A.aoMapUv),y.push(A.bumpMapUv),y.push(A.normalMapUv),y.push(A.displacementMapUv),y.push(A.emissiveMapUv),y.push(A.metalnessMapUv),y.push(A.roughnessMapUv),y.push(A.anisotropyMapUv),y.push(A.clearcoatMapUv),y.push(A.clearcoatNormalMapUv),y.push(A.clearcoatRoughnessMapUv),y.push(A.iridescenceMapUv),y.push(A.iridescenceThicknessMapUv),y.push(A.sheenColorMapUv),y.push(A.sheenRoughnessMapUv),y.push(A.specularMapUv),y.push(A.specularColorMapUv),y.push(A.specularIntensityMapUv),y.push(A.transmissionMapUv),y.push(A.thicknessMapUv),y.push(A.combine),y.push(A.fogExp2),y.push(A.sizeAttenuation),y.push(A.morphTargetsCount),y.push(A.morphAttributeCount),y.push(A.numSunLights),y.push(A.numDirLights),y.push(A.numPointLights),y.push(A.numSpotLights),y.push(A.numSpotLightMaps),y.push(A.numHemiLights),y.push(A.numRectAreaLights),y.push(A.numSunLightShadows),y.push(A.numDirLightShadows),y.push(A.numPointLightShadows),y.push(A.numSpotLightShadows),y.push(A.numSpotLightShadowsWithMaps),y.push(A.numLightProbes),y.push(A.shadowMapType),y.push(A.toneMapping),y.push(A.numClippingPlanes),y.push(A.numClipIntersection),y.push(A.depthPacking)}function b(y,A){if(a.disableAll(),A.instancing)a.enable(0);if(A.instancingColor)a.enable(1);if(A.instancingMorph)a.enable(2);if(A.matcap)a.enable(3);if(A.envMap)a.enable(4);if(A.normalMapObjectSpace)a.enable(5);if(A.normalMapTangentSpace)a.enable(6);if(A.clearcoat)a.enable(7);if(A.iridescence)a.enable(8);if(A.alphaTest)a.enable(9);if(A.vertexColors)a.enable(10);if(A.vertexAlphas)a.enable(11);if(A.vertexUv1s)a.enable(12);if(A.vertexUv2s)a.enable(13);if(A.vertexUv3s)a.enable(14);if(A.vertexTangents)a.enable(15);if(A.anisotropy)a.enable(16);if(A.alphaHash)a.enable(17);if(A.batching)a.enable(18);if(A.dispersion)a.enable(19);if(A.retroreflection)a.enable(24);if(A.batchingColor)a.enable(20);if(A.gradientMap)a.enable(21);if(A.packedNormalMap)a.enable(22);if(A.vertexNormals)a.enable(23);if(y.push(a.mask),a.disableAll(),A.fog)a.enable(0);if(A.useFog)a.enable(1);if(A.flatShading)a.enable(2);if(A.logarithmicDepthBuffer)a.enable(3);if(A.reversedDepthBuffer)a.enable(4);if(A.skinning)a.enable(5);if(A.morphTargets)a.enable(6);if(A.morphNormals)a.enable(7);if(A.morphColors)a.enable(8);if(A.premultipliedAlpha)a.enable(9);if(A.shadowMapEnabled)a.enable(10);if(A.doubleSided)a.enable(11);if(A.flipSided)a.enable(12);if(A.useDepthPacking)a.enable(13);if(A.dithering)a.enable(14);if(A.transmission)a.enable(15);if(A.sheen)a.enable(16);if(A.opaque)a.enable(17);if(A.pointsUvs)a.enable(18);if(A.decodeVideoTexture)a.enable(19);if(A.decodeVideoTextureEmissive)a.enable(20);if(A.alphaToCoverage)a.enable(21);if(A.numLightProbeGrids>0)a.enable(22);if(A.hasPositionAttribute)a.enable(23);y.push(a.mask)}function S(y){let A=f[y.type],Z;if(A){let D=On[A];Z=Rf.clone(D.uniforms)}else Z=y.uniforms;return Z}function x(y,A){let Z=h.get(A);if(Z!==void 0)++Z.usedTimes;else Z=new BS(t,A,y,s),c.push(Z),h.set(A,Z);return Z}function T(y){if(--y.usedTimes===0){let A=c.indexOf(y);c[A]=c[c.length-1],c.pop(),h.delete(y.cacheKey),y.destroy()}}function I(y){o.remove(y)}function R(){o.dispose()}return{getParameters:_,getProgramCacheKey:m,getUniforms:S,acquireProgram:x,releaseProgram:T,releaseShaderCache:I,programs:c,dispose:R}}function HS(){let t=new WeakMap;function e(a){return t.has(a)}function n(a){let o=t.get(a);if(o===void 0)o={},t.set(a,o);return o}function i(a){t.delete(a)}function s(a,o,l){t.get(a)[o]=l}function r(){t=new WeakMap}return{has:e,get:n,remove:i,update:s,dispose:r}}function VS(t,e){if(t.groupOrder!==e.groupOrder)return t.groupOrder-e.groupOrder;else if(t.renderOrder!==e.renderOrder)return t.renderOrder-e.renderOrder;else if(t.material.id!==e.material.id)return t.material.id-e.material.id;else if(t.materialVariant!==e.materialVariant)return t.materialVariant-e.materialVariant;else if(t.z!==e.z)return t.z-e.z;else return t.id-e.id}function zp(t,e){if(t.groupOrder!==e.groupOrder)return t.groupOrder-e.groupOrder;else if(t.renderOrder!==e.renderOrder)return t.renderOrder-e.renderOrder;else if(t.z!==e.z)return e.z-t.z;else return t.id-e.id}function kp(){let t=[],e=0,n=[],i=[],s=[];function r(){e=0,n.length=0,i.length=0,s.length=0}function a(u){let f=0;if(u.isInstancedMesh)f+=2;if(u.isSkinnedMesh)f+=1;return f}function o(u,f,g,_,m,p){let b=t[e];if(b===void 0)b={id:u.id,object:u,geometry:f,material:g,materialVariant:a(u),groupOrder:_,renderOrder:u.renderOrder,z:m,group:p},t[e]=b;else b.id=u.id,b.object=u,b.geometry=f,b.material=g,b.materialVariant=a(u),b.groupOrder=_,b.renderOrder=u.renderOrder,b.z=m,b.group=p;return e++,b}function l(u,f,g,_,m,p,b){if(b.reversedDepth===!0)m=-m;let S=o(u,f,g,_,m,p);if(g.transmission>0)i.push(S);else if(g.transparent===!0)s.push(S);else n.push(S)}function c(u,f,g,_,m,p){let b=o(u,f,g,_,m,p);if(g.transmission>0)i.unshift(b);else if(g.transparent===!0)s.unshift(b);else n.unshift(b)}function h(u,f){if(n.length>1)n.sort(u||VS);if(i.length>1)i.sort(f||zp);if(s.length>1)s.sort(f||zp)}function d(){for(let u=e,f=t.length;u<f;u++){let g=t[u];if(g.id===null)break;g.id=null,g.object=null,g.geometry=null,g.material=null,g.group=null}}return{opaque:n,transmissive:i,transparent:s,init:r,push:l,unshift:c,finish:d,sort:h}}function WS(){let t=new WeakMap;function e(i,s){let r=t.get(i),a;if(r===void 0)a=new kp,t.set(i,[a]);else if(s>=r.length)a=new kp,r.push(a);else a=r[s];return a}function n(){t=new WeakMap}return{get:e,dispose:n}}function XS(){let t={};return{get:function(e){if(t[e.id]!==void 0)return t[e.id];let n;switch(e.type){case"SunLight":case"DirectionalLight":n={direction:new L,color:new Pt};break;case"SpotLight":n={position:new L,direction:new L,color:new Pt,distance:0,coneCos:0,penumbraCos:0,decay:0};break;case"PointLight":n={position:new L,color:new Pt,distance:0,decay:0};break;case"HemisphereLight":n={direction:new L,skyColor:new Pt,groundColor:new Pt};break;case"RectAreaLight":n={color:new Pt,position:new L,halfWidth:new L,halfHeight:new L};break}return t[e.id]=n,n}}}function qS(){let t={};return{get:function(e){if(t[e.id]!==void 0)return t[e.id];let n;switch(e.type){case"SunLight":case"DirectionalLight":n={shadowIntensity:1,shadowBias:0,shadowNormalBias:0,shadowRadius:1,shadowMapSize:new rt};break;case"SpotLight":n={shadowIntensity:1,shadowBias:0,shadowNormalBias:0,shadowRadius:1,shadowMapSize:new rt};break;case"PointLight":n={shadowIntensity:1,shadowBias:0,shadowNormalBias:0,shadowRadius:1,shadowMapSize:new rt,shadowCameraNear:1,shadowCameraFar:1000};break}return t[e.id]=n,n}}}var YS=0;function ZS(t,e){return(e.castShadow?2:0)-(t.castShadow?2:0)+(e.map?1:0)-(t.map?1:0)}function JS(t){let e=new XS,n=qS(),i={version:0,hash:{sunLength:-1,directionalLength:-1,pointLength:-1,spotLength:-1,rectAreaLength:-1,hemiLength:-1,numSunShadows:-1,numDirectionalShadows:-1,numPointShadows:-1,numSpotShadows:-1,numSpotMaps:-1,numLightProbes:-1},ambient:[0,0,0],probe:[],sun:[],sunShadow:[],sunShadowMap:[],sunShadowMatrix:[],sunShadowCascade:[],directional:[],directionalShadow:[],directionalShadowMap:[],directionalShadowMatrix:[],spot:[],spotLightMap:[],spotShadow:[],spotShadowMap:[],spotLightMatrix:[],rectArea:[],rectAreaLTC1:null,rectAreaLTC2:null,point:[],pointShadow:[],pointShadowMap:[],pointShadowMatrix:[],hemi:[],numSpotLightShadowsWithMaps:0,numLightProbes:0};for(let c=0;c<9;c++)i.probe.push(new L);let s=new L,r=new Kt,a=new Kt;function o(c){let h=0,d=0,u=0;for(let q=0;q<9;q++)i.probe[q].set(0,0,0);let f=0,g=0,_=0,m=0,p=0,b=0,S=0,x=0,T=0,I=0,R=0,y=0,A=0,Z=0;c.sort(ZS);for(let q=0,ot=c.length;q<ot;q++){let z=c[q],{color:et,intensity:nt,distance:B}=z,ft=null;if(z.shadow&&z.shadow.map)if(z.shadow.map.texture.format===ki)ft=z.shadow.map.texture;else ft=z.shadow.map.depthTexture||z.shadow.map.texture;if(z.isAmbientLight)h+=et.r*nt,d+=et.g*nt,u+=et.b*nt;else if(z.isLightProbe){for(let Q=0;Q<9;Q++)i.probe[Q].addScaledVector(z.sh.coefficients[Q],nt);Z++}else if(z.isSunLight){let Q=e.get(z);if(Q.color.copy(z.color).multiplyScalar(z.intensity),z.castShadow){let ct=z.shadow,dt=n.get(z);dt.shadowIntensity=ct.intensity,dt.shadowBias=ct.bias,dt.shadowNormalBias=ct.normalBias,dt.shadowRadius=ct.radius,dt.shadowMapSize.copy(ct.mapSize).multiply(ct.getFrameExtents()),i.sunShadow[g]=dt,i.sunShadowMap[g]=ft;let Bt=ct.getViewportCount();for(let Dt=0;Dt<Bt;Dt++)i.sunShadowMatrix[_+Dt]=ct.getMatrix(Dt),i.sunShadowCascade[_+Dt]=ct._cascadeData[Dt];_+=Bt,g++}i.sun[f]=Q,f++}else if(z.isDirectionalLight){let Q=e.get(z);if(Q.color.copy(z.color).multiplyScalar(z.intensity),z.castShadow){let ct=z.shadow,dt=n.get(z);dt.shadowIntensity=ct.intensity,dt.shadowBias=ct.bias,dt.shadowNormalBias=ct.normalBias,dt.shadowRadius=ct.radius,dt.shadowMapSize=ct.mapSize,i.directionalShadow[m]=dt,i.directionalShadowMap[m]=ft,i.directionalShadowMatrix[m]=z.shadow.matrix,T++}i.directional[m]=Q,m++}else if(z.isSpotLight){let Q=e.get(z);Q.position.setFromMatrixPosition(z.matrixWorld),Q.color.copy(et).multiplyScalar(nt),Q.distance=B,Q.coneCos=Math.cos(z.angle),Q.penumbraCos=Math.cos(z.angle*(1-z.penumbra)),Q.decay=z.decay,i.spot[b]=Q;let ct=z.shadow;if(z.map){if(i.spotLightMap[y]=z.map,y++,ct.updateMatrices(z),z.castShadow)A++}if(i.spotLightMatrix[b]=ct.matrix,z.castShadow){let dt=n.get(z);dt.shadowIntensity=ct.intensity,dt.shadowBias=ct.bias,dt.shadowNormalBias=ct.normalBias,dt.shadowRadius=ct.radius,dt.shadowMapSize=ct.mapSize,i.spotShadow[b]=dt,i.spotShadowMap[b]=ft,R++}b++}else if(z.isRectAreaLight){let Q=e.get(z);Q.color.copy(et).multiplyScalar(nt),Q.halfWidth.set(z.width*0.5,0,0),Q.halfHeight.set(0,z.height*0.5,0),i.rectArea[S]=Q,S++}else if(z.isPointLight){let Q=e.get(z);if(Q.color.copy(z.color).multiplyScalar(z.intensity),Q.distance=z.distance,Q.decay=z.decay,z.castShadow){let ct=z.shadow,dt=n.get(z);dt.shadowIntensity=ct.intensity,dt.shadowBias=ct.bias,dt.shadowNormalBias=ct.normalBias,dt.shadowRadius=ct.radius,dt.shadowMapSize=ct.mapSize,dt.shadowCameraNear=ct.camera.near,dt.shadowCameraFar=ct.camera.far,i.pointShadow[p]=dt,i.pointShadowMap[p]=ft,i.pointShadowMatrix[p]=z.shadow.matrix,I++}i.point[p]=Q,p++}else if(z.isHemisphereLight){let Q=e.get(z);Q.skyColor.copy(z.color).multiplyScalar(nt),Q.groundColor.copy(z.groundColor).multiplyScalar(nt),i.hemi[x]=Q,x++}}if(S>0)if(t.has("OES_texture_float_linear")===!0)i.rectAreaLTC1=At.LTC_FLOAT_1,i.rectAreaLTC2=At.LTC_FLOAT_2;else i.rectAreaLTC1=At.LTC_HALF_1,i.rectAreaLTC2=At.LTC_HALF_2;i.ambient[0]=h,i.ambient[1]=d,i.ambient[2]=u;let D=i.hash;if(D.sunLength!==f||D.directionalLength!==m||D.pointLength!==p||D.spotLength!==b||D.rectAreaLength!==S||D.hemiLength!==x||D.numSunShadows!==g||D.numDirectionalShadows!==T||D.numPointShadows!==I||D.numSpotShadows!==R||D.numSpotMaps!==y||D.numLightProbes!==Z)i.sun.length=f,i.directional.length=m,i.spot.length=b,i.rectArea.length=S,i.point.length=p,i.hemi.length=x,i.sunShadow.length=g,i.sunShadowMap.length=g,i.sunShadowMatrix.length=_,i.sunShadowCascade.length=_,i.directionalShadow.length=T,i.directionalShadowMap.length=T,i.directionalShadowMatrix.length=T,i.pointShadow.length=I,i.pointShadowMap.length=I,i.pointShadowMatrix.length=I,i.spotShadow.length=R,i.spotShadowMap.length=R,i.spotLightMatrix.length=R+y-A,i.spotLightMap.length=y,i.numSpotLightShadowsWithMaps=A,i.numLightProbes=Z,D.sunLength=f,D.directionalLength=m,D.pointLength=p,D.spotLength=b,D.rectAreaLength=S,D.hemiLength=x,D.numSunShadows=g,D.numDirectionalShadows=T,D.numPointShadows=I,D.numSpotShadows=R,D.numSpotMaps=y,D.numLightProbes=Z,i.version=YS++}function l(c,h){let d=0,u=0,f=0,g=0,_=0,m=0,p=h.matrixWorldInverse;for(let b=0,S=c.length;b<S;b++){let x=c[b];if(x.isSunLight){let T=i.sun[d];T.direction.setFromMatrixPosition(x.matrixWorld),T.direction.transformDirection(p),d++}else if(x.isDirectionalLight){let T=i.directional[u];T.direction.setFromMatrixPosition(x.matrixWorld),s.setFromMatrixPosition(x.target.matrixWorld),T.direction.sub(s),T.direction.transformDirection(p),u++}else if(x.isSpotLight){let T=i.spot[g];T.position.setFromMatrixPosition(x.matrixWorld),T.position.applyMatrix4(p),T.direction.setFromMatrixPosition(x.matrixWorld),s.setFromMatrixPosition(x.target.matrixWorld),T.direction.sub(s),T.direction.transformDirection(p),g++}else if(x.isRectAreaLight){let T=i.rectArea[_];T.position.setFromMatrixPosition(x.matrixWorld),T.position.applyMatrix4(p),a.identity(),r.copy(x.matrixWorld),r.premultiply(p),a.extractRotation(r),T.halfWidth.set(x.width*0.5,0,0),T.halfHeight.set(0,x.height*0.5,0),T.halfWidth.applyMatrix4(a),T.halfHeight.applyMatrix4(a),_++}else if(x.isPointLight){let T=i.point[f];T.position.setFromMatrixPosition(x.matrixWorld),T.position.applyMatrix4(p),f++}else if(x.isHemisphereLight){let T=i.hemi[m];T.direction.setFromMatrixPosition(x.matrixWorld),T.direction.transformDirection(p),m++}}}return{setup:o,setupView:l,state:i}}function Gp(t){let e=new JS(t),n=[],i=[],s=[];function r(u){d.camera=u,n.length=0,i.length=0,s.length=0}function a(u){n.push(u)}function o(u){i.push(u)}function l(u){s.push(u)}function c(){e.setup(n)}function h(u){e.setupView(n,u)}let d={lightsArray:n,shadowsArray:i,lightProbeGridArray:s,camera:null,lights:e,transmissionRenderTarget:{},textureUnits:0};return{init:r,state:d,setupLights:c,setupLightsView:h,pushLight:a,pushShadow:o,pushLightProbeGrid:l}}function $S(t){let e=new WeakMap;function n(s,r=0){let a=e.get(s),o;if(a===void 0)o=new Gp(t),e.set(s,[o]);else if(r>=a.length)o=new Gp(t),a.push(o);else o=a[r];return o}function i(){e=new WeakMap}return{get:n,dispose:i}}var KS=`void main() {
	gl_Position = vec4( position, 1.0 );
}`,jS=`uniform sampler2D shadow_pass;
uniform vec2 resolution;
uniform float radius;
void main() {
	const float samples = float( VSM_SAMPLES );
	float mean = 0.0;
	float squared_mean = 0.0;
	float uvStride = samples <= 1.0 ? 0.0 : 2.0 / ( samples - 1.0 );
	float uvStart = samples <= 1.0 ? 0.0 : - 1.0;
	for ( float i = 0.0; i < samples; i ++ ) {
		float uvOffset = uvStart + i * uvStride;
		#ifdef HORIZONTAL_PASS
			vec2 distribution = texture2D( shadow_pass, ( gl_FragCoord.xy + vec2( uvOffset, 0.0 ) * radius ) / resolution ).rg;
			mean += distribution.x;
			squared_mean += distribution.y * distribution.y + distribution.x * distribution.x;
		#else
			float depth = texture2D( shadow_pass, ( gl_FragCoord.xy + vec2( 0.0, uvOffset ) * radius ) / resolution ).r;
			mean += depth;
			squared_mean += depth * depth;
		#endif
	}
	mean = mean / samples;
	squared_mean = squared_mean / samples;
	float std_dev = sqrt( max( 0.0, squared_mean - mean * mean ) );
	gl_FragColor = vec4( mean, std_dev, 0.0, 1.0 );
}`,QS=[new L(1,0,0),new L(-1,0,0),new L(0,1,0),new L(0,-1,0),new L(0,0,1),new L(0,0,-1)],tM=[new L(0,-1,0),new L(0,-1,0),new L(0,0,1),new L(0,0,-1),new L(0,-1,0),new L(0,-1,0)],Hp=new Kt,Ir=new L,Eh=new L;function eM(t,e,n){let i=new ui,s=new rt,r=new rt,a=new fe,o=new So,l=new Mo,c={},h=n.maxTextureSize,d={[Ms]:Ve,[Ve]:Ms,[on]:on},u=new ln({defines:{VSM_SAMPLES:8},uniforms:{shadow_pass:{value:null},resolution:{value:new rt},radius:{value:4}},vertexShader:KS,fragmentShader:jS}),f=u.clone();f.defines.HORIZONTAL_PASS=1;let g=new jt;g.setAttribute("position",new de(new Float32Array([-1,-1,0.5,3,-1,0.5,-1,3,0.5]),3));let _=new pe(g,u),m=this;this.enabled=!1,this.autoUpdate=!0,this.needsUpdate=!1,this.type=sr;let p=this.type;this.render=function(I,R,y){if(m.enabled===!1)return;if(m.autoUpdate===!1&&m.needsUpdate===!1)return;if(I.length===0)return;if(this.type===sd)Mt("WebGLShadowMap: PCFSoftShadowMap has been removed. Using PCFShadowMap instead."),this.type=sr;let A=t.getRenderTarget(),Z=t.getActiveCubeFace(),D=t.getActiveMipmapLevel(),q=t.state;if(q.setBlending(Ln),q.buffers.depth.getReversed()===!0)q.buffers.color.setClear(0,0,0,0);else q.buffers.color.setClear(1,1,1,1);q.buffers.depth.setTest(!0),q.setScissorTest(!1);let ot=p!==this.type;if(ot)R.traverse(function(z){if(z.material)if(Array.isArray(z.material))z.material.forEach((et)=>et.needsUpdate=!0);else z.material.needsUpdate=!0});for(let z=0,et=I.length;z<et;z++){let nt=I[z],B=nt.shadow;if(B===void 0){Mt("WebGLShadowMap:",nt,"has no shadow.");continue}if(B.autoUpdate===!1&&B.needsUpdate===!1)continue;s.copy(B.mapSize);let ft=B.getFrameExtents();if(s.multiply(ft),r.copy(B.mapSize),s.x>h||s.y>h){if(s.x>h)r.x=Math.floor(h/ft.x),s.x=r.x*ft.x,B.mapSize.x=r.x;if(s.y>h)r.y=Math.floor(h/ft.y),s.y=r.y*ft.y,B.mapSize.y=r.y}let Q=t.state.buffers.depth.getReversed();if(B.camera._reversedDepth=Q,B.map===null||ot===!0){if(B.map!==null){if(B.map.depthTexture!==null)B.map.depthTexture.dispose(),B.map.depthTexture=null;B.map.dispose()}if(this.type===Ss){if(nt.isPointLight){Mt("WebGLShadowMap: VSM shadow maps are not supported for PointLights. Use PCF or BasicShadowMap instead.");continue}B.map=new Je(s.x,s.y,{format:ki,type:Nn,minFilter:Qe,magFilter:Qe,generateMipmaps:!1}),B.map.texture.name=nt.name+".shadowMap",B.map.depthTexture=new Hi(s.x,s.y,$n),B.map.depthTexture.name=nt.name+".shadowMapDepth",B.map.depthTexture.format=Bi,B.map.depthTexture.compareFunction=null,B.map.depthTexture.minFilter=pi,B.map.depthTexture.magFilter=pi}else{if(nt.isPointLight)B.map=new Ch(s.x),B.map.depthTexture=new Uc(s.x,mi);else B.map=new Je(s.x,s.y),B.map.depthTexture=new Hi(s.x,s.y,mi);if(B.map.depthTexture.name=nt.name+".shadowMap",B.map.depthTexture.format=Bi,this.type===sr)B.map.depthTexture.compareFunction=Q?Va:Ha,B.map.depthTexture.minFilter=Qe,B.map.depthTexture.magFilter=Qe;else B.map.depthTexture.compareFunction=null,B.map.depthTexture.minFilter=pi,B.map.depthTexture.magFilter=pi}B.camera.updateProjectionMatrix()}if(B.map.isWebGLCubeRenderTarget!==!0&&(B.map.width!==s.x||B.map.height!==s.y))B.map.setSize(s.x,s.y);let ct=B.map.isWebGLCubeRenderTarget?6:B.getViewportCount();if(nt.isPointLight!==!0)B.updateMatrices(nt,y);for(let dt=0;dt<ct;dt++){let Bt=B.getCamera(dt);if(nt.isPointLight){let{camera:Dt,matrix:ee}=B,Jt=nt.distance||Dt.far;if(Jt!==Dt.far)Dt.far=Jt,Dt.updateProjectionMatrix();Ir.setFromMatrixPosition(nt.matrixWorld),Dt.position.copy(Ir),Eh.copy(Dt.position),Eh.add(QS[dt]),Dt.up.copy(tM[dt]),Dt.lookAt(Eh),Dt.updateMatrixWorld(),ee.makeTranslation(-Ir.x,-Ir.y,-Ir.z),Hp.multiplyMatrices(Dt.projectionMatrix,Dt.matrixWorldInverse),B._frustum.setFromProjectionMatrix(Hp,Dt.coordinateSystem,Dt.reversedDepth)}if(B.map.isWebGLCubeRenderTarget)t.setRenderTarget(B.map,dt),t.clear();else{if(dt===0)t.setRenderTarget(B.map),t.clear();let Dt=B.getViewport(dt);a.set(r.x*Dt.x,r.y*Dt.y,r.x*Dt.z,r.y*Dt.w),q.viewport(a)}i=B.getFrustum(dt),x(R,y,Bt,nt,this.type)}if(B.isPointLightShadow!==!0&&this.type===Ss)b(B,y);B.needsUpdate=!1}p=this.type,m.needsUpdate=!1,t.setRenderTarget(A,Z,D)};function b(I,R){let y=e.update(_);if(u.defines.VSM_SAMPLES!==I.blurSamples)u.defines.VSM_SAMPLES=I.blurSamples,f.defines.VSM_SAMPLES=I.blurSamples,u.needsUpdate=!0,f.needsUpdate=!0;if(I.mapPass===null)I.mapPass=new Je(s.x,s.y,{format:ki,type:Nn});else if(I.mapPass.width!==I.map.width||I.mapPass.height!==I.map.height)I.mapPass.setSize(I.map.width,I.map.height);u.uniforms.shadow_pass.value=I.map.depthTexture,u.uniforms.resolution.value.set(I.map.width,I.map.height),u.uniforms.radius.value=I.radius,t.setRenderTarget(I.mapPass),t.clear(),t.renderBufferDirect(R,null,y,u,_,null),f.uniforms.shadow_pass.value=I.mapPass.texture,f.uniforms.resolution.value.set(I.map.width,I.map.height),f.uniforms.radius.value=I.radius,t.setRenderTarget(I.map),t.clear(),t.renderBufferDirect(R,null,y,f,_,null)}function S(I,R,y,A){let Z=null,D=y.isPointLight===!0?I.customDistanceMaterial:I.customDepthMaterial;if(D!==void 0)Z=D;else if(Z=y.isPointLight===!0?l:o,t.localClippingEnabled&&R.clipShadows===!0&&Array.isArray(R.clippingPlanes)&&R.clippingPlanes.length!==0||R.displacementMap&&R.displacementScale!==0||R.alphaMap&&R.alphaTest>0||R.map&&R.alphaTest>0||R.alphaToCoverage===!0){let q=Z.uuid,ot=R.uuid,z=c[q];if(z===void 0)z={},c[q]=z;let et=z[ot];if(et===void 0)et=Z.clone(),z[ot]=et,R.addEventListener("dispose",T);Z=et}if(Z.visible=R.visible,Z.wireframe=R.wireframe,A===Ss)Z.side=R.shadowSide!==null?R.shadowSide:R.side;else Z.side=R.shadowSide!==null?R.shadowSide:d[R.side];if(Z.alphaMap=R.alphaMap,Z.alphaTest=R.alphaToCoverage===!0?0.5:R.alphaTest,Z.map=R.map,Z.clipShadows=R.clipShadows,Z.clippingPlanes=R.clippingPlanes,Z.clipIntersection=R.clipIntersection,Z.displacementMap=R.displacementMap,Z.displacementScale=R.displacementScale,Z.displacementBias=R.displacementBias,Z.wireframeLinewidth=R.wireframeLinewidth,Z.linewidth=R.linewidth,y.isPointLight===!0&&Z.isMeshDistanceMaterial===!0){let q=t.properties.get(Z);q.light=y}return Z}function x(I,R,y,A,Z){if(I.visible===!1)return;if(I.layers.test(R.layers)&&(I.isMesh||I.isLine||I.isPoints)){if((I.castShadow||I.receiveShadow&&Z===Ss)&&(!I.frustumCulled||I.intersectsFrustum(i))){I.modelViewMatrix.multiplyMatrices(y.matrixWorldInverse,I.matrixWorld);let ot=e.update(I),z=I.material;if(Array.isArray(z)){let et=ot.groups;for(let nt=0,B=et.length;nt<B;nt++){let ft=et[nt],Q=z[ft.materialIndex];if(Q&&Q.visible){let ct=S(I,Q,A,Z);I.onBeforeShadow(t,I,R,y,ot,ct,ft),t.renderBufferDirect(y,null,ot,ct,I,ft),I.onAfterShadow(t,I,R,y,ot,ct,ft)}}}else if(z.visible){let et=S(I,z,A,Z);I.onBeforeShadow(t,I,R,y,ot,et,null),t.renderBufferDirect(y,null,ot,et,I,null),I.onAfterShadow(t,I,R,y,ot,et,null)}}}let q=I.children;for(let ot=0,z=q.length;ot<z;ot++)x(q[ot],R,y,A,Z)}function T(I){I.target.removeEventListener("dispose",T);for(let y in c){let A=c[y],Z=I.target.uuid;if(Z in A)A[Z].dispose(),delete A[Z]}}}function nM(t,e){function n(){let W=!1,Tt=new fe,ht=null,Et=new fe(0,0,0,0);return{setMask:function(Ut){if(ht!==Ut&&!W)t.colorMask(Ut,Ut,Ut,Ut),ht=Ut},setLocked:function(Ut){W=Ut},setClear:function(Ut,xt,Rt,Qt,ye){if(ye===!0)Ut*=Qt,xt*=Qt,Rt*=Qt;if(Tt.set(Ut,xt,Rt,Qt),Et.equals(Tt)===!1)t.clearColor(Ut,xt,Rt,Qt),Et.copy(Tt)},reset:function(){W=!1,ht=null,Et.set(-1,0,0,0)}}}function i(){let W=!1,Tt=!1,ht=null,Et=null,Ut=null;return{setReversed:function(xt){if(Tt!==xt){let Rt=e.get("EXT_clip_control");if(xt)Rt.clipControlEXT(Rt.LOWER_LEFT_EXT,Rt.ZERO_TO_ONE_EXT);else Rt.clipControlEXT(Rt.LOWER_LEFT_EXT,Rt.NEGATIVE_ONE_TO_ONE_EXT);Tt=xt;let Qt=Ut;Ut=null,this.setClear(Qt)}},getReversed:function(){return Tt},setTest:function(xt){if(xt)N(t.DEPTH_TEST);else P(t.DEPTH_TEST)},setMask:function(xt){if(ht!==xt&&!W)t.depthMask(xt),ht=xt},setFunc:function(xt){if(Tt)xt=sf[xt];if(Et!==xt){switch(xt){case Ed:t.depthFunc(t.NEVER);break;case Ad:t.depthFunc(t.ALWAYS);break;case wd:t.depthFunc(t.LESS);break;case Il:t.depthFunc(t.LEQUAL);break;case Rd:t.depthFunc(t.EQUAL);break;case Cd:t.depthFunc(t.GEQUAL);break;case Id:t.depthFunc(t.GREATER);break;case Pd:t.depthFunc(t.NOTEQUAL);break;default:t.depthFunc(t.LEQUAL)}Et=xt}},setLocked:function(xt){W=xt},setClear:function(xt){if(Ut!==xt){if(Ut=xt,Tt)xt=1-xt;t.clearDepth(xt)}},reset:function(){W=!1,ht=null,Et=null,Ut=null,Tt=!1}}}function s(){let W=!1,Tt=null,ht=null,Et=null,Ut=null,xt=null,Rt=null,Qt=null,ye=null;return{setTest:function(he){if(!W)if(he)N(t.STENCIL_TEST);else P(t.STENCIL_TEST)},setMask:function(he){if(Tt!==he&&!W)t.stencilMask(he),Tt=he},setFunc:function(he,wn,Bn){if(ht!==he||Et!==wn||Ut!==Bn)t.stencilFunc(he,wn,Bn),ht=he,Et=wn,Ut=Bn},setOp:function(he,wn,Bn){if(xt!==he||Rt!==wn||Qt!==Bn)t.stencilOp(he,wn,Bn),xt=he,Rt=wn,Qt=Bn},setLocked:function(he){W=he},setClear:function(he){if(ye!==he)t.clearStencil(he),ye=he},reset:function(){W=!1,Tt=null,ht=null,Et=null,Ut=null,xt=null,Rt=null,Qt=null,ye=null}}}let r=new n,a=new i,o=new s,l=new WeakMap,c=new WeakMap,h={},d={},u={},f=new WeakMap,g=[],_=null,m=!1,p=null,b=null,S=null,x=null,T=null,I=null,R=null,y=new Pt(0,0,0),A=0,Z=!1,D=null,q=null,ot=null,z=null,et=null,nt=t.getParameter(t.MAX_COMBINED_TEXTURE_IMAGE_UNITS),B=!1,ft=0,Q=t.getParameter(t.VERSION);if(Q.indexOf("WebGL")!==-1)ft=parseFloat(/^WebGL (\d)/.exec(Q)[1]),B=ft>=1;else if(Q.indexOf("OpenGL ES")!==-1)ft=parseFloat(/^OpenGL ES (\d)/.exec(Q)[1]),B=ft>=2;let ct=null,dt={},Bt=t.getParameter(t.SCISSOR_BOX),Dt=t.getParameter(t.VIEWPORT),ee=new fe().fromArray(Bt),Jt=new fe().fromArray(Dt);function F(W,Tt,ht,Et){let Ut=new Uint8Array(4),xt=t.createTexture();t.bindTexture(W,xt),t.texParameteri(W,t.TEXTURE_MIN_FILTER,t.NEAREST),t.texParameteri(W,t.TEXTURE_MAG_FILTER,t.NEAREST);for(let Rt=0;Rt<ht;Rt++)if(W===t.TEXTURE_3D||W===t.TEXTURE_2D_ARRAY)t.texImage3D(Tt,0,t.RGBA,1,1,Et,0,t.RGBA,t.UNSIGNED_BYTE,Ut);else t.texImage2D(Tt+Rt,0,t.RGBA,1,1,0,t.RGBA,t.UNSIGNED_BYTE,Ut);return xt}let O={};O[t.TEXTURE_2D]=F(t.TEXTURE_2D,t.TEXTURE_2D,1),O[t.TEXTURE_CUBE_MAP]=F(t.TEXTURE_CUBE_MAP,t.TEXTURE_CUBE_MAP_POSITIVE_X,6),O[t.TEXTURE_2D_ARRAY]=F(t.TEXTURE_2D_ARRAY,t.TEXTURE_2D_ARRAY,1,1),O[t.TEXTURE_3D]=F(t.TEXTURE_3D,t.TEXTURE_3D,1,1),r.setClear(0,0,0,1),a.setClear(1),o.setClear(0),N(t.DEPTH_TEST),a.setFunc(Il),tt(!1),at(Al),N(t.CULL_FACE),U(Ln);function N(W){if(h[W]!==!0)t.enable(W),h[W]=!0}function P(W){if(h[W]!==!1)t.disable(W),h[W]=!1}function G(W,Tt){if(u[W]!==Tt){if(t.bindFramebuffer(W,Tt),u[W]=Tt,W===t.DRAW_FRAMEBUFFER)u[t.FRAMEBUFFER]=Tt;if(W===t.FRAMEBUFFER)u[t.DRAW_FRAMEBUFFER]=Tt;return!0}return!1}function J(W,Tt){let ht=g,Et=!1;if(W){if(ht=f.get(Tt),ht===void 0)ht=[],f.set(Tt,ht);let Ut=W.textures;if(ht.length!==Ut.length||ht[0]!==t.COLOR_ATTACHMENT0){for(let xt=0,Rt=Ut.length;xt<Rt;xt++)ht[xt]=t.COLOR_ATTACHMENT0+xt;ht.length=Ut.length,Et=!0}}else if(ht[0]!==t.BACK)ht[0]=t.BACK,Et=!0;if(Et)t.drawBuffers(ht)}function _t(W){if(_!==W)return t.useProgram(W),_=W,!0;return!1}let C={[bs]:t.FUNC_ADD,[ad]:t.FUNC_SUBTRACT,[od]:t.FUNC_REVERSE_SUBTRACT};C[ld]=t.MIN,C[cd]=t.MAX;let H={[hd]:t.ZERO,[ud]:t.ONE,[dd]:t.SRC_COLOR,[pd]:t.SRC_ALPHA,[vd]:t.SRC_ALPHA_SATURATE,[xd]:t.DST_COLOR,[gd]:t.DST_ALPHA,[fd]:t.ONE_MINUS_SRC_COLOR,[md]:t.ONE_MINUS_SRC_ALPHA,[yd]:t.ONE_MINUS_DST_COLOR,[_d]:t.ONE_MINUS_DST_ALPHA,[Sd]:t.CONSTANT_COLOR,[Md]:t.ONE_MINUS_CONSTANT_COLOR,[bd]:t.CONSTANT_ALPHA,[Td]:t.ONE_MINUS_CONSTANT_ALPHA};function U(W,Tt,ht,Et,Ut,xt,Rt,Qt,ye,he){if(W===Ln){if(m===!0)P(t.BLEND),m=!1;return}if(m===!1)N(t.BLEND),m=!0;if(W!==rd){if(W!==p||he!==Z){if(b!==bs||T!==bs)t.blendEquation(t.FUNC_ADD),b=bs,T=bs;if(he)switch(W){case rr:t.blendFuncSeparate(t.ONE,t.ONE_MINUS_SRC_ALPHA,t.ONE,t.ONE_MINUS_SRC_ALPHA);break;case wl:t.blendFunc(t.ONE,t.ONE);break;case Rl:t.blendFuncSeparate(t.ZERO,t.ONE_MINUS_SRC_COLOR,t.ZERO,t.ONE);break;case Cl:t.blendFuncSeparate(t.DST_COLOR,t.ONE_MINUS_SRC_ALPHA,t.ZERO,t.ONE);break;default:Wt("WebGLState: Invalid blending: ",W);break}else switch(W){case rr:t.blendFuncSeparate(t.SRC_ALPHA,t.ONE_MINUS_SRC_ALPHA,t.ONE,t.ONE_MINUS_SRC_ALPHA);break;case wl:t.blendFuncSeparate(t.SRC_ALPHA,t.ONE,t.ONE,t.ONE);break;case Rl:Wt("WebGLState: SubtractiveBlending requires material.premultipliedAlpha = true");break;case Cl:Wt("WebGLState: MultiplyBlending requires material.premultipliedAlpha = true");break;default:Wt("WebGLState: Invalid blending: ",W);break}S=null,x=null,I=null,R=null,y.set(0,0,0),A=0,p=W,Z=he}return}if(Ut=Ut||Tt,xt=xt||ht,Rt=Rt||Et,Tt!==b||Ut!==T)t.blendEquationSeparate(C[Tt],C[Ut]),b=Tt,T=Ut;if(ht!==S||Et!==x||xt!==I||Rt!==R)t.blendFuncSeparate(H[ht],H[Et],H[xt],H[Rt]),S=ht,x=Et,I=xt,R=Rt;if(Qt.equals(y)===!1||ye!==A)t.blendColor(Qt.r,Qt.g,Qt.b,ye),y.copy(Qt),A=ye;p=W,Z=!1}function Y(W,Tt){W.side===on?P(t.CULL_FACE):N(t.CULL_FACE);let ht=W.side===Ve;if(Tt)ht=!ht;tt(ht),W.blending===rr&&W.transparent===!1?U(Ln):U(W.blending,W.blendEquation,W.blendSrc,W.blendDst,W.blendEquationAlpha,W.blendSrcAlpha,W.blendDstAlpha,W.blendColor,W.blendAlpha,W.premultipliedAlpha),a.setFunc(W.depthFunc),a.setTest(W.depthTest),a.setMask(W.depthWrite),r.setMask(W.colorWrite);let Et=W.stencilWrite;if(o.setTest(Et),Et)o.setMask(W.stencilWriteMask),o.setFunc(W.stencilFunc,W.stencilRef,W.stencilFuncMask),o.setOp(W.stencilFail,W.stencilZFail,W.stencilZPass);gt(W.polygonOffset,W.polygonOffsetFactor,W.polygonOffsetUnits),W.alphaToCoverage===!0?N(t.SAMPLE_ALPHA_TO_COVERAGE):P(t.SAMPLE_ALPHA_TO_COVERAGE)}function tt(W){if(D!==W){if(W)t.frontFace(t.CW);else t.frontFace(t.CCW);D=W}}function at(W){if(W!==nd){if(N(t.CULL_FACE),W!==q)if(W===Al)t.cullFace(t.BACK);else if(W===id)t.cullFace(t.FRONT);else t.cullFace(t.FRONT_AND_BACK)}else P(t.CULL_FACE);q=W}function mt(W){if(W!==ot){if(B)t.lineWidth(W);ot=W}}function gt(W,Tt,ht){if(W){if(N(t.POLYGON_OFFSET_FILL),z!==Tt||et!==ht){if(z=Tt,et=ht,a.getReversed())Tt=-Tt;t.polygonOffset(Tt,ht)}}else P(t.POLYGON_OFFSET_FILL)}function X(W){if(W)N(t.SCISSOR_TEST);else P(t.SCISSOR_TEST)}function v(W){if(W===void 0)W=t.TEXTURE0+nt-1;if(ct!==W)t.activeTexture(W),ct=W}function Lt(W,Tt,ht){if(ht===void 0)if(ct===null)ht=t.TEXTURE0+nt-1;else ht=ct;let Et=dt[ht];if(Et===void 0)Et={type:void 0,texture:void 0},dt[ht]=Et;if(Et.type!==W||Et.texture!==Tt){if(ct!==ht)t.activeTexture(ht),ct=ht;t.bindTexture(W,Tt||O[W]),Et.type=W,Et.texture=Tt}}function zt(){let W=dt[ct];if(W!==void 0&&W.type!==void 0)t.bindTexture(W.type,null),W.type=void 0,W.texture=void 0}function $t(){try{t.compressedTexImage2D(...arguments)}catch(W){Wt("WebGLState:",W)}}function w(){try{t.compressedTexImage3D(...arguments)}catch(W){Wt("WebGLState:",W)}}function M(){try{t.texSubImage2D(...arguments)}catch(W){Wt("WebGLState:",W)}}function k(){try{t.texSubImage3D(...arguments)}catch(W){Wt("WebGLState:",W)}}function it(){try{t.compressedTexSubImage2D(...arguments)}catch(W){Wt("WebGLState:",W)}}function pt(){try{t.compressedTexSubImage3D(...arguments)}catch(W){Wt("WebGLState:",W)}}function yt(){try{t.texStorage2D(...arguments)}catch(W){Wt("WebGLState:",W)}}function St(){try{t.texStorage3D(...arguments)}catch(W){Wt("WebGLState:",W)}}function lt(){try{t.texImage2D(...arguments)}catch(W){Wt("WebGLState:",W)}}function ut(){try{t.texImage3D(...arguments)}catch(W){Wt("WebGLState:",W)}}function wt(W){if(d[W]!==void 0)return d[W];else return t.getParameter(W)}function Ht(W,Tt){if(d[W]!==Tt)t.pixelStorei(W,Tt),d[W]=Tt}function bt(W){if(ee.equals(W)===!1)t.scissor(W.x,W.y,W.z,W.w),ee.copy(W)}function vt(W){if(Jt.equals(W)===!1)t.viewport(W.x,W.y,W.z,W.w),Jt.copy(W)}function qt(W,Tt){let ht=c.get(Tt);if(ht===void 0)ht=new WeakMap,c.set(Tt,ht);let Et=ht.get(W);if(Et===void 0)Et=t.getUniformBlockIndex(Tt,W.name),ht.set(W,Et)}function Yt(W,Tt){let Et=c.get(Tt).get(W);if(l.get(Tt)!==Et)t.uniformBlockBinding(Tt,Et,W.__bindingPointIndex),l.set(Tt,Et)}function ce(){t.disable(t.BLEND),t.disable(t.CULL_FACE),t.disable(t.DEPTH_TEST),t.disable(t.POLYGON_OFFSET_FILL),t.disable(t.SCISSOR_TEST),t.disable(t.STENCIL_TEST),t.disable(t.SAMPLE_ALPHA_TO_COVERAGE),t.blendEquation(t.FUNC_ADD),t.blendFunc(t.ONE,t.ZERO),t.blendFuncSeparate(t.ONE,t.ZERO,t.ONE,t.ZERO),t.blendColor(0,0,0,0),t.colorMask(!0,!0,!0,!0),t.clearColor(0,0,0,0),t.depthMask(!0),t.depthFunc(t.LESS),a.setReversed(!1),t.clearDepth(1),t.stencilMask(4294967295),t.stencilFunc(t.ALWAYS,0,4294967295),t.stencilOp(t.KEEP,t.KEEP,t.KEEP),t.clearStencil(0),t.cullFace(t.BACK),t.frontFace(t.CCW),t.polygonOffset(0,0),t.activeTexture(t.TEXTURE0),t.bindFramebuffer(t.FRAMEBUFFER,null),t.bindFramebuffer(t.DRAW_FRAMEBUFFER,null),t.bindFramebuffer(t.READ_FRAMEBUFFER,null),t.useProgram(null),t.lineWidth(1),t.scissor(0,0,t.canvas.width,t.canvas.height),t.viewport(0,0,t.canvas.width,t.canvas.height),t.pixelStorei(t.PACK_ALIGNMENT,4),t.pixelStorei(t.UNPACK_ALIGNMENT,4),t.pixelStorei(t.UNPACK_FLIP_Y_WEBGL,!1),t.pixelStorei(t.UNPACK_PREMULTIPLY_ALPHA_WEBGL,!1),t.pixelStorei(t.UNPACK_COLORSPACE_CONVERSION_WEBGL,t.BROWSER_DEFAULT_WEBGL),t.pixelStorei(t.PACK_ROW_LENGTH,0),t.pixelStorei(t.PACK_SKIP_PIXELS,0),t.pixelStorei(t.PACK_SKIP_ROWS,0),t.pixelStorei(t.UNPACK_ROW_LENGTH,0),t.pixelStorei(t.UNPACK_IMAGE_HEIGHT,0),t.pixelStorei(t.UNPACK_SKIP_PIXELS,0),t.pixelStorei(t.UNPACK_SKIP_ROWS,0),t.pixelStorei(t.UNPACK_SKIP_IMAGES,0),h={},d={},ct=null,dt={},u={},f=new WeakMap,g=[],_=null,m=!1,p=null,b=null,S=null,x=null,T=null,I=null,R=null,y=new Pt(0,0,0),A=0,Z=!1,D=null,q=null,ot=null,z=null,et=null,ee.set(0,0,t.canvas.width,t.canvas.height),Jt.set(0,0,t.canvas.width,t.canvas.height),r.reset(),a.reset(),o.reset()}return{buffers:{color:r,depth:a,stencil:o},enable:N,disable:P,bindFramebuffer:G,drawBuffers:J,useProgram:_t,setBlending:U,setMaterial:Y,setFlipSided:tt,setCullFace:at,setLineWidth:mt,setPolygonOffset:gt,setScissorTest:X,activeTexture:v,bindTexture:Lt,unbindTexture:zt,compressedTexImage2D:$t,compressedTexImage3D:w,texImage2D:lt,texImage3D:ut,pixelStorei:Ht,getParameter:wt,updateUBOMapping:qt,uniformBlockBinding:Yt,texStorage2D:yt,texStorage3D:St,texSubImage2D:M,texSubImage3D:k,compressedTexSubImage2D:it,compressedTexSubImage3D:pt,scissor:bt,viewport:vt,reset:ce}}function iM(t,e,n,i,s,r,a){let o=e.has("WEBGL_multisampled_render_to_texture")?e.get("WEBGL_multisampled_render_to_texture"):null,l=typeof navigator>"u"?!1:/OculusBrowser/g.test(navigator.userAgent),c=new rt,h=new WeakMap,d=new Set,u,f=new WeakMap,g=!1;try{g=typeof OffscreenCanvas<"u"&&new OffscreenCanvas(1,1).getContext("2d")!==null}catch(w){}function _(w,M){return g?new OffscreenCanvas(w,M):_s("canvas")}function m(w,M,k){let it=1,pt=$t(w);if(pt.width>k||pt.height>k)it=k/Math.max(pt.width,pt.height);if(it<1)if(typeof HTMLImageElement<"u"&&w instanceof HTMLImageElement||typeof HTMLCanvasElement<"u"&&w instanceof HTMLCanvasElement||typeof ImageBitmap<"u"&&w instanceof ImageBitmap||typeof VideoFrame<"u"&&w instanceof VideoFrame){let yt=Math.floor(it*pt.width),St=Math.floor(it*pt.height);if(u===void 0)u=_(yt,St);let lt=M?_(yt,St):u;return lt.width=yt,lt.height=St,lt.getContext("2d").drawImage(w,0,0,yt,St),Mt("WebGLRenderer: Texture has been resized from ("+pt.width+"x"+pt.height+") to ("+yt+"x"+St+")."),lt}else{if("data"in w)Mt("WebGLRenderer: Image in DataTexture is too big ("+pt.width+"x"+pt.height+").");return w}return w}function p(w){return w.generateMipmaps}function b(w){t.generateMipmap(w)}function S(w){if(w.isWebGLCubeRenderTarget)return t.TEXTURE_CUBE_MAP;if(w.isWebGL3DRenderTarget)return t.TEXTURE_3D;if(w.isWebGLArrayRenderTarget||w.isCompressedArrayTexture)return t.TEXTURE_2D_ARRAY;return t.TEXTURE_2D}function x(w,M,k,it,pt,yt=!1){if(w!==null){if(t[w]!==void 0)return t[w];Mt("WebGLRenderer: Attempt to use non-existing WebGL internal format '"+w+"'")}let St;if(it){if(St=e.get("EXT_texture_norm16"),!St)Mt("WebGLRenderer: Unable to use normalized textures without EXT_texture_norm16 extension")}let lt=M;if(M===t.RED){if(k===t.FLOAT)lt=t.R32F;if(k===t.HALF_FLOAT)lt=t.R16F;if(k===t.UNSIGNED_BYTE)lt=t.R8;if(k===t.UNSIGNED_SHORT&&St)lt=St.R16_EXT;if(k===t.SHORT&&St)lt=St.R16_SNORM_EXT}if(M===t.RED_INTEGER){if(k===t.UNSIGNED_BYTE)lt=t.R8UI;if(k===t.UNSIGNED_SHORT)lt=t.R16UI;if(k===t.UNSIGNED_INT)lt=t.R32UI;if(k===t.BYTE)lt=t.R8I;if(k===t.SHORT)lt=t.R16I;if(k===t.INT)lt=t.R32I}if(M===t.RG){if(k===t.FLOAT)lt=t.RG32F;if(k===t.HALF_FLOAT)lt=t.RG16F;if(k===t.UNSIGNED_BYTE)lt=t.RG8;if(k===t.UNSIGNED_SHORT&&St)lt=St.RG16_EXT;if(k===t.SHORT&&St)lt=St.RG16_SNORM_EXT}if(M===t.RG_INTEGER){if(k===t.UNSIGNED_BYTE)lt=t.RG8UI;if(k===t.UNSIGNED_SHORT)lt=t.RG16UI;if(k===t.UNSIGNED_INT)lt=t.RG32UI;if(k===t.BYTE)lt=t.RG8I;if(k===t.SHORT)lt=t.RG16I;if(k===t.INT)lt=t.RG32I}if(M===t.RGB_INTEGER){if(k===t.UNSIGNED_BYTE)lt=t.RGB8UI;if(k===t.UNSIGNED_SHORT)lt=t.RGB16UI;if(k===t.UNSIGNED_INT)lt=t.RGB32UI;if(k===t.BYTE)lt=t.RGB8I;if(k===t.SHORT)lt=t.RGB16I;if(k===t.INT)lt=t.RGB32I}if(M===t.RGBA_INTEGER){if(k===t.UNSIGNED_BYTE)lt=t.RGBA8UI;if(k===t.UNSIGNED_SHORT)lt=t.RGBA16UI;if(k===t.UNSIGNED_INT)lt=t.RGBA32UI;if(k===t.BYTE)lt=t.RGBA8I;if(k===t.SHORT)lt=t.RGBA16I;if(k===t.INT)lt=t.RGBA32I}if(M===t.RGB){if(k===t.UNSIGNED_SHORT&&St)lt=St.RGB16_EXT;if(k===t.SHORT&&St)lt=St.RGB16_SNORM_EXT;if(k===t.UNSIGNED_INT_5_9_9_9_REV)lt=t.RGB9_E5;if(k===t.UNSIGNED_INT_10F_11F_11F_REV)lt=t.R11F_G11F_B10F}if(M===t.RGBA){let ut=yt?Mc:re.getTransfer(pt);if(k===t.FLOAT)lt=t.RGBA32F;if(k===t.HALF_FLOAT)lt=t.RGBA16F;if(k===t.UNSIGNED_BYTE)lt=ut===xe?t.SRGB8_ALPHA8:t.RGBA8;if(k===t.UNSIGNED_SHORT&&St)lt=St.RGBA16_EXT;if(k===t.SHORT&&St)lt=St.RGBA16_SNORM_EXT;if(k===t.UNSIGNED_SHORT_4_4_4_4)lt=t.RGBA4;if(k===t.UNSIGNED_SHORT_5_5_5_1)lt=t.RGB5_A1}if(lt===t.R16F||lt===t.R32F||lt===t.RG16F||lt===t.RG32F||lt===t.RGBA16F||lt===t.RGBA32F)e.get("EXT_color_buffer_float");return lt}function T(w,M){let k;if(w){if(M===null||M===mi||M===As)k=t.DEPTH24_STENCIL8;else if(M===$n)k=t.DEPTH32F_STENCIL8;else if(M===lr)k=t.DEPTH24_STENCIL8,Mt("DepthTexture: 16 bit depth attachment is not supported with stencil. Using 24-bit attachment.")}else if(M===null||M===mi||M===As)k=t.DEPTH_COMPONENT24;else if(M===$n)k=t.DEPTH_COMPONENT32F;else if(M===lr)k=t.DEPTH_COMPONENT16;return k}function I(w,M){if(p(w)===!0||w.isFramebufferTexture&&w.minFilter!==pi&&w.minFilter!==Qe)return Math.log2(Math.max(M.width,M.height))+1;else if(w.mipmaps!==void 0&&w.mipmaps.length>0)return w.mipmaps.length;else if(w.isCompressedTexture&&Array.isArray(w.image))return M.mipmaps.length;else return 1}function R(w){let M=w.target;if(M.removeEventListener("dispose",R),A(M),M.isVideoTexture)h.delete(M);if(M.isHTMLTexture)d.delete(M)}function y(w){let M=w.target;M.removeEventListener("dispose",y),D(M)}function A(w){let M=i.get(w);if(M.__webglInit===void 0)return;let k=w.source,it=f.get(k);if(it){let pt=it[M.__cacheKey];if(pt.usedTimes--,pt.usedTimes===0)Z(w);if(Object.keys(it).length===0)f.delete(k)}i.remove(w)}function Z(w){let M=i.get(w);t.deleteTexture(M.__webglTexture);let k=w.source,it=f.get(k);delete it[M.__cacheKey],a.memory.textures--}function D(w){let M=i.get(w);if(w.depthTexture)w.depthTexture.dispose(),i.remove(w.depthTexture);if(w.isWebGLCubeRenderTarget)for(let it=0;it<6;it++){if(Array.isArray(M.__webglFramebuffer[it]))for(let pt=0;pt<M.__webglFramebuffer[it].length;pt++)t.deleteFramebuffer(M.__webglFramebuffer[it][pt]);else t.deleteFramebuffer(M.__webglFramebuffer[it]);if(M.__webglDepthbuffer)t.deleteRenderbuffer(M.__webglDepthbuffer[it])}else{if(Array.isArray(M.__webglFramebuffer))for(let it=0;it<M.__webglFramebuffer.length;it++)t.deleteFramebuffer(M.__webglFramebuffer[it]);else t.deleteFramebuffer(M.__webglFramebuffer);if(M.__webglDepthbuffer)t.deleteRenderbuffer(M.__webglDepthbuffer);if(M.__webglMultisampledFramebuffer)t.deleteFramebuffer(M.__webglMultisampledFramebuffer);if(M.__webglColorRenderbuffer){for(let it=0;it<M.__webglColorRenderbuffer.length;it++)if(M.__webglColorRenderbuffer[it])t.deleteRenderbuffer(M.__webglColorRenderbuffer[it])}if(M.__webglDepthRenderbuffer)t.deleteRenderbuffer(M.__webglDepthRenderbuffer)}let k=w.textures;for(let it=0,pt=k.length;it<pt;it++){let yt=i.get(k[it]);if(yt.__webglTexture)t.deleteTexture(yt.__webglTexture),a.memory.textures--;i.remove(k[it])}i.remove(w)}let q=0;function ot(){q=0}function z(){return q}function et(w){q=w}function nt(){let w=q;if(w>=s.maxTextures)Mt("WebGLTextures: Trying to use "+(w+1)+" texture units while this GPU supports only "+s.maxTextures);return q+=1,w}function B(w){let M=[];return M.push(w.wrapS),M.push(w.wrapT),M.push(w.wrapR||0),M.push(w.magFilter),M.push(w.minFilter),M.push(w.anisotropy),M.push(w.internalFormat),M.push(w.format),M.push(w.type),M.push(w.generateMipmaps),M.push(w.premultiplyAlpha),M.push(w.flipY),M.push(w.unpackAlignment),M.push(w.colorSpace),M.join()}function ft(w,M){let k=i.get(w);if(w.isVideoTexture)Lt(w);if(w.isRenderTargetTexture===!1&&w.isExternalTexture!==!0&&w.version>0&&k.__version!==w.version){let it=w.image;if(it===null)Mt("WebGLRenderer: Texture marked for update but no image data found.");else if(it.complete===!1)Mt("WebGLRenderer: Texture marked for update but image is incomplete");else{P(k,w,M);return}}else if(w.isExternalTexture)k.__webglTexture=w.sourceTexture?w.sourceTexture:null;n.bindTexture(t.TEXTURE_2D,k.__webglTexture,t.TEXTURE0+M)}function Q(w,M){let k=i.get(w);if(w.isRenderTargetTexture===!1&&w.version>0&&k.__version!==w.version){P(k,w,M);return}else if(w.isExternalTexture)k.__webglTexture=w.sourceTexture?w.sourceTexture:null;n.bindTexture(t.TEXTURE_2D_ARRAY,k.__webglTexture,t.TEXTURE0+M)}function ct(w,M){let k=i.get(w);if(w.isRenderTargetTexture===!1&&w.version>0&&k.__version!==w.version){P(k,w,M);return}n.bindTexture(t.TEXTURE_3D,k.__webglTexture,t.TEXTURE0+M)}function dt(w,M){let k=i.get(w);if(w.isCubeDepthTexture!==!0&&w.version>0&&k.__version!==w.version){G(k,w,M);return}n.bindTexture(t.TEXTURE_CUBE_MAP,k.__webglTexture,t.TEXTURE0+M)}let Bt={[La]:t.REPEAT,[Es]:t.CLAMP_TO_EDGE,[Na]:t.MIRRORED_REPEAT},Dt={[pi]:t.NEAREST,[Ud]:t.NEAREST_MIPMAP_NEAREST,[or]:t.NEAREST_MIPMAP_LINEAR,[Qe]:t.LINEAR,[Da]:t.LINEAR_MIPMAP_NEAREST,[Oi]:t.LINEAR_MIPMAP_LINEAR},ee={[Yd]:t.NEVER,[jd]:t.ALWAYS,[Zd]:t.LESS,[Ha]:t.LEQUAL,[Jd]:t.EQUAL,[Va]:t.GEQUAL,[$d]:t.GREATER,[Kd]:t.NOTEQUAL};function Jt(w,M){if(M.type===$n&&e.has("OES_texture_float_linear")===!1&&(M.magFilter===Qe||M.magFilter===Da||M.magFilter===or||M.magFilter===Oi||M.minFilter===Qe||M.minFilter===Da||M.minFilter===or||M.minFilter===Oi))Mt("WebGLRenderer: Unable to use linear filtering with floating point textures. OES_texture_float_linear not supported on this device.");if(t.texParameteri(w,t.TEXTURE_WRAP_S,Bt[M.wrapS]),t.texParameteri(w,t.TEXTURE_WRAP_T,Bt[M.wrapT]),w===t.TEXTURE_3D||w===t.TEXTURE_2D_ARRAY)t.texParameteri(w,t.TEXTURE_WRAP_R,Bt[M.wrapR]);if(t.texParameteri(w,t.TEXTURE_MAG_FILTER,Dt[M.magFilter]),t.texParameteri(w,t.TEXTURE_MIN_FILTER,Dt[M.minFilter]),M.compareFunction)t.texParameteri(w,t.TEXTURE_COMPARE_MODE,t.COMPARE_REF_TO_TEXTURE),t.texParameteri(w,t.TEXTURE_COMPARE_FUNC,ee[M.compareFunction]);if(e.has("EXT_texture_filter_anisotropic")===!0){if(M.magFilter===pi)return;if(M.minFilter!==or&&M.minFilter!==Oi)return;if(M.type===$n&&e.has("OES_texture_float_linear")===!1)return;if(M.anisotropy>1||i.get(M).__currentAnisotropy){let k=e.get("EXT_texture_filter_anisotropic");t.texParameterf(w,k.TEXTURE_MAX_ANISOTROPY_EXT,Math.min(M.anisotropy,s.getMaxAnisotropy())),i.get(M).__currentAnisotropy=M.anisotropy}}}function F(w,M){let k=!1;if(w.__webglInit===void 0)w.__webglInit=!0,M.addEventListener("dispose",R);let it=M.source,pt=f.get(it);if(pt===void 0)pt={},f.set(it,pt);let yt=B(M);if(yt!==w.__cacheKey){if(pt[yt]===void 0)pt[yt]={texture:t.createTexture(),usedTimes:0},a.memory.textures++,k=!0;pt[yt].usedTimes++;let St=pt[w.__cacheKey];if(St!==void 0){if(pt[w.__cacheKey].usedTimes--,St.usedTimes===0)Z(M)}w.__cacheKey=yt,w.__webglTexture=pt[yt].texture}return k}function O(w,M,k){return Math.floor(Math.floor(w/k)/M)}function N(w,M,k,it){let yt=w.updateRanges;if(yt.length===0)n.texSubImage2D(t.TEXTURE_2D,0,0,0,M.width,M.height,k,it,M.data);else{yt.sort((Ht,bt)=>Ht.start-bt.start);let St=0;for(let Ht=1;Ht<yt.length;Ht++){let bt=yt[St],vt=yt[Ht],qt=bt.start+bt.count,Yt=O(vt.start,M.width,4),ce=O(bt.start,M.width,4);if(vt.start<=qt+1&&Yt===ce&&O(vt.start+vt.count-1,M.width,4)===Yt)bt.count=Math.max(bt.count,vt.start+vt.count-bt.start);else++St,yt[St]=vt}yt.length=St+1;let lt=n.getParameter(t.UNPACK_ROW_LENGTH),ut=n.getParameter(t.UNPACK_SKIP_PIXELS),wt=n.getParameter(t.UNPACK_SKIP_ROWS);n.pixelStorei(t.UNPACK_ROW_LENGTH,M.width);for(let Ht=0,bt=yt.length;Ht<bt;Ht++){let vt=yt[Ht],qt=Math.floor(vt.start/4),Yt=Math.ceil(vt.count/4),ce=qt%M.width,W=Math.floor(qt/M.width),Tt=Yt,ht=1;n.pixelStorei(t.UNPACK_SKIP_PIXELS,ce),n.pixelStorei(t.UNPACK_SKIP_ROWS,W),n.texSubImage2D(t.TEXTURE_2D,0,ce,W,Tt,1,k,it,M.data)}w.clearUpdateRanges(),n.pixelStorei(t.UNPACK_ROW_LENGTH,lt),n.pixelStorei(t.UNPACK_SKIP_PIXELS,ut),n.pixelStorei(t.UNPACK_SKIP_ROWS,wt)}}function P(w,M,k){let it=t.TEXTURE_2D;if(M.isDataArrayTexture||M.isCompressedArrayTexture)it=t.TEXTURE_2D_ARRAY;if(M.isData3DTexture)it=t.TEXTURE_3D;let pt=F(w,M),yt=M.source;n.bindTexture(it,w.__webglTexture,t.TEXTURE0+k);let St=i.get(yt);if(yt.version!==St.__version||pt===!0){if(n.activeTexture(t.TEXTURE0+k),(typeof ImageBitmap<"u"&&M.image instanceof ImageBitmap)===!1){let ht=re.getPrimaries(re.workingColorSpace),Et=M.colorSpace===Gi?null:re.getPrimaries(M.colorSpace),Ut=M.colorSpace===Gi||ht===Et?t.NONE:t.BROWSER_DEFAULT_WEBGL;n.pixelStorei(t.UNPACK_FLIP_Y_WEBGL,M.flipY),n.pixelStorei(t.UNPACK_PREMULTIPLY_ALPHA_WEBGL,M.premultiplyAlpha),n.pixelStorei(t.UNPACK_COLORSPACE_CONVERSION_WEBGL,Ut)}n.pixelStorei(t.UNPACK_ALIGNMENT,M.unpackAlignment);let ut=m(M.image,!1,s.maxTextureSize);ut=zt(M,ut);let wt=r.convert(M.format,M.colorSpace),Ht=r.convert(M.type),bt=x(M.internalFormat,wt,Ht,M.normalized,M.colorSpace,M.isVideoTexture);Jt(it,M);let vt,qt=M.mipmaps,Yt=M.isVideoTexture!==!0,ce=St.__version===void 0||pt===!0,W=yt.dataReady,Tt=I(M,ut);if(M.isDepthTexture){if(bt=T(M.format===zi,M.type),ce)if(Yt)n.texStorage2D(t.TEXTURE_2D,1,bt,ut.width,ut.height);else n.texImage2D(t.TEXTURE_2D,0,bt,ut.width,ut.height,0,wt,Ht,null)}else if(M.isDataTexture)if(qt.length>0){if(Yt&&ce)n.texStorage2D(t.TEXTURE_2D,Tt,bt,qt[0].width,qt[0].height);for(let ht=0,Et=qt.length;ht<Et;ht++)if(vt=qt[ht],Yt){if(W)n.texSubImage2D(t.TEXTURE_2D,ht,0,0,vt.width,vt.height,wt,Ht,vt.data)}else n.texImage2D(t.TEXTURE_2D,ht,bt,vt.width,vt.height,0,wt,Ht,vt.data);M.generateMipmaps=!1}else if(Yt){if(ce)n.texStorage2D(t.TEXTURE_2D,Tt,bt,ut.width,ut.height);if(W)N(M,ut,wt,Ht)}else n.texImage2D(t.TEXTURE_2D,0,bt,ut.width,ut.height,0,wt,Ht,ut.data);else if(M.isCompressedTexture)if(M.isCompressedArrayTexture){if(Yt&&ce)n.texStorage3D(t.TEXTURE_2D_ARRAY,Tt,bt,qt[0].width,qt[0].height,ut.depth);for(let ht=0,Et=qt.length;ht<Et;ht++)if(vt=qt[ht],M.format!==Dn)if(wt!==null)if(Yt){if(W)if(M.layerUpdates.size>0){let Ut=Po(vt.width,vt.height,M.format,M.type);for(let xt of M.layerUpdates){let Rt=vt.data.subarray(xt*Ut/vt.data.BYTES_PER_ELEMENT,(xt+1)*Ut/vt.data.BYTES_PER_ELEMENT);n.compressedTexSubImage3D(t.TEXTURE_2D_ARRAY,ht,0,0,xt,vt.width,vt.height,1,wt,Rt)}}else n.compressedTexSubImage3D(t.TEXTURE_2D_ARRAY,ht,0,0,0,vt.width,vt.height,ut.depth,wt,vt.data)}else n.compressedTexImage3D(t.TEXTURE_2D_ARRAY,ht,bt,vt.width,vt.height,ut.depth,0,vt.data,0,0);else Mt("WebGLRenderer: Attempt to load unsupported compressed texture format in .uploadTexture()");else if(Yt){if(W)n.texSubImage3D(t.TEXTURE_2D_ARRAY,ht,0,0,0,vt.width,vt.height,ut.depth,wt,Ht,vt.data)}else n.texImage3D(t.TEXTURE_2D_ARRAY,ht,bt,vt.width,vt.height,ut.depth,0,wt,Ht,vt.data);if(M.layerUpdates.size>0)M.clearLayerUpdates()}else{if(Yt&&ce)n.texStorage2D(t.TEXTURE_2D,Tt,bt,qt[0].width,qt[0].height);for(let ht=0,Et=qt.length;ht<Et;ht++)if(vt=qt[ht],M.format!==Dn)if(wt!==null)if(Yt){if(W)n.compressedTexSubImage2D(t.TEXTURE_2D,ht,0,0,vt.width,vt.height,wt,vt.data)}else n.compressedTexImage2D(t.TEXTURE_2D,ht,bt,vt.width,vt.height,0,vt.data);else Mt("WebGLRenderer: Attempt to load unsupported compressed texture format in .uploadTexture()");else if(Yt){if(W)n.texSubImage2D(t.TEXTURE_2D,ht,0,0,vt.width,vt.height,wt,Ht,vt.data)}else n.texImage2D(t.TEXTURE_2D,ht,bt,vt.width,vt.height,0,wt,Ht,vt.data)}else if(M.isDataArrayTexture)if(Yt){if(ce)n.texStorage3D(t.TEXTURE_2D_ARRAY,Tt,bt,ut.width,ut.height,ut.depth);if(W)if(M.layerUpdates.size>0){let ht=Po(ut.width,ut.height,M.format,M.type);for(let Et of M.layerUpdates){let Ut=ut.data.subarray(Et*ht/ut.data.BYTES_PER_ELEMENT,(Et+1)*ht/ut.data.BYTES_PER_ELEMENT);n.texSubImage3D(t.TEXTURE_2D_ARRAY,0,0,0,Et,ut.width,ut.height,1,wt,Ht,Ut)}M.clearLayerUpdates()}else n.texSubImage3D(t.TEXTURE_2D_ARRAY,0,0,0,0,ut.width,ut.height,ut.depth,wt,Ht,ut.data)}else n.texImage3D(t.TEXTURE_2D_ARRAY,0,bt,ut.width,ut.height,ut.depth,0,wt,Ht,ut.data);else if(M.isData3DTexture)if(Yt){if(ce)n.texStorage3D(t.TEXTURE_3D,Tt,bt,ut.width,ut.height,ut.depth);if(W)n.texSubImage3D(t.TEXTURE_3D,0,0,0,0,ut.width,ut.height,ut.depth,wt,Ht,ut.data)}else n.texImage3D(t.TEXTURE_3D,0,bt,ut.width,ut.height,ut.depth,0,wt,Ht,ut.data);else if(M.isFramebufferTexture){if(ce)if(Yt)n.texStorage2D(t.TEXTURE_2D,Tt,bt,ut.width,ut.height);else{let ht=ut.width,Et=ut.height;for(let Ut=0;Ut<Tt;Ut++)n.texImage2D(t.TEXTURE_2D,Ut,bt,ht,Et,0,wt,Ht,null),ht>>=1,Et>>=1}}else if(M.isHTMLTexture){if("texElementImage2D"in t){let ht=t.canvas;if(!ht.hasAttribute("layoutsubtree"))ht.setAttribute("layoutsubtree","true");if(ut.parentNode!==ht){ht.appendChild(ut),d.add(M),ht.onpaint=(Et)=>{let Ut=Et.changedElements;for(let xt of d)if(Ut.includes(xt.image))xt.needsUpdate=!0},ht.requestPaint();return}if(t.texElementImage2D.length===3)t.texElementImage2D(t.TEXTURE_2D,t.RGBA8,ut);else{let{RGBA:Ut,RGBA:xt,UNSIGNED_BYTE:Rt}=t;t.texElementImage2D(t.TEXTURE_2D,0,Ut,xt,Rt,ut)}t.texParameteri(t.TEXTURE_2D,t.TEXTURE_MIN_FILTER,t.LINEAR),t.texParameteri(t.TEXTURE_2D,t.TEXTURE_WRAP_S,t.CLAMP_TO_EDGE),t.texParameteri(t.TEXTURE_2D,t.TEXTURE_WRAP_T,t.CLAMP_TO_EDGE)}}else if(qt.length>0){if(Yt&&ce){let ht=$t(qt[0]);n.texStorage2D(t.TEXTURE_2D,Tt,bt,ht.width,ht.height)}for(let ht=0,Et=qt.length;ht<Et;ht++)if(vt=qt[ht],Yt){if(W)n.texSubImage2D(t.TEXTURE_2D,ht,0,0,wt,Ht,vt)}else n.texImage2D(t.TEXTURE_2D,ht,bt,wt,Ht,vt);M.generateMipmaps=!1}else if(Yt){if(ce){let ht=$t(ut);n.texStorage2D(t.TEXTURE_2D,Tt,bt,ht.width,ht.height)}if(W)n.texSubImage2D(t.TEXTURE_2D,0,0,0,wt,Ht,ut)}else n.texImage2D(t.TEXTURE_2D,0,bt,wt,Ht,ut);if(p(M))b(it);if(St.__version=yt.version,M.onUpdate)M.onUpdate(M)}w.__version=M.version}function G(w,M,k){if(M.image.length!==6)return;let it=F(w,M),pt=M.source;n.bindTexture(t.TEXTURE_CUBE_MAP,w.__webglTexture,t.TEXTURE0+k);let yt=i.get(pt);if(pt.version!==yt.__version||it===!0){n.activeTexture(t.TEXTURE0+k);let St=re.getPrimaries(re.workingColorSpace),lt=M.colorSpace===Gi?null:re.getPrimaries(M.colorSpace),ut=M.colorSpace===Gi||St===lt?t.NONE:t.BROWSER_DEFAULT_WEBGL;n.pixelStorei(t.UNPACK_FLIP_Y_WEBGL,M.flipY),n.pixelStorei(t.UNPACK_PREMULTIPLY_ALPHA_WEBGL,M.premultiplyAlpha),n.pixelStorei(t.UNPACK_ALIGNMENT,M.unpackAlignment),n.pixelStorei(t.UNPACK_COLORSPACE_CONVERSION_WEBGL,ut);let wt=M.isCompressedTexture||M.image[0].isCompressedTexture,Ht=M.image[0]&&M.image[0].isDataTexture,bt=[];for(let xt=0;xt<6;xt++){if(!wt&&!Ht)bt[xt]=m(M.image[xt],!0,s.maxCubemapSize);else bt[xt]=Ht?M.image[xt].image:M.image[xt];bt[xt]=zt(M,bt[xt])}let vt=bt[0],qt=r.convert(M.format,M.colorSpace),Yt=r.convert(M.type),ce=x(M.internalFormat,qt,Yt,M.normalized,M.colorSpace),W=M.isVideoTexture!==!0,Tt=yt.__version===void 0||it===!0,ht=pt.dataReady,Et=I(M,vt);Jt(t.TEXTURE_CUBE_MAP,M);let Ut;if(wt){if(W&&Tt)n.texStorage2D(t.TEXTURE_CUBE_MAP,Et,ce,vt.width,vt.height);for(let xt=0;xt<6;xt++){Ut=bt[xt].mipmaps;for(let Rt=0;Rt<Ut.length;Rt++){let Qt=Ut[Rt];if(M.format!==Dn)if(qt!==null)if(W){if(ht)n.compressedTexSubImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt,0,0,Qt.width,Qt.height,qt,Qt.data)}else n.compressedTexImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt,ce,Qt.width,Qt.height,0,Qt.data);else Mt("WebGLRenderer: Attempt to load unsupported compressed texture format in .setTextureCube()");else if(W){if(ht)n.texSubImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt,0,0,Qt.width,Qt.height,qt,Yt,Qt.data)}else n.texImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt,ce,Qt.width,Qt.height,0,qt,Yt,Qt.data)}}}else{if(Ut=M.mipmaps,W&&Tt){if(Ut.length>0)Et++;let xt=$t(bt[0]);n.texStorage2D(t.TEXTURE_CUBE_MAP,Et,ce,xt.width,xt.height)}for(let xt=0;xt<6;xt++)if(Ht){if(W){if(ht)n.texSubImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,0,0,0,bt[xt].width,bt[xt].height,qt,Yt,bt[xt].data)}else n.texImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,0,ce,bt[xt].width,bt[xt].height,0,qt,Yt,bt[xt].data);for(let Rt=0;Rt<Ut.length;Rt++){let ye=Ut[Rt].image[xt].image;if(W){if(ht)n.texSubImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt+1,0,0,ye.width,ye.height,qt,Yt,ye.data)}else n.texImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt+1,ce,ye.width,ye.height,0,qt,Yt,ye.data)}}else{if(W){if(ht)n.texSubImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,0,0,0,qt,Yt,bt[xt])}else n.texImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,0,ce,qt,Yt,bt[xt]);for(let Rt=0;Rt<Ut.length;Rt++){let Qt=Ut[Rt];if(W){if(ht)n.texSubImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt+1,0,0,qt,Yt,Qt.image[xt])}else n.texImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+xt,Rt+1,ce,qt,Yt,Qt.image[xt])}}}if(p(M))b(t.TEXTURE_CUBE_MAP);if(yt.__version=pt.version,M.onUpdate)M.onUpdate(M)}w.__version=M.version}function J(w,M,k,it,pt,yt){let St=r.convert(k.format,k.colorSpace),lt=r.convert(k.type),ut=x(k.internalFormat,St,lt,k.normalized,k.colorSpace),wt=i.get(M),Ht=i.get(k);if(Ht.__renderTarget=M,!wt.__hasExternalTextures){let bt=Math.max(1,M.width>>yt),vt=Math.max(1,M.height>>yt);if(pt===t.TEXTURE_3D||pt===t.TEXTURE_2D_ARRAY)n.texImage3D(pt,yt,ut,bt,vt,M.depth,0,St,lt,null);else n.texImage2D(pt,yt,ut,bt,vt,0,St,lt,null)}if(n.bindFramebuffer(t.FRAMEBUFFER,w),v(M))o.framebufferTexture2DMultisampleEXT(t.FRAMEBUFFER,it,pt,Ht.__webglTexture,0,X(M));else if(pt===t.TEXTURE_2D||pt>=t.TEXTURE_CUBE_MAP_POSITIVE_X&&pt<=t.TEXTURE_CUBE_MAP_NEGATIVE_Z)t.framebufferTexture2D(t.FRAMEBUFFER,it,pt,Ht.__webglTexture,yt);n.bindFramebuffer(t.FRAMEBUFFER,null)}function _t(w,M,k){if(t.bindRenderbuffer(t.RENDERBUFFER,w),M.depthBuffer){let it=M.depthTexture,pt=it&&it.isDepthTexture?it.type:null,yt=T(M.stencilBuffer,pt),St=M.stencilBuffer?t.DEPTH_STENCIL_ATTACHMENT:t.DEPTH_ATTACHMENT;if(v(M))o.renderbufferStorageMultisampleEXT(t.RENDERBUFFER,X(M),yt,M.width,M.height);else if(k)t.renderbufferStorageMultisample(t.RENDERBUFFER,X(M),yt,M.width,M.height);else t.renderbufferStorage(t.RENDERBUFFER,yt,M.width,M.height);t.framebufferRenderbuffer(t.FRAMEBUFFER,St,t.RENDERBUFFER,w)}else{let it=M.textures;for(let pt=0;pt<it.length;pt++){let yt=it[pt],St=r.convert(yt.format,yt.colorSpace),lt=r.convert(yt.type),ut=x(yt.internalFormat,St,lt,yt.normalized,yt.colorSpace);if(v(M))o.renderbufferStorageMultisampleEXT(t.RENDERBUFFER,X(M),ut,M.width,M.height);else if(k)t.renderbufferStorageMultisample(t.RENDERBUFFER,X(M),ut,M.width,M.height);else t.renderbufferStorage(t.RENDERBUFFER,ut,M.width,M.height)}}t.bindRenderbuffer(t.RENDERBUFFER,null)}function C(w,M,k){let it=M.isWebGLCubeRenderTarget===!0;if(n.bindFramebuffer(t.FRAMEBUFFER,w),!(M.depthTexture&&M.depthTexture.isDepthTexture))throw Error("THREE.WebGLTextures: renderTarget.depthTexture must be an instance of THREE.DepthTexture.");let pt=i.get(M.depthTexture);if(pt.__renderTarget=M,!pt.__webglTexture||M.depthTexture.image.width!==M.width||M.depthTexture.image.height!==M.height)M.depthTexture.image.width=M.width,M.depthTexture.image.height=M.height,M.depthTexture.needsUpdate=!0;if(it){if(pt.__webglInit===void 0)pt.__webglInit=!0,M.depthTexture.addEventListener("dispose",R);if(pt.__webglTexture===void 0){pt.__webglTexture=t.createTexture(),n.bindTexture(t.TEXTURE_CUBE_MAP,pt.__webglTexture),Jt(t.TEXTURE_CUBE_MAP,M.depthTexture);let wt=r.convert(M.depthTexture.format),Ht=r.convert(M.depthTexture.type),bt;if(M.depthTexture.format===Bi)bt=t.DEPTH_COMPONENT24;else if(M.depthTexture.format===zi)bt=t.DEPTH24_STENCIL8;for(let vt=0;vt<6;vt++)t.texImage2D(t.TEXTURE_CUBE_MAP_POSITIVE_X+vt,0,bt,M.width,M.height,0,wt,Ht,null)}}else ft(M.depthTexture,0);let yt=pt.__webglTexture,St=X(M),lt=it?t.TEXTURE_CUBE_MAP_POSITIVE_X+k:t.TEXTURE_2D,ut=M.depthTexture.format===zi?t.DEPTH_STENCIL_ATTACHMENT:t.DEPTH_ATTACHMENT;if(M.depthTexture.format===Bi)if(v(M))o.framebufferTexture2DMultisampleEXT(t.FRAMEBUFFER,ut,lt,yt,0,St);else t.framebufferTexture2D(t.FRAMEBUFFER,ut,lt,yt,0);else if(M.depthTexture.format===zi)if(v(M))o.framebufferTexture2DMultisampleEXT(t.FRAMEBUFFER,ut,lt,yt,0,St);else t.framebufferTexture2D(t.FRAMEBUFFER,ut,lt,yt,0);else throw Error("THREE.WebGLTextures: Unknown depthTexture format.")}function H(w){let M=i.get(w),k=w.isWebGLCubeRenderTarget===!0;if(M.__boundDepthTexture!==w.depthTexture){let it=w.depthTexture;if(M.__depthDisposeCallback)M.__depthDisposeCallback();if(it){let pt=()=>{delete M.__boundDepthTexture,delete M.__depthDisposeCallback,it.removeEventListener("dispose",pt)};it.addEventListener("dispose",pt),M.__depthDisposeCallback=pt}M.__boundDepthTexture=it}if(w.depthTexture&&!M.__autoAllocateDepthBuffer)if(k)for(let it=0;it<6;it++)C(M.__webglFramebuffer[it],w,it);else{let it=w.texture.mipmaps;if(it&&it.length>0)C(M.__webglFramebuffer[0],w,0);else C(M.__webglFramebuffer,w,0)}else if(k){M.__webglDepthbuffer=[];for(let it=0;it<6;it++)if(n.bindFramebuffer(t.FRAMEBUFFER,M.__webglFramebuffer[it]),M.__webglDepthbuffer[it]===void 0)M.__webglDepthbuffer[it]=t.createRenderbuffer(),_t(M.__webglDepthbuffer[it],w,!1);else{let pt=w.stencilBuffer?t.DEPTH_STENCIL_ATTACHMENT:t.DEPTH_ATTACHMENT,yt=M.__webglDepthbuffer[it];t.bindRenderbuffer(t.RENDERBUFFER,yt),t.framebufferRenderbuffer(t.FRAMEBUFFER,pt,t.RENDERBUFFER,yt)}}else{let it=w.texture.mipmaps;if(it&&it.length>0)n.bindFramebuffer(t.FRAMEBUFFER,M.__webglFramebuffer[0]);else n.bindFramebuffer(t.FRAMEBUFFER,M.__webglFramebuffer);if(M.__webglDepthbuffer===void 0)M.__webglDepthbuffer=t.createRenderbuffer(),_t(M.__webglDepthbuffer,w,!1);else{let pt=w.stencilBuffer?t.DEPTH_STENCIL_ATTACHMENT:t.DEPTH_ATTACHMENT,yt=M.__webglDepthbuffer;t.bindRenderbuffer(t.RENDERBUFFER,yt),t.framebufferRenderbuffer(t.FRAMEBUFFER,pt,t.RENDERBUFFER,yt)}}n.bindFramebuffer(t.FRAMEBUFFER,null)}function U(w,M,k){let it=i.get(w);if(M!==void 0)J(it.__webglFramebuffer,w,w.texture,t.COLOR_ATTACHMENT0,t.TEXTURE_2D,0);if(k!==void 0)H(w)}function Y(w){let M=w.texture,k=i.get(w),it=i.get(M);w.addEventListener("dispose",y);let pt=w.textures,yt=w.isWebGLCubeRenderTarget===!0,St=pt.length>1;if(!St){if(it.__webglTexture===void 0)it.__webglTexture=t.createTexture();it.__version=M.version,a.memory.textures++}if(yt){k.__webglFramebuffer=[];for(let lt=0;lt<6;lt++)if(M.mipmaps&&M.mipmaps.length>0){k.__webglFramebuffer[lt]=[];for(let ut=0;ut<M.mipmaps.length;ut++)k.__webglFramebuffer[lt][ut]=t.createFramebuffer()}else k.__webglFramebuffer[lt]=t.createFramebuffer()}else{if(M.mipmaps&&M.mipmaps.length>0){k.__webglFramebuffer=[];for(let lt=0;lt<M.mipmaps.length;lt++)k.__webglFramebuffer[lt]=t.createFramebuffer()}else k.__webglFramebuffer=t.createFramebuffer();if(St)for(let lt=0,ut=pt.length;lt<ut;lt++){let wt=i.get(pt[lt]);if(wt.__webglTexture===void 0)wt.__webglTexture=t.createTexture(),a.memory.textures++}if(w.samples>0&&v(w)===!1){k.__webglMultisampledFramebuffer=t.createFramebuffer(),k.__webglColorRenderbuffer=[],n.bindFramebuffer(t.FRAMEBUFFER,k.__webglMultisampledFramebuffer);for(let lt=0;lt<pt.length;lt++){let ut=pt[lt];k.__webglColorRenderbuffer[lt]=t.createRenderbuffer(),t.bindRenderbuffer(t.RENDERBUFFER,k.__webglColorRenderbuffer[lt]);let wt=r.convert(ut.format,ut.colorSpace),Ht=r.convert(ut.type),bt=x(ut.internalFormat,wt,Ht,ut.normalized,ut.colorSpace,w.isXRRenderTarget===!0),vt=X(w);t.renderbufferStorageMultisample(t.RENDERBUFFER,vt,bt,w.width,w.height),t.framebufferRenderbuffer(t.FRAMEBUFFER,t.COLOR_ATTACHMENT0+lt,t.RENDERBUFFER,k.__webglColorRenderbuffer[lt])}if(t.bindRenderbuffer(t.RENDERBUFFER,null),w.depthBuffer)k.__webglDepthRenderbuffer=t.createRenderbuffer(),_t(k.__webglDepthRenderbuffer,w,!0);n.bindFramebuffer(t.FRAMEBUFFER,null)}}if(yt){n.bindTexture(t.TEXTURE_CUBE_MAP,it.__webglTexture),Jt(t.TEXTURE_CUBE_MAP,M);for(let lt=0;lt<6;lt++)if(M.mipmaps&&M.mipmaps.length>0)for(let ut=0;ut<M.mipmaps.length;ut++)J(k.__webglFramebuffer[lt][ut],w,M,t.COLOR_ATTACHMENT0,t.TEXTURE_CUBE_MAP_POSITIVE_X+lt,ut);else J(k.__webglFramebuffer[lt],w,M,t.COLOR_ATTACHMENT0,t.TEXTURE_CUBE_MAP_POSITIVE_X+lt,0);if(p(M))b(t.TEXTURE_CUBE_MAP);n.unbindTexture()}else if(St){for(let lt=0,ut=pt.length;lt<ut;lt++){let wt=pt[lt],Ht=i.get(wt),bt=t.TEXTURE_2D;if(w.isWebGL3DRenderTarget||w.isWebGLArrayRenderTarget)bt=w.isWebGL3DRenderTarget?t.TEXTURE_3D:t.TEXTURE_2D_ARRAY;if(n.bindTexture(bt,Ht.__webglTexture),Jt(bt,wt),J(k.__webglFramebuffer,w,wt,t.COLOR_ATTACHMENT0+lt,bt,0),p(wt))b(bt)}n.unbindTexture()}else{let lt=t.TEXTURE_2D;if(w.isWebGL3DRenderTarget||w.isWebGLArrayRenderTarget)lt=w.isWebGL3DRenderTarget?t.TEXTURE_3D:t.TEXTURE_2D_ARRAY;if(n.bindTexture(lt,it.__webglTexture),Jt(lt,M),M.mipmaps&&M.mipmaps.length>0)for(let ut=0;ut<M.mipmaps.length;ut++)J(k.__webglFramebuffer[ut],w,M,t.COLOR_ATTACHMENT0,lt,ut);else J(k.__webglFramebuffer,w,M,t.COLOR_ATTACHMENT0,lt,0);if(p(M))b(lt);n.unbindTexture()}if(w.depthBuffer)H(w)}function tt(w){let M=w.textures;for(let k=0,it=M.length;k<it;k++){let pt=M[k];if(p(pt)){let yt=S(w),St=i.get(pt).__webglTexture;n.bindTexture(yt,St),b(yt),n.unbindTexture()}}}let at=[],mt=[];function gt(w){if(w.samples>0){if(v(w)===!1){let{textures:M,width:k,height:it}=w,pt=t.COLOR_BUFFER_BIT,yt=w.stencilBuffer?t.DEPTH_STENCIL_ATTACHMENT:t.DEPTH_ATTACHMENT,St=i.get(w),lt=M.length>1;if(lt)for(let wt=0;wt<M.length;wt++)n.bindFramebuffer(t.FRAMEBUFFER,St.__webglMultisampledFramebuffer),t.framebufferRenderbuffer(t.FRAMEBUFFER,t.COLOR_ATTACHMENT0+wt,t.RENDERBUFFER,null),n.bindFramebuffer(t.FRAMEBUFFER,St.__webglFramebuffer),t.framebufferTexture2D(t.DRAW_FRAMEBUFFER,t.COLOR_ATTACHMENT0+wt,t.TEXTURE_2D,null,0);n.bindFramebuffer(t.READ_FRAMEBUFFER,St.__webglMultisampledFramebuffer);let ut=w.texture.mipmaps;if(ut&&ut.length>0)n.bindFramebuffer(t.DRAW_FRAMEBUFFER,St.__webglFramebuffer[0]);else n.bindFramebuffer(t.DRAW_FRAMEBUFFER,St.__webglFramebuffer);for(let wt=0;wt<M.length;wt++){if(w.resolveDepthBuffer){if(w.depthBuffer)pt|=t.DEPTH_BUFFER_BIT;if(w.stencilBuffer&&w.resolveStencilBuffer)pt|=t.STENCIL_BUFFER_BIT}if(lt){t.framebufferRenderbuffer(t.READ_FRAMEBUFFER,t.COLOR_ATTACHMENT0,t.RENDERBUFFER,St.__webglColorRenderbuffer[wt]);let Ht=i.get(M[wt]).__webglTexture;t.framebufferTexture2D(t.DRAW_FRAMEBUFFER,t.COLOR_ATTACHMENT0,t.TEXTURE_2D,Ht,0)}if(t.blitFramebuffer(0,0,k,it,0,0,k,it,pt,t.NEAREST),l===!0){if(at.length=0,mt.length=0,at.push(t.COLOR_ATTACHMENT0+wt),w.depthBuffer&&w.storeMultisampledDepthBuffer===!1)at.push(yt),mt.push(yt),t.invalidateFramebuffer(t.DRAW_FRAMEBUFFER,mt);t.invalidateFramebuffer(t.READ_FRAMEBUFFER,at)}}if(n.bindFramebuffer(t.READ_FRAMEBUFFER,null),n.bindFramebuffer(t.DRAW_FRAMEBUFFER,null),lt)for(let wt=0;wt<M.length;wt++){n.bindFramebuffer(t.FRAMEBUFFER,St.__webglMultisampledFramebuffer),t.framebufferRenderbuffer(t.FRAMEBUFFER,t.COLOR_ATTACHMENT0+wt,t.RENDERBUFFER,St.__webglColorRenderbuffer[wt]);let Ht=i.get(M[wt]).__webglTexture;n.bindFramebuffer(t.FRAMEBUFFER,St.__webglFramebuffer),t.framebufferTexture2D(t.DRAW_FRAMEBUFFER,t.COLOR_ATTACHMENT0+wt,t.TEXTURE_2D,Ht,0)}n.bindFramebuffer(t.DRAW_FRAMEBUFFER,St.__webglMultisampledFramebuffer)}else if(w.depthBuffer&&w.storeMultisampledDepthBuffer===!1&&l){let M=w.stencilBuffer?t.DEPTH_STENCIL_ATTACHMENT:t.DEPTH_ATTACHMENT;t.invalidateFramebuffer(t.DRAW_FRAMEBUFFER,[M])}}}function X(w){return Math.min(s.maxSamples,w.samples)}function v(w){let M=i.get(w);return w.samples>0&&e.has("WEBGL_multisampled_render_to_texture")===!0&&M.__useRenderToTexture!==!1}function Lt(w){let M=a.render.frame;if(h.get(w)!==M)h.set(w,M),w.update()}function zt(w,M){let{colorSpace:k,format:it,type:pt}=w;if(w.isCompressedTexture===!0||w.isVideoTexture===!0)return M;if(k!==Sc&&k!==Gi)if(re.getTransfer(k)===xe){if(it!==Dn||pt!==Tn)Mt("WebGLTextures: sRGB encoded textures have to use RGBAFormat and UnsignedByteType.")}else Wt("WebGLTextures: Unsupported texture color space:",k);return M}function $t(w){if(typeof HTMLImageElement<"u"&&w instanceof HTMLImageElement)c.width=w.naturalWidth||w.width,c.height=w.naturalHeight||w.height;else if(typeof VideoFrame<"u"&&w instanceof VideoFrame)c.width=w.displayWidth,c.height=w.displayHeight;else c.width=w.width,c.height=w.height;return c}this.allocateTextureUnit=nt,this.resetTextureUnits=ot,this.getTextureUnits=z,this.setTextureUnits=et,this.setTexture2D=ft,this.setTexture2DArray=Q,this.setTexture3D=ct,this.setTextureCube=dt,this.rebindTextures=U,this.setupRenderTarget=Y,this.updateRenderTargetMipmap=tt,this.updateMultisampleRenderTarget=gt,this.setupDepthRenderbuffer=H,this.setupFrameBufferTexture=J,this.useMultisampledRTT=v,this.isReversedDepthBuffer=function(){return n.buffers.depth.getReversed()}}function sM(t,e){function n(i,s=Gi){let r,a=re.getTransfer(s);if(i===Tn)return t.UNSIGNED_BYTE;if(i===zl)return t.UNSIGNED_SHORT_4_4_4_4;if(i===kl)return t.UNSIGNED_SHORT_5_5_5_1;if(i===Bd)return t.UNSIGNED_INT_5_9_9_9_REV;if(i===zd)return t.UNSIGNED_INT_10F_11F_11F_REV;if(i===Fd)return t.BYTE;if(i===Od)return t.SHORT;if(i===lr)return t.UNSIGNED_SHORT;if(i===Bl)return t.INT;if(i===mi)return t.UNSIGNED_INT;if(i===$n)return t.FLOAT;if(i===Nn)return t.HALF_FLOAT;if(i===kd)return t.ALPHA;if(i===Gd)return t.RGB;if(i===Dn)return t.RGBA;if(i===Bi)return t.DEPTH_COMPONENT;if(i===zi)return t.DEPTH_STENCIL;if(i===Hd)return t.RED;if(i===Gl)return t.RED_INTEGER;if(i===ki)return t.RG;if(i===Hl)return t.RG_INTEGER;if(i===Vl)return t.RGBA_INTEGER;if(i===Ua||i===Fa||i===Oa||i===Ba)if(a===xe)if(r=e.get("WEBGL_compressed_texture_s3tc_srgb"),r!==null){if(i===Ua)return r.COMPRESSED_SRGB_S3TC_DXT1_EXT;if(i===Fa)return r.COMPRESSED_SRGB_ALPHA_S3TC_DXT1_EXT;if(i===Oa)return r.COMPRESSED_SRGB_ALPHA_S3TC_DXT3_EXT;if(i===Ba)return r.COMPRESSED_SRGB_ALPHA_S3TC_DXT5_EXT}else return null;else if(r=e.get("WEBGL_compressed_texture_s3tc"),r!==null){if(i===Ua)return r.COMPRESSED_RGB_S3TC_DXT1_EXT;if(i===Fa)return r.COMPRESSED_RGBA_S3TC_DXT1_EXT;if(i===Oa)return r.COMPRESSED_RGBA_S3TC_DXT3_EXT;if(i===Ba)return r.COMPRESSED_RGBA_S3TC_DXT5_EXT}else return null;if(i===Wl||i===Xl||i===ql||i===Yl)if(r=e.get("WEBGL_compressed_texture_pvrtc"),r!==null){if(i===Wl)return r.COMPRESSED_RGB_PVRTC_4BPPV1_IMG;if(i===Xl)return r.COMPRESSED_RGB_PVRTC_2BPPV1_IMG;if(i===ql)return r.COMPRESSED_RGBA_PVRTC_4BPPV1_IMG;if(i===Yl)return r.COMPRESSED_RGBA_PVRTC_2BPPV1_IMG}else return null;if(i===Zl||i===Jl||i===$l||i===Kl||i===jl||i===za||i===Ql)if(r=e.get("WEBGL_compressed_texture_etc"),r!==null){if(i===Zl||i===Jl)return a===xe?r.COMPRESSED_SRGB8_ETC2:r.COMPRESSED_RGB8_ETC2;if(i===$l)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ETC2_EAC:r.COMPRESSED_RGBA8_ETC2_EAC;if(i===Kl)return r.COMPRESSED_R11_EAC;if(i===jl)return r.COMPRESSED_SIGNED_R11_EAC;if(i===za)return r.COMPRESSED_RG11_EAC;if(i===Ql)return r.COMPRESSED_SIGNED_RG11_EAC}else return null;if(i===tc||i===ec||i===nc||i===ic||i===sc||i===rc||i===ac||i===oc||i===lc||i===cc||i===hc||i===uc||i===dc||i===fc)if(r=e.get("WEBGL_compressed_texture_astc"),r!==null){if(i===tc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR:r.COMPRESSED_RGBA_ASTC_4x4_KHR;if(i===ec)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR:r.COMPRESSED_RGBA_ASTC_5x4_KHR;if(i===nc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR:r.COMPRESSED_RGBA_ASTC_5x5_KHR;if(i===ic)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR:r.COMPRESSED_RGBA_ASTC_6x5_KHR;if(i===sc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR:r.COMPRESSED_RGBA_ASTC_6x6_KHR;if(i===rc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR:r.COMPRESSED_RGBA_ASTC_8x5_KHR;if(i===ac)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR:r.COMPRESSED_RGBA_ASTC_8x6_KHR;if(i===oc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR:r.COMPRESSED_RGBA_ASTC_8x8_KHR;if(i===lc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR:r.COMPRESSED_RGBA_ASTC_10x5_KHR;if(i===cc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR:r.COMPRESSED_RGBA_ASTC_10x6_KHR;if(i===hc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR:r.COMPRESSED_RGBA_ASTC_10x8_KHR;if(i===uc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR:r.COMPRESSED_RGBA_ASTC_10x10_KHR;if(i===dc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR:r.COMPRESSED_RGBA_ASTC_12x10_KHR;if(i===fc)return a===xe?r.COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR:r.COMPRESSED_RGBA_ASTC_12x12_KHR}else return null;if(i===pc||i===mc||i===gc)if(r=e.get("EXT_texture_compression_bptc"),r!==null){if(i===pc)return a===xe?r.COMPRESSED_SRGB_ALPHA_BPTC_UNORM_EXT:r.COMPRESSED_RGBA_BPTC_UNORM_EXT;if(i===mc)return r.COMPRESSED_RGB_BPTC_SIGNED_FLOAT_EXT;if(i===gc)return r.COMPRESSED_RGB_BPTC_UNSIGNED_FLOAT_EXT}else return null;if(i===_c||i===xc||i===ka||i===yc)if(r=e.get("EXT_texture_compression_rgtc"),r!==null){if(i===_c)return r.COMPRESSED_RED_RGTC1_EXT;if(i===xc)return r.COMPRESSED_SIGNED_RED_RGTC1_EXT;if(i===ka)return r.COMPRESSED_RED_GREEN_RGTC2_EXT;if(i===yc)return r.COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT}else return null;if(i===As)return t.UNSIGNED_INT_24_8;return t[i]!==void 0?t[i]:null}return{convert:n}}var rM=`
void main() {

	gl_Position = vec4( position, 1.0 );

}`,aM=`
uniform sampler2DArray depthColor;
uniform float depthWidth;
uniform float depthHeight;

void main() {

	vec2 coord = vec2( gl_FragCoord.x / depthWidth, gl_FragCoord.y / depthHeight );

	if ( coord.x >= 1.0 ) {

		gl_FragDepth = texture( depthColor, vec3( coord.x - 1.0, coord.y, 1 ) ).r;

	} else {

		gl_FragDepth = texture( depthColor, vec3( coord.x, coord.y, 0 ) ).r;

	}

}`;class tm{constructor(){this.texture=null,this.mesh=null,this.depthNear=0,this.depthFar=0}init(t,e){if(this.texture===null){let n=new eo(t.texture);if(t.depthNear!==e.depthNear||t.depthFar!==e.depthFar)this.depthNear=t.depthNear,this.depthFar=t.depthFar;this.texture=n}}getMesh(t){if(this.texture!==null){if(this.mesh===null){let e=t.cameras[0].viewport,n=new ln({vertexShader:rM,fragmentShader:aM,uniforms:{depthColor:{value:this.texture},depthWidth:{value:e.z},depthHeight:{value:e.w}}});this.mesh=new pe(new Is(20,20),n)}}return this.mesh}reset(){this.texture=null,this.mesh=null}getDepthTexture(){return this.texture}}class em extends pn{constructor(t,e){super();let n=this,i=null,s=1,r=null,a="local-floor",o=1,l=null,c=null,h=null,d=null,u=null,f=null,g=typeof XRWebGLBinding<"u",_=new tm,m={},p=e.getContextAttributes(),b=null,S=null,x=[],T=[],I=new rt,R=null,y=null,A=new De;A.viewport=new fe;let Z=new De;Z.viewport=new fe;let D=[A,Z],q=new uh,ot=null,z=null;this.cameraAutoUpdate=!0,this.enabled=!1,this.isPresenting=!1,this.getController=function(F){let O=x[F];if(O===void 0)O=new dr,x[F]=O;return O.getTargetRaySpace()},this.getControllerGrip=function(F){let O=x[F];if(O===void 0)O=new dr,x[F]=O;return O.getGripSpace()},this.getHand=function(F){let O=x[F];if(O===void 0)O=new dr,x[F]=O;return O.getHandSpace()};function et(F){let O=T.indexOf(F.inputSource);if(O===-1)return;let N=x[O];if(N!==void 0)N.update(F.inputSource,F.frame,l||r),N.dispatchEvent({type:F.type,data:F.inputSource})}function nt(){i.removeEventListener("select",et),i.removeEventListener("selectstart",et),i.removeEventListener("selectend",et),i.removeEventListener("squeeze",et),i.removeEventListener("squeezestart",et),i.removeEventListener("squeezeend",et),i.removeEventListener("end",nt),i.removeEventListener("inputsourceschange",B);for(let F=0;F<x.length;F++){let O=T[F];if(O===null)continue;T[F]=null,x[F].disconnect(O)}ot=null,z=null,_.reset();for(let F in m)delete m[F];if(t.setRenderTarget(b),u=null,d=null,h=null,i=null,S=null,Jt.stop(),n.isPresenting=!1,t.setPixelRatio(R),t.setSize(I.width,I.height,!1),y!==null){let F=y.camera;F.fov=y.fov,F.zoom=y.zoom,F.updateProjectionMatrix(),y=null}n.dispatchEvent({type:"sessionend"})}this.setFramebufferScaleFactor=function(F){if(s=F,n.isPresenting===!0)Mt("WebXRManager: Cannot change framebuffer scale while presenting.")},this.setReferenceSpaceType=function(F){if(a=F,n.isPresenting===!0)Mt("WebXRManager: Cannot change reference space type while presenting.")},this.getReferenceSpace=function(){return l||r},this.setReferenceSpace=function(F){l=F},this.getBaseLayer=function(){return d!==null?d:u},this.getBinding=function(){if(h===null&&g)h=new XRWebGLBinding(i,e);return h},this.getFrame=function(){return f},this.getSession=function(){return i},this.setSession=async function(F){if(i=F,i!==null){if(b=t.getRenderTarget(),i.addEventListener("select",et),i.addEventListener("selectstart",et),i.addEventListener("selectend",et),i.addEventListener("squeeze",et),i.addEventListener("squeezestart",et),i.addEventListener("squeezeend",et),i.addEventListener("end",nt),i.addEventListener("inputsourceschange",B),p.xrCompatible!==!0)await e.makeXRCompatible();if(R=t.getPixelRatio(),t.getSize(I),!(g&&("createProjectionLayer"in XRWebGLBinding.prototype))){let N={antialias:p.antialias,alpha:!0,depth:p.depth,stencil:p.stencil,framebufferScaleFactor:s};u=new XRWebGLLayer(i,e,N),i.updateRenderState({baseLayer:u}),t.setPixelRatio(1),t.setSize(u.framebufferWidth,u.framebufferHeight,!1),S=new Je(u.framebufferWidth,u.framebufferHeight,{format:Dn,type:Tn,colorSpace:t.outputColorSpace,stencilBuffer:p.stencil,resolveDepthBuffer:u.ignoreDepthValues===!1,resolveStencilBuffer:u.ignoreDepthValues===!1,storeMultisampledDepthBuffer:u.ignoreDepthValues===!1,storeMultisampledStencilBuffer:u.ignoreDepthValues===!1})}else{let N=null,P=null,G=null;if(p.depth)G=p.stencil?e.DEPTH24_STENCIL8:e.DEPTH_COMPONENT24,N=p.stencil?zi:Bi,P=p.stencil?As:mi;let J={colorFormat:e.RGBA8,depthFormat:G,scaleFactor:s};h=this.getBinding(),d=h.createProjectionLayer(J),i.updateRenderState({layers:[d]}),t.setPixelRatio(1),t.setSize(d.textureWidth,d.textureHeight,!1),S=new Je(d.textureWidth,d.textureHeight,{format:Dn,type:Tn,depthTexture:new Hi(d.textureWidth,d.textureHeight,P,void 0,void 0,void 0,void 0,void 0,void 0,N),stencilBuffer:p.stencil,colorSpace:t.outputColorSpace,samples:p.antialias?4:0,resolveDepthBuffer:d.ignoreDepthValues===!1,resolveStencilBuffer:d.ignoreDepthValues===!1,storeMultisampledDepthBuffer:d.ignoreDepthValues===!1,storeMultisampledStencilBuffer:d.ignoreDepthValues===!1})}S.isXRRenderTarget=!0,this.setFoveation(o),l=null,r=await i.requestReferenceSpace(a),Jt.setContext(i),Jt.start(),n.isPresenting=!0,n.dispatchEvent({type:"sessionstart"})}},this.getEnvironmentBlendMode=function(){if(i!==null)return i.environmentBlendMode},this.getDepthTexture=function(){return _.getDepthTexture()};function B(F){for(let O=0;O<F.removed.length;O++){let N=F.removed[O],P=T.indexOf(N);if(P>=0)T[P]=null,x[P].disconnect(N)}for(let O=0;O<F.added.length;O++){let N=F.added[O],P=T.indexOf(N);if(P===-1){for(let J=0;J<x.length;J++)if(J>=T.length){T.push(N),P=J;break}else if(T[J]===null){T[J]=N,P=J;break}if(P===-1)break}let G=x[P];if(G)G.connect(N)}}let ft=new L,Q=new L;function ct(F,O,N){ft.setFromMatrixPosition(O.matrixWorld),Q.setFromMatrixPosition(N.matrixWorld);let P=ft.distanceTo(Q),G=O.projectionMatrix.elements,J=N.projectionMatrix.elements,_t=G[14]/(G[10]-1),C=G[14]/(G[10]+1),H=(G[9]+1)/G[5],U=(G[9]-1)/G[5],Y=(G[8]-1)/G[0],tt=(J[8]+1)/J[0],at=_t*Y,mt=_t*tt,gt=P/(-Y+tt),X=gt*-Y;if(O.matrixWorld.decompose(F.position,F.quaternion,F.scale),F.translateX(X),F.translateZ(gt),F.matrixWorld.compose(F.position,F.quaternion,F.scale),F.matrixWorldInverse.copy(F.matrixWorld).invert(),G[10]===-1)F.projectionMatrix.copy(O.projectionMatrix),F.projectionMatrixInverse.copy(O.projectionMatrixInverse);else{let v=_t+gt,Lt=C+gt,zt=at-X,$t=mt+(P-X),w=H*C/Lt*v,M=U*C/Lt*v;F.projectionMatrix.makePerspective(zt,$t,w,M,v,Lt),F.projectionMatrixInverse.copy(F.projectionMatrix).invert()}}function dt(F,O){if(O===null)F.matrixWorld.copy(F.matrix);else F.matrixWorld.multiplyMatrices(O.matrixWorld,F.matrix);F.matrixWorldInverse.copy(F.matrixWorld).invert()}this.updateCamera=function(F){if(i===null)return;let{near:O,far:N}=F;if(_.texture!==null){if(_.depthNear>0)O=_.depthNear;if(_.depthFar>0)N=_.depthFar}if(q.near=Z.near=A.near=O,q.far=Z.far=A.far=N,ot!==q.near||z!==q.far)i.updateRenderState({depthNear:q.near,depthFar:q.far}),ot=q.near,z=q.far;q.layers.mask=F.layers.mask|6,A.layers.mask=q.layers.mask&-5,Z.layers.mask=q.layers.mask&-3;let P=F.parent,G=q.cameras;dt(q,P);for(let J=0;J<G.length;J++)dt(G[J],P);if(G.length===2)ct(q,A,Z);else q.projectionMatrix.copy(A.projectionMatrix);if(y===null&&F.isPerspectiveCamera)y={camera:F,fov:F.fov,zoom:F.zoom};Bt(F,q,P)};function Bt(F,O,N){if(N===null)F.matrix.copy(O.matrixWorld);else F.matrix.copy(N.matrixWorld),F.matrix.invert(),F.matrix.multiply(O.matrixWorld);if(F.matrix.decompose(F.position,F.quaternion,F.scale),F.updateMatrixWorld(!0),F.projectionMatrix.copy(O.projectionMatrix),F.projectionMatrixInverse.copy(O.projectionMatrixInverse),F.isPerspectiveCamera)F.fov=Di*2*Math.atan(1/F.projectionMatrix.elements[5]),F.zoom=1}this.getCamera=function(){return q},this.getFoveation=function(){if(d===null&&u===null)return;return o},this.setFoveation=function(F){if(o=F,d!==null)d.fixedFoveation=F;if(u!==null&&u.fixedFoveation!==void 0)u.fixedFoveation=F},this.hasDepthSensing=function(){return _.texture!==null},this.getDepthSensingMesh=function(){return _.getMesh(q)},this.getCameraTexture=function(F){return m[F]};let Dt=null;function ee(F,O){if(c=O.getViewerPose(l||r),f=O,c!==null){let N=c.views;if(u!==null)t.setRenderTargetFramebuffer(S,u.framebuffer),t.setRenderTarget(S);let P=!1;if(N.length!==q.cameras.length)q.cameras.length=0,P=!0;for(let C=0;C<N.length;C++){let H=N[C],U=null;if(u!==null)U=u.getViewport(H);else{let tt=h.getViewSubImage(d,H);if(U=tt.viewport,C===0)t.setRenderTargetTextures(S,tt.colorTexture,tt.depthStencilTexture),t.setRenderTarget(S)}let Y=D[C];if(Y===void 0)Y=new De,Y.layers.enable(C),Y.viewport=new fe,D[C]=Y;if(Y.matrix.fromArray(H.transform.matrix),Y.matrix.decompose(Y.position,Y.quaternion,Y.scale),Y.projectionMatrix.fromArray(H.projectionMatrix),Y.projectionMatrixInverse.copy(Y.projectionMatrix).invert(),Y.viewport.set(U.x,U.y,U.width,U.height),C===0)q.matrix.copy(Y.matrix),q.matrix.decompose(q.position,q.quaternion,q.scale);if(P===!0)q.cameras.push(Y)}let G=i.enabledFeatures;if(G&&G.includes("depth-sensing")&&i.depthUsage=="gpu-optimized"&&g){h=n.getBinding();let C=h.getDepthInformation(N[0]);if(C&&C.isValid&&C.texture)_.init(C,i.renderState)}if(G&&G.includes("camera-access")&&g){t.state.unbindTexture(),h=n.getBinding();for(let C=0;C<N.length;C++){let H=N[C].camera;if(H){let U=m[H];if(!U)U=new eo,m[H]=U;let Y=h.getCameraImage(H);U.sourceTexture=Y}}}}for(let N=0;N<x.length;N++){let P=T[N],G=x[N];if(P!==null&&G!==void 0)G.update(P,O,l||r)}if(Dt)Dt(F,O);if(O.detectedPlanes)n.dispatchEvent({type:"planesdetected",data:O});f=null}let Jt=new Vp;Jt.setAnimationLoop(ee),this.setAnimationLoop=function(F){Dt=F},this.dispose=function(){}}}var oM=new Kt,nm=new Xt;nm.set(-1,0,0,0,1,0,0,0,1);function lM(t,e){function n(m,p){if(m.matrixAutoUpdate===!0)m.updateMatrix();p.value.copy(m.matrix)}function i(m,p){if(p.color.getRGB(m.fogColor.value,Xc(t)),p.isFog)m.fogNear.value=p.near,m.fogFar.value=p.far;else if(p.isFogExp2)m.fogDensity.value=p.density}function s(m,p,b,S,x){if(p.isNodeMaterial)p.uniformsNeedUpdate=!1;else if(p.isMeshBasicMaterial)r(m,p);else if(p.isMeshLambertMaterial){if(r(m,p),p.envMap)m.envMapIntensity.value=p.envMapIntensity}else if(p.isMeshToonMaterial)r(m,p),d(m,p);else if(p.isMeshPhongMaterial){if(r(m,p),h(m,p),p.envMap)m.envMapIntensity.value=p.envMapIntensity}else if(p.isMeshStandardMaterial){if(r(m,p),u(m,p),p.isMeshPhysicalMaterial)f(m,p,x)}else if(p.isMeshMatcapMaterial)r(m,p),g(m,p);else if(p.isMeshDepthMaterial)r(m,p);else if(p.isMeshDistanceMaterial)r(m,p),_(m,p);else if(p.isMeshNormalMaterial)r(m,p);else if(p.isLineBasicMaterial){if(a(m,p),p.isLineDashedMaterial)o(m,p)}else if(p.isPointsMaterial)l(m,p,b,S);else if(p.isSpriteMaterial)c(m,p);else if(p.isShadowMaterial)m.color.value.copy(p.color),m.opacity.value=p.opacity;else if(p.isShaderMaterial)p.uniformsNeedUpdate=!1}function r(m,p){if(m.opacity.value=p.opacity,p.color)m.diffuse.value.copy(p.color);if(p.emissive)m.emissive.value.copy(p.emissive).multiplyScalar(p.emissiveIntensity);if(p.map)m.map.value=p.map,n(p.map,m.mapTransform);if(p.alphaMap)m.alphaMap.value=p.alphaMap,n(p.alphaMap,m.alphaMapTransform);if(p.bumpMap){if(m.bumpMap.value=p.bumpMap,n(p.bumpMap,m.bumpMapTransform),m.bumpScale.value=p.bumpScale,p.side===Ve)m.bumpScale.value*=-1}if(p.normalMap){if(m.normalMap.value=p.normalMap,n(p.normalMap,m.normalMapTransform),m.normalScale.value.copy(p.normalScale),p.side===Ve)m.normalScale.value.negate()}if(p.displacementMap)m.displacementMap.value=p.displacementMap,n(p.displacementMap,m.displacementMapTransform),m.displacementScale.value=p.displacementScale,m.displacementBias.value=p.displacementBias;if(p.emissiveMap)m.emissiveMap.value=p.emissiveMap,n(p.emissiveMap,m.emissiveMapTransform);if(p.specularMap)m.specularMap.value=p.specularMap,n(p.specularMap,m.specularMapTransform);if(p.alphaTest>0)m.alphaTest.value=p.alphaTest;let b=e.get(p),{envMap:S,envMapRotation:x}=b;if(S){if(m.envMap.value=S,m.envMapRotation.value.setFromMatrix4(oM.makeRotationFromEuler(x)).transpose(),S.isCubeTexture&&S.isRenderTargetTexture===!1)m.envMapRotation.value.premultiply(nm);m.reflectivity.value=p.reflectivity,m.ior.value=p.ior,m.refractionRatio.value=p.refractionRatio}if(p.lightMap)m.lightMap.value=p.lightMap,m.lightMapIntensity.value=p.lightMapIntensity,n(p.lightMap,m.lightMapTransform);if(p.aoMap)m.aoMap.value=p.aoMap,m.aoMapIntensity.value=p.aoMapIntensity,n(p.aoMap,m.aoMapTransform)}function a(m,p){if(m.diffuse.value.copy(p.color),m.opacity.value=p.opacity,p.map)m.map.value=p.map,n(p.map,m.mapTransform)}function o(m,p){m.dashSize.value=p.dashSize,m.totalSize.value=p.dashSize+p.gapSize,m.scale.value=p.scale}function l(m,p,b,S){if(m.diffuse.value.copy(p.color),m.opacity.value=p.opacity,m.size.value=p.size*b,m.scale.value=S*0.5,p.map)m.map.value=p.map,n(p.map,m.uvTransform);if(p.alphaMap)m.alphaMap.value=p.alphaMap,n(p.alphaMap,m.alphaMapTransform);if(p.alphaTest>0)m.alphaTest.value=p.alphaTest}function c(m,p){if(m.diffuse.value.copy(p.color),m.opacity.value=p.opacity,m.rotation.value=p.rotation,p.map)m.map.value=p.map,n(p.map,m.mapTransform);if(p.alphaMap)m.alphaMap.value=p.alphaMap,n(p.alphaMap,m.alphaMapTransform);if(p.alphaTest>0)m.alphaTest.value=p.alphaTest}function h(m,p){m.specular.value.copy(p.specular),m.shininess.value=Math.max(p.shininess,0.0001)}function d(m,p){if(p.gradientMap)m.gradientMap.value=p.gradientMap}function u(m,p){if(m.metalness.value=p.metalness,p.metalnessMap)m.metalnessMap.value=p.metalnessMap,n(p.metalnessMap,m.metalnessMapTransform);if(m.roughness.value=p.roughness,p.roughnessMap)m.roughnessMap.value=p.roughnessMap,n(p.roughnessMap,m.roughnessMapTransform);if(p.envMap)m.envMapIntensity.value=p.envMapIntensity}function f(m,p,b){if(m.ior.value=p.ior,p.sheen>0){if(m.sheenColor.value.copy(p.sheenColor).multiplyScalar(p.sheen),m.sheenRoughness.value=p.sheenRoughness,p.sheenColorMap)m.sheenColorMap.value=p.sheenColorMap,n(p.sheenColorMap,m.sheenColorMapTransform);if(p.sheenRoughnessMap)m.sheenRoughnessMap.value=p.sheenRoughnessMap,n(p.sheenRoughnessMap,m.sheenRoughnessMapTransform)}if(p.clearcoat>0){if(m.clearcoat.value=p.clearcoat,m.clearcoatRoughness.value=p.clearcoatRoughness,p.clearcoatMap)m.clearcoatMap.value=p.clearcoatMap,n(p.clearcoatMap,m.clearcoatMapTransform);if(p.clearcoatRoughnessMap)m.clearcoatRoughnessMap.value=p.clearcoatRoughnessMap,n(p.clearcoatRoughnessMap,m.clearcoatRoughnessMapTransform);if(p.clearcoatNormalMap){if(m.clearcoatNormalMap.value=p.clearcoatNormalMap,n(p.clearcoatNormalMap,m.clearcoatNormalMapTransform),m.clearcoatNormalScale.value.copy(p.clearcoatNormalScale),p.side===Ve)m.clearcoatNormalScale.value.negate()}}if(p.dispersion>0)m.dispersion.value=p.dispersion;if(p.retroreflectivity>0)m.retroreflectivity.value=p.retroreflectivity;if(p.iridescence>0){if(m.iridescence.value=p.iridescence,m.iridescenceIOR.value=p.iridescenceIOR,m.iridescenceThicknessMinimum.value=p.iridescenceThicknessRange[0],m.iridescenceThicknessMaximum.value=p.iridescenceThicknessRange[1],p.iridescenceMap)m.iridescenceMap.value=p.iridescenceMap,n(p.iridescenceMap,m.iridescenceMapTransform);if(p.iridescenceThicknessMap)m.iridescenceThicknessMap.value=p.iridescenceThicknessMap,n(p.iridescenceThicknessMap,m.iridescenceThicknessMapTransform)}if(p.transmission>0){if(m.transmission.value=p.transmission,m.transmissionSamplerMap.value=b.texture,m.transmissionSamplerSize.value.set(b.width,b.height),p.transmissionMap)m.transmissionMap.value=p.transmissionMap,n(p.transmissionMap,m.transmissionMapTransform);if(m.thickness.value=p.thickness,p.thicknessMap)m.thicknessMap.value=p.thicknessMap,n(p.thicknessMap,m.thicknessMapTransform);m.attenuationDistance.value=p.attenuationDistance,m.attenuationColor.value.copy(p.attenuationColor)}if(p.anisotropy>0){if(m.anisotropyVector.value.set(p.anisotropy*Math.cos(p.anisotropyRotation),p.anisotropy*Math.sin(p.anisotropyRotation)),p.anisotropyMap)m.anisotropyMap.value=p.anisotropyMap,n(p.anisotropyMap,m.anisotropyMapTransform)}if(m.specularIntensity.value=p.specularIntensity,m.specularColor.value.copy(p.specularColor),p.specularColorMap)m.specularColorMap.value=p.specularColorMap,n(p.specularColorMap,m.specularColorMapTransform);if(p.specularIntensityMap)m.specularIntensityMap.value=p.specularIntensityMap,n(p.specularIntensityMap,m.specularIntensityMapTransform)}function g(m,p){if(p.matcap)m.matcap.value=p.matcap}function _(m,p){let b=e.get(p).light;m.referencePosition.value.setFromMatrixPosition(b.matrixWorld),m.nearDistance.value=b.shadow.camera.near,m.farDistance.value=b.shadow.camera.far}return{refreshFogUniforms:i,refreshMaterialUniforms:s}}function cM(t,e,n,i){let s={},r={},a=[],o=t.getParameter(t.MAX_UNIFORM_BUFFER_BINDINGS);function l(x,T){let I=T.program;i.uniformBlockBinding(x,I)}function c(x,T){let I=s[x.id];if(I===void 0)m(x),I=h(x),s[x.id]=I,x.addEventListener("dispose",b);let R=T.program;i.updateUBOMapping(x,R);let y=e.render.frame;if(r[x.id]!==y)u(x),r[x.id]=y}function h(x){let T=d();x.__bindingPointIndex=T;let I=t.createBuffer(),{__size:R,usage:y}=x;return t.bindBuffer(t.UNIFORM_BUFFER,I),t.bufferData(t.UNIFORM_BUFFER,R,y),t.bindBuffer(t.UNIFORM_BUFFER,null),t.bindBufferBase(t.UNIFORM_BUFFER,T,I),I}function d(){for(let x=0;x<o;x++)if(a.indexOf(x)===-1)return a.push(x),x;return Wt("WebGLRenderer: Maximum number of simultaneously usable uniforms groups reached."),0}function u(x){let T=s[x.id],{uniforms:I,__cache:R}=x;t.bindBuffer(t.UNIFORM_BUFFER,T);for(let y=0,A=I.length;y<A;y++){let Z=I[y];if(Array.isArray(Z))for(let D=0,q=Z.length;D<q;D++)f(Z[D],y,D,R);else f(Z,y,0,R)}t.bindBuffer(t.UNIFORM_BUFFER,null)}function f(x,T,I,R){if(_(x,T,I,R)===!0){let{__offset:y,value:A}=x;if(Array.isArray(A)){let Z=0;for(let D=0;D<A.length;D++){let q=A[D],ot=p(q);if(g(q,x.__data,Z),typeof q!=="number"&&typeof q!=="boolean"&&!q.isMatrix3&&!ArrayBuffer.isView(q))Z+=ot.storage/Float32Array.BYTES_PER_ELEMENT}}else g(A,x.__data,0);t.bufferSubData(t.UNIFORM_BUFFER,y,x.__data)}}function g(x,T,I){if(typeof x==="number"||typeof x==="boolean")T[0]=x;else if(x.isMatrix3)T[0]=x.elements[0],T[1]=x.elements[1],T[2]=x.elements[2],T[3]=0,T[4]=x.elements[3],T[5]=x.elements[4],T[6]=x.elements[5],T[7]=0,T[8]=x.elements[6],T[9]=x.elements[7],T[10]=x.elements[8],T[11]=0;else if(ArrayBuffer.isView(x))T.set(new x.constructor(x.buffer,x.byteOffset,T.length));else x.toArray(T,I)}function _(x,T,I,R){let y=x.value,A=T+"_"+I;if(R[A]===void 0){if(typeof y==="number"||typeof y==="boolean")R[A]=y;else if(ArrayBuffer.isView(y))R[A]=y.slice();else R[A]=y.clone();return!0}else{let Z=R[A];if(typeof y==="number"||typeof y==="boolean"){if(Z!==y)return R[A]=y,!0}else if(ArrayBuffer.isView(y))return!0;else if(Z.equals(y)===!1)return Z.copy(y),!0}return!1}function m(x){let T=x.uniforms,I=0,R=16;for(let A=0,Z=T.length;A<Z;A++){let D=Array.isArray(T[A])?T[A]:[T[A]];for(let q=0,ot=D.length;q<ot;q++){let z=D[q],et=Array.isArray(z.value)?z.value:[z.value];for(let nt=0,B=et.length;nt<B;nt++){let ft=et[nt],Q=p(ft),ct=I%R,dt=ct%Q.boundary,Bt=ct+dt;if(I+=dt,Bt!==0&&R-Bt<Q.storage)I+=R-Bt;z.__data=new Float32Array(Q.storage/Float32Array.BYTES_PER_ELEMENT),z.__offset=I,I+=Q.storage}}}let y=I%R;if(y>0)I+=R-y;return x.__size=I,x.__cache={},this}function p(x){let T={boundary:0,storage:0};if(typeof x==="number"||typeof x==="boolean")T.boundary=4,T.storage=4;else if(x.isVector2)T.boundary=8,T.storage=8;else if(x.isVector3||x.isColor)T.boundary=16,T.storage=12;else if(x.isVector4)T.boundary=16,T.storage=16;else if(x.isMatrix3)T.boundary=48,T.storage=48;else if(x.isMatrix4)T.boundary=64,T.storage=64;else if(x.isTexture)Mt("WebGLRenderer: Texture samplers can not be part of an uniforms group.");else if(ArrayBuffer.isView(x))T.boundary=16,T.storage=x.byteLength;else Mt("WebGLRenderer: Unsupported uniform value type.",x);return T}function b(x){let T=x.target;T.removeEventListener("dispose",b);let I=a.indexOf(T.__bindingPointIndex);a.splice(I,1),t.deleteBuffer(s[T.id]),delete s[T.id],delete r[T.id]}function S(){for(let x in s)t.deleteBuffer(s[x]);a=[],s={},r={}}return{bind:l,update:c,dispose:S}}var hM=new Uint16Array([12469,15057,12620,14925,13266,14620,13807,14376,14323,13990,14545,13625,14713,13328,14840,12882,14931,12528,14996,12233,15039,11829,15066,11525,15080,11295,15085,10976,15082,10705,15073,10495,13880,14564,13898,14542,13977,14430,14158,14124,14393,13732,14556,13410,14702,12996,14814,12596,14891,12291,14937,11834,14957,11489,14958,11194,14943,10803,14921,10506,14893,10278,14858,9960,14484,14039,14487,14025,14499,13941,14524,13740,14574,13468,14654,13106,14743,12678,14818,12344,14867,11893,14889,11509,14893,11180,14881,10751,14852,10428,14812,10128,14765,9754,14712,9466,14764,13480,14764,13475,14766,13440,14766,13347,14769,13070,14786,12713,14816,12387,14844,11957,14860,11549,14868,11215,14855,10751,14825,10403,14782,10044,14729,9651,14666,9352,14599,9029,14967,12835,14966,12831,14963,12804,14954,12723,14936,12564,14917,12347,14900,11958,14886,11569,14878,11247,14859,10765,14828,10401,14784,10011,14727,9600,14660,9289,14586,8893,14508,8533,15111,12234,15110,12234,15104,12216,15092,12156,15067,12010,15028,11776,14981,11500,14942,11205,14902,10752,14861,10393,14812,9991,14752,9570,14682,9252,14603,8808,14519,8445,14431,8145,15209,11449,15208,11451,15202,11451,15190,11438,15163,11384,15117,11274,15055,10979,14994,10648,14932,10343,14871,9936,14803,9532,14729,9218,14645,8742,14556,8381,14461,8020,14365,7603,15273,10603,15272,10607,15267,10619,15256,10631,15231,10614,15182,10535,15118,10389,15042,10167,14963,9787,14883,9447,14800,9115,14710,8665,14615,8318,14514,7911,14411,7507,14279,7198,15314,9675,15313,9683,15309,9712,15298,9759,15277,9797,15229,9773,15166,9668,15084,9487,14995,9274,14898,8910,14800,8539,14697,8234,14590,7790,14479,7409,14367,7067,14178,6621,15337,8619,15337,8631,15333,8677,15325,8769,15305,8871,15264,8940,15202,8909,15119,8775,15022,8565,14916,8328,14804,8009,14688,7614,14569,7287,14448,6888,14321,6483,14088,6171,15350,7402,15350,7419,15347,7480,15340,7613,15322,7804,15287,7973,15229,8057,15148,8012,15046,7846,14933,7611,14810,7357,14682,7069,14552,6656,14421,6316,14251,5948,14007,5528,15356,5942,15356,5977,15353,6119,15348,6294,15332,6551,15302,6824,15249,7044,15171,7122,15070,7050,14949,6861,14818,6611,14679,6349,14538,6067,14398,5651,14189,5311,13935,4958,15359,4123,15359,4153,15356,4296,15353,4646,15338,5160,15311,5508,15263,5829,15188,6042,15088,6094,14966,6001,14826,5796,14678,5543,14527,5287,14377,4985,14133,4586,13869,4257,15360,1563,15360,1642,15358,2076,15354,2636,15341,3350,15317,4019,15273,4429,15203,4732,15105,4911,14981,4932,14836,4818,14679,4621,14517,4386,14359,4156,14083,3795,13808,3437,15360,122,15360,137,15358,285,15355,636,15344,1274,15322,2177,15281,2765,15215,3223,15120,3451,14995,3569,14846,3567,14681,3466,14511,3305,14344,3121,14037,2800,13753,2467,15360,0,15360,1,15359,21,15355,89,15346,253,15325,479,15287,796,15225,1148,15133,1492,15008,1749,14856,1882,14685,1886,14506,1783,14324,1608,13996,1398,13702,1183]),Fn=null;function uM(){if(Fn===null)Fn=new an(hM,16,16,ki,Nn),Fn.name="DFG_LUT",Fn.minFilter=Qe,Fn.magFilter=Qe,Fn.wrapS=Es,Fn.wrapT=Es,Fn.generateMipmaps=!1,Fn.needsUpdate=!0;return Fn}class dM{constructor(t={}){let{canvas:e=tf(),context:n=null,depth:i=!0,stencil:s=!1,alpha:r=!1,antialias:a=!1,premultipliedAlpha:o=!0,preserveDrawingBuffer:l=!1,powerPreference:c="default",failIfMajorPerformanceCaveat:h=!1,reversedDepthBuffer:d=!1,outputBufferType:u=Tn}=t;this.isWebGLRenderer=!0;let f;if(n!==null){if(typeof WebGLRenderingContext<"u"&&n instanceof WebGLRenderingContext)throw Error("THREE.WebGLRenderer: WebGL 1 is not supported since r163.");f=n.getContextAttributes().alpha}else f=r;let g=u,_=new Set([Vl,Hl,Gl]),m=new Set([Tn,mi,lr,As,zl,kl]),p=new Uint32Array(4),b=new Int32Array(4),S=new L,x=null,T=null,I=[],R=[],y=null;this.domElement=e,this.debug={checkShaderErrors:!0,diagnostics:{keywords:!1},onShaderError:null},this.autoClear=!0,this.autoClearColor=!0,this.autoClearDepth=!0,this.autoClearStencil=!0,this.sortObjects=!0,this.clippingPlanes=[],this.localClippingEnabled=!1,this.toneMapping=bn,this.toneMappingExposure=1,this.transmissionResolutionScale=1;let A=this,Z=!1,D=null,q=null,ot=null,z=null;this._outputColorSpace=Ga;let et=0,nt=0,B=null,ft=-1,Q=null,ct=new fe,dt=new fe,Bt=null,Dt=new Pt(0),ee=0,{width:Jt,height:F}=e,O=1,N=null,P=null,G=new fe(0,0,Jt,F),J=new fe(0,0,Jt,F),_t=!1,C=new ui,H=!1,U=!1,Y=new Kt,tt=new L,at=new fe,mt={background:null,fog:null,environment:null,overrideMaterial:null,isScene:!0},gt=!1;function X(){return B===null?O:1}let v=n;function Lt(E,V){return e.getContext(E,V)}let zt,$t,w,M,k,it,pt,yt,St,lt,ut,wt,Ht,bt,vt,qt,Yt,ce,W,Tt,ht,Et,Ut;try{let E={alpha:!0,depth:i,stencil:s,antialias:a,premultipliedAlpha:o,preserveDrawingBuffer:l,powerPreference:c,failIfMajorPerformanceCaveat:h};if("setAttribute"in e)e.setAttribute("data-engine",`three.js r${ed}`);if(e.addEventListener("webglcontextlost",Qt,!1),e.addEventListener("webglcontextrestored",ye,!1),e.addEventListener("webglcontextcreationerror",he,!1),v===null){if(v=Lt("webgl2",E),v===null)if(Lt("webgl2"))throw Error("THREE.WebGLRenderer: Error creating WebGL context with your selected attributes.");else throw Error("THREE.WebGLRenderer: Error creating WebGL context.")}xt()}catch(E){throw e.removeEventListener("webglcontextlost",Qt,!1),e.removeEventListener("webglcontextrestored",ye,!1),e.removeEventListener("webglcontextcreationerror",he,!1),Wt("WebGLRenderer: "+E.message),E}function xt(){if(zt=new xv(v),zt.init(),ht=new sM(v,zt),$t=new lv(v,zt,t,ht),w=new nM(v,zt),$t.reversedDepthBuffer&&d)w.buffers.depth.setReversed(!0);q=v.createFramebuffer(),ot=v.createFramebuffer(),z=v.createFramebuffer(),M=new Sv(v),k=new HS,it=new iM(v,zt,w,k,$t,ht,M),pt=new _v(A),yt=new b_(v),Et=new av(v,yt),St=new yv(v,yt,M,Et),lt=new bv(v,St,yt,Et,M),ce=new Mv(v,$t,it),vt=new cv(k),ut=new GS(A,pt,zt,$t,Et,vt),wt=new lM(A,k),Ht=new WS,bt=new $S(zt),Yt=new rv(A,pt,w,lt,f,o),qt=new eM(A,lt,$t),Ut=new cM(v,M,$t,w),W=new ov(v,zt,M),Tt=new vv(v,zt,M),M.programs=ut.programs,A.capabilities=$t,A.extensions=zt,A.properties=k,A.renderLists=Ht,A.shadowMap=qt,A.state=w,A.info=M}if(g!==Tn)y=new Ev(g,e.width,e.height,a,i,s);let Rt=new em(A,v);this.xr=Rt,this.getContext=function(){return v},this.getContextAttributes=function(){return v.getContextAttributes()},this.forceContextLoss=function(){let E=zt.get("WEBGL_lose_context");if(E)E.loseContext()},this.forceContextRestore=function(){let E=zt.get("WEBGL_lose_context");if(E)E.restoreContext()},this.getPixelRatio=function(){return O},this.setPixelRatio=function(E){if(E===void 0)return;O=E,this.setSize(Jt,F,!1)},this.getSize=function(E){return E.set(Jt,F)},this.setSize=function(E,V,st=!0){if(Rt.isPresenting){Mt("WebGLRenderer: Can't change size while VR device is presenting.");return}if(Jt=E,F=V,e.width=Math.floor(E*O),e.height=Math.floor(V*O),st===!0)e.style.width=E+"px",e.style.height=V+"px";if(y!==null)y.setSize(e.width,e.height);this.setViewport(0,0,E,V)},this.getDrawingBufferSize=function(E){return E.set(Jt*O,F*O).floor()},this.setDrawingBufferSize=function(E,V,st){Jt=E,F=V,O=st,e.width=Math.floor(E*st),e.height=Math.floor(V*st),this.setViewport(0,0,E,V)},this.setEffects=function(E){if(g===Tn){Wt("WebGLRenderer: setEffects() requires outputBufferType set to HalfFloatType or FloatType.");return}if(E){for(let V=0;V<E.length;V++)if(E[V].isOutputPass===!0){Mt("WebGLRenderer: OutputPass is not needed in setEffects(). Tone mapping and color space conversion are applied automatically.");break}}y.setEffects(E||[])},this.getCurrentViewport=function(E){return E.copy(ct)},this.getViewport=function(E){return E.copy(G)},this.setViewport=function(E,V,st,K){if(E.isVector4)G.set(E.x,E.y,E.z,E.w);else G.set(E,V,st,K);w.viewport(ct.copy(G).multiplyScalar(O).round())},this.getScissor=function(E){return E.copy(J)},this.setScissor=function(E,V,st,K){if(E.isVector4)J.set(E.x,E.y,E.z,E.w);else J.set(E,V,st,K);w.scissor(dt.copy(J).multiplyScalar(O).round())},this.getScissorTest=function(){return _t},this.setScissorTest=function(E){w.setScissorTest(_t=E)},this.setOpaqueSort=function(E){N=E},this.setTransparentSort=function(E){P=E},this.getClearColor=function(E){return E.copy(Yt.getClearColor())},this.setClearColor=function(){Yt.setClearColor(...arguments)},this.getClearAlpha=function(){return Yt.getClearAlpha()},this.setClearAlpha=function(){Yt.setClearAlpha(...arguments)},this.clear=function(E=!0,V=!0,st=!0){let K=0;if(E){let j=!1;if(B!==null){let It=B.texture.format;j=_.has(It)}if(j){let It=B.texture.type,Ot=m.has(It),Ct=Yt.getClearColor(),kt=Yt.getClearAlpha(),{r:Vt,g:ne,b:se}=Ct;if(Ot)p[0]=Vt,p[1]=ne,p[2]=se,p[3]=kt,v.clearBufferuiv(v.COLOR,0,p);else b[0]=Vt,b[1]=ne,b[2]=se,b[3]=kt,v.clearBufferiv(v.COLOR,0,b)}else K|=v.COLOR_BUFFER_BIT}if(V)K|=v.DEPTH_BUFFER_BIT,this.state.buffers.depth.setMask(!0);if(st)K|=v.STENCIL_BUFFER_BIT,this.state.buffers.stencil.setMask(4294967295);if(K!==0)v.clear(K)},this.clearColor=function(){this.clear(!0,!1,!1)},this.clearDepth=function(){this.clear(!1,!0,!1)},this.clearStencil=function(){this.clear(!1,!1,!0)},this.setNodesHandler=function(E){E.setRenderer(this),D=E},this.dispose=function(){e.removeEventListener("webglcontextlost",Qt,!1),e.removeEventListener("webglcontextrestored",ye,!1),e.removeEventListener("webglcontextcreationerror",he,!1),Yt.dispose(),Ht.dispose(),bt.dispose(),k.dispose(),pt.dispose(),lt.dispose(),Et.dispose(),Ut.dispose(),ut.dispose(),Rt.dispose(),Rt.removeEventListener("sessionstart",Uh),Rt.removeEventListener("sessionend",Fh),yi.stop()};function Qt(E){E.preventDefault(),tr("WebGLRenderer: Context Lost."),Z=!0}function ye(){tr("WebGLRenderer: Context Restored."),Z=!1;let E=M.autoReset,V=qt.enabled,st=qt.autoUpdate,K=qt.needsUpdate,j=qt.type;xt(),M.autoReset=E,qt.enabled=V,qt.autoUpdate=st,qt.needsUpdate=K,qt.type=j}function he(E){Wt("WebGLRenderer: A WebGL context could not be created. Reason: ",E.statusMessage)}function wn(E){let V=E.target;V.removeEventListener("dispose",wn),Bn(V)}function Bn(E){um(E),k.remove(E)}function um(E){let V=k.get(E).programs;if(V!==void 0){if(V.forEach(function(st){ut.releaseProgram(st)}),E.isShaderMaterial)ut.releaseShaderCache(E)}}this.renderBufferDirect=function(E,V,st,K,j,It){if(V===null)V=mt;let Ot=j.isMesh&&j.matrixWorld.determinantAffine()<0,Ct=pm(E,V,st,K,j);w.setMaterial(K,Ot);let kt=st.index,Vt=1;if(K.wireframe===!0){if(kt=St.getWireframeAttribute(st),kt===void 0)return;Vt=2}let ne=st.drawRange,se=st.attributes.position,Gt=ne.start*Vt,ue=(ne.start+ne.count)*Vt;if(It!==null)Gt=Math.max(Gt,It.start*Vt),ue=Math.min(ue,(It.start+It.count)*Vt);if(kt!==null)Gt=Math.max(Gt,0),ue=Math.min(ue,kt.count);else if(se!==void 0&&se!==null)Gt=Math.max(Gt,0),ue=Math.min(ue,se.count);let we=ue-Gt;if(we<0||we===1/0)return;Et.setup(j,K,Ct,st,kt);let Se,_e=W;if(kt!==null)Se=yt.get(kt),_e=Tt,_e.setIndex(Se);if(j.isMesh)if(K.wireframe===!0)w.setLineWidth(K.wireframeLinewidth*X()),_e.setMode(v.LINES);else _e.setMode(v.TRIANGLES);else if(j.isLine){let ze=K.linewidth;if(ze===void 0)ze=1;if(w.setLineWidth(ze*X()),j.isLineSegments)_e.setMode(v.LINES);else if(j.isLineLoop)_e.setMode(v.LINE_LOOP);else _e.setMode(v.LINE_STRIP)}else if(j.isPoints)_e.setMode(v.POINTS);else if(j.isSprite)_e.setMode(v.TRIANGLES);if(j.isBatchedMesh)if(!zt.get("WEBGL_multi_draw")){let{_multiDrawStarts:ze,_multiDrawCounts:Ft,_multiDrawCount:qe}=j,oe=kt?yt.get(kt).bytesPerElement:1,hn=k.get(K).currentProgram.getUniforms();for(let Rn=0;Rn<qe;Rn++)hn.setValue(v,"_gl_DrawID",Rn),_e.render(ze[Rn]/oe,Ft[Rn])}else _e.renderMultiDraw(j._multiDrawStarts,j._multiDrawCounts,j._multiDrawCount);else if(j.isInstancedMesh)_e.renderInstances(Gt,we,j.count);else if(st.isInstancedBufferGeometry){let ze=st._maxInstanceCount!==void 0?st._maxInstanceCount:1/0,Ft=Math.min(st.instanceCount,ze);_e.renderInstances(Gt,we,Ft)}else _e.render(Gt,we)};function Dh(E,V,st,K){if(D!==null&&E.isNodeMaterial)D.setObject(K,E);if(H===!0)vt.setState(E,st,!1);if(E.transparent===!0&&E.side===on&&E.forceSinglePass===!1)E.side=Ve,E.needsUpdate=!0,Fr(E,V,K),E.side=Ms,E.needsUpdate=!0,Fr(E,V,K),E.side=on;else Fr(E,V,K)}this.compile=function(E,V,st=null){if(st===null)st=E;if(D!==null)D.renderStart(E,V,st);if(T=bt.get(st),T.init(V),R.push(T),st.traverseVisible(function(j){if(j.isLight&&j.layers.test(V.layers)){if(T.pushLight(j),j.castShadow)T.pushShadow(j)}}),E!==st)E.traverseVisible(function(j){if(j.isLight&&j.layers.test(V.layers)){if(T.pushLight(j),j.castShadow)T.pushShadow(j)}});if(T.setupLights(),D!==null)D.updateLights(T.state.lightsArray);if(U=this.localClippingEnabled,H=vt.init(this.clippingPlanes,U),H===!0)vt.setGlobalState(this.clippingPlanes,V);if(D!==null)qt.render(T.state.shadowsArray,st,V);let K=new Set;if(E.traverse(function(j){if(!(j.isMesh||j.isPoints||j.isLine||j.isSprite))return;let It=j.material;if(It)if(Array.isArray(It))for(let Ot=0;Ot<It.length;Ot++){let Ct=It[Ot];Dh(Ct,st,V,j),K.add(Ct)}else Dh(It,st,V,j),K.add(It)}),T=R.pop(),D!==null)D.renderEnd();return K},this.compileAsync=function(E,V,st=null){let K=this.compile(E,V,st);return new Promise((j)=>{function It(){if(K.forEach(function(Ot){let kt=k.get(Ot).currentProgram;if(kt===void 0||kt.isReady())K.delete(Ot)}),K.size===0){j(E);return}setTimeout(It,10)}if(zt.get("KHR_parallel_shader_compile")!==null)It();else setTimeout(It,10)})};let Oo=null;function dm(E){if(Oo)Oo(E)}function Uh(){yi.stop()}function Fh(){yi.start()}let yi=new Vp;if(yi.setAnimationLoop(dm),typeof self<"u")yi.setContext(self);this.setAnimationLoop=function(E){Oo=E,Rt.setAnimationLoop(E),E===null?yi.stop():yi.start()},Rt.addEventListener("sessionstart",Uh),Rt.addEventListener("sessionend",Fh),this.render=function(E,V){if(V!==void 0&&V.isCamera!==!0){Wt("WebGLRenderer.render: camera is not an instance of THREE.Camera.");return}if(Z===!0)return;if(D!==null)D.renderStart(E,V);let st=Rt.enabled===!0&&Rt.isPresenting===!0,K=y!==null&&(B===null||st)&&y.begin(A,B);if(E.matrixWorldAutoUpdate===!0)E.updateMatrixWorld();if(V.parent===null&&V.matrixWorldAutoUpdate===!0)V.updateMatrixWorld();if(Rt.enabled===!0&&Rt.isPresenting===!0&&(y===null||y.isCompositing()===!1)){if(Rt.cameraAutoUpdate===!0)Rt.updateCamera(V);V=Rt.getCamera()}if(E.isScene===!0)E.onBeforeRender(A,E,V,B);if(T=bt.get(E,R.length),T.init(V),T.state.textureUnits=it.getTextureUnits(),R.push(T),Y.multiplyMatrices(V.projectionMatrix,V.matrixWorldInverse),C.setFromProjectionMatrix(Y,Tc,V.reversedDepth),U=this.localClippingEnabled,H=vt.init(this.clippingPlanes,U),x=Ht.get(E,I.length),x.init(),I.push(x),Rt.enabled===!0&&Rt.isPresenting===!0){let Ot=A.xr.getDepthSensingMesh();if(Ot!==null)Bo(Ot,V,-1/0,A.sortObjects)}if(Bo(E,V,0,A.sortObjects),x.finish(),D!==null)D.updateLights(T.state.lightsArray);if(A.sortObjects===!0)x.sort(N,P);if(gt=Rt.enabled===!1||Rt.isPresenting===!1||Rt.hasDepthSensing()===!1,gt)Yt.addToRenderList(x,E);if(this.info.render.frame++,this.info.autoReset===!0)this.info.reset();if(H===!0)vt.beginShadows();let j=T.state.shadowsArray;if(qt.render(j,E,V),H===!0)vt.endShadows();if((K&&y.hasRenderPass())===!1){let Ot=x.opaque,Ct=x.transmissive;if(T.setupLights(),V.isArrayCamera){let kt=V.cameras;if(Ct.length>0)for(let Vt=0,ne=kt.length;Vt<ne;Vt++){let se=kt[Vt];Bh(Ot,Ct,E,se)}if(gt)Yt.render(E);for(let Vt=0,ne=kt.length;Vt<ne;Vt++){let se=kt[Vt];Oh(x,E,se,se.viewport)}}else{if(Ct.length>0)Bh(Ot,Ct,E,V);if(gt)Yt.render(E);Oh(x,E,V)}}if(B!==null&&nt===0)it.updateMultisampleRenderTarget(B),it.updateRenderTargetMipmap(B);if(K)y.end(A);if(E.isScene===!0)E.onAfterRender(A,E,V);if(Et.resetDefaultState(),ft=-1,Q=null,R.pop(),R.length>0){if(T=R[R.length-1],it.setTextureUnits(T.state.textureUnits),H===!0)vt.setGlobalState(A.clippingPlanes,T.state.camera)}else T=null;if(I.pop(),I.length>0)x=I[I.length-1];else x=null;if(D!==null)D.renderEnd()};function Bo(E,V,st,K){if(E.visible===!1)return;if(E.layers.test(V.layers)){if(E.isGroup)st=E.renderOrder;else if(E.isLOD){if(E.autoUpdate===!0)E.update(V)}else if(E.isLightProbeGrid)T.pushLightProbeGrid(E);else if(E.isLight){if(T.pushLight(E),E.castShadow)T.pushShadow(E)}else if(E.isSprite){if(!E.frustumCulled||E.intersectsFrustum(C)){if(K)at.setFromMatrixPosition(E.matrixWorld).applyMatrix4(Y);let Ot=lt.update(E),Ct=E.material;if(Ct.visible)x.push(E,Ot,Ct,st,at.z,null,V)}}else if(E.isMesh||E.isLine||E.isPoints){if(!E.frustumCulled||E.intersectsFrustum(C)){let Ot=lt.update(E),Ct=E.material;if(K){if(E.boundingSphere!==void 0){if(E.boundingSphere===null)E.computeBoundingSphere();at.copy(E.boundingSphere.center)}else{if(Ot.boundingSphere===null)Ot.computeBoundingSphere();at.copy(Ot.boundingSphere.center)}at.applyMatrix4(E.matrixWorld).applyMatrix4(Y)}if(Array.isArray(Ct)){let kt=Ot.groups;for(let Vt=0,ne=kt.length;Vt<ne;Vt++){let se=kt[Vt],Gt=Ct[se.materialIndex];if(Gt&&Gt.visible)x.push(E,Ot,Gt,st,at.z,se,V)}}else if(Ct.visible)x.push(E,Ot,Ct,st,at.z,null,V)}}}let It=E.children;for(let Ot=0,Ct=It.length;Ot<Ct;Ot++)Bo(It[Ot],V,st,K)}function Oh(E,V,st,K){let{opaque:j,transmissive:It,transparent:Ot}=E;if(T.setupLightsView(st),H===!0)vt.setGlobalState(A.clippingPlanes,st);if(K)w.viewport(ct.copy(K));if(j.length>0)Ur(j,V,st);if(It.length>0)Ur(It,V,st);if(Ot.length>0)Ur(Ot,V,st);w.buffers.depth.setTest(!0),w.buffers.depth.setMask(!0),w.buffers.color.setMask(!0),w.setPolygonOffset(!1)}function Bh(E,V,st,K){if((st.isScene===!0?st.overrideMaterial:null)!==null)return;if(T.state.transmissionRenderTarget[K.id]===void 0){let Gt=zt.has("EXT_color_buffer_half_float")||zt.has("EXT_color_buffer_float");T.state.transmissionRenderTarget[K.id]=new Je(1,1,{generateMipmaps:!0,type:Gt?Nn:Tn,minFilter:Oi,samples:Math.max(4,$t.samples),stencilBuffer:s,resolveDepthBuffer:!1,resolveStencilBuffer:!1,storeMultisampledDepthBuffer:!1,storeMultisampledStencilBuffer:!1,colorSpace:re.workingColorSpace})}let It=T.state.transmissionRenderTarget[K.id],Ot=K.viewport||ct;It.setSize(Ot.z*A.transmissionResolutionScale,Ot.w*A.transmissionResolutionScale);let Ct=A.getRenderTarget(),kt=A.getActiveCubeFace(),Vt=A.getActiveMipmapLevel();if(A.setRenderTarget(It),A.getClearColor(Dt),ee=A.getClearAlpha(),ee<1)A.setClearColor(16777215,0.5);if(A.clear(),gt)Yt.render(st);let ne=A.toneMapping;A.toneMapping=bn;let se=K.viewport;if(K.viewport!==void 0)K.viewport=void 0;if(T.setupLightsView(K),H===!0)vt.setGlobalState(A.clippingPlanes,K);if(Ur(E,st,K),it.updateMultisampleRenderTarget(It),it.updateRenderTargetMipmap(It),zt.has("WEBGL_multisampled_render_to_texture")===!1){let Gt=!1;for(let ue=0,we=V.length;ue<we;ue++){let Se=V[ue],{object:_e,geometry:ze,material:Ft,group:qe}=Se;if(Ft.side===on&&_e.layers.test(K.layers)){let oe=Ft.side;Ft.side=Ve,Ft.needsUpdate=!0,zh(_e,st,K,ze,Ft,qe),Ft.side=oe,Ft.needsUpdate=!0,Gt=!0}}if(Gt===!0)it.updateMultisampleRenderTarget(It),it.updateRenderTargetMipmap(It)}if(A.setRenderTarget(Ct,kt,Vt),A.setClearColor(Dt,ee),se!==void 0)K.viewport=se;A.toneMapping=ne}function Ur(E,V,st){let K=V.isScene===!0?V.overrideMaterial:null;for(let j=0,It=E.length;j<It;j++){let Ot=E[j],{object:Ct,geometry:kt,group:Vt}=Ot,ne=Ot.material;if(ne.allowOverride===!0&&K!==null)ne=K;if(Ct.layers.test(st.layers))zh(Ct,V,st,kt,ne,Vt)}}function zh(E,V,st,K,j,It){if(D!==null&&j.isNodeMaterial)D.setObject(E,j);if(E.onBeforeRender(A,V,st,K,j,It),E.modelViewMatrix.multiplyMatrices(st.matrixWorldInverse,E.matrixWorld),E.normalMatrix.getNormalMatrix(E.modelViewMatrix),j.onBeforeRender(A,V,st,K,E,It),j.transparent===!0&&j.side===on&&j.forceSinglePass===!1)j.side=Ve,j.needsUpdate=!0,A.renderBufferDirect(st,V,K,j,E,It),j.side=Ms,j.needsUpdate=!0,A.renderBufferDirect(st,V,K,j,E,It),j.side=on;else A.renderBufferDirect(st,V,K,j,E,It);E.onAfterRender(A,V,st,K,j,It)}function Fr(E,V,st){if(V.isScene!==!0)V=mt;let K=k.get(E),j=T.state.lights,It=T.state.shadowsArray,Ot=j.state.version,Ct=ut.getParameters(E,j.state,It,V,st,T.state.lightProbeGridArray),kt=ut.getProgramCacheKey(Ct),Vt=K.programs;K.environment=E.isMeshStandardMaterial||E.isMeshLambertMaterial||E.isMeshPhongMaterial?V.environment:null,K.fog=V.fog;let ne=E.isMeshStandardMaterial||E.isMeshLambertMaterial&&!E.envMap||E.isMeshPhongMaterial&&!E.envMap;if(K.envMap=pt.get(E.envMap||K.environment,ne),K.envMapRotation=K.environment!==null&&E.envMap===null?V.environmentRotation:E.envMapRotation,Vt===void 0)E.addEventListener("dispose",wn),Vt=new Map,K.programs=Vt;let se=Vt.get(kt);if(se!==void 0){if(K.currentProgram===se&&K.lightsStateVersion===Ot)return Gh(E,Ct),se}else{if(Ct.uniforms=ut.getUniforms(E),D!==null&&E.isNodeMaterial)D.build(E,st,Ct);E.onBeforeCompile(Ct,A),se=ut.acquireProgram(Ct,kt),Vt.set(kt,se),K.uniforms=Ct.uniforms}let Gt=K.uniforms;if(!E.isShaderMaterial&&!E.isRawShaderMaterial||E.clipping===!0)Gt.clippingPlanes=vt.uniform;if(Gh(E,Ct),K.needsLights=gm(E),K.lightsStateVersion=Ot,K.needsLights)Gt.ambientLightColor.value=j.state.ambient,Gt.lightProbe.value=j.state.probe,Gt.sunLights.value=j.state.sun,Gt.sunLightShadows.value=j.state.sunShadow,Gt.directionalLights.value=j.state.directional,Gt.directionalLightShadows.value=j.state.directionalShadow,Gt.spotLights.value=j.state.spot,Gt.spotLightShadows.value=j.state.spotShadow,Gt.rectAreaLights.value=j.state.rectArea,Gt.ltc_1.value=j.state.rectAreaLTC1,Gt.ltc_2.value=j.state.rectAreaLTC2,Gt.pointLights.value=j.state.point,Gt.pointLightShadows.value=j.state.pointShadow,Gt.hemisphereLights.value=j.state.hemi,Gt.sunShadowMatrix.value=j.state.sunShadowMatrix,Gt.sunShadowCascade.value=j.state.sunShadowCascade,Gt.directionalShadowMatrix.value=j.state.directionalShadowMatrix,Gt.spotLightMatrix.value=j.state.spotLightMatrix,Gt.spotLightMap.value=j.state.spotLightMap,Gt.pointShadowMatrix.value=j.state.pointShadowMatrix;return K.lightProbeGrid=T.state.lightProbeGridArray.length>0,K.currentProgram=se,K.uniformsList=null,se}function kh(E){if(E.uniformsList===null){let V=E.currentProgram.getUniforms();E.uniformsList=Lr.seqWithValue(V.seq,E.uniforms)}return E.uniformsList}function Gh(E,V){let st=k.get(E);st.outputColorSpace=V.outputColorSpace,st.batching=V.batching,st.batchingColor=V.batchingColor,st.instancing=V.instancing,st.instancingColor=V.instancingColor,st.instancingMorph=V.instancingMorph,st.skinning=V.skinning,st.morphTargets=V.morphTargets,st.morphNormals=V.morphNormals,st.morphColors=V.morphColors,st.morphTargetsCount=V.morphTargetsCount,st.numClippingPlanes=V.numClippingPlanes,st.numIntersection=V.numClipIntersection,st.vertexAlphas=V.vertexAlphas,st.vertexTangents=V.vertexTangents,st.toneMapping=V.toneMapping}function fm(E,V){if(E.length===0)return null;if(E.length===1)return E[0].texture!==null?E[0]:null;S.setFromMatrixPosition(V.matrixWorld);for(let st=0,K=E.length;st<K;st++){let j=E[st];if(j.texture!==null&&j.boundingBox.containsPoint(S))return j}return null}function pm(E,V,st,K,j){if(V.isScene!==!0)V=mt;it.resetTextureUnits();let It=V.fog,Ot=K.isMeshStandardMaterial||K.isMeshLambertMaterial||K.isMeshPhongMaterial?V.environment:null,Ct=B===null?A.outputColorSpace:B.isXRRenderTarget===!0?B.texture.colorSpace:re.workingColorSpace,kt=K.isMeshStandardMaterial||K.isMeshLambertMaterial&&!K.envMap||K.isMeshPhongMaterial&&!K.envMap,Vt=pt.get(K.envMap||Ot,kt),ne=K.vertexColors===!0&&!!st.attributes.color&&st.attributes.color.itemSize===4,se=!!st.attributes.tangent&&(!!K.normalMap||K.anisotropy>0),Gt=!!st.morphAttributes.position,ue=!!st.morphAttributes.normal,we=!!st.morphAttributes.color,Se=bn;if(K.toneMapped){if(B===null||B.isXRRenderTarget===!0)Se=A.toneMapping}let _e=st.morphAttributes.position||st.morphAttributes.normal||st.morphAttributes.color,ze=_e!==void 0?_e.length:0,Ft=k.get(K),qe=T.state.lights;if(H===!0){if(U===!0||E!==Q){let ve=E===Q&&K.id===ft;vt.setState(K,E,ve)}}let oe=!1;if(K.version===Ft.__version){if(Ft.needsLights&&Ft.lightsStateVersion!==qe.state.version)oe=!0;else if(Ft.outputColorSpace!==Ct)oe=!0;else if(j.isBatchedMesh&&Ft.batching===!1)oe=!0;else if(!j.isBatchedMesh&&Ft.batching===!0)oe=!0;else if(j.isBatchedMesh&&Ft.batchingColor===!0&&j._colorsTexture===null)oe=!0;else if(j.isBatchedMesh&&Ft.batchingColor===!1&&j._colorsTexture!==null)oe=!0;else if(j.isInstancedMesh&&Ft.instancing===!1)oe=!0;else if(!j.isInstancedMesh&&Ft.instancing===!0)oe=!0;else if(j.isSkinnedMesh&&Ft.skinning===!1)oe=!0;else if(!j.isSkinnedMesh&&Ft.skinning===!0)oe=!0;else if(j.isInstancedMesh&&Ft.instancingColor===!0&&j.instanceColor===null)oe=!0;else if(j.isInstancedMesh&&Ft.instancingColor===!1&&j.instanceColor!==null)oe=!0;else if(j.isInstancedMesh&&Ft.instancingMorph===!0&&j.morphTexture===null)oe=!0;else if(j.isInstancedMesh&&Ft.instancingMorph===!1&&j.morphTexture!==null)oe=!0;else if(Ft.envMap!==Vt)oe=!0;else if(K.fog===!0&&Ft.fog!==It)oe=!0;else if(Ft.numClippingPlanes!==void 0&&(Ft.numClippingPlanes!==vt.numPlanes||Ft.numIntersection!==vt.numIntersection))oe=!0;else if(Ft.vertexAlphas!==ne)oe=!0;else if(Ft.vertexTangents!==se)oe=!0;else if(Ft.morphTargets!==Gt)oe=!0;else if(Ft.morphNormals!==ue)oe=!0;else if(Ft.morphColors!==we)oe=!0;else if(Ft.toneMapping!==Se)oe=!0;else if(Ft.morphTargetsCount!==ze)oe=!0;else if(!!Ft.lightProbeGrid!==T.state.lightProbeGridArray.length>0)oe=!0}else oe=!0,Ft.__version=K.version;let hn=Ft.currentProgram;if(oe===!0){if(hn=Fr(K,V,j),D&&K.isNodeMaterial)D.onUpdateProgram(K,hn,Ft)}let Rn=!1,Qn=!1,Yi=!1,ge=hn.getUniforms(),Ee=Ft.uniforms;if(w.useProgram(hn.program))Rn=!0,Qn=!0,Yi=!0;if(K.id!==ft)ft=K.id,Qn=!0;if(Ft.needsLights){let ve=fm(T.state.lightProbeGridArray,j);if(Ft.lightProbeGrid!==ve)Ft.lightProbeGrid=ve,Qn=!0}if(Rn||Q!==E){if(w.buffers.depth.getReversed()&&E.reversedDepth!==!0)E._reversedDepth=!0,E.updateProjectionMatrix();ge.setValue(v,"projectionMatrix",E.projectionMatrix),ge.setValue(v,"viewMatrix",E.matrixWorldInverse);let ei=ge.map.cameraPosition;if(ei!==void 0)ei.setValue(v,tt.setFromMatrixPosition(E.matrixWorld));if($t.logarithmicDepthBuffer)ge.setValue(v,"logDepthBufFC",2/(Math.log(E.far+1)/Math.LN2));if(K.isMeshPhongMaterial||K.isMeshToonMaterial||K.isMeshLambertMaterial||K.isMeshBasicMaterial||K.isMeshStandardMaterial||K.isShaderMaterial)ge.setValue(v,"isOrthographic",E.isOrthographicCamera===!0);if(Q!==E)Q=E,Qn=!0,Yi=!0}if(Ft.needsLights){if(qe.state.sunShadowMap.length>0)ge.setValue(v,"sunShadowMap",qe.state.sunShadowMap,it);if(qe.state.directionalShadowMap.length>0)ge.setValue(v,"directionalShadowMap",qe.state.directionalShadowMap,it);if(qe.state.spotShadowMap.length>0)ge.setValue(v,"spotShadowMap",qe.state.spotShadowMap,it);if(qe.state.pointShadowMap.length>0)ge.setValue(v,"pointShadowMap",qe.state.pointShadowMap,it)}if(j.isSkinnedMesh){ge.setOptional(v,j,"bindMatrix"),ge.setOptional(v,j,"bindMatrixInverse");let ve=j.skeleton;if(ve){if(ve.boneTexture===null)ve.computeBoneTexture();ge.setValue(v,"boneTexture",ve.boneTexture,it)}}if(j.isBatchedMesh){if(ge.setOptional(v,j,"batchingTexture"),ge.setValue(v,"batchingTexture",j._matricesTexture,it),ge.setOptional(v,j,"batchingIdTexture"),ge.setValue(v,"batchingIdTexture",j._indirectTexture,it),ge.setOptional(v,j,"batchingColorTexture"),j._colorsTexture!==null)ge.setValue(v,"batchingColorTexture",j._colorsTexture,it)}let ti=st.morphAttributes;if(ti.position!==void 0||ti.normal!==void 0||ti.color!==void 0)ce.update(j,st,hn);if(Qn||Ft.receiveShadow!==j.receiveShadow)Ft.receiveShadow=j.receiveShadow,ge.setValue(v,"receiveShadow",j.receiveShadow);if((K.isMeshStandardMaterial||K.isMeshLambertMaterial||K.isMeshPhongMaterial)&&K.envMap===null&&V.environment!==null)Ee.envMapIntensity.value=V.environmentIntensity;if(Ee.dfgLUT!==void 0)Ee.dfgLUT.value=uM();if(Qn){if(ge.setValue(v,"toneMappingExposure",A.toneMappingExposure),Ft.needsLights)mm(Ee,Yi);if(It&&K.fog===!0)wt.refreshFogUniforms(Ee,It);if(wt.refreshMaterialUniforms(Ee,K,O,F,T.state.transmissionRenderTarget[E.id]),Ft.needsLights&&Ft.lightProbeGrid){let ve=Ft.lightProbeGrid;Ee.probesSH.value=ve.texture,Ee.probesMin.value.copy(ve.boundingBox.min),Ee.probesMax.value.copy(ve.boundingBox.max),Ee.probesResolution.value.copy(ve.resolution)}Lr.upload(v,kh(Ft),Ee,it)}if(K.isShaderMaterial&&K.uniformsNeedUpdate===!0)Lr.upload(v,kh(Ft),Ee,it),K.uniformsNeedUpdate=!1;if(K.isSpriteMaterial)ge.setValue(v,"center",j.center);if(ge.setValue(v,"modelViewMatrix",j.modelViewMatrix),ge.setValue(v,"normalMatrix",j.normalMatrix),ge.setValue(v,"modelMatrix",j.matrixWorld),K.uniformsGroups!==void 0){let ve=K.uniformsGroups;for(let ei=0,Zi=ve.length;ei<Zi;ei++){let Vh=ve[ei];Ut.update(Vh,hn),Ut.bind(Vh,hn)}}return hn}function mm(E,V){E.ambientLightColor.needsUpdate=V,E.lightProbe.needsUpdate=V,E.sunLights.needsUpdate=V,E.sunLightShadows.needsUpdate=V,E.directionalLights.needsUpdate=V,E.directionalLightShadows.needsUpdate=V,E.pointLights.needsUpdate=V,E.pointLightShadows.needsUpdate=V,E.spotLights.needsUpdate=V,E.spotLightShadows.needsUpdate=V,E.rectAreaLights.needsUpdate=V,E.hemisphereLights.needsUpdate=V}function gm(E){return E.isMeshLambertMaterial||E.isMeshToonMaterial||E.isMeshPhongMaterial||E.isMeshStandardMaterial||E.isShadowMaterial||E.isShaderMaterial&&E.lights===!0}this.getActiveCubeFace=function(){return et},this.getActiveMipmapLevel=function(){return nt},this.getRenderTarget=function(){return B},this.setRenderTargetTextures=function(E,V,st){let K=k.get(E);if(K.__autoAllocateDepthBuffer=E.resolveDepthBuffer===!1,K.__autoAllocateDepthBuffer===!1)K.__useRenderToTexture=!1;k.get(E.texture).__webglTexture=V,k.get(E.depthTexture).__webglTexture=K.__autoAllocateDepthBuffer?void 0:st,K.__hasExternalTextures=!0},this.setRenderTargetFramebuffer=function(E,V){let st=k.get(E);st.__webglFramebuffer=V,st.__useDefaultFramebuffer=V===void 0},this.setRenderTarget=function(E,V=0,st=0){B=E,et=V,nt=st;let K=null,j=!1,It=!1;if(E){let Ct=k.get(E);if(Ct.__useDefaultFramebuffer!==void 0){w.bindFramebuffer(v.FRAMEBUFFER,Ct.__webglFramebuffer),ct.copy(E.viewport),dt.copy(E.scissor),Bt=E.scissorTest,w.viewport(ct),w.scissor(dt),w.setScissorTest(Bt),ft=-1;return}else if(Ct.__webglFramebuffer===void 0)it.setupRenderTarget(E);else if(Ct.__hasExternalTextures)it.rebindTextures(E,k.get(E.texture).__webglTexture,k.get(E.depthTexture).__webglTexture);else if(E.depthBuffer){let ne=E.depthTexture;if(Ct.__boundDepthTexture!==ne){if(ne!==null&&k.has(ne)&&(E.width!==ne.image.width||E.height!==ne.image.height))throw Error("THREE.WebGLRenderer: Attached DepthTexture is initialized to the incorrect size.");it.setupDepthRenderbuffer(E)}}let kt=E.texture;if(kt.isData3DTexture||kt.isDataArrayTexture||kt.isCompressedArrayTexture)It=!0;let Vt=k.get(E).__webglFramebuffer;if(E.isWebGLCubeRenderTarget){if(Array.isArray(Vt[V]))K=Vt[V][st];else K=Vt[V];j=!0}else if(E.samples>0&&it.useMultisampledRTT(E)===!1)K=k.get(E).__webglMultisampledFramebuffer;else if(Array.isArray(Vt))K=Vt[st];else K=Vt;ct.copy(E.viewport),dt.copy(E.scissor),Bt=E.scissorTest}else ct.copy(G).multiplyScalar(O).floor(),dt.copy(J).multiplyScalar(O).floor(),Bt=_t;if(st!==0)K=q;if(w.bindFramebuffer(v.FRAMEBUFFER,K))w.drawBuffers(E,K);if(w.viewport(ct),w.scissor(dt),w.setScissorTest(Bt),j){let Ct=k.get(E.texture);v.framebufferTexture2D(v.FRAMEBUFFER,v.COLOR_ATTACHMENT0,v.TEXTURE_CUBE_MAP_POSITIVE_X+V,Ct.__webglTexture,st)}else if(It){let Ct=V;for(let kt=0;kt<E.textures.length;kt++){let Vt=k.get(E.textures[kt]);v.framebufferTextureLayer(v.FRAMEBUFFER,v.COLOR_ATTACHMENT0+kt,Vt.__webglTexture,st,Ct)}}else if(E!==null&&st!==0){let Ct=k.get(E.texture);v.framebufferTexture2D(v.FRAMEBUFFER,v.COLOR_ATTACHMENT0,v.TEXTURE_2D,Ct.__webglTexture,st)}ft=-1};function Hh(E){let V=k.get(E);if(V.__readFormat!==E.format||V.__readType!==E.type)V.__readFormat=E.format,V.__readType=E.type,V.__formatReadable=$t.textureFormatReadable(E.format),V.__typeReadable=$t.textureTypeReadable(E.type);return V}if(this.readRenderTargetPixels=function(E,V,st,K,j,It,Ot,Ct=0){if(!(E&&E.isWebGLRenderTarget)){Wt("WebGLRenderer.readRenderTargetPixels: renderTarget is not THREE.WebGLRenderTarget.");return}let kt=k.get(E).__webglFramebuffer;if(E.isWebGLCubeRenderTarget&&Ot!==void 0)kt=kt[Ot];if(kt){w.bindFramebuffer(v.FRAMEBUFFER,kt);try{let Vt=E.textures[Ct],{format:ne,type:se}=Vt;if(E.textures.length>1)v.readBuffer(v.COLOR_ATTACHMENT0+Ct);let Gt=Hh(Vt);if(Gt.__formatReadable===!1){Wt("WebGLRenderer.readRenderTargetPixels: renderTarget is not in RGBA or implementation defined format.");return}if(Gt.__typeReadable===!1){Wt("WebGLRenderer.readRenderTargetPixels: renderTarget is not in UnsignedByteType or implementation defined type.");return}if(V>=0&&V<=E.width-K&&(st>=0&&st<=E.height-j))v.readPixels(V,st,K,j,ht.convert(ne),ht.convert(se),It)}finally{let Vt=B!==null?k.get(B).__webglFramebuffer:null;w.bindFramebuffer(v.FRAMEBUFFER,Vt)}}},this.readRenderTargetPixelsAsync=async function(E,V,st,K,j,It,Ot,Ct=0){if(!(E&&E.isWebGLRenderTarget))throw Error("THREE.WebGLRenderer.readRenderTargetPixels: renderTarget is not THREE.WebGLRenderTarget.");let kt=k.get(E).__webglFramebuffer;if(E.isWebGLCubeRenderTarget&&Ot!==void 0)kt=kt[Ot];if(kt)if(V>=0&&V<=E.width-K&&(st>=0&&st<=E.height-j)){w.bindFramebuffer(v.FRAMEBUFFER,kt);let Vt=E.textures[Ct],{format:ne,type:se}=Vt;if(E.textures.length>1)v.readBuffer(v.COLOR_ATTACHMENT0+Ct);let Gt=Hh(Vt);if(Gt.__formatReadable===!1)throw Error("THREE.WebGLRenderer.readRenderTargetPixelsAsync: renderTarget is not in RGBA or implementation defined format.");if(Gt.__typeReadable===!1)throw Error("THREE.WebGLRenderer.readRenderTargetPixelsAsync: renderTarget is not in UnsignedByteType or implementation defined type.");let ue=v.createBuffer();v.bindBuffer(v.PIXEL_PACK_BUFFER,ue),v.bufferData(v.PIXEL_PACK_BUFFER,It.byteLength,v.STREAM_READ),v.readPixels(V,st,K,j,ht.convert(ne),ht.convert(se),0),v.bindBuffer(v.PIXEL_PACK_BUFFER,null);let we=B!==null?k.get(B).__webglFramebuffer:null;w.bindFramebuffer(v.FRAMEBUFFER,we);let Se=v.fenceSync(v.SYNC_GPU_COMMANDS_COMPLETE,0);return v.flush(),await nf(v,Se,4),v.bindBuffer(v.PIXEL_PACK_BUFFER,ue),v.getBufferSubData(v.PIXEL_PACK_BUFFER,0,It),v.bindBuffer(v.PIXEL_PACK_BUFFER,null),v.deleteBuffer(ue),v.deleteSync(Se),It}else throw Error("THREE.WebGLRenderer.readRenderTargetPixelsAsync: requested read bounds are out of range.")},this.copyFramebufferToTexture=function(E,V=null,st=0){let K=Math.pow(2,-st),j=Math.floor(E.image.width*K),It=Math.floor(E.image.height*K),Ot=V!==null?V.x:0,Ct=V!==null?V.y:0;it.setTexture2D(E,0),v.copyTexSubImage2D(v.TEXTURE_2D,st,0,0,Ot,Ct,j,It),w.unbindTexture()},this.copyTextureToTexture=function(E,V,st=null,K=null,j=0,It=0){let Ot,Ct,kt,Vt,ne,se,Gt,ue,we,Se=E.isCompressedTexture?E.mipmaps[It]:E.image;if(st!==null)Ot=st.max.x-st.min.x,Ct=st.max.y-st.min.y,kt=st.isBox3?st.max.z-st.min.z:1,Vt=st.min.x,ne=st.min.y,se=st.isBox3?st.min.z:0;else{let Ee=Math.pow(2,-j);if(Ot=Math.floor(Se.width*Ee),Ct=Math.floor(Se.height*Ee),E.isDataArrayTexture)kt=Se.depth;else if(E.isData3DTexture)kt=Math.floor(Se.depth*Ee);else kt=1;Vt=0,ne=0,se=0}if(K!==null)Gt=K.x,ue=K.y,we=K.z;else Gt=0,ue=0,we=0;let _e=ht.convert(V.format),ze=ht.convert(V.type),Ft;if(V.isData3DTexture)it.setTexture3D(V,0),Ft=v.TEXTURE_3D;else if(V.isDataArrayTexture||V.isCompressedArrayTexture)it.setTexture2DArray(V,0),Ft=v.TEXTURE_2D_ARRAY;else it.setTexture2D(V,0),Ft=v.TEXTURE_2D;w.activeTexture(v.TEXTURE0),w.pixelStorei(v.UNPACK_FLIP_Y_WEBGL,V.flipY),w.pixelStorei(v.UNPACK_PREMULTIPLY_ALPHA_WEBGL,V.premultiplyAlpha),w.pixelStorei(v.UNPACK_ALIGNMENT,V.unpackAlignment);let qe=w.getParameter(v.UNPACK_ROW_LENGTH),oe=w.getParameter(v.UNPACK_IMAGE_HEIGHT),hn=w.getParameter(v.UNPACK_SKIP_PIXELS),Rn=w.getParameter(v.UNPACK_SKIP_ROWS),Qn=w.getParameter(v.UNPACK_SKIP_IMAGES);w.pixelStorei(v.UNPACK_ROW_LENGTH,Se.width),w.pixelStorei(v.UNPACK_IMAGE_HEIGHT,Se.height),w.pixelStorei(v.UNPACK_SKIP_PIXELS,Vt),w.pixelStorei(v.UNPACK_SKIP_ROWS,ne),w.pixelStorei(v.UNPACK_SKIP_IMAGES,se);let Yi=E.isDataArrayTexture||E.isData3DTexture,ge=V.isDataArrayTexture||V.isData3DTexture;if(E.isDepthTexture){let Ee=k.get(E),ti=k.get(V),ve=k.get(Ee.__renderTarget),ei=k.get(ti.__renderTarget);w.bindFramebuffer(v.READ_FRAMEBUFFER,ve.__webglFramebuffer),w.bindFramebuffer(v.DRAW_FRAMEBUFFER,ei.__webglFramebuffer);for(let Zi=0;Zi<kt;Zi++){if(Yi)v.framebufferTextureLayer(v.READ_FRAMEBUFFER,v.COLOR_ATTACHMENT0,k.get(E).__webglTexture,j,se+Zi),v.framebufferTextureLayer(v.DRAW_FRAMEBUFFER,v.COLOR_ATTACHMENT0,k.get(V).__webglTexture,It,we+Zi);v.blitFramebuffer(Vt,ne,Ot,Ct,Gt,ue,Ot,Ct,v.DEPTH_BUFFER_BIT,v.NEAREST)}w.bindFramebuffer(v.READ_FRAMEBUFFER,null),w.bindFramebuffer(v.DRAW_FRAMEBUFFER,null)}else if(j!==0||E.isRenderTargetTexture||k.has(E)){let Ee=k.get(E),ti=k.get(V);w.bindFramebuffer(v.READ_FRAMEBUFFER,ot),w.bindFramebuffer(v.DRAW_FRAMEBUFFER,z);for(let ve=0;ve<kt;ve++){if(Yi)v.framebufferTextureLayer(v.READ_FRAMEBUFFER,v.COLOR_ATTACHMENT0,Ee.__webglTexture,j,se+ve);else v.framebufferTexture2D(v.READ_FRAMEBUFFER,v.COLOR_ATTACHMENT0,v.TEXTURE_2D,Ee.__webglTexture,j);if(ge)v.framebufferTextureLayer(v.DRAW_FRAMEBUFFER,v.COLOR_ATTACHMENT0,ti.__webglTexture,It,we+ve);else v.framebufferTexture2D(v.DRAW_FRAMEBUFFER,v.COLOR_ATTACHMENT0,v.TEXTURE_2D,ti.__webglTexture,It);if(j!==0)v.blitFramebuffer(Vt,ne,Ot,Ct,Gt,ue,Ot,Ct,v.COLOR_BUFFER_BIT,v.NEAREST);else if(ge)v.copyTexSubImage3D(Ft,It,Gt,ue,we+ve,Vt,ne,Ot,Ct);else v.copyTexSubImage2D(Ft,It,Gt,ue,Vt,ne,Ot,Ct)}w.bindFramebuffer(v.READ_FRAMEBUFFER,null),w.bindFramebuffer(v.DRAW_FRAMEBUFFER,null)}else if(ge)if(E.isDataTexture||E.isData3DTexture)v.texSubImage3D(Ft,It,Gt,ue,we,Ot,Ct,kt,_e,ze,Se.data);else if(V.isCompressedArrayTexture)v.compressedTexSubImage3D(Ft,It,Gt,ue,we,Ot,Ct,kt,_e,Se.data);else v.texSubImage3D(Ft,It,Gt,ue,we,Ot,Ct,kt,_e,ze,Se);else if(E.isDataTexture)v.texSubImage2D(v.TEXTURE_2D,It,Gt,ue,Ot,Ct,_e,ze,Se.data);else if(E.isCompressedTexture)v.compressedTexSubImage2D(v.TEXTURE_2D,It,Gt,ue,Se.width,Se.height,_e,Se.data);else v.texSubImage2D(v.TEXTURE_2D,It,Gt,ue,Ot,Ct,_e,ze,Se);if(w.pixelStorei(v.UNPACK_ROW_LENGTH,qe),w.pixelStorei(v.UNPACK_IMAGE_HEIGHT,oe),w.pixelStorei(v.UNPACK_SKIP_PIXELS,hn),w.pixelStorei(v.UNPACK_SKIP_ROWS,Rn),w.pixelStorei(v.UNPACK_SKIP_IMAGES,Qn),It===0&&V.generateMipmaps)v.generateMipmap(Ft);w.unbindTexture()},this.initRenderTarget=function(E){if(k.get(E).__webglFramebuffer===void 0)it.setupRenderTarget(E)},this.initTexture=function(E){if(E.isCubeTexture)it.setTextureCube(E,0);else if(E.isData3DTexture)it.setTexture3D(E,0);else if(E.isDataArrayTexture||E.isCompressedArrayTexture)it.setTexture2DArray(E,0);else it.setTexture2D(E,0);w.unbindTexture()},this.resetState=function(){et=0,nt=0,B=null,w.reset(),Et.reset()},typeof __THREE_DEVTOOLS__<"u")__THREE_DEVTOOLS__.dispatchEvent(new CustomEvent("observe",{detail:this}))}get coordinateSystem(){return Tc}get outputColorSpace(){return this._outputColorSpace}set outputColorSpace(t){this._outputColorSpace=t;let e=this.getContext();e.drawingBufferColorSpace=re._getDrawingBufferColorSpace(t),e.unpackColorSpace=re._getUnpackColorSpace()}}var Nr=new L;function _n(t,e,n,i,s,r){let a=2*Math.PI*s/4,o=Math.max(r-2*s,0),l=Math.PI/4;Nr.copy(e),Nr[i]=0,Nr.normalize();let c=0.5*a/(a+o),h=1-Nr.angleTo(t)/l;if(Math.sign(Nr[n])===1)return h*c;else return o/(a+o)+c+c*(1-h)}class Ih extends An{constructor(t=1,e=1,n=1,i=2,s=0.1){let r=i*2+1;s=Math.min(t/2,e/2,n/2,s);super(1,1,1,r,r,r);if(this.type="RoundedBoxGeometry",this.parameters={width:t,height:e,depth:n,segments:i,radius:s},r===1)return;let a=this.toNonIndexed();this.index=null,this.attributes.position=a.attributes.position,this.attributes.normal=a.attributes.normal,this.attributes.uv=a.attributes.uv;let o=new L,l=new L,c=new L(t,e,n).divideScalar(2).subScalar(s),h=this.attributes.position.array,d=this.attributes.normal.array,u=this.attributes.uv.array,f=h.length/6,g=new L,_=0.5/r;for(let m=0,p=0;m<h.length;m+=3,p+=2)switch(o.fromArray(h,m),l.copy(o),l.x-=Math.sign(l.x)*_,l.y-=Math.sign(l.y)*_,l.z-=Math.sign(l.z)*_,l.normalize(),h[m+0]=c.x*Math.sign(o.x)+l.x*s,h[m+1]=c.y*Math.sign(o.y)+l.y*s,h[m+2]=c.z*Math.sign(o.z)+l.z*s,d[m+0]=l.x,d[m+1]=l.y,d[m+2]=l.z,Math.floor(m/f)){case 0:g.set(1,0,0),u[p+0]=_n(g,l,"z","y",s,n),u[p+1]=1-_n(g,l,"y","z",s,e);break;case 1:g.set(-1,0,0),u[p+0]=1-_n(g,l,"z","y",s,n),u[p+1]=1-_n(g,l,"y","z",s,e);break;case 2:g.set(0,1,0),u[p+0]=1-_n(g,l,"x","z",s,t),u[p+1]=_n(g,l,"z","x",s,n);break;case 3:g.set(0,-1,0),u[p+0]=1-_n(g,l,"x","z",s,t),u[p+1]=1-_n(g,l,"z","x",s,n);break;case 4:g.set(0,0,1),u[p+0]=1-_n(g,l,"x","y",s,t),u[p+1]=1-_n(g,l,"y","x",s,e);break;case 5:g.set(0,0,-1),u[p+0]=_n(g,l,"x","y",s,t),u[p+1]=1-_n(g,l,"y","x",s,e);break}}static fromJSON(t){return new Ih(t.width,t.height,t.depth,t.segments,t.radius)}}var Dr=Ga;class Ph extends Be{constructor(t){super(t);this.defaultDPI=90,this.defaultUnit="px"}load(t,e,n,i){let s=this,r=new fn(s.manager);r.setPath(s.path),r.setRequestHeader(s.requestHeader),r.setWithCredentials(s.withCredentials),r.load(t,function(a){try{e(s.parse(a))}catch(o){if(i)i(o);else console.error(o);s.manager.itemError(t)}},n,i)}parse(t){let e=this;function n(O,N){if(O.nodeType!==1)return;if(O.hasAttribute("filter"))console.warn("THREE.SVGLoader: Filters are not supported.");let P=x(O),G=!1,J=null;switch(O.nodeName){case"svg":N=g(O,N);break;case"style":s(O);break;case"g":N=g(O,N);break;case"path":if(N=g(O,N),O.hasAttribute("d"))J=i(O);break;case"rect":N=g(O,N),J=o(O);break;case"polygon":N=g(O,N),J=l(O);break;case"polyline":N=g(O,N),J=c(O);break;case"circle":N=g(O,N),J=h(O);break;case"ellipse":N=g(O,N),J=d(O);break;case"line":N=g(O,N),J=u(O);break;case"defs":G=!0;break;case"use":N=g(O,N);let H=(O.getAttributeNS("http://www.w3.org/1999/xlink","href")||"").substring(1),U=O.viewportElement.getElementById(H);if(U)n(U,N);else console.warn("SVGLoader: 'use node' references non-existent node id: "+H);break;default:}if(J){if(N.fill!==void 0&&N.fill!=="none"&&!N.fill.startsWith("url"))J.color.setStyle(N.fill,Dr);R(J,ee),z.push(J);let C=Object.assign({},N);C.strokeWidth=N.strokeWidth*q(ee),J.userData={node:O,style:C,transform:ee.clone(),gradients:nt}}let _t=O.childNodes;for(let C=0;C<_t.length;C++){let H=_t[C];if(G&&H.nodeName!=="style"&&H.nodeName!=="defs")continue;n(H,N)}if(P)if(B.pop(),B.length>0)ee.copy(B[B.length-1]);else ee.identity()}function i(O){let N=new Un,P=new rt,G=new rt,J=new rt,_t=!0,C=!1,H=O.getAttribute("d");if(H===""||H==="none")return null;let U=H.match(/[a-df-z][^a-df-z]*/ig);for(let Y=0,tt=U.length;Y<tt;Y++){let at=U[Y],mt=at.charAt(0),gt=at.slice(1).trim();if(_t===!0)C=!0,_t=!1;let X;switch(mt){case"M":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=2){if(P.x=X[v+0],P.y=X[v+1],G.x=P.x,G.y=P.y,v===0)N.moveTo(P.x,P.y);else N.lineTo(P.x,P.y);if(v===0)J.copy(P)}break;case"H":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v++)if(P.x=X[v],G.x=P.x,G.y=P.y,N.lineTo(P.x,P.y),v===0&&C===!0)J.copy(P);break;case"V":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v++)if(P.y=X[v],G.x=P.x,G.y=P.y,N.lineTo(P.x,P.y),v===0&&C===!0)J.copy(P);break;case"L":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=2)if(P.x=X[v+0],P.y=X[v+1],G.x=P.x,G.y=P.y,N.lineTo(P.x,P.y),v===0&&C===!0)J.copy(P);break;case"C":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=6)if(N.bezierCurveTo(X[v+0],X[v+1],X[v+2],X[v+3],X[v+4],X[v+5]),G.x=X[v+2],G.y=X[v+3],P.x=X[v+4],P.y=X[v+5],v===0&&C===!0)J.copy(P);break;case"S":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=4)if(N.bezierCurveTo(_(P.x,G.x),_(P.y,G.y),X[v+0],X[v+1],X[v+2],X[v+3]),G.x=X[v+0],G.y=X[v+1],P.x=X[v+2],P.y=X[v+3],v===0&&C===!0)J.copy(P);break;case"Q":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=4)if(N.quadraticCurveTo(X[v+0],X[v+1],X[v+2],X[v+3]),G.x=X[v+0],G.y=X[v+1],P.x=X[v+2],P.y=X[v+3],v===0&&C===!0)J.copy(P);break;case"T":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=2){let zt=_(P.x,G.x),$t=_(P.y,G.y);if(N.quadraticCurveTo(zt,$t,X[v+0],X[v+1]),G.x=zt,G.y=$t,P.x=X[v+0],P.y=X[v+1],v===0&&C===!0)J.copy(P)}break;case"A":X=m(gt,[3,4],7);for(let v=0,Lt=X.length;v<Lt;v+=7){if(X[v+5]==P.x&&X[v+6]==P.y)continue;let zt=P.clone();if(P.x=X[v+5],P.y=X[v+6],G.x=P.x,G.y=P.y,r(N,X[v],X[v+1],X[v+2],X[v+3],X[v+4],zt,P),v===0&&C===!0)J.copy(P)}break;case"m":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=2){if(P.x+=X[v+0],P.y+=X[v+1],G.x=P.x,G.y=P.y,v===0)N.moveTo(P.x,P.y);else N.lineTo(P.x,P.y);if(v===0)J.copy(P)}break;case"h":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v++)if(P.x+=X[v],G.x=P.x,G.y=P.y,N.lineTo(P.x,P.y),v===0&&C===!0)J.copy(P);break;case"v":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v++)if(P.y+=X[v],G.x=P.x,G.y=P.y,N.lineTo(P.x,P.y),v===0&&C===!0)J.copy(P);break;case"l":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=2)if(P.x+=X[v+0],P.y+=X[v+1],G.x=P.x,G.y=P.y,N.lineTo(P.x,P.y),v===0&&C===!0)J.copy(P);break;case"c":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=6)if(N.bezierCurveTo(P.x+X[v+0],P.y+X[v+1],P.x+X[v+2],P.y+X[v+3],P.x+X[v+4],P.y+X[v+5]),G.x=P.x+X[v+2],G.y=P.y+X[v+3],P.x+=X[v+4],P.y+=X[v+5],v===0&&C===!0)J.copy(P);break;case"s":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=4)if(N.bezierCurveTo(_(P.x,G.x),_(P.y,G.y),P.x+X[v+0],P.y+X[v+1],P.x+X[v+2],P.y+X[v+3]),G.x=P.x+X[v+0],G.y=P.y+X[v+1],P.x+=X[v+2],P.y+=X[v+3],v===0&&C===!0)J.copy(P);break;case"q":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=4)if(N.quadraticCurveTo(P.x+X[v+0],P.y+X[v+1],P.x+X[v+2],P.y+X[v+3]),G.x=P.x+X[v+0],G.y=P.y+X[v+1],P.x+=X[v+2],P.y+=X[v+3],v===0&&C===!0)J.copy(P);break;case"t":X=m(gt);for(let v=0,Lt=X.length;v<Lt;v+=2){let zt=_(P.x,G.x),$t=_(P.y,G.y);if(N.quadraticCurveTo(zt,$t,P.x+X[v+0],P.y+X[v+1]),G.x=zt,G.y=$t,P.x=P.x+X[v+0],P.y=P.y+X[v+1],v===0&&C===!0)J.copy(P)}break;case"a":X=m(gt,[3,4],7);for(let v=0,Lt=X.length;v<Lt;v+=7){if(X[v+5]==0&&X[v+6]==0)continue;let zt=P.clone();if(P.x+=X[v+5],P.y+=X[v+6],G.x=P.x,G.y=P.y,r(N,X[v],X[v+1],X[v+2],X[v+3],X[v+4],zt,P),v===0&&C===!0)J.copy(P)}break;case"Z":case"z":if(N.currentPath.autoClose=!0,N.currentPath.curves.length>0)P.copy(J),N.currentPath.currentPoint.copy(P),_t=!0;break;default:console.warn(at)}C=!1}return N}function s(O){if(!O.sheet||!O.sheet.cssRules||!O.sheet.cssRules.length)return;for(let N=0;N<O.sheet.cssRules.length;N++){let P=O.sheet.cssRules[N];if(P.type!==1)continue;let G=P.selectorText.split(/,/gm).filter(Boolean).map((J)=>J.trim());for(let J=0;J<G.length;J++){let _t=Object.fromEntries(Object.entries(P.style).filter(([,C])=>C!==""));et[G[J]]=Object.assign(et[G[J]]||{},_t)}}}function r(O,N,P,G,J,_t,C,H){if(N==0||P==0){O.lineTo(H.x,H.y);return}G=G*Math.PI/180,N=Math.abs(N),P=Math.abs(P);let U=(C.x-H.x)/2,Y=(C.y-H.y)/2,tt=Math.cos(G)*U+Math.sin(G)*Y,at=-Math.sin(G)*U+Math.cos(G)*Y,mt=N*N,gt=P*P,X=tt*tt,v=at*at,Lt=X/mt+v/gt;if(Lt>1){let lt=Math.sqrt(Lt);N=lt*N,P=lt*P,mt=N*N,gt=P*P}let zt=mt*v+gt*X,$t=(mt*gt-zt)/zt,w=Math.sqrt(Math.max(0,$t));if(J===_t)w=-w;let M=w*N*at/P,k=-w*P*tt/N,it=Math.cos(G)*M-Math.sin(G)*k+(C.x+H.x)/2,pt=Math.sin(G)*M+Math.cos(G)*k+(C.y+H.y)/2,yt=a(1,0,(tt-M)/N,(at-k)/P),St=a((tt-M)/N,(at-k)/P,(-tt-M)/N,(-at-k)/P)%(Math.PI*2);O.currentPath.absellipse(it,pt,N,P,yt,yt+St,_t===0,G)}function a(O,N,P,G){let J=O*P+N*G,_t=Math.sqrt(O*O+N*N)*Math.sqrt(P*P+G*G),C=Math.acos(Math.max(-1,Math.min(1,J/_t)));if(O*G-N*P<0)C=-C;return C}function o(O){let N=S(O.getAttribute("x")||0),P=S(O.getAttribute("y")||0),G=S(O.getAttribute("rx")||O.getAttribute("ry")||0),J=S(O.getAttribute("ry")||O.getAttribute("rx")||0),_t=S(O.getAttribute("width")),C=S(O.getAttribute("height")),H=0.448084975506,U=new Un;if(U.moveTo(N+G,P),U.lineTo(N+_t-G,P),G!==0||J!==0)U.bezierCurveTo(N+_t-G*0.448084975506,P,N+_t,P+J*0.448084975506,N+_t,P+J);if(U.lineTo(N+_t,P+C-J),G!==0||J!==0)U.bezierCurveTo(N+_t,P+C-J*0.448084975506,N+_t-G*0.448084975506,P+C,N+_t-G,P+C);if(U.lineTo(N+G,P+C),G!==0||J!==0)U.bezierCurveTo(N+G*0.448084975506,P+C,N,P+C-J*0.448084975506,N,P+C-J);if(U.lineTo(N,P+J),G!==0||J!==0)U.bezierCurveTo(N,P+J*0.448084975506,N+G*0.448084975506,P,N+G,P);return U}function l(O){function N(_t,C,H){let U=S(C),Y=S(H);if(J===0)G.moveTo(U,Y);else G.lineTo(U,Y);J++}let P=/([+-]?\d*\.?\d+(?:e[+-]?\d+)?)(?:,|\s)([+-]?\d*\.?\d+(?:e[+-]?\d+)?)/g,G=new Un,J=0;return O.getAttribute("points").replace(P,N),G.currentPath.autoClose=!0,G}function c(O){function N(_t,C,H){let U=S(C),Y=S(H);if(J===0)G.moveTo(U,Y);else G.lineTo(U,Y);J++}let P=/([+-]?\d*\.?\d+(?:e[+-]?\d+)?)(?:,|\s)([+-]?\d*\.?\d+(?:e[+-]?\d+)?)/g,G=new Un,J=0;return O.getAttribute("points").replace(P,N),G.currentPath.autoClose=!1,G}function h(O){let N=S(O.getAttribute("cx")||0),P=S(O.getAttribute("cy")||0),G=S(O.getAttribute("r")||0),J=new Jn;J.absarc(N,P,G,0,Math.PI*2);let _t=new Un;return _t.subPaths.push(J),_t}function d(O){let N=S(O.getAttribute("cx")||0),P=S(O.getAttribute("cy")||0),G=S(O.getAttribute("rx")||0),J=S(O.getAttribute("ry")||0),_t=new Jn;_t.absellipse(N,P,G,J,0,Math.PI*2);let C=new Un;return C.subPaths.push(_t),C}function u(O){let N=S(O.getAttribute("x1")||0),P=S(O.getAttribute("y1")||0),G=S(O.getAttribute("x2")||0),J=S(O.getAttribute("y2")||0),_t=new Un;return _t.moveTo(N,P),_t.lineTo(G,J),_t.currentPath.autoClose=!1,_t}function f(O){let P=O.querySelectorAll("linearGradient, radialGradient"),G=["x1","y1","x2","y2","cx","cy","r","fx","fy","gradientUnits","gradientTransform","spreadMethod"],J={};for(let C of P){let H=C.getAttribute("id");if(!H)continue;let U={type:C.nodeName==="radialGradient"?"radialGradient":"linearGradient",attrs:{},stops:null,href:null},Y=C.getAttributeNS("http://www.w3.org/1999/xlink","href")||C.getAttribute("href")||"";if(Y.startsWith("#"))U.href=Y.substring(1);for(let at of G)if(C.hasAttribute(at))U.attrs[at]=C.getAttribute(at);let tt=C.querySelectorAll("stop");if(tt.length>0){U.stops=[];for(let at of tt){let mt=at.getAttribute("stop-color");if(!mt&&at.style)mt=at.style["stop-color"];if(!mt)mt="#000";let gt=at.getAttribute("stop-opacity");if((gt===null||gt==="")&&at.style)gt=at.style["stop-opacity"];gt=gt===null||gt===""||gt===void 0?1:Math.max(0,Math.min(1,parseFloat(gt)));let X=Math.max(0,Math.min(1,parseFloat(at.getAttribute("offset")||"0")));U.stops.push({offset:X,color:mt,opacity:gt})}}J[H]=U}function _t(C,H){let U=J[C];if(!U||H.has(C))return U;if(H.add(C),U.href&&J[U.href]){let Y=_t(U.href,H);if(Y){if(!U.stops)U.stops=Y.stops;for(let tt in Y.attrs)if(!(tt in U.attrs))U.attrs[tt]=Y.attrs[tt]}}return U}for(let C in J)_t(C,new Set);for(let C in J){let at=function(mt){if(typeof mt!=="string")return 0;if(mt.endsWith("%"))return parseFloat(mt)/100;return S(mt)},H=J[C],U=H.attrs,Y=U.gradientUnits==="userSpaceOnUse"?"userSpaceOnUse":"objectBoundingBox",tt={type:H.type,gradientUnits:Y,spreadMethod:U.spreadMethod==="reflect"||U.spreadMethod==="repeat"?U.spreadMethod:"pad",gradientTransform:null,stops:(H.stops||[]).slice().sort((mt,gt)=>mt.offset-gt.offset)};if(U.gradientTransform)tt.gradientTransform=new Xt,I(U.gradientTransform,tt.gradientTransform);if(H.type==="linearGradient")tt.x1=U.x1!==void 0?at(U.x1):0,tt.y1=U.y1!==void 0?at(U.y1):0,tt.x2=U.x2!==void 0?at(U.x2):Y==="objectBoundingBox"?1:0,tt.y2=U.y2!==void 0?at(U.y2):0;else{let mt=Y==="objectBoundingBox"?0.5:0,gt=Y==="objectBoundingBox"?0.5:0;tt.cx=U.cx!==void 0?at(U.cx):mt,tt.cy=U.cy!==void 0?at(U.cy):mt,tt.r=U.r!==void 0?at(U.r):gt,tt.fx=U.fx!==void 0?at(U.fx):tt.cx,tt.fy=U.fy!==void 0?at(U.fy):tt.cy}nt[C]=tt}}function g(O,N){N=Object.assign({},N);let P={};if(O.hasAttribute("class")){let C=O.getAttribute("class").split(/\s/).filter(Boolean).map((H)=>H.trim());for(let H=0;H<C.length;H++)P=Object.assign(P,et["."+C[H]])}if(O.hasAttribute("id"))P=Object.assign(P,et["#"+O.getAttribute("id")]);function G(C,H,U){if(U===void 0)U=function(tt){return tt};if(O.hasAttribute(C))N[H]=U(O.getAttribute(C));if(P[H])N[H]=U(P[H]);if(O.style&&O.style[C]!=="")N[H]=U(O.style[C])}function J(C){return Math.max(0,Math.min(1,S(C)))}function _t(C){return Math.max(0,S(C))}return G("fill","fill"),G("fill-opacity","fillOpacity",J),G("fill-rule","fillRule"),G("opacity","opacity",J),G("stroke","stroke"),G("stroke-opacity","strokeOpacity",J),G("stroke-width","strokeWidth",_t),G("stroke-linejoin","strokeLineJoin"),G("stroke-linecap","strokeLineCap"),G("stroke-miterlimit","strokeMiterLimit",_t),G("visibility","visibility"),N}function _(O,N){return O-(N-O)}function m(O,N,P){if(typeof O!=="string")throw TypeError("Invalid input: "+typeof O);let G={SEPARATOR:/[ \t\r\n\,.\-+]/,WHITESPACE:/[ \t\r\n]/,DIGIT:/[\d]/,SIGN:/[-+]/,POINT:/\./,COMMA:/,/,EXP:/e/i,FLAGS:/[01]/},J=0,_t=1,C=2,H=3,U=J,Y=!0,tt="",at="",mt=[];function gt(zt,$t,w){let M=SyntaxError('Unexpected character "'+zt+'" at index '+$t+".");throw M.partial=w,M}function X(){if(tt!=="")if(at==="")mt.push(Number(tt));else mt.push(Number(tt)*Math.pow(10,Number(at)));tt="",at=""}let v,Lt=O.length;for(let zt=0;zt<Lt;zt++){if(v=O[zt],Array.isArray(N)&&N.includes(mt.length%P)&&G.FLAGS.test(v)){U=_t,tt=v,X();continue}if(U===J){if(G.WHITESPACE.test(v))continue;if(G.DIGIT.test(v)||G.SIGN.test(v)){U=_t,tt=v;continue}if(G.POINT.test(v)){U=C,tt=v;continue}if(G.COMMA.test(v)){if(Y)gt(v,zt,mt);Y=!0}}if(U===_t){if(G.DIGIT.test(v)){tt+=v;continue}if(G.POINT.test(v)){tt+=v,U=C;continue}if(G.EXP.test(v)){U=H;continue}if(G.SIGN.test(v)&&tt.length===1&&G.SIGN.test(tt[0]))gt(v,zt,mt)}if(U===C){if(G.DIGIT.test(v)){tt+=v;continue}if(G.EXP.test(v)){U=H;continue}if(G.POINT.test(v)&&tt[tt.length-1]===".")gt(v,zt,mt)}if(U===H){if(G.DIGIT.test(v)){at+=v;continue}if(G.SIGN.test(v)){if(at===""){at+=v;continue}if(at.length===1&&G.SIGN.test(at))gt(v,zt,mt)}}if(G.WHITESPACE.test(v))X(),U=J,Y=!1;else if(G.COMMA.test(v))X(),U=J,Y=!0;else if(G.SIGN.test(v))X(),U=_t,tt=v;else if(G.POINT.test(v))X(),U=C,tt=v;else gt(v,zt,mt)}return X(),mt}let p=["mm","cm","in","pt","pc","px"],b={mm:{mm:1,cm:0.1,in:0.03937007874015748,pt:2.834645669291339,pc:0.2362204724409449,px:-1},cm:{mm:10,cm:1,in:0.39370078740157477,pt:28.346456692913385,pc:2.3622047244094486,px:-1},in:{mm:25.4,cm:2.54,in:1,pt:72,pc:6,px:-1},pt:{mm:0.35277777777777775,cm:0.035277777777777776,in:0.013888888888888888,pt:1,pc:0.08333333333333333,px:-1},pc:{mm:4.233333333333333,cm:0.42333333333333334,in:0.16666666666666666,pt:12,pc:1,px:-1},px:{px:1}};function S(O){let N="px";if(typeof O==="string"||O instanceof String)for(let G=0,J=p.length;G<J;G++){let _t=p[G];if(O.endsWith(_t)){N=_t,O=O.substring(0,O.length-_t.length);break}}let P=void 0;if(N==="px"&&e.defaultUnit!=="px")P=b.in[e.defaultUnit]/e.defaultDPI;else if(P=b[N][e.defaultUnit],P<0)P=b[N].in*e.defaultDPI;return P*parseFloat(O)}function x(O){if(!(O.hasAttribute("transform")||O.nodeName==="use"&&(O.hasAttribute("x")||O.hasAttribute("y"))))return null;let N=T(O);if(B.length>0)N.premultiply(B[B.length-1]);return ee.copy(N),B.push(N),N}function T(O){let N=new Xt;if(O.nodeName==="use"&&(O.hasAttribute("x")||O.hasAttribute("y"))){let P=S(O.getAttribute("x")||0),G=S(O.getAttribute("y")||0);N.makeTranslation(P,G)}if(O.hasAttribute("transform"))I(O.getAttribute("transform"),N);return N}function I(O,N){let P=ft,G=O.split(")");for(let J=G.length-1;J>=0;J--){let _t=G[J].trim();if(_t==="")continue;let C=_t.indexOf("("),H=_t.length;if(C>0&&C<H){let U=_t.slice(0,C),Y=m(_t.slice(C+1));switch(P.identity(),U){case"translate":if(Y.length>=1){let tt=Y[0],at=0;if(Y.length>=2)at=Y[1];P.makeTranslation(tt,at)}break;case"rotate":if(Y.length>=1){let tt=0,at=0,mt=0;if(tt=Y[0]*Math.PI/180,Y.length>=3)at=Y[1],mt=Y[2];Q.makeTranslation(-at,-mt),ct.makeRotation(tt),dt.multiplyMatrices(ct,Q),Q.makeTranslation(at,mt),P.multiplyMatrices(Q,dt)}break;case"scale":if(Y.length>=1){let tt=Y[0],at=tt;if(Y.length>=2)at=Y[1];P.makeScale(tt,at)}break;case"skewX":if(Y.length===1)P.set(1,Math.tan(Y[0]*Math.PI/180),0,0,1,0,0,0,1);break;case"skewY":if(Y.length===1)P.set(1,0,0,Math.tan(Y[0]*Math.PI/180),1,0,0,0,1);break;case"matrix":if(Y.length===6)P.set(Y[0],Y[2],Y[4],Y[1],Y[3],Y[5],0,0,1);break}N.premultiply(P)}}return N}function R(O,N){function P(C){Dt.set(C.x,C.y,1).applyMatrix3(N),C.set(Dt.x,Dt.y)}function G(C){let{xRadius:H,yRadius:U}=C,Y=Math.cos(C.aRotation),tt=Math.sin(C.aRotation),at=new L(H*Y,H*tt,0),mt=new L(-U*tt,U*Y,0),gt=at.applyMatrix3(N),X=mt.applyMatrix3(N),v=ft.set(gt.x,X.x,0,gt.y,X.y,0,0,0,1),Lt=Q.copy(v).invert(),w=ct.copy(Lt).transpose().multiply(Lt).elements,M=ot(w[0],w[1],w[4]),k=Math.sqrt(M.rt1),it=Math.sqrt(M.rt2);if(C.xRadius=1/k,C.yRadius=1/it,C.aRotation=Math.atan2(M.sn,M.cs),!((C.aEndAngle-C.aStartAngle)%(2*Math.PI)<Number.EPSILON)){let yt=Q.set(k,0,0,0,it,0,0,0,1),St=ct.set(M.cs,M.sn,0,-M.sn,M.cs,0,0,0,1),lt=yt.multiply(St).multiply(v),ut=(wt)=>{let{x:Ht,y:bt}=new L(Math.cos(wt),Math.sin(wt),0).applyMatrix3(lt);return Math.atan2(bt,Ht)};if(C.aStartAngle=ut(C.aStartAngle),C.aEndAngle=ut(C.aEndAngle),y(N))C.aClockwise=!C.aClockwise}}function J(C){let H=Z(N),U=D(N);C.xRadius*=H,C.yRadius*=U;let Y=H>Number.EPSILON?Math.atan2(N.elements[1],N.elements[0]):Math.atan2(-N.elements[3],N.elements[4]);if(C.aRotation+=Y,y(N))C.aStartAngle*=-1,C.aEndAngle*=-1,C.aClockwise=!C.aClockwise}let _t=O.subPaths;for(let C=0,H=_t.length;C<H;C++){let Y=_t[C].curves;for(let tt=0;tt<Y.length;tt++){let at=Y[tt];if(at.isLineCurve)P(at.v1),P(at.v2);else if(at.isCubicBezierCurve)P(at.v0),P(at.v1),P(at.v2),P(at.v3);else if(at.isQuadraticBezierCurve)P(at.v0),P(at.v1),P(at.v2);else if(at.isEllipseCurve)if(Bt.set(at.aX,at.aY),P(Bt),at.aX=Bt.x,at.aY=Bt.y,A(N))G(at);else J(at)}}}function y(O){let N=O.elements;return N[0]*N[4]-N[1]*N[3]<0}function A(O){let N=O.elements,P=N[0]*N[3]+N[1]*N[4];if(P===0)return!1;let G=Z(O),J=D(O);return Math.abs(P/(G*J))>Number.EPSILON}function Z(O){let N=O.elements;return Math.sqrt(N[0]*N[0]+N[1]*N[1])}function D(O){let N=O.elements;return Math.sqrt(N[3]*N[3]+N[4]*N[4])}function q(O){let N=O.elements,P=N[0]*N[4]-N[1]*N[3];return Math.sqrt(Math.abs(P))}function ot(O,N,P){let G,J,_t,C,H,U=O+P,Y=O-P,tt=Math.sqrt(Y*Y+4*N*N);if(U>0)G=0.5*(U+tt),H=1/G,J=O*H*P-N*H*N;else if(U<0)J=0.5*(U-tt);else G=0.5*tt,J=-0.5*tt;if(Y>0)_t=Y+tt;else _t=Y-tt;if(Math.abs(_t)>2*Math.abs(N))H=-2*N/_t,C=1/Math.sqrt(1+H*H),_t=H*C;else if(Math.abs(N)===0)_t=1,C=0;else H=-0.5*_t/N,_t=1/Math.sqrt(1+H*H),C=H*_t;if(Y>0)H=_t,_t=-C,C=H;return{rt1:G,rt2:J,cs:_t,sn:C}}let z=[],et={},nt={},B=[],ft=new Xt,Q=new Xt,ct=new Xt,dt=new Xt,Bt=new rt,Dt=new L,ee=new Xt,Jt=new DOMParser().parseFromString(t,"image/svg+xml");return f(Jt),n(Jt.documentElement,{fill:"#000",fillOpacity:1,strokeOpacity:1,strokeWidth:1,strokeLineJoin:"miter",strokeLineCap:"butt",strokeMiterLimit:4}),{paths:z,gradients:nt,xml:Jt.documentElement}}static createFillMaterial(t){let e=t.userData.style;if(e.fill===void 0||e.fill==="none")return null;let n=t.color,i=null,s=im.exec(e.fill);if(s){let a=t.userData.gradients&&t.userData.gradients[s[1]];i=fM(a,t)}let r=new mn({opacity:e.fillOpacity*(e.opacity||1),transparent:!0,side:on,depthWrite:!1});if(i!==null)r.map=i;else r.color=n;return r}static createStrokeMaterial(t){let e=t.userData.style;if(e.stroke===void 0||e.stroke==="none")return null;if(im.test(e.stroke))console.warn("THREE.SVGLoader: Gradient strokes are not supported.");return new mn({color:new Pt().setStyle(e.stroke,Dr),opacity:e.strokeOpacity*(e.opacity||1),transparent:!0,side:on,depthWrite:!1})}static createShapes(t){return console.warn("SVGLoader: createShapes() is deprecated. Use shapePath.toShapes() instead."),t.toShapes()}static getStrokeStyle(t,e,n,i,s){return t=t!==void 0?t:1,e=e!==void 0?e:"#000",n=n!==void 0?n:"miter",i=i!==void 0?i:"butt",s=s!==void 0?s:4,{strokeColor:e,strokeWidth:t,strokeLineJoin:n,strokeLineCap:i,strokeMiterLimit:s}}static pointsToStroke(t,e,n,i){let s=[],r=[],a=[];if(Ph.pointsToStrokeWithBuffers(t,e,n,i,s,r,a)===0)return null;let o=new jt;return o.setAttribute("position",new Nt(s,3)),o.setAttribute("normal",new Nt(r,3)),o.setAttribute("uv",new Nt(a,2)),o}static pointsToStrokeWithBuffers(t,e,n,i,s,r,a,o){let l=new rt,c=new rt,h=new rt,d=new rt,u=new rt,f=new rt,g=new rt,_=new rt,m=new rt,p=new rt,b=new rt,S=new rt,x=new rt,T=new rt,I=new rt,R=new rt,y=new rt;n=n!==void 0?n:12,i=i!==void 0?i:0.001,o=o!==void 0?o:0,t=_t(t);let A=t.length;if(A<2)return 0;let Z=t[0].equals(t[A-1]),D,q=t[0],ot,z=e.strokeWidth/2,et=1/(A-1),nt=0,B,ft,Q,ct,dt=!1,Bt=0,Dt=o*3,ee=o*2;Jt(t[0],t[1],l).multiplyScalar(z),_.copy(t[0]).sub(l),m.copy(t[0]).add(l),p.copy(_),b.copy(m);for(let C=1;C<A;C++){if(D=t[C],C===A-1)if(Z)ot=t[1];else ot=void 0;else ot=t[C+1];let H=l;if(Jt(q,D,H),h.copy(H).multiplyScalar(z),S.copy(D).sub(h),x.copy(D).add(h),B=nt+et,ft=!1,ot!==void 0){if(Jt(D,ot,c),h.copy(c).multiplyScalar(z),T.copy(D).sub(h),I.copy(D).add(h),Q=!0,h.subVectors(ot,q),H.dot(h)<0)Q=!1;if(C===1)dt=Q;h.subVectors(ot,D),h.normalize();let U=Math.abs(H.dot(h));if(U>Number.EPSILON){let Y=z/U;h.multiplyScalar(-Y),d.subVectors(D,q),u.copy(d).setLength(Y).add(h),R.copy(u).negate();let tt=u.length(),at=d.length();d.divideScalar(at),f.subVectors(ot,D);let mt=f.length();if(f.divideScalar(mt),d.dot(R)<at&&f.dot(R)<mt)ft=!0;if(y.copy(u).add(D),R.add(D),ft){let gt=Q?m:_,X=(y.x-gt.x)*(R.y-gt.y)-(y.y-gt.y)*(R.x-gt.x);if(Q&&X<0||!Q&&X>0)R.copy(gt)}if(ct=!1,ft)if(Q)I.copy(R),x.copy(R);else T.copy(R),S.copy(R);else N();switch(e.strokeLineJoin){case"bevel":P(Q,ft,B);break;case"round":if(G(Q,ft),Q)O(D,S,T,B,0);else O(D,I,x,B,1);break;case"miter":case"miter-clip":default:let gt=z*e.strokeMiterLimit/tt;if(gt<1)if(e.strokeLineJoin!=="miter-clip"){P(Q,ft,B);break}else if(G(Q,ft),Q)f.subVectors(y,S).multiplyScalar(gt).add(S),g.subVectors(y,T).multiplyScalar(gt).add(T),F(S,B,0),F(f,B,0),F(D,B,0.5),F(D,B,0.5),F(f,B,0),F(g,B,0),F(D,B,0.5),F(g,B,0),F(T,B,0);else f.subVectors(y,x).multiplyScalar(gt).add(x),g.subVectors(y,I).multiplyScalar(gt).add(I),F(x,B,1),F(f,B,1),F(D,B,0.5),F(D,B,0.5),F(f,B,1),F(g,B,1),F(D,B,0.5),F(g,B,1),F(I,B,1);else{if(ft){if(Q)F(m,nt,1),F(_,nt,0),F(y,B,0),F(m,nt,1),F(y,B,0),F(R,B,1);else F(m,nt,1),F(_,nt,0),F(y,B,1),F(_,nt,0),F(R,B,0),F(y,B,1);if(Q)T.copy(y);else I.copy(y)}else if(Q)F(S,B,0),F(y,B,0),F(D,B,0.5),F(D,B,0.5),F(y,B,0),F(T,B,0);else F(x,B,1),F(y,B,1),F(D,B,0.5),F(D,B,0.5),F(y,B,1),F(I,B,1);ct=!0}break}}else N()}else N();if(!Z&&C===A-1)J(t[0],p,b,Q,!0,nt);nt=B,q=D,_.copy(T),m.copy(I)}if(!Z)J(D,S,x,Q,!1,B);else if(ft&&s){let C=y,H=R;if(dt!==Q)C=R,H=y;if(Q){if(ct||dt){if(H.toArray(s,0),H.toArray(s,9),ct)C.toArray(s,3)}}else if(ct||!dt){if(H.toArray(s,3),H.toArray(s,9),ct)C.toArray(s,0)}}if(s){let C=[new rt,new rt,new rt],H=o*3;for(let U=H;U<Dt;U+=9)if(C[0].set(s[U],s[U+1]),C[1].set(s[U+3],s[U+4]),C[2].set(s[U+6],s[U+7]),sn.area(C)<0)s[U+3]=C[0].x,s[U+4]=C[0].y}return Bt;function Jt(C,H,U){return U.subVectors(H,C),U.set(-U.y,U.x).normalize()}function F(C,H,U){if(s){if(s[Dt]=C.x,s[Dt+1]=C.y,s[Dt+2]=0,r)r[Dt]=0,r[Dt+1]=0,r[Dt+2]=1;if(Dt+=3,a)a[ee]=H,a[ee+1]=U,ee+=2}Bt+=3}function O(C,H,U,Y,tt){l.copy(H).sub(C).normalize(),c.copy(U).sub(C).normalize();let at=Math.PI,mt=l.dot(c);if(Math.abs(mt)<1)at=Math.abs(Math.acos(mt));at/=n,h.copy(H);for(let gt=0,X=n-1;gt<X;gt++)d.copy(h).rotateAround(C,at),F(h,Y,tt),F(d,Y,tt),F(C,Y,0.5),h.copy(d);F(h,Y,tt),F(U,Y,tt),F(C,Y,0.5)}function N(){F(m,nt,1),F(_,nt,0),F(S,B,0),F(m,nt,1),F(S,B,0),F(x,B,1)}function P(C,H,U){if(H)if(C)F(m,nt,1),F(_,nt,0),F(S,B,0),F(m,nt,1),F(S,B,0),F(R,B,1),F(S,U,0),F(T,U,0),F(R,U,0.5);else F(m,nt,1),F(_,nt,0),F(x,B,1),F(_,nt,0),F(R,B,0),F(x,B,1),F(x,U,1),F(R,U,0),F(I,U,1);else if(C)F(S,U,0),F(T,U,0),F(D,U,0.5);else F(x,U,1),F(I,U,0),F(D,U,0.5)}function G(C,H){if(H)if(C)F(m,nt,1),F(_,nt,0),F(S,B,0),F(m,nt,1),F(S,B,0),F(R,B,1),F(S,nt,0),F(D,B,0.5),F(R,B,1),F(D,B,0.5),F(T,nt,0),F(R,B,1);else F(m,nt,1),F(_,nt,0),F(x,B,1),F(_,nt,0),F(R,B,0),F(x,B,1),F(x,nt,1),F(R,B,0),F(D,B,0.5),F(D,B,0.5),F(R,B,0),F(I,nt,1)}function J(C,H,U,Y,tt,at){switch(e.strokeLineCap){case"round":if(tt)O(C,U,H,at,0.5);else O(C,H,U,at,0.5);break;case"square":if(tt)if(l.subVectors(H,C),c.set(l.y,-l.x),h.addVectors(l,c).add(C),d.subVectors(c,l).add(C),Y)h.toArray(s,3),d.toArray(s,0),d.toArray(s,9);else h.toArray(s,3),a[7]===1?d.toArray(s,9):h.toArray(s,9),d.toArray(s,0);else{l.subVectors(U,C),c.set(l.y,-l.x),h.addVectors(l,c).add(C),d.subVectors(c,l).add(C);let mt=s.length;if(Y)h.toArray(s,mt-3),d.toArray(s,mt-6),d.toArray(s,mt-12);else d.toArray(s,mt-6),h.toArray(s,mt-3),d.toArray(s,mt-12)}break;case"butt":default:break}}function _t(C){let H=!1;for(let Y=1,tt=C.length-1;Y<tt;Y++)if(C[Y].distanceTo(C[Y+1])<i){H=!0;break}if(!H)return C;let U=[];U.push(C[0]);for(let Y=1,tt=C.length-1;Y<tt;Y++)if(C[Y].distanceTo(C[Y+1])>=i)U.push(C[Y]);return U.push(C[C.length-1]),U}}}var im=/^\s*url\(\s*(?:["']\s*)?#([^)'"\s]+)(?:\s*["'])?\s*\)\s*$/;function fM(t,e,n=256){if(!t||!Array.isArray(t.stops)||t.stops.length===0)return null;let i=e.userData.transform,s=t.gradientUnits==="objectBoundingBox",r=null;if(s){if(r=pM(e,i),r===null)return null}function a(d,u,f){if(f.set(d,u,1),t.gradientTransform)f.applyMatrix3(t.gradientTransform);if(s)f.set(r.minX+f.x*r.width,r.minY+f.y*r.height,1);if(i)f.applyMatrix3(i)}let o=document.createElement("canvas"),l;if(t.type==="linearGradient"){o.width=n,o.height=1;let d=o.getContext("2d"),u=d.createLinearGradient(0,0,n,0);sm(u,t.stops),d.fillStyle=u,d.fillRect(0,0,n,1);let f=new L,g=new L;a(t.x1,t.y1,f),a(t.x2,t.y2,g);let _=g.x-f.x,m=g.y-f.y,p=_*_+m*m||0.00000000000000000001,b=_/p,S=m/p,x=-(b*f.x+S*f.y);l=new Xt().set(b,S,x,0,0,0.5,0,0,1)}else{let{cx:d,cy:u,fx:f,fy:g,r:_}=t;if(t.gradientTransform){let y=new L;y.set(d,u,1).applyMatrix3(t.gradientTransform),d=y.x,u=y.y,y.set(f,g,1).applyMatrix3(t.gradientTransform),f=y.x,g=y.y}if(s)d=r.minX+d*r.width,u=r.minY+u*r.height,f=r.minX+f*r.width,g=r.minY+g*r.height,_=_*Math.sqrt((r.width*r.width+r.height*r.height)/2);if(_<=0)return null;o.width=n,o.height=n;let m=o.getContext("2d"),p=d-_,b=u-_,S=2*_,x=n/S;m.setTransform(x,0,0,x,-p*x,-b*x);let T=m.createRadialGradient(f,g,0,d,u,_);sm(T,t.stops),m.fillStyle=T,m.fillRect(p,b,S,S);let I=i?i.clone().invert():new Xt;l=new Xt().set(1/S,0,-p/S,0,1/S,-b/S,0,0,1).multiply(I)}let c=new to(o);c.colorSpace=Dr,c.flipY=!1,c.matrixAutoUpdate=!1,c.matrix=l;let h=t.spreadMethod==="reflect"?Na:t.spreadMethod==="repeat"?La:Es;return c.wrapS=h,c.wrapT=h,c}function pM(t,e){let n=e?e.clone().invert():null,i=new rt,s=new Rr;for(let r of t.subPaths)for(let a of r.getPoints()){if(i.copy(a),n)i.applyMatrix3(n);s.expandByPoint(i)}if(s.isEmpty())return null;return{minX:s.min.x,minY:s.min.y,width:s.max.x-s.min.x,height:s.max.y-s.min.y}}function sm(t,e){let n=new Pt;for(let i of e){let s=i.color;if(i.opacity<1){n.setStyle(i.color,Dr);let r=/rgb\(([^)]+)\)/.exec(n.getStyle(Dr));if(r)s=`rgba(${r[1]},${i.opacity})`}t.addColorStop(Math.max(0,Math.min(1,i.offset)),s)}}function mM(t,e=!1){let n=t[0].index!==null,i=new Set(Object.keys(t[0].attributes)),s=new Set(Object.keys(t[0].morphAttributes)),r={},a={},o=t[0].morphTargetsRelative,l=new jt,c=0;for(let h=0;h<t.length;++h){let d=t[h],u=0;if(n!==(d.index!==null))return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed with geometry at index "+h+". All geometries must have compatible attributes; make sure index attribute exists among all geometries, or in none of them."),null;for(let f in d.attributes){if(!i.has(f))return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed with geometry at index "+h+'. All geometries must have compatible attributes; make sure "'+f+'" attribute exists among all geometries, or in none of them.'),null;if(r[f]===void 0)r[f]=[];r[f].push(d.attributes[f]),u++}if(u!==i.size)return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed with geometry at index "+h+". Make sure all geometries have the same number of attributes."),null;if(o!==d.morphTargetsRelative)return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed with geometry at index "+h+". .morphTargetsRelative must be consistent throughout all geometries."),null;for(let f in d.morphAttributes){if(!s.has(f))return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed with geometry at index "+h+".  .morphAttributes must be consistent throughout all geometries."),null;if(a[f]===void 0)a[f]=[];a[f].push(d.morphAttributes[f])}if(e){let f;if(n)f=d.index.count;else if(d.attributes.position!==void 0)f=d.attributes.position.count;else return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed with geometry at index "+h+". The geometry must have either an index or a position attribute"),null;l.addGroup(c,f,h),c+=f}}if(n){let h=0,d=[];for(let u=0;u<t.length;++u){let f=t[u].index;for(let g=0;g<f.count;++g)d.push(f.getX(g)+h);h+=t[u].attributes.position.count}l.setIndex(d)}for(let h in r){let d=rm(r[h]);if(!d)return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed while trying to merge the "+h+" attribute."),null;l.setAttribute(h,d)}for(let h in a){let d=a[h][0].length;if(d===0)continue;l.morphAttributes=l.morphAttributes||{},l.morphAttributes[h]=[];for(let u=0;u<d;++u){let f=[];for(let _=0;_<a[h].length;++_)f.push(a[h][_][u]);let g=rm(f);if(!g)return console.error("THREE.BufferGeometryUtils: .mergeGeometries() failed while trying to merge the "+h+" morphAttribute."),null;l.morphAttributes[h].push(g)}}return l}function rm(t){let e,n,i,s=-1,r=0;for(let c=0;c<t.length;++c){let h=t[c];if(e===void 0)e=h.array.constructor;if(e!==h.array.constructor)return console.error("THREE.BufferGeometryUtils: .mergeAttributes() failed. BufferAttribute.array must be of consistent array types across matching attributes."),null;if(n===void 0)n=h.itemSize;if(n!==h.itemSize)return console.error("THREE.BufferGeometryUtils: .mergeAttributes() failed. BufferAttribute.itemSize must be consistent across matching attributes."),null;if(i===void 0)i=h.normalized;if(i!==h.normalized)return console.error("THREE.BufferGeometryUtils: .mergeAttributes() failed. BufferAttribute.normalized must be consistent across matching attributes."),null;if(s===-1)s=h.gpuType;if(s!==h.gpuType)return console.error("THREE.BufferGeometryUtils: .mergeAttributes() failed. BufferAttribute.gpuType must be consistent across matching attributes."),null;r+=h.count*n}let a=new e(r),o=new de(a,n,i),l=0;for(let c=0;c<t.length;++c){let h=t[c];if(h.isInterleavedBufferAttribute){let d=l/n;for(let u=0,f=h.count;u<f;u++)for(let g=0;g<n;g++){let _=h.getComponent(u,g);o.setComponent(u+d,g,_)}}else a.set(h.array,l);l+=h.count*n}if(s!==void 0)o.gpuType=s;return o}class am extends fr{constructor(){super();this.name="RoomEnvironment",this.position.y=-3.5;let t=new An;t.deleteAttribute("uv");let e=new Wi({side:Ve}),n=new Wi,i=new Ar(16777215,900,28,2);i.position.set(0.418,16.199,0.3),this.add(i);let s=new pe(t,e);s.position.set(-0.757,13.219,0.717),s.scale.set(31.713,28.305,28.591),this.add(s);let r=new pr(t,n,6),a=new ae;a.position.set(-10.906,2.009,1.846),a.rotation.set(0,-0.195,0),a.scale.set(2.328,7.905,4.651),a.updateMatrix(),r.setMatrixAt(0,a.matrix),a.position.set(-5.607,-0.754,-0.758),a.rotation.set(0,0.994,0),a.scale.set(1.97,1.534,3.955),a.updateMatrix(),r.setMatrixAt(1,a.matrix),a.position.set(6.167,0.857,7.803),a.rotation.set(0,0.561,0),a.scale.set(3.927,6.285,3.687),a.updateMatrix(),r.setMatrixAt(2,a.matrix),a.position.set(-2.017,0.018,6.124),a.rotation.set(0,0.333,0),a.scale.set(2.002,4.566,2.064),a.updateMatrix(),r.setMatrixAt(3,a.matrix),a.position.set(2.291,-0.756,-2.621),a.rotation.set(0,-0.286,0),a.scale.set(1.546,1.552,1.496),a.updateMatrix(),r.setMatrixAt(4,a.matrix),a.position.set(-2.193,-0.369,-5.547),a.rotation.set(0,0.516,0),a.scale.set(3.875,3.487,2.986),a.updateMatrix(),r.setMatrixAt(5,a.matrix),this.add(r);let o=new pe(t,Us(50));o.position.set(-16.116,14.37,8.208),o.scale.set(0.1,2.428,2.739),this.add(o);let l=new pe(t,Us(50));l.position.set(-16.109,18.021,-8.207),l.scale.set(0.1,2.425,2.751),this.add(l);let c=new pe(t,Us(17));c.position.set(14.904,12.198,-1.832),c.scale.set(0.15,4.265,6.331),this.add(c);let h=new pe(t,Us(43));h.position.set(-0.462,8.89,14.52),h.scale.set(4.38,5.441,0.088),this.add(h);let d=new pe(t,Us(20));d.position.set(3.235,11.486,-12.541),d.scale.set(2.5,2,0.1),this.add(d);let u=new pe(t,Us(100));u.position.set(0,20,0),u.scale.set(1,0.1,1),this.add(u)}dispose(){let t=new Set;this.traverse((e)=>{if(e.isMesh)t.add(e.geometry),t.add(e.material)});for(let e of t)e.dispose()}}function Us(t){return new Sr({color:0,emissive:16777215,emissiveIntensity:t})}var om={type:"change"},Nh={type:"start"},cm={type:"end"},Fo=new Kn,lm=new dn,gM=Math.cos(70*Ac.DEG2RAD),Le=new L,tn=2*Math.PI,me={NONE:-1,ROTATE:0,DOLLY:1,PAN:2,TOUCH_ROTATE:3,TOUCH_PAN:4,TOUCH_DOLLY_PAN:5,TOUCH_DOLLY_ROTATE:6},Lh=0.000001;class hm extends Io{constructor(t,e=null){super(t,e);if(this.state=me.NONE,this.target=new L,this.cursor=new L,this.minDistance=0,this.maxDistance=1/0,this.minZoom=0,this.maxZoom=1/0,this.minTargetRadius=0,this.maxTargetRadius=1/0,this.minPolarAngle=0,this.maxPolarAngle=Math.PI,this.minAzimuthAngle=-1/0,this.maxAzimuthAngle=1/0,this.enableDamping=!1,this.dampingFactor=0.05,this.enableZoom=!0,this.zoomSpeed=1,this.enableRotate=!0,this.rotateSpeed=1,this.keyRotateSpeed=1,this.enablePan=!0,this.panSpeed=1,this.screenSpacePanning=!0,this.keyPanSpeed=7,this.zoomToCursor=!1,this.autoRotate=!1,this.autoRotateSpeed=2,this.keys={LEFT:"ArrowLeft",UP:"ArrowUp",RIGHT:"ArrowRight",BOTTOM:"ArrowDown"},this.mouseButtons={LEFT:di.ROTATE,MIDDLE:di.DOLLY,RIGHT:di.PAN},this.touches={ONE:fi.ROTATE,TWO:fi.DOLLY_PAN},this.target0=this.target.clone(),this.position0=this.object.position.clone(),this.zoom0=this.object.zoom,this._cursorStyle="auto",this._domElementKeyEvents=null,this._lastPosition=new L,this._lastQuaternion=new Ce,this._lastTargetPosition=new L,this._quat=new Ce().setFromUnitVectors(t.up,new L(0,1,0)),this._quatInverse=this._quat.clone().invert(),this._spherical=new wr,this._sphericalDelta=new wr,this._scale=1,this._panOffset=new L,this._rotateStart=new rt,this._rotateEnd=new rt,this._rotateDelta=new rt,this._panStart=new rt,this._panEnd=new rt,this._panDelta=new rt,this._dollyStart=new rt,this._dollyEnd=new rt,this._dollyDelta=new rt,this._dollyDirection=new L,this._mouse=new rt,this._performCursorZoom=!1,this._pointers=[],this._pointerPositions={},this._controlActive=!1,this._onPointerMove=xM.bind(this),this._onPointerDown=_M.bind(this),this._onPointerUp=yM.bind(this),this._onContextMenu=AM.bind(this),this._onMouseWheel=MM.bind(this),this._onKeyDown=bM.bind(this),this._onTouchStart=TM.bind(this),this._onTouchMove=EM.bind(this),this._onMouseDown=vM.bind(this),this._onMouseMove=SM.bind(this),this._interceptControlDown=wM.bind(this),this._interceptControlUp=RM.bind(this),this.domElement!==null)this.connect(this.domElement);this.update()}set cursorStyle(t){if(this._cursorStyle=t,t==="grab")this.domElement.style.cursor="grab";else this.domElement.style.cursor="auto"}get cursorStyle(){return this._cursorStyle}connect(t){super.connect(t),this.domElement.addEventListener("pointerdown",this._onPointerDown),this.domElement.addEventListener("pointercancel",this._onPointerUp),this.domElement.addEventListener("contextmenu",this._onContextMenu),this.domElement.addEventListener("wheel",this._onMouseWheel,{passive:!1}),this.domElement.getRootNode().addEventListener("keydown",this._interceptControlDown,{passive:!0,capture:!0}),this.domElement.style.touchAction="none"}disconnect(){this.state=me.NONE,this.domElement.removeEventListener("pointerdown",this._onPointerDown),this.domElement.ownerDocument.removeEventListener("pointermove",this._onPointerMove),this.domElement.ownerDocument.removeEventListener("pointerup",this._onPointerUp),this.domElement.removeEventListener("pointercancel",this._onPointerUp),this.domElement.removeEventListener("wheel",this._onMouseWheel),this.domElement.removeEventListener("contextmenu",this._onContextMenu),this.stopListenToKeyEvents();let t=this.domElement.getRootNode();t.removeEventListener("keydown",this._interceptControlDown,{capture:!0}),t.removeEventListener("keyup",this._interceptControlUp,{capture:!0}),this._controlActive=!1,this._pointers.length=0,this._pointerPositions={},this.domElement.style.touchAction="",this.domElement.style.cursor="auto"}dispose(){this.disconnect()}getPolarAngle(){return this._spherical.phi}getAzimuthalAngle(){return this._spherical.theta}getDistance(){return this.object.position.distanceTo(this.target)}listenToKeyEvents(t){t.addEventListener("keydown",this._onKeyDown),this._domElementKeyEvents=t}stopListenToKeyEvents(){if(this._domElementKeyEvents!==null)this._domElementKeyEvents.removeEventListener("keydown",this._onKeyDown),this._domElementKeyEvents=null}saveState(){this.target0.copy(this.target),this.position0.copy(this.object.position),this.zoom0=this.object.zoom}reset(){this.target.copy(this.target0),this.object.position.copy(this.position0),this.object.zoom=this.zoom0,this.object.updateProjectionMatrix(),this.dispatchEvent(om),this.update(),this.state=me.NONE}pan(t,e){this._pan(t,e),this.update()}dollyIn(t){this._dollyIn(t),this.update()}dollyOut(t){this._dollyOut(t),this.update()}rotateLeft(t){this._rotateLeft(t),this.update()}rotateUp(t){this._rotateUp(t),this.update()}update(t=null){let e=this.object.position;if(Le.copy(e).sub(this.target),Le.applyQuaternion(this._quat),this._spherical.setFromVector3(Le),this.autoRotate&&this.state===me.NONE)this._rotateLeft(this._getAutoRotationAngle(t));if(this.enableDamping)this._spherical.theta+=this._sphericalDelta.theta*this.dampingFactor,this._spherical.phi+=this._sphericalDelta.phi*this.dampingFactor;else this._spherical.theta+=this._sphericalDelta.theta,this._spherical.phi+=this._sphericalDelta.phi;let n=this.minAzimuthAngle,i=this.maxAzimuthAngle;if(isFinite(n)&&isFinite(i)){if(n<-Math.PI)n+=tn;else if(n>Math.PI)n-=tn;if(i<-Math.PI)i+=tn;else if(i>Math.PI)i-=tn;if(n<=i)this._spherical.theta=Math.max(n,Math.min(i,this._spherical.theta));else this._spherical.theta=this._spherical.theta>(n+i)/2?Math.max(n,this._spherical.theta):Math.min(i,this._spherical.theta)}if(this._spherical.phi=Math.max(this.minPolarAngle,Math.min(this.maxPolarAngle,this._spherical.phi)),this._spherical.makeSafe(),this.enableDamping===!0)this.target.addScaledVector(this._panOffset,this.dampingFactor);else this.target.add(this._panOffset);this.target.sub(this.cursor),this.target.clampLength(this.minTargetRadius,this.maxTargetRadius),this.target.add(this.cursor);let s=!1;if(this.zoomToCursor&&this._performCursorZoom||this.object.isOrthographicCamera)this._spherical.radius=this._clampDistance(this._spherical.radius);else{let r=this._spherical.radius;this._spherical.radius=this._clampDistance(this._spherical.radius*this._scale),s=r!=this._spherical.radius}if(Le.setFromSpherical(this._spherical),Le.applyQuaternion(this._quatInverse),e.copy(this.target).add(Le),this.object.lookAt(this.target),this.enableDamping===!0)this._sphericalDelta.theta*=1-this.dampingFactor,this._sphericalDelta.phi*=1-this.dampingFactor,this._panOffset.multiplyScalar(1-this.dampingFactor);else this._sphericalDelta.set(0,0,0),this._panOffset.set(0,0,0);if(this.zoomToCursor&&this._performCursorZoom){let r=null;if(this.object.isPerspectiveCamera){let a=Le.length();r=this._clampDistance(a*this._scale);let o=a-r;this.object.position.addScaledVector(this._dollyDirection,o),this.object.updateMatrixWorld(),s=!!o}else if(this.object.isOrthographicCamera){let a=new L(this._mouse.x,this._mouse.y,0);a.unproject(this.object);let o=this.object.zoom;this.object.zoom=Math.max(this.minZoom,Math.min(this.maxZoom,this.object.zoom/this._scale)),this.object.updateProjectionMatrix(),s=o!==this.object.zoom;let l=new L(this._mouse.x,this._mouse.y,0);l.unproject(this.object),this.object.position.sub(l).add(a),this.object.updateMatrixWorld(),r=Le.length()}else console.warn("WARNING: OrbitControls.js encountered an unknown camera type - zoom to cursor disabled."),this.zoomToCursor=!1;if(r!==null)if(this.screenSpacePanning)this.target.set(0,0,-1).transformDirection(this.object.matrix).multiplyScalar(r).add(this.object.position);else if(Fo.origin.copy(this.object.position),Fo.direction.set(0,0,-1).transformDirection(this.object.matrix),Math.abs(this.object.up.dot(Fo.direction))<gM)this.object.lookAt(this.target);else lm.setFromNormalAndCoplanarPoint(this.object.up,this.target),Fo.intersectPlane(lm,this.target)}else if(this.object.isOrthographicCamera){let r=this.object.zoom;if(this.object.zoom=Math.max(this.minZoom,Math.min(this.maxZoom,this.object.zoom/this._scale)),r!==this.object.zoom)this.object.updateProjectionMatrix(),s=!0}if(this._scale=1,this._performCursorZoom=!1,s||this._lastPosition.distanceToSquared(this.object.position)>Lh||8*(1-this._lastQuaternion.dot(this.object.quaternion))>Lh||this._lastTargetPosition.distanceToSquared(this.target)>Lh)return this.dispatchEvent(om),this._lastPosition.copy(this.object.position),this._lastQuaternion.copy(this.object.quaternion),this._lastTargetPosition.copy(this.target),!0;return!1}_getAutoRotationAngle(t){if(t!==null)return tn/60*this.autoRotateSpeed*t;else return tn/60/60*this.autoRotateSpeed}_getZoomScale(t){let e=Math.abs(t*0.01);return Math.pow(0.95,this.zoomSpeed*e)}_rotateLeft(t){this._sphericalDelta.theta-=t}_rotateUp(t){this._sphericalDelta.phi-=t}_panLeft(t,e){Le.setFromMatrixColumn(e,0),Le.multiplyScalar(-t),this._panOffset.add(Le)}_panUp(t,e){if(this.screenSpacePanning===!0)Le.setFromMatrixColumn(e,1);else Le.setFromMatrixColumn(e,0),Le.crossVectors(this.object.up,Le);Le.multiplyScalar(t),this._panOffset.add(Le)}_pan(t,e){let n=this.domElement;if(this.object.isPerspectiveCamera){let i=this.object.position;Le.copy(i).sub(this.target);let s=Le.length();s*=Math.tan(this.object.fov/2*Math.PI/180),this._panLeft(2*t*s/n.clientHeight,this.object.matrix),this._panUp(2*e*s/n.clientHeight,this.object.matrix)}else if(this.object.isOrthographicCamera)this._panLeft(t*(this.object.right-this.object.left)/this.object.zoom/n.clientWidth,this.object.matrix),this._panUp(e*(this.object.top-this.object.bottom)/this.object.zoom/n.clientHeight,this.object.matrix);else console.warn("WARNING: OrbitControls.js encountered an unknown camera type - pan disabled."),this.enablePan=!1}_dollyOut(t){if(this.object.isPerspectiveCamera||this.object.isOrthographicCamera)this._scale/=t;else console.warn("WARNING: OrbitControls.js encountered an unknown camera type - dolly/zoom disabled."),this.enableZoom=!1}_dollyIn(t){if(this.object.isPerspectiveCamera||this.object.isOrthographicCamera)this._scale*=t;else console.warn("WARNING: OrbitControls.js encountered an unknown camera type - dolly/zoom disabled."),this.enableZoom=!1}_updateZoomParameters(t,e){if(!this.zoomToCursor)return;this._performCursorZoom=!0;let n=this.domElement.getBoundingClientRect(),i=t-n.left,s=e-n.top,{width:r,height:a}=n;this._mouse.x=i/r*2-1,this._mouse.y=-(s/a)*2+1,this._dollyDirection.set(this._mouse.x,this._mouse.y,1).unproject(this.object).sub(this.object.position).normalize()}_clampDistance(t){return Math.max(this.minDistance,Math.min(this.maxDistance,t))}_handleMouseDownRotate(t){this._rotateStart.set(t.clientX,t.clientY)}_handleMouseDownDolly(t){this._updateZoomParameters(t.clientX,t.clientX),this._dollyStart.set(t.clientX,t.clientY)}_handleMouseDownPan(t){this._panStart.set(t.clientX,t.clientY)}_handleMouseMoveRotate(t){this._rotateEnd.set(t.clientX,t.clientY),this._rotateDelta.subVectors(this._rotateEnd,this._rotateStart).multiplyScalar(this.rotateSpeed);let e=this.domElement;this._rotateLeft(tn*this._rotateDelta.x/e.clientHeight),this._rotateUp(tn*this._rotateDelta.y/e.clientHeight),this._rotateStart.copy(this._rotateEnd),this.update()}_handleMouseMoveDolly(t){if(this._dollyEnd.set(t.clientX,t.clientY),this._dollyDelta.subVectors(this._dollyEnd,this._dollyStart),this._dollyDelta.y>0)this._dollyOut(this._getZoomScale(this._dollyDelta.y));else if(this._dollyDelta.y<0)this._dollyIn(this._getZoomScale(this._dollyDelta.y));this._dollyStart.copy(this._dollyEnd),this.update()}_handleMouseMovePan(t){this._panEnd.set(t.clientX,t.clientY),this._panDelta.subVectors(this._panEnd,this._panStart).multiplyScalar(this.panSpeed),this._pan(this._panDelta.x,this._panDelta.y),this._panStart.copy(this._panEnd),this.update()}_handleMouseWheel(t){if(this._updateZoomParameters(t.clientX,t.clientY),t.deltaY<0)this._dollyIn(this._getZoomScale(t.deltaY));else if(t.deltaY>0)this._dollyOut(this._getZoomScale(t.deltaY));this.update()}_handleKeyDown(t){let e=!1;switch(t.code){case this.keys.UP:if(t.ctrlKey||t.metaKey||t.shiftKey){if(this.enableRotate)this._rotateUp(tn*this.keyRotateSpeed/this.domElement.clientHeight)}else if(this.enablePan)this._pan(0,this.keyPanSpeed);e=!0;break;case this.keys.BOTTOM:if(t.ctrlKey||t.metaKey||t.shiftKey){if(this.enableRotate)this._rotateUp(-tn*this.keyRotateSpeed/this.domElement.clientHeight)}else if(this.enablePan)this._pan(0,-this.keyPanSpeed);e=!0;break;case this.keys.LEFT:if(t.ctrlKey||t.metaKey||t.shiftKey){if(this.enableRotate)this._rotateLeft(tn*this.keyRotateSpeed/this.domElement.clientHeight)}else if(this.enablePan)this._pan(this.keyPanSpeed,0);e=!0;break;case this.keys.RIGHT:if(t.ctrlKey||t.metaKey||t.shiftKey){if(this.enableRotate)this._rotateLeft(-tn*this.keyRotateSpeed/this.domElement.clientHeight)}else if(this.enablePan)this._pan(-this.keyPanSpeed,0);e=!0;break}if(e)t.preventDefault(),this.update()}_handleTouchStartRotate(t){if(this._pointers.length===1)this._rotateStart.set(t.pageX,t.pageY);else{let e=this._getSecondPointerPosition(t),n=0.5*(t.pageX+e.x),i=0.5*(t.pageY+e.y);this._rotateStart.set(n,i)}}_handleTouchStartPan(t){if(this._pointers.length===1)this._panStart.set(t.pageX,t.pageY);else{let e=this._getSecondPointerPosition(t),n=0.5*(t.pageX+e.x),i=0.5*(t.pageY+e.y);this._panStart.set(n,i)}}_handleTouchStartDolly(t){let e=this._getSecondPointerPosition(t),n=t.pageX-e.x,i=t.pageY-e.y,s=Math.sqrt(n*n+i*i);this._dollyStart.set(0,s)}_handleTouchStartDollyPan(t){if(this.enableZoom)this._handleTouchStartDolly(t);if(this.enablePan)this._handleTouchStartPan(t)}_handleTouchStartDollyRotate(t){if(this.enableZoom)this._handleTouchStartDolly(t);if(this.enableRotate)this._handleTouchStartRotate(t)}_handleTouchMoveRotate(t){if(this._pointers.length==1)this._rotateEnd.set(t.pageX,t.pageY);else{let n=this._getSecondPointerPosition(t),i=0.5*(t.pageX+n.x),s=0.5*(t.pageY+n.y);this._rotateEnd.set(i,s)}this._rotateDelta.subVectors(this._rotateEnd,this._rotateStart).multiplyScalar(this.rotateSpeed);let e=this.domElement;this._rotateLeft(tn*this._rotateDelta.x/e.clientHeight),this._rotateUp(tn*this._rotateDelta.y/e.clientHeight),this._rotateStart.copy(this._rotateEnd)}_handleTouchMovePan(t){if(this._pointers.length===1)this._panEnd.set(t.pageX,t.pageY);else{let e=this._getSecondPointerPosition(t),n=0.5*(t.pageX+e.x),i=0.5*(t.pageY+e.y);this._panEnd.set(n,i)}this._panDelta.subVectors(this._panEnd,this._panStart).multiplyScalar(this.panSpeed),this._pan(this._panDelta.x,this._panDelta.y),this._panStart.copy(this._panEnd)}_handleTouchMoveDolly(t){let e=this._getSecondPointerPosition(t),n=t.pageX-e.x,i=t.pageY-e.y,s=Math.sqrt(n*n+i*i);this._dollyEnd.set(0,s),this._dollyDelta.set(0,Math.pow(this._dollyEnd.y/this._dollyStart.y,this.zoomSpeed)),this._dollyOut(this._dollyDelta.y),this._dollyStart.copy(this._dollyEnd);let r=(t.pageX+e.x)*0.5,a=(t.pageY+e.y)*0.5;this._updateZoomParameters(r,a)}_handleTouchMoveDollyPan(t){if(this.enableZoom)this._handleTouchMoveDolly(t);if(this.enablePan)this._handleTouchMovePan(t)}_handleTouchMoveDollyRotate(t){if(this.enableZoom)this._handleTouchMoveDolly(t);if(this.enableRotate)this._handleTouchMoveRotate(t)}_addPointer(t){this._pointers.push(t.pointerId)}_removePointer(t){delete this._pointerPositions[t.pointerId];for(let e=0;e<this._pointers.length;e++)if(this._pointers[e]==t.pointerId){this._pointers.splice(e,1);return}}_isTrackingPointer(t){for(let e=0;e<this._pointers.length;e++)if(this._pointers[e]==t.pointerId)return!0;return!1}_trackPointer(t){let e=this._pointerPositions[t.pointerId];if(e===void 0)e=new rt,this._pointerPositions[t.pointerId]=e;e.set(t.pageX,t.pageY)}_getSecondPointerPosition(t){let e=t.pointerId===this._pointers[0]?this._pointers[1]:this._pointers[0];return this._pointerPositions[e]}_customWheelEvent(t){let e=t.deltaMode,n={clientX:t.clientX,clientY:t.clientY,deltaY:t.deltaY};switch(e){case 1:n.deltaY*=16;break;case 2:n.deltaY*=100;break}if(t.ctrlKey&&!this._controlActive)n.deltaY*=10;return n}}function _M(t){if(this.enabled===!1)return;if(this._pointers.length===0)this.domElement.setPointerCapture(t.pointerId),this.domElement.ownerDocument.addEventListener("pointermove",this._onPointerMove),this.domElement.ownerDocument.addEventListener("pointerup",this._onPointerUp);if(this._isTrackingPointer(t))return;if(this._addPointer(t),t.pointerType==="touch")this._onTouchStart(t);else this._onMouseDown(t);if(this._cursorStyle==="grab")this.domElement.style.cursor="grabbing"}function xM(t){if(this.enabled===!1)return;if(t.pointerType==="touch")this._onTouchMove(t);else this._onMouseMove(t)}function yM(t){switch(this._removePointer(t),this._pointers.length){case 0:if(this.domElement.releasePointerCapture(t.pointerId),this.domElement.ownerDocument.removeEventListener("pointermove",this._onPointerMove),this.domElement.ownerDocument.removeEventListener("pointerup",this._onPointerUp),this.dispatchEvent(cm),this.state=me.NONE,this._cursorStyle==="grab")this.domElement.style.cursor="grab";break;case 1:let e=this._pointers[0],n=this._pointerPositions[e];this._onTouchStart({pointerId:e,pageX:n.x,pageY:n.y});break}}function vM(t){let e;switch(t.button){case 0:e=this.mouseButtons.LEFT;break;case 1:e=this.mouseButtons.MIDDLE;break;case 2:e=this.mouseButtons.RIGHT;break;default:e=-1}switch(e){case di.DOLLY:if(this.enableZoom===!1)return;this._handleMouseDownDolly(t),this.state=me.DOLLY;break;case di.ROTATE:if(t.ctrlKey||t.metaKey||t.shiftKey){if(this.enablePan===!1)return;this._handleMouseDownPan(t),this.state=me.PAN}else{if(this.enableRotate===!1)return;this._handleMouseDownRotate(t),this.state=me.ROTATE}break;case di.PAN:if(t.ctrlKey||t.metaKey||t.shiftKey){if(this.enableRotate===!1)return;this._handleMouseDownRotate(t),this.state=me.ROTATE}else{if(this.enablePan===!1)return;this._handleMouseDownPan(t),this.state=me.PAN}break;default:this.state=me.NONE}if(this.state!==me.NONE)this.dispatchEvent(Nh)}function SM(t){switch(this.state){case me.ROTATE:if(this.enableRotate===!1)return;this._handleMouseMoveRotate(t);break;case me.DOLLY:if(this.enableZoom===!1)return;this._handleMouseMoveDolly(t);break;case me.PAN:if(this.enablePan===!1)return;this._handleMouseMovePan(t);break}}function MM(t){if(this.enabled===!1||this.enableZoom===!1||this.state!==me.NONE)return;t.preventDefault(),this.dispatchEvent(Nh),this._handleMouseWheel(this._customWheelEvent(t)),this.dispatchEvent(cm)}function bM(t){if(this.enabled===!1)return;this._handleKeyDown(t)}function TM(t){switch(this._trackPointer(t),this._pointers.length){case 1:switch(this.touches.ONE){case fi.ROTATE:if(this.enableRotate===!1)return;this._handleTouchStartRotate(t),this.state=me.TOUCH_ROTATE;break;case fi.PAN:if(this.enablePan===!1)return;this._handleTouchStartPan(t),this.state=me.TOUCH_PAN;break;default:this.state=me.NONE}break;case 2:switch(this.touches.TWO){case fi.DOLLY_PAN:if(this.enableZoom===!1&&this.enablePan===!1)return;this._handleTouchStartDollyPan(t),this.state=me.TOUCH_DOLLY_PAN;break;case fi.DOLLY_ROTATE:if(this.enableZoom===!1&&this.enableRotate===!1)return;this._handleTouchStartDollyRotate(t),this.state=me.TOUCH_DOLLY_ROTATE;break;default:this.state=me.NONE}break;default:this.state=me.NONE}if(this.state!==me.NONE)this.dispatchEvent(Nh)}function EM(t){switch(this._trackPointer(t),this.state){case me.TOUCH_ROTATE:if(this.enableRotate===!1)return;this._handleTouchMoveRotate(t),this.update();break;case me.TOUCH_PAN:if(this.enablePan===!1)return;this._handleTouchMovePan(t),this.update();break;case me.TOUCH_DOLLY_PAN:if(this.enableZoom===!1&&this.enablePan===!1)return;this._handleTouchMoveDollyPan(t),this.update();break;case me.TOUCH_DOLLY_ROTATE:if(this.enableZoom===!1&&this.enableRotate===!1)return;this._handleTouchMoveDollyRotate(t),this.update();break;default:this.state=me.NONE}}function AM(t){if(this.enabled===!1)return;t.preventDefault()}function wM(t){if(t.key==="Control")this._controlActive=!0,this.domElement.getRootNode().addEventListener("keyup",this._interceptControlUp,{passive:!0,capture:!0})}function RM(t){if(t.key==="Control")this._controlActive=!1,this.domElement.getRootNode().removeEventListener("keyup",this._interceptControlUp,{passive:!0,capture:!0})}export{Dl as ACESFilmicToneMapping,bs as AddEquation,Dd as AddOperation,zm as AdditiveAnimationBlendMode,wl as AdditiveBlending,Fl as AgXToneMapping,kd as AlphaFormat,jd as AlwaysCompare,Ad as AlwaysDepth,lg as AlwaysStencilFunc,rh as AmbientLight,_h as AnimationAction,ys as AnimationClip,Uf as AnimationLoader,jf as AnimationMixer,Kf as AnimationObjectGroup,Pf as AnimationUtils,Oc as ArcCurve,uh as ArrayCamera,yp as ArrowHelper,vm as AttachedBindMode,fh as Audio,Jf as AudioAnalyser,Co as AudioContext,Yf as AudioListener,Xf as AudioLoader,vp as AxesHelper,Ve as BackSide,km as BasicDepthPacking,xm as BasicShadowMap,Pc as BatchedMesh,th as BezierInterpolant,$a as Bone,_i as BooleanKeyframeTrack,Rr as Box2,Oe as Box3,_p as Box3Helper,An as BoxGeometry,gp as BoxHelper,de as BufferAttribute,jt as BufferGeometry,ch as BufferGeometryLoader,Fd as ByteType,Pn as Cache,Er as Camera,mp as CameraHelper,to as CanvasTexture,no as CapsuleGeometry,zc as CatmullRomCurve3,Nl as CineonToneMapping,io as CircleGeometry,Es as ClampToEdgeWrapping,sp as Clock,Pt as Color,To as ColorKeyframeTrack,re as ColorManagement,bg as Compatibility,vf as CompressedArrayTexture,Sf as CompressedCubeTexture,mr as CompressedTexture,Ff as CompressedTextureLoader,_r as ConeGeometry,bd as ConstantAlphaFactor,Sd as ConstantColorFactor,Io as Controls,hh as CubeCamera,Uc as CubeDepthTexture,Ts as CubeReflectionMapping,Fi as CubeRefractionMapping,Rs as CubeTexture,Of as CubeTextureLoader,ar as CubeUVReflectionMapping,ro as CubicBezierCurve,kc as CubicBezierCurve3,jc as CubicInterpolant,Al as CullFaceBack,id as CullFaceFront,_m as CullFaceFrontBack,nd as CullFaceNone,gn as Curve,Hc as CurvePath,rd as CustomBlending,Ul as CustomToneMapping,gr as CylinderGeometry,rp as Cylindrical,hr as Data3DTexture,cr as DataArrayTexture,an as DataTexture,Bf as DataTextureLoader,cf as DataUtils,Km as DecrementStencilOp,Qm as DecrementWrapStencilOp,Nf as DefaultLoadingManager,Bi as DepthFormat,zi as DepthStencilFormat,Hi as DepthTexture,Sm as DetachedBindMode,sh as DirectionalLight,pp as DirectionalLightHelper,Qc as DiscreteInterpolant,so as DodecahedronGeometry,on as DoubleSide,gd as DstAlphaFactor,xd as DstColorFactor,gg as DynamicCopyUsage,hg as DynamicDrawUsage,fg as DynamicReadUsage,Fc as EdgesGeometry,xr as EllipseCurve,Jd as EqualCompare,Rd as EqualDepth,ig as EqualStencilFunc,Ia as EquirectangularReflectionMapping,Pa as EquirectangularRefractionMapping,Mn as Euler,pn as EventDispatcher,eo as ExternalTexture,ho as ExtrudeGeometry,fn as FileLoader,mf as Float16BufferAttribute,Nt as Float32BufferAttribute,$n as FloatType,qa as Fog,Xa as FogExp2,yf as FramebufferTexture,Ms as FrontSide,ui as Frustum,ja as FrustumArray,np as GLBufferAttribute,xg as GLSL1,bc as GLSL3,$d as GreaterCompare,Id as GreaterDepth,Va as GreaterEqualCompare,Cd as GreaterEqualDepth,og as GreaterEqualStencilFunc,rg as GreaterStencilFunc,dp as GridHelper,Li as Group,Mf as HTMLTexture,Nn as HalfFloatType,nh as HemisphereLight,up as HemisphereLightHelper,uo as IcosahedronGeometry,Wf as ImageBitmapLoader,vs as ImageLoader,wc as ImageUtils,$m as IncrementStencilOp,jm as IncrementWrapStencilOp,hi as InstancedBufferAttribute,lh as InstancedBufferGeometry,ep as InstancedInterleavedBuffer,pr as InstancedMesh,ff as Int16BufferAttribute,pf as Int32BufferAttribute,hf as Int8BufferAttribute,Bl as IntType,ws as InterleavedBuffer,ci as InterleavedBufferAttribute,Xi as Interpolant,Dm as InterpolateBezier,Pm as InterpolateDiscrete,Lm as InterpolateLinear,Nm as InterpolateSmooth,Mg as InterpolationSamplingMode,Sg as InterpolationSamplingType,tg as InvertStencilOp,Zm as KeepStencilOp,cn as KeyframeTrack,Cc as LOD,fo as LatheGeometry,ur as Layers,Zd as LessCompare,wd as LessDepth,Ha as LessEqualCompare,Il as LessEqualDepth,sg as LessEqualStencilFunc,ng as LessStencilFunc,jn as Light,oh as LightProbe,Tr as LightShadow,Zn as Line,ap as Line3,We as LineBasicMaterial,ao as LineCurve,Gc as LineCurve3,Kc as LineDashedMaterial,Lc as LineLoop,En as LineSegments,Qe as LinearFilter,bo as LinearInterpolant,Am as LinearMipMapLinearFilter,Em as LinearMipMapNearestFilter,Oi as LinearMipmapLinearFilter,Da as LinearMipmapNearestFilter,Sc as LinearSRGBColorSpace,Pl as LinearToneMapping,Mc as LinearTransfer,Be as Loader,Ca as LoaderUtils,Ao as LoadingManager,Rm as LoopOnce,Im as LoopPingPong,Cm as LoopRepeat,di as MOUSE,Fe as Material,ym as MaterialBlending,Ro as MaterialLoader,Ac as MathUtils,yh as Matrix2,Xt as Matrix3,Kt as Matrix4,cd as MaxEquation,pe as Mesh,mn as MeshBasicMaterial,So as MeshDepthMaterial,Mo as MeshDistanceMaterial,Sr as MeshLambertMaterial,$c as MeshMatcapMaterial,Jc as MeshNormalMaterial,Yc as MeshPhongMaterial,qc as MeshPhysicalMaterial,Wi as MeshStandardMaterial,Zc as MeshToonMaterial,ld as MinEquation,Na as MirroredRepeatWrapping,Nd as MixOperation,Cl as MultiplyBlending,Ld as MultiplyOperation,pi as NearestFilter,Tm as NearestMipMapLinearFilter,bm as NearestMipMapNearestFilter,or as NearestMipmapLinearFilter,Ud as NearestMipmapNearestFilter,Ol as NeutralToneMapping,Yd as NeverCompare,Ed as NeverDepth,eg as NeverStencilFunc,Ln as NoBlending,Gi as NoColorSpace,Wm as NoNormalPacking,bn as NoToneMapping,Bm as NormalAnimationBlendMode,rr as NormalBlending,qm as NormalGAPacking,Xm as NormalRGPacking,Kd as NotEqualCompare,Pd as NotEqualDepth,ag as NotEqualStencilFunc,Mr as NumberKeyframeTrack,ae as Object3D,Vf as ObjectLoader,qd as ObjectSpaceNormalMap,yr as OctahedronGeometry,ud as OneFactor,Td as OneMinusConstantAlphaFactor,Md as OneMinusConstantColorFactor,_d as OneMinusDstAlphaFactor,yd as OneMinusDstColorFactor,md as OneMinusSrcAlphaFactor,fd as OneMinusSrcColorFactor,hm as OrbitControls,Ps as OrthographicCamera,sr as PCFShadowMap,sd as PCFSoftShadowMap,Ah as PMREMGenerator,Jn as Path,De as PerspectiveCamera,dn as Plane,Is as PlaneGeometry,xp as PlaneHelper,Ar as PointLight,hp as PointLightHelper,Nc as Points,Qa as PointsMaterial,fp as PolarGridHelper,gi as PolyhedronGeometry,Zf as PositionalAudio,le as PropertyBinding,ph as PropertyMixer,oo as QuadraticBezierCurve,lo as QuadraticBezierCurve3,Ce as Quaternion,br as QuaternionKeyframeTrack,eh as QuaternionLinearInterpolant,Kl as R11_EAC_Format,ka as RED_GREEN_RGTC2_Format,_c as RED_RGTC1_Format,ed as REVISION,za as RG11_EAC_Format,Gm as RGBADepthPacking,Dn as RGBAFormat,Vl as RGBAIntegerFormat,uc as RGBA_ASTC_10x10_Format,lc as RGBA_ASTC_10x5_Format,cc as RGBA_ASTC_10x6_Format,hc as RGBA_ASTC_10x8_Format,dc as RGBA_ASTC_12x10_Format,fc as RGBA_ASTC_12x12_Format,tc as RGBA_ASTC_4x4_Format,ec as RGBA_ASTC_5x4_Format,nc as RGBA_ASTC_5x5_Format,ic as RGBA_ASTC_6x5_Format,sc as RGBA_ASTC_6x6_Format,rc as RGBA_ASTC_8x5_Format,ac as RGBA_ASTC_8x6_Format,oc as RGBA_ASTC_8x8_Format,pc as RGBA_BPTC_Format,$l as RGBA_ETC2_EAC_Format,Yl as RGBA_PVRTC_2BPPV1_Format,ql as RGBA_PVRTC_4BPPV1_Format,Fa as RGBA_S3TC_DXT1_Format,Oa as RGBA_S3TC_DXT3_Format,Ba as RGBA_S3TC_DXT5_Format,Hm as RGBDepthPacking,Gd as RGBFormat,wm as RGBIntegerFormat,mc as RGB_BPTC_SIGNED_Format,gc as RGB_BPTC_UNSIGNED_Format,Zl as RGB_ETC1_Format,Jl as RGB_ETC2_Format,Xl as RGB_PVRTC_2BPPV1_Format,Wl as RGB_PVRTC_4BPPV1_Format,Ua as RGB_S3TC_DXT1_Format,Vm as RGDepthPacking,ki as RGFormat,Hl as RGIntegerFormat,vo as RawShaderMaterial,Kn as Ray,ip as Raycaster,ah as RectAreaLight,Hd as RedFormat,Gl as RedIntegerFormat,Ll as ReinhardToneMapping,Tg as RenderObjectRefreshType,Wa as RenderTarget,Qf as RenderTarget3D,La as RepeatWrapping,Jm as ReplaceStencilOp,od as ReverseSubtractEquation,po as RingGeometry,am as RoomEnvironment,Ih as RoundedBoxGeometry,jl as SIGNED_R11_EAC_Format,yc as SIGNED_RED_GREEN_RGTC2_Format,xc as SIGNED_RED_RGTC1_Format,Ql as SIGNED_RG11_EAC_Format,Ga as SRGBColorSpace,xe as SRGBTransfer,Ph as SVGLoader,fr as Scene,ie as ShaderChunk,On as ShaderLib,ln as ShaderMaterial,Wc as ShadowMaterial,Cs as Shape,mo as ShapeGeometry,Un as ShapePath,sn as ShapeUtils,Od as ShortType,Ka as Skeleton,lp as SkeletonHelper,Ic as SkinnedMesh,rf as Source,Ue as Sphere,vr as SphereGeometry,wr as Spherical,wo as SphericalHarmonics3,co as SplineCurve,ih as SpotLight,op as SpotLightHelper,Rc as Sprite,Ja as SpriteMaterial,pd as SrcAlphaFactor,vd as SrcAlphaSaturateFactor,dd as SrcColorFactor,mg as StaticCopyUsage,cg as StaticDrawUsage,dg as StaticReadUsage,qf as StereoCamera,_g as StreamCopyUsage,ug as StreamDrawUsage,pg as StreamReadUsage,xi as StringKeyframeTrack,ad as SubtractEquation,Rl as SubtractiveBlending,fi as TOUCH,vc as TangentSpaceNormalMap,go as TetrahedronGeometry,be as Texture,zf as TextureLoader,In as TextureSource,Sp as TextureUtils,dh as Timer,vg as TimestampQuery,_o as TorusGeometry,xo as TorusKnotGeometry,je as Triangle,Xd as TriangleFanDrawMode,Wd as TriangleStripDrawMode,Vd as TrianglesDrawMode,yo as TubeGeometry,Mm as UVMapping,Ya as Uint16BufferAttribute,Za as Uint32BufferAttribute,uf as Uint8BufferAttribute,df as Uint8ClampedBufferAttribute,xh as Uniform,tp as UniformsGroup,At as UniformsLib,Rf as UniformsUtils,Tn as UnsignedByteType,zd as UnsignedInt101111Type,As as UnsignedInt248Type,Bd as UnsignedInt5999Type,mi as UnsignedIntType,zl as UnsignedShort4444Type,kl as UnsignedShort5551Type,lr as UnsignedShortType,Ss as VSMShadowMap,rt as Vector2,L as Vector3,fe as Vector4,Eo as VectorKeyframeTrack,xf as VideoFrameTexture,Dc as VideoTexture,of as WebGL3DRenderTarget,af as WebGLArrayRenderTarget,Tc as WebGLCoordinateSystem,Ch as WebGLCubeRenderTarget,Je as WebGLRenderTarget,dM as WebGLRenderer,sM as WebGLUtils,yg as WebGPUCoordinateSystem,dr as WebXRController,Vc as WireframeGeometry,Om as WrapAroundEnding,Um as ZeroCurvatureEnding,hd as ZeroFactor,Fm as ZeroSlopeEnding,Ym as ZeroStencilOp,tf as createCanvasElement,Wt as error,Rg as getConsoleFunction,tr as log,mM as mergeGeometries,wg as setConsoleFunction,Mt as warn,qn as warnOnce};
