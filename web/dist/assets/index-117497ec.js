/*! 
 Build based on chatgpt-wechat-admin 
 Time : 1685686675000 */
import{as as e}from"./index-c0e4adfc.js";const t=100,o=600,n={beforeMount(n,a){const s=a.value,{interval:r=t,delay:d=o}=e(s)?{}:s;let i,l;const u=()=>e(s)?s():s.handler(),v=()=>{l&&(clearTimeout(l),l=void 0),i&&(clearInterval(i),i=void 0)};n.addEventListener("mousedown",(e=>{0===e.button&&(v(),u(),document.addEventListener("mouseup",(()=>v()),{once:!0}),l=setTimeout((()=>{i=setInterval((()=>{u()}),r)}),d))}))}};export{n as v};
