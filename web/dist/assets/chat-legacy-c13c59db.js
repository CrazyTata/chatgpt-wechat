/*! 
 Build based on chatgpt-wechat-admin 
 Time : 1685686675000 */
System.register(["./index-legacy-d5840ffd.js"],(function(t,e){"use strict";var r;return{setters:[function(t){r=t.v}],execute:function(){t("c",(function(t){return r({url:"/chat/createChat",method:"post",data:t})})),t("u",(function(t){return r({url:"/chat/updateChat",method:"put",data:t})})),t("f",(function(t){return r({url:"/chat/findChat",method:"get",params:t})})),t("g",(function(t){return r({url:"/chat/getChatList",method:"get",params:t})})),t("e",(function(t){return r({url:"/chat/exportChatList",method:"get",params:t})}))}}}));
