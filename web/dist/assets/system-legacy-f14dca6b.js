/*! 
 Build based on chatgpt-wechat-admin 
 Time : 1685686675000 */
System.register(["./index-legacy-d5840ffd.js"],(function(t,e){"use strict";var n;return{setters:[function(t){n=t.v}],execute:function(){t("a",(function(){return n({url:"/system/getSystemConfig",method:"post"})})),t("s",(function(t){return n({url:"/system/setSystemConfig",method:"post",data:t})})),t("g",(function(){return n({url:"/system/getServerInfo",method:"post",donNotShowLoading:!0})}))}}}));
