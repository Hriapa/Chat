"use strict";(self["webpackChunkwebserver"]=self["webpackChunkwebserver"]||[]).push([[467],{9467:function(e,t,s){s.r(t),s.d(t,{default:function(){return Fe}});var a=s(6768),r=s(4232);const n=(0,a.Lk)("h1",null,"Testing Protocol",-1),o={class:"container-test"},i={class:"test-button-block"},l=["disabled"],d=["disabled"],c=["disabled"];function u(e,t,s,u,m,g){const h=(0,a.g2)("ControlTest"),p=(0,a.g2)("DataTest"),k=(0,a.g2)("AcknowledgeTest");return(0,a.uX)(),(0,a.CE)(a.FK,null,[n,(0,a.Lk)("div",o,[(0,a.Lk)("div",i,[(0,a.Lk)("button",{disabled:m.isControlButtonEnable,class:(0,r.C4)({active:m.isControlButtonEnable}),onClick:t[0]||(t[0]=e=>g.pressControlButton())},"Control Test",10,l),(0,a.Lk)("button",{disabled:m.isDataButtonEnable,class:(0,r.C4)({active:m.isDataButtonEnable}),onClick:t[1]||(t[1]=e=>g.pressDataButton())},"Data Test",10,d),(0,a.Lk)("button",{disabled:m.isAcknowledgeEnable,class:(0,r.C4)({active:m.isAcknowledgeEnable}),onClick:t[2]||(t[2]=e=>g.pressAcknowledgeButton())},"Acknowelege Test",10,c)]),m.isControlButtonEnable?((0,a.uX)(),(0,a.Wv)(h,{key:0,CommandType:m.CommandType,UserId:m.UserId,UserName:m.UserName,UserParameter:m.UserParameter,UsersLlist:m.UsersLlist,codingRes:m.codingRes,origin:m.origin,onRunTest:t[3]||(t[3]=e=>g.testRequest(e))},null,8,["CommandType","UserId","UserName","UserParameter","UsersLlist","codingRes","origin"])):(0,a.Q3)("",!0),m.isDataButtonEnable?((0,a.uX)(),(0,a.Wv)(p,{key:1,DataType:m.DataType,DataFormat:m.DataFormat,IndexNumber:m.IndexNumber,UserId:m.UserId,RoomId:m.RoomId,Message:m.Message,codingRes:m.codingRes,origin:m.origin,onRunTest:t[4]||(t[4]=e=>g.testRequest(e))},null,8,["DataType","DataFormat","IndexNumber","UserId","RoomId","Message","codingRes","origin"])):(0,a.Q3)("",!0),m.isAcknowledgeEnable?((0,a.uX)(),(0,a.Wv)(k,{key:2,AckType:m.AcknowledgeType,IndexNumber:m.IndexNumber,UserId:m.UserId,RoomId:m.RoomId,codingRes:m.codingRes,origin:m.origin,onRunTest:t[5]||(t[5]=e=>g.testRequest(e))},null,8,["AckType","IndexNumber","UserId","RoomId","codingRes","origin"])):(0,a.Q3)("",!0)])],64)}s(4114),s(6573),s(8100),s(7936),s(7467),s(4732),s(9577);var m=s(5130);const g={class:"test-field"},h=(0,a.Fv)('<option disabled value="0">Select the test</option><option value="3">Connect</option><option value="4">Disconnect</option><option value="5">Register</option><option value="6">User Update</option><option value="7">User Info</option><option value="8">Users List</option>',7),p=[h],k=["disabled"],U={class:"test-decor"},v=(0,a.Lk)("h2",null," Command type test: ",-1),L={class:"test-decor"},I=(0,a.Lk)("h2",null,"User Id tets:",-1),y={class:"test-decor"},b=(0,a.Lk)("h2",null," User Name test:",-1),C={class:"test-decor"},T=(0,a.Lk)("h2",null,"User Parsmeter test:",-1),R={class:"test-decor"},N=(0,a.Lk)("h2",null," Users List Test:",-1),D={class:"users-list"},f={class:"item-in-list"},_={class:"test-decor"},F=(0,a.Lk)("h2",null," Coding result: ",-1);function A(e,t,s,n,o,i){return(0,a.uX)(),(0,a.CE)("div",g,[(0,a.bo)((0,a.Lk)("select",{"onUpdate:modelValue":t[0]||(t[0]=e=>o.test_num=e)},p,512),[[m.u1,o.test_num]]),(0,a.Lk)("button",{disabled:0==this.test_num,class:(0,r.C4)({active:0!=this.test_num}),onClick:t[1]||(t[1]=e=>i.testRequest())},"Run Test",10,k),(0,a.Lk)("div",U,[v,(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.CommandType})},"Command: "+(0,r.v_)(s.CommandType),3)]),(0,a.Lk)("div",L,[I,(0,a.Lk)("p",{class:(0,r.C4)({checked:0!=s.UserId})},"User Id: "+(0,r.v_)(s.UserId),3)]),(0,a.Lk)("div",y,[b,(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.UserName})}," User name: "+(0,r.v_)(s.UserName),3)]),(0,a.Lk)("div",C,[T,(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.UserParameter.Name})},"Name: "+(0,r.v_)(s.UserParameter.Name),3),(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.UserParameter.Surname})},"Surname: "+(0,r.v_)(s.UserParameter.Surname),3),(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.UserParameter.Familyname})},"Familyname: "+(0,r.v_)(s.UserParameter.Familyname),3),(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.UserParameter.Birthdate})},"Birthdate: "+(0,r.v_)(s.UserParameter.Birthdate),3)]),(0,a.Lk)("div",R,[N,(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.UsersLlist.Type})}," Command: "+(0,r.v_)(s.UsersLlist.Type),3),(0,a.Lk)("p",null," Number Of Items: "+(0,r.v_)(s.UsersLlist.NumOfItems),1),(0,a.Lk)("ul",D,[((0,a.uX)(!0),(0,a.CE)(a.FK,null,(0,a.pI)(s.UsersLlist.Items,(e=>((0,a.uX)(),(0,a.CE)("li",{v:"",key:e.id},[(0,a.Lk)("div",f,[(0,a.Lk)("p",null,(0,r.v_)(e.name),1),(0,a.Lk)("p",{class:(0,r.C4)({online:e.online,offline:!e.online})},(0,r.v_)(i.onlineStatus(e.online)),3),(0,a.Lk)("p",null,(0,r.v_)(e.nrm),1)])])))),128))])]),(0,a.Lk)("div",_,[F,(0,a.Lk)("p",null,"Recieve message: "+(0,r.v_)(s.origin),1),(0,a.Lk)("p",null,"Coding result: "+(0,r.v_)(s.codingRes),1)])])}var M={props:{CommandType:{type:String,default:""},UserId:{type:Number,default:0},UserName:{type:String,default:""},UserParameter:{type:Object,default:()=>({Name:"",Surname:"",Familyname:"",Birthdate:""})},UsersLlist:{type:Object,default:()=>({Type:"",NumOfItems:0,Items:[]})},codingRes:{type:Array,default:()=>[]},origin:{type:Array,default:()=>[]}},data(){return{test_num:0}},methods:{testRequest(){return this.$emit("runTest",this.test_num)},onlineStatus(e){return e?"online":"offline"}}},P=s(1241);const x=(0,P.A)(M,[["render",A]]);var B=x;const E={class:"test-field"},S=(0,a.Fv)('<option disabled value="0">Select the test</option><option value="1">Small text</option><option value="2">Big text</option><option value="3">First fragment</option><option value="4">Last ftagment</option><option value="5">Image</option>',6),w=[S],O=["disabled"],X={class:"test-decor"},q=(0,a.Lk)("h2",null,"Data Type test:",-1),Q={class:"test-decor"},j=(0,a.Lk)("h2",null,"Data Format test:",-1),$={class:"test-decor"},V=(0,a.Lk)("h2",null,"Index Number test:",-1),W={class:"test-decor"},K=(0,a.Lk)("h2",null,"User Id test:",-1),z={class:"test-decor"},G=(0,a.Lk)("h2",null,"Room Id test:",-1),H={class:"test-decor"},J=(0,a.Lk)("h2",null,"Message Data test:",-1),Y=["src"],Z={class:"test-decor"},ee=(0,a.Lk)("h2",null,"Coding result:",-1);function te(e,t,s,n,o,i){return(0,a.uX)(),(0,a.CE)("div",E,[(0,a.bo)((0,a.Lk)("select",{"onUpdate:modelValue":t[0]||(t[0]=e=>o.test_num=e)},w,512),[[m.u1,o.test_num]]),(0,a.Lk)("button",{disabled:0==this.test_num,class:(0,r.C4)({active:0!=this.test_num}),onClick:t[1]||(t[1]=e=>i.testRequest())},"Run Test",10,O),(0,a.Lk)("div",X,[q,(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.DataType})},"Type: "+(0,r.v_)(s.DataType),3)]),(0,a.Lk)("div",Q,[j,(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.DataFormat})},"Format: "+(0,r.v_)(s.DataFormat),3)]),(0,a.Lk)("div",$,[V,(0,a.Lk)("p",{class:(0,r.C4)({checked:0!=s.IndexNumber})},"Index Number: "+(0,r.v_)(s.IndexNumber),3)]),(0,a.Lk)("div",W,[K,(0,a.Lk)("p",{class:(0,r.C4)({checked:0!=s.UserId})},"User Id: "+(0,r.v_)(s.UserId),3)]),(0,a.Lk)("div",z,[G,(0,a.Lk)("p",{class:(0,r.C4)({checked:0==isNaN(s.RoomId)})},"Room Id: "+(0,r.v_)(s.RoomId),3)]),(0,a.Lk)("div",H,[J,(0,a.Lk)("p",{class:(0,r.C4)({checked:1==s.Message.Fragmentation})},"Fragmentation: "+(0,r.v_)(s.Message.Fragmentation),3),(0,a.Lk)("p",{class:(0,r.C4)({checked:1==s.Message.Fragmentation})},"Fragment Type: "+(0,r.v_)(s.Message.FragmentType),3),(0,a.Lk)("p",{class:(0,r.C4)({checked:1==s.Message.Fragmentation})},"Counter: "+(0,r.v_)(s.Message.Counter),3),"text"==s.DataFormat?((0,a.uX)(),(0,a.CE)("p",{key:0,class:(0,r.C4)({checked:""!=i.dataTextProcessing(s.Message.Data)})},"Text: "+(0,r.v_)(i.dataTextProcessing(s.Message.Data)),3)):(0,a.Q3)("",!0),"image"==s.DataFormat?((0,a.uX)(),(0,a.CE)("img",{key:1,src:i.dataImageProcessing(s.Message.Data),alt:""},null,8,Y)):(0,a.Q3)("",!0)]),(0,a.Lk)("div",Z,[ee,(0,a.Lk)("p",null,"Recieve message: "+(0,r.v_)(s.origin),1),(0,a.Lk)("p",null,"Coding result: "+(0,r.v_)(s.codingRes),1)])])}s(4979);var se={props:{DataType:{type:String,default:""},DataFormat:{type:String,default:""},IndexNumber:{type:Number,default:0},UserId:{type:Number,default:0},RoomId:{type:Number,default:NaN},Message:{type:Object,default:()=>({Fragmentation:!1,FragmentType:"",Counter:0,Data:Uint8Array})},codingRes:{type:Array,default:()=>[]},origin:{type:Array,default:()=>[]}},data(){return{test_num:0}},methods:{testRequest(){return this.$emit("runTest",this.test_num)},dataTextProcessing(e){let t=new TextDecoder;return t.decode(e)},dataImageProcessing(e){let t=btoa(String.fromCharCode(...e));return`data:image/jpeg;base64,${t}`}}};const ae=(0,P.A)(se,[["render",te]]);var re=ae;const ne={class:"test-field"},oe=(0,a.Lk)("option",{disabled:"",value:0},"Select the test",-1),ie=(0,a.Lk)("option",{value:1},"Send Result Test",-1),le=(0,a.Lk)("option",{value:2},"Recieve Result Test",-1),de=(0,a.Lk)("option",{value:3},"Read Result Test",-1),ce=[oe,ie,le,de],ue=["disabled"],me={class:"test-decor"},ge=(0,a.Lk)("h2",null,"Data Type test:",-1),he={class:"test-decor"},pe=(0,a.Lk)("h2",null,"Index Number test:",-1),ke={class:"test-decor"},Ue=(0,a.Lk)("h2",null,"User Id test:",-1),ve={class:"test-decor"},Le=(0,a.Lk)("h2",null,"Room Id test:",-1),Ie={class:"test-decor"},ye=(0,a.Lk)("h2",null,"Coding result:",-1);function be(e,t,s,n,o,i){return(0,a.uX)(),(0,a.CE)("div",ne,[(0,a.bo)((0,a.Lk)("select",{"onUpdate:modelValue":t[0]||(t[0]=e=>o.test_num=e)},ce,512),[[m.u1,o.test_num]]),(0,a.Lk)("button",{disabled:0==this.test_num,class:(0,r.C4)({active:0!=this.test_num}),onClick:t[1]||(t[1]=e=>i.testRequest())},"Run Test",10,ue),(0,a.Lk)("div",me,[ge,(0,a.Lk)("p",{class:(0,r.C4)({checked:""!=s.AckType})},"Type: "+(0,r.v_)(s.AckType),3)]),(0,a.Lk)("div",he,[pe,(0,a.Lk)("p",{class:(0,r.C4)({checked:0!=s.IndexNumber})},"Index Number: "+(0,r.v_)(s.IndexNumber),3)]),(0,a.Lk)("div",ke,[Ue,(0,a.Lk)("p",{class:(0,r.C4)({checked:0!=s.UserId})},"User Id: "+(0,r.v_)(s.UserId),3)]),(0,a.Lk)("div",ve,[Le,(0,a.Lk)("p",{class:(0,r.C4)({checked:0==isNaN(s.RoomId)})},"Room Id: "+(0,r.v_)(s.RoomId),3)]),(0,a.Lk)("div",Ie,[ye,(0,a.Lk)("p",null,"Recieve message: "+(0,r.v_)(s.origin),1),(0,a.Lk)("p",null,"Coding result: "+(0,r.v_)(s.codingRes),1)])])}var Ce={props:{AckType:{type:String,default:""},IndexNumber:{type:Number,default:0},UserId:{type:Number,default:0},RoomId:{type:Number,default:NaN},codingRes:{type:Array,default:()=>[]},origin:{type:Array,default:()=>[]}},data(){return{test_num:0}},methods:{testRequest(){return this.$emit("runTest",this.test_num)}}};const Te=(0,P.A)(Ce,[["render",be]]);var Re=Te,Ne=s(7416),De=s(4652),fe={components:{ControlTest:B,DataTest:re,AcknowledgeTest:Re},data(){return{isControlButtonEnable:!0,isDataButtonEnable:!1,isAcknowledgeEnable:!1,testMode:0,CommandType:"",DataType:"",DataFormat:"",AcknowledgeType:"",IndexNumber:0,UserId:0,RoomId:0,UserName:"",UserParameter:{Name:"",Surname:"",Familyname:"",Birthdate:""},UsersLlist:{Type:"",NumOfItems:0,Items:[]},Message:{Fragmentation:!1,FragmentType:"",Counter:0,Data:Uint8Array},codingRes:[],origin:[]}},methods:{pressControlButton(){this.isDataButtonEnable=!1,this.isControlButtonEnable=!0,this.isAcknowledgeEnable=!1,this.clearPreviosDataResult(),this.testMode=0},pressDataButton(){this.isDataButtonEnable=!0,this.isControlButtonEnable=!1,this.isAcknowledgeEnable=!1,this.clearPreviosControlResult(),this.testMode=1},pressAcknowledgeButton(){this.isDataButtonEnable=!1,this.isControlButtonEnable=!1,this.isAcknowledgeEnable=!0,this.clearPreviosAckResult(),this.testMode=2},testRequest(e){let t=new Uint8Array(2);t[0]=this.testMode,t[1]=e,(0,Ne.tQ)("/test",t).then((e=>{200==e.status?e.arrayBuffer().then((e=>{let t=new Uint8Array(e);switch(t[0]){case 0:this.controlCommandDecode(t.subarray(1));break;case 1:this.dataMessageDecode(t.subarray(1));break;case 2:this.ackMessageDecode(t.subarray(1));break}this.origin=[...t]})):console.log("error:"+e.status)})).catch((e=>console.error(e)))},controlCommandDecode(e){let t=De.T.ControlCommand;this.clearPreviosControlResult(),t.clear(),t.decode(e),this.processingControlCommand(t),this.codingRes=[...t.code()]},clearPreviosControlResult(){this.CommandType="",this.UserId=0,this.UserName="",this.UserParameter.Name="",this.UserParameter.Surname="",this.UserParameter.Familyname="",this.UserParameter.Birthdate="",this.UsersLlist.Type="",this.UsersLlist.NumOfItems=0,this.UsersLlist.Items.splice(0,this.UsersLlist.Items.length),this.origin.length=0,this.codingRes.length=0},processingControlCommand(e){switch(this.CommandType=e.ControlMessageType.toString(),e.ControlMessageType.toString()){case"connect":case"disconnect":this.UserId=e.UserId;break;case"registration":case"user update":this.UserId=e.UserId,this.UserName=e.UserName;break;case"user info":this.UserParameter.Name=e.UserParameter.Name,this.UserParameter.Surname=e.UserParameter.Surname,this.UserParameter.Familyname=e.UserParameter.Familyname,this.UserParameter.Birthdate=e.UserParameter.Birthdate;break;case"users list":this.UsersLlist.Type=e.UsersList.ListCommand,this.UsersLlist.NumOfItems=e.UsersList.NumberOfItems;for(let t of e.UsersList.ListElements)this.UsersLlist.Items.push({id:t[0],name:t[1].name,online:t[1].online,nrm:t[1].nrm});break}},dataMessageDecode(e){let t=De.T.Data;this.clearPreviosDataResult(),t.clear(),t.decode(e),this.processingDataMessage(t),this.codingRes=[...t.code()]},clearPreviosDataResult(){this.DataType="",this.DataFormat="",this.IndexNumber=0,this.UserId=0,this.RoomId=NaN,this.Message.Fragmentation=!1,this.Message.FragmentType="",this.Message.Counter=0,this.Message.Text="",this.origin.length=0,this.codingRes.length=0},processingDataMessage(e){this.DataType=e.DataMessageType.toString(),this.DataFormat=e.DataMessageFormat.toString(),this.IndexNumber=e.IndexNumber,this.UserId=e.UserId,this.RoomId=e.RoomId,this.Message.Fragmentation=e.UserData.Fragmentation.On,this.Message.FragmentType=De.T.FragmentTypeToString(e.UserData.Fragmentation.FragmentType),this.Message.Counter=e.UserData.Fragmentation.Counter,this.Message.Data=e.UserData.Data},ackMessageDecode(e){let t=De.T.Acknowledge;this.clearPreviosAckResult(),t.clear(),t.decode(e),this.processingAckMessage(t),this.codingRes=[...t.code()]},clearPreviosAckResult(){this.AcknowledgeType="",this.IndexNumber=0,this.UserId=0,this.RoomId=NaN,this.origin.length=0,this.codingRes.length=0},processingAckMessage(e){this.AcknowledgeType=e.AckMessageType.toString(),this.IndexNumber=e.IndexNumber,this.UserId=e.UserId,this.RoomId=e.RoomId}}};const _e=(0,P.A)(fe,[["render",u]]);var Fe=_e}}]);
//# sourceMappingURL=467.d178fd4c.js.map