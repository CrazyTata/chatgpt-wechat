/*! 
 Build based on chatgpt-wechat-admin 
 Time : 1685686675000 */
import{f as a}from"./sysDictionary-fa38e7c5.js";import{cn as t,r as i}from"./index-c0e4adfc.js";const s=t("dictionary",(()=>{const t=i({}),s=a=>{t.value={...t.value,...a}};return{dictionaryMap:t,setDictionaryMap:s,getDictionary:async i=>{if(t.value[i]&&t.value[i].length)return t.value[i];{const e=await a({type:i});if(0===e.code){const a={},r=[];return e.data.resysDictionary.sysDictionaryDetails&&e.data.resysDictionary.sysDictionaryDetails.forEach((a=>{r.push({label:a.label,value:a.value})})),a[e.data.resysDictionary.type]=r,s(a),t.value[i]}}}}}));export{s as u};
