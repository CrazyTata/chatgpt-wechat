/*! 
 Build based on chatgpt-wechat-admin 
 Time : 1685686675000 */
System.register([],(function(t,e){"use strict";return{execute:function(){t("f",(function(t,e){var s=new Date(t).Format("yyyy-MM-dd hh:mm:ss");return e&&(s=new Date(t).Format(e)),s.toLocaleString()})),Date.prototype.Format=function(t){var e={"M+":this.getMonth()+1,"d+":this.getDate(),"h+":this.getHours(),"m+":this.getMinutes(),"s+":this.getSeconds(),"q+":Math.floor((this.getMonth()+3)/3),S:this.getMilliseconds()};for(var s in/(y+)/.test(t)&&(t=t.replace(RegExp.$1,(this.getFullYear()+"").substr(4-RegExp.$1.length))),e)new RegExp("("+s+")").test(t)&&(t=t.replace(RegExp.$1,1===RegExp.$1.length?e[s]:("00"+e[s]).substr((""+e[s]).length)));return t}}}}));
