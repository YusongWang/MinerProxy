(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-53a89e5a"],{"2eab":function(t,e,a){"use strict";a("df51")},"8b11":function(t,e,a){},9406:function(t,e,a){"use strict";a.r(e);var s=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"dashboard-editor-container"},[a("panel-group",{attrs:{source:t.current,coin:t.coin},on:{onSelect:t.onSelect}}),a("status-group",{attrs:{source:t.current}}),a("el-row",{staticClass:"chart-box"},[a("div",{staticClass:"shadow"}),a("div",{staticClass:"title"},[a("span",[t._v(t._s(t.$t("home.poolData")))]),a("el-radio-group",{attrs:{size:"small"},model:{value:t.lineChart,callback:function(e){t.lineChart=e},expression:"lineChart"}},t._l(t.lineChartRadios,(function(e,s){return a("el-radio-button",{key:s,attrs:{label:e.label}},[t._v(" "+t._s(e.value)+" ")])})),1),a("el-select",{attrs:{size:"small"},model:{value:t.lineChartTime,callback:function(e){t.lineChartTime=e},expression:"lineChartTime"}},t._l(t.options,(function(t){return a("el-option",{key:t.value,attrs:{label:t.label,value:t.value}})})),1)],1),a("line-chart",{attrs:{"chart-data":t.lineChartData}})],1),a("el-row",{staticClass:"chart-box"},[a("div",{staticClass:"shadow"}),a("div",{staticClass:"title"},[a("span",[t._v(t._s(t.$t("home.machinePerformance")))]),a("el-radio-group",{attrs:{size:"small"},model:{value:t.performanceChart,callback:function(e){t.performanceChart=e},expression:"performanceChart"}},t._l(t.performanceChartRadios,(function(e,s){return a("el-radio-button",{key:s,attrs:{label:e.label}},[t._v(" "+t._s(e.value)+" ")])})),1),a("el-select",{attrs:{size:"small"},model:{value:t.performanceChartTime,callback:function(e){t.performanceChartTime=e},expression:"performanceChartTime"}},t._l(t.options,(function(t){return a("el-option",{key:t.value,attrs:{label:t.label,value:t.value}})})),1)],1),a("performance",{attrs:{"chart-data":t.performanceChartData,"class-name":"performance"}})],1)],1)},r=[],i=a("1da1"),n=a("5530"),o=(a("96cf"),a("d81d"),a("7db0"),a("d3b7"),function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"panel-group"},[a("el-row",{attrs:{gutter:40}},[a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:4}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.currency")))]),a("el-dropdown",{attrs:{trigger:"click"},on:{command:t.handleCommand}},[a("span",{staticClass:"value"},[t._v(" "+t._s(t.coin)+" "),a("i",{staticClass:"el-icon-arrow-down el-icon--right"})]),a("el-dropdown-menu",{attrs:{slot:"dropdown"},slot:"dropdown"},[a("el-dropdown-item",{attrs:{command:"ETH"}},[t._v("ETH")]),a("el-dropdown-item",{attrs:{command:"ETC"}},[t._v("ETC")])],1)],1)],1)]),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:3}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.poolNumber")))]),a("div",{staticClass:"value",staticStyle:{color:"#46bc59"}},[t._v(t._s(t.source&&t.source.pool_length||0)),a("span",{staticClass:"small"},[t._v(t._s(t.$t("home.piece")))])])])]),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:3}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.onlineMiner")))]),a("div",{staticClass:"value",staticStyle:{color:"#46bc59"}},[t._v(t._s(t.source&&+t.source.online_worker||0)),a("span",{staticClass:"small"},[t._v(t._s(t.$t("home.tower")))])])])]),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:3}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.offlineMiner")))]),a("div",{staticClass:"value",staticStyle:{color:"#fa6800"}},[t._v(t._s(t.source&&+t.source.offline_worker||0)),a("span",{staticClass:"small"},[t._v(t._s(t.$t("home.tower")))])])])]),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:4}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.totalComputingPower")))]),a("div",{staticClass:"value"},[t._v(t._s(t.source&&+t.source.total_hash?t.$utils.unitFilter(t.source.total_hash).value:0)),a("span",{staticClass:"small"},[t._v(t._s(t.source&&+t.source.total_hash?t.$utils.unitFilter(t.source.total_hash).unit:""))])])])]),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:6}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.runtime")))]),a("div",{staticClass:"value"},[t._v(t._s(t.source&&t.source.online_time?t.$utils.runtimeFilter(t.source.online_time):0))])])])],1)],1)}),l=[],c=a("5a0c"),u=a.n(c),h=a("d772");u.a.extend(h);var d={props:{source:{type:Object},coin:{type:String,default:"ETH"}},methods:{handleCommand:function(t){this.$emit("onSelect",t)}}},m=d,f=(a("2eab"),a("2877")),v=Object(f["a"])(m,o,l,!1,null,null,null),p=v.exports,_=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{class:t.className,style:{height:t.height,width:t.width}})},b=[],C=(a("b0c0"),a("313e")),g=a.n(C),$=a("ed08"),x={data:function(){return{$_sidebarElm:null,$_resizeHandler:null}},mounted:function(){var t=this;this.$_resizeHandler=Object($["debounce"])((function(){t.chart&&t.chart.resize()}),100),this.$_initResizeEvent(),this.$_initSidebarResizeEvent()},beforeDestroy:function(){this.$_destroyResizeEvent(),this.$_destroySidebarResizeEvent()},activated:function(){this.$_initResizeEvent(),this.$_initSidebarResizeEvent()},deactivated:function(){this.$_destroyResizeEvent(),this.$_destroySidebarResizeEvent()},methods:{$_initResizeEvent:function(){window.addEventListener("resize",this.$_resizeHandler)},$_destroyResizeEvent:function(){window.removeEventListener("resize",this.$_resizeHandler)},$_sidebarResizeHandler:function(t){"width"===t.propertyName&&this.$_resizeHandler()},$_initSidebarResizeEvent:function(){this.$_sidebarElm=document.getElementsByClassName("sidebar-container")[0],this.$_sidebarElm&&this.$_sidebarElm.addEventListener("transitionend",this.$_sidebarResizeHandler)},$_destroySidebarResizeEvent:function(){this.$_sidebarElm&&this.$_sidebarElm.removeEventListener("transitionend",this.$_sidebarResizeHandler)}}};a("817d"),a("a524");var w=function(t){return t.map((function(t){return{name:t.name,itemStyle:{normal:{color:t.color[0],lineStyle:{color:t.color[0],width:2},areaStyle:{color:t.color[1]}}},smooth:!1,type:"line",data:t.data,suffix:t.suffix,animationDuration:2800,animationEasing:"cubicInOut"}}))},y={mixins:[x],props:{className:{type:String,default:"chart"},width:{type:String,default:"100%"},height:{type:String,default:"350px"},autoResize:{type:Boolean,default:!0},chartData:{type:Object,required:!0}},data:function(){return{chart:null}},computed:{theme:function(){return this.$store.getters.theme}},watch:{chartData:{deep:!0,handler:function(t){this.chart.dispose(),this.chart=null,this.initChart({theme:"dark"===this.theme?"dark":"macarons"})}},theme:function(t){this.initChart({theme:"dark"===t?"dark":"macarons"})}},mounted:function(){var t=this;this.$nextTick((function(){t.initChart({theme:"dark"===t.theme?"dark":"macarons"})}))},beforeDestroy:function(){this.chart&&(this.chart.dispose(),this.chart=null)},methods:{initChart:function(t){var e=t.theme;this.chart=g.a.init(this.$el,e),this.setOptions(this.chartData)},setOptions:function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:{},e=t.series,a=t.xAxis;this.chart.setOption({backgroundColor:"",xAxis:{data:a,boundaryGap:!1,axisTick:{show:!1}},grid:{left:10,right:10,bottom:20,top:30,containLabel:!0},tooltip:{trigger:"axis",axisPointer:{type:"cross"},padding:[5,10],formatter:function(t){for(var a=t[0].name+"<br/>",s=0,r=t.length;s<r;s++)a+=t[s].marker+t[s].seriesName+"："+t[s].value+(e[t[s].seriesIndex].suffix||"")+"<br/>";return a}},yAxis:{axisTick:{show:!1}},series:w(e)})}}},k=y,E=Object(f["a"])(k,_,b,!1,null,null,null),S=E.exports,R=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"status-group-panel-group"},[a("el-row",{attrs:{gutter:40}},[a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:8}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.statusGroupOne")))]),a("div",{staticClass:"value"},[t._v(" "+t._s(t.source&&+t.source.total_shares||0)+" / "+t._s(t.source&&+t.source.total_diff?t.$utils.unitFilter(t.source.total_diff).value:0)+" "),a("span",{staticClass:"small"},[t._v(" "+t._s(t.source&&+t.source.total_diff?t.$utils.unitFilter(t.source.total_diff).unit:"")+" ")])])])]),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:8}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.statusGroupTwo")))]),a("div",{staticClass:"value"},[t._v(" "+t._s(t.source&&+t.source.fee_shares||0)+" / "+t._s(t.source&&+t.source.fee_diff?t.$utils.unitFilter(t.source.fee_diff).value:0)+" "),a("span",{staticClass:"small"},[t._v(" "+t._s(t.source&&+t.source.fee_diff?t.$utils.unitFilter(t.source.fee_diff).unit:"")+" ")]),t._v(" / "+t._s(t.source&&+t.source.fee_rate||0)+" "),a("span",{staticClass:"small"},[t._v("%")])])])]),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:8}},[a("div",{staticClass:"card-panel"},[a("div",{staticClass:"label"},[t._v(t._s(t.$t("home.statusGroupThree")))]),a("div",{staticClass:"value"},[t._v(" "+t._s(t.source&&+t.source.dev_shares||0)+" / "+t._s(t.source&&+t.source.dev_diff?t.$utils.unitFilter(t.source.dev_diff).value:0)+" "),a("span",{staticClass:"small"},[t._v(" "+t._s(t.source&&+t.source.dev_diff?t.$utils.unitFilter(t.source.dev_diff).unit:"")+" ")]),t._v(" / "+t._s(t.source&&+t.source.dev_rate||0)+" "),a("span",{staticClass:"small"},[t._v("%")])])])])],1)],1)},z=[],O={props:{source:{type:Object}},methods:{handleSetLineChartData:function(t){this.$emit("handleSetLineChartData",t)}}},D=O,T=(a("d873"),Object(f["a"])(D,R,z,!1,null,null,null)),j=T.exports,H=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{class:t.className,style:{height:t.height,width:t.width}})},M=[];a("817d"),a("a524");var A=function(t){return t.map((function(t){return{name:t.name,itemStyle:{normal:{color:t.color[0],lineStyle:{color:t.color[0],width:2},areaStyle:{color:t.color[1]}}},smooth:!1,type:"line",data:t.data,suffix:t.suffix,animationDuration:2800,animationEasing:"cubicInOut"}}))},P={mixins:[x],props:{className:{type:String,default:"chart"},width:{type:String,default:"100%"},height:{type:String,default:"350px"},autoResize:{type:Boolean,default:!0},chartData:{type:Object,required:!0}},data:function(){return{chart:null}},computed:{theme:function(){return this.$store.getters.theme}},watch:{chartData:{deep:!0,handler:function(t){this.chart.dispose(),this.chart=null,this.initChart({theme:"dark"===this.theme?"dark":"macarons"})}},theme:function(t){this.initChart({theme:"dark"===t?"dark":"macarons"})}},mounted:function(){var t=this;this.$nextTick((function(){t.initChart({theme:"dark"===t.theme?"dark":"macarons"})}))},beforeDestroy:function(){this.chart&&(this.chart.dispose(),this.chart=null)},methods:{initChart:function(t){var e=t.theme;this.chart=g.a.init(this.$el,e),this.setOptions(this.chartData)},setOptions:function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:{},e=t.series,a=t.xAxis;this.chart.setOption({backgroundColor:"",xAxis:{data:a,boundaryGap:!1,axisTick:{show:!1}},grid:{left:10,right:10,bottom:20,top:30,containLabel:!0},tooltip:{trigger:"axis",axisPointer:{type:"cross"},padding:[5,10],formatter:function(t){for(var a=t[0].name+"<br/>",s=0,r=t.length;s<r;s++)a+=t[s].marker+t[s].seriesName+"："+t[s].value+(e[t[s].seriesIndex].suffix||"")+"<br/>";return a}},yAxis:{axisTick:{show:!1}},series:A(e)})}}},L=P,F=Object(f["a"])(L,H,M,!1,null,null,null),G=F.exports,N=a("b775");function U(){return I.apply(this,arguments)}function I(){return I=Object(i["a"])(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,Object(N["a"])({url:"/api/dashborad",method:"get"});case 2:return t.abrupt("return",t.sent);case 3:case"end":return t.stop()}}),t)}))),I.apply(this,arguments)}var B=function(){for(var t=[],e=0;e<15;e++)t[e]=Object($["random"])(0,200);return t},q=function(t){var e=t.series;return e.map((function(t){return Object(n["a"])(Object(n["a"])({},t),{},{data:B()})}))},J=function(t){for(var e=[],a=u()(),s=0;s<15;s++)a=a.subtract(t,"minute"),e[s]=a.format("MM-DD HH:mm");return e},K={name:"DashboardAdmin",components:{Performance:G,PanelGroup:p,LineChart:S,StatusGroup:j},data:function(){return{dashborad:null,coin:"ETH",lineChart:"a",performanceChart:"a",lineChartTime:5,performanceChartTime:5,lineChartRadios:[{label:"a",value:this.$t("home.onlineMiner"),series:[{name:this.$t("home.onlineMiner"),color:["rgba(149,162,255)","rgba(149,162,255,0.7)"]},{name:this.$t("home.averageOnlineMiner"),color:["rgb(250,128,128)","rgba(250,128,128,0.7)"]}]},{label:"b",value:this.$t("home.machineComputingPower"),series:[{name:this.$t("home.machineComputingPower"),color:["rgba(255,192,118)","rgba(255,192,118,0.7)"],suffix:"GH/S"},{name:this.$t("home.averageMachineComputingPower"),color:["rgba(250,231,104)","rgba(250,231,104,0.7)"],suffix:"GH/S"}]}],performanceChartRadios:[{label:"a",value:"CPU",series:[{name:this.$t("home.CPUUsage"),color:["rgba(135,232,133)","rgba(135,232,133,0.7)"],suffix:"%"}]},{label:"b",value:this.$t("home.RAM"),series:[{name:this.$t("home.systemMemoryUsage"),color:["rgba(60,185,252)","rgba(60,185,252,0.7)"],suffix:"M"},{name:this.$t("home.processMemoryUsage"),color:["rgba(115,171,245)","rgba(115,171,245,0.7)"],suffix:"M"}]}],options:[{label:this.$t("home.minute",{num:5}),value:5},{label:this.$t("home.minute",{num:10}),value:10},{label:this.$t("home.minute",{num:30}),value:30},{label:this.$t("home.hour",{num:1}),value:60},{label:this.$t("home.hour",{num:3}),value:180},{label:this.$t("home.hour",{num:6}),value:360},{label:this.$t("home.hour",{num:12}),value:720},{label:this.$t("home.day",{num:1}),value:1440}]}},created:function(){var t=this;return Object(i["a"])(regeneratorRuntime.mark((function e(){return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,t.getData();case 2:case"end":return e.stop()}}),e)})))()},computed:{current:function(){return this.dashborad?this.dashborad[this.coin]:null},lineChartData:function(){var t=this;return console.log({series:q(this.lineChartRadios.find((function(e){return e.label===t.lineChart}))),xAxis:J(this.lineChartTime)}),{series:q(this.lineChartRadios.find((function(e){return e.label===t.lineChart}))),xAxis:J(this.lineChartTime)}},performanceChartData:function(){var t=this;return{series:q(this.performanceChartRadios.find((function(e){return e.label===t.performanceChart}))),xAxis:J(this.performanceChartTime)}}},methods:{handleSetLineChartData:function(t){this.lineChartData=q[t]},getData:function(){var t=this;return Object(i["a"])(regeneratorRuntime.mark((function e(){var a,s;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,U();case 2:a=e.sent,s=a.data,s&&(t.dashborad=s);case 5:case"end":return e.stop()}}),e)})))()},onSelect:function(t){this.coin=t}}},Q=K,V=(a("af54"),Object(f["a"])(Q,s,r,!1,null,null,null));e["default"]=V.exports},af54:function(t,e,a){"use strict";a("8b11")},d873:function(t,e,a){"use strict";a("ffa2")},df51:function(t,e,a){},ffa2:function(t,e,a){}}]);