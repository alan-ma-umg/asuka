package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chenset/asuka/database"
	"golang.org/x/net/html"
	"hash/crc32"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

var html200 = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html lang="zh-CN">
<head>
<meta name="Content-Type" content="text/html;charset=utf-8" />
<meta name="Referrer" content="unsafe-url" />
<meta content="True" name="HandheldFriendly" />
<meta name="theme-color" content="#333344" />
<meta name="apple-mobile-web-app-capable" content="yes" />
<meta name="mobile-web-app-capable" content="yes" />
<meta name="detectify-verification" content="d0264f228155c7a1f72c3d91c17ce8fb" />
<meta name="p:domain_verify" content="b87e3b55b409494aab88c1610b05a5f0" />
<meta name="alexaVerifyID" content="OFc8dmwZo7ttU4UCnDh1rKDtLlY" />
<meta name="baidu-site-verification" content="D00WizvYyr" />
<meta name="msvalidate.01" content="D9B08FEA08E3DA402BF07ABAB61D77DE" />
<meta property="wb:webmaster" content="f2f4cb229bda06a4" />
<meta name="google-site-verification" content="LM_cJR94XJIqcYJeOCscGVMWdaRUvmyz6cVOqkFplaU" />
<title>V2EX</title>
<link rel="dns-prefetch" href="//static.v2ex.com" />
<link rel="dns-prefetch" href="//cdn.v2ex.com" />
<link rel="dns-prefetch" href="//i.v2ex.co" />
<link rel="stylesheet" type="text/css" media="screen" href="/css/basic.css?v=102905:1566175659:3.9.8.3" />
<link rel="stylesheet" type="text/css" media="screen" href="/static/css/style.css?v=a2cc48571a034df380f58d1cbbdc2a0d" />
<link rel="stylesheet" type="text/css" media="screen" href="/css/desktop.css?v=3.9.8.3" />
<link rel="stylesheet" type="text/css" media="screen" href="/static/css/night.css?v=e3b286cc5b8dbc5bc2c7823315a19add" />
<link rel="stylesheet" href="/static/css/tomorrow-night.css?v=b0e517e91324470337a008bd08305fec" type="text/css" />
<link rel="icon" sizes="192x192" href="/static/img/v2ex_192.png" />
<link rel="shortcut icon" href="/static/img/icon_rayps_64.png" type="image/png" />
<link rel="stylesheet" type="text/css" href="/static/css/font-awesome.min.css?v=295235b28b6e649d99539a9d32b95d30" />
<script src="/static/dist/combo.js?v=050545155819e5b98d0f997c5cd17987" type="text/javascript" defer></script>
<link href="/static/css/jquery.textcomplete.css?v=5a041d39010ded8724744170cea6ce8d" rel="stylesheet" />
<link href="/static/css/select2.min.css?v=af22a7e2bfec4d2a82c4dde613a52fb1" rel="stylesheet" />
<link href="/static/js/selectboxit/selectboxit.css?v=5dc55d3860ef80ef1875d6800a5fbfa3" rel="stylesheet">
<link rel="stylesheet" type="text/css" href="/static/css/loading-bar.css?v=4a5569962dc210bac72543b39b347a0e" />
<meta name="description" content="" />
<link rel="canonical" href="https://www.v2ex.com/" />
<script>
    document.addEventListener("DOMContentLoaded", function(event) { 
		protectTraffic();
        $( "#MyNodes" ).sortable();
        $( "#MyNodes" ).disableSelection();
        $( "#MyNodes" ).sortable({
            stop: function( event, ui ) {
                var sorted = $( "#MyNodes" ).sortable( "serialize", { key: "n" } );
                $.post('/my/nodes/sorted', { sorted : sorted }, function(data) {

                });
            }
        });

        
        blocked = [215886];
        ignored_topics = [];
        $("#TopicsHot").children('.cell').each( function(index) {
            for (i in blocked) {
                if ($(this).hasClass('from_' + blocked[i])) {
                    $(this).css('display', 'none');
                }
            }
            for (i in ignored_topics) {
                css_class = 'hot_t_' + ignored_topics[i];
                if ($(this).hasClass(css_class)) {
                    $(this).css('display', 'none');
                }
            }
        });
        
    });
</script>
</head>
<body>
<div id="Top">
<div class="content">
<div style="padding-top: 6px;">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="110" align="left"><a href="/" name="top" title="way to explore"><div id="Logo"></div></a></td>
<td width="auto" align="left">
<div id="Search"><form action="https://www.google.com" onsubmit="return dispatch()" target="_blank"><div id="qbar"><input type="text" maxlength="40" name="q" id="q" value="" onfocus="$('#qbar').addClass('qbar_focus')" onblur="$('#qbar').removeClass('qbar_focus')" /></div></form></div>
</td>
<td width="570" align="right" style="padding-top: 2px; font-size: 13px;"><a href="/" class="top">Home</a>&nbsp;&nbsp;&nbsp;<a href="/member/chenset" class="top">chenset</a>&nbsp;&nbsp;&nbsp;<a href="/notes" class="top">Notes</a>&nbsp;&nbsp;&nbsp;<a href="/t" class="top">Timeline</a>&nbsp;&nbsp;&nbsp;<a href="/settings" class="top">Settings</a>&nbsp;&nbsp;&nbsp;<a href="#;" onclick="if (confirm('确定要从 V2EX 登出？')) { location.href= '/signout?once=42165'; }" class="top">Sign Out</a></td>
</tr>
</table>
</div>
</div>
</div>
<div id="Wrapper" class="Night">
<div class="content">
<div id="Leftbar"></div>
<div id="Rightbar">
<div class="sep20"></div>
<div class="box">
<div class="cell">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top"><a href="/member/chenset"><img src="//cdn.v2ex.com/avatar/36ca/f1aa/102905_large.png?m=1566175659" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10" valign="top"></td>
<td width="auto" align="left"><div class="fr"><a href="/settings/night/toggle?once=42165" class="light-toggle"><img src="/static/img/toggle-dark.png" align="absmiddle" height="10" alt="Dark" /></a></div><span class="bigger"><a href="/member/chenset">chenset</a></span>
</td>
</tr>
</table>
<div class="sep10"></div>
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="33%" align="center"><a href="/my/nodes" class="dark" style="display: block;"><span class="bigger">8</span><div class="sep3"></div><span class="fade">Nodes</span></a></td>
<td width="34%" style="border-left: 1px solid rgba(100, 100, 100, 0.4); border-right: 1px solid rgba(100, 100, 100, 0.4);" align="center"><a href="/my/topics" class="dark" style="display: block;"><span class="bigger">449</span><div class="sep3"></div><span class="fade">Topics</span></a></td>
<td width="33%" align="center"><a href="/my/following" class="dark" style="display: block;"><span class="bigger">23</span><div class="sep3"></div><span class="fade">Following</span></a></td>
</tr>
</table>
</div>
<div class="cell" id="member-activity">
<div class="member-activity-bar">
<div class="member-activity-start" style="width: 120px;"></div>
</div>
</div>
<div class="cell" style="padding: 5px;">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="28"><a href="/new"><img src="/static/img/essentials/compose.png?v=b9e1f045f4ad639733bf9f6dbc62ed4c" width="28" border="0" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle" align="left"><a href="/new">Compose</a></td>
</tr>
</table>
</div>
<div class="inner"><div class="fr" id="money" style="margin: -3px 0px 0px 0px"><a href="/balance" class="balance_area" style="">4 <img src="/static/img/gold@2x.png" height="16" alt="G" border="0" /> 59 <img src="/static/img/silver@2x.png" height="16" alt="S" border="0" /> 13 <img src="/static/img/bronze@2x.png" height="16" alt="B" border="0" /></a></div><a href="/notifications" class="fade">0 unread</a></div>
</div>
<div class="sep20"></div>
<div class="box">
<div class="inner">
<script async src="//pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"></script>

<ins class="adsbygoogle" style="display:inline-block;width:250px;height:250px" data-ad-client="ca-pub-3465543440750523" data-ad-slot="9619519096"></ins>
<script>
        (adsbygoogle = window.adsbygoogle || []).push({});
        </script>
</div>
</div>
<div class="sep20"></div>
<div class="box">
<div class="inner" style="padding: 5px;">
<div class="gray f12" style="padding: 5px;">My Favorite Nodes</div>
<div id="MyNodes">
<div class="node" id="n_201587"><div class="node_compose"><a href="/new/bb"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/bb"><img src="//cdn.v2ex.com/navatar/a3c6/5c29/108_normal.png?m=1571088660" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/bb">宽带症候群</a></div>
<div class="node" id="n_182340"><div class="node_compose"><a href="/new/ss"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/ss"><img src="/static/img/node_normal.png" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/ss">沉默的螺旋</a></div>
<div class="node" id="n_196031"><div class="node_compose"><a href="/new/jp"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/jp"><img src="//cdn.v2ex.com/navatar/fb7b/9ffa/357_normal.png?m=1360450842" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/jp">日本</a></div>
<div class="node" id="n_138515"><div class="node_compose"><a href="/new/bitcoin"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/bitcoin"><img src="//cdn.v2ex.com/navatar/bbf9/4b34/403_normal.png?m=1543441132" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/bitcoin">Bitcoin</a></div>
<div class="node" id="n_122112"><div class="node_compose"><a href="/new/flamewar"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/flamewar"><img src="//cdn.v2ex.com/navatar/7588/7499/314_normal.png?m=1563795094" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/flamewar">水深火热</a></div>
<div class="node" id="n_122113"><div class="node_compose"><a href="/new/chamber"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/chamber"><img src="//cdn.v2ex.com/navatar/bac9/162b/364_normal.png?m=1523199919" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/chamber">Chamber</a></div>
<div class="node" id="n_83790"><div class="node_compose"><a href="/new/pi"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/pi"><img src="//cdn.v2ex.com/navatar/884d/247c/650_normal.png?m=1484889842" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/pi">Raspberry Pi</a></div>
<div class="node" id="n_103985"><div class="node_compose"><a href="/new/deals"><img src="/static/img/compose.png" align="absmiddle" border="0" width="23" height="18" alt="New Topic" /></a></div><a href="/go/deals"><img src="//cdn.v2ex.com/navatar/8d31/7bdc/747_normal.png?m=1523201604" border="0" align="absmiddle" width="24" /></a>&nbsp; <a href="/go/deals">优惠信息</a></div>
</div>
</div>
</div>
<div class="sep20"></div>
<div class="box" id="TopicsHot">
<div class="cell"><span class="fade">Today Top 10</span></div>
<div class="cell from_222151 hot_t_613097">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/infra"><img src="//cdn.v2ex.com/avatar/2882/fb30/222151_normal.png?m=1491987872" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613097">这么多人关注扣扣 for Linux ，有多少银用 Linux 做桌面？</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_355191 hot_t_613021">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/fiypig"><img src="//cdn.v2ex.com/avatar/9578/84e6/355191_normal.png?m=1571627741" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613021">大佬们有辞职了以后还回原公司的吗</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_216646 hot_t_613041">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/Levi233"><img src="//cdn.v2ex.com/gravatar/da064924605953132a0baad75bbe1384?s=24&d=retro" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613041">xxx 强调，把区块链作为核心技术自主创新重要突破口</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_284747 hot_t_613069">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/xiaoming1992"><img src="//cdn.v2ex.com/avatar/c7e9/042b/284747_normal.png?m=1551399319" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613069">装机求助帖，在线等，挺急的😂</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_146356 hot_t_613125">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/yizhimamong"><img src="//cdn.v2ex.com/avatar/bd0a/d5b0/146356_normal.png?m=1451320379" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613125">试用期六个月和三个月的区别</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_315449 hot_t_613045">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/CSGO"><img src="//cdn.v2ex.com/gravatar/6a7906a28585a0b0ff656e13e1f04fde?s=24&d=retro" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613045">没人关注“一刻相册”APP 吗？</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_364531 hot_t_613124">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/xhqpp"><img src="//cdn.v2ex.com/gravatar/f26c95265b257b9aff0a07265daa2fef?s=24&d=retro" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613124">如果油管和奈飞接国内的 CDN，是不是可以极大缓解国际线路拥堵？</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_254353 hot_t_613159">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/FaiChou"><img src="//cdn.v2ex.com/gravatar/0a0b5ee72ddaab27625c5af407931a44?s=24&d=retro" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613159">帮忙起一个贸易公司名字</a>
</span>
</td>
</tr>
</table>
</div>
<div class="cell from_426816 hot_t_613184">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="24" valign="middle" align="center">
<a href="/member/JmingZhang"><img src="//cdn.v2ex.com/avatar/2278/5c86/426816_normal.png?m=1565918907" class="avatar" border="0" align="default" style="max-width: 24px; max-height: 24px;" /></a>
</td>
<td width="10"></td>
<td width="auto" valign="middle">
<span class="item_hot_topic_title">
<a href="/t/613184">如果你亲人犯法了你会举报他还是会等他自己被抓？</a>
</span>
</td>
</tr>
</table>
</div>
</div>
<div class="sep20"></div>
<div class="box">
<div class="cell"><div class="fr"></div><span class="fade">Hottest Nodes</span></div>
<div class="cell">
<a href="/go/qna" class="item_node">问与答</a><a href="/go/jobs" class="item_node">酷工作</a><a href="/go/programmer" class="item_node">程序员</a><a href="/go/share" class="item_node">分享发现</a><a href="/go/macos" class="item_node">macOS</a><a href="/go/create" class="item_node">分享创造</a><a href="/go/python" class="item_node">Python</a><a href="/go/apple" class="item_node">Apple</a><a href="/go/career" class="item_node">职场话题</a><a href="/go/android" class="item_node">Android</a><a href="/go/iphone" class="item_node">iPhone</a><a href="/go/bb" class="item_node">宽带症候群</a><a href="/go/gts" class="item_node">全球工单系统</a><a href="/go/cv" class="item_node">求职</a><a href="/go/mbp" class="item_node">MacBook Pro</a>
</div>
<div class="inner"><a href="/index.xml" target="_blank"><img src="/static/img/rss.png" align="absmiddle" border="0" style="margin-top:-3px;" /></a>&nbsp; <a href="/index.xml" target="_blank">RSS</a></div>
</div>
<div class="sep20"></div>
<div class="box">
<div class="cell"><div class="fr"></div><span class="fade">最近新增节点</span></div>
<div class="inner">
<a href="/go/ws" class="item_node">WebSocket</a><a href="/go/wg" class="item_node">WireGuard</a><a href="/go/zsh" class="item_node">Z shell</a><a href="/go/applearcade" class="item_node">Apple Arcade</a><a href="/go/darkmode" class="item_node">夜间模式</a><a href="/go/quake" class="item_node">雷神之锤系列</a><a href="/go/rss" class="item_node">RSS</a><a href="/go/jsonfeed" class="item_node">JSON Feed</a><a href="/go/vtuber" class="item_node">Virtual YouTubers</a><a href="/go/terraform" class="item_node">Terraform</a><a href="/go/remote" class="item_node">远程工作</a><a href="/go/weekly" class="item_node">写周报</a><a href="/go/cloudflare" class="item_node">Cloudflare</a><a href="/go/libra" class="item_node">Libra</a><a href="/go/typescript" class="item_node">TypeScript</a><a href="/go/tex" class="item_node">TeX</a><a href="/go/stadia" class="item_node">Stadia</a><a href="/go/apex" class="item_node">Apex Legends</a><a href="/go/bujo" class="item_node">子弹笔记</a><a href="/go/2019" class="item_node">2019</a>
</div>
</div>
<div class="sep20"></div>
<div class="box">
<div class="cell"><span class="fade">Community Stats</span></div>
<div class="cell">
<table cellpadding="5" cellspacing="0" border="0" width="100%">
<tr>
<td width="60" align="right"><span class="gray">注册会员</span></td>
<td width="auto" align="left"><strong>449321</strong></td>
</tr>
<tr>
<td width="60" align="right"><span class="gray">主题</span></td>
<td width="auto" align="left"><strong>613252</strong></td>
</tr>
<tr>
<td width="60" align="right"><span class="gray">回复</span></td>
<td width="auto" align="left"><strong>8089197</strong></td>
</tr>
</table>
</div>
<div class="inner">
<span class="chevron">›</span> <a href="/top/rich">财富排行榜</a>
<div class="sep5"></div>
<span class="chevron">›</span> <a href="/top/player">消费排行榜</a>
</div>
</div>
<div class="sep20"></div>
</div>
<div id="Main">
<div class="sep20"></div>
<div class="box">
<div class="cell" id="Tabs">
<a href="/?tab=tech" class="tab">技术</a><a href="/?tab=creative" class="tab">创意</a><a href="/?tab=play" class="tab">好玩</a><a href="/?tab=apple" class="tab">Apple</a><a href="/?tab=jobs" class="tab">酷工作</a><a href="/?tab=deals" class="tab">交易</a><a href="/?tab=city" class="tab">城市</a><a href="/?tab=qna" class="tab">问与答</a><a href="/?tab=hot" class="tab_current">最热</a><a href="/?tab=all" class="tab">全部</a><a href="/?tab=r2" class="tab">R2</a><a href="/?tab=nodes" class="tab">节点</a><a href="/?tab=members" class="tab">关注</a>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/yizhimamong"><img src="//cdn.v2ex.com/avatar/bd0a/d5b0/146356_large.png?m=1451320379" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613125#reply45" class="topic-link">试用期六个月和三个月的区别</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/career">职场话题</a> &nbsp;•&nbsp; <strong><a href="/member/yizhimamong">yizhimamong</a></strong> &nbsp;•&nbsp; 54 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/PHPer233">PHPer233</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613125#reply45" class="count_livid">45</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/xhqpp"><img src="//cdn.v2ex.com/gravatar/f26c95265b257b9aff0a07265daa2fef?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613124#reply43" class="topic-link">如果油管和奈飞接国内的 CDN，是不是可以极大缓解国际线路拥堵？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/bb">宽带症候群</a> &nbsp;•&nbsp; <strong><a href="/member/xhqpp">xhqpp</a></strong> &nbsp;•&nbsp; 31 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/love">love</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613124#reply43" class="count_livid">43</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/FaiChou"><img src="//cdn.v2ex.com/gravatar/0a0b5ee72ddaab27625c5af407931a44?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613159#reply40" class="topic-link">帮忙起一个贸易公司名字</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/programmer">程序员</a> &nbsp;•&nbsp; <strong><a href="/member/FaiChou">FaiChou</a></strong> &nbsp;•&nbsp; 39 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/berumotto">berumotto</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613159#reply40" class="count_livid">40</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/JmingZhang"><img src="//cdn.v2ex.com/avatar/2278/5c86/426816_large.png?m=1565918907" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613184#reply40" class="topic-link">如果你亲人犯法了你会举报他还是会等他自己被抓？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/JmingZhang">JmingZhang</a></strong> &nbsp;•&nbsp; 6 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/coderluan">coderluan</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613184#reply40" class="count_livid">40</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/noli"><img src="//cdn.v2ex.com/avatar/9545/99b8/75369_large.png?m=1550400134" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613192#reply36" class="topic-link">[zz 不正确] 自从印度人进入美国 IT 行业之后</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/programmer">程序员</a> &nbsp;•&nbsp; <strong><a href="/member/noli">noli</a></strong> &nbsp;•&nbsp; 24 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/darmau">darmau</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613192#reply36" class="count_livid">36</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/zhihupron"><img src="//cdn.v2ex.com/avatar/35cf/c4ac/410998_large.png?m=1557903286" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613126#reply32" class="topic-link">100m 宽带 60 元一个月算不算贵啊？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/bb">宽带症候群</a> &nbsp;•&nbsp; <strong><a href="/member/zhihupron">zhihupron</a></strong> &nbsp;•&nbsp; 2 h 0 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/devlnt">devlnt</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613126#reply32" class="count_livid">32</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/masonvip"><img src="//cdn.v2ex.com/avatar/669a/f7cd/434139_large.png?m=1568267513" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613196#reply29" class="topic-link">这次我站在川普一边</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/iphone">iPhone</a> &nbsp;•&nbsp; <strong><a href="/member/masonvip">masonvip</a></strong> &nbsp;•&nbsp; 1 h 20 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/archey">archey</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613196#reply29" class="count_livid">29</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/w4mxl"><img src="//cdn.v2ex.com/gravatar/3896b6baf91ec1933c38f370964647b7?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613127#reply28" class="topic-link">V2LF - 使用 Flutter 开发的开源的 V2EX 客户端，更新支持了 iOS（iPadOS）13 / Android 10 Dark Mode</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/create">分享创造</a> &nbsp;•&nbsp; <strong><a href="/member/w4mxl">w4mxl</a></strong> &nbsp;•&nbsp; 37 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/kooritea">kooritea</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613127#reply28" class="count_livid">28</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/skymxc"><img src="//cdn.v2ex.com/gravatar/e30a3d57b6d6ec705d4b1375fcf28150?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613173#reply26" class="topic-link">想换个 iPhone 11，但我是个 Android 开发呀！</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/programmer">程序员</a> &nbsp;•&nbsp; <strong><a href="/member/skymxc">skymxc</a></strong> &nbsp;•&nbsp; 1 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/agagega">agagega</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613173#reply26" class="count_livid">26</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/sherlockwhite"><img src="//cdn.v2ex.com/gravatar/9dbcdb46934bff356a958436ba9356f0?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613204#reply25" class="topic-link">谷歌要回归中国？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/share">分享发现</a> &nbsp;•&nbsp; <strong><a href="/member/sherlockwhite">sherlockwhite</a></strong> &nbsp;•&nbsp; 2 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/ClarkAbe">ClarkAbe</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613204#reply25" class="count_livid">25</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/Qmanman"><img src="//cdn.v2ex.com/gravatar/718c3d228b8e20c4dc4a868430fb2c30?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613134#reply23" class="topic-link">武汉本科电子信息科学与技术，秋招跑了一个多月没找到工作</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/jobs">酷工作</a> &nbsp;•&nbsp; <strong><a href="/member/Qmanman">Qmanman</a></strong> &nbsp;•&nbsp; 44 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/mpb">mpb</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613134#reply23" class="count_livid">23</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/caopi"><img src="//cdn.v2ex.com/avatar/6e49/1de7/348060_large.png?m=1540892093" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613179#reply20" class="topic-link">如何保留加班的证据和维护自身权益呢</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/career">职场话题</a> &nbsp;•&nbsp; <strong><a href="/member/caopi">caopi</a></strong> &nbsp;•&nbsp; 13 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/Thx4">Thx4</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613179#reply20" class="count_livid">20</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/snowspace"><img src="//cdn.v2ex.com/avatar/0f86/c3cb/103827_large.png?m=1434036537" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613246#reply20" class="topic-link">京东 App 点击个人头像可以查看消费情况，各位都消费多少？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/jd">京东</a> &nbsp;•&nbsp; <strong><a href="/member/snowspace">snowspace</a></strong> &nbsp;•&nbsp; Just Now &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/nathanw">nathanw</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613246#reply20" class="count_livid">20</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/penghh"><img src="//cdn.v2ex.com/gravatar/81759fdb23c7ff1fc630c5ab4f78a81f?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613172#reply19" class="topic-link">你们是在哪购买正版软件的？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/penghh">penghh</a></strong> &nbsp;•&nbsp; 30 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/JamesR">JamesR</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613172#reply19" class="count_livid">19</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/monsterlin"><img src="//cdn.v2ex.com/avatar/13dc/de5d/216884_large.png?m=1506076061" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613199#reply18" class="topic-link">今早看到产品提到一个打开 App 结束其他 App 的需求....</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/android">Android</a> &nbsp;•&nbsp; <strong><a href="/member/monsterlin">monsterlin</a></strong> &nbsp;•&nbsp; 13 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/HangoX">HangoX</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613199#reply18" class="count_livid">18</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/bagel"><img src="//cdn.v2ex.com/avatar/681b/305b/317118_large.png?m=1534689255" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613147#reply17" class="topic-link">有人和我一样担心不小心吃到食品包装的风险吗？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/bagel">bagel</a></strong> &nbsp;•&nbsp; 1 h 39 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/nieyujiang">nieyujiang</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613147#reply17" class="count_livid">17</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/kuyuzhiqi"><img src="//cdn.v2ex.com/gravatar/40d34c01993f1eb5990e11e369ff8475?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613153#reply15" class="topic-link">现在下载好难啊</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/kuyuzhiqi">kuyuzhiqi</a></strong> &nbsp;•&nbsp; 2 h 10 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/Johnny168">Johnny168</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613153#reply15" class="count_livid">15</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/Kcelone"><img src="//cdn.v2ex.com/avatar/635c/35dd/327132_large.png?m=1545212397" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613144#reply14" class="topic-link">git 技能复习进阶（开局一个键盘，内容全靠抄）</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/git">git</a> &nbsp;•&nbsp; <strong><a href="/member/Kcelone">Kcelone</a></strong> &nbsp;•&nbsp; 56 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/wysnylc">wysnylc</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613144#reply14" class="count_livid">14</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/lysS"><img src="//cdn.v2ex.com/avatar/57bc/f427/384142_large.png?m=1563786510" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613186#reply14" class="topic-link">关于化腾更新 Linux 版 QQ 的引申思考？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/lysS">lysS</a></strong> &nbsp;•&nbsp; 18 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/westoy">westoy</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613186#reply14" class="count_livid">14</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/jinYangChen"><img src="//cdn.v2ex.com/avatar/4359/8eaa/339946_large.png?m=1572049255" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613128#reply13" class="topic-link">套路云还有上车的没</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/programmer">程序员</a> &nbsp;•&nbsp; <strong><a href="/member/jinYangChen">jinYangChen</a></strong> &nbsp;•&nbsp; 28 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/xdaoo">xdaoo</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613128#reply13" class="count_livid">13</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/rlfz"><img src="//cdn.v2ex.com/gravatar/d50e2fa2975421b70477c5031966dc11?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613180#reply13" class="topic-link">大家一般如何清洗 apple watch</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/watch"> WATCH</a> &nbsp;•&nbsp; <strong><a href="/member/rlfz">rlfz</a></strong> &nbsp;•&nbsp; 1 h 28 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/wangchonglie">wangchonglie</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613180#reply13" class="count_livid">13</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/pppguest3962"><img src="//cdn.v2ex.com/avatar/1437/cf7c/214909_large.png?m=1487507092" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613191#reply13" class="topic-link">想找个比树莓派便宜，类似树莓派的&quot;平台&quot;...</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/pi">Raspberry Pi</a> &nbsp;•&nbsp; <strong><a href="/member/pppguest3962">pppguest3962</a></strong> &nbsp;•&nbsp; 3 h 10 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/IDCFAN">IDCFAN</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613191#reply13" class="count_livid">13</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/eason1874"><img src="//cdn.v2ex.com/avatar/3d27/4899/442190_large.png?m=1570079006" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613132#reply12" class="topic-link">这样设计网站日志数据库 MySQL 表格，有问题吗？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/programmer">程序员</a> &nbsp;•&nbsp; <strong><a href="/member/eason1874">eason1874</a></strong> &nbsp;•&nbsp; 5 h 37 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/eason1874">eason1874</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613132#reply12" class="count_livid">12</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/guisheng"><img src="//cdn.v2ex.com/avatar/373d/0b43/370999_large.png?m=1559010281" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613130#reply11" class="topic-link">macOS Catalina 抹盘重装后 竟然不会安装软件了？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/macos">macOS</a> &nbsp;•&nbsp; <strong><a href="/member/guisheng">guisheng</a></strong> &nbsp;•&nbsp; 3 h 17 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/guozozo">guozozo</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613130#reply11" class="count_livid">11</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/Windy00"><img src="//cdn.v2ex.com/avatar/5102/3f34/197520_large.png?m=1571038466" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613141#reply10" class="topic-link">git.io 被江苏电信指向 127.0.0.1</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/dns">DNS</a> &nbsp;•&nbsp; <strong><a href="/member/Windy00">Windy00</a></strong> &nbsp;•&nbsp; 1 h 2 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/janus77">janus77</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613141#reply10" class="count_livid">10</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/sagaxu"><img src="//cdn.v2ex.com/avatar/a4b7/3b82/200123_large.png?m=1481374171" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613148#reply10" class="topic-link">给电信 SDN 网关点个赞</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/bb">宽带症候群</a> &nbsp;•&nbsp; <strong><a href="/member/sagaxu">sagaxu</a></strong> &nbsp;•&nbsp; 2 h 14 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/cwbsw">cwbsw</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613148#reply10" class="count_livid">10</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/chickenJuicer"><img src="//cdn.v2ex.com/gravatar/7cd9a18b16733959a0e8f5c5a92ca235?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613162#reply10" class="topic-link">2020 福克斯 ST-Line 还是 2020 马自达 3 质雅???</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/car">汽车</a> &nbsp;•&nbsp; <strong><a href="/member/chickenJuicer">chickenJuicer</a></strong> &nbsp;•&nbsp; 2 h 5 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/duanran">duanran</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613162#reply10" class="count_livid">10</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/JaviDN"><img src="//cdn.v2ex.com/avatar/5a95/57da/417762_large.png?m=1559646063" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613163#reply10" class="topic-link">Python Flask 的问题</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/JaviDN">JaviDN</a></strong> &nbsp;•&nbsp; 3 h 36 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/SjwNo1">SjwNo1</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613163#reply10" class="count_livid">10</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/mdeadbird"><img src="//cdn.v2ex.com/avatar/3a44/a5f2/348218_large.png?m=1536626905" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613241#reply10" class="topic-link">马上 2020 年了，还是没有支持 4k60hz 的 typec 扩展坞吗？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/mbp">MacBook Pro</a> &nbsp;•&nbsp; <strong><a href="/member/mdeadbird">mdeadbird</a></strong> &nbsp;•&nbsp; 36 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/MrKerr">MrKerr</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613241#reply10" class="count_livid">10</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/wgh666"><img src="//cdn.v2ex.com/gravatar/dcdf09df345145522eced7d563351e65?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613169#reply9" class="topic-link">xs 256g 居然只要 6499 了</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/iphone">iPhone</a> &nbsp;•&nbsp; <strong><a href="/member/wgh666">wgh666</a></strong> &nbsp;•&nbsp; 4 h 10 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/CasualYours">CasualYours</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613169#reply9" class="count_livid">9</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/brblm"><img src="//cdn.v2ex.com/gravatar/54b0026202df5704d9895576ba48c79e?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613171#reply9" class="topic-link">被 qiang 过一次的 IP 是不是以后敏感期都会被 qiang 啊</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/brblm">brblm</a></strong> &nbsp;•&nbsp; 36 minutes ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/wangsongyan">wangsongyan</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613171#reply9" class="count_livid">9</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/KoukiGo"><img src="//cdn.v2ex.com/gravatar/0fae58a78ba246994540aaa20ab6b223?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613205#reply9" class="topic-link">有想去日本东京做 it 的吗 需要 Java 两年以上经验 不会日语可以</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/KoukiGo">KoukiGo</a></strong> &nbsp;•&nbsp; 1 h 37 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/gggxxxx">gggxxxx</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613205#reply9" class="count_livid">9</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/Awes0me"><img src="//cdn.v2ex.com/avatar/9ad2/67ac/111002_large.png?m=1484016881" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613146#reply8" class="topic-link">Catalina 的 AppStore 检测不到更新</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/macos">macOS</a> &nbsp;•&nbsp; <strong><a href="/member/Awes0me">Awes0me</a></strong> &nbsp;•&nbsp; 3 h 56 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/zhouweiluan">zhouweiluan</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613146#reply8" class="count_livid">8</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/qazwsxkevin"><img src="//cdn.v2ex.com/gravatar/bfefb99d6203d351791672a1d3fc936a?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613178#reply8" class="topic-link">请教 try 抛出异常后，有没有轻松回收（这个线程）内存的实现方法？</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/python">Python</a> &nbsp;•&nbsp; <strong><a href="/member/qazwsxkevin">qazwsxkevin</a></strong> &nbsp;•&nbsp; 1 h 12 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/ClericPy">ClericPy</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613178#reply8" class="count_livid">8</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/tantaolove"><img src="//cdn.v2ex.com/gravatar/34a4b32c6d447d98615731372faec6c8?s=48&d=retro" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613133#reply7" class="topic-link">数学题</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/qna">问与答</a> &nbsp;•&nbsp; <strong><a href="/member/tantaolove">tantaolove</a></strong> &nbsp;•&nbsp; 3 h 35 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/amai0w0">amai0w0</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613133#reply7" class="count_livid">7</a>
</td>
</tr>
</table>
</div>
<div class="cell item" style="">
<table cellpadding="0" cellspacing="0" border="0" width="100%">
<tr>
<td width="48" valign="top" align="center"><a href="/member/axb"><img src="//cdn.v2ex.com/avatar/9946/8002/55171_large.png?m=1456773867" class="avatar" border="0" align="default" style="max-width: 48px; max-height: 48px;" /></a></td>
<td width="10"></td>
<td width="auto" valign="middle"><span class="item_title"><a href="/t/613136#reply7" class="topic-link">[北京]微博视频平台团队招聘高级 Java 开发工程师/后台开发工程师</a></span>
<div class="sep5"></div>
<span class="topic_info"><div class="votes"></div><a class="node" href="/go/jobs">酷工作</a> &nbsp;•&nbsp; <strong><a href="/member/axb">axb</a></strong> &nbsp;•&nbsp; 3 h 13 m ago &nbsp;•&nbsp; Lastly replied by <strong><a href="/member/yizhimamong">yizhimamong</a></strong></span>
</td>
<td width="70" align="right" valign="middle">
<a href="/t/613136#reply7" class="count_livid">7</a>
</td>
</tr>
</table>
</div>
<div class="inner">
<span class="chevron">→</span> <a href="/recent">更多新主题</a>
</div>
</div>
<div class="sep20"></div>
<div class="box">
<div class="cell"><div class="fr"><a href="/planes">浏览全部节点</a></div><span class="fade"><strong>V2EX</strong> / 节点导航</span></div>
<div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">分享与探索</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/qna" style="font-size: 14px;">问与答</a>&nbsp; &nbsp; <a href="/go/share" style="font-size: 14px;">分享发现</a>&nbsp; &nbsp; <a href="/go/create" style="font-size: 14px;">分享创造</a>&nbsp; &nbsp; <a href="/go/ideas" style="font-size: 14px;">奇思妙想</a>&nbsp; &nbsp; <a href="/go/in" style="font-size: 14px;">分享邀请码</a>&nbsp; &nbsp; <a href="/go/autistic" style="font-size: 14px;">自言自语</a>&nbsp; &nbsp; <a href="/go/random" style="font-size: 14px;">随想</a>&nbsp; &nbsp; <a href="/go/design" style="font-size: 14px;">设计</a>&nbsp; &nbsp; <a href="/go/blog" style="font-size: 14px;">Blog</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">V2EX</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/v2ex" style="font-size: 14px;">V2EX</a>&nbsp; &nbsp; <a href="/go/dns" style="font-size: 14px;">DNS</a>&nbsp; &nbsp; <a href="/go/feedback" style="font-size: 14px;">反馈</a>&nbsp; &nbsp; <a href="/go/babel" style="font-size: 14px;">Project Babel</a>&nbsp; &nbsp; <a href="/go/guide" style="font-size: 14px;">使用指南</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">iOS</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/idev" style="font-size: 14px;">iDev</a>&nbsp; &nbsp; <a href="/go/icode" style="font-size: 14px;">iCode</a>&nbsp; &nbsp; <a href="/go/imarketing" style="font-size: 14px;">iMarketing</a>&nbsp; &nbsp; <a href="/go/iad" style="font-size: 14px;">iAd</a>&nbsp; &nbsp; <a href="/go/itransfer" style="font-size: 14px;">iTransfer</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">Geek</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/programmer" style="font-size: 14px;">程序员</a>&nbsp; &nbsp; <a href="/go/python" style="font-size: 14px;">Python</a>&nbsp; &nbsp; <a href="/go/android" style="font-size: 14px;">Android</a>&nbsp; &nbsp; <a href="/go/bb" style="font-size: 14px;">宽带症候群</a>&nbsp; &nbsp; <a href="/go/linux" style="font-size: 14px;">Linux</a>&nbsp; &nbsp; <a href="/go/php" style="font-size: 14px;">PHP</a>&nbsp; &nbsp; <a href="/go/cloud" style="font-size: 14px;">云计算</a>&nbsp; &nbsp; <a href="/go/outsourcing" style="font-size: 14px;">外包</a>&nbsp; &nbsp; <a href="/go/hardware" style="font-size: 14px;">硬件</a>&nbsp; &nbsp; <a href="/go/java" style="font-size: 14px;">Java</a>&nbsp; &nbsp; <a href="/go/nodejs" style="font-size: 14px;">Node.js</a>&nbsp; &nbsp; <a href="/go/server" style="font-size: 14px;">服务器</a>&nbsp; &nbsp; <a href="/go/bitcoin" style="font-size: 14px;">Bitcoin</a>&nbsp; &nbsp; <a href="/go/mysql" style="font-size: 14px;">MySQL</a>&nbsp; &nbsp; <a href="/go/programming" style="font-size: 14px;">编程</a>&nbsp; &nbsp; <a href="/go/car" style="font-size: 14px;">汽车</a>&nbsp; &nbsp; <a href="/go/docker" style="font-size: 14px;">Docker</a>&nbsp; &nbsp; <a href="/go/linode" style="font-size: 14px;">Linode</a>&nbsp; &nbsp; <a href="/go/designer" style="font-size: 14px;">设计师</a>&nbsp; &nbsp; <a href="/go/markdown" style="font-size: 14px;">Markdown</a>&nbsp; &nbsp; <a href="/go/kindle" style="font-size: 14px;">Kindle</a>&nbsp; &nbsp; <a href="/go/mongodb" style="font-size: 14px;">MongoDB</a>&nbsp; &nbsp; <a href="/go/redis" style="font-size: 14px;">Redis</a>&nbsp; &nbsp; <a href="/go/minecraft" style="font-size: 14px;">Minecraft</a>&nbsp; &nbsp; <a href="/go/tornado" style="font-size: 14px;">Tornado</a>&nbsp; &nbsp; <a href="/go/typography" style="font-size: 14px;">字体排印</a>&nbsp; &nbsp; <a href="/go/ror" style="font-size: 14px;">Ruby on Rails</a>&nbsp; &nbsp; <a href="/go/business" style="font-size: 14px;">商业模式</a>&nbsp; &nbsp; <a href="/go/math" style="font-size: 14px;">数学</a>&nbsp; &nbsp; <a href="/go/ruby" style="font-size: 14px;">Ruby</a>&nbsp; &nbsp; <a href="/go/photoshop" style="font-size: 14px;">Photoshop</a>&nbsp; &nbsp; <a href="/go/csharp" style="font-size: 14px;">C#</a>&nbsp; &nbsp; <a href="/go/sony" style="font-size: 14px;">SONY</a>&nbsp; &nbsp; <a href="/go/amazon" style="font-size: 14px;">Amazon</a>&nbsp; &nbsp; <a href="/go/nlp" style="font-size: 14px;">自然语言处理</a>&nbsp; &nbsp; <a href="/go/lego" style="font-size: 14px;">LEGO</a>&nbsp; &nbsp; <a href="/go/leetcode" style="font-size: 14px;">LeetCode</a>&nbsp; &nbsp; <a href="/go/ev" style="font-size: 14px;">电动汽车</a>&nbsp; &nbsp; <a href="/go/serverless" style="font-size: 14px;">Serverless</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">游戏</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/games" style="font-size: 14px;">游戏</a>&nbsp; &nbsp; <a href="/go/steam" style="font-size: 14px;">Steam</a>&nbsp; &nbsp; <a href="/go/ps4" style="font-size: 14px;">PlayStation 4</a>&nbsp; &nbsp; <a href="/go/lol" style="font-size: 14px;">英雄联盟</a>&nbsp; &nbsp; <a href="/go/igame" style="font-size: 14px;">iGame</a>&nbsp; &nbsp; <a href="/go/switch" style="font-size: 14px;">Nintendo Switch</a>&nbsp; &nbsp; <a href="/go/sc2" style="font-size: 14px;">StarCraft 2</a>&nbsp; &nbsp; <a href="/go/bf3" style="font-size: 14px;">Battlefield 3</a>&nbsp; &nbsp; <a href="/go/wow" style="font-size: 14px;">World of Warcraft</a>&nbsp; &nbsp; <a href="/go/5v5" style="font-size: 14px;">王者荣耀</a>&nbsp; &nbsp; <a href="/go/eve" style="font-size: 14px;">EVE</a>&nbsp; &nbsp; <a href="/go/gt" style="font-size: 14px;">Gran Turismo</a>&nbsp; &nbsp; <a href="/go/bf4" style="font-size: 14px;">Battlefield 4</a>&nbsp; &nbsp; <a href="/go/wiiu" style="font-size: 14px;">Wii U</a>&nbsp; &nbsp; <a href="/go/bfv" style="font-size: 14px;">Battlefield V</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">Apple</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/macos" style="font-size: 14px;">macOS</a>&nbsp; &nbsp; <a href="/go/iphone" style="font-size: 14px;">iPhone</a>&nbsp; &nbsp; <a href="/go/mbp" style="font-size: 14px;">MacBook Pro</a>&nbsp; &nbsp; <a href="/go/ipad" style="font-size: 14px;">iPad</a>&nbsp; &nbsp; <a href="/go/macbook" style="font-size: 14px;">MacBook</a>&nbsp; &nbsp; <a href="/go/accessory" style="font-size: 14px;">配件</a>&nbsp; &nbsp; <a href="/go/mba" style="font-size: 14px;">MacBook Air</a>&nbsp; &nbsp; <a href="/go/imac" style="font-size: 14px;">iMac</a>&nbsp; &nbsp; <a href="/go/macmini" style="font-size: 14px;">Mac mini</a>&nbsp; &nbsp; <a href="/go/macpro" style="font-size: 14px;">Mac Pro</a>&nbsp; &nbsp; <a href="/go/ipod" style="font-size: 14px;">iPod</a>&nbsp; &nbsp; <a href="/go/mobileme" style="font-size: 14px;">MobileMe</a>&nbsp; &nbsp; <a href="/go/iwork" style="font-size: 14px;">iWork</a>&nbsp; &nbsp; <a href="/go/ilife" style="font-size: 14px;">iLife</a>&nbsp; &nbsp; <a href="/go/garageband" style="font-size: 14px;">GarageBand</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">生活</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/all4all" style="font-size: 14px;">二手交易</a>&nbsp; &nbsp; <a href="/go/jobs" style="font-size: 14px;">酷工作</a>&nbsp; &nbsp; <a href="/go/afterdark" style="font-size: 14px;">天黑以后</a>&nbsp; &nbsp; <a href="/go/free" style="font-size: 14px;">免费赠送</a>&nbsp; &nbsp; <a href="/go/music" style="font-size: 14px;">音乐</a>&nbsp; &nbsp; <a href="/go/movie" style="font-size: 14px;">电影</a>&nbsp; &nbsp; <a href="/go/exchange" style="font-size: 14px;">物物交换</a>&nbsp; &nbsp; <a href="/go/tuan" style="font-size: 14px;">团购</a>&nbsp; &nbsp; <a href="/go/tv" style="font-size: 14px;">剧集</a>&nbsp; &nbsp; <a href="/go/invest" style="font-size: 14px;">投资</a>&nbsp; &nbsp; <a href="/go/creditcard" style="font-size: 14px;">信用卡</a>&nbsp; &nbsp; <a href="/go/travel" style="font-size: 14px;">旅行</a>&nbsp; &nbsp; <a href="/go/taste" style="font-size: 14px;">美酒与美食</a>&nbsp; &nbsp; <a href="/go/reading" style="font-size: 14px;">阅读</a>&nbsp; &nbsp; <a href="/go/photograph" style="font-size: 14px;">摄影</a>&nbsp; &nbsp; <a href="/go/pet" style="font-size: 14px;">宠物</a>&nbsp; &nbsp; <a href="/go/baby" style="font-size: 14px;">Baby</a>&nbsp; &nbsp; <a href="/go/soccer" style="font-size: 14px;">绿茵场</a>&nbsp; &nbsp; <a href="/go/coffee" style="font-size: 14px;">咖啡</a>&nbsp; &nbsp; <a href="/go/diary" style="font-size: 14px;">日记</a>&nbsp; &nbsp; <a href="/go/love" style="font-size: 14px;">非诚勿扰</a>&nbsp; &nbsp; <a href="/go/lohas" style="font-size: 14px;">乐活</a>&nbsp; &nbsp; <a href="/go/bike" style="font-size: 14px;">骑行</a>&nbsp; &nbsp; <a href="/go/plant" style="font-size: 14px;">植物</a>&nbsp; &nbsp; <a href="/go/mushroom" style="font-size: 14px;">蘑菇</a>&nbsp; &nbsp; <a href="/go/mileage" style="font-size: 14px;">行程控</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">Internet</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/google" style="font-size: 14px;">Google</a>&nbsp; &nbsp; <a href="/go/twitter" style="font-size: 14px;">Twitter</a>&nbsp; &nbsp; <a href="/go/coding" style="font-size: 14px;">Coding</a>&nbsp; &nbsp; <a href="/go/facebook" style="font-size: 14px;">Facebook</a>&nbsp; &nbsp; <a href="/go/wikipedia" style="font-size: 14px;">Wikipedia</a>&nbsp; &nbsp; <a href="/go/reddit" style="font-size: 14px;">reddit</a>&nbsp; &nbsp; </td></tr></table></div><div class="cell"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">城市</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/beijing" style="font-size: 14px;">北京</a>&nbsp; &nbsp; <a href="/go/shanghai" style="font-size: 14px;">上海</a>&nbsp; &nbsp; <a href="/go/shenzhen" style="font-size: 14px;">深圳</a>&nbsp; &nbsp; <a href="/go/hangzhou" style="font-size: 14px;">杭州</a>&nbsp; &nbsp; <a href="/go/chengdu" style="font-size: 14px;">成都</a>&nbsp; &nbsp; <a href="/go/guangzhou" style="font-size: 14px;">广州</a>&nbsp; &nbsp; <a href="/go/wuhan" style="font-size: 14px;">武汉</a>&nbsp; &nbsp; <a href="/go/kunming" style="font-size: 14px;">昆明</a>&nbsp; &nbsp; <a href="/go/tianjin" style="font-size: 14px;">天津</a>&nbsp; &nbsp; <a href="/go/qingdao" style="font-size: 14px;">青岛</a>&nbsp; &nbsp; <a href="/go/nyc" style="font-size: 14px;">New York</a>&nbsp; &nbsp; <a href="/go/sanfrancisco" style="font-size: 14px;">San Francisco</a>&nbsp; &nbsp; <a href="/go/la" style="font-size: 14px;">Los Angeles</a>&nbsp; &nbsp; <a href="/go/boston" style="font-size: 14px;">Boston</a>&nbsp; &nbsp; </td></tr></table></div><div class="inner"><table cellpadding="0" cellspacing="0" border="0"><tr><td align="right" width="80"><span class="fade">品牌</span></td><td style="line-height: 200%; padding-left: 10px; word-break: keep-all;"><a href="/go/uniqlo" style="font-size: 14px;">UNIQLO</a>&nbsp; &nbsp; <a href="/go/lamy" style="font-size: 14px;">Lamy</a>&nbsp; &nbsp; <a href="/go/ikea" style="font-size: 14px;">宜家</a>&nbsp; &nbsp; <a href="/go/muji" style="font-size: 14px;">无印良品</a>&nbsp; &nbsp; <a href="/go/nike" style="font-size: 14px;">Nike</a>&nbsp; &nbsp; <a href="/go/gap" style="font-size: 14px;">Gap</a>&nbsp; &nbsp; <a href="/go/moleskine" style="font-size: 14px;">Moleskine</a>&nbsp; &nbsp; <a href="/go/adidas" style="font-size: 14px;">Adidas</a>&nbsp; &nbsp; <a href="/go/gstar" style="font-size: 14px;">G-Star</a>&nbsp; &nbsp; </td></tr></table></div>
</div>
</div>
</div>
<div class="c"></div>
<div class="sep20"></div>
</div>
<div id="Bottom">
<div class="content">
<div class="inner">
<div class="sep10"></div>
<div class="fr">
<a href="https://www.digitalocean.com/?refcode=1b51f1a7651d" target="_blank"><div id="DigitalOcean"></div></a>
</div>
<strong><a href="/about" class="dark" target="_self">About</a> &nbsp; <span class="snow">·</span> &nbsp; <a href="/faq" class="dark" target="_self">FAQ</a> &nbsp; <span class="snow">·</span> &nbsp; <a href="/p/7v9TEc53" class="dark" target="_self">API</a> &nbsp; <span class="snow">·</span> &nbsp; <a href="/mission" class="dark" target="_self">Mission</a> &nbsp; <span class="snow">·</span> &nbsp; <a href="/advertise" class="dark" target="_self">Advertise</a> &nbsp; <span class="snow">·</span> &nbsp; <a href="/advertise/2017.html" class="dark" target="_self">Thanks</a> &nbsp; <span class="snow">·</span> &nbsp; <a href="/tools" class="dark" target="_self">Tools</a> &nbsp; <span class="snow">·</span> &nbsp; 2387 Online</strong> &nbsp; <span class="fade">Highest 5043</span> &nbsp; <span class="snow">·</span> &nbsp; <a href="/select/language" class="f11"><img src="/static/img/language.png?v=6a5cfa731dc71a3769f6daace6784739" width="16" align="absmiddle" id="ico-select-language" /> &nbsp; Select Language</a>
<div class="sep20"></div>
创意工作者们的社区
<div class="sep5"></div>
World is powered by solitude
<div class="sep20"></div>
<span class="small fade">VERSION: 3.9.8.3 · 10ms · UTC 09:36 · PVG 17:36 · LAX 02:36 · JFK 05:36<br />♥ Do have faith in what you're doing.</span>
<div class="sep10"></div>
</div>
</div>
</div>
<script>
	  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
	  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
	  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
	  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

	  ga('create', 'UA-11940834-2', 'v2ex.com');
	  ga('send', 'pageview');
      

	</script>
</body>
</html>`

var html76 = `


<!DOCTYPE html>
<html lang="zh-cmn-Hans" class="ua-windows ua-webkit book-new-nav">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <title>TIGER×DRAGON1! (豆瓣)</title>
  
<script>!function(e){var o=function(o,n,t){var c,i,r=new Date;n=n||30,t=t||"/",r.setTime(r.getTime()+24*n*60*60*1e3),c="; expires="+r.toGMTString();for(i in o)e.cookie=i+"="+o[i]+c+"; path="+t},n=function(o){var n,t,c,i=o+"=",r=e.cookie.split(";");for(t=0,c=r.length;t<c;t++)if(n=r[t].replace(/^\s+|\s+$/g,""),0==n.indexOf(i))return n.substring(i.length,n.length).replace(/\"/g,"");return null},t=e.write,c={"douban.com":1,"douban.fm":1,"google.com":1,"google.cn":1,"googleapis.com":1,"gmaptiles.co.kr":1,"gstatic.com":1,"gstatic.cn":1,"google-analytics.com":1,"googleadservices.com":1},i=function(e,o){var n=new Image;n.onload=function(){},n.src="https://www.douban.com/j/except_report?kind=ra022&reason="+encodeURIComponent(e)+"&environment="+encodeURIComponent(o)},r=function(o){try{t.call(e,o)}catch(e){t(o)}},a=/<script.*?src\=["']?([^"'\s>]+)/gi,g=/http:\/\/(.+?)\.([^\/]+).+/i;e.writeln=e.write=function(e){var t,l=a.exec(e);return l&&(t=g.exec(l[1]))?c[t[2]]?void r(e):void("tqs"!==n("hj")&&(i(l[1],location.href),o({hj:"tqs"},1),setTimeout(function(){location.replace(location.href)},50))):void r(e)}}(document);
</script>

  
  <meta http-equiv="Pragma" content="no-cache">
  <meta http-equiv="Expires" content="Sun, 6 Mar 2005 01:00:00 GMT">
  
<meta http-equiv="mobile-agent" content="format=html5; url=https://m.douban.com/book/subject/6438136/">
<meta name="keywords" content="TIGER×DRAGON1!,[日] 竹宫悠由子 著,[日] 泰 绘,湖南美术出版社,2011-7,简介,作者,书评,论坛,推荐,二手">
<meta name="description" content="图书TIGER×DRAGON1! 介绍、书评、论坛及推荐 ">

  <script>var _head_start = new Date();</script>
  
  <link href="https://img3.doubanio.com/f/book/e7612a013a5e76c7c680323c74748d21cd703ba0/css/book/master.css" rel="stylesheet" type="text/css">

  <link href="https://img3.doubanio.com/f/book/222a5c61e041638af8defc87cf97f4a863a77922/css/book/base/init.css" rel="stylesheet">
  <style type="text/css"></style>
  <script src="https://img3.doubanio.com/f/book/0495cb173e298c28593766009c7b0a953246c5b5/js/book/lib/jquery/jquery.js"></script>
  <script src="https://img3.doubanio.com/f/book/dfb2bae8d1c310e7d5d00e697aaa9891d2cb4eb8/js/book/master.js"></script>
  

  
  <link rel="stylesheet" href="https://img3.doubanio.com/f/book/1eabe3f4e416e77dbafd2ef9bc830fa2ac7a8d4c/css/book/subject.css">
  <link href="https://img3.doubanio.com/f/book/5d301503fbbd8e09f3114583859789884e942f47/css/book/annotation/like.css" rel="stylesheet">
  <script src="https://img3.doubanio.com/f/shire/3c6f2946669cfb2fc9ee4a4d1dcc41fc181cad92/js/lib/jquery.snippet.js"></script>
  <script src="https://img3.doubanio.com/f/shire/77323ae72a612bba8b65f845491513ff3329b1bb/js/do.js" data-cfg-autoload="false"></script>
  <script src="https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js"></script>
  <script src="https://img3.doubanio.com/f/book/2e421e5ec8f2869d31535206c0ac0322532be1f8/js/book/mod/hide.js"></script>
  <script src="https://img3.doubanio.com/f/book/cc6b1a77c3812c7dd20b0374332fade081e1c0b0/js/book/subject/unfold.js"></script>
    <link rel="alternate" href="https://book.douban.com/feed/subject/6438136/reviews" type="application/rss+xml" title="RSS">
  <style type="text/css"> h2 {color: #007722;} </style>
  <script type='text/javascript'>
    var _vds = _vds || [];
    (function(){ _vds.push(['setAccountId', '22c937bbd8ebd703f2d8e9445f7dfd03']);
        _vds.push(['setCS1','user_id','1']);
            (function() {var vds = document.createElement('script');
                vds.type='text/javascript';
                vds.async = true;
                vds.src = ('https:' == document.location.protocol ? 'https://' : 'http://') + 'dn-growing.qbox.me/vds.js';
                var s = document.getElementsByTagName('script')[0];
                s.parentNode.insertBefore(vds, s);
            })();
    })();
</script>

  
  <script type='text/javascript'>
    var _vwo_code=(function(){
      var account_id=249272,
          settings_tolerance=2000,
          library_tolerance=2500,
          use_existing_jquery=false,
          // DO NOT EDIT BELOW THIS LINE
          f=false,d=document;return{use_existing_jquery:function(){return use_existing_jquery;},library_tolerance:function(){return library_tolerance;},finish:function(){if(!f){f=true;var a=d.getElementById('_vis_opt_path_hides');if(a)a.parentNode.removeChild(a);}},finished:function(){return f;},load:function(a){var b=d.createElement('script');b.src=a;b.type='text/javascript';b.innerText;b.onerror=function(){_vwo_code.finish();};d.getElementsByTagName('head')[0].appendChild(b);},init:function(){settings_timer=setTimeout('_vwo_code.finish()',settings_tolerance);var a=d.createElement('style'),b='body{opacity:0 !important;filter:alpha(opacity=0) !important;background:none !important;}',h=d.getElementsByTagName('head')[0];a.setAttribute('id','_vis_opt_path_hides');a.setAttribute('type','text/css');if(a.styleSheet)a.styleSheet.cssText=b;else a.appendChild(d.createTextNode(b));h.appendChild(a);this.load('//dev.visualwebsiteoptimizer.com/j.php?a='+account_id+'&u='+encodeURIComponent(d.URL)+'&r='+Math.random());return settings_timer;}};}());_vwo_settings_timer=_vwo_code.init();
  </script>

  


<script type="application/ld+json">
{
  "@context":"http://schema.org",
  "@type":"Book",
  "workExample": [],
  "name" : "TIGER×DRAGON1!",
  "author": 
  [
    {
      "@type": "Person",
      "name": "[日] 竹宫悠由子 著"
    }
    ,
    {
      "@type": "Person",
      "name": "[日] 泰 绘"
    }
  ]
,
  "url" : "https://book.douban.com/subject/6438136/",
  "isbn" : "9787535644954",
  "sameAs": "https://book.douban.com/subject/6438136/"
}
</script>


  <script>  </script>
  <link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/3cd6e0524025c20.css">

  <link rel="shortcut icon" href="https://img3.doubanio.com/favicon.ico" type="image/x-icon">
</head>
<body>
  
    <script>var _body_start = new Date();</script>
    
  



    <link href="//img3.doubanio.com/dae/accounts/resources/0246c88/shire/bundle.css" rel="stylesheet" type="text/css">



<div id="db-global-nav" class="global-nav">
  <div class="bd">
    
<div class="top-nav-info">
  <ul>
    <li>
    <a id="top-nav-doumail-link" href="https://www.douban.com/doumail/">豆邮</a>
    </li>
    <li class="nav-user-account">
      <a target="_blank" href="https://accounts.douban.com/passport/setting/" class="bn-more">
        <span>chenset的帐号</span><span class="arrow"></span>
      </a>
      <div class="more-items">
        <table cellpadding="0" cellspacing="0">
          <tbody>
            <tr>
              <td>
                <a href="https://www.douban.com/mine/">个人主页</a>
              </td>
            </tr>
            <tr>
              <td>
                <a target="_blank" href="https://www.douban.com/mine/orders/">我的订单</a>
              </td>
            </tr>
            <tr>
              <td>
                <a target="_blank" href="https://www.douban.com/mine/wallet/">我的钱包</a>
              </td>
            </tr>
            <tr>
              <td>
                <a target="_blank" href="https://accounts.douban.com/passport/setting/">帐号管理</a>
              </td>
            </tr>
            <tr>
              <td>
                <a href="https://www.douban.com/accounts/logout?source=book&ck=AsnX">退出</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </li>
  </ul>
</div>

  <div class="top-nav-reminder">
    <a href="https://www.douban.com/notification/" class="lnk-remind">提醒</a>
    <div id="top-nav-notimenu" class="more-items">
      <div class="bd">
        <p>加载中...</p>
      </div>
    </div>
  </div>

    <div class="top-nav-doubanapp">
  <a href="https://www.douban.com/doubanapp/app?channel=top-nav" class="lnk-doubanapp">下载豆瓣客户端</a>
  <div id="doubanapp-tip">
    <a href="https://www.douban.com/doubanapp/app?channel=qipao" class="tip-link">豆瓣 <span class="version">6.0</span> 全新发布</a>
    <a href="javascript: void 0;" class="tip-close">×</a>
  </div>
  <div id="top-nav-appintro" class="more-items">
    <p class="appintro-title">豆瓣</p>
    <p class="qrcode">扫码直接下载</p>
    <div class="download">
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=iOS">iPhone</a>
      <span>·</span>
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=Android" class="download-android">Android</a>
    </div>
  </div>
</div>

    


<div class="global-nav-items">
  <ul>
    <li class="">
      <a href="https://www.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-main&quot;,&quot;uid&quot;:&quot;54106750&quot;}">豆瓣</a>
    </li>
    <li class="on">
      <a href="https://book.douban.com"  data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-book&quot;,&quot;uid&quot;:&quot;54106750&quot;}">读书</a>
    </li>
    <li class="">
      <a href="https://movie.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-movie&quot;,&quot;uid&quot;:&quot;54106750&quot;}">电影</a>
    </li>
    <li class="">
      <a href="https://music.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-music&quot;,&quot;uid&quot;:&quot;54106750&quot;}">音乐</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/location" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-location&quot;,&quot;uid&quot;:&quot;54106750&quot;}">同城</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/group" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-group&quot;,&quot;uid&quot;:&quot;54106750&quot;}">小组</a>
    </li>
    <li class="">
      <a href="https://read.douban.com&#47;?dcs=top-nav&amp;dcm=douban" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-read&quot;,&quot;uid&quot;:&quot;54106750&quot;}">阅读</a>
    </li>
    <li class="">
      <a href="https://douban.fm&#47;?from_=shire_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-fm&quot;,&quot;uid&quot;:&quot;54106750&quot;}">FM</a>
    </li>
    <li class="">
      <a href="https://time.douban.com&#47;?dt_time_source=douban-web_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-time&quot;,&quot;uid&quot;:&quot;54106750&quot;}">时间</a>
    </li>
    <li class="">
      <a href="https://market.douban.com&#47;?utm_campaign=douban_top_nav&amp;utm_source=douban&amp;utm_medium=pc_web" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-market&quot;,&quot;uid&quot;:&quot;54106750&quot;}">豆品</a>
    </li>
    <li>
      <a href="#more" class="bn-more"><span>更多</span></a>
      <div class="more-items">
        <table cellpadding="0" cellspacing="0">
          <tbody>
            <tr>
              <td>
                <a href="https://ypy.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-ypy&quot;,&quot;uid&quot;:&quot;54106750&quot;}">豆瓣摄影</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </li>
  </ul>
</div>

  </div>
</div>
<script>
  ;window._GLOBAL_NAV = {
    USER_ID: "54106750",
    UPLOAD_AUTH_TOKEN: "54106750:ab8488e4ea59a32c52ecea7f517cdb738d613760",
    SSE_TOKEN: "e50424ea3edf5e36839b22ccd5cc47c019cc2246",
    SSE_TIMESTAMP: "1550288086",
    DOUBAN_URL: "https://www.douban.com",
    N_NEW_NOTIS: 0,
    N_NEW_DOUMAIL: 0
  };
</script>



    <script src="//img3.doubanio.com/dae/accounts/resources/0246c88/shire/bundle.js" defer="defer"></script>




  



    <link href="//img3.doubanio.com/dae/accounts/resources/0246c88/book/bundle.css" rel="stylesheet" type="text/css">




<div id="db-nav-book" class="nav">
  <div class="nav-wrap">
  <div class="nav-primary">
    <div class="nav-logo">
      <a href="https:&#47;&#47;book.douban.com">豆瓣读书</a>
    </div>
    <div class="nav-search">
      <form action="https:&#47;&#47;book.douban.com/subject_search" method="get">
        <fieldset>
          <legend>搜索：</legend>
          <label for="inp-query">
          </label>
          <div class="inp"><input id="inp-query" name="search_text" size="22" maxlength="60" placeholder="书名、作者、ISBN" value=""></div>
          <div class="inp-btn"><input type="submit" value="搜索"></div>
          <input type="hidden" name="cat" value="1001" />
        </fieldset>
      </form>
    </div>
  </div>
  </div>
  <div class="nav-secondary">
    

<div class="nav-items">
  <ul>
    <li    ><a href="https://book.douban.com/mine"
     >我读</a>
    </li>
    <li    ><a href="https://book.douban.com/updates"
     >动态</a>
    </li>
    <li    ><a href="https://book.douban.com/recommended"
     >豆瓣猜</a>
    </li>
    <li    ><a href="https://book.douban.com/tag/"
     >分类浏览</a>
    </li>
    <li    ><a href="https://book.douban.com/cart/"
     >购书单</a>
    </li>
    <li    ><a href="https://read.douban.com/ebooks/?dcs=book-nav&dcm=douban"
            target="_blank"
     >电子图书</a>
    </li>
    <li    ><a href="https://market.douban.com/book?utm_campaign=book_nav_freyr&utm_source=douban&utm_medium=pc_web"
     >豆瓣书店</a>
    </li>
    <li    ><a href="https://book.douban.com/annual/2018?source=navigation"
            target="_blank"
     >2018年度榜单</a>
    </li>
    <li    ><a href="https://www.douban.com/standbyme/2018?source=navigation"
            target="_blank"
     >2018书影音报告</a>
    </li>
    <li          class=" book-cart"
    ><a href="https://market.douban.com/cart/?biz_type=book&utm_campaign=book_nav_cart&utm_source=douban&utm_medium=pc_web"
            target="_blank"
     >购物车</a>
    </li>
  </ul>
</div>

    <a href="https://book.douban.com/annual/2018?source=book_navigation" class="bookannual2018"></a>
  </div>
</div>

<script id="suggResult" type="text/x-jquery-tmpl">
  <li data-link="{{= url}}">
            <a href="{{= url}}" onclick="moreurl(this, {from:'book_search_sugg', query:'{{= keyword }}', subject_id:'{{= id}}', i: '{{= index}}', type: '{{= type}}'})">
            <img src="{{= pic}}" width="40" />
            <div>
                <em>{{= title}}</em>
                {{if year}}
                    <span>{{= year}}</span>
                {{/if}}
                <p>
                {{if type == "b"}}
                    {{= author_name}}
                {{else type == "a" }}
                    {{if en_name}}
                        {{= en_name}}
                    {{/if}}
                {{/if}}
                 </p>
            </div>
        </a>
        </li>
  </script>




    <script src="//img3.doubanio.com/dae/accounts/resources/0246c88/book/bundle.js" defer="defer"></script>





    <div id="wrapper">
        
    <div id="dale_book_subject_top_icon"></div>
<h1>
    <span property="v:itemreviewed">TIGER×DRAGON1!</span>
    <div class="clear"></div>
</h1>

        
  <div id="content">
    
    <div class="grid-16-8 clearfix">
      
      <div class="article">



<div class="indent">
  <div class="subjectwrap clearfix">
    



<div class="subject clearfix">
<div id="mainpic" class="">

  

  <a class="nbg"
      href="https://img3.doubanio.com/view/subject/l/public/s6439936.jpg" title="TIGER×DRAGON1!">
    <img src="https://img3.doubanio.com/view/subject/l/public/s6439936.jpg" title="点击看大图" alt="TIGER×DRAGON1!"
         rel="v:photo" style="width: 135px;max-height: 200px;">
  </a>

    <br>
      <p class="gact">
        <a href="https://book.douban.com/subject/6438136/edit">更新图书信息或封面</a>
      </p>


</div>





<div id="info" class="">



    
    
  
    <span>
      <span class="pl"> 作者</span>:
        
        
        <a class="" href="/search/%E7%AB%B9%E5%AE%AB%E6%82%A0%E7%94%B1%E5%AD%90%20%E8%91%97">[日] 竹宫悠由子 著</a>
        
           /
        
        <a class="" href="/search/%E6%B3%B0%20%E7%BB%98">[日] 泰 绘</a>
    </span><br/>

    
    
  
    <span class="pl">出版社:</span> 湖南美术出版社<br/>

    
    
  

    
    
  
    <span class="pl">副标题:</span> 龙与虎<br/>

    
    
  
    <span class="pl">原作名:</span> とらドラ!<br/>

    
    
  
    <span>
      <span class="pl"> 译者</span>:
        
        
        <a class="" href="/search/%E9%BB%84%E8%96%87%E5%AB%94">黄薇嫔</a>
    </span><br/>

    
    
  
    <span class="pl">出版年:</span> 2011-7<br/>

    
    
  
    <span class="pl">页数:</span> 211<br/>

    
    
  
    <span class="pl">定价:</span> 20.00元<br/>

    
    
  
    <span class="pl">装帧:</span> 平装<br/>

    
    
  
    <span class="pl">丛书:</span>&nbsp;<a href="https://book.douban.com/series/8715">TIGER×DRAGON! 龙与虎</a><br>

    
    
  
    
      
      <span class="pl">ISBN:</span> 9787535644954<br/>


</div>

</div>
























    





<div id="interest_sectl" class="">
  <div class="rating_wrap clearbox" rel="v:rating">
    <div class="rating_logo">豆瓣评分</div>
    <div class="rating_self clearfix" typeof="v:Rating">
      <strong class="ll rating_num " property="v:average"> 8.3 </strong>
      <span property="v:best" content="10.0"></span>
      <div class="rating_right ">
          <div class="ll bigstar40"></div>
            <div class="rating_sum">
                <span class="">
                    <a href="collections" class="rating_people"><span property="v:votes">194</span>人评价</a>
                </span>
            </div>


      </div>
    </div>
          
            
            
<span class="stars5 starstop" title="力荐">
    5星
</span>

            
<div class="power" style="width:64px"></div>

            <span class="rating_per">44.3%</span>
            <br>
            
            
<span class="stars4 starstop" title="推荐">
    4星
</span>

            
<div class="power" style="width:55px"></div>

            <span class="rating_per">38.7%</span>
            <br>
            
            
<span class="stars3 starstop" title="还行">
    3星
</span>

            
<div class="power" style="width:21px"></div>

            <span class="rating_per">14.9%</span>
            <br>
            
            
<span class="stars2 starstop" title="较差">
    2星
</span>

            
<div class="power" style="width:2px"></div>

            <span class="rating_per">2.1%</span>
            <br>
            
            
<span class="stars1 starstop" title="很差">
    1星
</span>

            
<div class="power" style="width:0px"></div>

            <span class="rating_per">0.0%</span>
            <br>
    </div>
        
</div>

  </div>
  





  
    
    <div id="interest_sect_level" class="clearfix">
        <a href="#" rel="nofollow" class="collect_btn colbutt ll" name="pbtn-6438136-wish">
          <span>
            
<form method="POST" action="https://book.douban.com/subject/6438136/?interest=wish" class="miniform"><div style="display:none;"><input type="hidden" name="ck" value="AsnX"/></div>
    <input type="submit" class="minisubmit j " value="想读" title="" />
</form>

          </span>
        </a>
        <a href="#" rel="nofollow" class="collect_btn colbutt ll" name="pbtn-6438136-do">
          <span>
            
<form method="POST" action="https://book.douban.com/subject/6438136/?interest=do" class="miniform"><div style="display:none;"><input type="hidden" name="ck" value="AsnX"/></div>
    <input type="submit" class="minisubmit j " value="在读" title="" />
</form>

          </span>
        </a>
        <a href="#" rel="nofollow" class="collect_btn colbutt ll" name="pbtn-6438136-collect">
          <span>
            
<form method="POST" action="https://book.douban.com/subject/6438136/?interest=collect" class="miniform"><div style="display:none;"><input type="hidden" name="ck" value="AsnX"/></div>
    <input type="submit" class="minisubmit j " value="读过" title="" />
</form>

          </span>
        </a>
      <div class="ll j a_stars">
        
    
    评价:
    <span id="rating"> <span id="stars" data-solid="https://img3.doubanio.com/f/shire/5a2327c04c0c231bced131ddf3f4467eb80c1c86/pics/rating_icons/star_onmouseover.png" data-hollow="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" data-solid-2x="https://img3.doubanio.com/f/shire/7258904022439076d57303c3b06ad195bf1dc41a/pics/rating_icons/star_onmouseover@2x.png" data-hollow-2x="https://img3.doubanio.com/f/shire/95cc2fa733221bb8edd28ad56a7145a5ad33383e/pics/rating_icons/star_hollow_hover@2x.png">

                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-6438136-collect-1">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star1" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-6438136-collect-2">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star2" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-6438136-collect-3">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star3" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-6438136-collect-4">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star4" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-6438136-collect-5">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star5" width="16" height="16"/></a>
    </span><span id="rateword" class="pl"></span>
    <input id="n_rating" type="hidden" value=""  />
    </span>

      </div>
      

    </div>



  
  <div class="gtleft">
    <ul class="ul_subject_menu bicelink color_gray pt6 clearfix">
        <li>
          <img src="https://img3.doubanio.com/f/shire/5bbf02b7b5ec12b23e214a580b6f9e481108488c/pics/add-review.gif" />&nbsp;
          <a class="" href="https://book.douban.com/annotation/write?sid=6438136" rel="nofollow">写笔记</a>
        </li>

          <li>
            <img src="https://img3.doubanio.com/f/shire/5bbf02b7b5ec12b23e214a580b6f9e481108488c/pics/add-review.gif" />&nbsp;<a class="" href="https://book.douban.com/subject/6438136/new_review" rel="nofollow">写书评</a>
          </li>

      <li>

  <span class="rr">
  

    <img src="https://img3.doubanio.com/pics/add-cart.gif"/>
        <a name="6438136" class="j a_add2cart" href="javascript:;">加入购书单</a>
  <span class="hidden">已在<a href="https://book.douban.com/cart">购书单</a></span>
</span><br class="clearfix" />
</li>

        <li>
          


    <div class="doulist-add-btn">
  

  

  
  <a href="javascript:void(0)"
     data-id="6438136"
     data-cate="1001"
     data-canview="True"
     data-url="https://book.douban.com/subject/6438136/"
     data-catename="图书"
     data-link="https://www.douban.com/people/chenset/doulists/all?add=6438136&amp;cat=1001"
     data-title="TIGER×DRAGON1!"
     data-picture="https://img3.doubanio.com/view/subject/l/public/s6439936.jpg"
     class="lnk-doulist-add"
     onclick="moreurl(this, { 'from':'doulist-btn-1001-6438136-54106750'})">
      <i></i>添加到豆列
  </a>
    </div>

        </li>

        
        
    
    <li class="rec" id="C-6438136">
        <a href="#" data-url="https://book.douban.com/subject/6438136/" data-desc="" data-title="书籍《TIGER×DRAGON1!》 (来自豆瓣) " data-pic="https://img3.doubanio.com/view/subject/l/public/s6439936.jpg" class="bn-sharing ">分享到</a> &nbsp;&nbsp;
    </li>
    <script>
    var __cache_url = __cache_url || {};
    (function(u){
        if(__cache_url[u]) return;
        __cache_url[u] = true;
        window.DoubanShareIcons = 'https://img3.doubanio.com/f/shire/d15ffd71f3f10a7210448fec5a68eaec66e7f7d0/pics/ic_shares.png';
        var initShareButton = function() {
          $.ajax({url:u,dataType:'script',cache:true});
        };
        if (typeof Do == 'function' && 'ready' in Do) {
          Do('https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
            'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
            initShareButton);
        } else if(typeof Douban == 'object' && 'loader' in Douban) {
          Douban.loader.batch(
            'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
            'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js'
          ).done(initShareButton);
        }
    })('https://img3.doubanio.com/f/shire/6e6a5f21daeec19bbb41bf48c07fccaa4dad4d98/js/lib/sharebutton.js');
    </script>

    </ul>
  </div>


    






    <script>
        var rec_url = 'https://www.douban.com/share/recommend?'
    </script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/9238d0db7d5fc042186ec54ebc5b4e64653f5e46/js/dshare.js"></script>

<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/share/recommend?sanity_key=_8c67a&amp;apikey=&amp;name=TIGER%C3%97DRAGON1%21&amp;image=https%3A%2F%2Fimg3.doubanio.com%2Fview%2Fsubject%2Fs%2Fpublic%2Fs6439936.jpg&amp;redir=https%3A%2F%2Fbook.douban.com%2Fsubject%2F6438136%2F&amp;href=https%3A%2F%2Fbook.douban.com%2Fsubject%2F6438136%2F&amp;curl=&amp;type=book&amp;properties=%7B%22rating%22%3A%228.3%22%7D&amp;desc=%E9%BE%99%E4%B8%8E%E8%99%8E+%2F+%5B%E6%97%A5%5D+%E7%AB%B9%E5%AE%AB%E6%82%A0%E7%94%B1%E5%AD%90+%E8%91%97+%2F+%E6%B9%96%E5%8D%97%E7%BE%8E%E6%9C%AF%E5%87%BA%E7%89%88%E7%A4%BE" share-id="6438136" data-mode="plain" data-name="TIGER×DRAGON1!" data-type="book" data-desc="龙与虎 / [日] 竹宫悠由子 著 / 湖南美术出版社" data-href="https://book.douban.com/subject/6438136/" data-image="https://img3.doubanio.com/view/subject/s/public/s6439936.jpg" data-properties="{&#34;rating&#34;:&#34;8.3&#34;}" data-redir="https://book.douban.com/subject/6438136/vote?ck=AsnX" data-text="" data-apikey="" data-curl="" data-count="10" data-object_kind="1001" data-object_id="6438136" data-target_type="rec" data-target_action="0" data-action_props="{&#34;subject_url&#34;:&#34;https:\/\/book.douban.com\/subject\/6438136\/&#34;,&#34;subject_title&#34;:&#34;TIGER×DRAGON1!&#34;}" data-sanity_key="_8c67a" class="lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>


<script>
  //bind events for collection button.
  $('.collect_btn', '#interest_sect_level').each(function(){
      Douban.init_collect_btn(this);
  });
</script>








</div>

<br clear="all">
<div id="collect_form_6438136"></div>
<div class="related_info">
  






  

  <h2>
    <span class="">内容简介</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>



<div class="indent" id="link-report">
    
      <div class="">
        <style type="text/css" media="screen">
.intro p{text-indent:2em;word-break:normal;}
</style>
<div class="intro">
    <p>樱花飞舞的四月，全新的高中二年级，新的班级……</p>    <p>眼神凶恶的普通男孩——高须龙儿，遇上了众人害怕、凶暴残忍的“掌中老虎”逢坂大河，还知道了她不欲为人知的秘密。这就是龙虎相争的爱与战斗的序幕！</p>    <p>再加上总是笑眯眯的超自我女孩栉枝实乃梨，以及文武双全、勤勉正直，不过却有点怪怪的眼镜班长北村佑作。这群怪家伙的恋爱哪有可能这么顺利呢……</p>    <p>由《我们俩的田村同学》作者竹宫悠由子以及插画家泰送上的超强恋爱喜剧在此登场！</p></div>

      </div>
    
<link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/c4c6dd266f58b41cbeebc9e4e6d7dd7b2a5c3711/css/report.css" />
<link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css" />
<link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/b45aa277f8b8df40596b96582dafb1ed0a899a64/css/report_dialog.css" />
<script type="text/javascript" src="https://img3.doubanio.com/f/shire/a14501790b4a2db257dc5be5e37d820e600703c6/js/report_dialog.js"></script>
<style>
    #link-report .report { text-align: right; font-size: 12px; visibility: hidden; }
    #link-report .report a { color: #BBB; }
    #link-report .report a:hover { color: #FFF; background-color: #BBB; }
</style>
<script>
    Do = (typeof Do === 'undefined')? $ : Do;
    Do(function(){
    $("body").delegate("#link-report", 'mouseenter mouseleave', function(e){
      switch (e.type) {
        case "mouseenter":
          $(this).find(".report").css('visibility', 'visible');
          break;
        case "mouseleave":
          $(this).find(".report").css('visibility', 'hidden');
          break;
      }
    });
    $("#link-report").delegate(".report a", 'click', function(e){
        e.preventDefault();
        var opt = "";
        var obj = $(e.target).closest('#link-report');
        var id = obj.length != 0 ? obj.data("id") : undefined;
        var params = (opt&&id) ? ''.concat(opt, '=', id) : '';

        var url = "https://book.douban.com/subject/6438136/";
        url += (~url.indexOf('?') ? '&' : '?') + params
        url = url.replace(/\&+/g, '&')
        generate_report_dialog({report_url: url, type: 'subject'});
    });

    $("#link-report").append('<div class="report"><a rel="nofollow" href="#">举报</a></div>');
  });
</script>

</div>

  

























    
  

  <h2>
    <span class="">作者简介</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>



      <div class="indent ">
          
            <div class="">
            <style type="text/css" media="screen">
.intro p{text-indent:2em;word-break:normal;}
</style>
<div class="intro">
    <p>作者简介：</p>    <p>竹宫悠由子，女，日本著名轻小说家，1978年2月24日出生于日本枥木县，目前居住在东京。过去曾任职于电脑游戏公司，离职后朝作家之路发展，以《我们俩的田村同学》一书出道。善于描写平凡中带有不平凡的校园恋爱喜剧，并深受读者的喜爱。</p>    <p></p>    <p>插画家介绍：</p>    <p>泰，男，日本插画家，1984年8月3日出生于日本德岛县，现居住在东京。代表作有《我们俩的田村同学》系列、《龙与虎》系列和《狂恋声优宝贝》系列等。</p></div>

            </div>
      </div>











































  

  


  


<link rel="stylesheet" href="https://img3.doubanio.com/f/verify/16c7e943aee3b1dc6d65f600fcc0f6d62db7dfb4/entry_creator/dist/author_subject/style.css">

<div id="author_subject" class="author-wrapper">
  <div class="loading"></div>
</div>

<script type="text/javascript">
  var answerObj = {
    TYPE: 'book',
    SUBJECT_ID: '6438136',
    ISALL: 'False' || false,
    USER_ID: '54106750'
  }
</script>
<script src="https://img3.doubanio.com/f/book/61252f2f9b35f08b37f69d17dfe48310dd295347/js/book/lib/react/bundle.js"></script>
<script type="text/javascript" src="https://img3.doubanio.com/f/verify/ac140ef86262b845d2be7b859e352d8196f3f6d4/entry_creator/dist/author_subject/index.js"></script> 
  






<div id="db-tags-section" class="blank20">
  
  

  <h2>
    <span class="">豆瓣成员常用的标签(共33个)</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


  <div class="indent">    <span class="">
        <a class="  tag" href="/tag/轻小说">轻小说</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/龙与虎">龙与虎</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/竹宫悠由子">竹宫悠由子</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/日本">日本</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/天闻角川">天闻角川</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/青春">青春</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/纯爱">纯爱</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/小说">小说</a> &nbsp;    </span>
  </div>
</div>


  


<div class="subject_show block5">
<h2>丛书信息</h2>
<div>
　　<a href="https://book.douban.com/series/8715">TIGER×DRAGON! 龙与虎 (共9册)</a>,
这套丛书还有
《TIGER×DRAGON5!》,《TIGER×DRAGON8!》,《TIGER×DRAGON3!》,《TIGER×DRAGON7!》,《TIGER×DRAGON6!》    等。</div>
</div>
<script>
$(function(){$(".knnlike a").click(function(){return moreurl(this,{'from':'knnlike'})})})
</script>

  





<div id="db-rec-section" class="block5 subject_show knnlike">
  
  
  

  <h2>
    <span class="">喜欢读&#34;TIGER×DRAGON1!&#34;的人也喜欢</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


  <div class="content clearfix">
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/6774105/" onclick="moreurl(this, {'total': 10, 'clicked': '6774105', 'pos': 0, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s6829125.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/6774105/" onclick="moreurl(this, {'total': 10, 'clicked': '6774105', 'pos': 0, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON2!
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/6828272/" onclick="moreurl(this, {'total': 10, 'clicked': '6828272', 'pos': 1, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s6917574.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/6828272/" onclick="moreurl(this, {'total': 10, 'clicked': '6828272', 'pos': 1, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON3!
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/6946920/" onclick="moreurl(this, {'total': 10, 'clicked': '6946920', 'pos': 2, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s6986347.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/6946920/" onclick="moreurl(this, {'total': 10, 'clicked': '6946920', 'pos': 2, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON4!
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/10574336/" onclick="moreurl(this, {'total': 10, 'clicked': '10574336', 'pos': 3, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s9065875.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/10574336/" onclick="moreurl(this, {'total': 10, 'clicked': '10574336', 'pos': 3, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON5!
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/11502034/" onclick="moreurl(this, {'total': 10, 'clicked': '11502034', 'pos': 4, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s11126079.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/11502034/" onclick="moreurl(this, {'total': 10, 'clicked': '11502034', 'pos': 4, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON6!
          </a>
        </dd>
      </dl>
        <dl class="clear"></dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/20397802/" onclick="moreurl(this, {'total': 10, 'clicked': '20397802', 'pos': 5, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s24521542.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/20397802/" onclick="moreurl(this, {'total': 10, 'clicked': '20397802', 'pos': 5, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON8!
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/24853078/" onclick="moreurl(this, {'total': 10, 'clicked': '24853078', 'pos': 6, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s26837182.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/24853078/" onclick="moreurl(this, {'total': 10, 'clicked': '24853078', 'pos': 6, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON9!
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/11611958/" onclick="moreurl(this, {'total': 10, 'clicked': '11611958', 'pos': 7, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s22704582.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/11611958/" onclick="moreurl(this, {'total': 10, 'clicked': '11611958', 'pos': 7, 'identifier': 'book-rec-books'})" class="">
            TIGER×DRAGON7!
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/5383590/" onclick="moreurl(this, {'total': 10, 'clicked': '5383590', 'pos': 8, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s4555049.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/5383590/" onclick="moreurl(this, {'total': 10, 'clicked': '5383590', 'pos': 8, 'identifier': 'book-rec-books'})" class="">
            TIGER X DRAGON SPIN OFF 3 !
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/5988307/" onclick="moreurl(this, {'total': 10, 'clicked': '5988307', 'pos': 9, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s4643200.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/5988307/" onclick="moreurl(this, {'total': 10, 'clicked': '5988307', 'pos': 9, 'identifier': 'book-rec-books'})" class="">
            笨蛋、测验、召唤兽 02
          </a>
        </dd>
      </dl>
        <dl class="clear"></dl>
  </div>
</div>

  






    <link rel="stylesheet" href="https://img3.doubanio.com/f/book/3ec79645ad5a5d15c9ead3c58da97f5d662c7400/css/book/subject/comment.css"/>
    <div class="mod-hd">
        

        <a class="redbutt j a_collect_btn rr" name="cbtn-6438136" href="javascript:;" rel="nofollow">
            <span>我来说两句</span>
        </a>

            
  

  <h2>
    <span class="">短评</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;
      <span class="pl">&nbsp;(
          <a href="https://book.douban.com/subject/6438136/comments/">全部 46 条</a>
        ) </span>

  </h2>


    </div>
    <div class="nav-tab">
        
    <div class="tabs-wrapper  line">
        <a class="short-comment-tabs on-tab" href="hot" data-tab="hot">热门</a>
        <span>/</span>
        <a class="short-comment-tabs " href="new" data-tab="new">最新</a>
        <span>/</span>
        <a class="short-comment-tabs " href="follows" data-tab="follows">好友</a>
    </div>

    </div>
    <div id="comment-list-wrapper" class="indent">
        

<div id="comments" class="comment-list hot show">
        <ul>
                
    <li class="comment-item" data-cid="469767820">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-469767820" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-469767820" class="j vote-comment" data-cid="469767820">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/manill/">marill</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2011-12-15</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">温柔的龙儿啊.</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="415704109">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-415704109" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-415704109" class="j vote-comment" data-cid="415704109">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/zhangyunji/">gerald</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2011-10-26</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">爆笑校园纯爱喜剧真是太有爱了！嗷嗷！</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="638152154">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-638152154" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-638152154" class="j vote-comment" data-cid="638152154">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/yumaobujianle/">年纪轻轻就枫了</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2013-02-02</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">好甜~</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="479812083">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-479812083" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-479812083" class="j vote-comment" data-cid="479812083">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/valencia/">[已注销]</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2012-01-08</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">煮男龙儿对刘海的执念。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="834817126">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-834817126" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-834817126" class="j vote-comment" data-cid="834817126">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/83352328/">现実を嗤う</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2014-08-10</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">心中永远的恋爱喜剧No.1</span>
            </p>
        </div>
    </li>


        </ul>
</div>

        

<div id="comments" class="comment-list new noshow">
        <ul>
                
    <li class="comment-item" data-cid="1510093799">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1510093799" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1510093799" class="j vote-comment" data-cid="1510093799">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/55400900/">吊诡彼岸梦</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2018-11-21</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">补个档，顺便喷一下腾讯动漫，正版轻小说阅读体验非常差，可惜了第十册。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1429485800">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1429485800" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1429485800" class="j vote-comment" data-cid="1429485800">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/60041811/">方丝萦</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2018-08-07</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">动画1、2集。原著小说心理描写更多，更适合做人物分析</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1421278828">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1421278828" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1421278828" class="j vote-comment" data-cid="1421278828">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/132081842/">雾川</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2018-07-30</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">170320
非常妙！所谓的“脱宅神作”恐怕就是如此。
逢坂大河，高须龙儿，北村悠作，节枝实乃梨，川岛亚美，出场的每一个人物，在这短短的十卷里都有了自己的成长。
书中也有许多引人深思的话题。
“每个人都有自己的‘不说’，都各怀心事，都独自为痛苦努力着。”</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1347093072">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1347093072" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1347093072" class="j vote-comment" data-cid="1347093072">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/a3309102/">方糖</a>
                    <span>2018-04-08</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">弃。第一人称和第三人称视角混乱，还有许多莫名奇妙的修饰语，还是看动漫来得痛快。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1338663459">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1338663459" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1338663459" class="j vote-comment" data-cid="1338663459">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/172755573/">反串死一车妈</a>
                        <span class="user-stars allstar30 rating" title="还行"></span>
                    <span>2018-03-10</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">纯爱故事，整体还行。只是很怀念从前叫我「龙儿」的孩子（听起来像神雕侠侣里小龙女？！）。</span>
            </p>
        </div>
    </li>


        </ul>
</div>

            

<div id="comments" class="comment-list follows noshow">
        
        <ul><li class="comment-item"><p class="blank-tip">你关注的人还没有写短评呢</p></li></ul>
</div>

    </div>
        <p>&gt; <a href="https://book.douban.com/subject/6438136/comments/">更多短评 46 条</a></p>
    <script src="https://img3.doubanio.com/f/book/f334de0b97baf7506d0c181ff24dc61f9a7fca64/js/book/subject/short_comment_vote.js"></script>
    <script src="https://img3.doubanio.com/f/book/39eace58cab8aaeec45a44e878bf0ed06f2ed0a4/js/book/subject/short_comment_nav.js"></script>
    <script>
        (function(){
            $('.comment-list').delegate('.vote-comment', 'click', function(e) {
                vote_comment(e);
            }).delegate('.delete-comment', 'click', function(e) {
                if (confirm('确定删除吗？')) {
                    delete_comment(e);
                }
            });
        })();
    </script>

  

<link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/73ed658484f98d44.css">

<section class="topics mod">
    <header>
        <h2>
            TIGER×DRAGON1!的话题 · · · · · ·
            <span class="pl">( <span class="gallery_topics">全部 <span id="topic-count"></span> 条</span> )</span>
        </h2>
    </header>

    




<section class="subject-topics">
    <div class="topic-guide" id="topic-guide">
        <img class="ic_question" src="//img3.doubanio.com/f/ithildin/b1a3edea3d04805f899e9d77c0bfc0d158df10d5/pics/export/icon_question.png">
        <div class="tip_content">
            <div class="tip_title">什么是话题</div>
            <div class="tip_desc">
                <div>无论是一部作品、一个人，还是一件事，都往往可以衍生出许多不同的话题。将这些话题细分出来，分别进行讨论，会有更多收获。</div>
            </div>
        </div>
        <img class="ic_guide" src="//img3.doubanio.com/f/ithildin/529f46d86bc08f55cd0b1843d0492242ebbd22de/pics/export/icon_guide_arrow.png">
        <img class="ic_close" id="topic-guide-close" src="//img3.doubanio.com/f/ithildin/2eb4ad488cb0854644b23f20b6fa312404429589/pics/export/close@3x.png">
    </div>

    <div id="topic-items"></div>

    <script>
        window.subject_id = 6438136;
        window.join_label_text = '写书评参与';

        window.topic_display_count = 4;
        window.topic_item_display_count = 1;
        window.no_content_fun_call_name = "no_topic";

        window.guideNode = document.getElementById('topic-guide');
        window.guideNodeClose = document.getElementById('topic-guide-close');
    </script>
    
        <link rel="stylesheet" href="https://img3.doubanio.com/f/ithildin/f731c9ea474da58c516290b3a6b1dd1237c07c5e/css/export/subject_topics.css">
        <script src="https://img3.doubanio.com/f/ithildin/d3590fc6ac47b33c804037a1aa7eec49075428c8/js/export/moment-with-locales-only-zh.js"></script>
        <script src="https://img3.doubanio.com/f/ithildin/c600fdbe69e3ffa5a3919c81ae8c8b4140e99a3e/js/export/subject_topics.js"></script>

</section>

    <script>
        function no_topic(){
            $('#content .topics').remove();
        }
    </script>
</section>

<section class="reviews mod book-content">
    <header>
        <a href="new_review" rel="nofollow" class="create-review redbutt rr "
            data-isverify="True"
            data-verify-url="https://www.douban.com/accounts/phone/verify?redir=https://book.douban.com/subject/6438136/new_review">
            <span>我要写书评</span>
        </a>
        <h2>
            TIGER×DRAGON1!的书评 · · · · · ·
            <span class="pl">( <a href="reviews">全部 1 条</a> )</span>
        </h2>
    </header>

    

<style>
#gallery-topics-selection {
  position: fixed;
  width: 595px;
  padding: 40px 40px 33px 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 16px 0 rgba(0, 0, 0, 0.2);
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  z-index: 9999;
}
#gallery-topics-selection h1 {
  font-size: 18px;
  color: #007722;
  margin-bottom: 36px;
  padding: 0;
  line-height: 28px;
  font-weight: normal;
}
#gallery-topics-selection .gl_topics {
  border-bottom: 1px solid #dfdfdf;
  max-height: 298px;
  overflow-y: scroll;
}
#gallery-topics-selection .topic {
  margin-bottom: 24px;
}
#gallery-topics-selection .topic_name {
  font-size: 15px;
  color: #333;
  margin: 0;
  line-height: inherit;
}
#gallery-topics-selection .topic_meta {
  font-size: 13px;
  color: #999;
}
#gallery-topics-selection .topics_skip {
  display: block;
  cursor: pointer;
  font-size: 16px;
  color: #3377AA;
  text-align: center;
  margin-top: 33px;
}
#gallery-topics-selection .topics_skip:hover {
  background: transparent;
}
#gallery-topics-selection .close_selection {
  position: absolute;
  width: 30px;
  height: 20px;
  top: 46px;
  right: 40px;
  background: #fff;
  color: #999;
  text-align: right;
}
#gallery-topics-selection .close_selection:hover{
  background: #fff;
  color: #999;
}
</style>




        <div class="review_filter">
            <a href="javascript:;;" class="cur" data-sort="">热门</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="time">最新</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="follow">好友</a href="javascript:;;">
            
        </div>


        



<div class="review-list  ">
        
    

        
    
    <div data-cid="7656160">
        <div class="main review-item" id="7656160">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/disangel/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u31851074-4.jpg">
        </a>

        <a href="https://www.douban.com/people/disangel/" class="name">老傅</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2015-11-12" class="main-meta">2015-11-12 08:52:38</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/7656160/">二次元最真实的恋爱剧</a></h2>

                <div id="review_7656160_short" class="review-short" data-rid="7656160">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇书评可能有关键情节透露</p>

                        我不喜欢看恋爱剧，不论次元。三次元的恋爱剧过于现实，金钱家庭生存等等等等过于现实的问题将那份最纯真的感情冲得平谈无味。我们在现实生活中已经受够了诸如此类的折磨，到了放松的时刻为何还要去面对这些？而二次元的恋爱剧又过于YY，无论男主角有着一副怎样的穷形尽相哪怕...

                        &nbsp;(<a href="javascript:;" id="toggle-7656160-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_7656160_full" class="hidden">
                    <div id="review_7656160_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="7656160" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-7656160">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="7656160" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-7656160">
                                1
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/7656160/#comments" class="reply ">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>




    

    

    <script type="text/javascript">
    ;(function(){
        window.subject_id = window.subject_id || 6438136;
        "use strict";function _classCallCheck(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}function _possibleConstructorReturn(e,t){if(!e)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!t||"object"!=typeof t&&"function"!=typeof t?e:t}function _inherits(e,t){if("function"!=typeof t&&null!==t)throw new TypeError("Super expression must either be null or a function, not "+typeof t);e.prototype=Object.create(t&&t.prototype,{constructor:{value:e,enumerable:!1,writable:!0,configurable:!0}}),t&&(Object.setPrototypeOf?Object.setPrototypeOf(e,t):e.__proto__=t)}var _createClass=function(){function e(e,t){for(var n=0;n<t.length;n++){var c=t[n];c.enumerable=c.enumerable||!1,c.configurable=!0,"value"in c&&(c.writable=!0),Object.defineProperty(e,c.key,c)}}return function(t,n,c){return n&&e(t.prototype,n),c&&e(t,c),t}}(),GalleryTopicsSelection=function(e){function t(){_classCallCheck(this,t);var e=_possibleConstructorReturn(this,(t.__proto__||Object.getPrototypeOf(t)).call(this));return e.state={topics:[]},e}return _inherits(t,e),_createClass(t,[{key:"componentDidMount",value:function(){var e=this;$.ajax({url:"https://m.douban.com/rexxar/api/v2/gallery/subject_feed?start=0&subject_id="+window.subject_id+"&ck="+get_cookie("ck"),xhrFields:{withCredentials:!0},success:function(t){e.setState({topics:t.items}),t.total&&(window.has_gallery_topics=!0)}})}},{key:"render",value:function(){var e=this.state.topics;return React.createElement("section",null,React.createElement("a",{href:"new_review?from=gallery-topics-selection&click=close",rel:"nofollow",className:"close_selection"},"跳过"),React.createElement("h1",null,"下面是否有你想写的话题？"),e.length?React.createElement("ul",{className:"gl_topics"},e.map(function(e){return React.createElement("li",{className:"topic"},React.createElement("a",{href:"new_review?from=gallery-topics-selection&click=create&topic_id="+e.topic.id,className:"comment_btn write_review",rel:"nofollow"},React.createElement("img",{src:window.write_icon}),React.createElement("span",null,window.join_label_text)),React.createElement("h2",{className:"topic_name"},e.topic.name),React.createElement("div",{className:"topic_meta"},e.topic.card_subtitle))})):React.createElement("div",null,"加载中"),React.createElement("a",{href:"new_review?from=gallery-topics-selection&click=skip",rel:"nofollow",className:"topics_skip"},React.createElement("span",null,"上面没有我想写的话题，去写影评 ＞ ")))}}]),t}(React.Component);$(function(){window.renderGalleryTopicsSelection=function(){var e=document.getElementById("gallery-topics-selection");ReactDOM.render(React.createElement(GalleryTopicsSelection,null),e)}});;
    })();
</script><script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/7c12b40a1a0c867a.js"></script>
    <!-- COLLECTED CSS -->
</div>








            <p class="pl">
                &gt;
                <a href="reviews">
                    更多书评1篇
                </a>
            </p>
</section>

<!-- COLLECTED JS -->

  









<div class="ugc-mod reading-notes">
  <div class="hd">
  </div>
  


</div>



<script type="text/javascript">
  $(document).ready(function(){
    var TEMPL_ADD_COL = '<a href="" id="n-{NOTE_ID}" class="colbutt ll add-col"><span>收藏</span></a>',
      TEMPL_DEL_COL = '<span class="pl">已收藏 &gt;<a href="" id="n-{NOTE_ID}" class="del-col">取消收藏</a></span>';

    $('body').delegate('.add-col', 'click', function(e){
      e.preventDefault();
      var nnid = $(this).attr('id').split('-')[1],
        bn_add = $(this);
      $.post_withck("/j/book/annotation_collect",{nid:nnid},function(){
        var a = $(TEMPL_DEL_COL.replace('{NOTE_ID}', nnid));
        bn_add.before(a);
        bn_add.remove();
      })
    });

    $('body').delegate('.del-col', 'click', function(e){
      e.preventDefault();
      var nnid = $(this).attr('id').split('-')[1],
        bn_del = $(this).parent();
      $.post_withck("/j/book/annotation_uncollect", {nid: nnid}, function() {
        var a = $(TEMPL_ADD_COL.replace('{NOTE_ID}', nnid));
        bn_del.before(a);
        bn_del.remove();
      })
    });

    $("pre.source").each(function(){
      var cn = $(this).attr('class').split(' ');
      l = cn[1];
      s = 'rand01';
      n = cn[2];
      $(this).snippet(n,{style: s, showNum: l});
    });

    var annotationMod = $('.reading-notes .bd')
      , annotationTabs = annotationMod.find('.inline-tabs li')
      , annotationTabLinks = annotationTabs.find('a')
      , annotationTabContents = annotationMod.find('ul.comments');

    annotationTabLinks.click(function(e){
      e.preventDefault();
      var el = $(this)
        , kind = el.attr('id');

      annotationTabs.removeClass('on');
      el.parent().addClass('on');
      annotationTabContents.hide();
      annotationTabContents.filter('.' + kind).show();
    });
  });
</script>

<script type="text/x-mathjax-config">
MathJax.Hub.Config({
	jax: ["input/TeX", "output/HTML-CSS"],
    extensions: ["tex2jax.js","TeX/AMSmath.js","TeX/AMSsymbols.js","TeX/noUndefined.js"],
    tex2jax: {
		inlineMath: [ ["($", "$)"], ['\\(','\\)'] ],
		displayMath: [ ["($$","$$)"], ['\\[','\\]']],
		skipTags: ["script","noscript","style","textarea"],
		processEscapes: true,
		processEnvironments: true,
		preview: "TeX"
	},
	showProcessingMessages: false
  });
</script>


  






<div id="db-discussion-section" class="indent ugc-mod">




        
    <a class="redbutt  rr" href="https://book.douban.com/subject/6438136/discussion/create" rel="nofollow"><span>在这本书的论坛里发言</span></a>



</div>




</div>
<script type="text/javascript">
$(function(){if($.browser.msie && $.browser.version == 6.0){
    var maxWidth = parseInt($('#info').css('max-width'));
    if($('#info').width() > maxWidth)
        $('#info').width(maxWidth)
}});
</script>
</div>
      <div class="aside">
        
  
  






  <div id="dale_book_subject_top_right"></div>

  






<style type="text/css" media="screen">
  .add2cartContainer{overflow:hidden;vertical-align:bottom;line-height:1}.add2cartContainer .add2cart{margin-right:0;display:inline-block}#buyinfo .bs{margin:0}#buyinfo li{padding:10px 0;position:relative;line-height:1;border-bottom:1px solid #eaeaea}#buyinfo li a:hover{background-image:none !important}#buyinfo li a:hover .buylink-price{background:#37a}#buyinfo li .publish,#buyinfo li .other-activity{margin-top:5px}#buyinfo li .publish a,#buyinfo li .other-activity a{color:#999}#buyinfo li .publish a:hover,#buyinfo li .other-activity a:hover{color:#37a;background:none;opacity:0.5;filter:alpha(opacity=50)}#buyinfo li .buylink-price{position:absolute;right:90px;text-align:right}#buyinfo .more-info{color:#aaa;margin:6px 0 -2px 0}#buyinfo .more-ebooks{padding:10px 0;color:#37a;cursor:pointer}#buyinfo .price-page{border-bottom:0;padding:15px 0 0}#buyinfo .saved-price{display:none;margin-left:5px}#buyinfo .cart-tip{float:right;padding-right:5px}#buyinfo #buyinfo-ebook{margin-bottom:15px}#buyinfo #buyinfo-ebook .buylink-price{display:inline}#buyinfo #buyinfo-ebook li.no-border{border:0}#buyinfo-printed{margin-bottom:15px}#buyinfo-printed.no-border{border-bottom:0}#buyinfo-printed .more-ebooks{line-height:1;padding:10px 0;color:#37a;cursor:pointer;padding:10px 0 0}#buyinfo-report{display:none}#buyinfo-report .lnk-close-report{float:right;margin-top:-30px;line-height:14px}#buyinfo-report .item{margin-bottom:10px}#buyinfo-report .item input{vertical-align:text-bottom;*vertical-align:middle}#buyinfo-report .item label{margin:0 5px 0 2px}#buyinfo-report .item-submit .bn-flat{margin-right:10px}#buyinfo-report .item-price input{width:220px;border:1px solid #ccc;padding:4px}#buyinfo-report form{margin:5px 0 10px}#bi-report-btn{float:right;margin:2px 0;line-height:14px;display:none}.bi-vendor-report{color:#aaa}.bi-vendor-report-form{display:none;color:#111;margin:0 5px;line-height:25px}.gray_ad{padding:30px 20px 25px 20px;background:#f6f6f1}.gray_ad h2{margin-bottom:6px;font-size:15px}.gray_ad .ebook-tag{margin-top:5px;color:#999;font-size:12px}.bs.more-after{margin-bottom:0px}@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi){#buyinfo li a:hover{background-image:url(https://img3.doubanio.com/f/book/fc4ff7f0a3a7f452f06d586540284b9738f2fe87/pics/book/cart/icon-brown@2x.png);background-size:16px 12px}}#intervenor-buyinfo .bs{margin:0}#intervenor-buyinfo li{position:relative;border-bottom:1px solid #eaeaea;padding:10px 0;line-height:1}#intervenor-buyinfo li .basic-info{color:#494949;font-size:14px;line-height:18px}#intervenor-buyinfo li a:hover .comment{color:#f67;opacity:0.75;filter:alpha(opacity=75)}#intervenor-buyinfo li a:hover .buy-btn{background:#fff;border:1px solid #e97e7e;border-radius:2px;color:#e97e7e}#intervenor-buyinfo li a:hover .buylink-price{background:#37a}#intervenor-buyinfo li .buylink-price{position:absolute;right:90px;text-align:right}#intervenor-buyinfo li .publish,#intervenor-buyinfo li .other-activity{margin-top:5px}#intervenor-buyinfo li .publish a,#intervenor-buyinfo li .other-activity a{color:#999}#intervenor-buyinfo li .publish a:hover,#intervenor-buyinfo li .other-activity a:hover{color:#37a;background:none;opacity:0.5;filter:alpha(opacity=50)}#intervenor-buyinfo .jd-buy-icon{float:left;margin-right:3px}#intervenor-buyinfo .buy-btn{float:right;position:absolute;right:10px;bottom:3px;color:#9c9c9c;padding:0 12px;border:1px solid transparent}#intervenor-buyinfo .comment{color:#FF8C9C;margin:6px 0 -2px 0}#intervenor-buyinfo .price-page a{display:inline-block;line-height:16px !important}#intervenor-buyinfo .price-page{border-bottom:0;padding:15px 0 0}#intervenor-buyinfo .saved-price{display:none;margin-left:5px}#intervenor-buyinfo .cart-tip{float:right;padding-right:5px}#intervenor-buyinfo #buyinfo-ebook{margin-bottom:15px}#intervenor-buyinfo #buyinfo-ebook .buylink-price{display:inline}#intervenor-buyinfo #buyinfo-ebook li.no-border{border:0}#buyinfo-printed .presale-indicator{margin:0;width:auto;color:#999;text-indent:0;background:none}

</style>

      <div class="gray_ad" id="buyinfo">
      <div id="buyinfo-printed">
        
  

  <h2>
    <span class="">在哪儿买这本书</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


        <ul class="bs noline more-after ">
          
                
                <li class="">
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=1300&amp;pre=0&amp;vendor=jingdong&amp;srcpage=subject&amp;price=1380&amp;pos=1&amp;url=http%3A%2F%2Funion-click.jd.com%2Fjdc%3Fe%3D%26p%3DAyIHZR5aEQISA1AYUyUCEw5VHFIVAiJDCkMFSjJLQhBaUAscSkIBR0ROVw1VC0dFFQMbB1ISWxUdS0IJRmsdVmZ6FBghfGdlQwlNXF1bVGxdchJTDh5pVhpbEx4SBEkaWAQCGgxUG1AUByIAUBxdEgsRAGUbXxQGFjdlG1olUHzf462DsLMO0%252F%252BUjp2VIgZlG18VAREHVh5cEAIQA2Uca8yCpdH3iYKLiMePwitrJTIiN1QrWg%253D%253D%26t%3DW1dCFBBFC1pXUwkEAEAdQFkJBVsUCxIAXBtbCltXWwg%253D&amp;cntvendor=2&amp;srcsubj=6438136&amp;type=bkbuy&amp;subject=6438136" class="">
                      <span class="">京东商城</span>
                    </a>
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=1300&amp;pre=0&amp;vendor=jingdong&amp;srcpage=subject&amp;price=1380&amp;pos=1&amp;url=http%3A%2F%2Funion-click.jd.com%2Fjdc%3Fe%3D%26p%3DAyIHZR5aEQISA1AYUyUCEw5VHFIVAiJDCkMFSjJLQhBaUAscSkIBR0ROVw1VC0dFFQMbB1ISWxUdS0IJRmsdVmZ6FBghfGdlQwlNXF1bVGxdchJTDh5pVhpbEx4SBEkaWAQCGgxUG1AUByIAUBxdEgsRAGUbXxQGFjdlG1olUHzf462DsLMO0%252F%252BUjp2VIgZlG18VAREHVh5cEAIQA2Uca8yCpdH3iYKLiMePwitrJTIiN1QrWg%253D%253D%26t%3DW1dCFBBFC1pXUwkEAEAdQFkJBVsUCxIAXBtbCltXWwg%253D&amp;cntvendor=2&amp;srcsubj=6438136&amp;type=bkbuy&amp;subject=6438136" class="buylink-price ">
                      <span class="">
                        13.80 元
                      </span>
                    </a>

                      
                </li>
                
                <li class="">
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=1300&amp;pre=0&amp;vendor=bookschina&amp;srcpage=subject&amp;price=1300&amp;pos=2&amp;url=http%3A%2F%2Fwww.bookschina.com%2Funion%2Fubook.asp%3Fadservice%3D354872%26tourl%3Dhttp%3A%2F%2Fwww.bookschina.com%2F5237361.htm&amp;cntvendor=2&amp;link2key=b589156548&amp;srcsubj=6438136&amp;type=bkbuy&amp;subject=6438136" class="">
                      <span class="">中国图书网</span>
                    </a>
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=1300&amp;pre=0&amp;vendor=bookschina&amp;srcpage=subject&amp;price=1300&amp;pos=2&amp;url=http%3A%2F%2Fwww.bookschina.com%2Funion%2Fubook.asp%3Fadservice%3D354872%26tourl%3Dhttp%3A%2F%2Fwww.bookschina.com%2F5237361.htm&amp;cntvendor=2&amp;link2key=b589156548&amp;srcsubj=6438136&amp;type=bkbuy&amp;subject=6438136" class="buylink-price ">
                      <span class="">
                        13.00 元
                      </span>
                    </a>

                      
                </li>
          <li class="price-page">
            <a href="https://book.douban.com/subject/6438136/buylinks">
              &gt; 查看2家网店价格
                (13.00 元起)
            </a>
          </li>
        </ul>
      </div>
      
  <div class="add2cartContainer ft no-border">
    
  <span class="add2cartWidget ll">
        
        <a href="javascript:;" class="j  a_add2cart add2cart"
          name="6438136">
          <span>+ 加入购书单</span></a>
        <span class="color_gary book-in-cart hidden"
          >
            已在购书单&nbsp;
            <a href="https://book.douban.com/cart">查看</a>
          <a class="delete-cart-item" rel="6438136" href="https://book.douban.com/cart">删除</a>
        </span>
  </span>
  

  </div>

  </div>
  <script type="text/javascript">
  $('.more-ebooks').on('click', function() {
    var $this = $(this),
      $li = $this.siblings('ul').find('li');
    if ($this.hasClass('isShow')) {
      $(this).text('展开更多').removeClass('isShow');
      $li.not(':first').addClass('hide');
    }else{
      $(this).text('收起').addClass('isShow');
      $li.removeClass('hide');
    }
    
  })
  </script>

<style class="text/css">
  .presale-indicator {
  display: inline-block;
  *display: inline;
  *zoom: 1;
  width: 24px;
  height: 15px;
  line-height: 15px;
  background: url(https://img3.doubanio.com/f/book/1679c65572eac1371f9872807199dea6e55a7f06/pics/book/cart/presale_text.gif) no-repeat;
  text-indent: -9999px;
  vertical-align: middle;
  *vertical-align: 0px;
  _vertical-align: 2px;
  margin-left: 0.5em;
}

</style>



  





<div class="gray_ad" id="borrowinfo">
  
  

  <h2>
    <span class="">在哪儿借这本书</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


  <ul class="bs more-after">
      
      <li style="border: none">
        <a href="https://www.douban.com/link2/?url=http%3A%2F%2Fipac.library.sh.cn%2Fipac20%2Fipac.jsp%3Faspect%3Dbasic_search%26profile%3Dsl%26index%3DISBN%26term%3D9787535644954&amp;subject=7535644953&amp;type=borrow&amp;library=10012&amp;link2key=b589156548" target="_blank">上海市中心图书馆(8)</a>
      </li>
      
      <li style="border: none">
        <a href="https://www.douban.com/link2/?url=http%3A%2F%2Fopac.gzlib.gov.cn%2Fopac%2Fsearch%3Frows%3D10%26curlibcode%3DGT%26hasholding%3D1%26searchWay0%3Dmarc%26q0%3D%26logical0%3DAND%26q%3D9787535644954%26searchWay%3Disbn%26scWay%3Ddim%26searchSource%3Dreader&amp;subject=7535644953&amp;type=borrow&amp;library=10022&amp;link2key=b589156548" target="_blank">广州图书馆(1)</a>
      </li>
      
      <li style="border: none">
        <a href="https://www.douban.com/link2/?url=http%3A%2F%2Fopac3.wzlib.cn%2Fopac%2Fsearch%3FsearchWay%3Disbn%26q%3D978-7-5356-4495-4%26booktype%3D%26marcformat%3D%26sortWay%3Dscore%26sortOrder%3Ddesc%26startPubdate%3D%26endPubdate%3D%26rows%3D10&amp;subject=7535644953&amp;type=borrow&amp;library=10020&amp;link2key=b589156548" target="_blank">温州市图书馆</a>
      </li>
  </ul>
  <div class="clearfix"></div>
  <!--<div class="ft pl">-->
    <!--<a class="rr"  href="https://book.douban.com/library_invitation">&gt; 图书馆合作</a>-->
    <!--找不到你需要的图书馆？-->
  <!--</div>-->
</div>

  <div id="dale_book_subject_top_middle"></div>
  





  

  
  

  <h2>
    <span class="">这本书的其他版本 </span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;
      <span class="pl">&nbsp;(
          <a href="https://book.douban.com/works/1011799">全部4</a>
        ) </span>

  </h2>


  <div class="indent">
    <ul>
        <li class="mb8 pl">
          <a href="https://book.douban.com/subject/30426492/">未知出版社版</a>
          2018-5-8 / 有售
        </li>
        <li class="mb8 pl">
          <a href="https://book.douban.com/subject/3257906/">台湾角川版</a>
          2006-12-26 / 194人读过
        </li>
        <li class="mb8 pl">
          <a href="https://book.douban.com/subject/2244504/">メディアワークス版</a>
          2006-03-25 / 58人读过 / 有售
        </li>
    </ul>
  </div>


  



      
  

  <h2>
    <span class="">以下豆列推荐</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;
      <span class="pl">&nbsp;(
          <a href="https://book.douban.com/subject/6438136/doulists">全部</a>
        ) </span>

  </h2>


    <div id="db-doulist-section" class="indent">
      <ul class="bs">
          <li class=""><a class="" href="https://www.douban.com/doulist/725624/" target="_blank">天闻角川图书</a>
                <span class="pl">(demongod)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/725726/" target="_blank">简体版轻小说首卷索引</a>
                <span class="pl">(田乐馒)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/821266/" target="_blank">2011年这本轻小说真厉害！ BEST 60</a>
                <span class="pl">(田乐馒)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/1152847/" target="_blank">豆瓣评星很高很强大【轻小说/网络/悬疑】</a>
                <span class="pl">(爱玛·包法利)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/3464897/" target="_blank">最近读的不错的轻小说。</a>
                <span class="pl">(Megnha)</span>
            </li>
      </ul>
    </div>

  <div id="dale_book_subject_middle_mini"></div>
  






  <h2>谁读这本书?</h2>
  <div class="indent" id="collector">

    

<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/84971956/"><img src="https://img3.doubanio.com/icon/u84971956-2.jpg" class="pil" alt="HelloKent" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/84971956/">HelloKent</a><br/>
      <div class="pl ll">          昨天          读过      </div>

        <span class="allstar50" title="力荐"></span>
      <br/>

    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>
<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/191454368/"><img src="https://img3.doubanio.com/icon/u191454368-1.jpg" class="pil" alt="__梦幻泡影__" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/191454368/">__梦幻泡影__</a><br/>
      <div class="pl ll">          2月11日          读过      </div>

        <span class="allstar50" title="力荐"></span>
      <br/>

    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>
<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/140508991/"><img src="https://img3.doubanio.com/icon/u140508991-1.jpg" class="pil" alt="华亭" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/140508991/">华亭</a><br/>
      <div class="pl ll">          2月8日          读过      </div>

        <span class="allstar40" title="推荐"></span>
      <br/>

    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>
<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/57755053/"><img src="https://img3.doubanio.com/icon/u57755053-4.jpg" class="pil" alt="Kagami" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/57755053/">Kagami</a><br/>
      <div class="pl ll">          1月6日          读过      </div>

        <span class="allstar50" title="力荐"></span>
      <br/>

      <span class="pl">tags:轻小说 日本</span><br />
    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>


        <p class="pl">&gt; <a href="https://book.douban.com/subject/6438136/doings">12人在读</a></p>
        <p class="pl">&gt; <a href="https://book.douban.com/subject/6438136/collections">212人读过</a></p>
        <p class="pl">&gt; <a href="https://book.douban.com/subject/6438136/wishes">66人想读</a></p>

  </div>





  
<!-- douban ad begin -->
<div id="dale_book_subject_middle_right"></div>
<script type="text/javascript">
    (function (global) {
        if(!document.getElementsByClassName) {
            document.getElementsByClassName = function(className) {
                return this.querySelectorAll("." + className);
            };
            Element.prototype.getElementsByClassName = document.getElementsByClassName;

        }
        var articles = global.document.getElementsByClassName('article'),
            asides = global.document.getElementsByClassName('aside');

        if (articles.length > 0 && asides.length > 0 && articles[0].offsetHeight >= asides[0].offsetHeight) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_book_subject_middle_right');
        }
    })(this);
</script>
<!-- douban ad end -->

  





  

  <h2>二手市场</h2>
  <div class="indent">
    <ul class="bs">
    <li class="">
          <a href='https://book.douban.com/subject/6438136/new_offer' class="rr">&gt; 点这儿转让</a>

        有66人想读,手里有一本闲着?
      </li>
    </ul>
  </div>

  
<p class="pl">订阅关于TIGER×DRAGON1!的评论: <br/><span class="feed">
    <a href="https://book.douban.com/feed/subject/6438136/reviews"> feed: rss 2.0</a></span></p>


      </div>
      <div class="extra">
        
  
<!-- douban ad begin -->
<div id="dale_book_subject_bottom_super_banner"></div>
<script type="text/javascript">
    (function (global) {
        var body = global.document.body,
            html = global.document.documentElement;

        var height = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight, html.scrollHeight, html.offsetHeight);
        if (height >= 2000) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_book_subject_bottom_super_banner');
        }
    })(this);
</script>
<!-- douban ad end -->


      </div>
    </div>
  </div>

        
<div id="footer">

<span id="icp" class="fleft gray-link">
    &copy; 2005－2019 douban.com, all rights reserved 北京豆网科技有限公司
</span>

<a href="https://www.douban.com/hnypt/variformcyst.py" style="display: none;"></a>

<span class="fright">
    <a href="https://www.douban.com/about">关于豆瓣</a>
    · <a href="https://www.douban.com/jobs">在豆瓣工作</a>
    · <a href="https://www.douban.com/about?topic=contactus">联系我们</a>
    · <a href="https://www.douban.com/about?policy=disclaimer">免责声明</a>
    
    · <a href="https://help.douban.com/?app=book" target="_blank">帮助中心</a>
    · <a href="https://book.douban.com/library_invitation">图书馆合作</a>
    · <a href="https://www.douban.com/doubanapp/">移动应用</a>
    · <a href="https://www.douban.com/partner/">豆瓣广告</a>
</span>

</div>

    </div>
      
  

    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/3b736fcab9961d20.js"></script>
    <!-- mako -->
    
  








    
<script type="text/javascript">
    (function (global) {
        var newNode = global.document.createElement('script'),
            existingNode = global.document.getElementsByTagName('script')[0],
            adSource = '//erebor.douban.com/',
            userId = '54106750',
            browserId = 'Fr2Y9YzxI8s',
            criteria = '7:轻小说|7:龙与虎|7:竹宫悠由子|7:日本|7:天闻角川|7:青春|7:纯爱|7:小说|7:大河俺の嫁！|7:輕小說|3:/subject/6438136/',
            preview = '',
            debug = false,
            adSlots = ['dale_book_subject_top_icon', 'dale_book_subject_top_right', 'dale_book_subject_top_middle', 'dale_book_subject_middle_mini'];

        global.DoubanAdRequest = {src: adSource, uid: userId, bid: browserId, crtr: criteria, prv: preview, debug: debug};
        global.DoubanAdSlots = (global.DoubanAdSlots || []).concat(adSlots);

        newNode.setAttribute('type', 'text/javascript');
        newNode.setAttribute('src', 'https://img3.doubanio.com/f/adjs/dcee110fbad16bf4ea756300ebd9551cead60e96/ad.release.js');
        newNode.setAttribute('async', true);
        existingNode.parentNode.insertBefore(newNode, existingNode);
    })(this);
</script>












    
  

<script type="text/javascript">
  var _paq = _paq || [];
  _paq.push(['trackPageView']);
  _paq.push(['enableLinkTracking']);
  (function() {
    var p=(('https:' == document.location.protocol) ? 'https' : 'http'), u=p+'://fundin.douban.com/';
    _paq.push(['setTrackerUrl', u+'piwik']);
    _paq.push(['setSiteId', '100001']);
    var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0]; 
    g.type='text/javascript';
    g.defer=true; 
    g.async=true; 
    g.src=p+'://s.doubanio.com/dae/fundin/piwik.js';
    s.parentNode.insertBefore(g,s);
  })();
</script>

<script type="text/javascript">
var setMethodWithNs = function(namespace) {
  var ns = namespace ? namespace + '.' : ''
    , fn = function(string) {
        if(!ns) {return string}
        return ns + string
      }
  return fn
}

var gaWithNamespace = function(fn, namespace) {
  var method = setMethodWithNs(namespace)
  fn.call(this, method)
}

var _gaq = _gaq || []
  , accounts = [
      { id: 'UA-7019765-1', namespace: 'douban' }
    , { id: 'UA-7019765-16', namespace: '' }
    ]
  , gaInit = function(account) {
      gaWithNamespace(function(method) {
        gaInitFn.call(this, method, account)
      }, account.namespace)
    }
  , gaInitFn = function(method, account) {
      _gaq.push([method('_setAccount'), account.id])

      
  _gaq.push([method('_addOrganic'), 'google', 'q'])
  _gaq.push([method('_addOrganic'), 'baidu', 'wd'])
  _gaq.push([method('_addOrganic'), 'soso', 'w'])
  _gaq.push([method('_addOrganic'), 'youdao', 'q'])
  _gaq.push([method('_addOrganic'), 'so.360.cn', 'q'])
  _gaq.push([method('_addOrganic'), 'sogou', 'query'])
  if (account.namespace) {
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣'])
    _gaq.push([method('_addIgnoredOrganic'), 'douban'])
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣网'])
    _gaq.push([method('_addIgnoredOrganic'), 'www.douban.com'])
  }

      if (account.namespace === 'douban') {
        _gaq.push([method('_setDomainName'), '.douban.com'])
      }

        _gaq.push([method('_setCustomVar'), 1, 'responsive_view_mode', 'desktop', 3])

        _gaq.push([method('_setCustomVar'), 2, 'login_status', '1', 2]);

      _gaq.push([method('_trackPageview')])
    }

for(var i = 0, l = accounts.length; i < l; i++) {
  var account = accounts[i]
  gaInit(account)
}


;(function() {
    var ga = document.createElement('script');
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    ga.setAttribute('async', 'true');
    document.documentElement.firstChild.appendChild(ga);
})()
</script>








    <!-- anson42-docker-->

</body>
</html>




































`
var html66 = `

<!DOCTYPE html>
<html lang="zh-cmn-Hans" class="ua-windows ua-webkit">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="renderer" content="webkit">
    <meta name="referrer" content="always">
    <meta name="google-site-verification" content="ok0wCgT20tBBgo9_zat2iAcimtN4Ftf5ccsh092Xeyw" />
    <title>
        三少爷的剑 (豆瓣)
</title>
    
    <meta name="baidu-site-verification" content="cZdR4xxR7RxmM4zE" />
    <meta http-equiv="Pragma" content="no-cache">
    <meta http-equiv="Expires" content="Sun, 6 Mar 2005 01:00:00 GMT">
    
    <link rel="apple-touch-icon" href="https://img3.doubanio.com/f/movie/d59b2715fdea4968a450ee5f6c95c7d7a2030065/pics/movie/apple-touch-icon.png">
    <link href="https://img3.doubanio.com/f/shire/bf61b1fa02f564a4a8f809da7c7179b883a56146/css/douban.css" rel="stylesheet" type="text/css">
    <link href="https://img3.doubanio.com/f/shire/ae3f5a3e3085968370b1fc63afcecb22d3284848/css/separation/_all.css" rel="stylesheet" type="text/css">
    <link href="https://img3.doubanio.com/f/movie/8864d3756094f5272d3c93e30ee2e324665855b0/css/movie/base/init.css" rel="stylesheet">
    <script type="text/javascript">var _head_start = new Date();</script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/movie/0495cb173e298c28593766009c7b0a953246c5b5/js/movie/lib/jquery.js"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/f010949d3f23dd7c972ad7cb40b800bf70723c93/js/douban.js"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/0efdc63b77f895eaf85281fb0e44d435c6239a3f/js/separation/_all.js"></script>
    
    <meta name="keywords" content="三少爷的剑,三少爷的剑,三少爷的剑影评,剧情介绍,图片,论坛">
    <meta name="description" content="三少爷的剑电视剧简介和剧情介绍,三少爷的剑影评、图片、论坛">
    <meta name="mobile-agent" content="format=html5; url=http://m.douban.com/movie/subject/2279835/"/>
    <link rel="alternate" href="android-app://com.douban.movie/doubanmovie/subject/2279835/" />
    <link rel="stylesheet" href="https://img3.doubanio.com/dae/cdnlib/libs/LikeButton/1.0.5/style.min.css">
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/77323ae72a612bba8b65f845491513ff3329b1bb/js/do.js" data-cfg-autoload="false"></script>
    <script type="text/javascript">
      Do.add('dialog', {path: 'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js', type: 'js'});
      Do.add('dialog-css', {path: 'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css', type: 'css'});
      Do.add('handlebarsjs', {path: 'https://img3.doubanio.com/f/movie/3d4f8e4a8918718256450eb6e57ec8e1f7a2e14b/js/movie/lib/handlebars.current.js', type: 'js'});
    </script>
    
  <script type='text/javascript'>
  var _vwo_code = (function() {
    var account_id = 249272,
      settings_tolerance = 0,
      library_tolerance = 2500,
      use_existing_jquery = false,
      // DO NOT EDIT BELOW THIS LINE
      f=false,d=document;return{use_existing_jquery:function(){return use_existing_jquery;},library_tolerance:function(){return library_tolerance;},finish:function(){if(!f){f=true;var a=d.getElementById('_vis_opt_path_hides');if(a)a.parentNode.removeChild(a);}},finished:function(){return f;},load:function(a){var b=d.createElement('script');b.src=a;b.type='text/javascript';b.innerText;b.onerror=function(){_vwo_code.finish();};d.getElementsByTagName('head')[0].appendChild(b);},init:function(){settings_timer=setTimeout('_vwo_code.finish()',settings_tolerance);var a=d.createElement('style'),b='body{opacity:0 !important;filter:alpha(opacity=0) !important;background:none !important;}',h=d.getElementsByTagName('head')[0];a.setAttribute('id','_vis_opt_path_hides');a.setAttribute('type','text/css');if(a.styleSheet)a.styleSheet.cssText=b;else a.appendChild(d.createTextNode(b));h.appendChild(a);this.load('//dev.visualwebsiteoptimizer.com/j.php?a='+account_id+'&u='+encodeURIComponent(d.URL)+'&r='+Math.random());return settings_timer;}};}());

  +function () {
    var bindEvent = function (el, type, handler) {
        var $ = window.jQuery || window.Zepto || window.$
       if ($ && $.fn && $.fn.on) {
           $(el).on(type, handler)
       } else if($ && $.fn && $.fn.bind) {
           $(el).bind(type, handler)
       } else if (el.addEventListener){
         el.addEventListener(type, handler, false);
       } else if (el.attachEvent){
         el.attachEvent("on" + type, handler);
       } else {
         el["on" + type] = handler;
       }
     }

    var _origin_load = _vwo_code.load
    _vwo_code.load = function () {
      var args = [].slice.call(arguments)
      bindEvent(window, 'load', function () {
        _origin_load.apply(_vwo_code, args)
      })
    }
  }()

  _vwo_settings_timer = _vwo_code.init();
  </script>


    


<script type="application/ld+json">
{
  "@context": "http://schema.org",
  "name": "三少爷的剑",
  "url": "/subject/2279835/",
  "image": "https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2256387990.webp",
  "director": 
  [
    {
      "@type": "Person",
      "url": "/celebrity/1320042/",
      "name": "靳德茂 De-mao Jin"
    }
  ]
,
  "author": 
  [
    {
      "@type": "Person",
      "url": "/celebrity/1315799/",
      "name": "古龙 Lung Ku"
    }
  ]
,
  "actor": 
  [
    {
      "@type": "Person",
      "url": "/celebrity/1314497/",
      "name": "何中华 Zhonghua He "
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1014618/",
      "name": "俞飞鸿 Faye Yu"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1314478/",
      "name": "陈龙 Long Chen"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1318452/",
      "name": "陈继铭 Jiming Chen"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1323125/",
      "name": "张伊函 Yihan Zhang"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1000145/",
      "name": "霍思燕 Siyan Huo"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1276110/",
      "name": "刘莉莉 Lili Liu"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1274488/",
      "name": "杨若兮 Ruoxi Yang"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1274566/",
      "name": "石小满 Xiaoman Shi"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1313611/",
      "name": "岳跃利 Yueli Yue"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1274559/",
      "name": "戴春荣 Chunrong Dai"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1274564/",
      "name": "陈莹 Ying Chen"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1016668/",
      "name": "张静初 Jingchu Zhang"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1314844/",
      "name": "赵毅 Yi Zhao"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1322790/",
      "name": "刘大刚 Dagang Liu"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1318619/",
      "name": "刘莉莉 Lili Liu "
    }
  ]
,
  "datePublished": "",
  "genre": ["\u6b66\u4fa0", "\u53e4\u88c5"],
  "duration": "PT0H47M",
  "description": "江湖纷争，麻烦不断。神剑山庄的三少爷谢晓峰（何中华 饰）和慕容世家的独女慕容秋荻（俞飞鸿 饰）在百般波折之下，终于迎来了二人要大喜的日子。可大婚当天，谢晓峰却不得不与夺命十三剑——燕十三（王冰 饰）进...",
  "@type": "TVSeries",
  "aggregateRating": {
    "@type": "AggregateRating",
    "ratingCount": "5475",
    "bestRating": "10",
    "worstRating": "2",
    "ratingValue": "6.9"
  }
}
</script>


    <style type="text/css">
  
</style>
    <style type="text/css">img { max-width: 100%; }</style>
    <script type="text/javascript"></script>
    <link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/703752ab692b1cba.css">

    <link rel="shortcut icon" href="https://img3.doubanio.com/favicon.ico" type="image/x-icon">
</head>

<body>
  
    <script type="text/javascript">var _body_start = new Date();</script>

    
    



    <link href="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.css" rel="stylesheet" type="text/css">



<div id="db-global-nav" class="global-nav">
  <div class="bd">
    
<div class="top-nav-info">
  <a href="https://www.douban.com/accounts/login?source=movie" class="nav-login" rel="nofollow">登录</a>
  <a href="https://www.douban.com/accounts/register?source=movie" class="nav-register" rel="nofollow">注册</a>
</div>


    <div class="top-nav-doubanapp">
  <a href="https://www.douban.com/doubanapp/app?channel=top-nav" class="lnk-doubanapp">下载豆瓣客户端</a>
  <div id="doubanapp-tip">
    <a href="https://www.douban.com/doubanapp/app?channel=qipao" class="tip-link">豆瓣 <span class="version">6.0</span> 全新发布</a>
    <a href="javascript: void 0;" class="tip-close">×</a>
  </div>
  <div id="top-nav-appintro" class="more-items">
    <p class="appintro-title">豆瓣</p>
    <p class="qrcode">扫码直接下载</p>
    <div class="download">
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=iOS">iPhone</a>
      <span>·</span>
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=Android" class="download-android">Android</a>
    </div>
  </div>
</div>

    


<div class="global-nav-items">
  <ul>
    <li class="">
      <a href="https://www.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-main&quot;,&quot;uid&quot;:&quot;0&quot;}">豆瓣</a>
    </li>
    <li class="">
      <a href="https://book.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-book&quot;,&quot;uid&quot;:&quot;0&quot;}">读书</a>
    </li>
    <li class="on">
      <a href="https://movie.douban.com"  data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-movie&quot;,&quot;uid&quot;:&quot;0&quot;}">电影</a>
    </li>
    <li class="">
      <a href="https://music.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-music&quot;,&quot;uid&quot;:&quot;0&quot;}">音乐</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/location" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-location&quot;,&quot;uid&quot;:&quot;0&quot;}">同城</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/group" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-group&quot;,&quot;uid&quot;:&quot;0&quot;}">小组</a>
    </li>
    <li class="">
      <a href="https://read.douban.com&#47;?dcs=top-nav&amp;dcm=douban" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-read&quot;,&quot;uid&quot;:&quot;0&quot;}">阅读</a>
    </li>
    <li class="">
      <a href="https://douban.fm&#47;?from_=shire_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-fm&quot;,&quot;uid&quot;:&quot;0&quot;}">FM</a>
    </li>
    <li class="">
      <a href="https://time.douban.com&#47;?dt_time_source=douban-web_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-time&quot;,&quot;uid&quot;:&quot;0&quot;}">时间</a>
    </li>
    <li class="">
      <a href="https://market.douban.com&#47;?utm_campaign=douban_top_nav&amp;utm_source=douban&amp;utm_medium=pc_web" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-market&quot;,&quot;uid&quot;:&quot;0&quot;}">豆品</a>
    </li>
    <li>
      <a href="#more" class="bn-more"><span>更多</span></a>
      <div class="more-items">
        <table cellpadding="0" cellspacing="0">
          <tbody>
            <tr>
              <td>
                <a href="https://ypy.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-ypy&quot;,&quot;uid&quot;:&quot;0&quot;}">豆瓣摄影</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </li>
  </ul>
</div>

  </div>
</div>
<script>
  ;window._GLOBAL_NAV = {
    DOUBAN_URL: "https://www.douban.com",
    N_NEW_NOTIS: 0,
    N_NEW_DOUMAIL: 0
  };
</script>



    <script src="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.js" defer="defer"></script>




    



    <link href="//img3.doubanio.com/dae/accounts/resources/8c80301/movie/bundle.css" rel="stylesheet" type="text/css">




<div id="db-nav-movie" class="nav">
  <div class="nav-wrap">
  <div class="nav-primary">
    <div class="nav-logo">
      <a href="https:&#47;&#47;movie.douban.com">豆瓣电影</a>
    </div>
    <div class="nav-search">
      <form action="https:&#47;&#47;movie.douban.com/subject_search" method="get">
        <fieldset>
          <legend>搜索：</legend>
          <label for="inp-query">
          </label>
          <div class="inp"><input id="inp-query" name="search_text" size="22" maxlength="60" placeholder="搜索电影、电视剧、综艺、影人" value=""></div>
          <div class="inp-btn"><input type="submit" value="搜索"></div>
          <input type="hidden" name="cat" value="1002" />
        </fieldset>
      </form>
    </div>
  </div>
  </div>
  <div class="nav-secondary">
    

<div class="nav-items">
  <ul>
    <li    ><a href="https://movie.douban.com/cinema/nowplaying/"
     >影讯&购票</a>
    </li>
    <li    ><a href="https://movie.douban.com/explore"
     >选电影</a>
    </li>
    <li    ><a href="https://movie.douban.com/tv/"
     >电视剧</a>
    </li>
    <li    ><a href="https://movie.douban.com/chart"
     >排行榜</a>
    </li>
    <li    ><a href="https://movie.douban.com/tag/"
     >分类</a>
    </li>
    <li    ><a href="https://movie.douban.com/review/best/"
     >影评</a>
    </li>
    <li    ><a href="https://movie.douban.com/annual/2018?source=navigation"
     >2018年度榜单</a>
    </li>
    <li    ><a href="https://www.douban.com/standbyme/2018?source=navigation"
     >2018书影音报告</a>
    </li>
  </ul>
</div>

    <a href="https://movie.douban.com/annual/2018?source=movie_navigation" class="movieannual2018"></a>
  </div>
</div>

<script id="suggResult" type="text/x-jquery-tmpl">
  <li data-link="{{= url}}">
            <a href="{{= url}}" onclick="moreurl(this, {from:'movie_search_sugg', query:'{{= keyword }}', subject_id:'{{= id}}', i: '{{= index}}', type: '{{= type}}'})">
            <img src="{{= img}}" width="40" />
            <p>
                <em>{{= title}}</em>
                {{if year}}
                    <span>{{= year}}</span>
                {{/if}}
                {{if sub_title}}
                    <br /><span>{{= sub_title}}</span>
                {{/if}}
                {{if address}}
                    <br /><span>{{= address}}</span>
                {{/if}}
                {{if episode}}
                    {{if episode=="unknow"}}
                        <br /><span>集数未知</span>
                    {{else}}
                        <br /><span>共{{= episode}}集</span>
                    {{/if}}
                {{/if}}
            </p>
        </a>
        </li>
  </script>




    <script src="//img3.doubanio.com/dae/accounts/resources/8c80301/movie/bundle.js" defer="defer"></script>





    
    <div id="wrapper">
        

        
    <div id="content">
        

    <div id="dale_movie_subject_top_icon"></div>
    <h1>
        <span property="v:itemreviewed">三少爷的剑</span>
            <span class="year">(2000)</span>
    </h1>

        <div class="grid-16-8 clearfix">
            

            
            <div class="article">
                
    

    





        <div class="indent clearfix">
            <div class="subjectwrap clearfix">
                <div class="subject clearfix">
                    



<div id="mainpic" class="">
    <a class="nbgnbg" href="https://movie.douban.com/subject/2279835/photos?type=R" title="点击看更多海报">
        <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2256387990.webp" title="点击看更多海报" alt="三少爷的剑" rel="v:image" />
   </a>
</div>

                    


<div id="info">
        <span ><span class='pl'>导演</span>: <span class='attrs'><a href="/celebrity/1320042/" rel="v:directedBy">靳德茂</a></span></span><br/>
        <span ><span class='pl'>编剧</span>: <span class='attrs'><a href="/subject_search?search_text=%E8%B5%B5%E5%BF%97%E7%BA%A2">赵志红</a> / <a href="/celebrity/1315799/">古龙</a></span></span><br/>
        <span class="actor"><span class='pl'>主演</span>: <span class='attrs'><a href="/celebrity/1314497/" rel="v:starring">何中华</a> / <a href="/celebrity/1014618/" rel="v:starring">俞飞鸿</a> / <a href="/celebrity/1314478/" rel="v:starring">陈龙</a> / <a href="/celebrity/1318452/" rel="v:starring">陈继铭</a> / <a href="/celebrity/1323125/" rel="v:starring">张伊函</a> / <a href="/celebrity/1000145/" rel="v:starring">霍思燕</a> / <a href="/celebrity/1276110/" rel="v:starring">刘莉莉</a> / <a href="/celebrity/1274488/" rel="v:starring">杨若兮</a> / <a href="/celebrity/1274566/" rel="v:starring">石小满</a> / <a href="/celebrity/1313611/" rel="v:starring">岳跃利</a> / <a href="/celebrity/1274559/" rel="v:starring">戴春荣</a> / <a href="/celebrity/1274564/" rel="v:starring">陈莹</a> / <a href="/celebrity/1016668/" rel="v:starring">张静初</a> / <a href="/celebrity/1314844/" rel="v:starring">赵毅</a></span></span><br/>
        <span class="pl">类型:</span> <span property="v:genre">武侠</span> / <span property="v:genre">古装</span><br/>
        
        <span class="pl">制片国家/地区:</span> 中国大陆<br/>
        <span class="pl">语言:</span> 汉语普通话<br/>
        <span class="pl">首播:</span> <span property="v:initialReleaseDate" content="2001">2001</span><br/>
        
        <span class="pl">集数:</span> 34<br/>
        <span class="pl">单集片长:</span> 47分钟<br/>
        
        

</div>




                </div>
                    


<div id="interest_sectl">
    <div class="rating_wrap clearbox" rel="v:rating">
        <div class="clearfix">
          <div class="rating_logo ll">豆瓣评分</div>
          <div class="output-btn-wrap rr" style="display:none">
            <img src="https://img3.doubanio.com/f/movie/692e86756648f29457847c5cc5e161d6f6b8aaac/pics/movie/reference.png" />
            <a class="download-output-image" href="#">引用</a>
          </div>
          
          
        </div>
        


<div class="rating_self clearfix" typeof="v:Rating">
    <strong class="ll rating_num" property="v:average">6.9</strong>
    <span property="v:best" content="10.0"></span>
    <div class="rating_right ">
        <div class="ll bigstar bigstar35"></div>
        <div class="rating_sum">
                <a href="collections" class="rating_people"><span property="v:votes">5475</span>人评价</a>
        </div>
    </div>
</div>
<div class="ratings-on-weight">
    
        <div class="item">
        
        <span class="stars5 starstop" title="力荐">
            5星
        </span>
        <div class="power" style="width:15px"></div>
        <span class="rating_per">11.5%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars4 starstop" title="推荐">
            4星
        </span>
        <div class="power" style="width:43px"></div>
        <span class="rating_per">32.7%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars3 starstop" title="还行">
            3星
        </span>
        <div class="power" style="width:64px"></div>
        <span class="rating_per">47.7%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars2 starstop" title="较差">
            2星
        </span>
        <div class="power" style="width:8px"></div>
        <span class="rating_per">6.6%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars1 starstop" title="很差">
            1星
        </span>
        <div class="power" style="width:2px"></div>
        <span class="rating_per">1.5%</span>
        <br />
        </div>
</div>

    </div>
</div>


                
            </div>
                




<div id="interest_sect_level" class="clearfix">
        
            <a href="https://www.douban.com/reason=collectwish&amp;ck=" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-2279835-wish">
                <span>想看</span>
            </a>
            <a href="https://www.douban.com/reason=collectdo&amp;ck=" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-2279835-do">
                <span>在看</span>
            </a>
            <a href="https://www.douban.com/reason=collectcollect&amp;ck=" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-2279835-collect">
                <span>看过</span>
            </a>
        <div class="ll j a_stars">
            
    
    评价:
    <span id="rating"> <span id="stars" data-solid="https://img3.doubanio.com/f/shire/5a2327c04c0c231bced131ddf3f4467eb80c1c86/pics/rating_icons/star_onmouseover.png" data-hollow="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" data-solid-2x="https://img3.doubanio.com/f/shire/7258904022439076d57303c3b06ad195bf1dc41a/pics/rating_icons/star_onmouseover@2x.png" data-hollow-2x="https://img3.doubanio.com/f/shire/95cc2fa733221bb8edd28ad56a7145a5ad33383e/pics/rating_icons/star_hollow_hover@2x.png">

            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-2279835-1">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star1" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-2279835-2">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star2" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-2279835-3">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star3" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-2279835-4">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star4" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-2279835-5">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star5" width="16" height="16"/></a>
    </span><span id="rateword" class="pl"></span>
    <input id="n_rating" type="hidden" value=""  />
    </span>

        </div>
    

</div>


            


















<div class="gtleft">
    <ul class="ul_subject_menu bicelink color_gray pt6 clearfix">
        
    
        
                <li> 
    <img src="https://img3.doubanio.com/f/shire/cc03d0fcf32b7ce3af7b160a0b85e5e66b47cc42/pics/short-comment.gif" />&nbsp;
        <a onclick="moreurl(this, {from:'mv_sbj_wr_cmnt_login'})" class="j a_show_login" href="https://www.douban.com/register?reason=review" rel="nofollow">写短评</a>
 </li>
                    <li> 
    
    <img src="https://img3.doubanio.com/f/shire/5bbf02b7b5ec12b23e214a580b6f9e481108488c/pics/add-review.gif" />&nbsp;
        <a onclick="moreurl(this, {from:'mv_sbj_wr_rv_login'})" class="j a_show_login" href="https://www.douban.com/register?reason=review" rel="nofollow">写影评</a>
 </li>
                <li> 
    



 </li>
                <li> 
   

   
    
    <span class="rec" id="电视剧-2279835">
    <a href= "#"
        data-type="电视剧"
        data-url="https://movie.douban.com/subject/2279835/"
        data-desc="电视剧《三少爷的剑》 (来自豆瓣) "
        data-title="电视剧《三少爷的剑》 (来自豆瓣) "
        data-pic="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2256387990.jpeg"
        class="bn-sharing ">
        分享到
    </a> &nbsp;&nbsp;
    </span>

    <script>
        if (!window.DoubanShareMenuList) {
            window.DoubanShareMenuList = [];
        }
        var __cache_url = __cache_url || {};

        (function(u){
            if(__cache_url[u]) return;
            __cache_url[u] = true;
            window.DoubanShareIcons = 'https://img3.doubanio.com/f/shire/d15ffd71f3f10a7210448fec5a68eaec66e7f7d0/pics/ic_shares.png';

            var initShareButton = function() {
                $.ajax({url:u,dataType:'script',cache:true});
            };

            if (typeof Do == 'function' && 'ready' in Do) {
                Do(
                    'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
                    'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
                    'https://img3.doubanio.com/f/movie/c4ab132ff4d3d64a83854c875ea79b8b541faf12/js/movie/lib/qrcode.min.js',
                    initShareButton
                );
            } else if(typeof Douban == 'object' && 'loader' in Douban) {
                Douban.loader.batch(
                    'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
                    'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
                    'https://img3.doubanio.com/f/movie/c4ab132ff4d3d64a83854c875ea79b8b541faf12/js/movie/lib/qrcode.min.js'
                ).done(initShareButton);
            }

        })('https://img3.doubanio.com/f/movie/32be6727ed3ad8f6c4a417d8a086355c3e7d1d27/js/movie/lib/sharebutton.js');
    </script>


  </li>
            

    </ul>

    <script type="text/javascript">
        $(function(){
            $(".ul_subject_menu li.rec .bn-sharing").bind("click", function(){
                $.get("/blank?sbj_page_click=bn_sharing");
            });
            $(".ul_subject_menu .create_from_menu").bind("click", function(e){
                e.preventDefault();
                var $el = $(this);
                var glRoot = document.getElementById('gallery-topics-selection');
                if (window.has_gallery_topics && glRoot) {
                    // 判断是否有话题
                    glRoot.style.display = 'block';
                } else {
                    location.href = $el.attr('href');
                }
            });
        });
    </script>
</div>




                





<div class="rec-sec">
<span class="rec">
    <script id="movie-share" type="text/x-html-snippet">
        
    <form class="movie-share" action="/j/share" method="POST">
        <div class="clearfix form-bd">
            <div class="input-area">
                <textarea name="text" class="share-text" cols="72" data-mention-api="https://api.douban.com/shuo/in/complete?alt=xd&amp;callback=?"></textarea>
                <input type="hidden" name="target-id" value="2279835">
                <input type="hidden" name="target-type" value="0">
                <input type="hidden" name="title" value="三少爷的剑‎ (2001)">
                <input type="hidden" name="desc" value="导演 靳德茂 主演 何中华 / 俞飞鸿 / 中国大陆 / 6.9分(5475评价)">
                <input type="hidden" name="redir" value=""/>
                <div class="mentioned-highlighter"></div>
            </div>

            <div class="info-area">
                    <img class="media" src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2256387990.webp" />
                <strong>三少爷的剑‎ (2001)</strong>
                <p>导演 靳德茂 主演 何中华 / 俞飞鸿 / 中国大陆 / 6.9分(5475评价)</p>
                <p class="error server-error">&nbsp;</p>
            </div>
        </div>
        <div class="form-ft">
            <div class="form-ft-inner">
                



                <span class="avail-num-indicator">140</span>
                <span class="bn-flat">
                    <input type="submit" value="推荐" />
                </span>
            </div>
        </div>
    </form>
    
    <div id="suggest-mention-tmpl" style="display:none;">
        <ul>
            {{#users}}
            <li id="{{uid}}">
              <img src="{{avatar}}">{{{username}}}&nbsp;<span>({{{uid}}})</span>
            </li>
            {{/users}}
        </ul>
    </div>


    </script>

        
        <a href="/accounts/register?reason=recommend"  class="j a_show_login lnk-sharing" share-id="2279835" data-mode="plain" data-name="三少爷的剑‎ (2001)" data-type="movie" data-desc="导演 靳德茂 主演 何中华 / 俞飞鸿 / 中国大陆 / 6.9分(5475评价)" data-href="https://movie.douban.com/subject/2279835/" data-image="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2256387990.webp" data-properties="{}" data-redir="" data-text="" data-apikey="" data-curl="" data-count="10" data-object_kind="1002" data-object_id="2279835" data-target_type="rec" data-target_action="1" data-action_props="{&#34;subject_url&#34;:&#34;https:\/\/movie.douban.com\/subject\/2279835\/&#34;,&#34;subject_title&#34;:&#34;三少爷的剑‎ (2001)&#34;}">推荐</a>
</span>


</div>






            <script type="text/javascript">
                $(function() {
                    $('.collect_btn', '#interest_sect_level').each(function() {
                        Douban.init_collect_btn(this);
                    });
                    $('html').delegate(".indent .rec-sec .lnk-sharing", "click", function() {
                        moreurl(this, {
                            from : 'mv_sbj_db_share'
                        });
                    });
                });
            </script>
        </div>
            


    <div id="collect_form_2279835"></div>

        







    <h2>
        <i class="">三少爷的剑的分集短评</i>
              · · · · · ·
    </h2>

    


    
    <div class="episode_list">


            

            <a class=" item" href="/subject/2279835/episode/1/">1集</a>
            

            <a class=" item" href="/subject/2279835/episode/2/">2集</a>
            

            <a class=" item" href="/subject/2279835/episode/3/">3集</a>
            

            <a class=" item" href="/subject/2279835/episode/4/">4集</a>
            

            <a class=" item" href="/subject/2279835/episode/5/">5集</a>
            

            <a class=" item" href="/subject/2279835/episode/6/">6集</a>
            

            <a class=" item" href="/subject/2279835/episode/7/">7集</a>
            

            <a class=" item" href="/subject/2279835/episode/8/">8集</a>
            

            <a class=" item" href="/subject/2279835/episode/9/">9集</a>
            

            <a class=" item" href="/subject/2279835/episode/10/">10集</a>
            

            <a class=" item" href="/subject/2279835/episode/11/">11集</a>
            

            <a class=" item" href="/subject/2279835/episode/12/">12集</a>
            

            <a class=" item" href="/subject/2279835/episode/13/">13集</a>
            

            <a class=" item" href="/subject/2279835/episode/14/">14集</a>
            

            <a class=" item" href="/subject/2279835/episode/15/">15集</a>
            

            <a class=" item" href="/subject/2279835/episode/16/">16集</a>
            

            <a class=" item" href="/subject/2279835/episode/17/">17集</a>
            

            <a class=" item" href="/subject/2279835/episode/18/">18集</a>
            

            <a class=" item" href="/subject/2279835/episode/19/">19集</a>
            

            <a class="hide item" href="/subject/2279835/episode/20/">20集</a>
            

            <a class="hide item" href="/subject/2279835/episode/21/">21集</a>
            

            <a class="hide item" href="/subject/2279835/episode/22/">22集</a>
            

            <a class="hide item" href="/subject/2279835/episode/23/">23集</a>
            

            <a class="hide item" href="/subject/2279835/episode/24/">24集</a>
            

            <a class="hide item" href="/subject/2279835/episode/25/">25集</a>
            

            <a class="hide item" href="/subject/2279835/episode/26/">26集</a>
            

            <a class="hide item" href="/subject/2279835/episode/27/">27集</a>
            

            <a class="hide item" href="/subject/2279835/episode/28/">28集</a>
            

            <a class="hide item" href="/subject/2279835/episode/29/">29集</a>
            

            <a class="hide item" href="/subject/2279835/episode/30/">30集</a>
            

            <a class="hide item" href="/subject/2279835/episode/31/">31集</a>
            

            <a class="hide item" href="/subject/2279835/episode/32/">32集</a>
            

            <a class="hide item" href="/subject/2279835/episode/33/">33集</a>
            

            <a class="hide item" href="/subject/2279835/episode/34/">34集</a>

            <a href="#" class="ep_more"><span></span></a>

    </div>




        



<div class="related-info" style="margin-bottom:-10px;">
    <a name="intro"></a>
    
        
            
            
    <h2>
        <i class="">三少爷的剑的剧情简介</i>
              · · · · · ·
    </h2>

            <div class="indent" id="link-report">
                    
                        <span property="v:summary" class="">
                                　　江湖纷争，麻烦不断。神剑山庄的三少爷谢晓峰（何中华 饰）和慕容世家的独女慕容秋荻（俞飞鸿 饰）在百般波折之下，终于迎来了二人要大喜的日子。可大婚当天，谢晓峰却不得不与夺命十三剑——燕十三（王冰 饰）进行一场比试，而这时意外发生，谢晓峰被他的兄弟铁铉之妻吕香华带走。而这背后竟牵扯到攸关国家生死存亡的大事。
                                    <br />
                                　　燕王朱棣发生兵变，为了保护太子朱文奎的出逃，无奈之举下，谢晓峰答应了铁铉的请求，把他们一家杀死，让太子得以出逃。可这真相外人却不知，于是，谢晓峰背下了天下的骂名，被逐出家门，也和慕容秋荻分散。他们的孩子在凄凉中诞生，这样兜兜转转数十年，一切阴谋与真相到了不得不揭晓的时刻……
                        </span>
                        <span class="pl"><a href="https://movie.douban.com/help/movie#t0-qs">&copy;豆瓣</a></span>
            </div>
</div>


    








<div id="celebrities" class="celebrities related-celebrities">

  
    <h2>
        <i class="">三少爷的剑的演职员</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="/subject/2279835/celebrities">全部 20</a>
            )
            </span>
    </h2>


  <ul class="celebrities-list from-subject __oneline">
        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1320042/" title="靳德茂 De-mao Jin" class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1405557781.07.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1320042/" title="靳德茂 De-mao Jin" class="name">靳德茂</a></span>

      <span class="role" title="导演">导演</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1314497/" title="何中华 Zhonghua He " class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p24707.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1314497/" title="何中华 Zhonghua He " class="name">何中华</a></span>

      <span class="role" title="演员">演员</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1014618/" title="俞飞鸿 Faye Yu" class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1364106535.09.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1014618/" title="俞飞鸿 Faye Yu" class="name">俞飞鸿</a></span>

      <span class="role" title="演员">演员</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1314478/" title="陈龙 Long Chen" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p44205.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1314478/" title="陈龙 Long Chen" class="name">陈龙</a></span>

      <span class="role" title="演员">演员</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1318452/" title="陈继铭 Jiming Chen" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1408451397.63.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1318452/" title="陈继铭 Jiming Chen" class="name">陈继铭</a></span>

      <span class="role" title="演员">演员</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1323125/" title="张伊函 Yihan Zhang" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p55526.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1323125/" title="张伊函 Yihan Zhang" class="name">张伊函</a></span>

      <span class="role" title="演员">演员</span>

    </div>
  </li>


  </ul>
</div>


    


<link rel="stylesheet" href="https://img3.doubanio.com/f/verify/16c7e943aee3b1dc6d65f600fcc0f6d62db7dfb4/entry_creator/dist/author_subject/style.css">
<div id="author_subject" class="author-wrapper">
    <div class="loading"></div>
</div>
<script type="text/javascript">
    var answerObj = {
      ISALL: 'False',
      TYPE: 'tv',
      SUBJECT_ID: '2279835',
      USER_ID: 'None'
    }
</script>
<script type="text/javascript" src="https://img3.doubanio.com/f/movie/61252f2f9b35f08b37f69d17dfe48310dd295347/js/movie/lib/react/15.4/bundle.js"></script>
<script type="text/javascript" src="https://img3.doubanio.com/f/verify/ac140ef86262b845d2be7b859e352d8196f3f6d4/entry_creator/dist/author_subject/index.js"></script>


    









    
    <div id="related-pic" class="related-pic">
        
    
    
    <h2>
        <i class="">三少爷的剑的图片</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="/video/create?subject_id=2279835">添加视频评论</a>&nbsp;|&nbsp;<a href="https://movie.douban.com/subject/2279835/all_photos">图片92</a>&nbsp;·&nbsp;<a href="https://movie.douban.com/subject/2279835/mupload">添加</a>
            )
            </span>
    </h2>


        <ul class="related-pic-bd  ">
                <li>
                    <a href="https://movie.douban.com/photos/photo/1804805724/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p1804805724.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/1804799930/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p1804799930.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/1972569287/"><img src="https://img1.doubanio.com/view/photo/sqxs/public/p1972569287.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2175339426/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p2175339426.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2224943330/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p2224943330.webp" alt="图片" /></a>
                </li>
        </ul>
    </div>



      








<div class="mod">
<div class="hd-ops">
  
  <a class="comment_btn j a_show_login" href="https://www.douban.com/register?reason=discussion" rel="nofollow">
      <span>发起新的讨论</span>
  </a>

</div>

    <h2>
        <i class="">讨论区</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/2279835/discussion/">全部</a>
            )
            </span>
    </h2>

<div class="bd">
<div class="mv-discussion-nav">
<a href="https://movie.douban.com/subject/2279835/discussion/" class="on">最新</a>
<a href="https://movie.douban.com/subject/2279835/discussion/?sort=vote" data-epid="hot">热门</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=1" data-epid="44772" data-num="1">1集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=2" data-epid="44773" data-num="2">2集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=3" data-epid="44774" data-num="3">3集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=4" data-epid="44775" data-num="4">4集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=5" data-epid="44776" data-num="5">5集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=6" data-epid="44777" data-num="6">6集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=7" data-epid="44778" data-num="7">7集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/" data-epid="more" title="更多">&#8230;</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=8" data-epid="44779" data-num="8" class="more-item">8集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=9" data-epid="44780" data-num="9" class="more-item">9集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=10" data-epid="44781" data-num="10" class="more-item">10集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=11" data-epid="44782" data-num="11" class="more-item">11集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=12" data-epid="44783" data-num="12" class="more-item">12集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=13" data-epid="44784" data-num="13" class="more-item">13集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=14" data-epid="44785" data-num="14" class="more-item">14集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=15" data-epid="44786" data-num="15" class="more-item">15集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=16" data-epid="44787" data-num="16" class="more-item">16集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=17" data-epid="44788" data-num="17" class="more-item">17集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=18" data-epid="44789" data-num="18" class="more-item">18集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=19" data-epid="44790" data-num="19" class="more-item">19集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=20" data-epid="44791" data-num="20" class="more-item">20集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=21" data-epid="44792" data-num="21" class="more-item">21集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=22" data-epid="44793" data-num="22" class="more-item">22集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=23" data-epid="44794" data-num="23" class="more-item">23集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=24" data-epid="44795" data-num="24" class="more-item">24集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=25" data-epid="44796" data-num="25" class="more-item">25集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=26" data-epid="44797" data-num="26" class="more-item">26集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=27" data-epid="44798" data-num="27" class="more-item">27集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=28" data-epid="44799" data-num="28" class="more-item">28集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=29" data-epid="44800" data-num="29" class="more-item">29集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=30" data-epid="44801" data-num="30" class="more-item">30集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=31" data-epid="44802" data-num="31" class="more-item">31集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=32" data-epid="44803" data-num="32" class="more-item">32集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=33" data-epid="44804" data-num="33" class="more-item">33集</a>
  <a href="https://movie.douban.com/subject/2279835/discussion/?ep_num=34" data-epid="44805" data-num="34" class="more-item">34集</a>
</div>

<div class="mv-discussion-list discussion-list">
  

<table>
  <thead>
  <tr>
    <td>讨论</td>
    <td>作者</td>
    <td nowrap="nowrap">回应</td>
    <td align="right">最后回应</td>
  </tr>
  </thead>
  <tbody>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/614890854/" title="看到演员表有张静初">看到演员表有张静初</a>
        <span class="with-pic">[图]</span>
    </td>
    <td><a href="https://www.douban.com/people/no1guangming/">獨孤求敗</a></td>
    <td class="reply-num">2</td>
    <td class="time">2018-09-12 09:34</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/17232361/" title="这也许，是天意吧">这也许，是天意吧</a>
    </td>
    <td><a href="https://www.douban.com/people/androtommy/">Andro</a></td>
    <td class="reply-num">2</td>
    <td class="time">2016-11-06 18:53</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/59579786/" title="主演何中华为何不注明？陈龙未出演好吧？！气愤！">主演何中华为何不注明？陈龙未出演好吧？！气愤！</a>
    </td>
    <td><a href="https://www.douban.com/people/luojianxun/">拉来的小提琴手</a></td>
    <td class="reply-num">3</td>
    <td class="time">2015-07-07 22:23</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/21828148/" title="为了俞飞鸿我能忍受这部烂剧...">为了俞飞鸿我能忍受这部烂剧...</a>
    </td>
    <td><a href="https://www.douban.com/people/iamguoguo/">断翅诺言</a></td>
    <td class="reply-num">7</td>
    <td class="time">2012-01-13 21:19</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/24907470/" title="霍思燕才是亮点~~">霍思燕才是亮点~~</a>
    </td>
    <td><a href="https://www.douban.com/people/3693020/">居里</a></td>
    <td class="reply-num">5</td>
    <td class="time">2011-12-14 00:09</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/18621850/" title="重温了……">重温了……</a>
    </td>
    <td><a href="https://www.douban.com/people/angelclaudia/">时有锦绣</a></td>
    <td class="reply-num">4</td>
    <td class="time">2011-06-11 14:54</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/13481276/" title="喜欢俞飞鸿">喜欢俞飞鸿</a>
    </td>
    <td><a href="https://www.douban.com/people/grace2046/">grace</a></td>
    <td class="reply-num">7</td>
    <td class="time">2011-06-11 14:47</td>
  </tr>
  </tbody>
</table>

<a href="https://movie.douban.com/subject/2279835/discussion/">&gt; 全部讨论7条</a>
</div>

<div class="mv-hot-discussion-list hide">
  

<table>
  <thead>
  <tr>
    <td>讨论</td>
    <td>作者</td>
    <td nowrap="nowrap">回应</td>
    <td align="right">最后回应</td>
  </tr>
  </thead>
  <tbody>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/21828148/" title="为了俞飞鸿我能忍受这部烂剧...">为了俞飞鸿我能忍受这部烂剧...</a>
    </td>
    <td><a href="https://www.douban.com/people/iamguoguo/">断翅诺言</a></td>
    <td class="reply-num">7</td>
    <td class="time">2012-01-13 21:19</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/13481276/" title="喜欢俞飞鸿">喜欢俞飞鸿</a>
    </td>
    <td><a href="https://www.douban.com/people/grace2046/">grace</a></td>
    <td class="reply-num">7</td>
    <td class="time">2011-06-11 14:47</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/24907470/" title="霍思燕才是亮点~~">霍思燕才是亮点~~</a>
    </td>
    <td><a href="https://www.douban.com/people/3693020/">居里</a></td>
    <td class="reply-num">5</td>
    <td class="time">2011-12-14 00:09</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/18621850/" title="重温了……">重温了……</a>
    </td>
    <td><a href="https://www.douban.com/people/angelclaudia/">时有锦绣</a></td>
    <td class="reply-num">4</td>
    <td class="time">2011-06-11 14:54</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/59579786/" title="主演何中华为何不注明？陈龙未出演好吧？！气愤！">主演何中华为何不注明？陈龙未出演好吧？！气愤！</a>
    </td>
    <td><a href="https://www.douban.com/people/luojianxun/">拉来的小提琴手</a></td>
    <td class="reply-num">3</td>
    <td class="time">2015-07-07 22:23</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/614890854/" title="看到演员表有张静初">看到演员表有张静初</a>
        <span class="with-pic">[图]</span>
    </td>
    <td><a href="https://www.douban.com/people/no1guangming/">獨孤求敗</a></td>
    <td class="reply-num">2</td>
    <td class="time">2018-09-12 09:34</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/2279835/discussion/17232361/" title="这也许，是天意吧">这也许，是天意吧</a>
    </td>
    <td><a href="https://www.douban.com/people/androtommy/">Andro</a></td>
    <td class="reply-num">2</td>
    <td class="time">2016-11-06 18:53</td>
  </tr>
  </tbody>
</table>

<a href="https://movie.douban.com/subject/2279835/discussion/?sort=vote">&gt; 全部讨论7条</a>
</div>

</div>
</div>






    
    



<style type="text/css">
.award li { display: inline; margin-right: 5px }
.awards { margin-bottom: 20px }
.awards h2 { background: none; color: #000; font-size: 14px; padding-bottom: 5px; margin-bottom: 8px; border-bottom: 1px dashed #dddddd }
.awards .year { color: #666666; margin-left: -5px }
.mod { margin-bottom: 25px }
.mod .hd { margin-bottom: 10px }
.mod .hd h2 {margin:24px 0 3px 0}
</style>



    








    <div id="recommendations" class="">
        
        
    <h2>
        <i class="">喜欢这部剧集的人也喜欢</i>
              · · · · · ·
    </h2>

        
    
    <div class="recommendations-bd">
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/3055383/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2386881058.webp" alt="策马啸西风" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/3055383/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>策马啸西风</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/3114923/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2375000270.webp" alt="白发魔女" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/3114923/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>白发魔女</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/2282477/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2354199260.webp" alt="武林外史" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/2282477/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>武林外史</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/2279825/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2372424064.webp" alt="小李飞刀" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/2279825/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>小李飞刀</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/2347271/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2190567792.webp" alt="萧十一郎" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/2347271/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>萧十一郎</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/2311147/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2214654588.webp" alt="绝代双骄" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/2311147/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>绝代双骄</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/2295783/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p1928871095.webp" alt="碧血剑" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/2295783/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>碧血剑</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/2279816/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2373056002.webp" alt="少年张三丰" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/2279816/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>少年张三丰</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/3098693/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2374999078.webp" alt="金蚕丝雨" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/3098693/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>金蚕丝雨</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/3546595/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2539586282.webp" alt="圆月弯刀" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/3546595/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>圆月弯刀</a>
            </dd>
        </dl>
    </div>

    </div>



        


<script type="text/x-handlebar-tmpl" id="comment-tmpl">
    <div class="dummy-fold">
        {{#each comments}}
        <div class="comment-item" data-cid="id">
            <div class="comment">
                <h3>
                    <span class="comment-vote">
                            <span class="votes">{{votes}}</span>
                        <input value="{{id}}" type="hidden"/>
                        <a href="javascript:;" class="j {{#if ../if_logined}}a_vote_comment{{else}}a_show_login{{/if}}">有用</a>
                    </span>
                    <span class="comment-info">
                        <a href="{{user.path}}" class="">{{user.name}}</a>
                        {{#if rating}}
                        <span class="allstar{{rating}}0 rating" title="{{rating_word}}"></span>
                        {{/if}}
                        <span>
                            {{time}}
                        </span>
                        <p> {{content_tmpl content}} </p>
                    </span>
            </div>
        </div>
        {{/each}}
    </div>
</script>












    

    <div id="comments-section">
        <div class="mod-hd">
            
        <a class="comment_btn j a_show_login" href="https://www.douban.com/register?reason=review" rel="nofollow">
            <span>我要写短评</span>
        </a>

            
            
    <h2>
        <i class="">三少爷的剑的短评</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/2279835/comments?status=P">全部 1036 条</a>
            )
            </span>
    </h2>

        </div>
        <div class="mod-bd">
                
    <div class="tab-hd">
        <a id="hot-comments-tab" href="comments" data-id="hot" class="on">热门</a>&nbsp;/&nbsp;
        <a id="new-comments-tab" href="comments?sort=time" data-id="new">最新</a>&nbsp;/&nbsp;
        <a id="following-comments-tab" href="follows_comments" data-id="following"  class="j a_show_login">好友</a>
    </div>

    <div class="tab-bd">
        <div id="hot-comments" class="tab">
            
    
        
        <div class="comment-item" data-cid="147524057">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">59</span>
                <input value="147524057" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/smarttutu/" class="">十二</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2009-11-01 18:21:08">
                    2009-11-01
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">我只恨他们只有那么短暂的谈情说爱镜头。俞姐姐那双大眼啊，真是顾盼生辉。</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1192975988">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">9</span>
                <input value="1192975988" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/152434281/" class="">元宝宝</a>
                    <span>看过</span>
                    <span class="allstar30 rating" title="还行"></span>
                <span class="comment-time " title="2017-05-21 21:03:39">
                    2017-05-21
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">俞飞鸿简直美的不可方物！</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1020220923">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">2</span>
                <input value="1020220923" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/131809416/" class="">天意怜幽草</a>
                    <span>看过</span>
                    <span class="allstar30 rating" title="还行"></span>
                <span class="comment-time " title="2016-03-11 22:15:35">
                    2016-03-11
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">一个人背负着人人所不耻的骂名，忍辱偷生的过着流浪的生活，只为了大义，剧情是可圈可点，人设也好，演员也是很到位，武打场面挺好看，现今的特效过多毫无味道</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1192997251">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">3</span>
                <input value="1192997251" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/155986723/" class="">梦梦梦梦</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2017-05-21 21:45:49">
                    2017-05-21
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">小的时候看这个，感觉何中华特别帅，特别适合演古装正派角色，可惜一直没有太火。俞飞鸿真的很漂亮，喜欢这种古装武侠剧</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="440281402">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">0</span>
                <input value="440281402" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/32623137/" class="">段古母辛</a>
                    <span>看过</span>
                    <span class="allstar30 rating" title="还行"></span>
                <span class="comment-time " title="2011-09-27 16:40:08">
                    2011-09-27
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">里面几个女的都还蛮漂亮的，但那时候的古装剧好像是一个人拍出来，怎么总拍个大脸啊，真变态。</span>
        </p>
    </div>

        </div>



                
                &gt; <a href="comments?sort=new_score&status=P" data-moreurl-dict={&#34;subject_id&#34;:&#34;2279835&#34;,&#34;from&#34;:&#34;tv-more-comments&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>更多短评1036条</a>
        </div>
        <div id="new-comments" class="tab">
            <div id="normal">
            </div>
            <div class="fold-hd hide">
                <a class="qa" href="/help/opinion#t2-q0" target="_blank">为什么被折叠？</a>
                <a class="btn-unfold" href="#">有一些短评被折叠了</a>
                <div class="qa-tip">
                    评论被折叠，是因为发布这条评论的帐号行为异常。评论仍可以被展开阅读，对发布人的账号不造成其他影响。如果认为有问题，可以<a href="https://help.douban.com/help/ask?category=movie">联系</a>豆瓣电影。
                </div>
            </div>
            <div class="fold-bd">
            </div>
            <span id="total-num"></span>
        </div>
        <div id="following-comments" class="tab">
            
    


        <div class="comment-item">
            你关注的人还没写过短评
        </div>

        </div>
    </div>


            
            
        </div>
    </div>



        

<link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/73ed658484f98d44.css">

<section class="topics mod">
    <header>
        <h2>
            三少爷的剑的话题 · · · · · ·
            <span class="pl">( <span class="gallery_topics">全部 <span id="topic-count"></span> 条</span> )</span>
        </h2>
    </header>

    




<section class="subject-topics">
    <div class="topic-guide" id="topic-guide">
        <img class="ic_question" src="//img3.doubanio.com/f/ithildin/b1a3edea3d04805f899e9d77c0bfc0d158df10d5/pics/export/icon_question.png">
        <div class="tip_content">
            <div class="tip_title">什么是话题</div>
            <div class="tip_desc">
                <div>无论是一部作品、一个人，还是一件事，都往往可以衍生出许多不同的话题。将这些话题细分出来，分别进行讨论，会有更多收获。</div>
            </div>
        </div>
        <img class="ic_guide" src="//img3.doubanio.com/f/ithildin/529f46d86bc08f55cd0b1843d0492242ebbd22de/pics/export/icon_guide_arrow.png">
        <img class="ic_close" id="topic-guide-close" src="//img3.doubanio.com/f/ithildin/2eb4ad488cb0854644b23f20b6fa312404429589/pics/export/close@3x.png">
    </div>

    <div id="topic-items"></div>

    <script>
        window.subject_id = 2279835;
        window.join_label_text = '写剧评参与';

        window.topic_display_count = 4;
        window.topic_item_display_count = 1;
        window.no_content_fun_call_name = "no_topic";

        window.guideNode = document.getElementById('topic-guide');
        window.guideNodeClose = document.getElementById('topic-guide-close');
    </script>
    
        <link rel="stylesheet" href="https://img3.doubanio.com/f/ithildin/f731c9ea474da58c516290b3a6b1dd1237c07c5e/css/export/subject_topics.css">
        <script src="https://img3.doubanio.com/f/ithildin/d3590fc6ac47b33c804037a1aa7eec49075428c8/js/export/moment-with-locales-only-zh.js"></script>
        <script src="https://img3.doubanio.com/f/ithildin/c600fdbe69e3ffa5a3919c81ae8c8b4140e99a3e/js/export/subject_topics.js"></script>

</section>

    <script>
        function no_topic(){
            $('#content .topics').remove();
        }
    </script>
</section>

<section class="reviews mod movie-content">
    <header>
        <a href="new_review" rel="nofollow" class="create-review comment_btn"
            data-isverify="False"
            data-verify-url="https://www.douban.com/accounts/phone/verify?redir=http://movie.douban.com/subject/2279835/new_review">
            <span>我要写剧评</span>
        </a>
        <h2>
            三少爷的剑的剧评 · · · · · ·
            <span class="pl">( <a href="reviews">全部 9 条</a> )</span>
        </h2>
    </header>

    

<style>
#gallery-topics-selection {
  position: fixed;
  width: 595px;
  padding: 40px 40px 33px 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 16px 0 rgba(0, 0, 0, 0.2);
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  z-index: 9999;
}
#gallery-topics-selection h1 {
  font-size: 18px;
  color: #007722;
  margin-bottom: 36px;
  padding: 0;
  line-height: 28px;
  font-weight: normal;
}
#gallery-topics-selection .gl_topics {
  border-bottom: 1px solid #dfdfdf;
  max-height: 298px;
  overflow-y: scroll;
}
#gallery-topics-selection .topic {
  margin-bottom: 24px;
}
#gallery-topics-selection .topic_name {
  font-size: 15px;
  color: #333;
  margin: 0;
  line-height: inherit;
}
#gallery-topics-selection .topic_meta {
  font-size: 13px;
  color: #999;
}
#gallery-topics-selection .topics_skip {
  display: block;
  cursor: pointer;
  font-size: 16px;
  color: #3377AA;
  text-align: center;
  margin-top: 33px;
}
#gallery-topics-selection .topics_skip:hover {
  background: transparent;
}
#gallery-topics-selection .close_selection {
  position: absolute;
  width: 30px;
  height: 20px;
  top: 46px;
  right: 40px;
  background: #fff;
  color: #999;
  text-align: right;
}
#gallery-topics-selection .close_selection:hover{
  background: #fff;
  color: #999;
}
</style>




        <div class="review_filter">
            <a href="javascript:;;" class="cur" data-sort="">热门</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="time">最新</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="follow">好友</a href="javascript:;;">
            
        </div>


        



<div class="review-list  ">
        
    

        
    
    <div data-cid="2207301">
        <div class="main review-item" id="2207301">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/2531185/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/u2531185-67.jpg">
        </a>

        <a href="https://www.douban.com/people/2531185/" class="name">熏衣草的小香水</a>

            <span class="allstar20 main-title-rating" title="较差"></span>

        <span content="2009-08-08" class="main-meta">2009-08-08 17:51:16</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/2207301/">令人嫌弃的三少爷</a></h2>

                <div id="review_2207301_short" class="review-short" data-rid="2207301">
                    <div class="short-content">

                            在古龙笔下的男主角中，若说不喜欢的，三少爷谢晓峰可算是名列前茅。原因无他，只为了被他辜负的慕容秋荻们不值。    真人版的谢晓峰，总是被编剧加上了很多无可奈何的理由，如此片中的“为国为民”。可惜“政治”和“朝廷”这种事情，从来不是古龙世界里的江湖。    古龙...

                        &nbsp;(<a href="javascript:;" id="toggle-2207301-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_2207301_full" class="hidden">
                    <div id="review_2207301_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="2207301" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-2207301">
                                25
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="2207301" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-2207301">
                                12
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/2207301/#comments" class="reply">24回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="5083913">
        <div class="main review-item" id="5083913">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/43768284/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u43768284-3.jpg">
        </a>

        <a href="https://www.douban.com/people/43768284/" class="name">五月晴天</a>

            <span class="allstar30 main-title-rating" title="还行"></span>

        <span content="2011-09-01" class="main-meta">2011-09-01 16:12:48</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/5083913/">情感淡然</a></h2>

                <div id="review_5083913_short" class="review-short" data-rid="5083913">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇剧评可能有剧透</p>

                               十几岁的谢晓峰和秋荻有点做作，一开始对谢晓峰期望很高，但最后还是弃秋荻而去，为什么不跟她讲明当时的原因呢？一个深爱你的女人不管你做了什么错事，只要你还爱她、肯回改，她都会原谅你的；即使你已经不爱她了，也要说明真相，那样她才不会归因错误、走火入魔，这是...

                        &nbsp;(<a href="javascript:;" id="toggle-5083913-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_5083913_full" class="hidden">
                    <div id="review_5083913_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="5083913" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-5083913">
                                12
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="5083913" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-5083913">
                                3
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/5083913/#comments" class="reply">2回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="3498825">
        <div class="main review-item" id="3498825">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/1856756/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u1856756-2.jpg">
        </a>

        <a href="https://www.douban.com/people/1856756/" class="name">sandy</a>

            <span class="allstar10 main-title-rating" title="很差"></span>

        <span content="2010-08-03" class="main-meta">2010-08-03 16:31:35</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/3498825/">痛心</a></h2>

                <div id="review_3498825_short" class="review-short" data-rid="3498825">
                    <div class="short-content">

                        他最爱的三少爷，我最爱的俞飞鸿，怎么就弄成了这么个鬼样子！ 剧情更是乱七八糟、莫名其妙！ 明明三少爷是无羁绊的，而她也是个绝顶的女人。

                        &nbsp;(<a href="javascript:;" id="toggle-3498825-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_3498825_full" class="hidden">
                    <div id="review_3498825_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="3498825" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-3498825">
                                7
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="3498825" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-3498825">
                                3
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/3498825/#comments" class="reply">1回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9709907">
        <div class="main review-item" id="9709907">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/rywbl/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u125046742-3.jpg">
        </a>

        <a href="https://www.douban.com/people/rywbl/" class="name">媛心清烃</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2018-10-16" class="main-meta">2018-10-16 22:59:21</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9709907/">评 三少爷的剑 剧版</a></h2>

                <div id="review_9709907_short" class="review-short" data-rid="9709907">
                    <div class="short-content">

                        不是偶然的原因，选择去看这部剧其实我当初想追这部剧吧，是因为沉鱼的原因的，具体原因不便透露。对于武侠吧，也不怎么感冒，不是特意的喜欢古龙或者是武侠又或者是经典而去追这部剧的，但是一旦选择了这部剧，真的发现这部剧是挺好看的，十几年前的电视剧，不仅仅是经典，而...

                        &nbsp;(<a href="javascript:;" id="toggle-9709907-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9709907_full" class="hidden">
                    <div id="review_9709907_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9709907" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9709907">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9709907" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9709907">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9709907/#comments" class="reply">1回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="6569101">
        <div class="main review-item" id="6569101">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/2744224/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/u2744224-39.jpg">
        </a>

        <a href="https://www.douban.com/people/2744224/" class="name">加西亚</a>

            <span class="allstar10 main-title-rating" title="很差"></span>

        <span content="2014-03-01" class="main-meta">2014-03-01 22:27:10</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/6569101/">奇葩</a></h2>

                <div id="review_6569101_short" class="review-short" data-rid="6569101">
                    <div class="short-content">

                        看了简介就直接不想看了。  虽然书写的比较杯具和让人郁闷。但是这个电视剧，明显就是名字而已。其他的都是编剧自己写的，和小说完全无关。  我还是很想看小说版本的电视剧的。  谢晓峰虽然让人讨厌，因为很自私。可是里面还是有很多很精彩的剧情的。。。  这电视剧，您改名字...

                        &nbsp;(<a href="javascript:;" id="toggle-6569101-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_6569101_full" class="hidden">
                    <div id="review_6569101_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="6569101" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-6569101">
                                2
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="6569101" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-6569101">
                                4
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/6569101/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>


        <div class="fold-hd">
                <a class="btn-unfold" href="#">有一些影评被折叠了</a>
                    <a class="qa" href="https://help.douban.com/opinion?app=movie#t1-q2">为什么被折叠？</a>
            <div class="qa-tip">评论被折叠，是因为发布这条评论的帐号行为异常。评论仍可以被展开阅读，对发布人的账号不造成其他影响。如果认为有问题，可以<a href="https://help.douban.com/help/ask?category=movie">联系</a>豆瓣电影。</div>
        </div>
        <div class="fold-bd">
                
    
    <div data-cid="9564004">
        <div class="main review-item" id="9564004">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/161332326/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u161332326-5.jpg">
        </a>

        <a href="https://www.douban.com/people/161332326/" class="name">沙里晶</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2018-08-04" class="main-meta">2018-08-04 20:20:14</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9564004/">说的是剑，谈的却是人生</a></h2>

                <div id="review_9564004_short" class="review-short" data-rid="9564004">
                    <div class="short-content">

                        三少爷谢晓峰，是古龙笔下的经典人物，他用跌宕起伏的人生，验证了杀机重重的江湖，和喜怒无常的人心，这样一个出身于“一剑功成万骨枯”的绝世高手，其实是众多江湖传说中的巨大悲剧故事，既有人在高处不胜寒的孤独，也有放下武器为凡人的向往，但正如常说：人在江湖，身不由...

                        &nbsp;(<a href="javascript:;" id="toggle-9564004-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9564004_full" class="hidden">
                    <div id="review_9564004_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9564004" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9564004">
                                9
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9564004" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9564004">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9564004/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

                
    
    <div data-cid="6585279">
        <div class="main review-item" id="6585279">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/27776377/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/u27776377-8.jpg">
        </a>

        <a href="https://www.douban.com/people/27776377/" class="name">紫仓鼠</a>

            <span class="allstar30 main-title-rating" title="还行"></span>

        <span content="2014-03-12" class="main-meta">2014-03-12 22:28:33</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/6585279/">那个梦想的年代</a></h2>

                <div id="review_6585279_short" class="review-short" data-rid="6585279">
                    <div class="short-content">

                        这是一篇很水的随笔，想到哪写到哪。  近日《大丈夫》大火，人人惊叹俞美人容颜不老，可我还是觉得她的角色变了。俞美人演过一次又一次武林/天下第一美人，而我印象最深刻的角色是杨艳（尽管那里面武林第一美人应该是林诗音吧）。惊鸿仙子的眼中全是智慧与笃定，爱得自信而真切...

                        &nbsp;(<a href="javascript:;" id="toggle-6585279-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_6585279_full" class="hidden">
                    <div id="review_6585279_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="6585279" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-6585279">
                                6
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="6585279" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-6585279">
                                1
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/6585279/#comments" class="reply">1回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

                
    
    <div data-cid="8870656">
        <div class="main review-item" id="8870656">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/49886349/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/user_normal_f.jpg">
        </a>

        <a href="https://www.douban.com/people/49886349/" class="name">月舞于星渊</a>

            <span class="allstar30 main-title-rating" title="还行"></span>

        <span content="2017-10-17" class="main-meta">2017-10-17 12:51:47</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8870656/">童年回忆褪去，只是一部烂剧</a></h2>

                <div id="review_8870656_short" class="review-short" data-rid="8870656">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇剧评可能有剧透</p>

                        2017.10.17《三少爷的剑》        虽然有童年回忆做加持，但我不得不承认，这是一部很糟糕的剧。        重温这部剧的原因有两个，一个是看了飞鸿姐姐的《十三邀》，惊讶于她的美貌和风度，所以又想到了这部老剧；另一个是当年追剧的时候觉得自己根本就没看懂剧情内容，所以想...

                        &nbsp;(<a href="javascript:;" id="toggle-8870656-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8870656_full" class="hidden">
                    <div id="review_8870656_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8870656" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8870656">
                                1
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8870656" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8870656">
                                2
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8870656/#comments" class="reply">4回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

                
    
    <div data-cid="8583845">
        <div class="main review-item" id="8583845">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/137137237/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u137137237-1.jpg">
        </a>

        <a href="https://www.douban.com/people/137137237/" class="name">周小环</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2017-06-05" class="main-meta">2017-06-05 16:50:34</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8583845/">这就是个悲剧。。。</a></h2>

                <div id="review_8583845_short" class="review-short" data-rid="8583845">
                    <div class="short-content">

                        六月十七日，京城发生兵变，燕王朱棣率军攻陷都城，继承了帝位。铁铉带着年仅七岁的逊太子朱文奎出逃，被朱棣派的大批高手追捕。万般无奈之下，铁铉夫妇与谢晓峰合演了一场苦肉计，铁铉把太子的服装套在自己亲生儿子的身上，让谢晓峰杀死自己和儿子，以使敌人产生误解就此罢手...

                        &nbsp;(<a href="javascript:;" id="toggle-8583845-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8583845_full" class="hidden">
                    <div id="review_8583845_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8583845" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8583845">
                                1
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8583845" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8583845">
                                1
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8583845/#comments" class="reply">1回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        </div>


    

    

    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/1ccaa3188342e8c5.js"></script>
    <!-- COLLECTED CSS -->
</div>








            <p class="pl">
                &gt;
                <a href="reviews">
                    更多剧评9篇
                </a>
            </p>
</section>

<!-- COLLECTED JS -->

    <br/>

        <div class="section-discussion">
                <p class="discussion_link">
    <a href="https://movie.douban.com/subject/2279835/tv_discuss">&gt; 查看 三少爷的剑 的分集短评（全部9条）</a>
</p>

        </div>


    <script type="text/javascript">
        $(function(){if($.browser.msie && $.browser.version == 6.0){
            var $info = $('#info'),
                maxWidth = parseInt($info.css('max-width'));
            if($info.width() > maxWidth) {
                $info.width(maxWidth);
            }
        }});
    </script>


            </div>
            <div class="aside">
                


    








        






    

<script id="episode-tmpl" type="text/x-jsrender">
<div id="tv-play-source" class="play-source">
    <div class="cross">
        <span style="color:#494949; font-size:16px">{{:cn}}</span>
        <span style="cursor:pointer">✕</span>
    </div>
    <div class="episode-list">
        {{for playlist}}
            <a href="{{:play_link}}&episode={{:ep}}" target="_blank">{{:ep}}集</a>
        {{/for}}
     <div>
 </div>
</script>

<div class="gray_ad">
    
    <h2>
        在哪儿看这部剧集
            &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;
    </h2>

    
    <ul class="bs">
                <li>
                        <a class="playBtn" data-cn="优酷视频" href="https://www.douban.com/link2/?url=http%3A%2F%2Fv.youku.com%2Fv_show%2Fid_XMjE3OTI3NDQ1Mg%3D%3D.html%3Ftpa%3DdW5pb25faWQ9MTAzNTY1XzEwMDAwMV8wMV8wMQ&amp;subtype=3&amp;type=online-video" target="_blank">
                            优酷视频
                        </a>
                    <span class="buylink-price"><span>
                        免费观看
                    </span></span>
                </li>
                <li>
                        <a class="playBtn" data-cn="哔哩哔哩" data-source="8"  href="javascript: void 0;">
                            哔哩哔哩
                        </a>
                    <span class="buylink-price"><span>
                        免费观看
                    </span></span>
                </li>

    </ul>
</div>


    <!-- douban ad begin -->
    <div id="dale_movie_subject_top_right"></div>
    <div id="dale_movie_subject_top_middle"></div>
    <!-- douban ad end -->

    



<style type="text/css">
    .m4 {margin-bottom:8px; padding-bottom:8px;}
    .movieOnline {background:#FFF6ED; padding:10px; margin-bottom:20px;}
    .movieOnline h2 {margin:0 0 5px;}
    .movieOnline .sitename {line-height:2em; width:160px;}
    .movieOnline td,.movieOnline td a:link,.movieOnline td a:visited{color:#666;}
    .movieOnline td a:hover {color:#fff;}
    .link-bt:link,
    .link-bt:visited,
    .link-bt:hover,
    .link-bt:active {margin:5px 0 0; padding:2px 8px; background:#a8c598; color:#fff; -moz-border-radius: 3px; -webkit-border-radius: 3px; border-radius: 3px; display:inline-block;}
</style>



    







    
    <div class="tags">
        
        
    <h2>
        <i class="">豆瓣成员常用的标签</i>
              · · · · · ·
    </h2>

        <div class="tags-body">
                <a href="/tag/武侠" class="">武侠</a>
                <a href="/tag/电视剧" class="">电视剧</a>
                <a href="/tag/古龙" class="">古龙</a>
                <a href="/tag/俞飞鸿" class="">俞飞鸿</a>
                <a href="/tag/国产电视剧" class="">国产电视剧</a>
                <a href="/tag/古装" class="">古装</a>
                <a href="/tag/看过的电视剧" class="">看过的电视剧</a>
                <a href="/tag/华语标签" class="">华语标签</a>
        </div>
    </div>


    <div id="dale_movie_subject_inner_middle"></div>
    <div id="dale_movie_subject_download_middle"></div>
        








<div id="subject-doulist">
    
    
    <h2>
        <i class="">以下豆列推荐</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/2279835/doulists">全部</a>
            )
            </span>
    </h2>


    
    <ul>
            <li>
                <a href="https://www.douban.com/doulist/111448/" target="_blank">80末出生的人是这样成长的</a>
                <span>(沉歌)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/558514/" target="_blank">♡°那些年我们一起追过的「古装剧」</a>
                <span>(姽婳。)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/642764/" target="_blank">我看过的电视剧——古装情结</a>
                <span>(辰夕)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/40978873/" target="_blank">影视，良心制作大杂烩</a>
                <span>(李大虾)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/991297/" target="_blank">古装剧是场不能苛刻的梦。壹。</a>
                <span>(某J。624)</span>
            </li>
    </ul>

</div>

        








<div id="subject-others-interests">
    
    
    <h2>
        <i class="">谁在看这部剧集</i>
              · · · · · ·
    </h2>

    
    <ul class="">
            
            <li class="">
                <a href="https://www.douban.com/people/59949251/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u59949251-5.jpg" class="pil" alt="kobesdu">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/59949251/" class="">kobesdu</a>
                    <div class="">
                        今天下午
                        看过
                        
                    </div>
                </div>
            </li>
            
            <li class="">
                <a href="https://www.douban.com/people/148519898/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u148519898-1.jpg" class="pil" alt="忽然而已">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/148519898/" class="">忽然而已</a>
                    <div class="">
                        昨天
                        想看
                        
                    </div>
                </div>
            </li>
            
            <li class="">
                <a href="https://www.douban.com/people/159200542/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u159200542-1.jpg" class="pil" alt="默默">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/159200542/" class="">默默</a>
                    <div class="">
                        昨天
                        想看
                        
                    </div>
                </div>
            </li>
    </ul>

    
    <div class="subject-others-interests-ft">
        
            <a href="https://movie.douban.com/subject/2279835/doings">54人在看</a>
                &nbsp;/&nbsp;
            <a href="https://movie.douban.com/subject/2279835/collections">7034人看过</a>
                &nbsp;/&nbsp;
            <a href="https://movie.douban.com/subject/2279835/wishes">484人想看</a>
    </div>

</div>



    
    

<!-- douban ad begin -->
<div id="dale_movie_subject_middle_right"></div>
<script type="text/javascript">
    (function (global) {
        if(!document.getElementsByClassName) {
            document.getElementsByClassName = function(className) {
                return this.querySelectorAll("." + className);
            };
            Element.prototype.getElementsByClassName = document.getElementsByClassName;

        }
        var articles = global.document.getElementsByClassName('article'),
            asides = global.document.getElementsByClassName('aside');

        if (articles.length > 0 && asides.length > 0 && articles[0].offsetHeight >= asides[0].offsetHeight) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_movie_subject_middle_right');
        }
    })(this);
</script>
<!-- douban ad end -->



    <br/>

    
<p class="pl">订阅三少爷的剑的影评: <br/><span class="feed">
    <a href="https://movie.douban.com/feed/subject/2279835/reviews"> feed: rss 2.0</a></span></p>


            </div>
            <div class="extra">
                
    
<!-- douban ad begin -->
<div id="dale_movie_subject_bottom_super_banner"></div>
<script type="text/javascript">
    (function (global) {
        var body = global.document.body,
            html = global.document.documentElement;

        var height = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight, html.scrollHeight, html.offsetHeight);
        if (height >= 2000) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_movie_subject_bottom_super_banner');
        }
    })(this);
</script>
<!-- douban ad end -->


            </div>
        </div>
    </div>

        
    <div id="footer">
            <div class="footer-extra"></div>
        
<span id="icp" class="fleft gray-link">
    &copy; 2005－2019 douban.com, all rights reserved 北京豆网科技有限公司
</span>

<a href="https://www.douban.com/hnypt/variformcyst.py" style="display: none;"></a>

<span class="fright">
    <a href="https://www.douban.com/about">关于豆瓣</a>
    · <a href="https://www.douban.com/jobs">在豆瓣工作</a>
    · <a href="https://www.douban.com/about?topic=contactus">联系我们</a>
    · <a href="https://www.douban.com/about?policy=disclaimer">免责声明</a>
    
    · <a href="https://help.douban.com/?app=movie" target="_blank">帮助中心</a>
    · <a href="https://www.douban.com/doubanapp/">移动应用</a>
    · <a href="https://www.douban.com/partner/">豆瓣广告</a>
</span>

    </div>

    </div>
    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/13f9f0eb40433947.js"></script><script type="text/javascript">
        if (!Do.ready) {
            !function(){var t,e,n=document,r=window,i=window.__external_files_loaded=window.__external_files_loaded||{},o=window.__external_files_loading=window.__external_files_loading||{},a=function(t){return t.constructor===Array},s={autoLoad:!0,coreLib:["//img3.doubanio.com/js/jquery.min.js"],mods:{}},c=n.getElementsByTagName("script"),d=c[c.length-1],u=[],l=!1,f=[],h=function(t,e,r,a,s){var d=c[0];if(t){if(i[t])return o[t]=!1,void(a&&a(t,s));if(o[t])return void setTimeout(function(){h(t,e,r,a,s)},1);o[t]=!0;var u,l=e||t.toLowerCase().split(/\./).pop().replace(/[\?#].*/,"");if("js"===l?(u=n.createElement("script"),u.setAttribute("type","text/javascript"),u.setAttribute("src",t),u.setAttribute("async",!0)):"css"===l&&(u=n.createElement("link"),u.setAttribute("type","text/css"),u.setAttribute("rel","stylesheet"),u.setAttribute("href",t),i[t]=!0),u){if(r&&(u.charset=r),"css"===l)return d.parentNode.insertBefore(u,d),void(a&&a(t,s));u.onload=u.onreadystatechange=function(){this.readyState&&"loaded"!==this.readyState&&"complete"!==this.readyState||(i[this.getAttribute("src")]=!0,a&&a(this.getAttribute("src"),s),u.onload=u.onreadystatechange=null)},d.parentNode.insertBefore(u,d)}}},p=function(t){if(t&&a(t)){for(var e,n=0,r=[],i=s.mods,o=[],c={},d=function(t){var e,n,r=0;if(c[t])return o;if(c[t]=!0,i[t].requires){for(n=i[t].requires;"undefined"!=typeof(e=n[r++]);)i[e]?(d(e),o.push(e)):o.push(e);return o}return o};"undefined"!=typeof(e=t[n++]);)i[e]&&i[e].requires&&i[e].requires[0]&&(o=[],c={},r=r.concat(d(e))),r.push(e);return r}},y=function(){l=!0,u.length>0&&(e.apply(this,u),u=[])},m=function(){n.addEventListener?n.removeEventListener("DOMContentLoaded",m,!1):n.attachEvent&&n.detachEvent("onreadystatechange",m),y()},v=function(){if(!l){try{n.documentElement.doScroll("left")}catch(t){return r.setTimeout(v,1)}y()}},g=function(){if("complete"===n.readyState)return r.setTimeout(y,1);var t=!1;if(n.addEventListener)n.addEventListener("DOMContentLoaded",m,!1),r.addEventListener("load",y,!1);else if(n.attachEvent){n.attachEvent("onreadystatechange",m),r.attachEvent("onload",y);try{t=null===r.frameElement}catch(t){}document.documentElement.doScroll&&t&&v()}},E=function(t){t&&a(t)&&(this.queue=t,this.current=null)};E.prototype={_interval:10,start:function(){return this.current=this.next(),this.current?void this.run():void(this.end=!0)},run:function(){var t,e=this,n=this.current;return"function"==typeof n?(n(),void this.start()):void("string"==typeof n&&(s.mods[n]?(t=s.mods[n],h(t.path,t.type,t.charset,function(t){e.start()},e)):/\.js|\.css/i.test(n)?h(n,"","",function(t,e){e.start()},e):this.start()))},next:function(){return this.queue.shift()}},t=d.getAttribute("data-cfg-autoload"),"string"==typeof t&&(s.autoLoad="true"===t.toLowerCase()),t=d.getAttribute("data-cfg-corelib"),"string"==typeof t&&(s.coreLib=t.split(",")),e=function(){var t,e=[].slice.call(arguments);f.length>0&&(e=f.concat(e)),s.autoLoad&&(e=s.coreLib.concat(e)),t=new E(p(e)),t.start()},e.add=function(t,e){t&&e&&e.path&&(s.mods[t]=e)},e.delay=function(){var t=[].slice.call(arguments),n=t.shift();r.setTimeout(function(){e.apply(this,t)},n)},e.global=function(){var t=[].slice.call(arguments);f=f.concat(t)},e.ready=function(){var t=[].slice.call(arguments);return l?e.apply(this,t):void(u=u.concat(t))},e.css=function(t){var e=n.getElementById("do-inline-css");e||(e=n.createElement("style"),e.type="text/css",e.id="do-inline-css",n.getElementsByTagName("head")[0].appendChild(e)),e.styleSheet?e.styleSheet.cssText=e.styleSheet.cssText+t:e.appendChild(n.createTextNode(t))},s.autoLoad&&e(s.coreLib),e.define=e.add,e._config=s,e._mods=s.mods,this.Do=e,g()}();

        }
        Do.ready(
            'https://img3.doubanio.com/f/movie/b2a06a0332fc1526f4caaf8c76c2717d24da408d/js/movie/lib/jsrender.min.js',
            function(){
                $(document).on('click', '.cross span', function(e) {
                    var $this = $(this);
                    var $dialog = $this.parents('#tv-play-source');
                    $dialog.remove();
                });
                $('body').bind('click', function(e) {
                    var $this = $(e.target),
                        $source = $('.play-source');
                    if (!$this.is('.playBtn') && !$this.parents('.play-source').length) {
                        $source.remove();
                    }
                });
                var sources = {};
                sources[3] = [
                            {play_link: "https://www.douban.com/link2/?url=http%3A%2F%2Fv.youku.com%2Fv_show%2Fid_XMjE3OTI3NDQ1Mg%3D%3D.html%3Ftpa%3DdW5pb25faWQ9MTAzNTY1XzEwMDAwMV8wMV8wMQ&amp;subtype=3&amp;type=online-video", ep: "1"},
                ];
                sources[8] = [
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143183%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "1"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143184%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "2"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143185%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "3"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143186%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "4"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143187%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "5"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143188%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "6"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143189%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "7"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143190%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "8"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143191%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "9"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143192%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "10"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143193%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "11"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143194%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "12"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143195%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "13"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143196%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "14"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143197%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "15"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143198%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "16"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143199%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "17"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143200%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "18"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143201%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "19"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143202%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "20"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143203%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "21"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143204%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "22"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143205%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "23"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143206%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "24"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143207%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "25"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143208%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "26"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143209%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "27"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143210%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "28"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143211%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "29"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143212%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "30"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143213%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "31"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143214%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "32"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143215%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "33"},
                            {play_link: "https://www.douban.com/link2/?url=https%3A%2F%2Fwww.bilibili.com%2Fbangumi%2Fplay%2Fep143216%3Fbsource%3Ddouban&amp;subtype=8&amp;type=online-video&amp;link2key=9d72eb3373", ep: "34"},
                ];
                $('.playBtn').click(function(e){
                    $('.play-source').remove();

                    var height, width, tmpl, cn, source;
                    var $dialog = $('#tv-play-source');
                    var $this = $(this);
                    source = $this.data('source');
                    if(!source)return;
                    cn = $this.data('cn');

                    tmpl = $.templates('#episode-tmpl');
                    $dialog = $(tmpl.render({playlist: sources[source], cn: cn}));

                    $dialog.hide();
                    $('body').append($dialog);
                    width = $dialog.outerWidth();

                    $dialog.css({
                        marginLeft: -width / 2,
                        top: $this.offset().top + $this.outerHeight()
                    }).show();
                });
            });
    </script><script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/3255d8b0f0133dae.js"></script>
        
        
    <link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css" />
    <link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/movie/1d829b8605b9e81435b127cbf3d16563aaa51838/css/movie/mod/reg_login_pop.css" />
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/77323ae72a612bba8b65f845491513ff3329b1bb/js/do.js" data-cfg-autoload="false"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js"></script>
    <script type="text/javascript">
        var HTTPS_DB='https://www.douban.com';
var account_pop={open:function(o,e){e?referrer="?referrer="+encodeURIComponent(e):referrer="?referrer="+window.location.href;var n="",i="",t=382;"reg"===o?(n="用户注册",i="https://accounts.douban.com/popup/login?source=movie#popup_register",t=480):"login"===o&&(n="用户登录",i="https://accounts.douban.com/popup/login?source=movie");var r=document.location.protocol+"//"+document.location.hostname,a=dui.Dialog({width:478,title:n,height:t,cls:"account_pop",isHideTitle:!0,modal:!0,content:"<iframe scrolling='no' frameborder='0' width='478' height='"+t+"' src='"+i+"' name='"+r+"'></iframe>"},!0),c=a.node;if(c.undelegate(),c.delegate(".dui-dialog-close","click",function(){var o=$("body");o.find("#login_msk").hide(),o.find(".account_pop").remove()}),$(window).width()<478){var u="";"reg"===o?u=HTTPS_DB+"/accounts/register"+referrer:"login"===o&&(u=HTTPS_DB+"/accounts/login"+referrer),window.location.href=u}else a.open();$(window).bind("message",function(o){"https://accounts.douban.com"===o.originalEvent.origin&&(c.find("iframe").css("height",o.originalEvent.data),c.height(o.originalEvent.data),a.update())})}};Douban&&Douban.init_show_login&&(Douban.init_show_login=function(o){var e=$(o);e.click(function(){var o=e.data("ref")||"";return account_pop.open("login",o),!1})}),Do(function(){$("body").delegate(".pop_register","click",function(o){o.preventDefault();var e=$(this).data("ref")||"";return account_pop.open("reg",e),!1}),$("body").delegate(".pop_login","click",function(o){o.preventDefault();var e=$(this).data("ref")||"";return account_pop.open("login",e),!1})});
    </script>

    
    
    
    




    
<script type="text/javascript">
    (function (global) {
        var newNode = global.document.createElement('script'),
            existingNode = global.document.getElementsByTagName('script')[0],
            adSource = '//erebor.douban.com/',
            userId = '',
            browserId = 'auSsK8Dk5cg',
            criteria = '7:杨若兮|7:戴春荣|7:陈继铭|7:刘大刚|7:赵毅|7:武侠|7:古龙|7:电视剧|7:看过的电视剧|7:陈莹|7:张静初|7:国产电视剧|7:岳跃利|7:陈龙|7:霍思燕|7:何中华|7:石小满|7:刘莉莉|7:俞飞鸿|7:古装|7:中国|7:华语标签|7:张伊函|7:靳德茂|3:/subject/2279835/',
            preview = '',
            debug = false,
            adSlots = ['dale_movie_subject_top_icon', 'dale_movie_subject_top_right', 'dale_movie_subject_top_middle', 'dale_movie_subject_inner_middle', 'dale_movie_subject_download_middle'];

        global.DoubanAdRequest = {src: adSource, uid: userId, bid: browserId, crtr: criteria, prv: preview, debug: debug};
        global.DoubanAdSlots = (global.DoubanAdSlots || []).concat(adSlots);

        newNode.setAttribute('type', 'text/javascript');
        newNode.setAttribute('src', 'https://img3.doubanio.com/f/adjs/dd37385211bc8deb01376096bfa14d2c0436a98c/ad.release.js');
        newNode.setAttribute('async', true);
        existingNode.parentNode.insertBefore(newNode, existingNode);
    })(this);
</script>











    
  









<script type="text/javascript">
var _paq = _paq || [];
_paq.push(['trackPageView']);
_paq.push(['enableLinkTracking']);
(function() {
    var p=(('https:' == document.location.protocol) ? 'https' : 'http'), u=p+'://fundin.douban.com/';
    _paq.push(['setTrackerUrl', u+'piwik']);
    _paq.push(['setSiteId', '100001']);
    var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
    g.type='text/javascript';
    g.defer=true;
    g.async=true;
    g.src=p+'://img3.doubanio.com/dae/fundin/piwik.js';
    s.parentNode.insertBefore(g,s);
})();
</script>

<script type="text/javascript">
var setMethodWithNs = function(namespace) {
  var ns = namespace ? namespace + '.' : ''
    , fn = function(string) {
        if(!ns) {return string}
        return ns + string
      }
  return fn
}

var gaWithNamespace = function(fn, namespace) {
  var method = setMethodWithNs(namespace)
  fn.call(this, method)
}

var _gaq = _gaq || []
  , accounts = [
      { id: 'UA-7019765-1', namespace: 'douban' }
    , { id: 'UA-7019765-19', namespace: '' }
    ]
  , gaInit = function(account) {
      gaWithNamespace(function(method) {
        gaInitFn.call(this, method, account)
      }, account.namespace)
    }
  , gaInitFn = function(method, account) {
      _gaq.push([method('_setAccount'), account.id]);
      _gaq.push([method('_setSampleRate'), '5']);

      
  _gaq.push([method('_addOrganic'), 'google', 'q'])
  _gaq.push([method('_addOrganic'), 'baidu', 'wd'])
  _gaq.push([method('_addOrganic'), 'soso', 'w'])
  _gaq.push([method('_addOrganic'), 'youdao', 'q'])
  _gaq.push([method('_addOrganic'), 'so.360.cn', 'q'])
  _gaq.push([method('_addOrganic'), 'sogou', 'query'])
  if (account.namespace) {
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣'])
    _gaq.push([method('_addIgnoredOrganic'), 'douban'])
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣网'])
    _gaq.push([method('_addIgnoredOrganic'), 'www.douban.com'])
  }

      if (account.namespace === 'douban') {
        _gaq.push([method('_setDomainName'), '.douban.com'])
      }

        _gaq.push([method('_setCustomVar'), 1, 'responsive_view_mode', 'desktop', 3])

        _gaq.push([method('_setCustomVar'), 2, 'login_status', '0', 2]);

      _gaq.push([method('_trackPageview')])
    }

for(var i = 0, l = accounts.length; i < l; i++) {
  var account = accounts[i]
  gaInit(account)
}


;(function() {
    var ga = document.createElement('script');
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    ga.setAttribute('async', 'true');
    document.documentElement.firstChild.appendChild(ga);
})()
</script>








      
    

    <!-- brand10-docker-->

  <script>_SPLITTEST=''</script>
</body>

</html>


`
var ZhiHuByte = []byte(`
<!doctype html>
<html lang="zh" data-hairline="true" data-theme="light"><head><meta charSet="utf-8"/><title data-react-helmet="true">你为什么支持死刑？ - 知乎</title><meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"/><meta name="renderer" content="webkit"/><meta name="force-rendering" content="webkit"/><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/><meta name="google-site-verification" content="FTeR0c8arOPKh8c5DYh_9uu98_zJbaWw53J-Sch9MTg"/><title>知乎 - 有问题上知乎</title><meta name="description" content="有问题，上知乎。知乎是中文互联网知名知识分享平台，以「知识连接一切」为愿景，致力于构建一个人人都可以便捷接入的知识分享网络，让人们便捷地与世界分享知识、经验和见解，发现更大的世界。"/><link rel="shortcut icon" type="image/x-icon" href="https://static.zhihu.com/static/favicon.ico"/><link rel="search" type="application/opensearchdescription+xml" href="https://static.zhihu.com/static/search.xml" title="知乎"/><link rel="dns-prefetch" href="//static.zhimg.com"/><link rel="dns-prefetch" href="//pic1.zhimg.com"/><link rel="dns-prefetch" href="//pic2.zhimg.com"/><link rel="dns-prefetch" href="//pic3.zhimg.com"/><link rel="dns-prefetch" href="//pic4.zhimg.com"/><link href="https://static.zhihu.com/heifetz/main.app.65cf65ee96a255525453.css" rel="stylesheet"/><script defer="" crossOrigin="anonymous" src="https://unpkg.zhimg.com/@cfe/sentry-script@latest/dist/init.js" data-sentry-config="{&quot;dsn&quot;:&quot;https://65e244586890460588f00f2987137aa8@crash2.zhihu.com/193&quot;,&quot;sampleRate&quot;:0.1,&quot;release&quot;:&quot;848-2324b460&quot;,&quot;ignoreErrorNames&quot;:[&quot;NetworkError&quot;,&quot;SecurityError&quot;],&quot;ignoreErrors&quot;:[&quot;origin message&quot;,&quot;Network request failed&quot;,&quot;Loading chunk&quot;,&quot;这个系统不支持该功能。&quot;,&quot;Can&#x27;t find variable: webkit&quot;,&quot;Can&#x27;t find variable: $&quot;,&quot;内存不足&quot;,&quot;out of memory&quot;,&quot;DOM Exception 18&quot;,&quot;zfeedback sdk 初始化失败！&quot;,&quot;zfeedback sdk 加载失败！&quot;,&quot;The operation is insecure&quot;,&quot;[object Event]&quot;,&quot;[object FileError]&quot;,&quot;[object DOMError]&quot;,&quot;[object Object]&quot;,&quot;拒绝访问。&quot;,&quot;Maximum call stack size exceeded&quot;,&quot;UploadError&quot;,&quot;无法 fetch&quot;,&quot;draft-js&quot;,&quot;缺少 JavaScript 对象&quot;,&quot;componentWillEnter&quot;,&quot;componentWillLeave&quot;,&quot;componentWillAppear&quot;,&quot;getInlineStyleAt&quot;,&quot;getCharacterList&quot;],&quot;whitelistUrls&quot;:[&quot;static.zhihu.com&quot;]}"></script><script nonce="021b92cb-d90a-482e-a6a1-1cbeec6f244d">if (window.requestAnimationFrame) {    window.requestAnimationFrame(function() {      window.FIRST_ANIMATION_FRAME = Date.now();    });  }</script></head><body class="Entry-body"><div id="root"><div data-zop-usertoken="{}" data-reactroot=""><div class="LoadingBar"></div><div><header role="banner" class="Sticky AppHeader" data-za-module="TopNavBar"><div class="AppHeader-inner"><a href="//www.zhihu.com" aria-label="知乎"><svg viewBox="0 0 200 91" class="Icon ZhihuLogo ZhihuLogo--blue Icon--logo" style="height:30px;width:64px" width="64" height="30" aria-hidden="true"><title></title><g><path d="M53.29 80.035l7.32.002 2.41 8.24 13.128-8.24h15.477v-67.98H53.29v67.978zm7.79-60.598h22.756v53.22h-8.73l-8.718 5.473-1.587-5.46-3.72-.012v-53.22zM46.818 43.162h-16.35c.545-8.467.687-16.12.687-22.955h15.987s.615-7.05-2.68-6.97H16.807c1.09-4.1 2.46-8.332 4.1-12.708 0 0-7.523 0-10.085 6.74-1.06 2.78-4.128 13.48-9.592 24.41 1.84-.2 7.927-.37 11.512-6.94.66-1.84.785-2.08 1.605-4.54h9.02c0 3.28-.374 20.9-.526 22.95H6.51c-3.67 0-4.863 7.38-4.863 7.38H22.14C20.765 66.11 13.385 79.24 0 89.62c6.403 1.828 12.784-.29 15.937-3.094 0 0 7.182-6.53 11.12-21.64L43.92 85.18s2.473-8.402-.388-12.496c-2.37-2.788-8.768-10.33-11.496-13.064l-4.57 3.627c1.363-4.368 2.183-8.61 2.46-12.71H49.19s-.027-7.38-2.372-7.38zm128.752-.502c6.51-8.013 14.054-18.302 14.054-18.302s-5.827-4.625-8.556-1.27c-1.874 2.548-11.51 15.063-11.51 15.063l6.012 4.51zm-46.903-18.462c-2.814-2.577-8.096.667-8.096.667s12.35 17.2 12.85 17.953l6.08-4.29s-8.02-11.752-10.83-14.33zM199.99 46.5c-6.18 0-40.908.292-40.953.292v-31.56c1.503 0 3.882-.124 7.14-.376 12.773-.753 21.914-1.25 27.427-1.504 0 0 3.817-8.496-.185-10.45-.96-.37-7.24 1.43-7.24 1.43s-51.63 5.153-72.61 5.64c.5 2.756 2.38 5.336 4.93 6.11 4.16 1.087 7.09.53 15.36.277 7.76-.5 13.65-.76 17.66-.76v31.19h-41.71s.88 6.97 7.97 7.14h33.73v22.16c0 4.364-3.498 6.87-7.65 6.6-4.4.034-8.15-.36-13.027-.566.623 1.24 1.977 4.496 6.035 6.824 3.087 1.502 5.054 2.053 8.13 2.053 9.237 0 14.27-5.4 14.027-14.16V53.93h38.235c3.026 0 2.72-7.432 2.72-7.432z" fill-rule="evenodd"/></g></svg></a><nav role="navigation" class="AppHeader-nav"><a class="AppHeader-navItem" href="//www.zhihu.com/" data-za-not-track-link="true">首页</a><a class="AppHeader-navItem" href="//www.zhihu.com/explore" data-za-not-track-link="true">发现</a><a href="//www.zhihu.com/topic" class="AppHeader-navItem" data-za-not-track-link="true">话题</a></nav><div class="SearchBar" role="search" data-za-module="PresetWordItem"><div class="SearchBar-toolWrapper"><form class="SearchBar-tool"><div><div class="Popover"><div class="SearchBar-input Input-wrapper Input-wrapper--grey"><input type="text" maxLength="100" value="" autoComplete="off" role="combobox" aria-expanded="false" aria-autocomplete="list" aria-activedescendant="null--1" id="null-toggle" aria-haspopup="true" aria-owns="null-content" class="Input" placeholder=""/><div class="Input-after"><button aria-label="搜索" type="button" class="Button SearchBar-searchIcon Button--primary"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Search" fill="currentColor" viewBox="0 0 24 24" width="18" height="18"><path d="M17.068 15.58a8.377 8.377 0 0 0 1.774-5.159 8.421 8.421 0 1 0-8.42 8.421 8.38 8.38 0 0 0 5.158-1.774l3.879 3.88c.957.573 2.131-.464 1.488-1.49l-3.879-3.878zm-6.647 1.157a6.323 6.323 0 0 1-6.316-6.316 6.323 6.323 0 0 1 6.316-6.316 6.323 6.323 0 0 1 6.316 6.316 6.323 6.323 0 0 1-6.316 6.316z" fill-rule="evenodd"></path></svg></span></button></div></div></div></div></form></div></div><div class="AppHeader-userInfo"><div class="AppHeader-profile"><div><button type="button" class="Button AppHeader-login Button--blue">登录</button><button type="button" class="Button Button--primary Button--blue">加入知乎</button></div></div></div></div></header></div><main role="main" class="App-main"><div class="QuestionPage" itemscope="" itemType="http://schema.org/Question"><meta itemProp="name" content="你为什么支持死刑？"/><meta itemProp="url" content="https://www.zhihu.com/question/25084350"/><meta itemProp="keywords" content="法律,死刑,刑法,法学"/><meta itemProp="answerCount" content="12470"/><meta itemProp="commentCount" content="323"/><meta itemProp="dateCreated" content="2014-09-02T08:33:21.000Z"/><meta itemProp="dateModified" content="2018-11-30T09:48:19.000Z"/><meta itemProp="zhihu:visitsCount"/><meta itemProp="zhihu:followerCount" content="44528"/><script type="application/ld+json">
        {
          &quot;@context&quot;: &quot;https://ziyuan.baidu.com/contexts/cambrian.jsonld&quot;,
          &quot;@id&quot;: &quot;https://www.zhihu.com/question/25084350/answer/571315682&quot;,
          &quot;appid&quot;: &quot;否&quot;,
          &quot;pubDate&quot;: &quot;2019-01-11T10:37:06&quot;,
          &quot;upDate&quot;: &quot;2019-01-13T17:23:05&quot;
        }</script><div data-zop-question="{&quot;title&quot;:&quot;你为什么支持死刑？&quot;,&quot;topics&quot;:[{&quot;name&quot;:&quot;法律&quot;,&quot;id&quot;:&quot;19550874&quot;},{&quot;name&quot;:&quot;死刑&quot;,&quot;id&quot;:&quot;19562985&quot;},{&quot;name&quot;:&quot;刑法&quot;,&quot;id&quot;:&quot;19591312&quot;},{&quot;name&quot;:&quot;法学&quot;,&quot;id&quot;:&quot;19604890&quot;}],&quot;id&quot;:25084350,&quot;isEditable&quot;:false}"><div class="QuestionHeader"><div class="QuestionHeader-content"><div class="QuestionHeader-main"><div class="QuestionHeader-tags"><div class="QuestionHeader-topics"><div class="Tag QuestionTopic"><span class="Tag-content"><a class="TopicLink" href="//www.zhihu.com/topic/19550874" target="_blank"><div class="Popover"><div id="null-toggle" aria-haspopup="true" aria-expanded="false" aria-owns="null-content">法律</div></div></a></span></div><div class="Tag QuestionTopic"><span class="Tag-content"><a class="TopicLink" href="//www.zhihu.com/topic/19562985" target="_blank"><div class="Popover"><div id="null-toggle" aria-haspopup="true" aria-expanded="false" aria-owns="null-content">死刑</div></div></a></span></div><div class="Tag QuestionTopic"><span class="Tag-content"><a class="TopicLink" href="//www.zhihu.com/topic/19591312" target="_blank"><div class="Popover"><div id="null-toggle" aria-haspopup="true" aria-expanded="false" aria-owns="null-content">刑法</div></div></a></span></div><div class="Tag QuestionTopic"><span class="Tag-content"><a class="TopicLink" href="//www.zhihu.com/topic/19604890" target="_blank"><div class="Popover"><div id="null-toggle" aria-haspopup="true" aria-expanded="false" aria-owns="null-content">法学</div></div></a></span></div></div></div><h1 class="QuestionHeader-title">你为什么支持死刑？</h1><div><div class="QuestionHeader-detail"><div class="QuestionRichText QuestionRichText--expandable QuestionRichText--collapsed"><div><span class="RichText ztext" itemProp="text">鉴于目前对于死刑的态度两极化，且比例悬殊，在现有问题中很难找出双方各自有价值的回答，特分为两个问题。反对死刑、支持废除死刑的朋友请到问题“<a href="http://www.zhihu.com/question/25084336" class="internal">在中国，你反对死刑的原因是什么？</a>”中回答，谢谢。（已被合併） 回答请尽量避免“理所当然”的口吻，多陈述理由。 相关问题：<a href="http://www.zhihu.com/question/25084336" class="internal">在中国，你反对死刑的原因是什么？ - 调查类问题</a></span><button type="button" class="Button QuestionRichText-more Button--plain">显示全部<svg viewBox="0 0 10 6" class="Icon QuestionRichText-more-icon Icon--arrow" style="height:16px;width:10px" width="10" height="16" aria-hidden="true"><title></title><g><path d="M8.716.217L5.002 4 1.285.218C.99-.072.514-.072.22.218c-.294.29-.294.76 0 1.052l4.25 4.512c.292.29.77.29 1.063 0L9.78 1.27c.293-.29.293-.76 0-1.052-.295-.29-.77-.29-1.063 0z"/></g></svg></button></div></div></div></div></div><div class="QuestionHeader-side"><div class="QuestionHeader-follow-status"><div class="QuestionFollowStatus"><div class="NumberBoard QuestionFollowStatus-counts NumberBoard--divider"><div class="NumberBoard-item"><div class="NumberBoard-itemInner"><div class="NumberBoard-itemName">关注者</div><strong class="NumberBoard-itemValue" title="44528">44,528</strong></div></div><div class="NumberBoard-item"><div class="NumberBoard-itemInner"><div class="NumberBoard-itemName">被浏览</div><strong class="NumberBoard-itemValue" title="70785747">70,785,747</strong></div></div></div></div></div></div></div><div class="QuestionHeader-footer"><div class="QuestionHeader-footer-inner"><div class="QuestionHeader-main QuestionHeader-footer-main"><div class="QuestionButtonGroup"><button type="button" class="Button FollowButton Button--primary Button--blue">关注问题</button><button type="button" class="Button Button--blue"><svg viewBox="0 0 12 12" class="Icon Button-icon Icon--modify" style="height:16px;width:14px" width="14" height="16" aria-hidden="true"><title></title><g><path d="M.423 10.32L0 12l1.667-.474 1.55-.44-2.4-2.33-.394 1.564zM10.153.233c-.327-.318-.85-.31-1.17.018l-.793.817 2.49 2.414.792-.814c.318-.328.312-.852-.017-1.17l-1.3-1.263zM3.84 10.536L1.35 8.122l6.265-6.46 2.49 2.414-6.265 6.46z" fill-rule="evenodd"/></g></svg>写回答</button></div><div class="QuestionHeaderActions"><button style="margin-right:16px" type="button" class="Button Button--grey Button--withIcon Button--withLabel"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Invite Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M4 10V8a1 1 0 1 1 2 0v2h2a1 1 0 0 1 0 2H6v2a1 1 0 0 1-2 0v-2H2a1 1 0 0 1 0-2h2zm10.455 2c-2.21 0-4-1.79-4-4s1.79-4 4-4 4 1.79 4 4-1.79 4-4 4zm-7 6c0-2.66 4.845-4 7.272-4C17.155 14 22 15.34 22 18v1.375c0 .345-.28.625-.625.625H8.08a.625.625 0 0 1-.625-.625V18z" fill-rule="evenodd"></path></svg></span>邀请回答</button><div class="QuestionHeader-Comment"><button type="button" class="Button Button--plain Button--withIcon Button--withLabel"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Comment Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M10.241 19.313a.97.97 0 0 0-.77.2 7.908 7.908 0 0 1-3.772 1.482.409.409 0 0 1-.38-.637 5.825 5.825 0 0 0 1.11-2.237.605.605 0 0 0-.227-.59A7.935 7.935 0 0 1 3 11.25C3 6.7 7.03 3 12 3s9 3.7 9 8.25-4.373 9.108-10.759 8.063z" fill-rule="evenodd"></path></svg></span>323 条评论</button></div><div class="Popover ShareMenu"><div class="ShareMenu-toggler" id="null-toggle" aria-haspopup="true" aria-expanded="false" aria-owns="null-content"><button type="button" class="Button Button--plain Button--withIcon Button--withLabel"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Share Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M2.931 7.89c-1.067.24-1.275 1.669-.318 2.207l5.277 2.908 8.168-4.776c.25-.127.477.198.273.39L9.05 14.66l.927 5.953c.18 1.084 1.593 1.376 2.182.456l9.644-15.242c.584-.892-.212-2.029-1.234-1.796L2.93 7.89z" fill-rule="evenodd"></path></svg></span>分享</button></div></div><div class="Popover"><button aria-label="更多" type="button" id="null-toggle" aria-haspopup="true" aria-expanded="false" aria-owns="null-content" class="Button Button--plain Button--withIcon Button--iconOnly"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Dots Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M5 14a2 2 0 1 1 0-4 2 2 0 0 1 0 4zm7 0a2 2 0 1 1 0-4 2 2 0 0 1 0 4zm7 0a2 2 0 1 1 0-4 2 2 0 0 1 0 4z" fill-rule="evenodd"></path></svg></span></button></div></div><div class="QuestionHeader-actions"></div></div></div></div></div></div><div class="Question-main"><div class="Question-mainColumn" data-zop-questionanswerlist="true"><div class="Card"><a class="QuestionMainAction" data-za-detail-view-element_name="ViewAll" href="/question/25084350">查看全部 12,470 个回答</a></div><div class="Card AnswerCard"><div class="QuestionAnswer-content"><div class="ContentItem AnswerItem" data-zop="{&quot;authorName&quot;:&quot;知乎用户&quot;,&quot;itemId&quot;:571315682,&quot;title&quot;:&quot;你为什么支持死刑？&quot;,&quot;type&quot;:&quot;answer&quot;}" name="571315682" itemProp="mainEntityOfPage" itemType="http://schema.org/Answer" itemscope=""><div class="ContentItem-meta"><div class="AuthorInfo AnswerItem-authorInfo AnswerItem-authorInfo--related" itemProp="author" itemscope="" itemType="http://schema.org/Person"><meta itemProp="name" content="知乎用户"/><meta itemProp="image" content="https://pic4.zhimg.com/da8e974dc_is.jpg"/><meta itemProp="url" content="https://www.zhihu.com/people/"/><meta itemProp="zhihu:followerCount"/><span class="UserLink AuthorInfo-avatarWrapper"><img class="Avatar AuthorInfo-avatar" width="38" height="38" src="https://pic4.zhimg.com/da8e974dc_xs.jpg" srcSet="https://pic4.zhimg.com/da8e974dc_l.jpg 2x" alt="知乎用户"/></span><div class="AuthorInfo-content"><div class="AuthorInfo-head"><span class="UserLink AuthorInfo-name">知乎用户</span></div><div class="AuthorInfo-detail"><div class="AuthorInfo-badge"></div></div></div></div><div class="LabelContainer"></div><div class="AnswerItem-extraInfo"><span class="Voters"><button type="button" class="Button Button--plain">1,320 人<!-- -->赞同了该回答</button></span></div></div><meta itemProp="image" content=""/><meta itemProp="upvoteCount" content="1320"/><meta itemProp="url" content="https://www.zhihu.com/question/25084350/answer/571315682"/><meta itemProp="dateCreated" content="2019-01-11T02:37:06.000Z"/><meta itemProp="dateModified" content="2019-01-13T09:23:05.000Z"/><meta itemProp="commentCount" content="226"/><div class="RichContent RichContent--unescapable"><div class="RichContent-inner"><span class="RichText ztext CopyrightRichText-richText" itemProp="text"><p>我有一个儿时伙伴，美丽的像天使一样，十六岁左右被邻家七十多的老头堵在厕所里给强了。</p><p>那个老头还强过别的女孩，也强过他自己的儿媳。因为这些事，被人打被人揍还要赔偿，他的一个女儿被婆家人看不起，也上吊了。</p><p>这个千刀万剐的老头还是活得好好的，也没有被抓起来。农村不想张扬，赔点钱就算了。可是，我这个伙伴没有熬过这件事，最后心理阴影太重，大年夜喝药了。</p><p>几十年过去了，我依旧不能忘记那个皮肤白皙，眼睛似水，高挑美丽善良的女孩。可怜她就这么没了</p><p>为什么支持死刑？我不仅支持死刑，我还支持千刀万剐呢。一枪没了太幸运，还什么注射，切。</p><p>===============</p><p>因为我的儿时伙伴给我的印象太深，几十年过去我一直对她的离世耿耿于怀，所以，有时候也会问父亲当时的情况，父亲总是不耐烦的回避这个问题。今天在知乎又提起这件事，我还是跟父亲多了一些交谈。结果大吵了一架</p><p>总体的过程是；这个老头害人不是害了一个，从我的伙伴之前的十多年间，害了很多个。第一个竟然是我大姑奶奶家的十几岁小女儿，在果园里看园的时候被这个老头给强了。当时他六十四五岁。有很多留言的说不相信这么大年纪的能有这么大的力气，那是因为你没有在农村待过，壮得跟头牛一样的老头，你以为随便就能一脚撂倒了？</p><p>那个女孩后来生了个孩子，又过了几年就嫁到了别的地方。我生气的是，如果第一个被害人就去报警就制裁，还会有后面那么多的被害人么？我可怜的小伙伴也不会经受这样的痛苦后离开。而我父亲的回答则是事不关己，完全就是与自己无关啊，多丢人的事啊，农村这种事能宣扬吗？从那个老头六十多岁开始犯事，到现在已经快三十多年了。三十多年前我无法去指责那些人的法律意识能强到什么地步了</p><p>更让我气愤的是，人生的不公平。此老头活了八十四五岁才死。明明都知道他不是好人，可是，唯一会做的只是躲着他，不跟他接触而己。我跟父亲对话的时候产生了极度愤慨难过悲愤的心理，无法抑制自己的情绪。</p><p>我执着的询问当时到底对那个老头怎么处理的，父亲说第一次的时候，他跪下来道歉，也被打的够呛，发誓不会再这样了。结果以后还这样，但都是一个村的，怕说出去女儿不好嫁人，也就私下里打一顿，赔点钱了事 ；最后他自己的女儿也不堪忍受父亲这样的丑事，一直被婆家挤兑，加上赔钱的时候他要求自己的几个女儿帮他出钱，人家婆家自然不答应，他自己的女儿自杀了。</p><p>不公平啊，为什么这个老头祸害了这么多女孩，妇女，他竟然能活到八十多岁，哪来的公平？</p><p>听父亲说这老头的儿子，也是到了六十多岁开始有这样的迹象了，六十岁之前父子俩都是好人的。好在现在不比以前，村里人的思想不比过去了</p><p>我还是无法释怀，因为那个儿时伙伴的灿烂笑容一直在我脑海里，几十年都没有抹去，所以那些劝人宽心原谅的，我只能说，那不是你亲身遇到，你无法体会当事人的心理感受，请不要随便开口劝解什么</p><p>所以那些犯了大罪在监狱里改造然后开始什么心理辅导，治疗，什么悔过啊，什么挖掘犯罪分子的内心啊，什么为他们开启新的人生啊，什么让他们从头再来啊。。。拜托，他们这些圣人难道没看到有很多出了监狱就接着犯罪的吗？管个球用啊？气死了。还联系受害人达成谅解啊，还要去想办法关注关心劝解罪犯啊。。。你妹啊，死的人呢？受过伤害的人呢？人家的家属呢？人家的人生呢？我去。。。气死了。</p></span></div><div><div class="ContentItem-time"><a target="_blank" href="/question/25084350/answer/571315682"><span data-tooltip="发布于 2019-01-11 10:37">编辑于昨天 17:23</span></a></div></div><div class="ContentItem-actions RichContent-actions"><span><button aria-label="赞同" type="button" class="Button VoteButton VoteButton--up"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--TriangleUp VoteButton-TriangleUp" fill="currentColor" viewBox="0 0 24 24" width="10" height="10"><path d="M2 18.242c0-.326.088-.532.237-.896l7.98-13.203C10.572 3.57 11.086 3 12 3c.915 0 1.429.571 1.784 1.143l7.98 13.203c.15.364.236.57.236.896 0 1.386-.875 1.9-1.955 1.9H3.955c-1.08 0-1.955-.517-1.955-1.9z" fill-rule="evenodd"></path></svg></span>赞同 <!-- -->1.3K</button><button aria-label="反对" type="button" class="Button VoteButton VoteButton--down"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--TriangleDown" fill="currentColor" viewBox="0 0 24 24" width="10" height="10"><path d="M20.044 3H3.956C2.876 3 2 3.517 2 4.9c0 .326.087.533.236.896L10.216 19c.355.571.87 1.143 1.784 1.143s1.429-.572 1.784-1.143l7.98-13.204c.149-.363.236-.57.236-.896 0-1.386-.876-1.9-1.956-1.9z" fill-rule="evenodd"></path></svg></span></button></span><button type="button" class="Button ContentItem-action Button--plain Button--withIcon Button--withLabel"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Comment Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M10.241 19.313a.97.97 0 0 0-.77.2 7.908 7.908 0 0 1-3.772 1.482.409.409 0 0 1-.38-.637 5.825 5.825 0 0 0 1.11-2.237.605.605 0 0 0-.227-.59A7.935 7.935 0 0 1 3 11.25C3 6.7 7.03 3 12 3s9 3.7 9 8.25-4.373 9.108-10.759 8.063z" fill-rule="evenodd"></path></svg></span>226 条评论</button><div class="Popover ShareMenu ContentItem-action"><div class="ShareMenu-toggler" id="null-toggle" aria-haspopup="true" aria-expanded="false" aria-owns="null-content"><button type="button" class="Button Button--plain Button--withIcon Button--withLabel"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Share Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M2.931 7.89c-1.067.24-1.275 1.669-.318 2.207l5.277 2.908 8.168-4.776c.25-.127.477.198.273.39L9.05 14.66l.927 5.953c.18 1.084 1.593 1.376 2.182.456l9.644-15.242c.584-.892-.212-2.029-1.234-1.796L2.93 7.89z" fill-rule="evenodd"></path></svg></span>分享</button></div></div><button type="button" class="Button ContentItem-action Button--plain Button--withIcon Button--withLabel"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Star Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M5.515 19.64l.918-5.355-3.89-3.792c-.926-.902-.639-1.784.64-1.97L8.56 7.74l2.404-4.871c.572-1.16 1.5-1.16 2.072 0L15.44 7.74l5.377.782c1.28.186 1.566 1.068.64 1.97l-3.89 3.793.918 5.354c.219 1.274-.532 1.82-1.676 1.218L12 18.33l-4.808 2.528c-1.145.602-1.896.056-1.677-1.218z" fill-rule="evenodd"></path></svg></span>收藏</button><button type="button" class="Button ContentItem-action Button--plain Button--withIcon Button--withLabel"><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--Heart Button-zi" fill="currentColor" viewBox="0 0 24 24" width="1.2em" height="1.2em"><path d="M2 8.437C2 5.505 4.294 3.094 7.207 3 9.243 3 11.092 4.19 12 6c.823-1.758 2.649-3 4.651-3C19.545 3 22 5.507 22 8.432 22 16.24 13.842 21 12 21 10.158 21 2 16.24 2 8.437z" fill-rule="evenodd"></path></svg></span>感谢</button><button data-zop-retract-question="true" type="button" class="Button ContentItem-action ContentItem-rightButton Button--plain"><span class="RichContent-collapsedText">收起</span><span style="display:inline-flex;align-items:center">​<svg class="Zi Zi--ArrowDown ContentItem-arrowIcon is-active" fill="currentColor" viewBox="0 0 24 24" width="24" height="24"><path d="M12 13L8.285 9.218a.758.758 0 0 0-1.064 0 .738.738 0 0 0 0 1.052l4.249 4.512a.758.758 0 0 0 1.064 0l4.246-4.512a.738.738 0 0 0 0-1.052.757.757 0 0 0-1.063 0L12.002 13z" fill-rule="evenodd"></path></svg></span></button></div></div></div></div></div><div class="Card"><a class="QuestionMainAction" data-za-detail-view-element_name="ViewAll" href="/question/25084350">查看全部 12,470 个回答</a></div></div></div></div></main></div></div><script id="js-clientConfig" type="text/json">{"host":"zhihu.com","protocol":"https:","wwwHost":"www.zhihu.com","zhuanlanHost":"zhuanlan.zhihu.com"}</script><script id="js-initialData" type="text/json">{"initialState":{"common":{"ask":{}},"privacy":{"showPrivacy":false},"loading":{"global":{"count":0},"local":{"question\u002Fget\u002F":false,"answer\u002Fget\u002F":false}},"entities":{"users":{},"questions":{"25084350":{"type":"question","id":25084350,"title":"你为什么支持死刑？","questionType":"normal","created":1409646801,"updatedTime":1543571299,"url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Fquestions\u002F25084350","isMuted":false,"isNormal":true,"isEditable":false,"adminClosedComment":false,"hasPublishingDraft":false,"answerCount":12470,"visitCount":70785747,"commentCount":323,"followerCount":44528,"collapsedAnswerCount":664,"excerpt":"鉴于目前对于死刑的态度两极化，且比例悬殊，在现有问题中很难找出双方各自有价值的回答，特分为两个问题。反对死刑、支持废除死刑的朋友请到问题“\u003Ca href=\"http:\u002F\u002Fwww.zhihu.com\u002Fquestion\u002F25084336\" class=\"internal\"\u003E在中国，你反对死刑的原因是什么？\u003C\u002Fa\u003E”中回答，谢谢。（已被合併） 回答请尽量避免“理所当然”的口吻，多陈述理由。 相关问题：\u003Ca href=\"http:\u002F\u002Fwww.zhihu.com\u002Fquestion\u002F25084336\" class=\"internal\"\u003E在中国，你反对死刑的原因是什么？ - 调查类问题\u003C\u002Fa\u003E","commentPermission":"all","detail":"鉴于目前对于死刑的态度两极化，且比例悬殊，在现有问题中很难找出双方各自有价值的回答，特分为两个问题。反对死刑、支持废除死刑的朋友请到问题“\u003Ca href=\"http:\u002F\u002Fwww.zhihu.com\u002Fquestion\u002F25084336\" class=\"internal\"\u003E在中国，你反对死刑的原因是什么？\u003C\u002Fa\u003E”中回答，谢谢。（已被合併）\u003Cbr\u003E\u003Cbr\u003E回答请尽量避免“理所当然”的口吻，多陈述理由。\u003Cbr\u003E相关问题：\u003Ca href=\"http:\u002F\u002Fwww.zhihu.com\u002Fquestion\u002F25084336\" class=\"internal\"\u003E在中国，你反对死刑的原因是什么？ - 调查类问题\u003C\u002Fa\u003E","editableDetail":"鉴于目前对于死刑的态度两极化，且比例悬殊，在现有问题中很难找出双方各自有价值的回答，特分为两个问题。反对死刑、支持废除死刑的朋友请到问题“\u003Ca href=\"http:\u002F\u002Fwww.zhihu.com\u002Fquestion\u002F25084336\" class=\"internal\"\u003E在中国，你反对死刑的原因是什么？\u003C\u002Fa\u003E”中回答，谢谢。（已被合併）\u003Cbr\u003E\u003Cbr\u003E回答请尽量避免“理所当然”的口吻，多陈述理由。\u003Cbr\u003E相关问题：\u003Ca href=\"http:\u002F\u002Fwww.zhihu.com\u002Fquestion\u002F25084336\" class=\"internal\"\u003E在中国，你反对死刑的原因是什么？ - 调查类问题\u003C\u002Fa\u003E","status":{"isLocked":false,"isClose":false,"isEvaluate":false,"isSuggest":false},"relationship":{"isAuthor":false,"isFollowing":false,"isAnonymous":false,"canLock":false,"canStickAnswers":false,"canCollapseAnswers":false},"topics":[{"id":"19550874","type":"topic","url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Ftopics\u002F19550874","name":"法律","avatarUrl":"https:\u002F\u002Fpic2.zhimg.com\u002Ff126e096216e4554289b0996539b79b5_is.jpg","excerpt":"法律，是国家的产物，是指统治阶级（泛指政治、经济、思想形态上占支配地位的阶级），为了实现统治并管理国家的目的，经过一定立法程序，所颁布的基本法律和普通法律。法律是统治阶级意志的体现，国家的统治工具。法律是由享有立法权的立法机关（全国人民代表大会和全国人民代表大会常务委员会）行使国家立法权，依照法定程序制定、修改并颁布，并由国家强制力保证实施的基本法律和普通法律总称。包括基本法律、普通法律。法，可…","introduction":"法律，是国家的产物，是指统治阶级（泛指政治、经济、思想形态上占支配地位的阶级），为了实现统治并管理国家的目的，经过一定立法程序，所颁布的基本法律和普通法律。法律是统治阶级意志的体现，国家的统治工具。法律是由享有立法权的立法机关（全国人民代表大会和全国人民代表大会常务委员会）行使国家立法权，依照法定程序制定、修改并颁布，并由国家强制力保证实施的基本法律和普通法律总称。包括基本法律、普通法律。法，可划分为1、宪法，2、法律，3、行政法规，4、地方性法规，5、自治条例和单行条例。宪法是高于其它法律部门（法律、行政法规、地方性法规、自治条例和单行条例）的国家根本大法，它规定国家制度和社会制度最基本的原则，公民基本权利和义务，国家机构的组织及其活动的原则等。法律是从属于宪法的强制性规范，是宪法的具体化。宪法是国家法的基础与核心，法律则是国家法的重要组成部分。法律可划分为基本法律（如刑法、刑事诉讼法、民法通则、民事诉讼法、行政诉讼法、行政法、商法、国际法等）和普通法律（如商标法、文物保护法等）。行政法规，是国家行政机关（国务院）根据宪法和法律，制定的行政规范的总称。"},{"id":"19562985","type":"topic","url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Ftopics\u002F19562985","name":"死刑","avatarUrl":"https:\u002F\u002Fpic4.zhimg.com\u002F9434c6a11_is.jpg","excerpt":"死刑，也称为极刑、处决、生命刑，世界上最古老的刑罚之一，指行刑者基于法律所赋予的权力，结束一个犯人的生命。遭受这种剥夺生命的刑罚方法的有关犯人通常都在当地犯了严重罪行。尽管这“严重罪行”的定义时常有争议，但在现时保有死刑的国家中，一般来说，“谋杀”必然是犯人被判死刑的其中一个重要理由。《刑法修正案》9中：执行死刑条件由如果故意犯罪，查证属实，等发生重大犯罪，手段极其残忍，社会影响极其恶劣，如杀人…","introduction":"死刑，也称为极刑、处决、生命刑，世界上最古老的刑罚之一，指行刑者基于法律所赋予的权力，结束一个犯人的生命。遭受这种剥夺生命的刑罚方法的有关犯人通常都在当地犯了严重罪行。尽管这“严重罪行”的定义时常有争议，但在现时保有死刑的国家中，一般来说，“谋杀”必然是犯人被判死刑的其中一个重要理由。《刑法修正案》9中：执行死刑条件由如果故意犯罪，查证属实，等发生重大犯罪，手段极其残忍，社会影响极其恶劣，如杀人，勒索绑架，抢劫，强奸，等危害国家刑法都有可能会执行死刑。"},{"id":"19591312","type":"topic","url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Ftopics\u002F19591312","name":"刑法","avatarUrl":"https:\u002F\u002Fpic1.zhimg.com\u002Fv2-5bc811c2495ed9439b44d1cece595e48_is.jpg","excerpt":"刑法是规定犯罪、刑事责任和刑罚的法律，是掌握政权的统治阶级为了维护本阶级政治上的统治和各阶级经济上的利益，根据自己的意志，规定哪些行为是犯罪并且应当负何种刑事责任 ，并给予犯罪嫌疑人何种刑事处罚的法律规范的总称。刑法有广义与狭义之分。广义刑法是一切刑事法律规范的总称，狭义刑法仅指刑法典，在我国即《中华人民共和国刑法》。与广义刑法、狭义刑法相联系的，刑法还可区分为普通刑法和特别刑法。普通刑法指具有…","introduction":"刑法是规定犯罪、刑事责任和刑罚的法律，是掌握政权的统治阶级为了维护本阶级政治上的统治和各阶级经济上的利益，根据自己的意志，规定哪些行为是犯罪并且应当负何种刑事责任 ，并给予犯罪嫌疑人何种刑事处罚的法律规范的总称。刑法有广义与狭义之分。广义刑法是一切刑事法律规范的总称，狭义刑法仅指刑法典，在我国即《中华人民共和国刑法》。与广义刑法、狭义刑法相联系的，刑法还可区分为普通刑法和特别刑法。普通刑法指具有普遍使用效力的刑法，实际上即指刑法典。特别刑法指仅使用于特定的人、时、地、事（犯罪）的刑法。在我国，也叫单行刑法和附属刑法。2015年8月29日，十二届全国人大常委会十六次会议表决通过刑法修正案（九）。修改后的刑法自2015年11月1日开始施行。这也是继1997年全面修订刑法后通过的第九个刑法修正案。"},{"id":"19604890","type":"topic","url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Ftopics\u002F19604890","name":"法学","avatarUrl":"https:\u002F\u002Fpic1.zhimg.com\u002F4551f0a33_is.jpg","excerpt":"法学，又称法律学、法律科学，是以法律、法律现象以及其规律性为研究内容的科学，它是研究与法相关问题的专门学问，是关于法律问题的知识和理论体系。法学，是关于法律的科学。法律作为社会的强制性规范，其直接目的在于维持社会秩序，并通过秩序的构建与维护，实现社会公正。作为以法律为研究对象的法学，其核心就在对于秩序与公正的研究，是秩序与公正之学。法学是世界各国高等学校普遍开设的大类，也是中国大学的十大学科体系…","introduction":"法学，又称法律学、法律科学，是以法律、法律现象以及其规律性为研究内容的科学，它是研究与法相关问题的专门学问，是关于法律问题的知识和理论体系。法学，是关于法律的科学。法律作为社会的强制性规范，其直接目的在于维持社会秩序，并通过秩序的构建与维护，实现社会公正。作为以法律为研究对象的法学，其核心就在对于秩序与公正的研究，是秩序与公正之学。法学是世界各国高等学校普遍开设的大类，也是中国大学的十大学科体系之一，包括法学、政治学、公安学、社会学四个主要组成部分。在社会上，很多人习惯将法学专业称之为法律专业。在中国，法学思想最早源于春秋战国时期的法家哲学思想，法学一词，在中国先秦时代被称为“刑名之学”，从汉代开始有“律学”的名称。在西方，古罗马法学家乌尔比安（Ulpianus）对“法学”（古代拉丁语中的Jurisprudentia）一词的定义是：人和神的事务的概念，正义和非正义之学。"}],"author":{"id":"6654eb10f1ccce9bb6d3d202db8b8d9e","urlToken":"walter-white-83","name":"Walter White","avatarUrl":"https:\u002F\u002Fpic4.zhimg.com\u002F3c52e17c4_is.jpg","avatarUrlTemplate":"https:\u002F\u002Fpic4.zhimg.com\u002F3c52e17c4_{size}.jpg","isOrg":false,"type":"people","url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Fpeople\u002F6654eb10f1ccce9bb6d3d202db8b8d9e","userType":"people","headline":"情绪不太稳定，以讲话难听著称","badge":[],"gender":1,"isAdvertiser":false,"isPrivacy":false},"canComment":{"status":true,"reason":""},"reviewInfo":{"type":"","tips":"","editTips":"","isReviewing":false},"relatedCards":[],"muteInfo":{"type":""}}},"answers":{"571315682":{"id":571315682,"type":"answer","answerType":"normal","question":{"type":"question","id":25084350,"title":"你为什么支持死刑？","questionType":"normal","created":1409646801,"updatedTime":1543571299,"url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Fquestions\u002F25084350","relationship":{}},"author":{"id":"ac1bbf41913ac36f2e4428390821bb14","urlToken":"","name":"知乎用户","avatarUrl":"https:\u002F\u002Fpic4.zhimg.com\u002Fda8e974dc_is.jpg","avatarUrlTemplate":"https:\u002F\u002Fpic4.zhimg.com\u002Fda8e974dc_{size}.jpg","isOrg":false,"type":"people","url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Fpeople\u002F0","userType":"people","headline":"","badge":[],"gender":-1,"isAdvertiser":false,"isFollowed":false,"isPrivacy":true},"url":"https:\u002F\u002Fwww.zhihu.com\u002Fapi\u002Fv4\u002Fanswers\u002F571315682","thumbnail":"","isCollapsed":false,"createdTime":1547174226,"updatedTime":1547371385,"extras":"","isCopyable":true,"isNormal":true,"voteupCount":1320,"commentCount":226,"isSticky":false,"adminClosedComment":false,"commentPermission":"all","canComment":{"reason":"","status":true},"reshipmentSettings":"allowed","content":"\u003Cp\u003E我有一个儿时伙伴，美丽的像天使一样，十六岁左右被邻家七十多的老头堵在厕所里给强了。\u003C\u002Fp\u003E\u003Cp\u003E那个老头还强过别的女孩，也强过他自己的儿媳。因为这些事，被人打被人揍还要赔偿，他的一个女儿被婆家人看不起，也上吊了。\u003C\u002Fp\u003E\u003Cp\u003E这个千刀万剐的老头还是活得好好的，也没有被抓起来。农村不想张扬，赔点钱就算了。可是，我这个伙伴没有熬过这件事，最后心理阴影太重，大年夜喝药了。\u003C\u002Fp\u003E\u003Cp\u003E几十年过去了，我依旧不能忘记那个皮肤白皙，眼睛似水，高挑美丽善良的女孩。可怜她就这么没了\u003C\u002Fp\u003E\u003Cp\u003E为什么支持死刑？我不仅支持死刑，我还支持千刀万剐呢。一枪没了太幸运，还什么注射，切。\u003C\u002Fp\u003E\u003Cp\u003E===============\u003C\u002Fp\u003E\u003Cp\u003E因为我的儿时伙伴给我的印象太深，几十年过去我一直对她的离世耿耿于怀，所以，有时候也会问父亲当时的情况，父亲总是不耐烦的回避这个问题。今天在知乎又提起这件事，我还是跟父亲多了一些交谈。结果大吵了一架\u003C\u002Fp\u003E\u003Cp\u003E总体的过程是；这个老头害人不是害了一个，从我的伙伴之前的十多年间，害了很多个。第一个竟然是我大姑奶奶家的十几岁小女儿，在果园里看园的时候被这个老头给强了。当时他六十四五岁。有很多留言的说不相信这么大年纪的能有这么大的力气，那是因为你没有在农村待过，壮得跟头牛一样的老头，你以为随便就能一脚撂倒了？\u003C\u002Fp\u003E\u003Cp\u003E那个女孩后来生了个孩子，又过了几年就嫁到了别的地方。我生气的是，如果第一个被害人就去报警就制裁，还会有后面那么多的被害人么？我可怜的小伙伴也不会经受这样的痛苦后离开。而我父亲的回答则是事不关己，完全就是与自己无关啊，多丢人的事啊，农村这种事能宣扬吗？从那个老头六十多岁开始犯事，到现在已经快三十多年了。三十多年前我无法去指责那些人的法律意识能强到什么地步了\u003C\u002Fp\u003E\u003Cp\u003E更让我气愤的是，人生的不公平。此老头活了八十四五岁才死。明明都知道他不是好人，可是，唯一会做的只是躲着他，不跟他接触而己。我跟父亲对话的时候产生了极度愤慨难过悲愤的心理，无法抑制自己的情绪。\u003C\u002Fp\u003E\u003Cp\u003E我执着的询问当时到底对那个老头怎么处理的，父亲说第一次的时候，他跪下来道歉，也被打的够呛，发誓不会再这样了。结果以后还这样，但都是一个村的，怕说出去女儿不好嫁人，也就私下里打一顿，赔点钱了事 ；最后他自己的女儿也不堪忍受父亲这样的丑事，一直被婆家挤兑，加上赔钱的时候他要求自己的几个女儿帮他出钱，人家婆家自然不答应，他自己的女儿自杀了。\u003C\u002Fp\u003E\u003Cp\u003E不公平啊，为什么这个老头祸害了这么多女孩，妇女，他竟然能活到八十多岁，哪来的公平？\u003C\u002Fp\u003E\u003Cp\u003E听父亲说这老头的儿子，也是到了六十多岁开始有这样的迹象了，六十岁之前父子俩都是好人的。好在现在不比以前，村里人的思想不比过去了\u003C\u002Fp\u003E\u003Cp\u003E我还是无法释怀，因为那个儿时伙伴的灿烂笑容一直在我脑海里，几十年都没有抹去，所以那些劝人宽心原谅的，我只能说，那不是你亲身遇到，你无法体会当事人的心理感受，请不要随便开口劝解什么\u003C\u002Fp\u003E\u003Cp\u003E所以那些犯了大罪在监狱里改造然后开始什么心理辅导，治疗，什么悔过啊，什么挖掘犯罪分子的内心啊，什么为他们开启新的人生啊，什么让他们从头再来啊。。。拜托，他们这些圣人难道没看到有很多出了监狱就接着犯罪的吗？管个球用啊？气死了。还联系受害人达成谅解啊，还要去想办法关注关心劝解罪犯啊。。。你妹啊，死的人呢？受过伤害的人呢？人家的家属呢？人家的人生呢？我去。。。气死了。\u003C\u002Fp\u003E","editableContent":"","excerpt":"我有一个儿时伙伴，美丽的像天使一样，十六岁左右被邻家七十多的老头堵在厕所里给强了。那个老头还强过别的女孩，也强过他自己的儿媳。因为这些事，被人打被人揍还要赔偿，他的一个女儿被婆家人看不起，也上吊了。这个千刀万剐的老头还是活得好好的，也没有…","collapsedBy":"nobody","collapseReason":"","annotationAction":[],"markInfos":[],"relevantInfo":{"isRelevant":false,"relevantType":"","relevantText":""},"suggestEdit":{"reason":"","status":false,"tip":"","title":"","unnormalDetails":{"status":"","description":"","reason":"","reasonId":0,"note":""},"url":""},"isLabeled":false,"rewardInfo":{"canOpenReward":false,"isRewardable":false,"rewardMemberCount":0,"rewardTotalMoney":0,"tagline":""},"relationship":{"isAuthor":false,"isAuthorized":false,"isNothelp":false,"isThanked":false,"voting":0,"upvotedFollowees":[]}}},"articles":{},"columns":{},"topics":{},"roundtables":{},"favlists":{},"comments":{},"notifications":{},"ebooks":{},"activities":{},"feeds":{},"pins":{},"promotions":{},"drafts":{}},"currentUser":"","account":{"lockLevel":{},"unlockTicketStatus":false,"unlockTicket":null,"challenge":[],"errorStatus":false,"message":"","isFetching":false,"accountInfo":{},"urlToken":{"loading":false}},"settings":{"socialBind":null,"inboxMsg":null,"notification":{},"privacyFlag":null,"blockedUsers":{"isFetching":false,"paging":{"pageNo":1,"pageSize":6},"data":[]},"blockedFollowees":{"isFetching":false,"paging":{"pageNo":1,"pageSize":6},"data":[]},"ignoredTopics":{"isFetching":false,"paging":{"pageNo":1,"pageSize":6},"data":[]},"restrictedTopics":null,"laboratory":{}},"notification":{},"people":{"profileStatus":{},"activitiesByUser":{},"answersByUser":{},"answersSortByVotesByUser":{},"answersIncludedByUser":{},"votedAnswersByUser":{},"thankedAnswersByUser":{},"voteAnswersByUser":{},"thankAnswersByUser":{},"topicAnswersByUser":{},"articlesByUser":{},"articlesSortByVotesByUser":{},"articlesIncludedByUser":{},"pinsByUser":{},"questionsByUser":{},"commercialQuestionsByUser":{},"favlistsByUser":{},"followingByUser":{},"followersByUser":{},"mutualsByUser":{},"followingColumnsByUser":{},"followingQuestionsByUser":{},"followingFavlistsByUser":{},"followingTopicsByUser":{},"publicationsByUser":{},"columnsByUser":{},"allFavlistsByUser":{},"brands":null,"creationsByUser":{},"creationsSortByVotesByUser":{}},"env":{"ab":{"config":{"experiments":[{"expId":"launch-ad_ios_lans-2","expPrefix":"ad_ios_lans","isDynamicallyUpdated":true,"isRuntime":true,"includeTriggerInfo":false},{"expId":"launch-ad_uiweb_js-2","expPrefix":"ad_uiweb_js","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-ad_uiweb_open-2","expPrefix":"ad_uiweb_open","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-ad_web_js_track2-2","expPrefix":"ad_web_js_track2","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_adr_dkts-11","expPrefix":"gw_adr_dkts","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_adr_wbtp-2","expPrefix":"gw_adr_wbtp","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_adr_wxfb-2","expPrefix":"gw_adr_wxfb","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_anr_wxbk-2","expPrefix":"gw_anr_wxbk","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_ios_dkts-8","expPrefix":"gw_ios_dkts","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_ios_tk_d-2","expPrefix":"gw_ios_tk_d","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_ios_wxb-2","expPrefix":"gw_ios_wxb","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_ios_wxfb-1","expPrefix":"gw_ios_wxfb","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-gw_wbtp-2","expPrefix":"gw_wbtp","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-mp_amap_ios-1","expPrefix":"mp_amap_ios","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-mp_apm-1","expPrefix":"mp_apm","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-mp_hb_si-3","expPrefix":"mp_hb_si","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-mp_httpdns_ios-4","expPrefix":"mp_httpdns_ios","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-mp_ios_bvc-2","expPrefix":"mp_ios_bvc","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-mp_ios_webp-2","expPrefix":"mp_ios_webp","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-mp_video_feed-2","expPrefix":"mp_video_feed","isDynamicallyUpdated":false,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-se_ios_topsearch-2","expPrefix":"se_ios_topsearch","isDynamicallyUpdated":false,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-se_searchbox-4","expPrefix":"se_searchbox","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-se_search_tab-2","expPrefix":"se_search_tab","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-top_feed_card-1","expPrefix":"top_feed_card","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-top_rfsh_all-2","expPrefix":"top_rfsh_all","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-us_mobile_login-2","expPrefix":"us_mobile_login","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-us_telecom_login-2","expPrefix":"us_telecom_login","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-vd_adrupload_cdn-2","expPrefix":"vd_adrupload_cdn","isDynamicallyUpdated":false,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-vd_ppt_enter_2-2","expPrefix":"vd_ppt_enter_2","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-vd_upload_cdn-2","expPrefix":"vd_upload_cdn","isDynamicallyUpdated":false,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-vd_video_agent-2","expPrefix":"vd_video_agent","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-vd_v_upload_core-2","expPrefix":"vd_v_upload_core","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"launch-vd_zm_core-2","expPrefix":"vd_zm_core","isDynamicallyUpdated":false,"isRuntime":false,"includeTriggerInfo":false},{"expId":"vd_challenge_p-3","expPrefix":"vd_challenge_p","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"ad_lto_new-2","expPrefix":"ad_lto_new","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"top_rfsh-1","expPrefix":"top_rfsh","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false},{"expId":"top_v_album-2","expPrefix":"top_v_album","isDynamicallyUpdated":true,"isRuntime":false,"includeTriggerInfo":false}],"params":[{"id":"top_recall_tb_long","type":"String","value":"51","chainId":"_all_"},{"id":"top_wonderful","type":"String","value":"1","chainId":"_all_"},{"id":"se_filter","type":"String","value":"0","chainId":"_all_"},{"id":"top_rerank_reformat","type":"String","value":"-1","chainId":"_all_"},{"id":"adr_guest_refresh","type":"String","value":"1"},{"id":"hb_billboard","type":"String","value":"0"},{"id":"ios_hybrid_editor_v3","type":"String","value":"0"},{"id":"ios_video_feed","type":"String","value":"1"},{"id":"se_mfq","type":"String","value":"0","chainId":"_all_"},{"id":"top_billpic","type":"String","value":"0","chainId":"_all_"},{"id":"top_test_4_liguangyi","type":"String","value":"1","chainId":"_all_"},{"id":"adr_traffic_leak","type":"String","value":"false"},{"id":"ios_launch_timeout","type":"String","value":"-1"},{"id":"ios_question_answer_preload","type":"String","value":"0"},{"id":"li_gbdt","type":"String","value":"default","chainId":"_all_"},{"id":"tp_discussion_feed_type_android","type":"String","value":"0","chainId":"_all_"},{"id":"adr_sqtc","type":"String","value":"1"},{"id":"adr_video_upload_cdn","type":"String","value":"1"},{"id":"ios_cilssq","type":"String","value":"1"},{"id":"se_time_search","type":"String","value":"origin","chainId":"_all_"},{"id":"pf_newguide_vertical","type":"String","value":"0","chainId":"_all_"},{"id":"top_question_ask","type":"String","value":"1","chainId":"_all_"},{"id":"adr_ad_local_js_sw","type":"String","value":" open"},{"id":"adr_pdf","type":"String","value":"n"},{"id":"adr_ques_comment","type":"String","value":"0"},{"id":"adr_rlp","type":"String","value":"0"},{"id":"top_newuser_feed","type":"String","value":"0","chainId":"_all_"},{"id":"adr_challenge_plan","type":"String","value":"2"},{"id":"hb_stream_render","type":"String","value":"0"},{"id":"se_consulting_switch","type":"String","value":"off","chainId":"_all_"},{"id":"top_core_session","type":"String","value":"-1","chainId":"_all_"},{"id":"adr_guest_login","type":"String","value":"0"},{"id":"top_universalebook","type":"String","value":"1","chainId":"_all_"},{"id":"ios_x_z_c0","type":"String","value":"1"},{"id":"se_websearch","type":"String","value":"0","chainId":"_all_"},{"id":"top_raf","type":"String","value":"y","chainId":"_all_"},{"id":"top_recall_tb","type":"String","value":"1","chainId":"_all_"},{"id":"adr_cta","type":"String","value":"0"},{"id":"adr_laws_length","type":"String","value":"-1"},{"id":"adr_mobile_login","type":"String","value":"0"},{"id":"adr_new_roundtable","type":"String","value":"true"},{"id":"tp_sft","type":"String","value":"a","chainId":"_all_"},{"id":"zr_video_rec","type":"String","value":"zr_video_rec:base","chainId":"_all_"},{"id":"soc_zero_follow","type":"String","value":"0","chainId":"_all_"},{"id":"top_recall_core_interest","type":"String","value":"81","chainId":"_all_"},{"id":"adr_dsa","type":"String","value":"0"},{"id":"adr_medal","type":"String","value":"0"},{"id":"adr_ps","type":"String","value":"0"},{"id":"ios_video_continuous","type":"String","value":"0"},{"id":"top_recall_tb_follow","type":"String","value":"71","chainId":"_all_"},{"id":"adr_android_video_continuous","type":"String","value":"0"},{"id":"hb_consulting_price","type":"String","value":"np"},{"id":"ios_consultation","type":"String","value":"0"},{"id":"ios_spic","type":"String","value":"0"},{"id":"se_daxuechuisou","type":"String","value":"new","chainId":"_all_"},{"id":"hb_search_video_icon","type":"String","value":"0"},{"id":"ios_ad_web_js_track","type":"String","value":"1"},{"id":"ios_q_o_b","type":"String","value":"0"},{"id":"ls_new_video","type":"String","value":"0","chainId":"_all_"},{"id":"ios_magitab","type":"String","value":"0"},{"id":"pin_efs","type":"String","value":"orig","chainId":"_all_"},{"id":"top_accm_ab","type":"String","value":"1","chainId":"_all_"},{"id":"top_followtop","type":"String","value":"1","chainId":"_all_"},{"id":"adr_book_is_card","type":"String","value":"0"},{"id":"adr_mqtt_5_24_0","type":"String","value":"0"},{"id":"adr_new_hybrid","type":"String","value":"0"},{"id":"ios_httpdns","type":"String","value":"zhihu"},{"id":"top_yhgc","type":"String","value":"0","chainId":"_all_"},{"id":"top_no_weighing","type":"String","value":"1","chainId":"_all_"},{"id":"top_search_query","type":"String","value":"0","chainId":"_all_"},{"id":"tp_qa_metacard_top","type":"String","value":"0","chainId":"_all_"},{"id":"ios_dns_hyb_300002","type":"String","value":"0"},{"id":"ios_ps","type":"String","value":"0"},{"id":"qa_web_answerlist_ad","type":"String","value":"0","chainId":"_all_"},{"id":"top_brand","type":"String","value":"1","chainId":"_all_"},{"id":"ios_webp","type":"String","value":"1"},{"id":"ios_ydyq","type":"String","value":"X"},{"id":"se_major_onebox","type":"String","value":"major","chainId":"_all_"},{"id":"top_rank","type":"String","value":"0","chainId":"_all_"},{"id":"adr_mbv","type":"String","value":"false"},{"id":"adr_use_cashier","type":"String","value":"new"},{"id":"hb_new_upvote","type":"String","value":"online_upvote"},{"id":"ios_comment","type":"String","value":"0"},{"id":"top_tagextend","type":"String","value":"1","chainId":"_all_"},{"id":"top_video_score","type":"String","value":"1","chainId":"_all_"},{"id":"tp_discussion_feed_card_type","type":"String","value":"0","chainId":"_all_"},{"id":"tp_answer_meta_guide","type":"String","value":"0","chainId":"_all_"},{"id":"hb_follow_guide_v2","type":"String","value":"X"},{"id":"ios_cashier_color","type":"String","value":"1"},{"id":"se_click2","type":"String","value":"0","chainId":"_all_"},{"id":"top_promo","type":"String","value":"1","chainId":"_all_"},{"id":"ios_add_question_v2","type":"String","value":"0"},{"id":"top_billvideo","type":"String","value":"0","chainId":"_all_"},{"id":"top_hotlist","type":"String","value":"1","chainId":"_all_"},{"id":"web_heifetz_grow_ad","type":"String","value":"1"},{"id":"adr_question_editor","type":"String","value":"0"},{"id":"adr_use_gd","type":"String","value":"n"},{"id":"ios_searchbox","type":"String","value":"2"},{"id":"top_billab","type":"String","value":"0","chainId":"_all_"},{"id":"ios_pay_view","type":"String","value":"new"},{"id":"se_km_ad_locate","type":"String","value":"1","chainId":"_all_"},{"id":"adr_android_launch_ad_mp4","type":"String","value":"open"},{"id":"adr_telecom_login","type":"String","value":"1"},{"id":"hb_unfollow_reason","type":"String","value":"0"},{"id":"ios_ad_web_cache","type":"String","value":"0"},{"id":"top_recall_follow_user","type":"String","value":"91","chainId":"_all_"},{"id":"top_sj","type":"String","value":"2","chainId":"_all_"},{"id":"top_user_gift","type":"String","value":"0","chainId":"_all_"},{"id":"web_answerlist_ad","type":"String","value":"0"},{"id":"adr_dkts","type":"String","value":"20"},{"id":"adr_mqtt","type":"String","value":"0"},{"id":"adr_regis_avatar","type":"String","value":"1"},{"id":"ios_article_misc_panel","type":"String","value":"0"},{"id":"web_stream_render","type":"String","value":"0"},{"id":"web_top_hp_thanks","type":"String","value":"0"},{"id":"top_ab_validate","type":"String","value":"0","chainId":"_all_"},{"id":"top_reason","type":"String","value":"1","chainId":"_all_"},{"id":"top_thank","type":"String","value":"1","chainId":"_all_"},{"id":"web_zhuanlan_api2_w","type":"String","value":"0"},{"id":"adr_topsearch","type":"String","value":"2"},{"id":"hb_column_v3_topbar","type":"String","value":"old"},{"id":"hb_recommend_column","type":"String","value":"1"},{"id":"ios_video_agent_4_28","type":"String","value":"false"},{"id":"top_new_user_gift","type":"String","value":"0","chainId":"_all_"},{"id":"adr_launch_ad_new_strategy","type":"String","value":"open"},{"id":"adr_spic","type":"String","value":"0"},{"id":"adr_zmcore","type":"String","value":"1"},{"id":"ios_le_nav","type":"String","value":"0"},{"id":"top_native_answer","type":"String","value":"1","chainId":"_all_"},{"id":"top_tr","type":"String","value":"0","chainId":"_all_"},{"id":"web_new_comment","type":"String","value":"1"},{"id":"adr_article_new_comment","type":"String","value":"0"},{"id":"ios_no_re_t","type":"String","value":"0"},{"id":"ios_roundtable","type":"String","value":"B"},{"id":"se_correct_ab","type":"String","value":"0","chainId":"_all_"},{"id":"top_yc","type":"String","value":"0","chainId":"_all_"},{"id":"adr_new_reader","type":"String","value":"1"},{"id":"adr_preload_question","type":"String","value":"0"},{"id":"adr_profile_medal","type":"String","value":"0"},{"id":"qa_test","type":"String","value":"0","chainId":"_all_"},{"id":"ios_new_player","type":"String","value":"0"},{"id":"ios_new_reader","type":"String","value":"1"},{"id":"se_ios_spb309bugfix","type":"String","value":"0","chainId":"_all_"},{"id":"se_minor_onebox","type":"String","value":"d","chainId":"_all_"},{"id":"adr_editor_version","type":"String","value":"V2"},{"id":"adr_video_ne_comment","type":"String","value":"0"},{"id":"hb_verification","type":"String","value":"0"},{"id":"ios_ad_cta","type":"String","value":"0"},{"id":"web_answer_list_ad","type":"String","value":"1"},{"id":"web_zhuanlan_api2","type":"String","value":"0"},{"id":"adr_easyza_on","type":"String","value":"0"},{"id":"ios_next_ans","type":"String","value":"N"},{"id":"ios_q_a_c","type":"String","value":"0"},{"id":"se_webrs","type":"String","value":"0","chainId":"_all_"},{"id":"ios_play_config_type","type":"String","value":"default"},{"id":"ios_q","type":"String","value":"0"},{"id":"top_scaled_score","type":"String","value":"0","chainId":"_all_"},{"id":"adr_android_p_type","type":"String","value":"default"},{"id":"adr_edit_question","type":"String","value":"0"},{"id":"adr_real_time_launch_http","type":"String","value":"http_off"},{"id":"ios_mini","type":"String","value":"0"},{"id":"adr_grow_guide_login_4","type":"String","value":"3"},{"id":"adr_hybrid_dns_v2","type":"String","value":"0"},{"id":"top_feedre_itemcf","type":"String","value":"31","chainId":"_all_"},{"id":"top_newfollow","type":"String","value":"0","chainId":"_all_"},{"id":"tp_dis_version","type":"String","value":"0","chainId":"_all_"},{"id":"adr_ans_video","type":"String","value":"N"},{"id":"ios_ge4","type":"String","value":"3"},{"id":"se_ios_spb309","type":"String","value":"0","chainId":"_all_"},{"id":"top_limit_num","type":"String","value":"0","chainId":"_all_"},{"id":"ios_adr_vid_vol","type":"String","value":"0"},{"id":"qa_video_answer_list","type":"String","value":"0","chainId":"_all_"},{"id":"se_entity","type":"String","value":"on","chainId":"_all_"},{"id":"web_heifetz_column_api2","type":"String","value":"0"},{"id":"adr_answer_dampen","type":"String","value":"0"},{"id":"adr_q_bar","type":"String","value":"NO"},{"id":"adr_unif","type":"String","value":"off"},{"id":"hb_liguangyi_test","type":"String","value":"1"},{"id":"ios_search_tab","type":"String","value":"1"},{"id":"se_search_feed","type":"String","value":"N","chainId":"_all_"},{"id":"top_root","type":"String","value":"0","chainId":"_all_"},{"id":"adr_bugly","type":"String","value":"n"},{"id":"adr_unfollow_reason","type":"String","value":"0"},{"id":"hb_majorob_style","type":"String","value":"0"},{"id":"ios_answer_preload","type":"String","value":"0"},{"id":"top_fqai","type":"String","value":"0","chainId":"_all_"},{"id":"top_sess_diversity","type":"String","value":"-1","chainId":"_all_"},{"id":"adr_android_medal_badge_view","type":"String","value":"false"},{"id":"adr_feed_video_continuous","type":"String","value":"0"},{"id":"adr_invite","type":"String","value":"false"},{"id":"se_gemini_service","type":"String","value":"content","chainId":"_all_"},{"id":"top_feedre","type":"String","value":"1","chainId":"_all_"},{"id":"web_question_invite","type":"String","value":"B"},{"id":"adr_sdk_data_switch","type":"String","value":"0"},{"id":"adr_wxbk","type":"String","value":"1"},{"id":"ios_ios_launch_mp4","type":"String","value":"1"},{"id":"soc_brandquestion","type":"String","value":"1","chainId":"_all_"},{"id":"top_30","type":"String","value":"0","chainId":"_all_"},{"id":"adr_float_video","type":"String","value":"1"},{"id":"adr_task_statistics","type":"String","value":"false"},{"id":"adr_wxfb","type":"String","value":"0"},{"id":"ios_db_n_e","type":"String","value":"0"},{"id":"se_bert","type":"String","value":"0","chainId":"_all_"},{"id":"top_round_table","type":"String","value":"0","chainId":"_all_"},{"id":"hb_active_answerer","type":"String","value":"0"},{"id":"hb_show_special_all","type":"String","value":"0"},{"id":"ios_wbtp","type":"String","value":"1"},{"id":"ls_is_use_zrec","type":"String","value":"0","chainId":"_all_"},{"id":"ios_more_editcard","type":"String","value":"true"},{"id":"top_ntr","type":"String","value":"1","chainId":"_all_"},{"id":"top_vd_gender","type":"String","value":"0","chainId":"_all_"},{"id":"adr_anr_watch","type":"String","value":"false"},{"id":"ios_lssq","type":"String","value":"0"},{"id":"ios_real_time_launch_http","type":"String","value":"http_off"},{"id":"top_distinction","type":"String","value":"0","chainId":"_all_"},{"id":"adr_osen_label","type":"String","value":"old"},{"id":"adr_question_invite_v2","type":"String","value":"0"},{"id":"hb_best_answerer","type":"String","value":"0"},{"id":"top_root_web","type":"String","value":"0","chainId":"_all_"},{"id":"tp_favsku","type":"String","value":"a","chainId":"_all_"},{"id":"tp_m_intro_re_topic","type":"String","value":"0","chainId":"_all_"},{"id":"tp_sticky_android","type":"String","value":"0","chainId":"_all_"},{"id":"adr_ydyq","type":"String","value":"X"},{"id":"ios_ad_skip_pos","type":"String","value":"up"},{"id":"ios_yhyq","type":"String","value":"C"},{"id":"top_hkc_test","type":"String","value":"1","chainId":"_all_"},{"id":"ios_video_agent_4_22","type":"String","value":"false"},{"id":"pin_ef","type":"String","value":"orig","chainId":"_all_"},{"id":"se_ad_index","type":"String","value":"10","chainId":"_all_"},{"id":"se_auto_syn","type":"String","value":"0","chainId":"_all_"},{"id":"ios_1752","type":"String","value":"0"},{"id":"ios_ff_cardtype","type":"String","value":"A"},{"id":"ios_hybrid_intercepting","type":"String","value":"1"},{"id":"ios_notif_new_invite","type":"String","value":"off"},{"id":"top_ydyq","type":"String","value":"X","chainId":"_all_"},{"id":"tp_header_style","type":"String","value":"0","chainId":"_all_"},{"id":"web_card_style","type":"String","value":"b"},{"id":"web_column_auto_invite","type":"String","value":"0"},{"id":"ios_question_new_comment","type":"String","value":"0"},{"id":"ios_q_bar","type":"String","value":"NO"},{"id":"top_f_r_nb","type":"String","value":"1","chainId":"_all_"},{"id":"top_gif","type":"String","value":"0","chainId":"_all_"},{"id":"adr_comment","type":"String","value":"false"},{"id":"ios_qtoc","type":"String","value":"0"},{"id":"top_follow_reason","type":"String","value":"0","chainId":"_all_"},{"id":"ug_zero_follow","type":"String","value":"0","chainId":"_all_"},{"id":"adr_profile_label","type":"String","value":"1"},{"id":"ios_telecom_login","type":"String","value":"0"},{"id":"ios_video_agent_4_32","type":"String","value":"true"},{"id":"top_nucc","type":"String","value":"0","chainId":"_all_"},{"id":"se_engine","type":"String","value":"0","chainId":"_all_"},{"id":"adr_consultation","type":"String","value":"0"},{"id":"hb_entity_ui","type":"String","value":"origin"},{"id":"ios_7324","type":"String","value":"0"},{"id":"ios_mlssq","type":"String","value":"0"},{"id":"adr_next_answer_btn","type":"String","value":"0"},{"id":"se_new_market_search","type":"String","value":"off","chainId":"_all_"},{"id":"top_cc_at","type":"String","value":"1","chainId":"_all_"},{"id":"web_follow_api_move","type":"String","value":"0"},{"id":"adr_liguangi_test","type":"String","value":"1"},{"id":"adr_prt","type":"String","value":"false"},{"id":"ios_vid_home","type":"String","value":"0"},{"id":"top_rerank_video","type":"String","value":"-1","chainId":"_all_"},{"id":"ios_quill_editor","type":"String","value":"0"},{"id":"ios_video_upload_cdn","type":"String","value":"1"},{"id":"ls_new_score","type":"String","value":"1","chainId":"_all_"},{"id":"pf_creator_card","type":"String","value":"1","chainId":"_all_"},{"id":"adr_add_account","type":"String","value":"1"},{"id":"adr_new_answer_pager","type":"String","value":"false"},{"id":"gw_guide","type":"String","value":"0","chainId":"_all_"},{"id":"ios_invite_ans","type":"String","value":"A"},{"id":"se_config","type":"String","value":"0","chainId":"_all_"},{"id":"se_majorob_style","type":"String","value":"0","chainId":"_all_"},{"id":"web_km_ab","type":"String","value":"1"},{"id":"se_webtimebox","type":"String","value":"0","chainId":"_all_"},{"id":"top_card","type":"String","value":"-1","chainId":"_all_"},{"id":"top_quality","type":"String","value":"0","chainId":"_all_"},{"id":"adr_audio_enable_exo","type":"String","value":"0"},{"id":"adr_hybrid_dns","type":"String","value":"0"},{"id":"adr_pre_load_html","type":"String","value":"0"},{"id":"ios_apm","type":"String","value":"1"},{"id":"qa_answerlist_ad","type":"String","value":"0","chainId":"_all_"},{"id":"se_colos","type":"String","value":"0","chainId":"_all_"},{"id":"top_feedre_cpt","type":"String","value":"101","chainId":"_all_"},{"id":"top_recall_tb_short","type":"String","value":"61","chainId":"_all_"},{"id":"ios_km_center","type":"String","value":"0"},{"id":"ios_launch_timeout_2","type":"String","value":"2000"},{"id":"ios_sw_regis_avatar","type":"String","value":"1"},{"id":"ios_tk_d","type":"String","value":"1"},{"id":"top_root_ac","type":"String","value":"1","chainId":"_all_"},{"id":"top_recall_exp_v1","type":"String","value":"1","chainId":"_all_"},{"id":"adr_member_switch","type":"String","value":"0"},{"id":"adr_ppt_enter","type":"String","value":"1"},{"id":"ios_hide_last_ac","type":"String","value":"0"},{"id":"ios_wxbk","type":"String","value":"1"},{"id":"tp_related_tps_movie","type":"String","value":"a","chainId":"_all_"},{"id":"adr_comment_new_editor","type":"String","value":"0"},{"id":"adr_traffic_monitor","type":"String","value":"false"},{"id":"ios_topsearch","type":"String","value":"1"},{"id":"ls_topic_is_use_zrec","type":"String","value":"0","chainId":"_all_"},{"id":"top_video_rerank","type":"String","value":"-1","chainId":"_all_"},{"id":"web_new_qa_related","type":"String","value":"1"},{"id":"zr_ans_rec","type":"String","value":"gbrank","chainId":"_all_"},{"id":"adr_hybrid_offline","type":"String","value":"0"},{"id":"adr_more_hyb_card","type":"String","value":"0"},{"id":"ios_article_nav","type":"String","value":"0"},{"id":"se_backsearch","type":"String","value":"0","chainId":"_all_"},{"id":"se_likebutton","type":"String","value":"0","chainId":"_all_"},{"id":"top_nad","type":"String","value":"1","chainId":"_all_"},{"id":"top_source","type":"String","value":"0","chainId":"_all_"},{"id":"top_topic_feedre","type":"String","value":"21","chainId":"_all_"},{"id":"adr_enable_agent","type":"String","value":"0"},{"id":"hb_follow_guide_wl","type":"String","value":"0"},{"id":"ios_vid_qt","type":"String","value":"0"},{"id":"ios_wxfb","type":"String","value":"0"},{"id":"top_new_user_rec","type":"String","value":"0","chainId":"_all_"},{"id":"adr_perm","type":"String","value":"0"},{"id":"ios_question_invite_v2","type":"String","value":"0"},{"id":"se_expired_ob","type":"String","value":"0","chainId":"_all_"},{"id":"top_is_gr","type":"String","value":"0","chainId":"_all_"},{"id":"adr_android_gdt","type":"String","value":"open"},{"id":"se_prf","type":"String","value":"0","chainId":"_all_"},{"id":"top_v_album","type":"String","value":"1","chainId":"_all_"},{"id":"tp_qa_metacard","type":"String","value":"0","chainId":"_all_"},{"id":"se_spb309","type":"String","value":"0","chainId":"_all_"},{"id":"adr_httpdns","type":"String","value":"aliyun"},{"id":"ios_ad_uiweb_open","type":"String","value":"1"},{"id":"ios_article_recommend_column","type":"String","value":"1"},{"id":"ios_vm_subject_type","type":"String","value":"0"},{"id":"zr_article_rec_rank","type":"String","value":"base","chainId":"_all_"},{"id":"zr_art_rec_rank","type":"String","value":"base","chainId":"_all_"},{"id":"zr_infinity","type":"String","value":"zr_infinity_close","chainId":"_all_"},{"id":"hb_live_btn_color","type":"String","value":"default_color"},{"id":"ios_ad_uiweb_js","type":"String","value":"1"},{"id":"ios_cmcc_login","type":"String","value":"1"},{"id":"ios_input_image","type":"String","value":"1"},{"id":"ios_profile_sig","type":"String","value":"true"},{"id":"se_consulting_price","type":"String","value":"n","chainId":"_all_"},{"id":"se_webmajorob","type":"String","value":"0","chainId":"_all_"},{"id":"tp_write_pin_guide","type":"String","value":"3","chainId":"_all_"},{"id":"adr_cashier_color","type":"String","value":"1"},{"id":"adr_upload_core","type":"String","value":"1"},{"id":"ios_answer_hybrid_preload","type":"String","value":"0"},{"id":"ios_pdf","type":"String","value":"n"},{"id":"ios_dkts","type":"String","value":"20"},{"id":"top_bill","type":"String","value":"0","chainId":"_all_"},{"id":"top_mt","type":"String","value":"0","chainId":"_all_"},{"id":"top_new_feed","type":"String","value":"1","chainId":"_all_"},{"id":"adr_recommend_column","type":"String","value":"1"},{"id":"adr_refresh_token","type":"String","value":"1"},{"id":"hb_major_onebox","type":"String","value":"0"},{"id":"hb_report","type":"String","value":"0"},{"id":"top_root_mg","type":"String","value":"1","chainId":"_all_"},{"id":"tp_related_topics","type":"String","value":" a","chainId":"_all_"},{"id":"ios_show_edit_image","type":"String","value":"1"},{"id":"se_billboardsearch","type":"String","value":"0","chainId":"_all_"},{"id":"top_ebook","type":"String","value":"0","chainId":"_all_"},{"id":"top_newfollowans","type":"String","value":"0","chainId":"_all_"},{"id":"adr_wbtp","type":"String","value":"1"},{"id":"ios_article_new_comment","type":"String","value":"0"},{"id":"ios_lans","type":"String","value":"close"},{"id":"ios_profile_badge","type":"String","value":"true"},{"id":"top_recall_exp_v2","type":"String","value":"1","chainId":"_all_"},{"id":"adr_anp","type":"String","value":"android_answer_pager_off"},{"id":"adr_editor_enabled","type":"String","value":"1"},{"id":"adr_traffic_threshold","type":"String","value":"314572800"},{"id":"ios_book_is_card","type":"String","value":"1"},{"id":"top_recall","type":"String","value":"0","chainId":"_all_"},{"id":"adr_video_topic_volume_control","type":"String","value":"0"},{"id":"top_recall_deep_user","type":"String","value":"1","chainId":"_all_"},{"id":"top_root_few_topic","type":"String","value":"0","chainId":"_all_"},{"id":"gue_new_special_page","type":"String","value":"0"},{"id":"ios_amap","type":"String","value":"y"},{"id":"se_wiki_box","type":"String","value":"1","chainId":"_all_"},{"id":"top_login_card","type":"String","value":"1","chainId":"_all_"},{"id":"top_billupdate1","type":"String","value":"2","chainId":"_all_"},{"id":"top_feedre_rtt","type":"String","value":"41","chainId":"_all_"},{"id":"zr_art_rec","type":"String","value":"base","chainId":"_all_"},{"id":"adr_hybrid_longc","type":"String","value":"0"},{"id":"adr_mini","type":"String","value":"0"},{"id":"ios_asp","type":"String","value":"off"},{"id":"ios_medal_badge_view","type":"String","value":"false"}],"chains":[{"chainId":"_all_"}]},"triggers":{}},"userAgent":{"Edge":false,"Wechat":false,"Weibo":false,"QQ":false,"Qzone":false,"Mobile":false,"Android":false,"iOS":false,"isAppleDevice":false,"Zhihu":false,"ZhihuHybrid":false,"isBot":false,"Tablet":false,"UC":false,"Sogou":false,"Qihoo":false,"Baidu":false,"BaiduApp":false,"Safari":false,"isWebView":false,"origin":"Mozilla\u002F5.0 (Windows NT 10.0; Win64; x64) AppleWebKit\u002F537.36 (KHTML, like Gecko) Chrome\u002F71.0.3578.98 Safari\u002F537.36"},"trafficSource":"production","edition":{"baidu":false,"sogou":false,"baiduBeijing":false,"yidianzixun":false},"theme":"light","referer":"","conf":{},"ipInfo":{},"logged":false,"tdkInfo":{}},"me":{"accountInfoLoadStatus":{},"organizationProfileStatus":{},"columnContributions":[]},"label":{},"comments":{"pagination":{},"collapsed":{},"reverse":{},"reviewing":{},"conversation":{},"parent":{}},"commentsV2":{"stickers":[],"commentWithPicPermission":{},"notificationsComments":{},"pagination":{},"collapsed":{},"reverse":{},"reviewing":{},"conversation":{},"conversationMore":{},"parent":{}},"pushNotifications":{"default":{"isFetching":false,"isDrained":false,"ids":[]},"follow":{"isFetching":false,"isDrained":false,"ids":[]},"vote_thank":{"isFetching":false,"isDrained":false,"ids":[]},"currentTab":"default","notificationsCount":{"default":0,"follow":0,"vote_thank":0}},"messages":{"data":{},"currentTab":"common","messageCount":0},"register":{"registerValidateSucceeded":null,"registerValidateErrors":{},"registerConfirmError":null,"sendDigitsError":null,"registerConfirmSucceeded":null},"login":{"loginUnregisteredError":false,"loginBindWechatError":false,"loginConfirmError":null,"sendDigitsError":null,"validateDigitsError":false,"loginConfirmSucceeded":null,"qrcodeLoginToken":"","qrcodeLoginScanStatus":0,"qrcodeLoginError":null,"qrcodeLoginReturnNewToken":false},"active":{"sendDigitsError":null,"activeConfirmSucceeded":null,"activeConfirmError":null},"switches":{},"captcha":{"captchaNeeded":false,"captchaValidated":false,"captchaBase64String":null,"captchaValidationMessage":null,"loginCaptchaExpires":false},"sms":{"supportedCountries":[]},"coupon":{"isRedeemingCoupon":false},"question":{"followers":{},"concernedFollowers":{},"answers":{},"hiddenAnswers":{},"updatedAnswers":{},"collapsedAnswers":{},"notificationAnswers":{},"invitationCandidates":{},"inviters":{},"invitees":{},"similarQuestions":{},"relatedCommodities":{},"recommendReadings":{},"bio":{},"brand":{},"permission":{},"adverts":{},"advancedStyle":{},"commonAnswerCount":0,"hiddenAnswerCount":0,"meta":{},"autoInvitation":{},"simpleConcernedFollowers":{}},"shareTexts":{},"answers":{"voters":{},"copyrightApplicants":{},"favlists":{},"newAnswer":{},"concernedUpvoters":{},"simpleConcernedUpvoters":{}},"banner":{},"topic":{"bios":{},"hot":{},"newest":{},"top":{},"unanswered":{},"questions":{},"followers":{},"contributors":{},"parent":{},"children":{},"bestAnswerers":{},"wikiMeta":{},"index":{},"intro":{},"meta":{},"schema":{},"creatorWall":{}},"explore":{"recommendations":{}},"articles":{"voters":{}},"favlists":{"relations":{}},"pins":{"voters":{}},"topstory":{"topstorys":{"isFetching":false,"isDrained":false,"afterId":0,"items":[],"next":null},"recommend":{"isFetching":false,"isDrained":false,"afterId":0,"items":[],"next":null},"follow":{"isFetching":false,"isDrained":false,"afterId":0,"items":[],"next":null},"followWonderful":{"isFetching":false,"isDrained":false,"afterId":0,"items":[],"next":null},"sidebar":null,"announcement":{},"hotList":[],"guestFeeds":{"isFetching":false,"isDrained":false,"afterId":0,"items":[],"next":null},"followExtra":{"isNewUser":null,"isFetched":false,"followCount":0,"followers":[]}},"upload":{},"video":{"data":{},"shareVideoDetail":{},"last":{}},"guide":{"guide":{"isFetching":false,"isShowGuide":false}},"reward":{"answer":{},"article":{},"question":{}},"search":{"recommendSearch":[],"topSearch":{},"attachedInfo":{},"nextOffset":{},"topicReview":{},"generalByQuery":{},"generalByQueryInADay":{},"generalByQueryInAWeek":{},"generalByQueryInThreeMonths":{},"peopleByQuery":{},"topicByQuery":{},"columnByQuery":{},"liveByQuery":{},"albumByQuery":{},"eBookByQuery":{}},"creator":{"currentCreatorUrlToken":null,"tools":{"question":{"invitationCount":{"questionFolloweeCount":0,"questionTotalCount":0},"goodatTopics":[]},"customPromotion":{"itemLists":{}},"recommend":{"recommendTimes":{}}},"explore":{"academy":{"tabs":[],"article":{}}},"rights":[],"rightsStatus":{},"levelUpperLimit":10,"account":{"growthLevel":{}}},"publicEditPermission":{},"readStatus":{}},"subAppName":"main"}</script><script src="https://static.zhihu.com/heifetz/vendor.4709996994b6c965ecab.js"></script><script src="https://static.zhihu.com/heifetz/main.app.79955df8146609ce4eea.js"></script><script src="https://static.zhihu.com/heifetz/main.question-routes.20480862074944c4aa07.js"></script></body></html>
`)

var html22 = `

<!DOCTYPE html>
<html lang="zh-cmn-Hans" class="ua-windows ua-webkit">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="renderer" content="webkit">
    <meta name="referrer" content="always">
    <meta name="google-site-verification" content="ok0wCgT20tBBgo9_zat2iAcimtN4Ftf5ccsh092Xeyw" />
    <title>
        大黄蜂 (豆瓣)
</title>
    
    <meta name="baidu-site-verification" content="cZdR4xxR7RxmM4zE" />
    <meta http-equiv="Pragma" content="no-cache">
    <meta http-equiv="Expires" content="Sun, 6 Mar 2005 01:00:00 GMT">
    
    <link rel="apple-touch-icon" href="https://img3.doubanio.com/f/movie/d59b2715fdea4968a450ee5f6c95c7d7a2030065/pics/movie/apple-touch-icon.png">
    <link href="https://img3.doubanio.com/f/shire/bf61b1fa02f564a4a8f809da7c7179b883a56146/css/douban.css" rel="stylesheet" type="text/css">
    <link href="https://img3.doubanio.com/f/shire/ae3f5a3e3085968370b1fc63afcecb22d3284848/css/separation/_all.css" rel="stylesheet" type="text/css">
    <link href="https://img3.doubanio.com/f/movie/8864d3756094f5272d3c93e30ee2e324665855b0/css/movie/base/init.css" rel="stylesheet">
    <script type="text/javascript">var _head_start = new Date();</script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/movie/0495cb173e298c28593766009c7b0a953246c5b5/js/movie/lib/jquery.js"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/f010949d3f23dd7c972ad7cb40b800bf70723c93/js/douban.js"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/0efdc63b77f895eaf85281fb0e44d435c6239a3f/js/separation/_all.js"></script>
    
    <meta name="keywords" content="大黄蜂,Bumblebee,大黄蜂影评,剧情介绍,电影图片,预告片,影讯,在线购票,论坛">
    <meta name="description" content="大黄蜂电影简介和剧情介绍,大黄蜂影评、图片、预告片、影讯、论坛、在线购票">
    <meta name="mobile-agent" content="format=html5; url=http://m.douban.com/movie/subject/26394152/"/>
    <link rel="alternate" href="android-app://com.douban.movie/doubanmovie/subject/26394152/" />
    <link rel="stylesheet" href="https://img3.doubanio.com/dae/cdnlib/libs/LikeButton/1.0.5/style.min.css">
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/77323ae72a612bba8b65f845491513ff3329b1bb/js/do.js" data-cfg-autoload="false"></script>
    <script type="text/javascript">
      Do.add('dialog', {path: 'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js', type: 'js'});
      Do.add('dialog-css', {path: 'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css', type: 'css'});
      Do.add('handlebarsjs', {path: 'https://img3.doubanio.com/f/movie/3d4f8e4a8918718256450eb6e57ec8e1f7a2e14b/js/movie/lib/handlebars.current.js', type: 'js'});
    </script>
    
  <script type='text/javascript'>
  var _vwo_code = (function() {
    var account_id = 249272,
      settings_tolerance = 0,
      library_tolerance = 2500,
      use_existing_jquery = false,
      // DO NOT EDIT BELOW THIS LINE
      f=false,d=document;return{use_existing_jquery:function(){return use_existing_jquery;},library_tolerance:function(){return library_tolerance;},finish:function(){if(!f){f=true;var a=d.getElementById('_vis_opt_path_hides');if(a)a.parentNode.removeChild(a);}},finished:function(){return f;},load:function(a){var b=d.createElement('script');b.src=a;b.type='text/javascript';b.innerText;b.onerror=function(){_vwo_code.finish();};d.getElementsByTagName('head')[0].appendChild(b);},init:function(){settings_timer=setTimeout('_vwo_code.finish()',settings_tolerance);var a=d.createElement('style'),b='body{opacity:0 !important;filter:alpha(opacity=0) !important;background:none !important;}',h=d.getElementsByTagName('head')[0];a.setAttribute('id','_vis_opt_path_hides');a.setAttribute('type','text/css');if(a.styleSheet)a.styleSheet.cssText=b;else a.appendChild(d.createTextNode(b));h.appendChild(a);this.load('//dev.visualwebsiteoptimizer.com/j.php?a='+account_id+'&u='+encodeURIComponent(d.URL)+'&r='+Math.random());return settings_timer;}};}());

  +function () {
    var bindEvent = function (el, type, handler) {
        var $ = window.jQuery || window.Zepto || window.$
       if ($ && $.fn && $.fn.on) {
           $(el).on(type, handler)
       } else if($ && $.fn && $.fn.bind) {
           $(el).bind(type, handler)
       } else if (el.addEventListener){
         el.addEventListener(type, handler, false);
       } else if (el.attachEvent){
         el.attachEvent("on" + type, handler);
       } else {
         el["on" + type] = handler;
       }
     }

    var _origin_load = _vwo_code.load
    _vwo_code.load = function () {
      var args = [].slice.call(arguments)
      bindEvent(window, 'load', function () {
        _origin_load.apply(_vwo_code, args)
      })
    }
  }()

  _vwo_settings_timer = _vwo_code.init();
  </script>


    


<script type="application/ld+json">
{
  "@context": "http://schema.org",
  "name": "大黄蜂 Bumblebee",
  "url": "/subject/26394152/",
  "image": "https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2541662397.webp",
  "director": 
  [
    {
      "@type": "Person",
      "url": "/celebrity/1305796/",
      "name": "特拉维斯·奈特 Travis Knight"
    }
  ]
,
  "author": 
  [
    {
      "@type": "Person",
      "url": "/celebrity/1364682/",
      "name": "克里斯蒂娜·霍德森 Christina Hodson"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1009507/",
      "name": "阿齐瓦·高斯曼 Akiva Goldsman"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1376970/",
      "name": "小豪尔赫·兰登伯格 Jorge Lendeborg Jr."
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1249805/",
      "name": "罗伯特·柯克曼 Robert Kirkman"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1293241/",
      "name": "肯·诺兰 Ken Nolan"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1306770/",
      "name": "阿特·马库姆 Art Marcum"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1286819/",
      "name": "马特·霍洛维 Matt Holloway"
    }
  ]
,
  "actor": 
  [
    {
      "@type": "Person",
      "url": "/celebrity/1312964/",
      "name": "海莉·斯坦菲尔德 Hailee Steinfeld"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1376970/",
      "name": "小豪尔赫·兰登伯格 Jorge Lendeborg Jr."
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1044883/",
      "name": "约翰·塞纳 John Cena"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1361062/",
      "name": "杰森·德鲁克 Jason Drucker"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1143044/",
      "name": "帕梅拉·阿德龙 Pamela Adlon"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1208831/",
      "name": "斯蒂芬·施耐德 Stephen Schneider"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1236720/",
      "name": "里卡多·霍约斯 Ricardo Hoyos"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1018110/",
      "name": "约翰·奥提兹 John Ortiz"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1007099/",
      "name": "格林·特鲁曼 Glynn Turman"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1027942/",
      "name": "兰·卡琉 Len Cariou"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1341549/",
      "name": "格蕾丝·达斯恩妮 Gracie Dzienny"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1009514/",
      "name": "弗里德·杜莱尔 Fred Dryer"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1239775/",
      "name": "蓝尼·雅各布森  Lenny Jacobson"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1137115/",
      "name": "梅金·普莱斯 Megyn Price"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1394730/",
      "name": "萨钦·巴特 Sachin Bhatt"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1360890/",
      "name": "蒂姆·马丁·格里森 Tim Martin Gleason"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1073264/",
      "name": "安东尼奥·查丽蒂 Antonio D. Charity"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1004934/",
      "name": "艾德文·霍德吉 Edwin Hodge"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1329869/",
      "name": "拉斯·斯兰德 Lars Slind"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1314963/",
      "name": "迪伦·奥布莱恩 Dylan O&#39;Brien"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1025212/",
      "name": "彼特·库伦 Peter Cullen"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1025214/",
      "name": "安吉拉·贝塞特 Angela Bassett"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1025199/",
      "name": "贾斯汀·塞洛克斯 Justin Theroux"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1153825/",
      "name": "大卫·索博洛夫 David Sobolov"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1086522/",
      "name": "格蕾·德丽斯勒 Grey DeLisle"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1064164/",
      "name": "史蒂夫·布卢姆 Steve Blum"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1375868/",
      "name": "安德鲁·莫尔加多 Andrew Morgado"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1393373/",
      "name": "威廉·W·巴伯 William W. Barbour"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1078881/",
      "name": "罗伯特·切斯纳特 Robert Chestnut"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1394733/",
      "name": "米歇尔·方 Michelle Fang"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1394734/",
      "name": "克里斯蒂安·哈切森 Christian Hutcherson"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1393387/",
      "name": "里克·理查森 Rick Richardson"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1330973/",
      "name": "瓦内萨·罗斯 Vanessa Ross"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1394731/",
      "name": "波士顿·拉什·弗里曼 Boston Rush Freeman"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1394729/",
      "name": "托尼·托斯特 Tony Toste"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1394732/",
      "name": "迪娜·特鲁迪 Deena Trudy"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1373172/",
      "name": "艾蒂安·维克 Etienne Vick"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1027505/",
      "name": "肯尼斯·崔 Kenneth Choi"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1393377/",
      "name": "玛塞拉·布拉吉奥 Marcella Bragio"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1339562/",
      "name": "蕾切尔·克劳 Rachel Crow"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1386474/",
      "name": "艾比·奎因 Abby Quinn"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1323535/",
      "name": "迈克尔·马西尼 Michael Masini"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1394735/",
      "name": "盖尔·甘布尔 Gail Gamble"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1002740/",
      "name": "马丁·肖特 Martin Short"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1349377/",
      "name": "特里萨·纳瓦罗 Teresa Navarro"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1016698/",
      "name": "杰斯·哈梅尔 Jess Harnell"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1393368/",
      "name": "尼娜·奇克 Nina Cheek"
    }
  ]
,
  "datePublished": "2018-12-21",
  "genre": ["\u52a8\u4f5c", "\u79d1\u5e7b", "\u5192\u9669"],
  "duration": "PT1H54M",
  "description": "本片故事设定在1987年，正值青春期的18岁少女查理Charlie（海莉·斯坦菲尔德 饰）在加州海边小镇的废弃场里发现了伤痕累累的大黄蜂，他们之间会发生怎样的故事呢？让我们拭目以待！",
  "@type": "Movie",
  "aggregateRating": {
    "@type": "AggregateRating",
    "ratingCount": "127869",
    "bestRating": "10",
    "worstRating": "2",
    "ratingValue": "7.2"
  }
}
</script>


    <style type="text/css">
  
</style>
    <style type="text/css">img { max-width: 100%; }</style>
    <script type="text/javascript"></script>
    <link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/639975bd8cedf840.css">

    <link rel="shortcut icon" href="https://img3.doubanio.com/favicon.ico" type="image/x-icon">
</head>

<body>
  
    <script type="text/javascript">var _body_start = new Date();</script>

    
    



    <link href="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.css" rel="stylesheet" type="text/css">



<div id="db-global-nav" class="global-nav">
  <div class="bd">
    
<div class="top-nav-info">
  <ul>
    <li>
    <a id="top-nav-doumail-link" href="https://www.douban.com/doumail/">豆邮</a>
    </li>
    <li class="nav-user-account">
      <a target="_blank" href="https://www.douban.com/accounts/" class="bn-more">
        <span>chenset的帐号</span><span class="arrow"></span>
      </a>
      <div class="more-items">
        <table cellpadding="0" cellspacing="0">
          <tbody>
            <tr>
              <td>
                <a href="https://www.douban.com/mine/">个人主页</a>
              </td>
            </tr>
            <tr>
              <td>
                <a target="_blank" href="https://www.douban.com/mine/orders/">我的订单</a>
              </td>
            </tr>
            <tr>
              <td>
                <a target="_blank" href="https://www.douban.com/mine/wallet/">我的钱包</a>
              </td>
            </tr>
            <tr>
              <td>
                <a target="_blank" href="https://www.douban.com/accounts/">帐号管理</a>
              </td>
            </tr>
            <tr>
              <td>
                <a href="https://www.douban.com/accounts/logout?source=movie&ck=8GU1">退出</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </li>
  </ul>
</div>

  <div class="top-nav-reminder">
    <a href="https://www.douban.com/notification/" class="lnk-remind">提醒</a>
    <div id="top-nav-notimenu" class="more-items">
      <div class="bd">
        <p>加载中...</p>
      </div>
    </div>
  </div>

    <div class="top-nav-doubanapp">
  <a href="https://www.douban.com/doubanapp/app?channel=top-nav" class="lnk-doubanapp">下载豆瓣客户端</a>
  <div id="doubanapp-tip">
    <a href="https://www.douban.com/doubanapp/app?channel=qipao" class="tip-link">豆瓣 <span class="version">6.0</span> 全新发布</a>
    <a href="javascript: void 0;" class="tip-close">×</a>
  </div>
  <div id="top-nav-appintro" class="more-items">
    <p class="appintro-title">豆瓣</p>
    <p class="qrcode">扫码直接下载</p>
    <div class="download">
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=iOS">iPhone</a>
      <span>·</span>
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=Android" class="download-android">Android</a>
    </div>
  </div>
</div>

    


<div class="global-nav-items">
  <ul>
    <li class="">
      <a href="https://www.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-main&quot;,&quot;uid&quot;:&quot;54106750&quot;}">豆瓣</a>
    </li>
    <li class="">
      <a href="https://book.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-book&quot;,&quot;uid&quot;:&quot;54106750&quot;}">读书</a>
    </li>
    <li class="on">
      <a href="https://movie.douban.com"  data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-movie&quot;,&quot;uid&quot;:&quot;54106750&quot;}">电影</a>
    </li>
    <li class="">
      <a href="https://music.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-music&quot;,&quot;uid&quot;:&quot;54106750&quot;}">音乐</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/location" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-location&quot;,&quot;uid&quot;:&quot;54106750&quot;}">同城</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/group" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-group&quot;,&quot;uid&quot;:&quot;54106750&quot;}">小组</a>
    </li>
    <li class="">
      <a href="https://read.douban.com&#47;?dcs=top-nav&amp;dcm=douban" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-read&quot;,&quot;uid&quot;:&quot;54106750&quot;}">阅读</a>
    </li>
    <li class="">
      <a href="https://douban.fm&#47;?from_=shire_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-fm&quot;,&quot;uid&quot;:&quot;54106750&quot;}">FM</a>
    </li>
    <li class="">
      <a href="https://time.douban.com&#47;?dt_time_source=douban-web_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-time&quot;,&quot;uid&quot;:&quot;54106750&quot;}">时间</a>
    </li>
    <li class="">
      <a href="https://market.douban.com&#47;?utm_campaign=douban_top_nav&amp;utm_source=douban&amp;utm_medium=pc_web" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-market&quot;,&quot;uid&quot;:&quot;54106750&quot;}">豆品</a>
    </li>
    <li>
      <a href="#more" class="bn-more"><span>更多</span></a>
      <div class="more-items">
        <table cellpadding="0" cellspacing="0">
          <tbody>
            <tr>
              <td>
                <a href="https://ypy.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-ypy&quot;,&quot;uid&quot;:&quot;54106750&quot;}">豆瓣摄影</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </li>
  </ul>
</div>

  </div>
</div>
<script>
  ;window._GLOBAL_NAV = {
    USER_ID: "54106750",
    UPLOAD_AUTH_TOKEN: "54106750:80176e2308350ff6c78297171d6301a7a79f323c",
    SSE_TOKEN: "c00642e03d8d6587e172139ce55d9c3d998e3079",
    SSE_TIMESTAMP: "1547951968",
    DOUBAN_URL: "https://www.douban.com",
    N_NEW_NOTIS: 0,
    N_NEW_DOUMAIL: 0
  };
</script>



    <script src="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.js" defer="defer"></script>




    



    <link href="//img3.doubanio.com/dae/accounts/resources/8c80301/movie/bundle.css" rel="stylesheet" type="text/css">




<div id="db-nav-movie" class="nav">
  <div class="nav-wrap">
  <div class="nav-primary">
    <div class="nav-logo">
      <a href="https:&#47;&#47;movie.douban.com">豆瓣电影</a>
    </div>
    <div class="nav-search">
      <form action="https:&#47;&#47;movie.douban.com/subject_search" method="get">
        <fieldset>
          <legend>搜索：</legend>
          <label for="inp-query">
          </label>
          <div class="inp"><input id="inp-query" name="search_text" size="22" maxlength="60" placeholder="搜索电影、电视剧、综艺、影人" value=""></div>
          <div class="inp-btn"><input type="submit" value="搜索"></div>
          <input type="hidden" name="cat" value="1002" />
        </fieldset>
      </form>
    </div>
  </div>
  </div>
  <div class="nav-secondary">
    

<div class="nav-items">
  <ul>
    <li    ><a href="https://movie.douban.com/mine"
     >我看</a>
    </li>
    <li    ><a href="https://movie.douban.com/cinema/nowplaying/"
     >影讯&购票</a>
    </li>
    <li    ><a href="https://movie.douban.com/explore"
     >选电影</a>
    </li>
    <li    ><a href="https://movie.douban.com/tv/"
     >电视剧</a>
    </li>
    <li    ><a href="https://movie.douban.com/chart"
     >排行榜</a>
    </li>
    <li    ><a href="https://movie.douban.com/tag/"
     >分类</a>
    </li>
    <li    ><a href="https://movie.douban.com/review/best/"
     >影评</a>
    </li>
    <li    ><a href="https://movie.douban.com/annual/2018?source=navigation"
     >2018年度榜单</a>
    </li>
    <li    ><a href="https://www.douban.com/standbyme/2018?source=navigation"
     >2018书影音报告</a>
    </li>
  </ul>
</div>

    <a href="https://movie.douban.com/annual/2018?source=movie_navigation" class="movieannual2018"></a>
  </div>
</div>

<script id="suggResult" type="text/x-jquery-tmpl">
  <li data-link="{{= url}}">
            <a href="{{= url}}" onclick="moreurl(this, {from:'movie_search_sugg', query:'{{= keyword }}', subject_id:'{{= id}}', i: '{{= index}}', type: '{{= type}}'})">
            <img src="{{= img}}" width="40" />
            <p>
                <em>{{= title}}</em>
                {{if year}}
                    <span>{{= year}}</span>
                {{/if}}
                {{if sub_title}}
                    <br /><span>{{= sub_title}}</span>
                {{/if}}
                {{if address}}
                    <br /><span>{{= address}}</span>
                {{/if}}
                {{if episode}}
                    {{if episode=="unknow"}}
                        <br /><span>集数未知</span>
                    {{else}}
                        <br /><span>共{{= episode}}集</span>
                    {{/if}}
                {{/if}}
            </p>
        </a>
        </li>
  </script>




    <script src="//img3.doubanio.com/dae/accounts/resources/8c80301/movie/bundle.js" defer="defer"></script>





    
    <div id="wrapper">
        

        
    <div id="content">
        

    <div id="dale_movie_subject_top_icon"></div>
    <h1>
        <span property="v:itemreviewed">大黄蜂 Bumblebee</span>
            <span class="year">(2018)</span>
    </h1>

        <div class="grid-16-8 clearfix">
            

            
            <div class="article">
                
    

    





        <div class="indent clearfix">
            <div class="subjectwrap clearfix">
                <div class="subject clearfix">
                    



<div id="mainpic" class="">
    <a class="nbgnbg" href="https://movie.douban.com/subject/26394152/photos?type=R" title="点击看更多海报">
        <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2541662397.webp" title="点击看更多海报" alt="Bumblebee" rel="v:image" />
   </a>
                <p class="gact"><a href="https://movie.douban.com/subject/26394152/edit">更新描述或海报</a></p>
</div>

                    


<div id="info">
        <span ><span class='pl'>导演</span>: <span class='attrs'><a href="/celebrity/1305796/" rel="v:directedBy">特拉维斯·奈特</a></span></span><br/>
        <span ><span class='pl'>编剧</span>: <span class='attrs'><a href="/celebrity/1364682/">克里斯蒂娜·霍德森</a></span></span><br/>
        <span class="actor"><span class='pl'>主演</span>: <span class='attrs'><a href="/celebrity/1312964/" rel="v:starring">海莉·斯坦菲尔德</a> / <a href="/celebrity/1376970/" rel="v:starring">小豪尔赫·兰登伯格</a> / <a href="/celebrity/1044883/" rel="v:starring">约翰·塞纳</a> / <a href="/celebrity/1361062/" rel="v:starring">杰森·德鲁克</a> / <a href="/celebrity/1143044/" rel="v:starring">帕梅拉·阿德龙</a> / <a href="/celebrity/1208831/" rel="v:starring">斯蒂芬·施耐德</a> / <a href="/celebrity/1236720/" rel="v:starring">里卡多·霍约斯</a> / <a href="/celebrity/1018110/" rel="v:starring">约翰·奥提兹</a> / <a href="/celebrity/1007099/" rel="v:starring">格林·特鲁曼</a> / <a href="/celebrity/1027942/" rel="v:starring">兰·卡琉</a> / <a href="/celebrity/1341549/" rel="v:starring">格蕾丝·达斯恩妮</a> / <a href="/celebrity/1009514/" rel="v:starring">弗里德·杜莱尔</a> / <a href="/celebrity/1239775/" rel="v:starring">蓝尼·雅各布森 </a> / <a href="/celebrity/1137115/" rel="v:starring">梅金·普莱斯</a> / <a href="/celebrity/1394730/" rel="v:starring">萨钦·巴特</a> / <a href="/celebrity/1360890/" rel="v:starring">蒂姆·马丁·格里森</a> / <a href="/celebrity/1073264/" rel="v:starring">安东尼奥·查丽蒂</a> / <a href="/celebrity/1004934/" rel="v:starring">艾德文·霍德吉</a> / <a href="/celebrity/1329869/" rel="v:starring">拉斯·斯兰德</a> / <a href="/celebrity/1314963/" rel="v:starring">迪伦·奥布莱恩</a> / <a href="/celebrity/1025212/" rel="v:starring">彼特·库伦</a> / <a href="/celebrity/1025214/" rel="v:starring">安吉拉·贝塞特</a> / <a href="/celebrity/1025199/" rel="v:starring">贾斯汀·塞洛克斯</a> / <a href="/celebrity/1153825/" rel="v:starring">大卫·索博洛夫</a> / <a href="/celebrity/1086522/" rel="v:starring">格蕾·德丽斯勒</a> / <a href="/celebrity/1064164/" rel="v:starring">史蒂夫·布卢姆</a> / <a href="/celebrity/1375868/" rel="v:starring">安德鲁·莫尔加多</a> / <a href="/celebrity/1393373/" rel="v:starring">威廉·W·巴伯</a> / <a href="/celebrity/1078881/" rel="v:starring">罗伯特·切斯纳特</a> / <a href="/celebrity/1394733/" rel="v:starring">米歇尔·方</a> / <a href="/celebrity/1394734/" rel="v:starring">克里斯蒂安·哈切森</a> / <a href="/celebrity/1393387/" rel="v:starring">里克·理查森</a> / <a href="/celebrity/1330973/" rel="v:starring">瓦内萨·罗斯</a> / <a href="/celebrity/1394731/" rel="v:starring">波士顿·拉什·弗里曼</a> / <a href="/celebrity/1394729/" rel="v:starring">托尼·托斯特</a> / <a href="/celebrity/1394732/" rel="v:starring">迪娜·特鲁迪</a> / <a href="/celebrity/1373172/" rel="v:starring">艾蒂安·维克</a></span></span><br/>
        <span class="pl">类型:</span> <span property="v:genre">动作</span> / <span property="v:genre">科幻</span> / <span property="v:genre">冒险</span><br/>
        
        <span class="pl">制片国家/地区:</span> 美国<br/>
        <span class="pl">语言:</span> 英语<br/>
        <span class="pl">上映日期:</span> <span property="v:initialReleaseDate" content="2019-01-04(中国大陆)">2019-01-04(中国大陆)</span> / <span property="v:initialReleaseDate" content="2018-12-21(美国)">2018-12-21(美国)</span><br/>
        <span class="pl">片长:</span> <span property="v:runtime" content="114">114分钟</span><br/>
        <span class="pl">又名:</span> 大黄蜂大电影 / 大黄蜂独立电影 / 变形金刚外传：大黄蜂 / 变形金刚外传大黄蜂 / Brighton Falls<br/>
        <span class="pl">IMDb链接:</span> <a href="http://www.imdb.com/title/tt4701182" target="_blank" rel="nofollow">tt4701182</a><br>

</div>




                </div>
                    


<div id="interest_sectl">
    <div class="rating_wrap clearbox" rel="v:rating">
        <div class="clearfix">
          <div class="rating_logo ll">豆瓣评分</div>
          <div class="output-btn-wrap rr" style="display:none">
            <img src="https://img3.doubanio.com/f/movie/692e86756648f29457847c5cc5e161d6f6b8aaac/pics/movie/reference.png" />
            <a class="download-output-image" href="#">引用</a>
          </div>
          
          
        </div>
        


<div class="rating_self clearfix" typeof="v:Rating">
    <strong class="ll rating_num" property="v:average">7.2</strong>
    <span property="v:best" content="10.0"></span>
    <div class="rating_right ">
        <div class="ll bigstar bigstar35"></div>
        <div class="rating_sum">
                <a href="collections" class="rating_people"><span property="v:votes">127928</span>人评价</a>
        </div>
    </div>
</div>
<div class="ratings-on-weight">
    
        <div class="item">
        
        <span class="stars5 starstop" title="力荐">
            5星
        </span>
        <div class="power" style="width:19px"></div>
        <span class="rating_per">13.0%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars4 starstop" title="推荐">
            4星
        </span>
        <div class="power" style="width:64px"></div>
        <span class="rating_per">42.2%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars3 starstop" title="还行">
            3星
        </span>
        <div class="power" style="width:58px"></div>
        <span class="rating_per">38.4%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars2 starstop" title="较差">
            2星
        </span>
        <div class="power" style="width:8px"></div>
        <span class="rating_per">5.6%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars1 starstop" title="很差">
            1星
        </span>
        <div class="power" style="width:1px"></div>
        <span class="rating_per">0.8%</span>
        <br />
        </div>
</div>

    </div>
        
            <div class="friends_rating_wrap clearbox">
                <div class="rating_logo_wrap">
                    <div class="content">好友评分</div>
                    <div class="rating_helper_wrap">
                        <span class="rating_helper_icon"></span>
                        <span class="rating_helper_content">你关注的人看过这部作品的平均分</span>
                    </div>
                </div>
                <div class="rating_content_wrap clearfix">
                    <strong class="rating_avg">7.0</strong>
                    <div class="friends">
                            <a class="avatar" title="是灼灼" href="javascript:;">
                                <img src="https://img3.doubanio.com/icon/u62054075-164.jpg" alt="是灼灼">
                            </a>
                            <a class="avatar" title="l ǐ" href="javascript:;">
                                <img src="https://img1.doubanio.com/icon/u35767782-278.jpg" alt="l ǐ">
                            </a>
                            <a class="avatar" title="饭" href="javascript:;">
                                <img src="https://img3.doubanio.com/icon/u50082009-184.jpg" alt="饭">
                            </a>
                    </div>
                    <a href="follows_comments" class="friends_count" target="_blank">4人评价</a>
                </div>
            </div>
        <div class="rating_betterthan">
            好于 <a href="/typerank?type_name=科幻&type=17&interval_id=70:60&action=">68% 科幻片</a><br/>
            好于 <a href="/typerank?type_name=动作&type=5&interval_id=70:60&action=">69% 动作片</a><br/>
        </div>
</div>


                
            </div>
                




<div id="interest_sect_level" class="clearfix">
        
            <a href="https://movie.douban.com/subject/26394152/?interest=wish&amp;ck=8GU1" rel="nofollow" class="collect_btn colbutt ll" name="pbtn-26394152-wish">
                <span>想看</span>
            </a>
            <a href="https://movie.douban.com/subject/26394152/?interest=collect&amp;ck=8GU1" rel="nofollow" class="collect_btn colbutt ll" name="pbtn-26394152-collect">
                <span>看过</span>
            </a>
        <div class="ll j a_stars">
            
    
    评价:
    <span id="rating"> <span id="stars" data-solid="https://img3.doubanio.com/f/shire/5a2327c04c0c231bced131ddf3f4467eb80c1c86/pics/rating_icons/star_onmouseover.png" data-hollow="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" data-solid-2x="https://img3.doubanio.com/f/shire/7258904022439076d57303c3b06ad195bf1dc41a/pics/rating_icons/star_onmouseover@2x.png" data-hollow-2x="https://img3.doubanio.com/f/shire/95cc2fa733221bb8edd28ad56a7145a5ad33383e/pics/rating_icons/star_hollow_hover@2x.png">

                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-26394152-collect-1">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star1" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-26394152-collect-2">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star2" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-26394152-collect-3">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star3" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-26394152-collect-4">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star4" width="16" height="16"/></a>
                    <a href="javascript:;" class="j a_collect_btn" name="pbtn-26394152-collect-5">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star5" width="16" height="16"/></a>
    </span><span id="rateword" class="pl"></span>
    <input id="n_rating" type="hidden" value=""  />
    </span>

        </div>
    

</div>


            


















<div class="gtleft">
    <ul class="ul_subject_menu bicelink color_gray pt6 clearfix">
        
    
        
                <li> 
    <img src="https://img3.doubanio.com/f/shire/cc03d0fcf32b7ce3af7b160a0b85e5e66b47cc42/pics/short-comment.gif" />&nbsp;
        <a onclick="moreurl(this, {from:'mv_sbj_wr_cmnt'})" href="javascript:;" class="j a_collect_btn" name="cbtn-26394152">写短评</a>
 </li>
                    <li> 
    
    <img src="https://img3.doubanio.com/f/shire/5bbf02b7b5ec12b23e214a580b6f9e481108488c/pics/add-review.gif" />&nbsp;
        <a onclick="moreurl(this, {from:'mv_sbj_wr_rv'})" class="create_from_menu" href="https://movie.douban.com/subject/26394152/new_review" rel="nofollow">写影评</a>
 </li>
                    <li> 
    <img src="https://img3.doubanio.com/f/shire/61cc48ba7c40e0272d46bb93fe0dc514f3b71ec5/pics/add-doulist.gif" />&nbsp;
    <a href="/subject/26394152/questions/ask?from=subject_top">提问题</a>
 </li>
                <li> 
    


    <div class="doulist-add-btn">
  

  

  
  <a href="javascript:void(0)"
     data-id="26394152"
     data-cate="1002"
     data-canview="True"
     data-url="https://movie.douban.com/subject/26394152/"
     data-catename="电影"
     data-link="https://www.douban.com/people/chenset/doulists/all?add=26394152&amp;cat=1002"
     data-title="大黄蜂 Bumblebee"
     data-picture="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2541662397.webp"
     class="lnk-doulist-add"
     onclick="moreurl(this, { 'from':'doulist-btn-1002-26394152-54106750'})">
      <i></i>添加到豆列
  </a>
    </div>

 </li>
                <li> 
   

   
    
    <span class="rec" id="电影-26394152">
    <a href= "#"
        data-type="电影"
        data-url="https://movie.douban.com/subject/26394152/"
        data-desc="电影《大黄蜂 Bumblebee》 (来自豆瓣) "
        data-title="电影《大黄蜂 Bumblebee》 (来自豆瓣) "
        data-pic="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2541662397.jpeg"
        class="bn-sharing ">
        分享到
    </a> &nbsp;&nbsp;
    </span>

    <script>
        if (!window.DoubanShareMenuList) {
            window.DoubanShareMenuList = [];
        }
        var __cache_url = __cache_url || {};

        (function(u){
            if(__cache_url[u]) return;
            __cache_url[u] = true;
            window.DoubanShareIcons = 'https://img3.doubanio.com/f/shire/d15ffd71f3f10a7210448fec5a68eaec66e7f7d0/pics/ic_shares.png';

            var initShareButton = function() {
                $.ajax({url:u,dataType:'script',cache:true});
            };

            if (typeof Do == 'function' && 'ready' in Do) {
                Do(
                    'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
                    'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
                    'https://img3.doubanio.com/f/movie/c4ab132ff4d3d64a83854c875ea79b8b541faf12/js/movie/lib/qrcode.min.js',
                    initShareButton
                );
            } else if(typeof Douban == 'object' && 'loader' in Douban) {
                Douban.loader.batch(
                    'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
                    'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
                    'https://img3.doubanio.com/f/movie/c4ab132ff4d3d64a83854c875ea79b8b541faf12/js/movie/lib/qrcode.min.js'
                ).done(initShareButton);
            }

        })('https://img3.doubanio.com/f/movie/32be6727ed3ad8f6c4a417d8a086355c3e7d1d27/js/movie/lib/sharebutton.js');
    </script>


  </li>
            

    </ul>

    <script type="text/javascript">
        $(function(){
            $(".ul_subject_menu li.rec .bn-sharing").bind("click", function(){
                $.get("/blank?sbj_page_click=bn_sharing");
            });
            $(".ul_subject_menu .create_from_menu").bind("click", function(e){
                e.preventDefault();
                var $el = $(this);
                var glRoot = document.getElementById('gallery-topics-selection');
                if (window.has_gallery_topics && glRoot) {
                    // 判断是否有话题
                    glRoot.style.display = 'block';
                } else {
                    location.href = $el.attr('href');
                }
            });
        });
    </script>
</div>




                




        <style type="text/css">
            
.suggestions-list li { position: relative; left: 0; top: 0; margin-bottom: 7px; height: 35px }
.suggestions-list li .user-thumb { display: inline-block; *display: inline; float: left; margin: 2px 5px 0 0; vertical-align: top }
.suggestions-list li .user-thumb img { width: 25px; height: 25px }
.suggestions-list li .user-name-info { display: inline-block; *display: inline; line-height: 1.4em }
.suggestions-list li .user-name-info .user-profile-link { color: #333; font-weight: 800 }
.suggestions-list li .user-name-info .user-profile-link:hover { color: #4b8dc5 }
.suggestions-list li .user-name-info p { color: #999 }
.suggestions-list li .user-name-info b { color: #4b8dc5; font-weight: normal; cursor: pointer }
.suggestions-list li .user-name-info b:hover { text-decoration: underline }
.suggestions-list li .dismiss { position: absolute }
.suggestions-list li .dismiss { color: #aaa; margin: 0 0 0 10px; top: 0; right: 0 }
.suggestions-list li .dismiss:hover { color: #333; text-decoration: none }


.suggest-overlay { position: absolute; z-index: 99; width: auto; background: #fff; border: 1px solid #c5c7d2;
    -moz-border-radius : 3px;
    -webkit-border-radius : 3px;
    border-radius: 3px
}
.suggest-overlay .bd { min-width: 220px; line-height: 1; background: #fafafa; color: #b3b3b3; padding: 5px;
    -moz-border-radius : 3px;
    -webkit-border-radius : 3px;
    border-radius: 3px
}
.suggest-overlay ul { color: #999; padding: 3px 0; min-width: 214px }
.suggest-overlay li { cursor: pointer; padding: 3px 7px }
.suggest-overlay li b { font-weight: bold }
.suggest-overlay li .username { color: #333 }
.suggest-overlay img { margin-right: 5px; width: 20px; height: 20px; vertical-align: middle }
.suggest-overlay .on { background: #e9f0f8 }

.mentioned-highlighter { font: 14px/20px "Helvetica Neue",Helvetica,Arial,sans-serif; position: absolute; left: 4px; top: 4px; font-size: 14px; height: 60px; width: 98.5%; overflow: hidden; background: #fff; white-space: pre-wrap; word-wrap: break-word; color: transparent }
.mentioned-highlighter b { font-weight: normal; background-color: #d2e1f3; color: transparent;
  -moz-border-radius: 2px;
  -webkit-border-radius: 2px;
  border-radius: 2px
}

            .movie-share-dialog .bn-flat input {
    font-size: 14px;
}
.movie-share-dialog {
    z-index: 100;
}
.movie-share-dialog .form-ft-inner{
    text-align: right;
}
.movie-share-dialog div.bd {
    padding: 0;
}

.movie-share .form-bd .input-area {
    position: relative;
    margin: 15px;
    height: 91px;
    zoom: 1;
}

.movie-share-no-media .form-bd {
    height: 140px;
}

.movie-share-dialog .share-text {
    height: 85px;
    position: absolute;
    z-index: 9;
    background: transparent;
    font: 14px/18px "Helvetica Neue",Helvetica,Arial,sans-serif;
    width: 98%;
    -webkit-border-radius: 4px 4px 4px 4px;
    border-radius: 4px 4px 4px 4px;
}

.movie-share-dialog .mentioned-highlighter {
    width: 483px;
    padding: 3px 4px 4px;
    color: white;
    position: absolute;
    top:0;
    left:0;
    z-index: 0;
}

.movie-share-dialog .mentioned-highlighter code {
    color: #D2E1F3;
    background: #D2E1F3;
    border-radius: 2px;
    padding-right: 1px;
    display: inline-block;
    font: 14px/18px "Helvetica Neue",Helvetica,Arial,sans-serif;
}


.movie-share .form-ft {
    background: #e9eef2;
    height: 25px;
    padding-top: 10px;
    padding-bottom: 10px;
}

.movie-share .form-ft-inner {
    height: 25px;
    padding-left: 15px;
    padding-right: 15px;
}

.movie-share-dialog .dialog-only-text {
    text-align: center;
    font-size: 14px;
    line-height: 1.5;
    padding-top: 30px;
    padding-bottom: 30px;
    color: #0c7823;
}

.movie-share-dialog .ll {
    float: left;
    display: inline;
}
.movie-share-dialog .share-label {
    width: auto;
    display: inline;
    float: none;
}

.movie-share-dialog .leading-label {
    _vertical-align: -2px;
}
.movie-share-dialog .media {
    float: left;
    margin-right: 10px;
    max-width: 100px;
    max-height: 100px;
    *width: 100px;
}
.movie-share-dialog .info-area{
    overflow: hidden;
    zoom: 1;
    margin: 0 15px 15px;
    height: 100px;
}
.movie-share-dialog .info-area strong{
    font-weight: bold;
}
.movie-share-dialog .info-area p{
    margin: 3px 0;
    color: #999;
}

.movie-share-dialog #sync-setting {
    _vertical-align: -5px;
    margin-left: 10px;
}

.movie-share-dialog .info-area .server-error {
    position: absolute;
    bottom: 45px;
    right: 15px;
    color: red;
}

.movie-share-dialog .avail-num-indicator {
    color: #aaa;
    font-weight: 800;
    padding-right: 3px;
}

.movie-share-dialog .bottom-setting {
    width: 432px;
}
.movie-share-dialog .input-checkbox {
    vertical-align: -2px;
    _vertical-align: -1px;
}

.movie-share-dialog #sync-setting img {
    _vertical-align: 2px;
}



.suggest-overlay {
    z-index: 2000;
}

.movie-bar {
    position: relative;
    margin-top: 10px;
}

.movie-bar-fav {
    position: absolute;
    top: 0;
    right: 0;
}

        </style>
        <script src="https://img3.doubanio.com/f/shire/a40c5220b3f40ce737b366c0030ecf810b37bfea/js/lib/mustache.js" type="text/javascript"></script>
        <script src="https://img3.doubanio.com/f/shire/1d985568f3cc434b145983919d9954e2ca627e9c/js/lib/textarea-mention.js" type="text/javascript"></script>
        <script src="https://img3.doubanio.com/f/movie/6b10694c6523e81ebdea9963901b757cf91387f6/js/movie/share.js" type="text/javascript"></script>

<div class="rec-sec">
<span class="rec">
    <script id="movie-share" type="text/x-html-snippet">
        
    <form class="movie-share" action="/j/share" method="POST"><div style="display:none;"><input type="hidden" name="ck" value="8GU1"/></div>
        <div class="clearfix form-bd">
            <div class="input-area">
                <textarea name="text" class="share-text" cols="72" data-mention-api="https://api.douban.com/shuo/in/complete?alt=xd&amp;callback=?"></textarea>
                <input type="hidden" name="target-id" value="26394152">
                <input type="hidden" name="target-type" value="0">
                <input type="hidden" name="title" value="大黄蜂 Bumblebee‎ (2018)">
                <input type="hidden" name="desc" value="导演 特拉维斯·奈特 主演 海莉·斯坦菲尔德 / 小豪尔赫·兰登伯格 / 美国 / 7.2分(127928评价)">
                <input type="hidden" name="redir" value=""/>
                <div class="mentioned-highlighter"></div>
            </div>

            <div class="info-area">
                    <img class="media" src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2541662397.webp" />
                <strong>大黄蜂 Bumblebee‎ (2018)</strong>
                <p>导演 特拉维斯·奈特 主演 海莉·斯坦菲尔德 / 小豪尔赫·兰登伯格 / 美国 / 7.2分(127928评价)</p>
                <p class="error server-error">&nbsp;</p>
            </div>
        </div>
        <div class="form-ft">
            <div class="form-ft-inner">
                




<div class="sync-setting pl">
    <label>分享到</label>



            <a id="lnk-sync-setting" class="no-visited no-hover" href="https://movie.douban.com/settings/sync" target="_blank"
                class="pl share-label"><img src="https://img3.doubanio.com/f/movie/9389c4e5cab0cd1089a189d607d296c31ddb1bc0/pics/movie/share_g.png"
                alt="去绑定新浪微博" />去绑定新浪微博</a>

</div>

<style type="text/css">
    .sync-setting {
        float: left;
    }
    #lnk-sync-setting,
    #lnk-sync-setting:hover,
    #lnk-sync-setting:visited {
        vertical-align: middle;
        color: #0192b5;
        background: none;
        line-height: 27px;
        margin-right: 8px;
    }
    #lnk-sync-setting img {
        vertical-align: baseline;
        *vertical-align: middle;
        opacity: .5;
        filter: alpha(opacity=50);
        display: inline-block;
        width: 10px;
        height: 10px;
        *display: inline;
        zoom: 1;
        position: relative;
        top: 1px;
        margin-left: 5px;
    }
    #lnk-sync-setting:hover img {
        opacity: .8;
        background: none;
        filter: alpha(opacity=80);
    }
    .share-label {
        margin: 8px;
        cursor: pointer;
        vertical-align: middle;
        *vertical-align: text-bottom;
    }
    .interest-form-ft .share-label input {
        margin-bottom: 0;
    }
    .interest-form-ft {
        text-align: right;
    }
    .interest-form-ft .bn-flat {
        float: none;
    }
    .interest-form-ft .sync-setting {
        float: left;
        line-height: 25px;
    }
</style>


                <span class="avail-num-indicator">140</span>
                <span class="bn-flat">
                    <input type="submit" value="推荐" />
                </span>
            </div>
        </div>
    </form>
    
    <div id="suggest-mention-tmpl" style="display:none;">
        <ul>
            {{#users}}
            <li id="{{uid}}">
              <img src="{{avatar}}">{{{username}}}&nbsp;<span>({{{uid}}})</span>
            </li>
            {{/users}}
        </ul>
    </div>


    </script>

        
        <a href="#" data-share-dialog="#movie-share" data-dialog-title="推荐电影" class="lnk-sharing" share-id="26394152" data-mode="plain" data-name="大黄蜂 Bumblebee‎ (2018)" data-type="movie" data-desc="导演 特拉维斯·奈特 主演 海莉·斯坦菲尔德 / 小豪尔赫·兰登伯格 / 美国 / 7.2分(127928评价)" data-href="https://movie.douban.com/subject/26394152/" data-image="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2541662397.webp" data-properties="{}" data-redir="" data-text="" data-apikey="" data-curl="" data-count="10" data-object_kind="1002" data-object_id="26394152" data-target_type="rec" data-target_action="0" data-action_props="{&#34;subject_url&#34;:&#34;https:\/\/movie.douban.com\/subject\/26394152\/&#34;,&#34;subject_title&#34;:&#34;大黄蜂 Bumblebee‎ (2018)&#34;}">推荐</a>
</span>


</div>






            <script type="text/javascript">
                $(function() {
                    $('.collect_btn', '#interest_sect_level').each(function() {
                        Douban.init_collect_btn(this);
                    });
                    $('html').delegate(".indent .rec-sec .lnk-sharing", "click", function() {
                        moreurl(this, {
                            from : 'mv_sbj_db_share'
                        });
                    });
                });
            </script>
        </div>
            


    <div id="collect_form_26394152"></div>


        



<div class="related-info" style="margin-bottom:-10px;">
    <a name="intro"></a>
    
        
            
            
    <h2>
        <i class="">大黄蜂的剧情简介</i>
              · · · · · ·
    </h2>

            <div class="indent" id="link-report">
                    
                        <span property="v:summary" class="">
                                　　本片故事设定在1987年，正值青春期的18岁少女查理Charlie（海莉·斯坦菲尔德 饰）在加州海边小镇的废弃场里发现了伤痕累累的大黄蜂，他们之间会发生怎样的故事呢？让我们拭目以待！
                        </span>
                        
<script type="text/javascript" src="https://img3.doubanio.com/f/shire/a14501790b4a2db257dc5be5e37d820e600703c6/js/report_dialog.js"></script>
<link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/b45aa277f8b8df40596b96582dafb1ed0a899a64/css/report_dialog.css" />



            </div>
</div>


    








<div id="celebrities" class="celebrities related-celebrities">

  
    <h2>
        <i class="">大黄蜂的演职员</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="/subject/26394152/celebrities">全部 61</a>
            )
            </span>
    </h2>


  <ul class="celebrities-list from-subject __oneline">
        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1305796/" title="特拉维斯·奈特 Travis Knight" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1471358307.31.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1305796/" title="特拉维斯·奈特 Travis Knight" class="name">特拉维斯·奈特</a></span>

      <span class="role" title="导演">导演</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1312964/" title="海莉·斯坦菲尔德 Hailee Steinfeld" class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p20419.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1312964/" title="海莉·斯坦菲尔德 Hailee Steinfeld" class="name">海莉·斯坦菲尔德</a></span>

      <span class="role" title="饰 夏琳·沃森 Charlie Watson">饰 夏琳·沃森 Charlie Watson</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1376970/" title="小豪尔赫·兰登伯格 Jorge Lendeborg Jr." class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1545624925.39.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1376970/" title="小豪尔赫·兰登伯格 Jorge Lendeborg Jr." class="name">小豪尔赫·兰登伯格</a></span>

      <span class="role" title="饰 梅莫 Memo">饰 梅莫 Memo</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1044883/" title="约翰·塞纳 John Cena" class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p23477.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1044883/" title="约翰·塞纳 John Cena" class="name">约翰·塞纳</a></span>

      <span class="role" title="饰 特工伯恩斯 Agent Burns">饰 特工伯恩斯 Agent Burns</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1361062/" title="杰森·德鲁克 Jason Drucker" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/pkgttz5tui54cel_avatar_uploaded1471074955.45.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1361062/" title="杰森·德鲁克 Jason Drucker" class="name">杰森·德鲁克</a></span>

      <span class="role" title="饰 奥蒂斯 Otis">饰 奥蒂斯 Otis</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1143044/" title="帕梅拉·阿德龙 Pamela Adlon" class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p21887.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1143044/" title="帕梅拉·阿德龙 Pamela Adlon" class="name">帕梅拉·阿德龙</a></span>

      <span class="role" title="饰 夏琳母亲 Charlie’s Mother">饰 夏琳母亲 Charlie’s Mother</span>

    </div>
  </li>


  </ul>
</div>


    


<link rel="stylesheet" href="https://img3.doubanio.com/f/verify/16c7e943aee3b1dc6d65f600fcc0f6d62db7dfb4/entry_creator/dist/author_subject/style.css">
<div id="author_subject" class="author-wrapper">
    <div class="loading"></div>
</div>
<script type="text/javascript">
    var answerObj = {
      ISALL: 'False',
      TYPE: 'movie',
      SUBJECT_ID: '26394152',
      USER_ID: '54106750'
    }
</script>
<script type="text/javascript" src="https://img3.doubanio.com/f/movie/61252f2f9b35f08b37f69d17dfe48310dd295347/js/movie/lib/react/15.4/bundle.js"></script>
<script type="text/javascript" src="https://img3.doubanio.com/f/verify/ac140ef86262b845d2be7b859e352d8196f3f6d4/entry_creator/dist/author_subject/index.js"></script>


    









    
    <div id="related-pic" class="related-pic">
        
    
    
    <h2>
        <i class="">大黄蜂的视频和图片</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/26394152/trailer#trailer">预告片53</a>&nbsp;|&nbsp;<a href="https://movie.douban.com/subject/26394152/trailer#short_video">视频评论2</a>&nbsp;·&nbsp;<a href="/video/create?subject_id=26394152">添加</a>&nbsp;|&nbsp;<a href="https://movie.douban.com/subject/26394152/all_photos">图片351</a>&nbsp;·&nbsp;<a href="https://movie.douban.com/subject/26394152/mupload">添加</a>
            )
            </span>
    </h2>


        <ul class="related-pic-bd  wide_videos">
                <li class="label-trailer">
                    <a class="related-pic-video" href="https://movie.douban.com/trailer/241374/#content" title="预告片" style="background-image:url(https://img3.doubanio.com/img/trailer/medium/2543888505.jpg?1545993957)">
                    </a>
                </li>
                
                <li class="label-short-video">
                    <a class="related-pic-video" href="https://movie.douban.com/video/102071/" title="视频评论" style="background-image:url(https://img3.doubanio.com/view/photo/photo/public/p2545684033.webp?)">
                    </a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2542656035/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p2542656035.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2524063732/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p2524063732.webp" alt="图片" /></a>
                </li>
        </ul>
    </div>




    
    



<style type="text/css">
.award li { display: inline; margin-right: 5px }
.awards { margin-bottom: 20px }
.awards h2 { background: none; color: #000; font-size: 14px; padding-bottom: 5px; margin-bottom: 8px; border-bottom: 1px dashed #dddddd }
.awards .year { color: #666666; margin-left: -5px }
.mod { margin-bottom: 25px }
.mod .hd { margin-bottom: 10px }
.mod .hd h2 {margin:24px 0 3px 0}
</style>



    








    <div id="recommendations" class="">
        
        
    <h2>
        <i class="">喜欢这部电影的人也喜欢</i>
              · · · · · ·
    </h2>

        
    
    <div class="recommendations-bd">
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/3168101/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2537158013.webp" alt="毒液：致命守护者" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/3168101/?from=subject-page" class="" >毒液：致命守护者</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/24773958/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2517753454.webp" alt="复仇者联盟3：无限战争" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/24773958/?from=subject-page" class="" >复仇者联盟3：无限战争</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/4920389/?from=subject-page" >
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2516578307.webp" alt="头号玩家" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/4920389/?from=subject-page" class="" >头号玩家</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/25820460/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2332503406.webp" alt="美国队长3" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/25820460/?from=subject-page" class="" >美国队长3</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26636712/?from=subject-page" >
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2529389608.webp" alt="蚁人2：黄蜂女现身" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26636712/?from=subject-page" class="" >蚁人2：黄蜂女现身</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/1432146/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p725871004.webp" alt="钢铁侠" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/1432146/?from=subject-page" class="" >钢铁侠</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/24753477/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2497756471.webp" alt="蜘蛛侠：英雄归来" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/24753477/?from=subject-page" class="" >蜘蛛侠：英雄归来</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/1794171/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p1188042816.webp" alt="变形金刚" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/1794171/?from=subject-page" class="" >变形金刚</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/25786060/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2352321614.webp" alt="X战警：天启" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/25786060/?from=subject-page" class="" >X战警：天启</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/25765735/?from=subject-page" >
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2431980130.webp" alt="金刚狼3：殊死一战" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/25765735/?from=subject-page" class="" >金刚狼3：殊死一战</a>
            </dd>
        </dl>
    </div>

    </div>



        


<script type="text/x-handlebar-tmpl" id="comment-tmpl">
    <div class="dummy-fold">
        {{#each comments}}
        <div class="comment-item" data-cid="id">
            <div class="comment">
                <h3>
                    <span class="comment-vote">
                            <span class="votes">{{votes}}</span>
                        <input value="{{id}}" type="hidden"/>
                        <a href="javascript:;" class="j {{#if ../if_logined}}a_vote_comment{{else}}a_show_login{{/if}}">有用</a>
                    </span>
                    <span class="comment-info">
                        <a href="{{user.path}}" class="">{{user.name}}</a>
                        {{#if rating}}
                        <span class="allstar{{rating}}0 rating" title="{{rating_word}}"></span>
                        {{/if}}
                        <span>
                            {{time}}
                        </span>
                        <p> {{content_tmpl content}} </p>
                    </span>
            </div>
        </div>
        {{/each}}
    </div>
</script>












    

    <div id="comments-section">
        <div class="mod-hd">
            
        <a class="comment_btn j a_collect_btn" name="cbtn-26394152" href="javascript:;" rel="nofollow">
            <span>我要写短评</span>
        </a>

            
            
    <h2>
        <i class="">大黄蜂的短评</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/26394152/comments?status=P">全部 49223 条</a>
            )
            </span>
    </h2>

        </div>
        <div class="mod-bd">
                
    <div class="tab-hd">
        <a id="hot-comments-tab" href="comments" data-id="hot" class="on">热门</a>&nbsp;/&nbsp;
        <a id="new-comments-tab" href="comments?sort=time" data-id="new">最新</a>&nbsp;/&nbsp;
        <a id="following-comments-tab" href="follows_comments" data-id="following" >好友</a>
    </div>

    <div class="tab-bd">
        <div id="hot-comments" class="tab">
            
    
        
        <div class="comment-item" data-cid="1580657493">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">5306</span>
                <input value="1580657493" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/lingrui1995/" class="">凌睿</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2018-12-20 09:55:09">
                    2018-12-20
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">ofo小黄车濒临死亡之际，查莉继承杨永信衣钵，使用“电击疗法”使其死而复生，和TF男孩（Transformers）上演变形金刚版《水形物语》。爱如潮水，它将你我包围。
导演：幸亏拍得不错，不用回去继承亿万家产了。</span>
                
                <a class="source-icon" href="https://www.douban.com/doubanapp/" target="_blank"><img src="https://img3.doubanio.com/f/shire/f62b2d2de3fc4a56d176b01cc3bbd47d2681fb38/pics/comment/android.png" title="发自Android" alt="Android" rel="v:image"/></a>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1567862725">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">1146</span>
                <input value="1567862725" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/3540441/" class="">同志亦凡人中文站</a>
                    <span>看过</span>
                    <span class="allstar30 rating" title="还行"></span>
                <span class="comment-time " title="2018-12-12 15:59:54">
                    2018-12-12
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">没有了卖拷贝的狂轰乱炸后，特效和剧情都节制了许多。能把变形金刚电影拍成少女成长日记，派拉蒙和腾讯影业是卯足了劲想收割女性观众群啊。Bumblebee真的是太萌了太萌了太萌了，到最后已经分不清是大黄蜂还是小黄人了，试问哪个妹子不想拥有这样一只乖巧的赛博宠物呢~~~</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1611249285">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">575</span>
                <input value="1611249285" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/helizhanglao/" class="">河狸</a>
                    <span>看过</span>
                    <span class="allstar50 rating" title="力荐"></span>
                <span class="comment-time " title="2019-01-04 16:34:54">
                    2019-01-04
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">撸猫撸狗都弱爆了，撸铁的妹子才是人生赢家。</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1588789020">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">201</span>
                <input value="1588789020" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/AlohA0407/" class="">荆棘</a>
                    <span>看过</span>
                    <span class="allstar20 rating" title="较差"></span>
                <span class="comment-time " title="2018-12-24 06:18:38">
                    2018-12-24
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">果然还是魔弦传说的那个少爷导演…但是去电影院看变形金刚的人想看一个半小时的感情戏然后打10分钟？！不懂为什么有一颗文艺心的要拍这种题材，然后用不美的女主和没完没了的文戏折磨观众。说实话之前大黄蜂的适度萌很讨喜，但过度萌就真的有点腻歪。</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1571773220">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">159</span>
                <input value="1571773220" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/dreamfox/" class="">乌鸦火堂</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2018-12-15 00:31:50">
                    2018-12-15
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">首映，4星，变形金刚的壳，铁巨人的核，新导演风格大变，几乎成为萌宠电影，温情和诙谐为主，大黄蜂卖萌+功夫高手，打斗不多算养眼，远谈不上狂轰滥炸，就一个小品剧。但拍得还不错，耐克哥是真爱粉，就为了赛博坦之战的CG，就为擎天柱和多位汽车人霸天虎的G1造型，还有那首《the touch》，情怀燃爆，比卖拷贝那堆破铜烂铁强太多！加分！</span>
        </p>
    </div>

        </div>



                
                &gt; <a href="comments?sort=new_score&status=P" >更多短评49223条</a>
        </div>
        <div id="new-comments" class="tab">
            <div id="normal">
            </div>
            <div class="fold-hd hide">
                <a class="qa" href="/help/opinion#t2-q0" target="_blank">为什么被折叠？</a>
                <a class="btn-unfold" href="#">有一些短评被折叠了</a>
                <div class="qa-tip">
                    评论被折叠，是因为发布这条评论的帐号行为异常。评论仍可以被展开阅读，对发布人的账号不造成其他影响。如果认为有问题，可以<a href="https://help.douban.com/help/ask?category=movie">联系</a>豆瓣电影。
                </div>
            </div>
            <div class="fold-bd">
            </div>
            <span id="total-num"></span>
        </div>
        <div id="following-comments" class="tab">
            
    
        
        <div class="comment-item" data-cid="1632078797">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">0</span>
                <input value="1632078797" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/62054075/" class="">是灼灼</a>
                    <span>看过</span>
                    <span class="allstar30 rating" title="还行"></span>
                <span class="comment-time " title="2019-01-17 23:52:48">
                    2019-01-17
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">？？？怎么回事？核心记忆芯片不是受损了？自动修复了？女主怎么知道点击大黄蜂的哪个部位就可以救活它？擎天柱竟然是个大货车？？？（无知如我）</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1629065386">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">0</span>
                <input value="1629065386" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/lff121/" class="">l ǐ</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2019-01-15 18:22:40">
                    2019-01-15
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">试问谁不想要一个会卖萌的大黄蜂呢。</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1613159102">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">7</span>
                <input value="1613159102" type="hidden"/>
                <a href="javascript:;" class="j a_vote_comment" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/itzhaoxiangyu/" class="">心生</a>
                    <span>看过</span>
                    <span class="allstar30 rating" title="还行"></span>
                <span class="comment-time " title="2019-01-05 16:55:50">
                    2019-01-05
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">7分。这次的大黄蜂可爱到了新高度，蹲着坐着的时候更是各种萌。成长之路上的少女与重振之路上的汽车人，很好的陪伴与互补。赛博坦之战虽短，但还是看得激动。人类男女主都表现可以，但是配角John Cena的表演却是一脸的笨拙，看着都糟心。PS：片尾还是看到了那个熟悉的名字——迈克尔·贝（制片人）。</span>
        </p>
    </div>

        </div>



        </div>
    </div>


            
            
        </div>
    </div>



        

<link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/73ed658484f98d44.css">

<section class="topics mod">
    <header>
        <h2>
            大黄蜂的话题 · · · · · ·
            <span class="pl">( <span class="gallery_topics">全部 <span id="topic-count"></span> 条</span> )</span>
        </h2>
    </header>

    




<section class="subject-topics">
    <div class="topic-guide" id="topic-guide">
        <img class="ic_question" src="//img3.doubanio.com/f/ithildin/b1a3edea3d04805f899e9d77c0bfc0d158df10d5/pics/export/icon_question.png">
        <div class="tip_content">
            <div class="tip_title">什么是话题</div>
            <div class="tip_desc">
                <div>无论是一部作品、一个人，还是一件事，都往往可以衍生出许多不同的话题。将这些话题细分出来，分别进行讨论，会有更多收获。</div>
            </div>
        </div>
        <img class="ic_guide" src="//img3.doubanio.com/f/ithildin/529f46d86bc08f55cd0b1843d0492242ebbd22de/pics/export/icon_guide_arrow.png">
        <img class="ic_close" id="topic-guide-close" src="//img3.doubanio.com/f/ithildin/2eb4ad488cb0854644b23f20b6fa312404429589/pics/export/close@3x.png">
    </div>

    <div id="topic-items"></div>

    <script>
        window.subject_id = 26394152;
        window.join_label_text = '写影评参与';

        window.topic_display_count = 4;
        window.topic_item_display_count = 1;
        window.no_content_fun_call_name = "no_topic";

        window.guideNode = document.getElementById('topic-guide');
        window.guideNodeClose = document.getElementById('topic-guide-close');
    </script>
    
        <link rel="stylesheet" href="https://img3.doubanio.com/f/ithildin/f731c9ea474da58c516290b3a6b1dd1237c07c5e/css/export/subject_topics.css">
        <script src="https://img3.doubanio.com/f/ithildin/d3590fc6ac47b33c804037a1aa7eec49075428c8/js/export/moment-with-locales-only-zh.js"></script>
        <script src="https://img3.doubanio.com/f/ithildin/c600fdbe69e3ffa5a3919c81ae8c8b4140e99a3e/js/export/subject_topics.js"></script>

</section>

    <script>
        function no_topic(){
            $('#content .topics').remove();
        }
    </script>
</section>

<section class="reviews mod movie-content">
    <header>
        <a href="new_review" rel="nofollow" class="create-review comment_btn"
            data-isverify="True"
            data-verify-url="https://www.douban.com/accounts/phone/verify?redir=http://movie.douban.com/subject/26394152/new_review">
            <span>我要写影评</span>
        </a>
        <h2>
            大黄蜂的影评 · · · · · ·
            <span class="pl">( <a href="reviews">全部 1027 条</a> )</span>
        </h2>
    </header>

    

<style>
#gallery-topics-selection {
  position: fixed;
  width: 595px;
  padding: 40px 40px 33px 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 16px 0 rgba(0, 0, 0, 0.2);
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  z-index: 9999;
}
#gallery-topics-selection h1 {
  font-size: 18px;
  color: #007722;
  margin-bottom: 36px;
  padding: 0;
  line-height: 28px;
  font-weight: normal;
}
#gallery-topics-selection .gl_topics {
  border-bottom: 1px solid #dfdfdf;
  max-height: 298px;
  overflow-y: scroll;
}
#gallery-topics-selection .topic {
  margin-bottom: 24px;
}
#gallery-topics-selection .topic_name {
  font-size: 15px;
  color: #333;
  margin: 0;
  line-height: inherit;
}
#gallery-topics-selection .topic_meta {
  font-size: 13px;
  color: #999;
}
#gallery-topics-selection .topics_skip {
  display: block;
  cursor: pointer;
  font-size: 16px;
  color: #3377AA;
  text-align: center;
  margin-top: 33px;
}
#gallery-topics-selection .topics_skip:hover {
  background: transparent;
}
#gallery-topics-selection .close_selection {
  position: absolute;
  width: 30px;
  height: 20px;
  top: 46px;
  right: 40px;
  background: #fff;
  color: #999;
  text-align: right;
}
#gallery-topics-selection .close_selection:hover{
  background: #fff;
  color: #999;
}
</style>




        <div class="review_filter">
            <a href="javascript:;;" class="cur" data-sort="">热门</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="time">最新</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="follow">好友</a href="javascript:;;">
            
        </div>


        



<div class="review-list  ">
        
    

        
    
    <div data-cid="9870822">
        <div class="main review-item" id="9870822">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/dreamfox/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u2297669-12.jpg">
        </a>

        <a href="https://www.douban.com/people/dreamfox/" class="name">乌鸦火堂</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-04" class="main-meta">2019-01-04 03:03:39</span>

            <a class="rel-topic" target="_blank" href="//www.douban.com/gallery/topic/《大黄蜂》中有哪些致敬经典G1动画的细节？">#《大黄蜂》中有哪些致敬经典G1动画的细节？</a>

    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9870822/">《大黄蜂》疑难解惑+设定梗+全彩蛋终极整理</a></h2>

                <div id="review_9870822_short" class="review-short" data-rid="9870822">
                    <div class="short-content">

                        情怀不能当饭吃。但没有情怀，有时候也会味同嚼蜡。 从2007年“爆炸贝”的《变形金刚》，到2019年的《大黄蜂》（虽然人家在美国是2018年上映，但我们大多数还是在2019年看的），12年，真人版《变形金刚》走过了“一个生肖轮回”。 我还记得07年在大银幕上看到《变形金刚》电影...

                        &nbsp;(<a href="javascript:;" id="toggle-9870822-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9870822_full" class="hidden">
                    <div id="review_9870822_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9870822" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9870822">
                                689
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9870822" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9870822">
                                59
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9870822/#comments" class="reply">146回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9871460">
        <div class="main review-item" id="9871460">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/lolalireader/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u4682629-6.jpg">
        </a>

        <a href="https://www.douban.com/people/lolalireader/" class="name">Sylvia</a>

            <span class="allstar20 main-title-rating" title="较差"></span>

        <span content="2019-01-04" class="main-meta">2019-01-04 13:36:53</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9871460/">“通人性”而非“有人性”的主角是无法打动观众的啊。</a></h2>

                <div id="review_9871460_short" class="review-short" data-rid="9871460">
                    <div class="short-content">

                        近年来好莱坞大片里拍的非人类都被拍的要么像狗要么像猫。 反正就都是犯蠢卖萌这一卦的，永远是“通人性”而不是“有人性”。 从驯龙高手里的龙，再到超能陆战队的大白，再到这部大黄蜂，甚至银河护卫队里的groot，都是如此。（超能陆战队或许稍好一些） 好听点叫人类的伙伴，...

                        &nbsp;(<a href="javascript:;" id="toggle-9871460-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9871460_full" class="hidden">
                    <div id="review_9871460_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9871460" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9871460">
                                327
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9871460" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9871460">
                                63
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9871460/#comments" class="reply">131回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9883223">
        <div class="main review-item" id="9883223">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/51665133/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u51665133-140.jpg">
        </a>

        <a href="https://www.douban.com/people/51665133/" class="name">你大豪爷</a>

            <span class="allstar30 main-title-rating" title="还行"></span>

        <span content="2019-01-08" class="main-meta">2019-01-08 23:48:10</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9883223/">大黄蜂得亏没坠落在中国</a></h2>

                <div id="review_9883223_short" class="review-short" data-rid="9883223">
                    <div class="short-content">

                        一个设想： 如果1987年大黄蜂从塞伯坦逃亡到地球没有落在加州海岸，而是北京西城区一个胡同里，附在一个叫李东宝的出租车司机的黄面的上，那么这会是一部什么气质的影片？ 12月某一天凌晨两点，一辆黄面的在空旷寂静的二环上匀速开着，司机李东宝双手捧着热茶杯坐在驾驶位上一...

                        &nbsp;(<a href="javascript:;" id="toggle-9883223-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9883223_full" class="hidden">
                    <div id="review_9883223_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9883223" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9883223">
                                239
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9883223" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9883223">
                                11
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9883223/#comments" class="reply">45回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9871657">
        <div class="main review-item" id="9871657">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/127299151/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u127299151-3.jpg">
        </a>

        <a href="https://www.douban.com/people/127299151/" class="name">绿毛水怪</a>

            <span class="allstar30 main-title-rating" title="还行"></span>

        <span content="2019-01-04" class="main-meta">2019-01-04 15:11:00</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9871657/">《大黄蜂》拍成《哈士奇》，美国王思聪可以考虑回家继承耐克了</a></h2>

                <div id="review_9871657_short" class="review-short" data-rid="9871657">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇影评可能有剧透</p>

                        《变形金刚》自从1988年把动画版引入中国以后，成为很多中国小朋友的童年回忆。我小时候不看《变形金刚》，身边也有很多人在谈论什么“汽车人变身”、“擎天柱”、“霸天虎”，后来甚至大黄蜂同款跑车都火了。在我那个年代，哪个小朋友家里没有一个变形金刚玩具都不好意思出来...

                        &nbsp;(<a href="javascript:;" id="toggle-9871657-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9871657_full" class="hidden">
                    <div id="review_9871657_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9871657" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9871657">
                                381
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9871657" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9871657">
                                73
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9871657/#comments" class="reply">267回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9871978">
        <div class="main review-item" id="9871978">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/loyoyo_O/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/u67858795-18.jpg">
        </a>

        <a href="https://www.douban.com/people/loyoyo_O/" class="name">木由</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2019-01-04" class="main-meta">2019-01-04 17:41:07</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9871978/">有史以来最温情的变形金刚和一些细节思考</a></h2>

                <div id="review_9871978_short" class="review-short" data-rid="9871978">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇影评可能有剧透</p>

                        先给大家讲一个美国版田螺姑娘的故事。 话说Charlie小妹妹是一个瓷实的姑娘，体格倍儿棒，性格爽朗，从小跟着爸爸修车，跳水，也曾是个被宠爱的小公主。可爸爸心脏病突发去世，快乐的生活被打烂。妈妈另结新欢，还给她添了个小弟弟。 这天Charlie十八岁生日，在家里委屈地接受...

                        &nbsp;(<a href="javascript:;" id="toggle-9871978-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9871978_full" class="hidden">
                    <div id="review_9871978_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9871978" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9871978">
                                140
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9871978" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9871978">
                                22
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9871978/#comments" class="reply">33回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9865966">
        <div class="main review-item" id="9865966">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/mr_tree/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u60689518-41.jpg">
        </a>

        <a href="https://www.douban.com/people/mr_tree/" class="name">凹凸</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-01" class="main-meta">2019-01-01 23:28:49</span>

            <a class="rel-topic" target="_blank" href="//www.douban.com/gallery/topic/《大黄蜂》有哪些值得推荐的看点？">#《大黄蜂》有哪些值得推荐的看点？</a>

    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9865966/">接棒“卖拷贝”的“美国王思聪”拍出了《变形金刚》系列最佳</a></h2>

                <div id="review_9865966_short" class="review-short" data-rid="9865966">
                    <div class="short-content">

                        未经授权，严禁转载！！！ 1984年，美国玩具厂商HASBRO公司，从日本TAKARA公司收购了可变成机器人的原模合金汽车和飞机，起名为“变形金刚”，并制作了同名动画在美国各大电视台播出，以吸引更多的小朋友购买这些玩具。恐怕连HASBRO公司也没想到，“变形金刚”会受到如此大的欢...

                        &nbsp;(<a href="javascript:;" id="toggle-9865966-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9865966_full" class="hidden">
                    <div id="review_9865966_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9865966" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9865966">
                                151
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9865966" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9865966">
                                53
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9865966/#comments" class="reply">94回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9855770">
        <div class="main review-item" id="9855770">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/N.B./" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u1473794-226.jpg">
        </a>

        <a href="https://www.douban.com/people/N.B./" class="name">无非🏳️🌈</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2018-12-28" class="main-meta">2018-12-28 18:16:12</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9855770/">《大黄蜂》与科幻电影对古典神话母题的再生产</a></h2>

                <div id="review_9855770_short" class="review-short" data-rid="9855770">
                    <div class="short-content">

                        大概所有人都还记得《变形金刚》真人版在2007年被迈克尔贝启动时口碑和票房双赢的空前盛况。不过在经历了5部的狂轰乱炸后，《大黄蜂》这部番外似乎是派拉蒙寻求新方向的调整和试水之作。与之前的机甲酣战不同，虽然不乏博派与狂派的星际大乱斗，但《大黄蜂》更像一部以人类和异...

                        &nbsp;(<a href="javascript:;" id="toggle-9855770-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9855770_full" class="hidden">
                    <div id="review_9855770_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9855770" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9855770">
                                122
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9855770" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9855770">
                                21
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9855770/#comments" class="reply">12回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9866325">
        <div class="main review-item" id="9866325">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/166472455/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u166472455-1.jpg">
        </a>

        <a href="https://www.douban.com/people/166472455/" class="name">Mr.White</a>


        <span content="2019-01-02" class="main-meta">2019-01-02 01:37:18</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9866325/">《大黄蜂》含剧透吐槽</a></h2>

                <div id="review_9866325_short" class="review-short" data-rid="9866325">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇影评可能有剧透</p>

                        我是一个自认为观影品味不差的影迷，也是钢丝，下面作为这双重身份混杂胡乱吐槽一下《大黄蜂》这部电影（以下吐槽不理性不客观不中立，来自一个只会抬杠不会写影评的傻逼）： 1.开场赛博坦之战还不如不要。且不论赛博坦场景那游戏质感的CG，这场大战戏的调度混乱程度堪称灾难，...

                        &nbsp;(<a href="javascript:;" id="toggle-9866325-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9866325_full" class="hidden">
                    <div id="review_9866325_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9866325" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9866325">
                                129
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9866325" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9866325">
                                18
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9866325/#comments" class="reply">72回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9872771">
        <div class="main review-item" id="9872771">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/178131964/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u178131964-2.jpg">
        </a>

        <a href="https://www.douban.com/people/178131964/" class="name">影探</a>

            <span class="allstar30 main-title-rating" title="还行"></span>

        <span content="2019-01-04" class="main-meta">2019-01-04 23:57:35</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9872771/">你所不知道的大黄蜂：你永远都是我的第一辆车</a></h2>

                <div id="review_9872771_short" class="review-short" data-rid="9872771">
                    <div class="short-content">

                        首发于公众号“影探”ID：ttyingtan 微博：影探探长 作者：探长 转载请注明出处 2007年，面对自动打开的车门，一个男孩对身边犹豫要不要上车的女生说： “50年后，当你回望今生，你会不会后悔今天没胆上这辆车?” 然后两人对视2秒，一同上了那辆黄色跑车。 两人就此从陌生到熟...

                        &nbsp;(<a href="javascript:;" id="toggle-9872771-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9872771_full" class="hidden">
                    <div id="review_9872771_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9872771" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9872771">
                                76
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9872771" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9872771">
                                8
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9872771/#comments" class="reply">12回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9871000">
        <div class="main review-item" id="9871000">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/firo/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/u1799384-69.jpg">
        </a>

        <a href="https://www.douban.com/people/firo/" class="name">银谷</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-04" class="main-meta">2019-01-04 09:21:29</span>

            <a class="rel-topic" target="_blank" href="//www.douban.com/gallery/topic/《大黄蜂》在变形金刚系列电影中有哪些独特之处？">#《大黄蜂》在变形金刚系列电影中有哪些独特之处？</a>

    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/9871000/">变形金刚版ET，这是史上最温情，最有人情味的一部变形金刚。</a></h2>

                <div id="review_9871000_short" class="review-short" data-rid="9871000">
                    <div class="short-content">

                        在经历了迈克尔贝几部无脑视效疲劳轰炸的变形金刚后，耐克贵公子奈特是如何让这一系列焕发生机的？看了点映后发现，《大黄蜂》北美上映后烂番茄新鲜度高度93%不是没有原因的。 特拉维斯·奈特是一个非常有想法的导演，之前担任美术设计和制片人的高口碑动画《僵尸新娘》《鬼妈...

                        &nbsp;(<a href="javascript:;" id="toggle-9871000-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9871000_full" class="hidden">
                    <div id="review_9871000_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9871000" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9871000">
                                51
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9871000" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9871000">
                                16
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/9871000/#comments" class="reply">18回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>




    

    

    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/60bd61e3276b01f8.js"></script>
    <!-- COLLECTED CSS -->
</div>








            <p class="pl">
                &gt;
                <a href="reviews">
                    更多影评1027篇
                </a>
            </p>
</section>

<!-- COLLECTED JS -->

    <br/>

        <div class="section-discussion">
                
                <div class="mod-hd">
                        <a class="comment_btn" href="/subject/26394152/discussion/create" rel="nofollow"><span>添加新讨论</span></a>
                    
    <h2>
        讨论区
         &nbsp; &middot;&nbsp; &middot;&nbsp; &middot;&nbsp; &middot;&nbsp; &middot;&nbsp; &middot;
    </h2>

                </div>
                
  <table class="olt"><tr><td><td><td><td></tr>
        
        <tr>
          <td class="pl"><a href="https://movie.douban.com/subject/26394152/discussion/615944163/" title="好像有资源了">好像有资源了</a></td>
          <td class="pl"><span>来自</span><a href="https://www.douban.com/people/madfishy1/">白白</a></td>
          <td class="pl"><span>4 回应</span></td>
          <td class="pl"><span>2019-01-20</span></td>
        </tr>
        
        <tr>
          <td class="pl"><a href="https://movie.douban.com/subject/26394152/discussion/615945069/" title="转发这个爆痘少女，下一个杨超越就是你！">转发这个爆痘少女，下一个杨超越就是你！</a></td>
          <td class="pl"><span>来自</span><a href="https://www.douban.com/people/1786804/">XD|醒来。你在。</a></td>
          <td class="pl"><span></span></td>
          <td class="pl"><span>2019-01-20</span></td>
        </tr>
        
        <tr>
          <td class="pl"><a href="https://movie.douban.com/subject/26394152/discussion/615944566/" title="没有人觉得这部电影BUG还是很多的么？比如：">没有人觉得这部电影BUG还是很多的么？比如：</a></td>
          <td class="pl"><span>来自</span><a href="https://www.douban.com/people/175343924/">yeafine</a></td>
          <td class="pl"><span>2 回应</span></td>
          <td class="pl"><span>2019-01-20</span></td>
        </tr>
        
        <tr>
          <td class="pl"><a href="https://movie.douban.com/subject/26394152/discussion/615937689/" title="到底好不好看">到底好不好看</a></td>
          <td class="pl"><span>来自</span><a href="https://www.douban.com/people/150369352/">RUH</a></td>
          <td class="pl"><span>17 回应</span></td>
          <td class="pl"><span>2019-01-20</span></td>
        </tr>
        
        <tr>
          <td class="pl"><a href="https://movie.douban.com/subject/26394152/discussion/615944472/" title="有人扒一下女主的衣服吗？我觉得好好看呀！">有人扒一下女主的衣服吗？我觉得好好看呀！</a></td>
          <td class="pl"><span>来自</span><a href="https://www.douban.com/people/173077711/">-卿羽-</a></td>
          <td class="pl"><span>3 回应</span></td>
          <td class="pl"><span>2019-01-20</span></td>
        </tr>
  </table>

                <p class="pl" align="right">
                    <a href="/subject/26394152/discussion/" rel="nofollow">
                        &gt; 去这部影片的讨论区（全部372条）
                    </a>
                </p>
        </div>

        
    
        
                





<div id="askmatrix">
    <div class="mod-hd">
        <h2>
            关于《大黄蜂》的问题
            · · · · · ·
            <span class="pl">
                (<a href='https://movie.douban.com/subject/26394152/questions/?from=subject'>
                    全部10个
                </a>)
            </span>
        </h2>


        
    
    <a class=' comment_btn'
        href='https://movie.douban.com/subject/26394152/questions/ask/?from=subject'>我来提问</a>

    </div>

    <div class="mod-bd">
        <ul class="">
            <li class="">
                <span class="tit">
                    <a href="https://movie.douban.com/subject/26394152/questions/821035/?from=subject" class="">
                            这里面有多少水军？
                    </a>
                </span>
                <span class="meta">
                    5人回答
                </span>
            </li>
            <li class="">
                <span class="tit">
                    <a href="https://movie.douban.com/subject/26394152/questions/820661/?from=subject" class="">
                            大黄蜂的雪佛兰科迈罗？
                    </a>
                </span>
                <span class="meta">
                    3人回答
                </span>
            </li>
        </ul>

        <p>&gt;
            <a href='https://movie.douban.com/subject/26394152/questions/?from=subject'>
                全部10个问题
            </a>
        </p>

    </div>
</div>



            


    <script type="text/javascript">
        $(function(){if($.browser.msie && $.browser.version == 6.0){
            var $info = $('#info'),
                maxWidth = parseInt($info.css('max-width'));
            if($info.width() > maxWidth) {
                $info.width(maxWidth);
            }
        }});
    </script>


            </div>
            <div class="aside">
                


    







            






<div class="ticket">
        <a class="ticket-btn" href="https://movie.douban.com/ticket/redirect/?url=https%3A%2F%2Fm.maoyan.com%2Fcinema%2Fmovie%2F1206875%3F_v_%3Dyes%26merCode%3D1000011">购票</a>
</div>



    <!-- douban ad begin -->
    <div id="dale_movie_subject_top_right"></div>
    <div id="dale_movie_subject_top_middle"></div>
    <!-- douban ad end -->

    



<style type="text/css">
    .m4 {margin-bottom:8px; padding-bottom:8px;}
    .movieOnline {background:#FFF6ED; padding:10px; margin-bottom:20px;}
    .movieOnline h2 {margin:0 0 5px;}
    .movieOnline .sitename {line-height:2em; width:160px;}
    .movieOnline td,.movieOnline td a:link,.movieOnline td a:visited{color:#666;}
    .movieOnline td a:hover {color:#fff;}
    .link-bt:link,
    .link-bt:visited,
    .link-bt:hover,
    .link-bt:active {margin:5px 0 0; padding:2px 8px; background:#a8c598; color:#fff; -moz-border-radius: 3px; -webkit-border-radius: 3px; border-radius: 3px; display:inline-block;}
</style>



    







    
    <div class="tags">
        
        
    <h2>
        <i class="">豆瓣成员常用的标签</i>
              · · · · · ·
    </h2>

        <div class="tags-body">
                <a href="/tag/科幻" class="">科幻</a>
                <a href="/tag/美国" class="">美国</a>
                <a href="/tag/温情" class="">温情</a>
                <a href="/tag/动作" class="">动作</a>
                <a href="/tag/超级英雄" class="">超级英雄</a>
                <a href="/tag/漫画改编" class="">漫画改编</a>
                <a href="/tag/2018" class="">2018</a>
                <a href="/tag/成长" class="">成长</a>
        </div>
    </div>


    <div id="dale_movie_subject_inner_middle"></div>
    <div id="dale_movie_subject_download_middle"></div>
        








<div id="subject-doulist">
    
    
    <h2>
        <i class="">以下豆列推荐</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/26394152/doulists">全部</a>
            )
            </span>
    </h2>


    
    <ul>
            <li>
                <a href="https://www.douban.com/doulist/30299/" target="_blank">豆瓣电影【口碑榜】2018-12-20更新</a>
                <span>(影志)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/41896950/" target="_blank">想看的电影太多怕忘了</a>
                <span>(J.D.)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/1504454/" target="_blank">ღ♩♪生活有这些期待很有动力♫♬ღ</a>
                <span>(freedom♪)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/43556971/" target="_blank">始终会看的电影</a>
                <span>(可可)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/267907/" target="_blank">ღ『如何让你遇见我，在我最美丽的时刻。』</a>
                <span>(ღ 狐不悔，)</span>
            </li>
    </ul>

</div>

        








<div id="subject-others-interests">
    
    
    <h2>
        <i class="">谁在看这部电影</i>
              · · · · · ·
    </h2>

    
    <ul class="">
            
            <li class="">
                <a href="https://www.douban.com/people/43980523/" class="others-interest-avatar">
                    <img src="https://img1.doubanio.com/icon/u43980523-19.jpg" class="pil" alt="掮客">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/43980523/" class="">掮客</a>
                    <div class="">
                        刚刚
                        看过
                        <span class="allstar30" title="还行"></span>
                    </div>
                </div>
            </li>
            
            <li class="">
                <a href="https://www.douban.com/people/190310004/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u190310004-1.jpg" class="pil" alt="㎕࿐.">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/190310004/" class="">㎕࿐.</a>
                    <div class="">
                        刚刚
                        想看
                        
                    </div>
                </div>
            </li>
            
            <li class="">
                <a href="https://www.douban.com/people/165279758/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u165279758-1.jpg" class="pil" alt="三九天喝汽水">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/165279758/" class="">三九天喝汽水</a>
                    <div class="">
                        刚刚
                        想看
                        
                    </div>
                </div>
            </li>
    </ul>

    
    <div class="subject-others-interests-ft">
        
            <a href="https://movie.douban.com/subject/26394152/collections">244616人看过</a>
                &nbsp;/&nbsp;
            <a href="https://movie.douban.com/subject/26394152/wishes">193692人想看</a>
    </div>

</div>



    
    

<!-- douban ad begin -->
<div id="dale_movie_subject_middle_right"></div>
<script type="text/javascript">
    (function (global) {
        if(!document.getElementsByClassName) {
            document.getElementsByClassName = function(className) {
                return this.querySelectorAll("." + className);
            };
            Element.prototype.getElementsByClassName = document.getElementsByClassName;

        }
        var articles = global.document.getElementsByClassName('article'),
            asides = global.document.getElementsByClassName('aside');

        if (articles.length > 0 && asides.length > 0 && articles[0].offsetHeight >= asides[0].offsetHeight) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_movie_subject_middle_right');
        }
    })(this);
</script>
<!-- douban ad end -->



    <br/>

    
<p class="pl">订阅大黄蜂的评论: <br/><span class="feed">
    <a href="https://movie.douban.com/feed/subject/26394152/reviews"> feed: rss 2.0</a></span></p>


            </div>
            <div class="extra">
                
    
<!-- douban ad begin -->
<div id="dale_movie_subject_bottom_super_banner"></div>
<script type="text/javascript">
    (function (global) {
        var body = global.document.body,
            html = global.document.documentElement;

        var height = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight, html.scrollHeight, html.offsetHeight);
        if (height >= 2000) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_movie_subject_bottom_super_banner');
        }
    })(this);
</script>
<!-- douban ad end -->


            </div>
        </div>
    </div>

        
    <div id="footer">
            <div class="footer-extra"></div>
        
<span id="icp" class="fleft gray-link">
    &copy; 2005－2019 douban.com, all rights reserved 北京豆网科技有限公司
</span>

<a href="https://www.douban.com/hnypt/variformcyst.py" style="display: none;"></a>

<span class="fright">
    <a href="https://www.douban.com/about">关于豆瓣</a>
    · <a href="https://www.douban.com/jobs">在豆瓣工作</a>
    · <a href="https://www.douban.com/about?topic=contactus">联系我们</a>
    · <a href="https://www.douban.com/about?policy=disclaimer">免责声明</a>
    
    · <a href="https://help.douban.com/?app=movie" target="_blank">帮助中心</a>
    · <a href="https://www.douban.com/doubanapp/">移动应用</a>
    · <a href="https://www.douban.com/partner/">豆瓣广告</a>
</span>

    </div>

    </div>
    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/1958db76bb07bfe6.js"></script>
    
    
    
    




    
<script type="text/javascript">
    (function (global) {
        var newNode = global.document.createElement('script'),
            existingNode = global.document.getElementsByTagName('script')[0],
            adSource = '//erebor.douban.com/',
            userId = '54106750',
            browserId = 'auSsK8Dk5cg',
            criteria = '7:格林·特鲁曼|7:动作|7:约翰·塞纳|7:搞笑|7:斯蒂芬·施耐德|7:萨钦·巴特|7:弗里德·杜莱尔|7:小豪尔赫·兰登伯格|7:约翰·奥提兹|7:漫画改编|7:安东尼奥·查丽蒂|7:海莉·斯坦菲尔德|7:2018|7:兰·卡琉|7:成长|7:里卡多·霍约斯|7:帕梅拉·阿德龙|7:超级英雄|7:杰森·德鲁克|7:特拉维斯·奈特|7:迪伦·奥布莱恩|7:温情|7:拉斯·斯兰德|7:冒险|7:蓝尼·雅各布森 |7:科幻|7:美国|7:青春|7:艾德文·霍德吉|7:梅金·普莱斯|7:格蕾丝·达斯恩妮|7:蒂姆·马丁·格里森|3:/subject/26394152/',
            preview = '',
            debug = false,
            adSlots = ['dale_movie_subject_top_icon', 'dale_movie_subject_top_right', 'dale_movie_subject_top_middle', 'dale_movie_subject_inner_middle', 'dale_movie_subject_download_middle'];

        global.DoubanAdRequest = {src: adSource, uid: userId, bid: browserId, crtr: criteria, prv: preview, debug: debug};
        global.DoubanAdSlots = (global.DoubanAdSlots || []).concat(adSlots);

        newNode.setAttribute('type', 'text/javascript');
        newNode.setAttribute('src', 'https://img3.doubanio.com/f/adjs/dd37385211bc8deb01376096bfa14d2c0436a98c/ad.release.js');
        newNode.setAttribute('async', true);
        existingNode.parentNode.insertBefore(newNode, existingNode);
    })(this);
</script>











    
  









<script type="text/javascript">
var _paq = _paq || [];
_paq.push(['trackPageView']);
_paq.push(['enableLinkTracking']);
(function() {
    var p=(('https:' == document.location.protocol) ? 'https' : 'http'), u=p+'://fundin.douban.com/';
    _paq.push(['setTrackerUrl', u+'piwik']);
    _paq.push(['setSiteId', '100001']);
    var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
    g.type='text/javascript';
    g.defer=true;
    g.async=true;
    g.src=p+'://img3.doubanio.com/dae/fundin/piwik.js';
    s.parentNode.insertBefore(g,s);
})();
</script>

<script type="text/javascript">
var setMethodWithNs = function(namespace) {
  var ns = namespace ? namespace + '.' : ''
    , fn = function(string) {
        if(!ns) {return string}
        return ns + string
      }
  return fn
}

var gaWithNamespace = function(fn, namespace) {
  var method = setMethodWithNs(namespace)
  fn.call(this, method)
}

var _gaq = _gaq || []
  , accounts = [
      { id: 'UA-7019765-1', namespace: 'douban' }
    , { id: 'UA-7019765-19', namespace: '' }
    ]
  , gaInit = function(account) {
      gaWithNamespace(function(method) {
        gaInitFn.call(this, method, account)
      }, account.namespace)
    }
  , gaInitFn = function(method, account) {
      _gaq.push([method('_setAccount'), account.id]);
      _gaq.push([method('_setSampleRate'), '5']);

      
  _gaq.push([method('_addOrganic'), 'google', 'q'])
  _gaq.push([method('_addOrganic'), 'baidu', 'wd'])
  _gaq.push([method('_addOrganic'), 'soso', 'w'])
  _gaq.push([method('_addOrganic'), 'youdao', 'q'])
  _gaq.push([method('_addOrganic'), 'so.360.cn', 'q'])
  _gaq.push([method('_addOrganic'), 'sogou', 'query'])
  if (account.namespace) {
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣'])
    _gaq.push([method('_addIgnoredOrganic'), 'douban'])
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣网'])
    _gaq.push([method('_addIgnoredOrganic'), 'www.douban.com'])
  }

      if (account.namespace === 'douban') {
        _gaq.push([method('_setDomainName'), '.douban.com'])
      }

        _gaq.push([method('_setCustomVar'), 1, 'responsive_view_mode', 'desktop', 3])

        _gaq.push([method('_setCustomVar'), 2, 'login_status', '1', 2]);

      _gaq.push([method('_trackPageview')])
    }

for(var i = 0, l = accounts.length; i < l; i++) {
  var account = accounts[i]
  gaInit(account)
}


;(function() {
    var ga = document.createElement('script');
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    ga.setAttribute('async', 'true');
    document.documentElement.firstChild.appendChild(ga);
})()
</script>








      
    

    <!-- brand10-docker-->

  <script>_SPLITTEST=''</script>
</body>

</html>

`

var html33 = `



<!DOCTYPE html>
<html lang="zh-cmn-Hans" class="ua-windows ua-webkit book-new-nav">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <title>四个春天 (豆瓣)</title>
  
<script>!function(e){var o=function(o,n,t){var c,i,r=new Date;n=n||30,t=t||"/",r.setTime(r.getTime()+24*n*60*60*1e3),c="; expires="+r.toGMTString();for(i in o)e.cookie=i+"="+o[i]+c+"; path="+t},n=function(o){var n,t,c,i=o+"=",r=e.cookie.split(";");for(t=0,c=r.length;t<c;t++)if(n=r[t].replace(/^\s+|\s+$/g,""),0==n.indexOf(i))return n.substring(i.length,n.length).replace(/\"/g,"");return null},t=e.write,c={"douban.com":1,"douban.fm":1,"google.com":1,"google.cn":1,"googleapis.com":1,"gmaptiles.co.kr":1,"gstatic.com":1,"gstatic.cn":1,"google-analytics.com":1,"googleadservices.com":1},i=function(e,o){var n=new Image;n.onload=function(){},n.src="https://www.douban.com/j/except_report?kind=ra022&reason="+encodeURIComponent(e)+"&environment="+encodeURIComponent(o)},r=function(o){try{t.call(e,o)}catch(e){t(o)}},a=/<script.*?src\=["']?([^"'\s>]+)/gi,g=/http:\/\/(.+?)\.([^\/]+).+/i;e.writeln=e.write=function(e){var t,l=a.exec(e);return l&&(t=g.exec(l[1]))?c[t[2]]?void r(e):void("tqs"!==n("hj")&&(i(l[1],location.href),o({hj:"tqs"},1),setTimeout(function(){location.replace(location.href)},50))):void r(e)}}(document);
</script>

  
  <meta http-equiv="Pragma" content="no-cache">
  <meta http-equiv="Expires" content="Sun, 6 Mar 2005 01:00:00 GMT">
  
<meta http-equiv="mobile-agent" content="format=html5; url=https://m.douban.com/book/subject/30389935/">
<meta name="keywords" content="四个春天,陆庆屹,南海出版公司,2019-1-1,简介,作者,书评,论坛,推荐,二手">
<meta name="description" content="图书四个春天 介绍、书评、论坛及推荐 ">

  <script>var _head_start = new Date();</script>
  
  <link href="https://img3.doubanio.com/f/book/e7612a013a5e76c7c680323c74748d21cd703ba0/css/book/master.css" rel="stylesheet" type="text/css">

  <link href="https://img3.doubanio.com/f/book/222a5c61e041638af8defc87cf97f4a863a77922/css/book/base/init.css" rel="stylesheet">
  <style type="text/css"></style>
  <script src="https://img3.doubanio.com/f/book/0495cb173e298c28593766009c7b0a953246c5b5/js/book/lib/jquery/jquery.js"></script>
  <script src="https://img3.doubanio.com/f/book/7d36c07b1b7a7a386c3bff538ae46d9d0f0990a0/js/book/master.js"></script>
  

  
  <link rel="stylesheet" href="https://img3.doubanio.com/f/book/1eabe3f4e416e77dbafd2ef9bc830fa2ac7a8d4c/css/book/subject.css">
  <link href="https://img3.doubanio.com/f/book/5d301503fbbd8e09f3114583859789884e942f47/css/book/annotation/like.css" rel="stylesheet">
  <script src="https://img3.doubanio.com/f/shire/3c6f2946669cfb2fc9ee4a4d1dcc41fc181cad92/js/lib/jquery.snippet.js"></script>
  <script src="https://img3.doubanio.com/f/shire/77323ae72a612bba8b65f845491513ff3329b1bb/js/do.js" data-cfg-autoload="false"></script>
  <script src="https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js"></script>
  <script src="https://img3.doubanio.com/f/book/2e421e5ec8f2869d31535206c0ac0322532be1f8/js/book/mod/hide.js"></script>
  <script src="https://img3.doubanio.com/f/book/cc6b1a77c3812c7dd20b0374332fade081e1c0b0/js/book/subject/unfold.js"></script>
    <link rel="alternate" href="https://book.douban.com/feed/subject/30389935/reviews" type="application/rss+xml" title="RSS">
  <style type="text/css"> h2 {color: #007722;} </style>
  <script type='text/javascript'>
    var _vds = _vds || [];
    (function(){ _vds.push(['setAccountId', '22c937bbd8ebd703f2d8e9445f7dfd03']);
        _vds.push(['setCS1','user_id','0']);
            (function() {var vds = document.createElement('script');
                vds.type='text/javascript';
                vds.async = true;
                vds.src = ('https:' == document.location.protocol ? 'https://' : 'http://') + 'dn-growing.qbox.me/vds.js';
                var s = document.getElementsByTagName('script')[0];
                s.parentNode.insertBefore(vds, s);
            })();
    })();
</script>

  
  <script type='text/javascript'>
    var _vwo_code=(function(){
      var account_id=249272,
          settings_tolerance=2000,
          library_tolerance=2500,
          use_existing_jquery=false,
          // DO NOT EDIT BELOW THIS LINE
          f=false,d=document;return{use_existing_jquery:function(){return use_existing_jquery;},library_tolerance:function(){return library_tolerance;},finish:function(){if(!f){f=true;var a=d.getElementById('_vis_opt_path_hides');if(a)a.parentNode.removeChild(a);}},finished:function(){return f;},load:function(a){var b=d.createElement('script');b.src=a;b.type='text/javascript';b.innerText;b.onerror=function(){_vwo_code.finish();};d.getElementsByTagName('head')[0].appendChild(b);},init:function(){settings_timer=setTimeout('_vwo_code.finish()',settings_tolerance);var a=d.createElement('style'),b='body{opacity:0 !important;filter:alpha(opacity=0) !important;background:none !important;}',h=d.getElementsByTagName('head')[0];a.setAttribute('id','_vis_opt_path_hides');a.setAttribute('type','text/css');if(a.styleSheet)a.styleSheet.cssText=b;else a.appendChild(d.createTextNode(b));h.appendChild(a);this.load('//dev.visualwebsiteoptimizer.com/j.php?a='+account_id+'&u='+encodeURIComponent(d.URL)+'&r='+Math.random());return settings_timer;}};}());_vwo_settings_timer=_vwo_code.init();
  </script>

  


<script type="application/ld+json">
{
  "@context":"http://schema.org",
  "@type":"Book",
  "workExample": [],
  "name" : "四个春天",
  "author": 
  [
    {
      "@type": "Person",
      "name": "陆庆屹"
    }
  ]
,
  "url" : "https://book.douban.com/subject/30389935/",
  "isbn" : "9787544294881",
  "sameAs": "https://book.douban.com/subject/30389935/"
}
</script>


  <script>  </script>
  <link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/5ef2a22fba637133.css">

  <link rel="shortcut icon" href="https://img3.doubanio.com/favicon.ico" type="image/x-icon">
</head>
<body>
  
    <script>var _body_start = new Date();</script>
    
  



    <link href="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.css" rel="stylesheet" type="text/css">



<div id="db-global-nav" class="global-nav">
  <div class="bd">
    
<div class="top-nav-info">
  <a href="https://www.douban.com/accounts/login?source=book" class="nav-login" rel="nofollow">登录</a>
  <a href="https://www.douban.com/accounts/register?source=book" class="nav-register" rel="nofollow">注册</a>
</div>


    <div class="top-nav-doubanapp">
  <a href="https://www.douban.com/doubanapp/app?channel=top-nav" class="lnk-doubanapp">下载豆瓣客户端</a>
  <div id="doubanapp-tip">
    <a href="https://www.douban.com/doubanapp/app?channel=qipao" class="tip-link">豆瓣 <span class="version">6.0</span> 全新发布</a>
    <a href="javascript: void 0;" class="tip-close">×</a>
  </div>
  <div id="top-nav-appintro" class="more-items">
    <p class="appintro-title">豆瓣</p>
    <p class="qrcode">扫码直接下载</p>
    <div class="download">
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=iOS">iPhone</a>
      <span>·</span>
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=Android" class="download-android">Android</a>
    </div>
  </div>
</div>

    


<div class="global-nav-items">
  <ul>
    <li class="">
      <a href="https://www.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-main&quot;,&quot;uid&quot;:&quot;0&quot;}">豆瓣</a>
    </li>
    <li class="on">
      <a href="https://book.douban.com"  data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-book&quot;,&quot;uid&quot;:&quot;0&quot;}">读书</a>
    </li>
    <li class="">
      <a href="https://movie.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-movie&quot;,&quot;uid&quot;:&quot;0&quot;}">电影</a>
    </li>
    <li class="">
      <a href="https://music.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-music&quot;,&quot;uid&quot;:&quot;0&quot;}">音乐</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/location" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-location&quot;,&quot;uid&quot;:&quot;0&quot;}">同城</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/group" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-group&quot;,&quot;uid&quot;:&quot;0&quot;}">小组</a>
    </li>
    <li class="">
      <a href="https://read.douban.com&#47;?dcs=top-nav&amp;dcm=douban" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-read&quot;,&quot;uid&quot;:&quot;0&quot;}">阅读</a>
    </li>
    <li class="">
      <a href="https://douban.fm&#47;?from_=shire_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-fm&quot;,&quot;uid&quot;:&quot;0&quot;}">FM</a>
    </li>
    <li class="">
      <a href="https://time.douban.com&#47;?dt_time_source=douban-web_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-time&quot;,&quot;uid&quot;:&quot;0&quot;}">时间</a>
    </li>
    <li class="">
      <a href="https://market.douban.com&#47;?utm_campaign=douban_top_nav&amp;utm_source=douban&amp;utm_medium=pc_web" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-market&quot;,&quot;uid&quot;:&quot;0&quot;}">豆品</a>
    </li>
    <li>
      <a href="#more" class="bn-more"><span>更多</span></a>
      <div class="more-items">
        <table cellpadding="0" cellspacing="0">
          <tbody>
            <tr>
              <td>
                <a href="https://ypy.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-ypy&quot;,&quot;uid&quot;:&quot;0&quot;}">豆瓣摄影</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </li>
  </ul>
</div>

  </div>
</div>
<script>
  ;window._GLOBAL_NAV = {
    DOUBAN_URL: "https://www.douban.com",
    N_NEW_NOTIS: 0,
    N_NEW_DOUMAIL: 0
  };
</script>



    <script src="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.js" defer="defer"></script>




  



    <link href="//img3.doubanio.com/dae/accounts/resources/8c80301/book/bundle.css" rel="stylesheet" type="text/css">




<div id="db-nav-book" class="nav">
  <div class="nav-wrap">
  <div class="nav-primary">
    <div class="nav-logo">
      <a href="https:&#47;&#47;book.douban.com">豆瓣读书</a>
    </div>
    <div class="nav-search">
      <form action="https:&#47;&#47;book.douban.com/subject_search" method="get">
        <fieldset>
          <legend>搜索：</legend>
          <label for="inp-query">
          </label>
          <div class="inp"><input id="inp-query" name="search_text" size="22" maxlength="60" placeholder="书名、作者、ISBN" value=""></div>
          <div class="inp-btn"><input type="submit" value="搜索"></div>
          <input type="hidden" name="cat" value="1001" />
        </fieldset>
      </form>
    </div>
  </div>
  </div>
  <div class="nav-secondary">
    

<div class="nav-items">
  <ul>
    <li    ><a href="https://book.douban.com/cart/"
     >购书单</a>
    </li>
    <li    ><a href="https://read.douban.com/ebooks/?dcs=book-nav&dcm=douban"
            target="_blank"
     >电子图书</a>
    </li>
    <li    ><a href="https://market.douban.com/book?utm_campaign=book_nav_freyr&utm_source=douban&utm_medium=pc_web"
     >豆瓣书店</a>
    </li>
    <li    ><a href="https://book.douban.com/annual/2018?source=navigation"
            target="_blank"
     >2018年度榜单</a>
    </li>
    <li    ><a href="https://www.douban.com/standbyme/2018?source=navigation"
            target="_blank"
     >2018书影音报告</a>
    </li>
    <li          class=" book-cart"
    ><a href="https://market.douban.com/cart/?biz_type=book&utm_campaign=book_nav_cart&utm_source=douban&utm_medium=pc_web"
            target="_blank"
     >购物车</a>
    </li>
  </ul>
</div>

    <a href="https://book.douban.com/annual/2018?source=book_navigation" class="bookannual2018"></a>
  </div>
</div>

<script id="suggResult" type="text/x-jquery-tmpl">
  <li data-link="{{= url}}">
            <a href="{{= url}}" onclick="moreurl(this, {from:'book_search_sugg', query:'{{= keyword }}', subject_id:'{{= id}}', i: '{{= index}}', type: '{{= type}}'})">
            <img src="{{= pic}}" width="40" />
            <div>
                <em>{{= title}}</em>
                {{if year}}
                    <span>{{= year}}</span>
                {{/if}}
                <p>
                {{if type == "b"}}
                    {{= author_name}}
                {{else type == "a" }}
                    {{if en_name}}
                        {{= en_name}}
                    {{/if}}
                {{/if}}
                 </p>
            </div>
        </a>
        </li>
  </script>




    <script src="//img3.doubanio.com/dae/accounts/resources/8c80301/book/bundle.js" defer="defer"></script>





    <div id="wrapper">
        
    <div id="dale_book_subject_top_icon"></div>
<h1>
    <span property="v:itemreviewed">四个春天</span>
    <div class="clear"></div>
</h1>

        
  <div id="content">
    
    <div class="grid-16-8 clearfix">
      
      <div class="article">



<div class="indent">
  <div class="subjectwrap clearfix">
    



<div class="subject clearfix">
<div id="mainpic" class="">

  

  <a class="nbg"
      href="https://img3.doubanio.com/view/subject/l/public/s29957035.jpg" title="四个春天">
    <img src="https://img3.doubanio.com/view/subject/l/public/s29957035.jpg" title="点击看大图" alt="四个春天"
         rel="v:photo" style="width: 135px;max-height: 200px;">
  </a>



</div>





<div id="info" class="">



    
    
  
    <span>
      <span class="pl"> 作者</span>:
        
            
            <a class="" href="/search/%E9%99%86%E5%BA%86%E5%B1%B9">陆庆屹</a>
    </span><br/>

    
    
  
    <span class="pl">出版社:</span> 南海出版公司<br/>

    
    
  
    <span class="pl">出品方:</span>&nbsp;<a href="https://book.douban.com/series/39059?brand=1">新经典文化</a><br>

    
    
  

    
    
  

    
    
  

    
    
  
    <span class="pl">出版年:</span> 2019-1-1<br/>

    
    
  

    
    
  
    <span class="pl">定价:</span> 49.00元<br/>

    
    
  
    <span class="pl">装帧:</span> 平装<br/>

    
    
  

    
    
  
    
      
      <span class="pl">ISBN:</span> 9787544294881<br/>


</div>

</div>
























    





<div id="interest_sectl" class="">
  <div class="rating_wrap clearbox" rel="v:rating">
    <div class="rating_logo">豆瓣评分</div>
    <div class="rating_self clearfix" typeof="v:Rating">
      <strong class="ll rating_num " property="v:average"> 8.3 </strong>
      <span property="v:best" content="10.0"></span>
      <div class="rating_right ">
          <div class="ll bigstar40"></div>
            <div class="rating_sum">
                <span class="">
                    <a href="collections" class="rating_people"><span property="v:votes">311</span>人评价</a>
                </span>
            </div>


      </div>
    </div>
          
            
            
<span class="stars5 starstop" title="力荐">
    5星
</span>

            
<div class="power" style="width:51px"></div>

            <span class="rating_per">38.9%</span>
            <br>
            
            
<span class="stars4 starstop" title="推荐">
    4星
</span>

            
<div class="power" style="width:64px"></div>

            <span class="rating_per">48.6%</span>
            <br>
            
            
<span class="stars3 starstop" title="还行">
    3星
</span>

            
<div class="power" style="width:16px"></div>

            <span class="rating_per">12.2%</span>
            <br>
            
            
<span class="stars2 starstop" title="较差">
    2星
</span>

            
<div class="power" style="width:0px"></div>

            <span class="rating_per">0.3%</span>
            <br>
            
            
<span class="stars1 starstop" title="很差">
    1星
</span>

            
<div class="power" style="width:0px"></div>

            <span class="rating_per">0.0%</span>
            <br>
    </div>
</div>

  </div>
  





  
    
    <div id="interest_sect_level" class="clearfix">
        <a href="#" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-30389935-wish">
          <span>
            
<form method="POST" action="https://www.douban.com/register?reason=collectwish" class="miniform">
    <input type="submit" class="minisubmit j " value="想读" title="" />
</form>

          </span>
        </a>
        <a href="#" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-30389935-do">
          <span>
            
<form method="POST" action="https://www.douban.com/register?reason=collectdo" class="miniform">
    <input type="submit" class="minisubmit j " value="在读" title="" />
</form>

          </span>
        </a>
        <a href="#" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-30389935-collect">
          <span>
            
<form method="POST" action="https://www.douban.com/register?reason=collectcollect" class="miniform">
    <input type="submit" class="minisubmit j " value="读过" title="" />
</form>

          </span>
        </a>
      <div class="ll j a_stars">
        
    
    评价:
    <span id="rating"> <span id="stars" data-solid="https://img3.doubanio.com/f/shire/5a2327c04c0c231bced131ddf3f4467eb80c1c86/pics/rating_icons/star_onmouseover.png" data-hollow="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" data-solid-2x="https://img3.doubanio.com/f/shire/7258904022439076d57303c3b06ad195bf1dc41a/pics/rating_icons/star_onmouseover@2x.png" data-hollow-2x="https://img3.doubanio.com/f/shire/95cc2fa733221bb8edd28ad56a7145a5ad33383e/pics/rating_icons/star_hollow_hover@2x.png">

            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-30389935-1">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star1" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-30389935-2">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star2" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-30389935-3">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star3" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-30389935-4">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star4" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-30389935-5">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star5" width="16" height="16"/></a>
    </span><span id="rateword" class="pl"></span>
    <input id="n_rating" type="hidden" value=""  />
    </span>

      </div>
      

    </div>



  
  <div class="gtleft">
    <ul class="ul_subject_menu bicelink color_gray pt6 clearfix">
        <li>
          <img src="https://img3.doubanio.com/f/shire/5bbf02b7b5ec12b23e214a580b6f9e481108488c/pics/add-review.gif" />&nbsp;<a href="https://www.douban.com/register?reason=annotate" class="j a_show_login" rel="nofollow">写笔记</a>
        </li>

          <li>
            <img src="https://img3.doubanio.com/f/shire/5bbf02b7b5ec12b23e214a580b6f9e481108488c/pics/add-review.gif" />&nbsp;<a class="j a_show_login" href="https://www.douban.com/register?reason=review" rel="nofollow">写书评</a>
          </li>

      <li>

  <span class="rr">
  

    <img src="https://img3.doubanio.com/pics/add-cart.gif"/>
      <a class="j a_show_login" href="http://http://www.douban.com/register?reason=addbook2cart" rel="nofollow">加入购书单</a>
  <span class="hidden">已在<a href="https://book.douban.com/cart">购书单</a></span>
</span><br class="clearfix" />
</li>


        
        
    
    <li class="rec" id="C-30389935">
        <a href="#" data-url="https://book.douban.com/subject/30389935/" data-desc="" data-title="书籍《四个春天》 (来自豆瓣) " data-pic="https://img3.doubanio.com/view/subject/l/public/s29957035.jpg" class="bn-sharing ">分享到</a> &nbsp;&nbsp;
    </li>
    <script>
    var __cache_url = __cache_url || {};
    (function(u){
        if(__cache_url[u]) return;
        __cache_url[u] = true;
        window.DoubanShareIcons = 'https://img3.doubanio.com/f/shire/d15ffd71f3f10a7210448fec5a68eaec66e7f7d0/pics/ic_shares.png';
        var initShareButton = function() {
          $.ajax({url:u,dataType:'script',cache:true});
        };
        if (typeof Do == 'function' && 'ready' in Do) {
          Do('https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
            'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
            initShareButton);
        } else if(typeof Douban == 'object' && 'loader' in Douban) {
          Douban.loader.batch(
            'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
            'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js'
          ).done(initShareButton);
        }
    })('https://img3.doubanio.com/f/shire/6e6a5f21daeec19bbb41bf48c07fccaa4dad4d98/js/lib/sharebutton.js');
    </script>

    </ul>
  </div>


    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>


<script>
  //bind events for collection button.
  $('.collect_btn', '#interest_sect_level').each(function(){
      Douban.init_collect_btn(this);
  });
</script>








</div>

<br clear="all">
<div id="collect_form_30389935"></div>
<div class="related_info">
  






  

  <h2>
    <span class="">内容简介</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>



<div class="indent" id="link-report">
    
      <div class="">
        <style type="text/css" media="screen">
.intro p{text-indent:2em;word-break:normal;}
</style>
<div class="intro">
    <p>《四个春天》是陆庆屹首部文字作品，他用深情质朴的文字和饱含温度的摄影，记录下父母、故乡、旧识……在书页间搭建起西南小城中充满烟火气、人情味，同时充盈着诗意的生活景象。</p>    <p>在外的人，只能在春节时回家，和父母共处的日子，大都在春天里——相濡以沫半个世纪的父母，“土味”却饱含智慧的乡土人，细碎平常的片段，柔软浪漫的小事。清水白菜式的记录，简单却有热力。虽不是自己的故事却发生在我们生命中的每一天里。</p>    <p>风物景致或许不尽相同，但是对亲情的温润感悟，对世事变迁的杂陈体味，对时光与故人的怀恋，是相通的。</p>    <p>写给天下的父母和每一个游子。诗意，不在远方，在身旁。</p></div>

      </div>
    
<link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/c4c6dd266f58b41cbeebc9e4e6d7dd7b2a5c3711/css/report.css" />
<link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css" />
<link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/b45aa277f8b8df40596b96582dafb1ed0a899a64/css/report_dialog.css" />
<script type="text/javascript" src="https://img3.doubanio.com/f/shire/a14501790b4a2db257dc5be5e37d820e600703c6/js/report_dialog.js"></script>
<style>
    #link-report .report { text-align: right; font-size: 12px; visibility: hidden; }
    #link-report .report a { color: #BBB; }
    #link-report .report a:hover { color: #FFF; background-color: #BBB; }
</style>
<script>
    Do = (typeof Do === 'undefined')? $ : Do;
    Do(function(){
    $("body").delegate("#link-report", 'mouseenter mouseleave', function(e){
      switch (e.type) {
        case "mouseenter":
          $(this).find(".report").css('visibility', 'visible');
          break;
        case "mouseleave":
          $(this).find(".report").css('visibility', 'hidden');
          break;
      }
    });
    $("#link-report").delegate(".report a", 'click', function(e){
        e.preventDefault();
        var opt = "";
        var obj = $(e.target).closest('#link-report');
        var id = obj.length != 0 ? obj.data("id") : undefined;
        var params = (opt&&id) ? ''.concat(opt, '=', id) : '';

        var url = "https://book.douban.com/subject/30389935/";
        url += (~url.indexOf('?') ? '&' : '?') + params
        url = url.replace(/\&+/g, '&')
        generate_report_dialog({report_url: url, type: 'subject'});
    });

    $("#link-report").append('<div class="report"><a rel="nofollow" href="#">举报</a></div>');
  });
</script>

</div>

  

























    
  

  <h2>
    <span class="">作者简介</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>



      <div class="indent ">
          
            <div class="">
            <style type="text/css" media="screen">
.intro p{text-indent:2em;word-break:normal;}
</style>
<div class="intro">
    <p>陆庆屹，1973年生于贵州独山。15岁离家，曾做过足球运动员、酒吧歌手、矿工、摄影师，现为独立电影制作人。</p>    <p>电影拍摄零基础的他，耗时6年完成了导演处女作《四个春天》，记录下家乡年迈父母寻常生活中的诗意。无论影像还是文字，他观察日常，却能剥离日常中的庸碌琐碎，为平凡的人与 事赋予温度与质感。</p></div>

            </div>
      </div>











































  

  
    




  

  <h2>
    <span class="">目录</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>



<div class="indent" id="dir_30389935_short">
        代序 关于《四个春天》的一些小事<br/>
        光阴寂寂<br/>
        老帅<br/>
        吴叔<br/>
        老哥俩儿<br/>
        祖方舅<br/>
    · · · · · ·
    (<a href="javascript:$('#dir_30389935_short').hide();$('#dir_30389935_full').show();$.get('/j/subject/j_dir_count',{id:30389935});void(0);">更多</a>)
</div>

<div class="indent" id="dir_30389935_full" style="display:none">
        代序 关于《四个春天》的一些小事<br/>
        光阴寂寂<br/>
        老帅<br/>
        吴叔<br/>
        老哥俩儿<br/>
        祖方舅<br/>
        老三<br/>
        重返少年<br/>
        老家人<br/>
        赶场天<br/>
        麻尾记忆<br/>
        糖与蜜<br/>
        “童工”时代<br/>
        我爸<br/>
        我妈<br/>
        爸的书房<br/>
        后园<br/>
        打野菜<br/>
        送别<br/>
        速写<br/>
        城南一夜<br/>
        新居<br/>
        意外的清晨<br/>
        山居几日<br/>
        想做就去做<br/>
     · · · · · ·     (<a href="javascript:$('#dir_30389935_full').hide();$('#dir_30389935_short').show();void(0);">收起</a>)
</div>

    





  

  <h2>
    <span class="">&#34;四个春天&#34;试读</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


<div class="indent">
<p>我妈天生暴脾气，见不得不平事，眼睛一瞪，路灯都要黯淡几分；又争强好胜不服输，眉头下从没写过“困难”二字。外公生前逢人就说：这丫头投错胎了，要是个男娃就刚好！
我家在贵州南部的一个小县城。十年前，姐到沈阳工作，那时家里穷，坐火车属于巨额花费，爸妈想去看看女儿很不容易，一般春节才能团聚。后来，我姐在公司当了领导，收入涨了，想让爸妈直接从贵阳坐飞机到沈阳，爸晕车很...</p>
<ul class="col2-list clearfix">

<li>
    <a href="https://book.douban.com/reading/67713784/">我 妈 </a>
</li>
</ul>
<div align="right">  · · · · · ·    (<a href="https://book.douban.com/subject/30389935/reading/">查看全部试读</a>)</div>
</div>



  


<link rel="stylesheet" href="https://img3.doubanio.com/f/verify/16c7e943aee3b1dc6d65f600fcc0f6d62db7dfb4/entry_creator/dist/author_subject/style.css">

<div id="author_subject" class="author-wrapper">
  <div class="loading"></div>
</div>

<script type="text/javascript">
  var answerObj = {
    TYPE: 'book',
    SUBJECT_ID: '30389935',
    ISALL: 'False' || false,
    USER_ID: 'None'
  }
</script>
<script src="https://img3.doubanio.com/f/book/61252f2f9b35f08b37f69d17dfe48310dd295347/js/book/lib/react/bundle.js"></script>
<script type="text/javascript" src="https://img3.doubanio.com/f/verify/ac140ef86262b845d2be7b859e352d8196f3f6d4/entry_creator/dist/author_subject/index.js"></script> 
  






<div id="db-tags-section" class="blank20">
  
  

  <h2>
    <span class="">豆瓣成员常用的标签(共61个)</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


  <div class="indent">    <span class="">
        <a class="  tag" href="/tag/四个春天">四个春天</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/随笔">随笔</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/陆庆屹">陆庆屹</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/散文随笔">散文随笔</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/文学">文学</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/2019">2019</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/杂文">杂文</a> &nbsp;    </span>
    <span class="">
        <a class="  tag" href="/tag/电影人">电影人</a> &nbsp;    </span>
  </div>
</div>


  


  












<div id="rec-ebook-section" class="block5 subject_show">
  

  
  

  <h2>
    <span class="">喜欢读&#34;四个春天&#34;的人也喜欢的电子书</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


  <div class="tips-mod">
    支持 Web、iPhone、iPad、Android 阅读器
  </div>
  <div class="content clearfix">
      
      <dl>
        <dt>
          <a href="https://read.douban.com/ebook/17402307/?dcs=subject-rec&amp;dcm=douban&amp;dct=30389935" target="_blank">
            <span class="cover-outer">
              <img src="https://img1.doubanio.com/view/ark_article_cover/retina/public/17402307.jpg?v=0">
            </span>
          </a>
        </dt>
        <dd>
          <div class="title">
              <a href="https://read.douban.com/ebook/17402307/" target="_blank">山水与日常</a>
          </div>
          <div class="price">
              1.99元
          </div>
        </dd>
      </dl>
      
      <dl>
        <dt>
          <a href="https://read.douban.com/ebook/2055412/?dcs=subject-rec&amp;dcm=douban&amp;dct=30389935" target="_blank">
            <span class="cover-outer">
              <img src="https://img3.doubanio.com/view/ark_article_cover/retina/public/2055412.jpg?v=0">
            </span>
          </a>
        </dt>
        <dd>
          <div class="title">
              <a href="https://read.douban.com/ebook/2055412/" target="_blank">京都第五年</a>
          </div>
          <div class="price">
              1.99元
          </div>
        </dd>
      </dl>
      
      <dl>
        <dt>
          <a href="https://read.douban.com/ebook/29591419/?dcs=subject-rec&amp;dcm=douban&amp;dct=30389935" target="_blank">
            <span class="cover-outer">
              <img src="https://img1.doubanio.com/view/ark_article_cover/retina/public/29591419.jpg?v=0">
            </span>
          </a>
        </dt>
        <dd>
          <div class="title">
              <a href="https://read.douban.com/ebook/29591419/" target="_blank">大裂</a>
          </div>
          <div class="price">
              5.25元
          </div>
        </dd>
      </dl>
      
      <dl>
        <dt>
          <a href="https://read.douban.com/ebook/17036628/?dcs=subject-rec&amp;dcm=douban&amp;dct=30389935" target="_blank">
            <span class="cover-outer">
              <img src="https://img1.doubanio.com/view/ark_article_cover/retina/public/17036628.jpg?v=0">
            </span>
          </a>
        </dt>
        <dd>
          <div class="title">
              <a href="https://read.douban.com/ebook/17036628/" target="_blank">人间采蜜记：李银河自传</a>
          </div>
          <div class="price">
              1.99元
          </div>
        </dd>
      </dl>
      
      <dl>
        <dt>
          <a href="https://read.douban.com/ebook/1073721/?dcs=subject-rec&amp;dcm=douban&amp;dct=30389935" target="_blank">
            <span class="cover-outer">
              <img src="https://img3.doubanio.com/view/ark_article_cover/retina/public/1073721.jpg?v=0">
            </span>
          </a>
        </dt>
        <dd>
          <div class="title">
              <a href="https://read.douban.com/ebook/1073721/" target="_blank">我讲个笑话，你可别哭啊</a>
          </div>
          <div class="price">
              0.99元
          </div>
        </dd>
      </dl>
  </div>
</div>

<div id="db-rec-section" class="block5 subject_show knnlike">
  
  
  

  <h2>
    <span class="">喜欢读&#34;四个春天&#34;的人也喜欢</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


  <div class="content clearfix">
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/30259150/" onclick="moreurl(this, {'total': 10, 'clicked': '30259150', 'pos': 0, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s29857786.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/30259150/" onclick="moreurl(this, {'total': 10, 'clicked': '30259150', 'pos': 0, 'identifier': 'book-rec-books'})" class="">
            失明的摄影师
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/30327470/" onclick="moreurl(this, {'total': 10, 'clicked': '30327470', 'pos': 1, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s29870997.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/30327470/" onclick="moreurl(this, {'total': 10, 'clicked': '30327470', 'pos': 1, 'identifier': 'book-rec-books'})" class="">
            衣的现象学
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/27183696/" onclick="moreurl(this, {'total': 10, 'clicked': '27183696', 'pos': 2, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s29611249.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/27183696/" onclick="moreurl(this, {'total': 10, 'clicked': '27183696', 'pos': 2, 'identifier': 'book-rec-books'})" class="">
            小鸟睡在我身旁
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/30379604/" onclick="moreurl(this, {'total': 10, 'clicked': '30379604', 'pos': 3, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s29924003.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/30379604/" onclick="moreurl(this, {'total': 10, 'clicked': '30379604', 'pos': 3, 'identifier': 'book-rec-books'})" class="">
            故乡的味道
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/30234603/" onclick="moreurl(this, {'total': 10, 'clicked': '30234603', 'pos': 4, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s29834751.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/30234603/" onclick="moreurl(this, {'total': 10, 'clicked': '30234603', 'pos': 4, 'identifier': 'book-rec-books'})" class="">
            任天堂哲学
          </a>
        </dd>
      </dl>
        <dl class="clear"></dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/30344780/" onclick="moreurl(this, {'total': 10, 'clicked': '30344780', 'pos': 5, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s29888237.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/30344780/" onclick="moreurl(this, {'total': 10, 'clicked': '30344780', 'pos': 5, 'identifier': 'book-rec-books'})" class="">
            观看王维的十九种方式
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/27616952/" onclick="moreurl(this, {'total': 10, 'clicked': '27616952', 'pos': 6, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s29685208.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/27616952/" onclick="moreurl(this, {'total': 10, 'clicked': '27616952', 'pos': 6, 'identifier': 'book-rec-books'})" class="">
            印度放浪
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/27201290/" onclick="moreurl(this, {'total': 10, 'clicked': '27201290', 'pos': 7, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s29669694.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/27201290/" onclick="moreurl(this, {'total': 10, 'clicked': '27201290', 'pos': 7, 'identifier': 'book-rec-books'})" class="">
            啊！这样就能辞职了
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/27075586/" onclick="moreurl(this, {'total': 10, 'clicked': '27075586', 'pos': 8, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img3.doubanio.com/view/subject/l/public/s29539862.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/27075586/" onclick="moreurl(this, {'total': 10, 'clicked': '27075586', 'pos': 8, 'identifier': 'book-rec-books'})" class="">
            东京漂流
          </a>
        </dd>
      </dl>
      
      <dl class="">
        <dt>
            <a href="https://book.douban.com/subject/30324264/" onclick="moreurl(this, {'total': 10, 'clicked': '30324264', 'pos': 9, 'identifier': 'book-rec-books'})"><img class="m_sub_img" src="https://img1.doubanio.com/view/subject/l/public/s29869087.jpg"/></a>
        </dt>
        <dd>
          <a href="https://book.douban.com/subject/30324264/" onclick="moreurl(this, {'total': 10, 'clicked': '30324264', 'pos': 9, 'identifier': 'book-rec-books'})" class="">
            天边一星子
          </a>
        </dd>
      </dl>
        <dl class="clear"></dl>
  </div>
</div>

  






    <link rel="stylesheet" href="https://img3.doubanio.com/f/book/3ec79645ad5a5d15c9ead3c58da97f5d662c7400/css/book/subject/comment.css"/>
    <div class="mod-hd">
        

        <a class="redbutt j a_show_login rr" href="https://www.douban.com/register?reason=review" rel="nofollow">
            <span> 我来说两句 </span>
        </a>

            
  

  <h2>
    <span class="">短评</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;
      <span class="pl">&nbsp;(
          <a href="https://book.douban.com/subject/30389935/comments/">全部 185 条</a>
        ) </span>

  </h2>


    </div>
    <div class="nav-tab">
        
    <div class="tabs-wrapper  line">
        <a class="short-comment-tabs on-tab" href="hot" data-tab="hot">热门</a>
        <span>/</span>
        <a class="short-comment-tabs " href="new" data-tab="new">最新</a>
        <span>/</span>
        <a class="j a_show_login " href="follows" data-tab="follows">好友</a>
    </div>

    </div>
    <div id="comment-list-wrapper" class="indent">
        

<div id="comments" class="comment-list hot show">
        <ul>
                
    <li class="comment-item" data-cid="1587199940">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1587199940" class="vote-count">5</span>
                        <a href="javascript:;" id="btn-1587199940" class="j a_show_login" data-cid="1587199940">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/149797382/">風_梢_螢_星_</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2019-01-05</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">按照《论摄影》中的论述，摄影并非“短暂、信息极少、无实体、虚弱的”真实的影子，任何一位真正的摄影师都不会忽略摄影背后的道德意义和美学使命。故而摄影极其磨练摄影师的观察力、美学意识和艺术知觉，在此领域颇有成就的导演，以文字为媒介通常也能缔造一些独到的画面感，黑泽明、伯格曼和塔氏如是。陆庆屹这本集子里的文章亦如是，喜欢他对光影的留恋和描述，笔触很像纳博科夫呀。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1610940332">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1610940332" class="vote-count">3</span>
                        <a href="javascript:;" id="btn-1610940332" class="j a_show_login" data-cid="1610940332">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/73857325/">青苔入心</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2019-01-07</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">一本写给生活的抒情诗，温柔细腻又动人。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1621257745">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1621257745" class="vote-count">3</span>
                        <a href="javascript:;" id="btn-1621257745" class="j a_show_login" data-cid="1621257745">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/159913149/">阜咚</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2019-01-10</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">想起了以前看过的《马桥词典》等乡土文学，非常真诚朴实，多少有些&#34;书生意气&#34;的感觉，很有意思。书里偷糖吃、进山挖菜等情节读了之后很有共鸣，我小时候也蹑手蹑脚地打开冰箱门偷吃东西然后快速“销毁罪证”，甚至把食品袋也小心摆好，害怕爸妈发现。作者通过双眼和镜头观察人们，而人们也通过这本书观察作者和自己，形成一种微妙的“看与被看”，很美。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1621429692">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1621429692" class="vote-count">16</span>
                        <a href="javascript:;" id="btn-1621429692" class="j a_show_login" data-cid="1621429692">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/122394728/">Vise Versa</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2019-01-10</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">导演对电影的斧凿很浅，几近克制。剩下的丰沛的情感都盛在这本书里。
&#34;若只是廉价的自我感动，在这茫茫人海中毫无意义，在这缄默的天地间更没有任何价值。&#34;这样的人怎么能不做出好的电影呢?</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1618184407">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1618184407" class="vote-count">24</span>
                        <a href="javascript:;" id="btn-1618184407" class="j a_show_login" data-cid="1618184407">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/3965296/">Dan</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2019-01-08</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">前几天看书枝在安利这部电影，正临期末没有时间去电影院看，就找来书读。
一家人的日常读起来就像是在看是枝裕和的电影，温情柔软，却笑中带泪。有一段印象最为深刻：妈妈带着大包小包的特产坐飞机去沈阳看姐姐，姐姐接机时看到妈妈头发被汗水打湿，东一片西一缕地贴在脸上，但是根本顾不上，只是四处张望着找她，扑过去帮忙时有两件行李重得提不动，打开一看竟是两大袋糯米粑，姐姐再也控制不住泪水，跌坐在机场哭起来……看到这...</span>
                <span class="hide-item full">前几天看书枝在安利这部电影，正临期末没有时间去电影院看，就找来书读。
一家人的日常读起来就像是在看是枝裕和的电影，温情柔软，却笑中带泪。有一段印象最为深刻：妈妈带着大包小包的特产坐飞机去沈阳看姐姐，姐姐接机时看到妈妈头发被汗水打湿，东一片西一缕地贴在脸上，但是根本顾不上，只是四处张望着找她，扑过去帮忙时有两件行李重得提不动，打开一看竟是两大袋糯米粑，姐姐再也控制不住泪水，跌坐在机场哭起来……看到这一段，想到之前看的影评中的剧透，我也哭得不能自已。
作者写得都是日常中再普通不过的小事，可是当这些小事以文字的形式展现在我面前，就很能戳中我的心。想起李银河在访谈中的一段话，她说，个人的生活对宇宙没有意义，人类整体的生活对宇宙也没有意义，唯有诗意地栖居，对自己有意义。</span>
                <span class="expand">(<a href="javascript:;">展开</a>)</span>
            </p>
        </div>
    </li>


        </ul>
</div>

        

<div id="comments" class="comment-list new noshow">
        <ul>
                
    <li class="comment-item" data-cid="1635116438">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1635116438" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1635116438" class="j a_show_login" data-cid="1635116438">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/ytxwz1978/">樱桃小丸子</a>
                        <span class="user-stars allstar30 rating" title="还行"></span>
                    <span>2019-01-20</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">普通人的生活片段感悟，绘景的文字清美，记人的语言质朴，但是整体结构琐碎，个人不太喜欢像手摇摄影机一样晃来晃去的架构，看电影会眼晕，读书会脑乱。

</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1635111939">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1635111939" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1635111939" class="j a_show_login" data-cid="1635111939">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/68825809/">海哈哈</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2019-01-20</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">温柔能带来这个世界上最美好的东西。
人是在什么时候开始变得温柔呢？
逐渐开始觉得，也许是在与「另一个人」的想法发生碰撞，倔强得不服输之后，因为更爱对方而选择妥协的那一刻。
先认输，然后再试着去交流，依旧无果也不重要，但对方一定会明白你的妥协，以及没说出口的温柔。
类似于此的感受，或许只有在失去之后才明白，所以从这个角度来说，缘分就是最好的温柔，留给了最合适的人。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1635034330">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1635034330" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1635034330" class="j a_show_login" data-cid="1635034330">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/92789223/">你都如何回忆我</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2019-01-20</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">内心的柔软与温柔足以抵御世间一切磨难</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1634784542">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1634784542" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1634784542" class="j a_show_login" data-cid="1634784542">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/52383485/">由酱就是要做梦</a>
                        <span class="user-stars allstar40 rating" title="推荐"></span>
                    <span>2019-01-20</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">3.5 写父母的部分真是让人动容 私心希望最后一篇写电影制作过程的能更长一些。</span>
            </p>
        </div>
    </li>

                
    <li class="comment-item" data-cid="1634661152">
        <div class="comment">
            <h3>
                <span class="comment-vote">
                    <span id="c-1634661152" class="vote-count">0</span>
                        <a href="javascript:;" id="btn-1634661152" class="j a_show_login" data-cid="1634661152">有用</a>
                </span>
                <span class="comment-info">
                    <a href="https://www.douban.com/people/171044360/">Hiker</a>
                        <span class="user-stars allstar50 rating" title="力荐"></span>
                    <span>2019-01-19</span>
                </span>
            </h3>
            <p class="comment-content">
            
                <span class="short">前几天看完的
陆庆屹的文字和他的镜头一样温柔啊</span>
            </p>
        </div>
    </li>


        </ul>
</div>

    </div>
        <p>&gt; <a href="https://book.douban.com/subject/30389935/comments/">更多短评 185 条</a></p>
    <script src="https://img3.doubanio.com/f/book/87c0b94bb0c3698250d8df98cebf5bb30e7f08fe/js/book/subject/short_comment_vote.js"></script>
    <script src="https://img3.doubanio.com/f/book/39eace58cab8aaeec45a44e878bf0ed06f2ed0a4/js/book/subject/short_comment_nav.js"></script>
    <script>
        (function(){
            $('.comment-list').delegate('.vote-comment', 'click', function(e) {
                vote_comment(e);
            }).delegate('.delete-comment', 'click', function(e) {
                if (confirm('确定删除吗？')) {
                    delete_comment(e);
                }
            });
        })();
    </script>

  

<link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/73ed658484f98d44.css">

<section class="topics mod">
    <header>
        <h2>
            四个春天的话题 · · · · · ·
            <span class="pl">( <span class="gallery_topics">全部 <span id="topic-count"></span> 条</span> )</span>
        </h2>
    </header>

    




<section class="subject-topics">
    <div class="topic-guide" id="topic-guide">
        <img class="ic_question" src="//img3.doubanio.com/f/ithildin/b1a3edea3d04805f899e9d77c0bfc0d158df10d5/pics/export/icon_question.png">
        <div class="tip_content">
            <div class="tip_title">什么是话题</div>
            <div class="tip_desc">
                <div>无论是一部作品、一个人，还是一件事，都往往可以衍生出许多不同的话题。将这些话题细分出来，分别进行讨论，会有更多收获。</div>
            </div>
        </div>
        <img class="ic_guide" src="//img3.doubanio.com/f/ithildin/529f46d86bc08f55cd0b1843d0492242ebbd22de/pics/export/icon_guide_arrow.png">
        <img class="ic_close" id="topic-guide-close" src="//img3.doubanio.com/f/ithildin/2eb4ad488cb0854644b23f20b6fa312404429589/pics/export/close@3x.png">
    </div>

    <div id="topic-items"></div>

    <script>
        window.subject_id = 30389935;
        window.join_label_text = '写书评参与';

        window.topic_display_count = 4;
        window.topic_item_display_count = 1;
        window.no_content_fun_call_name = "no_topic";

        window.guideNode = document.getElementById('topic-guide');
        window.guideNodeClose = document.getElementById('topic-guide-close');
    </script>
    
        <link rel="stylesheet" href="https://img3.doubanio.com/f/ithildin/f731c9ea474da58c516290b3a6b1dd1237c07c5e/css/export/subject_topics.css">
        <script src="https://img3.doubanio.com/f/ithildin/d3590fc6ac47b33c804037a1aa7eec49075428c8/js/export/moment-with-locales-only-zh.js"></script>
        <script src="https://img3.doubanio.com/f/ithildin/c600fdbe69e3ffa5a3919c81ae8c8b4140e99a3e/js/export/subject_topics.js"></script>

</section>

    <script>
        function no_topic(){
            $('#content .topics').remove();
        }
    </script>
</section>

<section class="reviews mod book-content">
    <header>
        <a href="new_review" rel="nofollow" class="create-review redbutt rr"
            data-isverify="False"
            data-verify-url="https://www.douban.com/accounts/phone/verify?redir=https://book.douban.com/subject/30389935/new_review">
            <span>我要写书评</span>
        </a>
        <h2>
            四个春天的书评 · · · · · ·
            <span class="pl">( <a href="reviews">全部 11 条</a> )</span>
        </h2>
    </header>

    

<style>
#gallery-topics-selection {
  position: fixed;
  width: 595px;
  padding: 40px 40px 33px 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 16px 0 rgba(0, 0, 0, 0.2);
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  z-index: 9999;
}
#gallery-topics-selection h1 {
  font-size: 18px;
  color: #007722;
  margin-bottom: 36px;
  padding: 0;
  line-height: 28px;
  font-weight: normal;
}
#gallery-topics-selection .gl_topics {
  border-bottom: 1px solid #dfdfdf;
  max-height: 298px;
  overflow-y: scroll;
}
#gallery-topics-selection .topic {
  margin-bottom: 24px;
}
#gallery-topics-selection .topic_name {
  font-size: 15px;
  color: #333;
  margin: 0;
  line-height: inherit;
}
#gallery-topics-selection .topic_meta {
  font-size: 13px;
  color: #999;
}
#gallery-topics-selection .topics_skip {
  display: block;
  cursor: pointer;
  font-size: 16px;
  color: #3377AA;
  text-align: center;
  margin-top: 33px;
}
#gallery-topics-selection .topics_skip:hover {
  background: transparent;
}
#gallery-topics-selection .close_selection {
  position: absolute;
  width: 30px;
  height: 20px;
  top: 46px;
  right: 40px;
  background: #fff;
  color: #999;
  text-align: right;
}
#gallery-topics-selection .close_selection:hover{
  background: #fff;
  color: #999;
}
</style>




        <div class="review_filter">
            <a href="javascript:;;" class="cur" data-sort="">热门</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="time">最新</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="follow">好友</a href="javascript:;;">
            
        </div>


        



<div class="review-list  ">
        
    

        
    
    <div data-cid="9875705">
        <div class="main review-item" id="9875705">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/zhaoxun69/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u1131801-3.jpg">
        </a>

        <a href="https://www.douban.com/people/zhaoxun69/" class="name">盘子</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-06" class="main-meta">2019-01-06 04:39:30</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9875705/">敏锐·记忆·一种天真和坦诚</a></h2>

                <div id="review_9875705_short" class="review-short" data-rid="9875705">
                    <div class="short-content">

                        作者出这本书之前犹豫了很久，觉得自己的文字远不够印刷。我作为作者电影处女作的工作人员，自己也有犹豫，怕出版一本同名图书会被读解成功利行为或仅仅变成影像的营销附属品——这对作者对文字都是不公平的。但又觉得不出有点可惜，锦上添花，有什么不好呢，后来也在劝他可以...

                        &nbsp;(<a href="javascript:;" id="toggle-9875705-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9875705_full" class="hidden">
                    <div id="review_9875705_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9875705" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9875705">
                                35
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9875705" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9875705">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9875705/#comments" class="reply">10回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9901890">
        <div class="main review-item" id="9901890">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/xiaoxiaodeyuanz/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u14597285-44.jpg">
        </a>

        <a href="https://www.douban.com/people/xiaoxiaodeyuanz/" class="name">远子</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-17" class="main-meta">2019-01-17 20:21:13</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9901890/">今天终于晴了</a></h2>

                <div id="review_9901890_short" class="review-short" data-rid="9901890">
                    <div class="short-content">

                        这本书与纪录片的内容重合的地方不多。总体上，前面写人的几篇要好过后面写景和写情的。看的时候有一些想法，本想写进短评，没想到放不下，只好搁这儿。和书的关系不大。 我和陆庆屹只见过三面，但感觉已经很熟悉了，相信很多和他接触过的人都会有这种感觉，大家都亲切地称他为...

                        &nbsp;(<a href="javascript:;" id="toggle-9901890-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9901890_full" class="hidden">
                    <div id="review_9901890_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9901890" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9901890">
                                25
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9901890" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9901890">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9901890/#comments" class="reply">3回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9872605">
        <div class="main review-item" id="9872605">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/jiayinxy/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u78857520-6.jpg">
        </a>

        <a href="https://www.douban.com/people/jiayinxy/" class="name">佳音</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-04" class="main-meta">2019-01-04 23:03:46</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9872605/">妈妈，你要好好的</a></h2>

                <div id="review_9872605_short" class="review-short" data-rid="9872605">
                    <div class="short-content">

                        书中没有特意的煽情，也没有宣传电影的商业化。 有的是身边人点点滴滴的琐事，是有迹可循的回忆。 一口气看下来，脑中想的总是《项脊轩志》，还有巴斯大学的杨娃娃写的那句：“项脊轩志的迷人之处大概就在于它的平淡。没有胸怀大志，没有家国天下，有的是琐碎日子里平平淡淡的...

                        &nbsp;(<a href="javascript:;" id="toggle-9872605-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9872605_full" class="hidden">
                    <div id="review_9872605_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9872605" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9872605">
                                12
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9872605" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9872605">
                                1
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9872605/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9885639">
        <div class="main review-item" id="9885639">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/157800227/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u157800227-2.jpg">
        </a>

        <a href="https://www.douban.com/people/157800227/" class="name">青衫</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-10" class="main-meta">2019-01-10 03:49:43</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9885639/">推书</a></h2>

                <div id="review_9885639_short" class="review-short" data-rid="9885639">
                    <div class="short-content">

                        也是一个想看电影而不得的人，知道出书了，然而人在国外，想看书也不得。谁曾想，居然一下看到微信读书上在推介！立时对微信读书的评价提高了不止一个档。 初读，觉得有汪曾祺的影子，那浓绿山水洞坑又有沈从文的湘西山水的意思。与电影或多或少的联系又有吴念真的感觉。然而一...

                        &nbsp;(<a href="javascript:;" id="toggle-9885639-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9885639_full" class="hidden">
                    <div id="review_9885639_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9885639" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9885639">
                                2
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9885639" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9885639">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9885639/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9897184">
        <div class="main review-item" id="9897184">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/164525950/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u164525950-30.jpg">
        </a>

        <a href="https://www.douban.com/people/164525950/" class="name">乌</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2019-01-15" class="main-meta">2019-01-15 12:46:22</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9897184/">小记</a></h2>

                <div id="review_9897184_short" class="review-short" data-rid="9897184">
                    <div class="short-content">

                        在车上读这本书的时候，书里对于家乡，对于童年的描写，共鸣太多。但是作者对于家乡的热爱，我远不及。字里行间的种种，让我感到作者应该是一个内心柔软的人，并且饱含情绪。 看山看水，路过遇见，总是勾出他心里的感慨，并且偏向于伤怀。文风很稳，情感线基本保持在一个淡淡叙...

                        &nbsp;(<a href="javascript:;" id="toggle-9897184-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9897184_full" class="hidden">
                    <div id="review_9897184_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9897184" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9897184">
                                1
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9897184" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9897184">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9897184/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9884210">
        <div class="main review-item" id="9884210">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/187381538/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u187381538-1.jpg">
        </a>

        <a href="https://www.douban.com/people/187381538/" class="name">晴天</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2019-01-09" class="main-meta">2019-01-09 15:17:24</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9884210/">春天春天春天春天春天春天</a></h2>

                <div id="review_9884210_short" class="review-short" data-rid="9884210">
                    <div class="short-content">

                        春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天春天...

                        &nbsp;(<a href="javascript:;" id="toggle-9884210-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9884210_full" class="hidden">
                    <div id="review_9884210_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9884210" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9884210">
                                1
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9884210" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9884210">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9884210/#comments" class="reply">8回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9905081">
        <div class="main review-item" id="9905081">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/52594056/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u52594056-5.jpg">
        </a>

        <a href="https://www.douban.com/people/52594056/" class="name">暖暖_jm</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-19" class="main-meta">2019-01-19 12:42:32</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9905081/">家是永远的安心之处</a></h2>

                <div id="review_9905081_short" class="review-short" data-rid="9905081">
                    <div class="short-content">

                        读完。 这次体验了微信读书的听书模式，断断续续听完，好处是可以闭着眼睛，边听边放松。不好说是听着容易走神，稍微不注意就漏掉大段内容而不自知，想来还是看书更记忆深刻。 看这本书是源于在豆瓣看到同名电影最近热映，好评如潮，于是关注了这个人，看了他的相册，又读了他...

                        &nbsp;(<a href="javascript:;" id="toggle-9905081-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9905081_full" class="hidden">
                    <div id="review_9905081_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9905081" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9905081">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9905081" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9905081">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9905081/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9903004">
        <div class="main review-item" id="9903004">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/175736025/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u175736025-1.jpg">
        </a>

        <a href="https://www.douban.com/people/175736025/" class="name">梨涡</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2019-01-18" class="main-meta">2019-01-18 11:24:46</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9903004/">四个春天</a></h2>

                <div id="review_9903004_short" class="review-short" data-rid="9903004">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇书评可能有关键情节透露</p>

                        四个春天 看完纪录片之后，我和朋友两个人都很沉默，感慨太多却不知从何说起。偶然间看到还有同名的书籍，便带着那时的感动开始阅读。 陆的成长和我最初设想的并不一样，单看电影，他的父母很有趣味，很勤劳能干，很乐观豁达，有自己的兴趣爱好，有文化，和骨子里的文艺气息，...

                        &nbsp;(<a href="javascript:;" id="toggle-9903004-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9903004_full" class="hidden">
                    <div id="review_9903004_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9903004" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9903004">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9903004" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9903004">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9903004/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9900530">
        <div class="main review-item" id="9900530">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/69117350/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u69117350-3.jpg">
        </a>

        <a href="https://www.douban.com/people/69117350/" class="name">造梦坊</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-17" class="main-meta">2019-01-17 00:09:32</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9900530/">摘抄</a></h2>

                <div id="review_9900530_short" class="review-short" data-rid="9900530">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇书评可能有关键情节透露</p>

                        回到家已是凌晨一点。躺下后，回想这个夜里，那片巨大的林子里，除了黑，还有什么呢？我实在想知道它的模样，于是提前半个小时起床----天亮后，我要去林子那头赶车，要把昨夜的路清清楚楚地走一遍。 忙到夜深，大家围坐在厨房炉边闲聊，谁都不忍心开口说出那句“去睡吧”。若看...

                        &nbsp;(<a href="javascript:;" id="toggle-9900530-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9900530_full" class="hidden">
                    <div id="review_9900530_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9900530" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9900530">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9900530" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9900530">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9900530/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="9895720">
        <div class="main review-item" id="9895720">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/189791192/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/user_normal.jpg">
        </a>

        <a href="https://www.douban.com/people/189791192/" class="name">观山海</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2019-01-14" class="main-meta">2019-01-14 18:16:24</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://book.douban.com/review/9895720/">读书的点点滴滴</a></h2>

                <div id="review_9895720_short" class="review-short" data-rid="9895720">
                    <div class="short-content">

                        但事实上所所所所所所所所哒哒哒哒哒哒多多多多多多多但事实上所所所所所所所所哒哒哒哒哒哒多多多多多多多但事实上所所所所所所所所哒哒哒哒哒哒多多多多多多多但事实上所所所所所所所所哒哒哒哒哒哒多多多多多多多但事实上所所所所所所所所哒哒哒哒哒哒多多多多多多多但事实...

                        &nbsp;(<a href="javascript:;" id="toggle-9895720-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_9895720_full" class="hidden">
                    <div id="review_9895720_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="9895720" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-9895720">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="9895720" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-9895720">
                        </span>
                    </a>
                    <a href="https://book.douban.com/review/9895720/#comments" class="reply">1回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>




    

    

    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/68afa4e1a12ae93d.js"></script>
    <!-- COLLECTED CSS -->
</div>








            <p class="pl">
                &gt;
                <a href="reviews">
                    更多书评11篇
                </a>
            </p>
</section>

<!-- COLLECTED JS -->

  









<div class="ugc-mod reading-notes">
  <div class="hd">
    <h2>
      读书笔记&nbsp;&nbsp;· · · · · ·&nbsp;
          <span class="pl">(<a href="https://book.douban.com/subject/30389935/annotation">共<span property="v:count">11</span>篇</a>)</span>

    </h2>

      <a class="redbutt rr j a_show_login" href="https://www.douban.com/register?reason=annotate" rel="nofollow"><span>我来写笔记</span></a>
  </div>
  

    <div class="bd">
      <ul class="inline-tabs">
        <li class="on"><a href="#rank" id="by_rank" >按有用程度</a></li>
        <li><a href="#page" id="by_page" >按页码先后</a></li>
        <li><a href="#time" id="by_time" >最新笔记</a></li>
      </ul>
      
  <ul class="comments by_rank" >
      
      <li class="ctsh clearfix" data-cid="68054575">
        <div class="ilst">
          <a href="https://www.douban.com/people/134853618/"><img src="https://img3.doubanio.com/icon/u134853618-2.jpg" alt="三三" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68054575/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68054575/" class="">第45页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/134853618/" class=" " title="三三">三三</a>
                (苏世独立，横而不流)
              
            </p>
            <div class="reading-note" data-cid="68054575">
              <div class="short">
                
                  <span class="">但走到这一步，我已经成功了，不奢望成就的话，就不会失败。</span>
                <p class="pl">
                  <span class="">2019-01-07 15:20</span>
                  
                    &nbsp;&nbsp;<span class="">1人喜欢</span>
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>但走到这一步，我已经成功了，不奢望成就的话，就不会失败。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68054575/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-07 15:20
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="69122799">
        <div class="ilst">
          <a href="https://www.douban.com/people/rockpiano/"><img src="https://img1.doubanio.com/icon/u1902777-18.jpg" alt="仙人掌C" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/69122799/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/69122799/" class="">第4页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/rockpiano/" class=" " title="仙人掌C">仙人掌C</a>
                (不合时宜的旁白小姐)
              
            </p>
            <div class="reading-note" data-cid="69122799">
              <div class="short">
                
                  <span class="">除夕停电，爸爸说：“哈哈哈，好玩。” 其实一些无伤大雅的“意外”，比起按部就班的“顺利”，更能让生活充满色彩。</span>
                <p class="pl">
                  <span class="">2019-01-20 01:54</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>除夕停电，爸爸说：“哈哈哈，好玩。”</p><p>其实一些无伤大雅的“意外”，比起按部就班的“顺利”，更能让生活充满色彩。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/69122799/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-20 01:54
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="68989392">
        <div class="ilst">
          <a href="https://www.douban.com/people/Aterego/"><img src="https://img1.doubanio.com/icon/u55667344-7.jpg" alt="小铁" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68989392/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68989392/" class="">第186页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/Aterego/" class=" " title="小铁">小铁</a>
                (万物自然生长。)
              
                <span class="allstar50" title="力荐"></span>
            </p>
            <div class="reading-note" data-cid="68989392">
              <div class="short">
                
                  <span class="">每年春节都是欢娱和阵痛的交织。有一年节后离家，刚到火车站我就收到了妈的短信：“早知道心里这么难受，你们明年干脆别回家过年了，我和你爸平时清清静静惯了。回来几天又走，家里刚一热闹又冷清下来，我们受不了。刚才想叫你下来吃面，才想起你已经走了。” 我一个壮如蛮牛的大老爷们，居然从进站口哭到了车上。放妥行李，坐下看滑过车窗的独山城，想起临别时爸跟我走到街角，妈直到我们拐弯仍然倚在门口，手扶铁门，我又忍不...</span>
                <p class="pl">
                  <span class="">2019-01-18 13:40</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>每年春节都是欢娱和阵痛的交织。有一年节后离家，刚到火车站我就收到了妈的短信：“早知道心里这么难受，你们明年干脆别回家过年了，我和你爸平时清清静静惯了。回来几天又走，家里刚一热闹又冷清下来，我们受不了。刚才想叫你下来吃面，才想起你已经走了。”</p><p>我一个壮如蛮牛的大老爷们，居然从进站口哭到了车上。放妥行李，坐下看滑过车窗的独山城，想起临别时爸跟我走到街角，妈直到我们拐弯仍然倚在门口，手扶铁门，我又忍不住泣不成声。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68989392/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-18 13:40
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="68989161">
        <div class="ilst">
          <a href="https://www.douban.com/people/Aterego/"><img src="https://img1.doubanio.com/icon/u55667344-7.jpg" alt="小铁" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68989161/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68989161/" class="">第155页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/Aterego/" class=" " title="小铁">小铁</a>
                (万物自然生长。)
              
                <span class="allstar50" title="力荐"></span>
            </p>
            <div class="reading-note" data-cid="68989161">
              <div class="short">
                
                  <span class="">我妈天生暴脾气，见不得不平事，眼睛一瞪，路灯都要黯淡几分；又争强好胜不服输，眉头下从没写过“困难”二字。外公生前逢人就说：这丫头投错胎了，要是个男娃就刚好！</span>
                <p class="pl">
                  <span class="">2019-01-18 13:30</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>我妈天生暴脾气，见不得不平事，眼睛一瞪，路灯都要黯淡几分；又争强好胜不服输，眉头下从没写过“困难”二字。外公生前逢人就说：这丫头投错胎了，要是个男娃就刚好！</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68989161/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-18 13:30
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
  </ul>
  

      
  <ul class="comments by_page"  style="display: none">
      
      <li class="ctsh clearfix" data-cid="68982376">
        <div class="ilst">
          <a href="https://www.douban.com/people/Aterego/"><img src="https://img1.doubanio.com/icon/u55667344-7.jpg" alt="小铁" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68982376/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68982376/" class="">第1页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/Aterego/" class=" " title="小铁">小铁</a>
                (万物自然生长。)
              
                <span class="allstar50" title="力荐"></span>
            </p>
            <div class="reading-note" data-cid="68982376">
              <div class="short">
                
                  <span class="">照片是我的日记 帮助有限的脑容量记录下走过的路 因为时间过得太快 我怕拥有的幸福太容易失去</span>
                <p class="pl">
                  <span class="">2019-01-18 10:14</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>照片是我的日记</p><p>帮助有限的脑容量记录下走过的路</p><p>因为时间过得太快</p><p>我怕拥有的幸福太容易失去</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68982376/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-18 10:14
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="68680058">
        <div class="ilst">
          <a href="https://www.douban.com/people/190074342/"><img src="https://img3.doubanio.com/icon/u190074342-1.jpg" alt="跨越龙门" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68680058/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68680058/" class="">第2页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/190074342/" class=" " title="跨越龙门">跨越龙门</a>
              
            </p>
            <div class="reading-note" data-cid="68680058">
              <div class="short">
                
                  <div class="ll">
                    <a href="https://book.douban.com/annotation/68680058/"><img src="https://img3.doubanio.com/view/page_note/small/public/p68680058-1.jpg"></a>
                  </div>
                  <span class="">在这科技时代高速发展的今天，黑网铺天盖地，很多黑网也逐渐浮出水面，据网友不完全统计目前黑网多达几百家，真是骇人听闻！遇到黑网该怎样维护自己的资金安全？这是一个很关键，也很值得我们深思的主题。 很多朋友打网投不给出款确实是一件让人头疼的事情，不管是谁遇到了都会脑瓜疼。但是真的有人可以帮你出款吗？其实，是没有的，不要天真了。那些说是能出款的都是骗你而已，他会给你各种理由说是不能出，这样你可能再次给黑...</span>
                  &nbsp;&nbsp;<a href="https://book.douban.com/annotation/68680058/">(<span class="">1回应</span>)</a>
                <p class="pl">
                  <span class="">2019-01-14 23:15</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>在这科技时代高速发展的今天，黑网铺天盖地，很多黑网也逐渐浮出水面，据网友不完全统计目前黑网多达几百家，真是骇人听闻！遇到黑网该怎样维护自己的资金安全？这是一个很关键，也很值得我们深思的主题。</p><p></p><div class="image-container image-float-center"><div class="image-wrapper"><img src="https://img3.doubanio.com/view/page_note/large/public/p68680058-1.jpg" width=""></div></div><p></p><p>很多朋友打网投不给出款确实是一件让人头疼的事情，不管是谁遇到了都会脑瓜疼。但是真的有人可以帮你出款吗？其实，是没有的，不要天真了。那些说是能出款的都是骗你而已，他会给你各种理由说是不能出，这样你可能再次给黑。因此，不要相信那些出黑的了。因为是没有人可以真正的出款的，出黑只是一个幌子。没有现场的平台就没有不黑人的，记住一定要学会辨别真伪，没有实体现场一切都不可信。那么，遇到这种情况，该如何应对？首先是保持冷静，不要心急。只要账号能正常登录、额度能够转换，就还有机会挽回。主题咨询扣193-966-098，很多朋友都成功的，任何事情都有解决的办法。</p><p></p><div class="image-container image-float-center"><div class="image-wrapper"><img src="https://img3.doubanio.com/view/page_note/large/public/p68680058-2.jpg" width=""></div></div><p></p><p>然后学会伪装：通过与客服交流的谈话中不经意的像客服透露自己的经济实力，（往不差钱的方向去说但要掌握分寸）让客服觉得你身上有很多他们想要的价值和有继续利用你能赚取更多利益的想法。（做到这一步后在适当的装傻迷惑客服上路即可）如果你成功完成了以上操作的话你的资金基本就出来了。还有就是装傻:完成以上操作之后你就可以引蛇出洞了，适当的装傻迷惑客服上路即可,如果你成功完成了以上操作的话你的资金基本就出来了。 老弟还是的说几句赌博是害人的，赶紧上岸吧！</p><p></p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68680058/#comments">1回应</a>&nbsp;&nbsp;
                  2019-01-14 23:15
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="69122799">
        <div class="ilst">
          <a href="https://www.douban.com/people/rockpiano/"><img src="https://img1.doubanio.com/icon/u1902777-18.jpg" alt="仙人掌C" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/69122799/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/69122799/" class="">第4页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/rockpiano/" class=" " title="仙人掌C">仙人掌C</a>
                (不合时宜的旁白小姐)
              
            </p>
            <div class="reading-note" data-cid="69122799">
              <div class="short">
                
                  <span class="">除夕停电，爸爸说：“哈哈哈，好玩。” 其实一些无伤大雅的“意外”，比起按部就班的“顺利”，更能让生活充满色彩。</span>
                <p class="pl">
                  <span class="">2019-01-20 01:54</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>除夕停电，爸爸说：“哈哈哈，好玩。”</p><p>其实一些无伤大雅的“意外”，比起按部就班的“顺利”，更能让生活充满色彩。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/69122799/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-20 01:54
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="68054575">
        <div class="ilst">
          <a href="https://www.douban.com/people/134853618/"><img src="https://img3.doubanio.com/icon/u134853618-2.jpg" alt="三三" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68054575/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68054575/" class="">第45页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/134853618/" class=" " title="三三">三三</a>
                (苏世独立，横而不流)
              
            </p>
            <div class="reading-note" data-cid="68054575">
              <div class="short">
                
                  <span class="">但走到这一步，我已经成功了，不奢望成就的话，就不会失败。</span>
                <p class="pl">
                  <span class="">2019-01-07 15:20</span>
                  
                    &nbsp;&nbsp;<span class="">1人喜欢</span>
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>但走到这一步，我已经成功了，不奢望成就的话，就不会失败。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68054575/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-07 15:20
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
  </ul>
  

      
  <ul class="comments by_time"  style="display: none">
      
      <li class="ctsh clearfix" data-cid="69122799">
        <div class="ilst">
          <a href="https://www.douban.com/people/rockpiano/"><img src="https://img1.doubanio.com/icon/u1902777-18.jpg" alt="仙人掌C" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/69122799/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/69122799/" class="">第4页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/rockpiano/" class=" " title="仙人掌C">仙人掌C</a>
                (不合时宜的旁白小姐)
              
            </p>
            <div class="reading-note" data-cid="69122799">
              <div class="short">
                
                  <span class="">除夕停电，爸爸说：“哈哈哈，好玩。” 其实一些无伤大雅的“意外”，比起按部就班的“顺利”，更能让生活充满色彩。</span>
                <p class="pl">
                  <span class="">2019-01-20 01:54</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>除夕停电，爸爸说：“哈哈哈，好玩。”</p><p>其实一些无伤大雅的“意外”，比起按部就班的“顺利”，更能让生活充满色彩。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/69122799/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-20 01:54
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="68989392">
        <div class="ilst">
          <a href="https://www.douban.com/people/Aterego/"><img src="https://img1.doubanio.com/icon/u55667344-7.jpg" alt="小铁" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68989392/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68989392/" class="">第186页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/Aterego/" class=" " title="小铁">小铁</a>
                (万物自然生长。)
              
                <span class="allstar50" title="力荐"></span>
            </p>
            <div class="reading-note" data-cid="68989392">
              <div class="short">
                
                  <span class="">每年春节都是欢娱和阵痛的交织。有一年节后离家，刚到火车站我就收到了妈的短信：“早知道心里这么难受，你们明年干脆别回家过年了，我和你爸平时清清静静惯了。回来几天又走，家里刚一热闹又冷清下来，我们受不了。刚才想叫你下来吃面，才想起你已经走了。” 我一个壮如蛮牛的大老爷们，居然从进站口哭到了车上。放妥行李，坐下看滑过车窗的独山城，想起临别时爸跟我走到街角，妈直到我们拐弯仍然倚在门口，手扶铁门，我又忍不...</span>
                <p class="pl">
                  <span class="">2019-01-18 13:40</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>每年春节都是欢娱和阵痛的交织。有一年节后离家，刚到火车站我就收到了妈的短信：“早知道心里这么难受，你们明年干脆别回家过年了，我和你爸平时清清静静惯了。回来几天又走，家里刚一热闹又冷清下来，我们受不了。刚才想叫你下来吃面，才想起你已经走了。”</p><p>我一个壮如蛮牛的大老爷们，居然从进站口哭到了车上。放妥行李，坐下看滑过车窗的独山城，想起临别时爸跟我走到街角，妈直到我们拐弯仍然倚在门口，手扶铁门，我又忍不住泣不成声。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68989392/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-18 13:40
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="68989161">
        <div class="ilst">
          <a href="https://www.douban.com/people/Aterego/"><img src="https://img1.doubanio.com/icon/u55667344-7.jpg" alt="小铁" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68989161/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68989161/" class="">第155页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/Aterego/" class=" " title="小铁">小铁</a>
                (万物自然生长。)
              
                <span class="allstar50" title="力荐"></span>
            </p>
            <div class="reading-note" data-cid="68989161">
              <div class="short">
                
                  <span class="">我妈天生暴脾气，见不得不平事，眼睛一瞪，路灯都要黯淡几分；又争强好胜不服输，眉头下从没写过“困难”二字。外公生前逢人就说：这丫头投错胎了，要是个男娃就刚好！</span>
                <p class="pl">
                  <span class="">2019-01-18 13:30</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>我妈天生暴脾气，见不得不平事，眼睛一瞪，路灯都要黯淡几分；又争强好胜不服输，眉头下从没写过“困难”二字。外公生前逢人就说：这丫头投错胎了，要是个男娃就刚好！</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68989161/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-18 13:30
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
      
      <li class="ctsh clearfix" data-cid="68989127">
        <div class="ilst">
          <a href="https://www.douban.com/people/Aterego/"><img src="https://img1.doubanio.com/icon/u55667344-7.jpg" alt="小铁" class="" /></a>
        </div>
        <div class="con">
          <div class="nlst">
            <h3>
              <div class="note-toggle rr">
                <a href="https://book.douban.com/annotation/68989127/" class="note-unfolder">展开</a>
                <a href="javascript:void(0);" class="note-folder">收起</a>
              </div>
              <a href="https://book.douban.com/annotation/68989127/" class="">第148页</a></h3>
          </div>
          <div class="clst">
            <p class="user"><a href="https://www.douban.com/people/Aterego/" class=" " title="小铁">小铁</a>
                (万物自然生长。)
              
                <span class="allstar50" title="力荐"></span>
            </p>
            <div class="reading-note" data-cid="68989127">
              <div class="short">
                
                  <span class="">再比如，有了你喜欢的食物，他看似不经意地把东西放在你面前就去做其他事了，什么都不说。哪怕这也是他最喜欢的，只要你爱吃，他就一口都不动，全都留给你。若是生病了，谁也不告诉，自己恹恹地去买药，病容却是掩藏不住的，我小时候曾见过他发高烧时往自己屁股上扎针。他不愿意让人担心，更不喜欢麻烦人，哪怕是自己的孩子。</span>
                <p class="pl">
                  <span class="">2019-01-18 13:27</span>
                  
                </p>
              </div>
              <div class="all hidden" style="display:none" >
                <p>再比如，有了你喜欢的食物，他看似不经意地把东西放在你面前就去做其他事了，什么都不说。哪怕这也是他最喜欢的，只要你爱吃，他就一口都不动，全都留给你。若是生病了，谁也不告诉，自己恹恹地去买药，病容却是掩藏不住的，我小时候曾见过他发高烧时往自己屁股上扎针。他不愿意让人担心，更不喜欢麻烦人，哪怕是自己的孩子。</p>
                  <div class="col-rec-con clearfix">
                    







<div class="rec-sec">

    <span class="rec">

<a href="https://www.douban.com/accounts/register?reason=collect" class="j a_show_login lnk-sharing lnk-douban-sharing">推荐</a>
</span>
</div>

                  </div>
                <div class="pl col-time">
                  <a href="https://book.douban.com/annotation/68989127/#comments">回应</a>&nbsp;&nbsp;
                  2019-01-18 13:27
                </div>
              </div>
            </div>
          </div>
        </div>
      </li>
  </ul>
  

    </div>
      <div class="ft">
        <p class="trr">&gt; <a href="https://book.douban.com/subject/30389935/annotation">更多读书笔记（共11篇）</a></p>
      </div>

</div>



<script type="text/javascript">
  $(document).ready(function(){
    var TEMPL_ADD_COL = '<a href="" id="n-{NOTE_ID}" class="colbutt ll add-col"><span>收藏</span></a>',
      TEMPL_DEL_COL = '<span class="pl">已收藏 &gt;<a href="" id="n-{NOTE_ID}" class="del-col">取消收藏</a></span>';

    $('body').delegate('.add-col', 'click', function(e){
      e.preventDefault();
      var nnid = $(this).attr('id').split('-')[1],
        bn_add = $(this);
      $.post_withck("/j/book/annotation_collect",{nid:nnid},function(){
        var a = $(TEMPL_DEL_COL.replace('{NOTE_ID}', nnid));
        bn_add.before(a);
        bn_add.remove();
      })
    });

    $('body').delegate('.del-col', 'click', function(e){
      e.preventDefault();
      var nnid = $(this).attr('id').split('-')[1],
        bn_del = $(this).parent();
      $.post_withck("/j/book/annotation_uncollect", {nid: nnid}, function() {
        var a = $(TEMPL_ADD_COL.replace('{NOTE_ID}', nnid));
        bn_del.before(a);
        bn_del.remove();
      })
    });

    $("pre.source").each(function(){
      var cn = $(this).attr('class').split(' ');
      l = cn[1];
      s = 'rand01';
      n = cn[2];
      $(this).snippet(n,{style: s, showNum: l});
    });

    var annotationMod = $('.reading-notes .bd')
      , annotationTabs = annotationMod.find('.inline-tabs li')
      , annotationTabLinks = annotationTabs.find('a')
      , annotationTabContents = annotationMod.find('ul.comments');

    annotationTabLinks.click(function(e){
      e.preventDefault();
      var el = $(this)
        , kind = el.attr('id');

      annotationTabs.removeClass('on');
      el.parent().addClass('on');
      annotationTabContents.hide();
      annotationTabContents.filter('.' + kind).show();
    });
  });
</script>

<script type="text/x-mathjax-config">
MathJax.Hub.Config({
	jax: ["input/TeX", "output/HTML-CSS"],
    extensions: ["tex2jax.js","TeX/AMSmath.js","TeX/AMSsymbols.js","TeX/noUndefined.js"],
    tex2jax: {
		inlineMath: [ ["($", "$)"], ['\\(','\\)'] ],
		displayMath: [ ["($$","$$)"], ['\\[','\\]']],
		skipTags: ["script","noscript","style","textarea"],
		processEscapes: true,
		processEnvironments: true,
		preview: "TeX"
	},
	showProcessingMessages: false
  });
</script>


  






<div id="db-discussion-section" class="indent ugc-mod">




        <span>
            <a class="redbutt rr j a_show_login" href="https://www.douban.com/register?reason=discuss">
              <span>在这本书的论坛里发言</span>
            </a>
        </span>


</div>




</div>
<script type="text/javascript">
$(function(){if($.browser.msie && $.browser.version == 6.0){
    var maxWidth = parseInt($('#info').css('max-width'));
    if($('#info').width() > maxWidth)
        $('#info').width(maxWidth)
}});
</script>
</div>
      <div class="aside">
        
  
  






  <div id="dale_book_subject_top_right"></div>
    
  
  
  <div class="gray_ad">
    
  

  <h2>
    <span class="">在豆瓣书店有售</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


    <div class="market-banner">
      <span class="title">
        纸质版&nbsp;
      </span>
      <span class="price"> 39.20元</span>
        <span class="price"> <del>49.00元</del></span>
      <span class="promotion-info">满48元包邮</span>
      <span class="actions">
        <a class="buy-btn buy" target="_blank" href="https://market.douban.com/cart/checkout/?sku_id=268945&utm_campaign=douban_book_subject_buy_btn&utm_source=douban&utm_medium=pc_web">
          去购买
        </a>
          <a class="j a_show_login buy-btn cart" href="javascript:;">加入购物车</a>
      </span>
    </div>
  </div>


  






<style type="text/css" media="screen">
  .add2cartContainer{overflow:hidden;vertical-align:bottom;line-height:1}.add2cartContainer .add2cart{margin-right:0;display:inline-block}#buyinfo .bs{margin:0}#buyinfo li{padding:10px 0;position:relative;line-height:1;border-bottom:1px solid #eaeaea}#buyinfo li a:hover{background-image:none !important}#buyinfo li a:hover .buylink-price{background:#37a}#buyinfo li .publish,#buyinfo li .other-activity{margin-top:5px}#buyinfo li .publish a,#buyinfo li .other-activity a{color:#999}#buyinfo li .publish a:hover,#buyinfo li .other-activity a:hover{color:#37a;background:none;opacity:0.5;filter:alpha(opacity=50)}#buyinfo li .buylink-price{position:absolute;right:90px;text-align:right}#buyinfo .more-info{color:#aaa;margin:6px 0 -2px 0}#buyinfo .more-ebooks{padding:10px 0;color:#37a;cursor:pointer}#buyinfo .price-page{border-bottom:0;padding:15px 0 0}#buyinfo .saved-price{display:none;margin-left:5px}#buyinfo .cart-tip{float:right;padding-right:5px}#buyinfo #buyinfo-ebook{margin-bottom:15px}#buyinfo #buyinfo-ebook .buylink-price{display:inline}#buyinfo #buyinfo-ebook li.no-border{border:0}#buyinfo-printed{margin-bottom:15px}#buyinfo-printed.no-border{border-bottom:0}#buyinfo-printed .more-ebooks{line-height:1;padding:10px 0;color:#37a;cursor:pointer;padding:10px 0 0}#buyinfo-report{display:none}#buyinfo-report .lnk-close-report{float:right;margin-top:-30px;line-height:14px}#buyinfo-report .item{margin-bottom:10px}#buyinfo-report .item input{vertical-align:text-bottom;*vertical-align:middle}#buyinfo-report .item label{margin:0 5px 0 2px}#buyinfo-report .item-submit .bn-flat{margin-right:10px}#buyinfo-report .item-price input{width:220px;border:1px solid #ccc;padding:4px}#buyinfo-report form{margin:5px 0 10px}#bi-report-btn{float:right;margin:2px 0;line-height:14px;display:none}.bi-vendor-report{color:#aaa}.bi-vendor-report-form{display:none;color:#111;margin:0 5px;line-height:25px}.gray_ad{padding:30px 20px 25px 20px;background:#f6f6f1}.gray_ad h2{margin-bottom:6px;font-size:15px}.gray_ad .ebook-tag{margin-top:5px;color:#999;font-size:12px}.bs.more-after{margin-bottom:0px}@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi){#buyinfo li a:hover{background-image:url(https://img3.doubanio.com/f/book/fc4ff7f0a3a7f452f06d586540284b9738f2fe87/pics/book/cart/icon-brown@2x.png);background-size:16px 12px}}#intervenor-buyinfo .bs{margin:0}#intervenor-buyinfo li{position:relative;border-bottom:1px solid #eaeaea;padding:10px 0;line-height:1}#intervenor-buyinfo li .basic-info{color:#494949;font-size:14px;line-height:18px}#intervenor-buyinfo li a:hover .comment{color:#f67;opacity:0.75;filter:alpha(opacity=75)}#intervenor-buyinfo li a:hover .buy-btn{background:#fff;border:1px solid #e97e7e;border-radius:2px;color:#e97e7e}#intervenor-buyinfo li a:hover .buylink-price{background:#37a}#intervenor-buyinfo li .buylink-price{position:absolute;right:90px;text-align:right}#intervenor-buyinfo li .publish,#intervenor-buyinfo li .other-activity{margin-top:5px}#intervenor-buyinfo li .publish a,#intervenor-buyinfo li .other-activity a{color:#999}#intervenor-buyinfo li .publish a:hover,#intervenor-buyinfo li .other-activity a:hover{color:#37a;background:none;opacity:0.5;filter:alpha(opacity=50)}#intervenor-buyinfo .jd-buy-icon{float:left;margin-right:3px}#intervenor-buyinfo .buy-btn{float:right;position:absolute;right:10px;bottom:3px;color:#9c9c9c;padding:0 12px;border:1px solid transparent}#intervenor-buyinfo .comment{color:#FF8C9C;margin:6px 0 -2px 0}#intervenor-buyinfo .price-page a{display:inline-block;line-height:16px !important}#intervenor-buyinfo .price-page{border-bottom:0;padding:15px 0 0}#intervenor-buyinfo .saved-price{display:none;margin-left:5px}#intervenor-buyinfo .cart-tip{float:right;padding-right:5px}#intervenor-buyinfo #buyinfo-ebook{margin-bottom:15px}#intervenor-buyinfo #buyinfo-ebook .buylink-price{display:inline}#intervenor-buyinfo #buyinfo-ebook li.no-border{border:0}#buyinfo-printed .presale-indicator{margin:0;width:auto;color:#999;text-indent:0;background:none}

</style>

      <div class="gray_ad" id="buyinfo">
      <div id="buyinfo-printed">
        
  

  <h2>
    <span class="">在哪儿买这本书</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;

  </h2>


        <ul class="bs noline more-after ">
          
                
                <li class="">
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=3820&amp;pre=0&amp;vendor=jingdong&amp;srcpage=subject&amp;price=3820&amp;pos=1&amp;url=https%3A%2F%2Funion-click.jd.com%2Fjdc%3Fe%3D%26p%3DAyIHZRtYFAcXBFIZWR0yEgRQHVMQBhQ3EUQDS10iXhBeGh4cDF8QTwcKWUcYB0UHCwIRAlMTXhEEDV4QRwYlWmBQKUEbTGlyAE9zGE5mb1sAeRNFch4LZRxeEgQVDlYcaxUGEwNRK2sVAyJVO8Pto9q3tknP8ZrXmpBlGmsVBhcAUx1ZHQESDl0caxIyy4fizfmH24yNgJPMJTIiN2UrWyUBIlgRRgYlAw%253D%253D%26t%3DW1dCFBBFC1pXUwkEAEAdQFkJBVsWBxQPUB9dCltXWwg%253D&amp;cntvendor=2&amp;srcsubj=30389935&amp;type=bkbuy&amp;subject=30389935" class="">
                      <span class="">京东商城</span>
                    </a>
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=3820&amp;pre=0&amp;vendor=jingdong&amp;srcpage=subject&amp;price=3820&amp;pos=1&amp;url=https%3A%2F%2Funion-click.jd.com%2Fjdc%3Fe%3D%26p%3DAyIHZRtYFAcXBFIZWR0yEgRQHVMQBhQ3EUQDS10iXhBeGh4cDF8QTwcKWUcYB0UHCwIRAlMTXhEEDV4QRwYlWmBQKUEbTGlyAE9zGE5mb1sAeRNFch4LZRxeEgQVDlYcaxUGEwNRK2sVAyJVO8Pto9q3tknP8ZrXmpBlGmsVBhcAUx1ZHQESDl0caxIyy4fizfmH24yNgJPMJTIiN2UrWyUBIlgRRgYlAw%253D%253D%26t%3DW1dCFBBFC1pXUwkEAEAdQFkJBVsWBxQPUB9dCltXWwg%253D&amp;cntvendor=2&amp;srcsubj=30389935&amp;type=bkbuy&amp;subject=30389935" class="buylink-price ">
                      <span class="">
                        38.20 元
                      </span>
                    </a>

                      
                </li>
                
                <li class="">
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=3820&amp;pre=0&amp;vendor=dangdang&amp;srcpage=subject&amp;price=3820&amp;pos=2&amp;url=http%3A%2F%2Funion.dangdang.com%2Ftransfer.php%3Ffrom%3DP-306226-0-s30389935%26backurl%3Dhttp%3A%2F%2Fproduct.dangdang.com%2Fproduct.aspx%3Fproduct_id%3D26445258&amp;cntvendor=2&amp;srcsubj=30389935&amp;type=bkbuy&amp;subject=30389935" class="">
                      <span class="">当当网</span>
                    </a>
                    <a target="_blank" href="https://book.douban.com/link2/?lowest=3820&amp;pre=0&amp;vendor=dangdang&amp;srcpage=subject&amp;price=3820&amp;pos=2&amp;url=http%3A%2F%2Funion.dangdang.com%2Ftransfer.php%3Ffrom%3DP-306226-0-s30389935%26backurl%3Dhttp%3A%2F%2Fproduct.dangdang.com%2Fproduct.aspx%3Fproduct_id%3D26445258&amp;cntvendor=2&amp;srcsubj=30389935&amp;type=bkbuy&amp;subject=30389935" class="buylink-price ">
                      <span class="">
                        38.20 元
                      </span>
                    </a>

                      
                        <div class="more-info">
                            <span class="buyinfo-promotion">
                              自出版加价购
                            </span>
                        </div>
                </li>
          <li class="price-page">
            <a href="https://book.douban.com/subject/30389935/buylinks">
              &gt; 查看2家网店价格
                (38.20 元起)
            </a>
          </li>
        </ul>
      </div>
      
  <div class="add2cartContainer ft no-border">
    
  <span class="add2cartWidget ll">
      <a class="j  add2cart a_show_login" href="https://www.douban.com/register?reason=addbook2cart" rel="nofollow">
        <span>+ 加入购书单</span></a>
  </span>
  

  </div>

  </div>
  <script type="text/javascript">
  $('.more-ebooks').on('click', function() {
    var $this = $(this),
      $li = $this.siblings('ul').find('li');
    if ($this.hasClass('isShow')) {
      $(this).text('展开更多').removeClass('isShow');
      $li.not(':first').addClass('hide');
    }else{
      $(this).text('收起').addClass('isShow');
      $li.removeClass('hide');
    }
    
  })
  </script>

<style class="text/css">
  .presale-indicator {
  display: inline-block;
  *display: inline;
  *zoom: 1;
  width: 24px;
  height: 15px;
  line-height: 15px;
  background: url(https://img3.doubanio.com/f/book/1679c65572eac1371f9872807199dea6e55a7f06/pics/book/cart/presale_text.gif) no-repeat;
  text-indent: -9999px;
  vertical-align: middle;
  *vertical-align: 0px;
  _vertical-align: 2px;
  margin-left: 0.5em;
}

</style>



  






  <div id="dale_book_subject_top_middle"></div>
  





  



  



      
  

  <h2>
    <span class="">以下豆列推荐</span>
      &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;
      <span class="pl">&nbsp;(
          <a href="https://book.douban.com/subject/30389935/doulists">全部</a>
        ) </span>

  </h2>


    <div id="db-doulist-section" class="indent">
      <ul class="bs">
          <li class=""><a class="" href="https://www.douban.com/doulist/36764655/" target="_blank">2015-2019好书新发现</a>
                <span class="pl">(Moon)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/27039041/" target="_blank">我的身体里有一个游荡的未来</a>
                <span class="pl">(诗菡)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/38991156/" target="_blank">生活中永远保持期待：好书追寻中</a>
                <span class="pl">(无心恋战)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/3349252/" target="_blank">购书单</a>
                <span class="pl">(波豆豆)</span>
            </li>
          <li class=""><a class="" href="https://www.douban.com/doulist/110803224/" target="_blank">2019华文电影图书</a>
                <span class="pl">(妖灵妖)</span>
            </li>
      </ul>
    </div>

  <div id="dale_book_subject_middle_mini"></div>
  






  <h2>谁读这本书?</h2>
  <div class="indent" id="collector">

    

<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/170684699/"><img src="https://img3.doubanio.com/icon/u170684699-2.jpg" class="pil" alt="深蓝" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/170684699/">深蓝</a><br/>
      <div class="pl ll">          32分钟前          想读      </div>

      <br/>

    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>
<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/103823087/"><img src="https://img3.doubanio.com/icon/u103823087-1.jpg" class="pil" alt="小Kristen" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/103823087/">小Kristen</a><br/>
      <div class="pl ll">          33分钟前          读过      </div>

        <span class="allstar40" title="推荐"></span>
      <br/>

    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>
<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/ytxwz1978/"><img src="https://img1.doubanio.com/icon/u2389146-37.jpg" class="pil" alt="樱桃小丸子" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/ytxwz1978/">樱桃小丸子</a><br/>
      <div class="pl ll">          47分钟前          读过      </div>

        <span class="allstar30" title="还行"></span>
      <br/>

    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>
<div class="">
    
    <div class="ll"><a href="https://www.douban.com/people/68825809/"><img src="https://img1.doubanio.com/icon/u68825809-17.jpg" class="pil" alt="海哈哈" /></a></div>
    <div style="padding-left:60px"><a class="" href="https://www.douban.com/people/68825809/">海哈哈</a><br/>
      <div class="pl ll">          52分钟前          读过      </div>

        <span class="allstar40" title="推荐"></span>
      <br/>

    </div>
    <div class="clear"></div><br/>
    <div class="ul" style="margin-bottom:12px;"></div>
</div>


        <p class="pl">&gt; <a href="https://book.douban.com/subject/30389935/doings">90人在读</a></p>
        <p class="pl">&gt; <a href="https://book.douban.com/subject/30389935/collections">298人读过</a></p>
        <p class="pl">&gt; <a href="https://book.douban.com/subject/30389935/wishes">1717人想读</a></p>

  </div>





  
<!-- douban ad begin -->
<div id="dale_book_subject_middle_right"></div>
<script type="text/javascript">
    (function (global) {
        if(!document.getElementsByClassName) {
            document.getElementsByClassName = function(className) {
                return this.querySelectorAll("." + className);
            };
            Element.prototype.getElementsByClassName = document.getElementsByClassName;

        }
        var articles = global.document.getElementsByClassName('article'),
            asides = global.document.getElementsByClassName('aside');

        if (articles.length > 0 && asides.length > 0 && articles[0].offsetHeight >= asides[0].offsetHeight) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_book_subject_middle_right');
        }
    })(this);
</script>
<!-- douban ad end -->

  





  

  <h2>二手市场</h2>
  <div class="indent">
    <ul class="bs">
    <li class="">
          <a class="rr j a_show_login" href="https://www.douban.com/register?reason=secondhand-offer&amp;cat=book"><span>&gt; 点这儿转让</span></a>

        有1717人想读,手里有一本闲着?
      </li>
    </ul>
  </div>

  
<p class="pl">订阅关于四个春天的评论: <br/><span class="feed">
    <a href="https://book.douban.com/feed/subject/30389935/reviews"> feed: rss 2.0</a></span></p>


      </div>
      <div class="extra">
        
  
<!-- douban ad begin -->
<div id="dale_book_subject_bottom_super_banner"></div>
<script type="text/javascript">
    (function (global) {
        var body = global.document.body,
            html = global.document.documentElement;

        var height = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight, html.scrollHeight, html.offsetHeight);
        if (height >= 2000) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_book_subject_bottom_super_banner');
        }
    })(this);
</script>
<!-- douban ad end -->


      </div>
    </div>
  </div>

        
<div id="footer">

<span id="icp" class="fleft gray-link">
    &copy; 2005－2019 douban.com, all rights reserved 北京豆网科技有限公司
</span>

<a href="https://www.douban.com/hnypt/variformcyst.py" style="display: none;"></a>

<span class="fright">
    <a href="https://www.douban.com/about">关于豆瓣</a>
    · <a href="https://www.douban.com/jobs">在豆瓣工作</a>
    · <a href="https://www.douban.com/about?topic=contactus">联系我们</a>
    · <a href="https://www.douban.com/about?policy=disclaimer">免责声明</a>
    
    · <a href="https://help.douban.com/?app=book" target="_blank">帮助中心</a>
    · <a href="https://book.douban.com/library_invitation">图书馆合作</a>
    · <a href="https://www.douban.com/doubanapp/">移动应用</a>
    · <a href="https://www.douban.com/partner/">豆瓣广告</a>
</span>

</div>

    </div>
      
  

    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/7dccdcfb50652052.js"></script>
    <!-- mako -->
    
  








    
<script type="text/javascript">
    (function (global) {
        var newNode = global.document.createElement('script'),
            existingNode = global.document.getElementsByTagName('script')[0],
            adSource = '//erebor.douban.com/',
            userId = '',
            browserId = 'auSsK8Dk5cg',
            criteria = '7:四个春天|7:随笔|7:陆庆屹|7:散文随笔|7:文学|7:2019|7:杂文|7:电影人|7:中国|7:新经典文库|3:/subject/30389935/?icn=index-latestbook-subject',
            preview = '',
            debug = false,
            adSlots = ['dale_book_subject_top_icon', 'dale_book_subject_top_right', 'dale_book_subject_top_middle', 'dale_book_subject_middle_mini'];

        global.DoubanAdRequest = {src: adSource, uid: userId, bid: browserId, crtr: criteria, prv: preview, debug: debug};
        global.DoubanAdSlots = (global.DoubanAdSlots || []).concat(adSlots);

        newNode.setAttribute('type', 'text/javascript');
        newNode.setAttribute('src', 'https://img3.doubanio.com/f/adjs/dd37385211bc8deb01376096bfa14d2c0436a98c/ad.release.js');
        newNode.setAttribute('async', true);
        existingNode.parentNode.insertBefore(newNode, existingNode);
    })(this);
</script>












    
  

<script type="text/javascript">
  var _paq = _paq || [];
  _paq.push(['trackPageView']);
  _paq.push(['enableLinkTracking']);
  (function() {
    var p=(('https:' == document.location.protocol) ? 'https' : 'http'), u=p+'://fundin.douban.com/';
    _paq.push(['setTrackerUrl', u+'piwik']);
    _paq.push(['setSiteId', '100001']);
    var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0]; 
    g.type='text/javascript';
    g.defer=true; 
    g.async=true; 
    g.src=p+'://s.doubanio.com/dae/fundin/piwik.js';
    s.parentNode.insertBefore(g,s);
  })();
</script>

<script type="text/javascript">
var setMethodWithNs = function(namespace) {
  var ns = namespace ? namespace + '.' : ''
    , fn = function(string) {
        if(!ns) {return string}
        return ns + string
      }
  return fn
}

var gaWithNamespace = function(fn, namespace) {
  var method = setMethodWithNs(namespace)
  fn.call(this, method)
}

var _gaq = _gaq || []
  , accounts = [
      { id: 'UA-7019765-1', namespace: 'douban' }
    , { id: 'UA-7019765-16', namespace: '' }
    ]
  , gaInit = function(account) {
      gaWithNamespace(function(method) {
        gaInitFn.call(this, method, account)
      }, account.namespace)
    }
  , gaInitFn = function(method, account) {
      _gaq.push([method('_setAccount'), account.id])

      
  _gaq.push([method('_addOrganic'), 'google', 'q'])
  _gaq.push([method('_addOrganic'), 'baidu', 'wd'])
  _gaq.push([method('_addOrganic'), 'soso', 'w'])
  _gaq.push([method('_addOrganic'), 'youdao', 'q'])
  _gaq.push([method('_addOrganic'), 'so.360.cn', 'q'])
  _gaq.push([method('_addOrganic'), 'sogou', 'query'])
  if (account.namespace) {
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣'])
    _gaq.push([method('_addIgnoredOrganic'), 'douban'])
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣网'])
    _gaq.push([method('_addIgnoredOrganic'), 'www.douban.com'])
  }

      if (account.namespace === 'douban') {
        _gaq.push([method('_setDomainName'), '.douban.com'])
      }

        _gaq.push([method('_setCustomVar'), 1, 'responsive_view_mode', 'desktop', 3])

        _gaq.push([method('_setCustomVar'), 2, 'login_status', '0', 2]);

      _gaq.push([method('_trackPageview')])
    }

for(var i = 0, l = accounts.length; i < l; i++) {
  var account = accounts[i]
  gaInit(account)
}


;(function() {
    var ga = document.createElement('script');
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    ga.setAttribute('async', 'true');
    document.documentElement.firstChild.appendChild(ga);
})()
</script>








    <!-- anson70-docker-->

</body>
</html>




































`
var html44 = `
<!DOCTYPE html>
<html lang="zh-cmn-Hans" class="ua-windows ua-webkit">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="renderer" content="webkit">
    <meta name="referrer" content="always">
    <meta name="google-site-verification" content="ok0wCgT20tBBgo9_zat2iAcimtN4Ftf5ccsh092Xeyw" />
    <title>
        梦想改造家 第一季 (豆瓣)
</title>
    
    <meta name="baidu-site-verification" content="cZdR4xxR7RxmM4zE" />
    <meta http-equiv="Pragma" content="no-cache">
    <meta http-equiv="Expires" content="Sun, 6 Mar 2005 01:00:00 GMT">
    
    <link rel="apple-touch-icon" href="https://img3.doubanio.com/f/movie/d59b2715fdea4968a450ee5f6c95c7d7a2030065/pics/movie/apple-touch-icon.png">
    <link href="https://img3.doubanio.com/f/shire/bf61b1fa02f564a4a8f809da7c7179b883a56146/css/douban.css" rel="stylesheet" type="text/css">
    <link href="https://img3.doubanio.com/f/shire/ae3f5a3e3085968370b1fc63afcecb22d3284848/css/separation/_all.css" rel="stylesheet" type="text/css">
    <link href="https://img3.doubanio.com/f/movie/8864d3756094f5272d3c93e30ee2e324665855b0/css/movie/base/init.css" rel="stylesheet">
    <script type="text/javascript">var _head_start = new Date();</script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/movie/0495cb173e298c28593766009c7b0a953246c5b5/js/movie/lib/jquery.js"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/f010949d3f23dd7c972ad7cb40b800bf70723c93/js/douban.js"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/0efdc63b77f895eaf85281fb0e44d435c6239a3f/js/separation/_all.js"></script>
    
    <meta name="keywords" content="梦想改造家 第一季,梦想改造家 第一季 Season 1,梦想改造家 第一季影评,剧情介绍,图片,论坛">
    <meta name="description" content="梦想改造家 第一季电视剧简介和剧情介绍,梦想改造家 第一季影评、图片、论坛">
    <meta name="mobile-agent" content="format=html5; url=http://m.douban.com/movie/subject/25942176/"/>
    <link rel="alternate" href="android-app://com.douban.movie/doubanmovie/subject/25942176/" />
    <link rel="stylesheet" href="https://img3.doubanio.com/dae/cdnlib/libs/LikeButton/1.0.5/style.min.css">
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/77323ae72a612bba8b65f845491513ff3329b1bb/js/do.js" data-cfg-autoload="false"></script>
    <script type="text/javascript">
      Do.add('dialog', {path: 'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js', type: 'js'});
      Do.add('dialog-css', {path: 'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css', type: 'css'});
      Do.add('handlebarsjs', {path: 'https://img3.doubanio.com/f/movie/3d4f8e4a8918718256450eb6e57ec8e1f7a2e14b/js/movie/lib/handlebars.current.js', type: 'js'});
    </script>
    
  <script type='text/javascript'>
  var _vwo_code = (function() {
    var account_id = 249272,
      settings_tolerance = 0,
      library_tolerance = 2500,
      use_existing_jquery = false,
      // DO NOT EDIT BELOW THIS LINE
      f=false,d=document;return{use_existing_jquery:function(){return use_existing_jquery;},library_tolerance:function(){return library_tolerance;},finish:function(){if(!f){f=true;var a=d.getElementById('_vis_opt_path_hides');if(a)a.parentNode.removeChild(a);}},finished:function(){return f;},load:function(a){var b=d.createElement('script');b.src=a;b.type='text/javascript';b.innerText;b.onerror=function(){_vwo_code.finish();};d.getElementsByTagName('head')[0].appendChild(b);},init:function(){settings_timer=setTimeout('_vwo_code.finish()',settings_tolerance);var a=d.createElement('style'),b='body{opacity:0 !important;filter:alpha(opacity=0) !important;background:none !important;}',h=d.getElementsByTagName('head')[0];a.setAttribute('id','_vis_opt_path_hides');a.setAttribute('type','text/css');if(a.styleSheet)a.styleSheet.cssText=b;else a.appendChild(d.createTextNode(b));h.appendChild(a);this.load('//dev.visualwebsiteoptimizer.com/j.php?a='+account_id+'&u='+encodeURIComponent(d.URL)+'&r='+Math.random());return settings_timer;}};}());

  +function () {
    var bindEvent = function (el, type, handler) {
        var $ = window.jQuery || window.Zepto || window.$
       if ($ && $.fn && $.fn.on) {
           $(el).on(type, handler)
       } else if($ && $.fn && $.fn.bind) {
           $(el).bind(type, handler)
       } else if (el.addEventListener){
         el.addEventListener(type, handler, false);
       } else if (el.attachEvent){
         el.attachEvent("on" + type, handler);
       } else {
         el["on" + type] = handler;
       }
     }

    var _origin_load = _vwo_code.load
    _vwo_code.load = function () {
      var args = [].slice.call(arguments)
      bindEvent(window, 'load', function () {
        _origin_load.apply(_vwo_code, args)
      })
    }
  }()

  _vwo_settings_timer = _vwo_code.init();
  </script>


    


<script type="application/ld+json">
{
  "@context": "http://schema.org",
  "name": "梦想改造家 第一季",
  "url": "/subject/25942176/",
  "image": "https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2261392790.webp",
  "director": 
  [
  ]
,
  "author": 
  [
  ]
,
  "actor": 
  [
    {
      "@type": "Person",
      "url": "/celebrity/1342671/",
      "name": "施琰 Yan Shi"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1342429/",
      "name": "骆新 Xin Luo"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1322341/",
      "name": "金星 Xing Jin"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1275430/",
      "name": "戴娇倩 Shirley Dai"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1276103/",
      "name": "佟瑞欣 Ruixin Tong"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1329801/",
      "name": "黄豆豆 Doudou Huang"
    }
    ,
    {
      "@type": "Person",
      "url": "/celebrity/1001393/",
      "name": "何赛飞 Saifei He"
    }
  ]
,
  "datePublished": "2014-07-30",
  "genre": ["\u771f\u4eba\u79c0"],
  "duration": "",
  "description": "《梦想改造家》是东方卫视打造的一档大型装修真人秀节目。节目在全国范围内，遴选有居住困难的普通家庭，在限定费用中，通过设计师的匠心巧思，完成看似不可能完成的家装梦想。
节目大胆突破了以色彩和软装为主的传...",
  "@type": "TVSeries",
  "aggregateRating": {
    "@type": "AggregateRating",
    "ratingCount": "3618",
    "bestRating": "10",
    "worstRating": "2",
    "ratingValue": "9.0"
  }
}
</script>


    <style type="text/css">
  
</style>
    <style type="text/css">img { max-width: 100%; }</style>
    <script type="text/javascript"></script>
    <link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/4b01fb68e0ca359b.css">

    <link rel="shortcut icon" href="https://img3.doubanio.com/favicon.ico" type="image/x-icon">
</head>

<body>
  
    <script type="text/javascript">var _body_start = new Date();</script>

    
    



    <link href="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.css" rel="stylesheet" type="text/css">



<div id="db-global-nav" class="global-nav">
  <div class="bd">
    
<div class="top-nav-info">
  <a href="https://www.douban.com/accounts/login?source=movie" class="nav-login" rel="nofollow">登录</a>
  <a href="https://www.douban.com/accounts/register?source=movie" class="nav-register" rel="nofollow">注册</a>
</div>


    <div class="top-nav-doubanapp">
  <a href="https://www.douban.com/doubanapp/app?channel=top-nav" class="lnk-doubanapp">下载豆瓣客户端</a>
  <div id="doubanapp-tip">
    <a href="https://www.douban.com/doubanapp/app?channel=qipao" class="tip-link">豆瓣 <span class="version">6.0</span> 全新发布</a>
    <a href="javascript: void 0;" class="tip-close">×</a>
  </div>
  <div id="top-nav-appintro" class="more-items">
    <p class="appintro-title">豆瓣</p>
    <p class="qrcode">扫码直接下载</p>
    <div class="download">
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=iOS">iPhone</a>
      <span>·</span>
      <a href="https://www.douban.com/doubanapp/redirect?channel=top-nav&direct_dl=1&download=Android" class="download-android">Android</a>
    </div>
  </div>
</div>

    


<div class="global-nav-items">
  <ul>
    <li class="">
      <a href="https://www.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-main&quot;,&quot;uid&quot;:&quot;0&quot;}">豆瓣</a>
    </li>
    <li class="">
      <a href="https://book.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-book&quot;,&quot;uid&quot;:&quot;0&quot;}">读书</a>
    </li>
    <li class="on">
      <a href="https://movie.douban.com"  data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-movie&quot;,&quot;uid&quot;:&quot;0&quot;}">电影</a>
    </li>
    <li class="">
      <a href="https://music.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-music&quot;,&quot;uid&quot;:&quot;0&quot;}">音乐</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/location" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-location&quot;,&quot;uid&quot;:&quot;0&quot;}">同城</a>
    </li>
    <li class="">
      <a href="https://www.douban.com/group" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-group&quot;,&quot;uid&quot;:&quot;0&quot;}">小组</a>
    </li>
    <li class="">
      <a href="https://read.douban.com&#47;?dcs=top-nav&amp;dcm=douban" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-read&quot;,&quot;uid&quot;:&quot;0&quot;}">阅读</a>
    </li>
    <li class="">
      <a href="https://douban.fm&#47;?from_=shire_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-fm&quot;,&quot;uid&quot;:&quot;0&quot;}">FM</a>
    </li>
    <li class="">
      <a href="https://time.douban.com&#47;?dt_time_source=douban-web_top_nav" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-time&quot;,&quot;uid&quot;:&quot;0&quot;}">时间</a>
    </li>
    <li class="">
      <a href="https://market.douban.com&#47;?utm_campaign=douban_top_nav&amp;utm_source=douban&amp;utm_medium=pc_web" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-market&quot;,&quot;uid&quot;:&quot;0&quot;}">豆品</a>
    </li>
    <li>
      <a href="#more" class="bn-more"><span>更多</span></a>
      <div class="more-items">
        <table cellpadding="0" cellspacing="0">
          <tbody>
            <tr>
              <td>
                <a href="https://ypy.douban.com" target="_blank" data-moreurl-dict="{&quot;from&quot;:&quot;top-nav-click-ypy&quot;,&quot;uid&quot;:&quot;0&quot;}">豆瓣摄影</a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </li>
  </ul>
</div>

  </div>
</div>
<script>
  ;window._GLOBAL_NAV = {
    DOUBAN_URL: "https://www.douban.com",
    N_NEW_NOTIS: 0,
    N_NEW_DOUMAIL: 0
  };
</script>



    <script src="//img3.doubanio.com/dae/accounts/resources/984c231/shire/bundle.js" defer="defer"></script>




    



    <link href="//img3.doubanio.com/dae/accounts/resources/8c80301/movie/bundle.css" rel="stylesheet" type="text/css">




<div id="db-nav-movie" class="nav">
  <div class="nav-wrap">
  <div class="nav-primary">
    <div class="nav-logo">
      <a href="https:&#47;&#47;movie.douban.com">豆瓣电影</a>
    </div>
    <div class="nav-search">
      <form action="https:&#47;&#47;movie.douban.com/subject_search" method="get">
        <fieldset>
          <legend>搜索：</legend>
          <label for="inp-query">
          </label>
          <div class="inp"><input id="inp-query" name="search_text" size="22" maxlength="60" placeholder="搜索电影、电视剧、综艺、影人" value=""></div>
          <div class="inp-btn"><input type="submit" value="搜索"></div>
          <input type="hidden" name="cat" value="1002" />
        </fieldset>
      </form>
    </div>
  </div>
  </div>
  <div class="nav-secondary">
    

<div class="nav-items">
  <ul>
    <li    ><a href="https://movie.douban.com/cinema/nowplaying/"
     >影讯&购票</a>
    </li>
    <li    ><a href="https://movie.douban.com/explore"
     >选电影</a>
    </li>
    <li    ><a href="https://movie.douban.com/tv/"
     >电视剧</a>
    </li>
    <li    ><a href="https://movie.douban.com/chart"
     >排行榜</a>
    </li>
    <li    ><a href="https://movie.douban.com/tag/"
     >分类</a>
    </li>
    <li    ><a href="https://movie.douban.com/review/best/"
     >影评</a>
    </li>
    <li    ><a href="https://movie.douban.com/annual/2018?source=navigation"
     >2018年度榜单</a>
    </li>
    <li    ><a href="https://www.douban.com/standbyme/2018?source=navigation"
     >2018书影音报告</a>
    </li>
  </ul>
</div>

    <a href="https://movie.douban.com/annual/2018?source=movie_navigation" class="movieannual2018"></a>
  </div>
</div>

<script id="suggResult" type="text/x-jquery-tmpl">
  <li data-link="{{= url}}">
            <a href="{{= url}}" onclick="moreurl(this, {from:'movie_search_sugg', query:'{{= keyword }}', subject_id:'{{= id}}', i: '{{= index}}', type: '{{= type}}'})">
            <img src="{{= img}}" width="40" />
            <p>
                <em>{{= title}}</em>
                {{if year}}
                    <span>{{= year}}</span>
                {{/if}}
                {{if sub_title}}
                    <br /><span>{{= sub_title}}</span>
                {{/if}}
                {{if address}}
                    <br /><span>{{= address}}</span>
                {{/if}}
                {{if episode}}
                    {{if episode=="unknow"}}
                        <br /><span>集数未知</span>
                    {{else}}
                        <br /><span>共{{= episode}}集</span>
                    {{/if}}
                {{/if}}
            </p>
        </a>
        </li>
  </script>




    <script src="//img3.doubanio.com/dae/accounts/resources/8c80301/movie/bundle.js" defer="defer"></script>





    
    <div id="wrapper">
        

        
    <div id="content">
        

    <div id="dale_movie_subject_top_icon"></div>
    <h1>
        <span property="v:itemreviewed">梦想改造家 第一季</span>
            <span class="year">(2014)</span>
    </h1>

        <div class="grid-16-8 clearfix">
            

            
            <div class="article">
                
    

    





        <div class="indent clearfix">
            <div class="subjectwrap clearfix">
                <div class="subject clearfix">
                    



<div id="mainpic" class="">
    <a class="nbgnbg" href="https://movie.douban.com/subject/25942176/photos?type=R" title="点击看更多海报">
        <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2261392790.webp" title="点击看更多海报" alt="梦想改造家 第一季" rel="v:image" />
   </a>
</div>

                    


<div id="info">
        
        
        <span class="actor"><span class='pl'>主演</span>: <span class='attrs'><a href="/celebrity/1342671/" rel="v:starring">施琰</a> / <a href="/celebrity/1342429/" rel="v:starring">骆新</a></span></span><br/>
        <span class="pl">类型:</span> <span property="v:genre">真人秀</span><br/>
        
        <span class="pl">制片国家/地区:</span> 中国大陆<br/>
        <span class="pl">语言:</span> 汉语普通话<br/>
        <span class="pl">首播:</span> <span property="v:initialReleaseDate" content="2014-07-30(中国大陆)">2014-07-30(中国大陆)</span><br/>
        <span class="pl">季数:</span> 1<br/>
        <span class="pl">集数:</span> 13<br/>
        
        
        

</div>




                </div>
                    


<div id="interest_sectl">
    <div class="rating_wrap clearbox" rel="v:rating">
        <div class="clearfix">
          <div class="rating_logo ll">豆瓣评分</div>
          <div class="output-btn-wrap rr" style="display:none">
            <img src="https://img3.doubanio.com/f/movie/692e86756648f29457847c5cc5e161d6f6b8aaac/pics/movie/reference.png" />
            <a class="download-output-image" href="#">引用</a>
          </div>
          
          
        </div>
        


<div class="rating_self clearfix" typeof="v:Rating">
    <strong class="ll rating_num" property="v:average">9.0</strong>
    <span property="v:best" content="10.0"></span>
    <div class="rating_right ">
        <div class="ll bigstar bigstar45"></div>
        <div class="rating_sum">
                <a href="collections" class="rating_people"><span property="v:votes">3618</span>人评价</a>
        </div>
    </div>
</div>
<div class="ratings-on-weight">
    
        <div class="item">
        
        <span class="stars5 starstop" title="力荐">
            5星
        </span>
        <div class="power" style="width:64px"></div>
        <span class="rating_per">60.2%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars4 starstop" title="推荐">
            4星
        </span>
        <div class="power" style="width:33px"></div>
        <span class="rating_per">31.5%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars3 starstop" title="还行">
            3星
        </span>
        <div class="power" style="width:8px"></div>
        <span class="rating_per">7.5%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars2 starstop" title="较差">
            2星
        </span>
        <div class="power" style="width:0px"></div>
        <span class="rating_per">0.6%</span>
        <br />
        </div>
        <div class="item">
        
        <span class="stars1 starstop" title="很差">
            1星
        </span>
        <div class="power" style="width:0px"></div>
        <span class="rating_per">0.2%</span>
        <br />
        </div>
</div>

    </div>
</div>


                
            </div>
                




<div id="interest_sect_level" class="clearfix">
        
            <a href="https://www.douban.com/reason=collectwish&amp;ck=" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-25942176-wish">
                <span>想看</span>
            </a>
            <a href="https://www.douban.com/reason=collectdo&amp;ck=" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-25942176-do">
                <span>在看</span>
            </a>
            <a href="https://www.douban.com/reason=collectcollect&amp;ck=" rel="nofollow" class="j a_show_login colbutt ll" name="pbtn-25942176-collect">
                <span>看过</span>
            </a>
        <div class="ll j a_stars">
            
    
    评价:
    <span id="rating"> <span id="stars" data-solid="https://img3.doubanio.com/f/shire/5a2327c04c0c231bced131ddf3f4467eb80c1c86/pics/rating_icons/star_onmouseover.png" data-hollow="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" data-solid-2x="https://img3.doubanio.com/f/shire/7258904022439076d57303c3b06ad195bf1dc41a/pics/rating_icons/star_onmouseover@2x.png" data-hollow-2x="https://img3.doubanio.com/f/shire/95cc2fa733221bb8edd28ad56a7145a5ad33383e/pics/rating_icons/star_hollow_hover@2x.png">

            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-25942176-1">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star1" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-25942176-2">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star2" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-25942176-3">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star3" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-25942176-4">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star4" width="16" height="16"/></a>
            <a href="https://www.douban.com/register?reason=rate" class="j a_show_login" name="pbtn-25942176-5">
        <img src="https://img3.doubanio.com/f/shire/2520c01967207a1735171056ec588c8c1257e5f8/pics/rating_icons/star_hollow_hover.png" id="star5" width="16" height="16"/></a>
    </span><span id="rateword" class="pl"></span>
    <input id="n_rating" type="hidden" value=""  />
    </span>

        </div>
    

</div>


            


















<div class="gtleft">
    <ul class="ul_subject_menu bicelink color_gray pt6 clearfix">
        
    
        
                <li> 
    <img src="https://img3.doubanio.com/f/shire/cc03d0fcf32b7ce3af7b160a0b85e5e66b47cc42/pics/short-comment.gif" />&nbsp;
        <a onclick="moreurl(this, {from:'mv_sbj_wr_cmnt_login'})" class="j a_show_login" href="https://www.douban.com/register?reason=review" rel="nofollow">写短评</a>
 </li>
                    <li> 
    
    <img src="https://img3.doubanio.com/f/shire/5bbf02b7b5ec12b23e214a580b6f9e481108488c/pics/add-review.gif" />&nbsp;
        <a onclick="moreurl(this, {from:'mv_sbj_wr_rv_login'})" class="j a_show_login" href="https://www.douban.com/register?reason=review" rel="nofollow">写影评</a>
 </li>
                <li> 
    



 </li>
                <li> 
   

   
    
    <span class="rec" id="电视剧-25942176">
    <a href= "#"
        data-type="电视剧"
        data-url="https://movie.douban.com/subject/25942176/"
        data-desc="电视剧《梦想改造家 第一季》 (来自豆瓣) "
        data-title="电视剧《梦想改造家 第一季》 (来自豆瓣) "
        data-pic="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2261392790.jpeg"
        class="bn-sharing ">
        分享到
    </a> &nbsp;&nbsp;
    </span>

    <script>
        if (!window.DoubanShareMenuList) {
            window.DoubanShareMenuList = [];
        }
        var __cache_url = __cache_url || {};

        (function(u){
            if(__cache_url[u]) return;
            __cache_url[u] = true;
            window.DoubanShareIcons = 'https://img3.doubanio.com/f/shire/d15ffd71f3f10a7210448fec5a68eaec66e7f7d0/pics/ic_shares.png';

            var initShareButton = function() {
                $.ajax({url:u,dataType:'script',cache:true});
            };

            if (typeof Do == 'function' && 'ready' in Do) {
                Do(
                    'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
                    'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
                    'https://img3.doubanio.com/f/movie/c4ab132ff4d3d64a83854c875ea79b8b541faf12/js/movie/lib/qrcode.min.js',
                    initShareButton
                );
            } else if(typeof Douban == 'object' && 'loader' in Douban) {
                Douban.loader.batch(
                    'https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css',
                    'https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js',
                    'https://img3.doubanio.com/f/movie/c4ab132ff4d3d64a83854c875ea79b8b541faf12/js/movie/lib/qrcode.min.js'
                ).done(initShareButton);
            }

        })('https://img3.doubanio.com/f/movie/32be6727ed3ad8f6c4a417d8a086355c3e7d1d27/js/movie/lib/sharebutton.js');
    </script>


  </li>
            

    </ul>

    <script type="text/javascript">
        $(function(){
            $(".ul_subject_menu li.rec .bn-sharing").bind("click", function(){
                $.get("/blank?sbj_page_click=bn_sharing");
            });
            $(".ul_subject_menu .create_from_menu").bind("click", function(e){
                e.preventDefault();
                var $el = $(this);
                var glRoot = document.getElementById('gallery-topics-selection');
                if (window.has_gallery_topics && glRoot) {
                    // 判断是否有话题
                    glRoot.style.display = 'block';
                } else {
                    location.href = $el.attr('href');
                }
            });
        });
    </script>
</div>




                





<div class="rec-sec">
<span class="rec">
    <script id="movie-share" type="text/x-html-snippet">
        
    <form class="movie-share" action="/j/share" method="POST">
        <div class="clearfix form-bd">
            <div class="input-area">
                <textarea name="text" class="share-text" cols="72" data-mention-api="https://api.douban.com/shuo/in/complete?alt=xd&amp;callback=?"></textarea>
                <input type="hidden" name="target-id" value="25942176">
                <input type="hidden" name="target-type" value="0">
                <input type="hidden" name="title" value="梦想改造家 第一季‎ (2014)">
                <input type="hidden" name="desc" value=" 主演 施琰 / 骆新 / 中国大陆 / 9.0分(3618评价)">
                <input type="hidden" name="redir" value=""/>
                <div class="mentioned-highlighter"></div>
            </div>

            <div class="info-area">
                    <img class="media" src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2261392790.webp" />
                <strong>梦想改造家 第一季‎ (2014)</strong>
                <p> 主演 施琰 / 骆新 / 中国大陆 / 9.0分(3618评价)</p>
                <p class="error server-error">&nbsp;</p>
            </div>
        </div>
        <div class="form-ft">
            <div class="form-ft-inner">
                



                <span class="avail-num-indicator">140</span>
                <span class="bn-flat">
                    <input type="submit" value="推荐" />
                </span>
            </div>
        </div>
    </form>
    
    <div id="suggest-mention-tmpl" style="display:none;">
        <ul>
            {{#users}}
            <li id="{{uid}}">
              <img src="{{avatar}}">{{{username}}}&nbsp;<span>({{{uid}}})</span>
            </li>
            {{/users}}
        </ul>
    </div>


    </script>

        
        <a href="/accounts/register?reason=recommend"  class="j a_show_login lnk-sharing" share-id="25942176" data-mode="plain" data-name="梦想改造家 第一季‎ (2014)" data-type="movie" data-desc=" 主演 施琰 / 骆新 / 中国大陆 / 9.0分(3618评价)" data-href="https://movie.douban.com/subject/25942176/" data-image="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2261392790.webp" data-properties="{}" data-redir="" data-text="" data-apikey="" data-curl="" data-count="10" data-object_kind="1002" data-object_id="25942176" data-target_type="rec" data-target_action="1" data-action_props="{&#34;subject_url&#34;:&#34;https:\/\/movie.douban.com\/subject\/25942176\/&#34;,&#34;subject_title&#34;:&#34;梦想改造家 第一季‎ (2014)&#34;}">推荐</a>
</span>


</div>






            <script type="text/javascript">
                $(function() {
                    $('.collect_btn', '#interest_sect_level').each(function() {
                        Douban.init_collect_btn(this);
                    });
                    $('html').delegate(".indent .rec-sec .lnk-sharing", "click", function() {
                        moreurl(this, {
                            from : 'mv_sbj_db_share'
                        });
                    });
                });
            </script>
        </div>
            


    <div id="collect_form_25942176"></div>

        







    <h2>
        <i class="">梦想改造家 第一季的分集短评</i>
              · · · · · ·
    </h2>

    


    
    <div class="episode_list">


            

            <a class=" item" href="/subject/25942176/episode/1/">1集</a>
            

            <a class=" item" href="/subject/25942176/episode/2/">2集</a>
            

            <a class=" item" href="/subject/25942176/episode/3/">3集</a>
            

            <a class=" item" href="/subject/25942176/episode/4/">4集</a>
            

            <a class=" item" href="/subject/25942176/episode/5/">5集</a>
            

            <a class=" item" href="/subject/25942176/episode/6/">6集</a>
            

            <a class=" item" href="/subject/25942176/episode/7/">7集</a>
            

            <a class=" item" href="/subject/25942176/episode/8/">8集</a>
            

            <a class=" item" href="/subject/25942176/episode/9/">9集</a>
            

            <a class=" item" href="/subject/25942176/episode/10/">10集</a>
            

            <a class=" item" href="/subject/25942176/episode/11/">11集</a>
            

            <a class=" item" href="/subject/25942176/episode/12/">12集</a>
            

            <a class=" item" href="/subject/25942176/episode/13/">13集</a>


    </div>




        



<div class="related-info" style="margin-bottom:-10px;">
    <a name="intro"></a>
    
        
            
            
    <h2>
        <i class="">梦想改造家 第一季的剧情简介</i>
              · · · · · ·
    </h2>

            <div class="indent" id="link-report">
                    
                        <span property="v:summary" class="">
                                　　《梦想改造家》是东方卫视打造的一档大型装修真人秀节目。节目在全国范围内，遴选有居住困难的普通家庭，在限定费用中，通过设计师的匠心巧思，完成看似不可能完成的家装梦想。
                                    <br />
                                　　节目大胆突破了以色彩和软装为主的传统装修形式，革命性的聚焦在空间改变和功能实现上，狭小的空间，奇葩的房型，看似无法解决的居住困难，都由顶尖设计师在有限的资金，有限的空间里根据委托人的特殊需求进行彻底的改造，真正解决委托家庭的住房难题。它通过聚焦全国范围内不同地域特色的建筑物，不同类型和背景的家庭故事，揭示家给予人的意义，见证家装改造给予人的幸福。
                        </span>
                        

            </div>
</div>


    








<div id="celebrities" class="celebrities related-celebrities">

  
    <h2>
        <i class="">梦想改造家 第一季的演职员</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="/subject/25942176/celebrities">全部 7</a>
            )
            </span>
    </h2>


  <ul class="celebrities-list from-subject __oneline">
        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1342671/" title="施琰 Yan Shi" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1510713635.43.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1342671/" title="施琰 Yan Shi" class="name">施琰</a></span>

      <span class="role" title="自己">自己</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1342429/" title="骆新 Xin Luo" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1408720453.34.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1342429/" title="骆新 Xin Luo" class="name">骆新</a></span>

      <span class="role" title="自己">自己</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1322341/" title="金星 Xing Jin" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1432001903.74.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1322341/" title="金星 Xing Jin" class="name">金星</a></span>

      <span class="role" title="自己">自己</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1275430/" title="戴娇倩 Shirley Dai" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p13632.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1275430/" title="戴娇倩 Shirley Dai" class="name">戴娇倩</a></span>

      <span class="role" title="自己">自己</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1276103/" title="佟瑞欣 Ruixin Tong" class="">
      <div class="avatar" style="background-image: url(https://img1.doubanio.com/view/celebrity/s_ratio_celebrity/public/p29367.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1276103/" title="佟瑞欣 Ruixin Tong" class="name">佟瑞欣</a></span>

      <span class="role" title="自己">自己</span>

    </div>
  </li>


        
    

    
  
  <li class="celebrity">
    

  <a href="https://movie.douban.com/celebrity/1329801/" title="黄豆豆 Doudou Huang" class="">
      <div class="avatar" style="background-image: url(https://img3.doubanio.com/view/celebrity/s_ratio_celebrity/public/p1370401643.61.webp)">
    </div>
  </a>

    <div class="info">
      <span class="name"><a href="https://movie.douban.com/celebrity/1329801/" title="黄豆豆 Doudou Huang" class="name">黄豆豆</a></span>

      <span class="role" title="自己">自己</span>

    </div>
  </li>


  </ul>
</div>


    


<link rel="stylesheet" href="https://img3.doubanio.com/f/verify/16c7e943aee3b1dc6d65f600fcc0f6d62db7dfb4/entry_creator/dist/author_subject/style.css">
<div id="author_subject" class="author-wrapper">
    <div class="loading"></div>
</div>
<script type="text/javascript">
    var answerObj = {
      ISALL: 'False',
      TYPE: 'tv',
      SUBJECT_ID: '25942176',
      USER_ID: 'None'
    }
</script>
<script type="text/javascript" src="https://img3.doubanio.com/f/movie/61252f2f9b35f08b37f69d17dfe48310dd295347/js/movie/lib/react/15.4/bundle.js"></script>
<script type="text/javascript" src="https://img3.doubanio.com/f/verify/ac140ef86262b845d2be7b859e352d8196f3f6d4/entry_creator/dist/author_subject/index.js"></script>


    









    
    <div id="related-pic" class="related-pic">
        
    
    
    <h2>
        <i class="">梦想改造家 第一季的图片</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="/video/create?subject_id=25942176">添加视频评论</a>&nbsp;|&nbsp;<a href="https://movie.douban.com/subject/25942176/all_photos">图片160</a>&nbsp;·&nbsp;<a href="https://movie.douban.com/subject/25942176/mupload">添加</a>
            )
            </span>
    </h2>


        <ul class="related-pic-bd  ">
                <li>
                    <a href="https://movie.douban.com/photos/photo/2228138376/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p2228138376.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2271081404/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p2271081404.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2265216429/"><img src="https://img1.doubanio.com/view/photo/sqxs/public/p2265216429.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2265217330/"><img src="https://img3.doubanio.com/view/photo/sqxs/public/p2265217330.webp" alt="图片" /></a>
                </li>
                <li>
                    <a href="https://movie.douban.com/photos/photo/2265217329/"><img src="https://img1.doubanio.com/view/photo/sqxs/public/p2265217329.webp" alt="图片" /></a>
                </li>
        </ul>
    </div>



      








<div class="mod">
<div class="hd-ops">
  
  <a class="comment_btn j a_show_login" href="https://www.douban.com/register?reason=discussion" rel="nofollow">
      <span>发起新的讨论</span>
  </a>

</div>

    <h2>
        <i class="">讨论区</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/25942176/discussion/">全部</a>
            )
            </span>
    </h2>

<div class="bd">
<div class="mv-discussion-nav">
<a href="https://movie.douban.com/subject/25942176/discussion/" class="on">最新</a>
<a href="https://movie.douban.com/subject/25942176/discussion/?sort=vote" data-epid="hot">热门</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=1" data-epid="821658" data-num="1">1集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=2" data-epid="821659" data-num="2">2集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=3" data-epid="821660" data-num="3">3集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=4" data-epid="821661" data-num="4">4集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=5" data-epid="821662" data-num="5">5集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=6" data-epid="821663" data-num="6">6集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=7" data-epid="821664" data-num="7">7集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/" data-epid="more" title="更多">&#8230;</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=8" data-epid="821665" data-num="8" class="more-item">8集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=9" data-epid="821666" data-num="9" class="more-item">9集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=10" data-epid="821667" data-num="10" class="more-item">10集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=11" data-epid="821668" data-num="11" class="more-item">11集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=12" data-epid="821669" data-num="12" class="more-item">12集</a>
  <a href="https://movie.douban.com/subject/25942176/discussion/?ep_num=13" data-epid="821670" data-num="13" class="more-item">13集</a>
</div>

<div class="mv-discussion-list discussion-list">
  

<table>
  <thead>
  <tr>
    <td>讨论</td>
    <td>作者</td>
    <td nowrap="nowrap">回应</td>
    <td align="right">最后回应</td>
  </tr>
  </thead>
  <tbody>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/614679321/" title="[第9集] 王仲平果然是温暖小天使">[第9集] 王仲平果然是温暖小天使</a>
    </td>
    <td><a href="https://www.douban.com/people/74728534/">加载中请稍后</a></td>
    <td class="reply-num">1</td>
    <td class="time">2018-10-07 11:05</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/615122650/" title="上海小夫妻两个价格控制要10万的">上海小夫妻两个价格控制要10万的</a>
    </td>
    <td><a href="https://www.douban.com/people/wangwanxu/">JILL</a></td>
    <td class="reply-num"></td>
    <td class="time">2017-11-30 19:28</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/614679306/" title="[第11集] 平淡而实用的设计">[第11集] 平淡而实用的设计</a>
    </td>
    <td><a href="https://www.douban.com/people/74728534/">加载中请稍后</a></td>
    <td class="reply-num"></td>
    <td class="time">2017-03-02 20:49</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/58942742/" title="[第4集] 武汉这期很厉害">[第4集] 武汉这期很厉害</a>
    </td>
    <td><a href="https://www.douban.com/people/40405464/">Near</a></td>
    <td class="reply-num">2</td>
    <td class="time">2017-01-27 11:03</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/612270924/" title="[第2集] 和我猜的价钱差不多，我猜大概要20万，没想到居然是19万5！">[第2集] 和我猜的价钱差不多，我猜大概要20万，没...</a>
    </td>
    <td><a href="https://www.douban.com/people/43693739/">Deluge Again</a></td>
    <td class="reply-num"></td>
    <td class="time">2015-08-12 01:02</td>
  </tr>
  </tbody>
</table>

<a href="https://movie.douban.com/subject/25942176/discussion/">&gt; 全部讨论5条</a>
</div>

<div class="mv-hot-discussion-list hide">
  

<table>
  <thead>
  <tr>
    <td>讨论</td>
    <td>作者</td>
    <td nowrap="nowrap">回应</td>
    <td align="right">最后回应</td>
  </tr>
  </thead>
  <tbody>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/58942742/" title="[第4集] 武汉这期很厉害">[第4集] 武汉这期很厉害</a>
    </td>
    <td><a href="https://www.douban.com/people/40405464/">Near</a></td>
    <td class="reply-num">2</td>
    <td class="time">2017-01-27 11:03</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/614679321/" title="[第9集] 王仲平果然是温暖小天使">[第9集] 王仲平果然是温暖小天使</a>
    </td>
    <td><a href="https://www.douban.com/people/74728534/">加载中请稍后</a></td>
    <td class="reply-num">1</td>
    <td class="time">2018-10-07 11:05</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/615122650/" title="上海小夫妻两个价格控制要10万的">上海小夫妻两个价格控制要10万的</a>
    </td>
    <td><a href="https://www.douban.com/people/wangwanxu/">JILL</a></td>
    <td class="reply-num"></td>
    <td class="time">2017-11-30 19:28</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/614679306/" title="[第11集] 平淡而实用的设计">[第11集] 平淡而实用的设计</a>
    </td>
    <td><a href="https://www.douban.com/people/74728534/">加载中请稍后</a></td>
    <td class="reply-num"></td>
    <td class="time">2017-03-02 20:49</td>
  </tr>
  
  <tr>
    <td class="title">
      <a href="https://movie.douban.com/subject/25942176/discussion/612270924/" title="[第2集] 和我猜的价钱差不多，我猜大概要20万，没想到居然是19万5！">[第2集] 和我猜的价钱差不多，我猜大概要20万，没...</a>
    </td>
    <td><a href="https://www.douban.com/people/43693739/">Deluge Again</a></td>
    <td class="reply-num"></td>
    <td class="time">2015-08-12 01:02</td>
  </tr>
  </tbody>
</table>

<a href="https://movie.douban.com/subject/25942176/discussion/?sort=vote">&gt; 全部讨论5条</a>
</div>

</div>
</div>






    
    



<style type="text/css">
.award li { display: inline; margin-right: 5px }
.awards { margin-bottom: 20px }
.awards h2 { background: none; color: #000; font-size: 14px; padding-bottom: 5px; margin-bottom: 8px; border-bottom: 1px dashed #dddddd }
.awards .year { color: #666666; margin-left: -5px }
.mod { margin-bottom: 25px }
.mod .hd { margin-bottom: 10px }
.mod .hd h2 {margin:24px 0 3px 0}
</style>



    








    <div id="recommendations" class="">
        
        
    <h2>
        <i class="">喜欢这部剧集的人也喜欢</i>
              · · · · · ·
    </h2>

        
    
    <div class="recommendations-bd">
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26576624/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2261392587.webp" alt="梦想改造家 第二季" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26576624/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>梦想改造家 第二季</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/25814941/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2187866062.webp" alt="全能住宅改造王 第二季" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/25814941/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>全能住宅改造王 第二季</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26292751/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2247455571.webp" alt="夕阳红" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26292751/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>夕阳红</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26613426/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2267246025.webp" alt="七巧板" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26613426/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>七巧板</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/25949779/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2227958120.webp" alt="正大综艺" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/25949779/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>正大综艺</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26611670/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2395173877.webp" alt="曲苑杂坛" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26611670/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>曲苑杂坛</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26611668/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2416459304.webp" alt="综艺大观" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26611668/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>综艺大观</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26177736/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2221154684.webp" alt="奇葩说 第一季" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26177736/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>奇葩说 第一季</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26946548/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2409449068.webp" alt="WOW新家" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26946548/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>WOW新家</a>
            </dd>
        </dl>
        <dl class="">
            <dt>
                <a href="https://movie.douban.com/subject/26387728/?from=subject-page" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>
                    <img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2264809475.webp" alt="极限挑战 第一季" class="" />
                </a>
            </dt>
            <dd>
                <a href="https://movie.douban.com/subject/26387728/?from=subject-page" class="" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-recommended-subject&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>极限挑战 第一季</a>
            </dd>
        </dl>
    </div>

    </div>



        


<script type="text/x-handlebar-tmpl" id="comment-tmpl">
    <div class="dummy-fold">
        {{#each comments}}
        <div class="comment-item" data-cid="id">
            <div class="comment">
                <h3>
                    <span class="comment-vote">
                            <span class="votes">{{votes}}</span>
                        <input value="{{id}}" type="hidden"/>
                        <a href="javascript:;" class="j {{#if ../if_logined}}a_vote_comment{{else}}a_show_login{{/if}}">有用</a>
                    </span>
                    <span class="comment-info">
                        <a href="{{user.path}}" class="">{{user.name}}</a>
                        {{#if rating}}
                        <span class="allstar{{rating}}0 rating" title="{{rating_word}}"></span>
                        {{/if}}
                        <span>
                            {{time}}
                        </span>
                        <p> {{content_tmpl content}} </p>
                    </span>
            </div>
        </div>
        {{/each}}
    </div>
</script>












    

    <div id="comments-section">
        <div class="mod-hd">
            
        <a class="comment_btn j a_show_login" href="https://www.douban.com/register?reason=review" rel="nofollow">
            <span>我要写短评</span>
        </a>

            
            
    <h2>
        <i class="">梦想改造家 第一季的短评</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/25942176/comments?status=P">全部 624 条</a>
            )
            </span>
    </h2>

        </div>
        <div class="mod-bd">
                
    <div class="tab-hd">
        <a id="hot-comments-tab" href="comments" data-id="hot" class="on">热门</a>&nbsp;/&nbsp;
        <a id="new-comments-tab" href="comments?sort=time" data-id="new">最新</a>&nbsp;/&nbsp;
        <a id="following-comments-tab" href="follows_comments" data-id="following"  class="j a_show_login">好友</a>
    </div>

    <div class="tab-bd">
        <div id="hot-comments" class="tab">
            
    
        
        <div class="comment-item" data-cid="1065907027">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">37</span>
                <input value="1065907027" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/movie24city/" class="">24city</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2016-07-25 10:13:23">
                    2016-07-25
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">第一季的败笔就是演播室环节，完全割裂了整个节目的风格统一性，而且请的嘉宾和专家也没什么实际的作用，3D模型的讲解不如第二季的电脑3D绘图，全方位展现及讲解，更让人直观理解。大部分都很好，不过这一季出了一个无良设计师廖开民，不止毁了一期节目，一所房子，还毁了一个老人的梦想。</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="967457324">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">8</span>
                <input value="967457324" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/Plutoyw/" class="">Pluto</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2015-10-13 10:17:10">
                    2015-10-13
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">对于空间的利用设计真是达到极致，尊重人性和尊严的意义被提上无比的高度，有几期风格的确趋于酒店化，个性不够突出，西湖边北山路上那套纯白的设计惊艳到我了，嘉宾太败笔，砖家和明星都没啥想法，全靠主持人在那儿蹦跶撑着场面，泪点总是在走过大半个世纪的老人看着儿孙满堂的笑容里</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="945591588">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">4</span>
                <input value="945591588" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/yue525/" class="">玥玥</a>
                    <span>看过</span>
                    <span class="allstar50 rating" title="力荐"></span>
                <span class="comment-time " title="2015-07-30 11:48:27">
                    2015-07-30
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">有点儿日本小屋改造那个节目的味道，虽然都故事化了，但拍得挺有情怀。人活一世，精神层面之下，衣食住行是基础，有个干净明亮温暖的自己的小窝，总是踏实欢喜。有几期邻居不乐意的状态，也特别中国，我们总是要在乱七八糟的环境里，努力让自己过得好。看见别人高兴的样子，也觉得心里真高兴。</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="961891186">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">4</span>
                <input value="961891186" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/scorpio0_0/" class="">Scorpio|千寻</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2015-09-25 15:11:43">
                    2015-09-25
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">模仿日本《全能改造王》很成功。设计师各显神通，小市民人生百态。1，最喜欢的依次是史南桥、沈雷、王平仲、何永明。后面几期比较平淡广告植入也多，最奇葩最恶心的是廖开民。2，通风、采光、灯光、防水。修旧如旧。3，阁楼、钢结构、沙发床、地柜、移门。</span>
        </p>
    </div>

        </div>
        
        <div class="comment-item" data-cid="1114256348">
            
    
    <div class="comment">
        <h3>
            <span class="comment-vote">
                <span class="votes">0</span>
                <input value="1114256348" type="hidden"/>
                <a href="javascript:;" class="j a_show_login" onclick="">有用</a>
            </span>
            <span class="comment-info">
                <a href="https://www.douban.com/people/cindyvvv/" class="">鏡花可可</a>
                    <span>看过</span>
                    <span class="allstar40 rating" title="推荐"></span>
                <span class="comment-time " title="2016-11-27 23:44:09">
                    2016-11-27
                </span>
            </span>
        </h3>
        <p class="">
            
                <span class="short">还需努力</span>
        </p>
    </div>

        </div>



                
                &gt; <a href="comments?sort=new_score&status=P" data-moreurl-dict={&#34;subject_id&#34;:&#34;25942176&#34;,&#34;from&#34;:&#34;tv-more-comments&#34;,&#34;bid&#34;:&#34;auSsK8Dk5cg&#34;}>更多短评624条</a>
        </div>
        <div id="new-comments" class="tab">
            <div id="normal">
            </div>
            <div class="fold-hd hide">
                <a class="qa" href="/help/opinion#t2-q0" target="_blank">为什么被折叠？</a>
                <a class="btn-unfold" href="#">有一些短评被折叠了</a>
                <div class="qa-tip">
                    评论被折叠，是因为发布这条评论的帐号行为异常。评论仍可以被展开阅读，对发布人的账号不造成其他影响。如果认为有问题，可以<a href="https://help.douban.com/help/ask?category=movie">联系</a>豆瓣电影。
                </div>
            </div>
            <div class="fold-bd">
            </div>
            <span id="total-num"></span>
        </div>
        <div id="following-comments" class="tab">
            
    


        <div class="comment-item">
            你关注的人还没写过短评
        </div>

        </div>
    </div>


            
            
        </div>
    </div>



        

<link rel="stylesheet" href="https://img3.doubanio.com/misc/mixed_static/73ed658484f98d44.css">

<section class="topics mod">
    <header>
        <h2>
            梦想改造家 第一季的话题 · · · · · ·
            <span class="pl">( <span class="gallery_topics">全部 <span id="topic-count"></span> 条</span> )</span>
        </h2>
    </header>

    




<section class="subject-topics">
    <div class="topic-guide" id="topic-guide">
        <img class="ic_question" src="//img3.doubanio.com/f/ithildin/b1a3edea3d04805f899e9d77c0bfc0d158df10d5/pics/export/icon_question.png">
        <div class="tip_content">
            <div class="tip_title">什么是话题</div>
            <div class="tip_desc">
                <div>无论是一部作品、一个人，还是一件事，都往往可以衍生出许多不同的话题。将这些话题细分出来，分别进行讨论，会有更多收获。</div>
            </div>
        </div>
        <img class="ic_guide" src="//img3.doubanio.com/f/ithildin/529f46d86bc08f55cd0b1843d0492242ebbd22de/pics/export/icon_guide_arrow.png">
        <img class="ic_close" id="topic-guide-close" src="//img3.doubanio.com/f/ithildin/2eb4ad488cb0854644b23f20b6fa312404429589/pics/export/close@3x.png">
    </div>

    <div id="topic-items"></div>

    <script>
        window.subject_id = 25942176;
        window.join_label_text = '写剧评参与';

        window.topic_display_count = 4;
        window.topic_item_display_count = 1;
        window.no_content_fun_call_name = "no_topic";

        window.guideNode = document.getElementById('topic-guide');
        window.guideNodeClose = document.getElementById('topic-guide-close');
    </script>
    
        <link rel="stylesheet" href="https://img3.doubanio.com/f/ithildin/f731c9ea474da58c516290b3a6b1dd1237c07c5e/css/export/subject_topics.css">
        <script src="https://img3.doubanio.com/f/ithildin/d3590fc6ac47b33c804037a1aa7eec49075428c8/js/export/moment-with-locales-only-zh.js"></script>
        <script src="https://img3.doubanio.com/f/ithildin/c600fdbe69e3ffa5a3919c81ae8c8b4140e99a3e/js/export/subject_topics.js"></script>

</section>

    <script>
        function no_topic(){
            $('#content .topics').remove();
        }
    </script>
</section>

<section class="reviews mod movie-content">
    <header>
        <a href="new_review" rel="nofollow" class="create-review comment_btn"
            data-isverify="False"
            data-verify-url="https://www.douban.com/accounts/phone/verify?redir=http://movie.douban.com/subject/25942176/new_review">
            <span>我要写剧评</span>
        </a>
        <h2>
            梦想改造家 第一季的剧评 · · · · · ·
            <span class="pl">( <a href="reviews">全部 10 条</a> )</span>
        </h2>
    </header>

    

<style>
#gallery-topics-selection {
  position: fixed;
  width: 595px;
  padding: 40px 40px 33px 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 16px 0 rgba(0, 0, 0, 0.2);
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  z-index: 9999;
}
#gallery-topics-selection h1 {
  font-size: 18px;
  color: #007722;
  margin-bottom: 36px;
  padding: 0;
  line-height: 28px;
  font-weight: normal;
}
#gallery-topics-selection .gl_topics {
  border-bottom: 1px solid #dfdfdf;
  max-height: 298px;
  overflow-y: scroll;
}
#gallery-topics-selection .topic {
  margin-bottom: 24px;
}
#gallery-topics-selection .topic_name {
  font-size: 15px;
  color: #333;
  margin: 0;
  line-height: inherit;
}
#gallery-topics-selection .topic_meta {
  font-size: 13px;
  color: #999;
}
#gallery-topics-selection .topics_skip {
  display: block;
  cursor: pointer;
  font-size: 16px;
  color: #3377AA;
  text-align: center;
  margin-top: 33px;
}
#gallery-topics-selection .topics_skip:hover {
  background: transparent;
}
#gallery-topics-selection .close_selection {
  position: absolute;
  width: 30px;
  height: 20px;
  top: 46px;
  right: 40px;
  background: #fff;
  color: #999;
  text-align: right;
}
#gallery-topics-selection .close_selection:hover{
  background: #fff;
  color: #999;
}
</style>




        <div class="review_filter">
            <a href="javascript:;;" class="cur" data-sort="">热门</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="time">最新</a href="javascript:;;"> /
            <a href="javascript:;;" data-sort="follow">好友</a href="javascript:;;">
            
        </div>


        



<div class="review-list  ">
        
    

        
    
    <div data-cid="8422005">
        <div class="main review-item" id="8422005">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/46153548/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u46153548-1.jpg">
        </a>

        <a href="https://www.douban.com/people/46153548/" class="name">土皮</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2017-03-17" class="main-meta">2017-03-17 23:59:05</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8422005/">所谓的风水风水，大意就是如此吧。</a></h2>

                <div id="review_8422005_short" class="review-short" data-rid="8422005">
                    <div class="short-content">

                        房子，家，很多人一辈子就可能困在这只有几十平米的地方生老病死，在广州，上海，北京这些高度城市化的地方，依然有很多原住民生活在城中村那样狭小的地方生存，在这样窘迫的地方，一家几口人，邻里之间的关系也会变的错综复杂，在上海看到的在阳台全裸洗澡的大叔，虽然搞笑，...

                        &nbsp;(<a href="javascript:;" id="toggle-8422005-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8422005_full" class="hidden">
                    <div id="review_8422005_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8422005" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8422005">
                                5
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8422005" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8422005">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8422005/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="8839152">
        <div class="main review-item" id="8839152">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/linyuchen/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u1713332-1.jpg">
        </a>

        <a href="https://www.douban.com/people/linyuchen/" class="name">若焉</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2017-09-30" class="main-meta">2017-09-30 15:26:44</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8839152/">户型没有人奇葩</a></h2>

                <div id="review_8839152_short" class="review-short" data-rid="8839152">
                    <div class="short-content">
                            <p class="spoiler-tip">这篇剧评可能有剧透</p>

                        完全不能想象在二十一世纪的上海和北京，居然还有屋内没有厕所，几代人挤在一起的房子，尤其装修时的人生百态，实是无奈。 第一集：上海黄浦区银行大楼，曾建龙。五十岁还没有结婚的小儿子，不知道从事什么样的工作，总归是啃老啃了一辈子，无法照顾父母，只有大姐打地铺。但依...

                        &nbsp;(<a href="javascript:;" id="toggle-8839152-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8839152_full" class="hidden">
                    <div id="review_8839152_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8839152" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8839152">
                                2
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8839152" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8839152">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8839152/#comments" class="reply">2回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="7527465">
        <div class="main review-item" id="7527465">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/moonhyde/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u1178102-4.jpg">
        </a>

        <a href="https://www.douban.com/people/moonhyde/" class="name">梅有荨</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2015-07-11" class="main-meta">2015-07-11 12:21:11</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/7527465/">不懂得买了日本的版权没有</a></h2>

                <div id="review_7527465_short" class="review-short" data-rid="7527465">
                    <div class="short-content">

                                目前是对日本《全能住宅改造王》还原度最高的一档节目，非常喜欢《全能住宅改造王》，去年发现《梦想改造家》时很开心，每集都看了。         罗里吧嗦的主持人和嘉宾讨论能去掉就最好了。         房间有一些细节的改造，比如楼梯的每级台阶贴上小孩子每一岁的照片，...

                        &nbsp;(<a href="javascript:;" id="toggle-7527465-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_7527465_full" class="hidden">
                    <div id="review_7527465_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="7527465" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-7527465">
                                6
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="7527465" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-7527465">
                                1
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/7527465/#comments" class="reply">9回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="8926654">
        <div class="main review-item" id="8926654">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/54631587/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u54631587-3.jpg">
        </a>

        <a href="https://www.douban.com/people/54631587/" class="name">少先队员周富贵</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2017-11-16" class="main-meta">2017-11-16 10:54:25</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8926654/">备忘</a></h2>

                <div id="review_8926654_short" class="review-short" data-rid="8926654">
                    <div class="short-content">

                        1、上海恒丰大楼，Gary曾建龙，儿子五十岁未婚，老人很恩爱，公共过道厨房全改造。 2、广州西关大屋，何永明，女儿就在广州读书却嫌家破不愿回来住，找个借口说没地方住，父母和奶奶却很朴实，镜子借光。 3、上海城隍庙，史南桥。八十多岁的老奶奶和无力买房儿子和女儿、女婿住...

                        &nbsp;(<a href="javascript:;" id="toggle-8926654-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8926654_full" class="hidden">
                    <div id="review_8926654_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8926654" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8926654">
                                7
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8926654" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8926654">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8926654/#comments" class="reply">2回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="6837529">
        <div class="main review-item" id="6837529">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/76508365/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u76508365-2.jpg">
        </a>

        <a href="https://www.douban.com/people/76508365/" class="name">小荤</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2014-08-21" class="main-meta">2014-08-21 20:31:41</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/6837529/">比老版更好</a></h2>

                <div id="review_6837529_short" class="review-short" data-rid="6837529">
                    <div class="short-content">

                        买的日本 全能改造王的版权，每集都叫……的家  从非常惠生活的《非常梦想家》开始看起，一直觉得不错，新版比老版好在新版的钱是每集主角自己出的，有明细，老版没有提这个问题，我猜是节目组出的。  日本原版的钱就是委托人自己出，我觉得这点很重要，《非常》里面很多家虽然...

                        &nbsp;(<a href="javascript:;" id="toggle-6837529-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_6837529_full" class="hidden">
                    <div id="review_6837529_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="6837529" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-6837529">
                                6
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="6837529" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-6837529">
                                3
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/6837529/#comments" class="reply">3回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="8999871">
        <div class="main review-item" id="8999871">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/123396490/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u123396490-5.jpg">
        </a>

        <a href="https://www.douban.com/people/123396490/" class="name">椒盐栗子糕</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2017-12-20" class="main-meta">2017-12-20 02:57:26</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8999871/">给自己的思考单</a></h2>

                <div id="review_8999871_short" class="review-short" data-rid="8999871">
                    <div class="short-content">

                        E01风格和挑高喜欢，楼梯思考改进方法，公共空间过道略压抑（宽高比不协调）； E02没太大亮点，对楼镜子维护成本（时间精力）高且容易掉下伤人，需要固定装置； E03非常酷了！觉得我需要了解一下建筑，中间层的床挡板有点矮怕摔下来，可以用医院那种可以收起来的挡板； E04真的...

                        &nbsp;(<a href="javascript:;" id="toggle-8999871-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8999871_full" class="hidden">
                    <div id="review_8999871_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8999871" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8999871">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8999871" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8999871">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8999871/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="8778110">
        <div class="main review-item" id="8778110">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/djaensnhyu/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/u2092712-7.jpg">
        </a>

        <a href="https://www.douban.com/people/djaensnhyu/" class="name">阡陌Jane</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2017-08-27" class="main-meta">2017-08-27 14:53:01</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8778110/">向天再借五百年</a></h2>

                <div id="review_8778110_short" class="review-short" data-rid="8778110">
                    <div class="short-content">

                          周末在家无聊，没书看的焦虑和刷手机的空虚感，让我在豆瓣想看中翻出之前有兴趣的影片，看到了梦想改造家。  这种类型赏心悦目又有意义，于是连着看了几期，最开始都是各种偷空间和架构改造，把小房子做出各种实用，还有大师般气场的设计师入住，看着精彩纷呈。  但看到今天...

                        &nbsp;(<a href="javascript:;" id="toggle-8778110-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8778110_full" class="hidden">
                    <div id="review_8778110_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8778110" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8778110">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8778110" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8778110">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8778110/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        
    
    <div data-cid="8748748">
        <div class="main review-item" id="8748748">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/144280826/" class="avator">
            <img width="24" height="24" src="https://img1.doubanio.com/icon/user_normal.jpg">
        </a>

        <a href="https://www.douban.com/people/144280826/" class="name">柳浪仙子</a>

            <span class="allstar50 main-title-rating" title="力荐"></span>

        <span content="2017-08-14" class="main-meta">2017-08-14 11:41:24</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8748748/">充满正能量感人的节目！</a></h2>

                <div id="review_8748748_short" class="review-short" data-rid="8748748">
                    <div class="short-content">

                        节目基本选残疾人家庭改造，切切实实帮助了这些弱势群体！很正能量！ 节目中家的那种温情，家人之间本来存在的隔阂、担心、理解,在屋子的重新设计装修的过程，设计师不仅考虑了屋子改造的硬件问题，还对亲人间的矛盾等进行了充分的思考并提出意见(例如有个妈妈开服装店儿子一直...

                        &nbsp;(<a href="javascript:;" id="toggle-8748748-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8748748_full" class="hidden">
                    <div id="review_8748748_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8748748" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8748748">
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8748748" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8748748">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8748748/#comments" class="reply">0回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>


        <div class="fold-hd">
                <a class="btn-unfold" href="#">有一些影评被折叠了</a>
                    <a class="qa" href="https://help.douban.com/opinion?app=movie#t1-q2">为什么被折叠？</a>
            <div class="qa-tip">评论被折叠，是因为发布这条评论的帐号行为异常。评论仍可以被展开阅读，对发布人的账号不造成其他影响。如果认为有问题，可以<a href="https://help.douban.com/help/ask?category=movie">联系</a>豆瓣电影。</div>
        </div>
        <div class="fold-bd">
                
    
    <div data-cid="7261967">
        <div class="main review-item" id="7261967">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/108213028/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u108213028-2.jpg">
        </a>

        <a href="https://www.douban.com/people/108213028/" class="name">╭⌒一世的姻缘</a>

            <span class="allstar20 main-title-rating" title="较差"></span>

        <span content="2014-12-20" class="main-meta">2014-12-20 11:11:36</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/7261967/">第六集</a></h2>

                <div id="review_7261967_short" class="review-short" data-rid="7261967">
                    <div class="short-content">

                        [第6集] 这一期比较失败。刚开始出现的廖设计师显然是一个不怎么负责人的设计师：1.对于老人栽培20年的树砍掉，而不是想办法保留下来；2.对于房屋结构没有一个准确的认识；3.对于室内外高差问题考虑不周全；4.对设计不负责任，不怎么经常去沟通，出现问题一味的指责施工方；5....

                        &nbsp;(<a href="javascript:;" id="toggle-7261967-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_7261967_full" class="hidden">
                    <div id="review_7261967_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="7261967" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-7261967">
                                8
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="7261967" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-7261967">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/7261967/#comments" class="reply">5回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

                
    
    <div data-cid="8332332">
        <div class="main review-item" id="8332332">

            
    
    <header class="main-hd">
        <a href="https://www.douban.com/people/64382079/" class="avator">
            <img width="24" height="24" src="https://img3.doubanio.com/icon/u64382079-6.jpg">
        </a>

        <a href="https://www.douban.com/people/64382079/" class="name">瑶啊瑶</a>

            <span class="allstar40 main-title-rating" title="推荐"></span>

        <span content="2017-02-02" class="main-meta">2017-02-02 22:38:12</span>


    </header>


            <div class="main-bd">

                <h2><a href="https://movie.douban.com/review/8332332/">几个感想</a></h2>

                <div id="review_8332332_short" class="review-short" data-rid="8332332">
                    <div class="short-content">

                        看节目过程中的几个感想： 1.妈呀，这样的房子竟然能装这么多东西（住这么多人）； 2.感触最深的是上海城隍庙和北京四合院那两家，那么破也不搬多半因为那是学区房吧……； 3.有些虽然房子破，但是装修费并不便宜，证明北京上海土著还是很有实力的； 4.没看全，但是发现上海及...

                        &nbsp;(<a href="javascript:;" id="toggle-8332332-copy" class="unfold" title="展开">展开</a>)
                    </div>
                </div>

                <div id="review_8332332_full" class="hidden">
                    <div id="review_8332332_full_content" class="full-content"></div>
                </div>

                <div class="action">
                    <a href="javascript:;" class="action-btn up" data-rid="8332332" title="有用">
                        <img src="https://img3.doubanio.com/f/zerkalo/536fd337139250b5fb3cf9e79cb65c6193f8b20b/pics/up.png" />
                        <span id="r-useful_count-8332332">
                                1
                        </span>
                    </a>
                    <a href="javascript:;" class="action-btn down" data-rid="8332332" title="没用">
                        <img src="https://img3.doubanio.com/f/zerkalo/68849027911140623cf338c9845893c4566db851/pics/down.png" />
                        <span id="r-useless_count-8332332">
                        </span>
                    </a>
                    <a href="https://movie.douban.com/review/8332332/#comments" class="reply">2回应</a>

                    <a href="javascript:;;" class="fold hidden">收起</a>
                </div>
            </div>
        </div>
    </div>

        </div>


    

    

    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/4af3abf4186d1238.js"></script>
    <!-- COLLECTED CSS -->
</div>








            <p class="pl">
                &gt;
                <a href="reviews">
                    更多剧评10篇
                </a>
            </p>
</section>

<!-- COLLECTED JS -->

    <br/>

        <div class="section-discussion">
                <p class="discussion_link">
    <a href="https://movie.douban.com/subject/25942176/tv_discuss">&gt; 查看 梦想改造家 第一季 的分集短评（全部72条）</a>
</p>

        </div>


    <script type="text/javascript">
        $(function(){if($.browser.msie && $.browser.version == 6.0){
            var $info = $('#info'),
                maxWidth = parseInt($info.css('max-width'));
            if($info.width() > maxWidth) {
                $info.width(maxWidth);
            }
        }});
    </script>


            </div>
            <div class="aside">
                


    








        






    

<script id="episode-tmpl" type="text/x-jsrender">
<div id="tv-play-source" class="play-source">
    <div class="cross">
        <span style="color:#494949; font-size:16px">{{:cn}}</span>
        <span style="cursor:pointer">✕</span>
    </div>
    <div class="episode-list">
        {{for playlist}}
            <a href="{{:play_link}}&episode={{:ep}}" target="_blank">{{:ep}}集</a>
        {{/for}}
     <div>
 </div>
</script>

<div class="gray_ad">
    
    <h2>
        在哪儿看这部剧集
            &nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;&nbsp;&middot;
    </h2>

    
    <ul class="bs">
                <li>
                        <a class="playBtn" data-cn="爱奇艺视频" data-source="9"  href="javascript: void 0;">
                            爱奇艺视频
                        </a>
                    <span class="buylink-price"><span>
                        免费观看
                    </span></span>
                </li>

    </ul>
</div>


    <!-- douban ad begin -->
    <div id="dale_movie_subject_top_right"></div>
    <div id="dale_movie_subject_top_middle"></div>
    <!-- douban ad end -->

    



<style type="text/css">
    .m4 {margin-bottom:8px; padding-bottom:8px;}
    .movieOnline {background:#FFF6ED; padding:10px; margin-bottom:20px;}
    .movieOnline h2 {margin:0 0 5px;}
    .movieOnline .sitename {line-height:2em; width:160px;}
    .movieOnline td,.movieOnline td a:link,.movieOnline td a:visited{color:#666;}
    .movieOnline td a:hover {color:#fff;}
    .link-bt:link,
    .link-bt:visited,
    .link-bt:hover,
    .link-bt:active {margin:5px 0 0; padding:2px 8px; background:#a8c598; color:#fff; -moz-border-radius: 3px; -webkit-border-radius: 3px; border-radius: 3px; display:inline-block;}
</style>



    







    
    <div class="tags">
        
        
    <h2>
        <i class="">豆瓣成员常用的标签</i>
              · · · · · ·
    </h2>

        <div class="tags-body">
                <a href="/tag/综艺" class="">综艺</a>
                <a href="/tag/装修" class="">装修</a>
                <a href="/tag/设计" class="">设计</a>
                <a href="/tag/真人秀" class="">真人秀</a>
                <a href="/tag/中国" class="">中国</a>
                <a href="/tag/大陆" class="">大陆</a>
                <a href="/tag/2014" class="">2014</a>
                <a href="/tag/家庭" class="">家庭</a>
        </div>
    </div>


    <div id="dale_movie_subject_inner_middle"></div>
    <div id="dale_movie_subject_download_middle"></div>
        








<div id="subject-doulist">
    
    
    <h2>
        <i class="">以下豆列推荐</i>
              · · · · · ·
            <span class="pl">
            (
                <a href="https://movie.douban.com/subject/25942176/doulists">全部</a>
            )
            </span>
    </h2>


    
    <ul>
            <li>
                <a href="https://www.douban.com/doulist/39460326/" target="_blank">中国大陆（精华）</a>
                <span>(Deffrro-ZI)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/20184555/" target="_blank">一个人看电视（二）</a>
                <span>(鹿小羽)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/40159250/" target="_blank">【那些优秀综艺佳作精选集】【收藏必备】</a>
                <span>(熊孩子)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/2609657/" target="_blank">2014年待看新片，惊喜或者地雷，看过才知晓</a>
                <span>(糯米女巫话痨蘇)</span>
            </li>
            <li>
                <a href="https://www.douban.com/doulist/3425671/" target="_blank">综艺</a>
                <span>(开文Vincent)</span>
            </li>
    </ul>

</div>

        








<div id="subject-others-interests">
    
    
    <h2>
        <i class="">谁在看这部剧集</i>
              · · · · · ·
    </h2>

    
    <ul class="">
            
            <li class="">
                <a href="https://www.douban.com/people/31347229/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u31347229-3.jpg" class="pil" alt="亦池">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/31347229/" class="">亦池</a>
                    <div class="">
                        今天上午
                        看过
                        <span class="allstar50" title="力荐"></span>
                    </div>
                </div>
            </li>
            
            <li class="">
                <a href="https://www.douban.com/people/zxx549527/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u120176442-3.jpg" class="pil" alt="鬼塚乱步">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/zxx549527/" class="">鬼塚乱步</a>
                    <div class="">
                        今天上午
                        想看
                        
                    </div>
                </div>
            </li>
            
            <li class="">
                <a href="https://www.douban.com/people/131959648/" class="others-interest-avatar">
                    <img src="https://img3.doubanio.com/icon/u131959648-1.jpg" class="pil" alt="Nirvana↘★">
                </a>
                <div class="others-interest-info">
                    <a href="https://www.douban.com/people/131959648/" class="">Nirvana↘★</a>
                    <div class="">
                        今天上午
                        看过
                        <span class="allstar40" title="推荐"></span>
                    </div>
                </div>
            </li>
    </ul>

    
    <div class="subject-others-interests-ft">
        
            <a href="https://movie.douban.com/subject/25942176/doings">261人在看</a>
                &nbsp;/&nbsp;
            <a href="https://movie.douban.com/subject/25942176/collections">3991人看过</a>
                &nbsp;/&nbsp;
            <a href="https://movie.douban.com/subject/25942176/wishes">1775人想看</a>
    </div>

</div>



    
    

<!-- douban ad begin -->
<div id="dale_movie_subject_middle_right"></div>
<script type="text/javascript">
    (function (global) {
        if(!document.getElementsByClassName) {
            document.getElementsByClassName = function(className) {
                return this.querySelectorAll("." + className);
            };
            Element.prototype.getElementsByClassName = document.getElementsByClassName;

        }
        var articles = global.document.getElementsByClassName('article'),
            asides = global.document.getElementsByClassName('aside');

        if (articles.length > 0 && asides.length > 0 && articles[0].offsetHeight >= asides[0].offsetHeight) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_movie_subject_middle_right');
        }
    })(this);
</script>
<!-- douban ad end -->



    <br/>

    
<p class="pl">订阅梦想改造家 第一季的影评: <br/><span class="feed">
    <a href="https://movie.douban.com/feed/subject/25942176/reviews"> feed: rss 2.0</a></span></p>


            </div>
            <div class="extra">
                
    
<!-- douban ad begin -->
<div id="dale_movie_subject_bottom_super_banner"></div>
<script type="text/javascript">
    (function (global) {
        var body = global.document.body,
            html = global.document.documentElement;

        var height = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight, html.scrollHeight, html.offsetHeight);
        if (height >= 2000) {
            (global.DoubanAdSlots = global.DoubanAdSlots || []).push('dale_movie_subject_bottom_super_banner');
        }
    })(this);
</script>
<!-- douban ad end -->


            </div>
        </div>
    </div>

        
    <div id="footer">
            <div class="footer-extra"></div>
        
<span id="icp" class="fleft gray-link">
    &copy; 2005－2019 douban.com, all rights reserved 北京豆网科技有限公司
</span>

<a href="https://www.douban.com/hnypt/variformcyst.py" style="display: none;"></a>

<span class="fright">
    <a href="https://www.douban.com/about">关于豆瓣</a>
    · <a href="https://www.douban.com/jobs">在豆瓣工作</a>
    · <a href="https://www.douban.com/about?topic=contactus">联系我们</a>
    · <a href="https://www.douban.com/about?policy=disclaimer">免责声明</a>
    
    · <a href="https://help.douban.com/?app=movie" target="_blank">帮助中心</a>
    · <a href="https://www.douban.com/doubanapp/">移动应用</a>
    · <a href="https://www.douban.com/partner/">豆瓣广告</a>
</span>

    </div>

    </div>
    <script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/4c2365f5127232f.js"></script><script type="text/javascript">
                var if_logined='false';
var subject_id='25942176';
Do.ready("handlebarsjs",function(){var e=$("#comment-tmpl").html(),n=Handlebars.compile(e),t=Handlebars.compile('&gt; <a href="comments?sort=time">更多短评{{num}}条</a>');if_logined="true"===if_logined,Handlebars.registerHelper("content_tmpl",function(e){var n="";return n=e.length>200?['<span class="short">'+e.substring(0,200)+"...</span>",'<span class="hide-item full">'+e+"</span>",'<span class="expand">(<a href="javascript:;">展开</a>)</span>'].join(""):"<span>"+e+"</span>",new Handlebars.SafeString(n)});var a=$("#comments-section");a.delegate(".tab-hd a","click",function(e){e.preventDefault();var n=$(this);return a.find(".tab-hd a").removeClass("on").end().find(".tab").hide().end().find("#"+n.data("id")+"-comments").show(),n.addClass("on"),$.get("/blank?track-"+n.attr("id")),!1}).delegate("#new-comments-tab","click",function(e){$(this).data("clicked")||($(this).data("clicked",!0),$.get("/j/subject/"+subject_id+"/comments",function(e){if(1===e.retcode){var a=(e.result,n({comments:e.result.normal,if_logined:if_logined}));$("#new-comments #normal").html(a),e.result.spammed.length>0&&($("#new-comments .fold-bd").append(n({comments:e.result.spammed,if_logined:if_logined})),$("#new-comments .fold-hd").removeClass("hide")),e.result.total_num>4&&$("#new-comments #total-num").html(t({num:e.result.total_num})),load_event_monitor($("#new-comments"))}}))})});
                $(function(){$("body").delegate(".btn-unfold","click",function(e){e.preventDefault();var t=$(e.target),d=t.parent(".fold-hd");d.slideUp().next().slideDown()}),$("body").delegate(".comment-item .expand a","click",function(e){e.preventDefault();var t=$(e.target),d=t.parents("p");$short=d.find(".short"),$hide=d.find(".hide-item"),t.hasClass("isfold")?(t.removeClass("isfold").text("展开"),$short.show(),$hide.hide()):(t.addClass("isfold").text("收起"),$short.hide(),$hide.show())})});
            </script><script type="text/javascript" src="https://img3.doubanio.com/misc/mixed_static/5be90cabe1ab4b46.js"></script>
        
        
    <link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/shire/8377b9498330a2e6f056d863987cc7a37eb4d486/css/ui/dialog.css" />
    <link rel="stylesheet" type="text/css" href="https://img3.doubanio.com/f/movie/1d829b8605b9e81435b127cbf3d16563aaa51838/css/movie/mod/reg_login_pop.css" />
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/77323ae72a612bba8b65f845491513ff3329b1bb/js/do.js" data-cfg-autoload="false"></script>
    <script type="text/javascript" src="https://img3.doubanio.com/f/shire/4ea3216519a6183c7bcd4f7d1a6d4fd57ce1a244/js/ui/dialog.js"></script>
    <script type="text/javascript">
        var HTTPS_DB='https://www.douban.com';
var account_pop={open:function(o,e){e?referrer="?referrer="+encodeURIComponent(e):referrer="?referrer="+window.location.href;var n="",i="",t=382;"reg"===o?(n="用户注册",i="https://accounts.douban.com/popup/login?source=movie#popup_register",t=480):"login"===o&&(n="用户登录",i="https://accounts.douban.com/popup/login?source=movie");var r=document.location.protocol+"//"+document.location.hostname,a=dui.Dialog({width:478,title:n,height:t,cls:"account_pop",isHideTitle:!0,modal:!0,content:"<iframe scrolling='no' frameborder='0' width='478' height='"+t+"' src='"+i+"' name='"+r+"'></iframe>"},!0),c=a.node;if(c.undelegate(),c.delegate(".dui-dialog-close","click",function(){var o=$("body");o.find("#login_msk").hide(),o.find(".account_pop").remove()}),$(window).width()<478){var u="";"reg"===o?u=HTTPS_DB+"/accounts/register"+referrer:"login"===o&&(u=HTTPS_DB+"/accounts/login"+referrer),window.location.href=u}else a.open();$(window).bind("message",function(o){"https://accounts.douban.com"===o.originalEvent.origin&&(c.find("iframe").css("height",o.originalEvent.data),c.height(o.originalEvent.data),a.update())})}};Douban&&Douban.init_show_login&&(Douban.init_show_login=function(o){var e=$(o);e.click(function(){var o=e.data("ref")||"";return account_pop.open("login",o),!1})}),Do(function(){$("body").delegate(".pop_register","click",function(o){o.preventDefault();var e=$(this).data("ref")||"";return account_pop.open("reg",e),!1}),$("body").delegate(".pop_login","click",function(o){o.preventDefault();var e=$(this).data("ref")||"";return account_pop.open("login",e),!1})});
    </script>

    
    
    
    




    
<script type="text/javascript">
    (function (global) {
        var newNode = global.document.createElement('script'),
            existingNode = global.document.getElementsByTagName('script')[0],
            adSource = '//erebor.douban.com/',
            userId = '',
            browserId = 'auSsK8Dk5cg',
            criteria = '7:建筑|7:何赛飞|7:骆新|7:设计|7:佟瑞欣|7:家庭|7:装修|7:金星|7:黄豆豆|7:中国|7:戴娇倩|7:施琰|7:真人秀|7:家装改造|7:2014|7:大陆|7:综艺|3:/subject/25942176/',
            preview = '',
            debug = false,
            adSlots = ['dale_movie_subject_top_icon', 'dale_movie_subject_top_right', 'dale_movie_subject_top_middle', 'dale_movie_subject_inner_middle', 'dale_movie_subject_download_middle'];

        global.DoubanAdRequest = {src: adSource, uid: userId, bid: browserId, crtr: criteria, prv: preview, debug: debug};
        global.DoubanAdSlots = (global.DoubanAdSlots || []).concat(adSlots);

        newNode.setAttribute('type', 'text/javascript');
        newNode.setAttribute('src', 'https://img3.doubanio.com/f/adjs/dd37385211bc8deb01376096bfa14d2c0436a98c/ad.release.js');
        newNode.setAttribute('async', true);
        existingNode.parentNode.insertBefore(newNode, existingNode);
    })(this);
</script>











    
  









<script type="text/javascript">
var _paq = _paq || [];
_paq.push(['trackPageView']);
_paq.push(['enableLinkTracking']);
(function() {
    var p=(('https:' == document.location.protocol) ? 'https' : 'http'), u=p+'://fundin.douban.com/';
    _paq.push(['setTrackerUrl', u+'piwik']);
    _paq.push(['setSiteId', '100001']);
    var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
    g.type='text/javascript';
    g.defer=true;
    g.async=true;
    g.src=p+'://img3.doubanio.com/dae/fundin/piwik.js';
    s.parentNode.insertBefore(g,s);
})();
</script>

<script type="text/javascript">
var setMethodWithNs = function(namespace) {
  var ns = namespace ? namespace + '.' : ''
    , fn = function(string) {
        if(!ns) {return string}
        return ns + string
      }
  return fn
}

var gaWithNamespace = function(fn, namespace) {
  var method = setMethodWithNs(namespace)
  fn.call(this, method)
}

var _gaq = _gaq || []
  , accounts = [
      { id: 'UA-7019765-1', namespace: 'douban' }
    , { id: 'UA-7019765-19', namespace: '' }
    ]
  , gaInit = function(account) {
      gaWithNamespace(function(method) {
        gaInitFn.call(this, method, account)
      }, account.namespace)
    }
  , gaInitFn = function(method, account) {
      _gaq.push([method('_setAccount'), account.id]);
      _gaq.push([method('_setSampleRate'), '5']);

      
  _gaq.push([method('_addOrganic'), 'google', 'q'])
  _gaq.push([method('_addOrganic'), 'baidu', 'wd'])
  _gaq.push([method('_addOrganic'), 'soso', 'w'])
  _gaq.push([method('_addOrganic'), 'youdao', 'q'])
  _gaq.push([method('_addOrganic'), 'so.360.cn', 'q'])
  _gaq.push([method('_addOrganic'), 'sogou', 'query'])
  if (account.namespace) {
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣'])
    _gaq.push([method('_addIgnoredOrganic'), 'douban'])
    _gaq.push([method('_addIgnoredOrganic'), '豆瓣网'])
    _gaq.push([method('_addIgnoredOrganic'), 'www.douban.com'])
  }

      if (account.namespace === 'douban') {
        _gaq.push([method('_setDomainName'), '.douban.com'])
      }

        _gaq.push([method('_setCustomVar'), 1, 'responsive_view_mode', 'desktop', 3])

        _gaq.push([method('_setCustomVar'), 2, 'login_status', '0', 2]);

      _gaq.push([method('_trackPageview')])
    }

for(var i = 0, l = accounts.length; i < l; i++) {
  var account = accounts[i]
  gaInit(account)
}


;(function() {
    var ga = document.createElement('script');
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    ga.setAttribute('async', 'true');
    document.documentElement.firstChild.appendChild(ga);
})()
</script>








      
    

    <!-- brand13-docker-->

  <script>_SPLITTEST=''</script>
</body>

</html>


`

func TestPageHtml(t *testing.T) {
	fmt.Println(time.Duration(rand.Float64() * 120e9).Seconds())
}
func Test_crc32(t *testing.T) {
	e := strings.Split(reflect.TypeOf(&JianShu{}).String(), ".")[1]
	fmt.Println(e)

	//crc32q := crc32.ChecksumIEEE([]byte("33861"))
	//fmt.Println(crc32q)
}

func TestZhiHu_PageHtml(t *testing.T) {

	var title, watch, view string

	var tag []string
	node, _ := html.Parse(ioutil.NopCloser(bytes.NewBuffer(ZhiHuByte)))
	PageHtml(node, &title, &watch, &view, &tag)

	fmt.Println(title, watch, view, tag)

	//log.Println(err)
	//log.Println(node)
}

func TestDouBanPageHtml(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//movie
	u, _ := url.Parse("https://movie.douban.com/subject/26394152/")
	node, err := html.Parse(ioutil.NopCloser(bytes.NewBuffer([]byte(html22))))
	if err != nil {
		return
	}

	douBanId, err := strconv.Atoi(strings.Split(u.Path, "/")[2])
	if err != nil {
		log.Fatal(u, err)
	}
	model := &AsukaDouBan{
		DouBanId: int64(douBanId),
	}

	if strings.HasPrefix(u.String(), "https://movie.douban.com") {
		model.Cate = "电影"
	}
	if strings.HasPrefix(u.String(), "https://book.douban.com") {
		model.Cate = "图书"
	}

	DouBanPageHtml(node, model)

	printAll(reflect.ValueOf(model).Elem())

	//update
	existsModel := &AsukaDouBan{
		UrlCrc32: int64(crc32.ChecksumIEEE([]byte(u.String()))),
		Url:      u.String(),
	}
	if ok, err := database.Mysql().Get(existsModel); ok && err == nil {
		model.Version = existsModel.Version
		database.Mysql().Id(existsModel.Id).Update(model)
	}

	fmt.Println("")
	//book
	u, _ = url.Parse("https://book.douban.com/subject/30389935/?icn=index-latestbook-subject")
	douBanId, _ = strconv.Atoi(strings.Split(u.Path, "/")[2])
	node, err = html.Parse(ioutil.NopCloser(bytes.NewBuffer([]byte(html33))))
	if err != nil {
		return
	}

	book := &AsukaDouBan{
		DouBanId: int64(douBanId),
	}

	if strings.HasPrefix(u.String(), "https://movie.douban.com") {
		book.Cate = "电影"
	}
	if strings.HasPrefix(u.String(), "https://book.douban.com") {
		book.Cate = "图书"
	}

	DouBanPageHtml(node, book)

	printAll(reflect.ValueOf(book).Elem())
	fmt.Println("")
	//movie
	u, _ = url.Parse("https://movie.douban.com/subject/25942176/")
	douBanId, _ = strconv.Atoi(strings.Split(u.Path, "/")[2])
	node, err = html.Parse(ioutil.NopCloser(bytes.NewBuffer([]byte(html44))))
	if err != nil {
		return
	}

	movie := &AsukaDouBan{
		DouBanId: int64(douBanId),
	}

	if strings.HasPrefix(u.String(), "https://movie.douban.com") {
		movie.Cate = "电影"
	}
	if strings.HasPrefix(u.String(), "https://book.douban.com") {
		movie.Cate = "图书"
	}

	DouBanPageHtml(node, movie)

	printAll(reflect.ValueOf(movie).Elem())

	fmt.Println("")
	//movie
	u, _ = url.Parse("https://movie.douban.com/subject/2279835/")
	douBanId, _ = strconv.Atoi(strings.Split(u.Path, "/")[2])
	node, err = html.Parse(ioutil.NopCloser(bytes.NewBuffer([]byte(html66))))
	if err != nil {
		return
	}

	movie3 := &AsukaDouBan{
		DouBanId: int64(douBanId),
	}

	if strings.HasPrefix(u.String(), "https://movie.douban.com") {
		movie3.Cate = "电影"
	}
	if strings.HasPrefix(u.String(), "https://book.douban.com") {
		movie3.Cate = "图书"
	}

	DouBanPageHtml(node, movie3)

	if movie3.DateStr == "" {
		DouBanPageHtmlSecondly(node, movie3)
	}

	printAll(reflect.ValueOf(movie3).Elem())
	fmt.Println("")
	//book
	u, _ = url.Parse("https://book.douban.com/subject/6438136/")
	douBanId, _ = strconv.Atoi(strings.Split(u.Path, "/")[2])
	node, err = html.Parse(ioutil.NopCloser(bytes.NewBuffer([]byte(html76))))
	if err != nil {
		return
	}

	movie3 = &AsukaDouBan{
		DouBanId: int64(douBanId),
	}

	if strings.HasPrefix(u.String(), "https://movie.douban.com") {
		movie3.Cate = "电影"
	}
	if strings.HasPrefix(u.String(), "https://book.douban.com") {
		movie3.Cate = "图书"
	}

	DouBanPageHtml(node, movie3)

	if movie3.DateStr == "" {
		DouBanPageHtmlSecondly(node, movie3)
	}

	printAll(reflect.ValueOf(movie3).Elem())

}

var isSubject = regexp.MustCompile(`douban.com/subject/[0-9]+/?$`).MatchString

func TestDouBan_EntryUrl(t *testing.T) {
	str := "https://book.douba.com/subject/27614904/"
	log.Println(isSubject(str))
	str = "https://book.douban.com/subject/27614904/"
	log.Println(isSubject(str))
	str = "https://book.douban.com/subject/27614904/123"
	log.Println(isSubject(str))
	str = "https://book.douban.com/subject/27614904/hot/2313"
	log.Println(isSubject(str))
	str = "https://book.douban.com/subject/27614904/hot"
	log.Println(isSubject(str))
	str = "https://book.douban.com/subject/27614904"
	log.Println(isSubject(str))
	str = "https://book.douban.com/subject/27614904/12312"
	log.Println(isSubject(str))
	str = "https://book.douban.com/subjet/27614904/"
	log.Println(isSubject(str))
}

func TestTime(t *testing.T) {
	jsonStr := []byte(`{
  "@context":"http://schema.org",
  "@type":"Book",
  "workExample": [], 
  "name" : "李天命思考艺术",
  "author": 
  [
    {   
      "@type": "Person",
      "name": "戎子由\"
    }   
    ,   
    {   
      "@type": "Person",
      "name": "梁沛霖" 
    }   
  ]
,
  "url" : "https://book.douban.com/subject/2298835/",
  "isbn" : "9789623572972",
  "sameAs": "https://book.douban.com/subject/2298835/"
}
`)

	//for i, ch := range jsonStr {
	//	if ch == 92 {
	//jsonStr[i] = '/'
	//}
	//}

	//log.Println(byte('\\'))
	//jsonStr = strings.Replace(jsonStr, `\`, `/`, len(jsonStr))

	fmt.Println(json.Valid(jsonStr))

	d := make(map[string]interface{})
	log.Println(json.Unmarshal(jsonStr, &d))
	//fmt.Println(d)

	fmt.Println(time.Parse("06/01/02", "92/12/11"))
	fmt.Println(time.Parse("060102", "921211"))
	fmt.Println(time.Parse("06/1/2", "92/12/11"))
	fmt.Println(time.Parse("0612", "921211"))
}

func TestEEEEE(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Write([]byte(`{"data":[{"directors":["郭帆"],"rate":"7.9","cover_x":1786,"star":"40","title":"流浪地球","url":"https:\/\/movie.douban.com\/subject\/26266893\/","casts":["屈楚萧","吴京","李光洁","吴孟达","赵今麦"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2545472803.webp","id":"26266893","cover_y":2500},{"directors":["韩寒"],"rate":"7.0","cover_x":7142,"star":"35","title":"飞驰人生","url":"https:\/\/movie.douban.com\/subject\/30163509\/","casts":["沈腾","黄景瑜","尹正","张本煜","尹昉"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2542973862.webp","id":"30163509","cover_y":10000},{"directors":["周星驰","邱礼涛"],"rate":"5.8","cover_x":1078,"star":"30","title":"新喜剧之王","url":"https:\/\/movie.douban.com\/subject\/4840388\/","casts":["王宝强","鄂靖文","张全蛋","景如洋","张琪"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2541240741.webp","id":"4840388","cover_y":1512},{"directors":["大卫·雷奇"],"rate":"7.3","cover_x":4618,"star":"35","title":"死侍2：我爱我家","url":"https:\/\/movie.douban.com\/subject\/26588308\/","casts":["瑞恩·雷诺兹","乔什·布洛林","朱利安·迪尼森","斯蒂芬·卡皮契奇","莎姬·贝兹"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2545479945.webp","id":"26588308","cover_y":6842},{"directors":["大卫·叶茨"],"rate":"7.1","cover_x":2048,"star":"35","title":"神奇动物：格林德沃之罪","url":"https:\/\/movie.douban.com\/subject\/26147417\/","casts":["埃迪·雷德梅恩","凯瑟琳·沃特斯顿","约翰尼·德普","裘德·洛","埃兹拉·米勒"],"cover":"https://img1.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2538826177.webp","id":"26147417","cover_y":3035},{"directors":["宁浩"],"rate":"6.4","cover_x":960,"star":"30","title":"疯狂的外星人","url":"https:\/\/movie.douban.com\/subject\/25986662\/","casts":["黄渤","沈腾","汤姆·派福瑞","马修·莫里森","徐峥"],"cover":"https://img1.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2541901817.webp","id":"25986662","cover_y":1359},{"directors":["温子仁"],"rate":"7.8","cover_x":6750,"star":"40","title":"海王","url":"https:\/\/movie.douban.com\/subject\/3878007\/","casts":["杰森·莫玛","艾梅柏·希尔德","威廉·达福","帕特里克·威尔森","妮可·基德曼"],"cover":"https://img1.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2541280047.webp","id":"3878007","cover_y":10000},{"directors":["特拉维斯·奈特"],"rate":"7.1","cover_x":1080,"star":"35","title":"大黄蜂","url":"https:\/\/movie.douban.com\/subject\/26394152\/","casts":["海莉·斯坦菲尔德","小豪尔赫·兰登伯格","约翰·塞纳","杰森·德鲁克","帕梅拉·阿德龙"],"cover":"https://img1.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2541662397.webp","id":"26394152","cover_y":1590},{"directors":["黄家康","赵霁"],"rate":"7.9","cover_x":5315,"star":"40","title":"白蛇：缘起","url":"https:\/\/movie.douban.com\/subject\/30331149\/","casts":["张喆","杨天翔","唐小喜","刘薇","张遥函"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2544313786.webp","id":"30331149","cover_y":7441},{"directors":["徐誉庭","许智彦"],"rate":"8.6","cover_x":1697,"star":"45","title":"谁先爱上他的","url":"https:\/\/movie.douban.com\/subject\/27119586\/","casts":["邱泽","谢盈萱","陈如山","黄圣球","周洺甫"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2546961193.webp","id":"27119586","cover_y":2500},{"directors":["饶晓志"],"rate":"8.1","cover_x":1429,"star":"40","title":"无名之辈","url":"https:\/\/movie.douban.com\/subject\/27110296\/","casts":["陈建斌","任素汐","潘斌龙","章宇","王砚辉"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2539661066.webp","id":"27110296","cover_y":2000},{"directors":["闫非","彭大魔"],"rate":"6.5","cover_x":679,"star":"35","title":"西虹市首富","url":"https:\/\/movie.douban.com\/subject\/27605698\/","casts":["沈腾","宋芸桦","张一鸣","张晨光","常远"],"cover":"https://img1.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2529206747.webp","id":"27605698","cover_y":950},{"directors":["庄文强"],"rate":"8.0","cover_x":714,"star":"40","title":"无双","url":"https:\/\/movie.douban.com\/subject\/26425063\/","casts":["周润发","郭富城","张静初","冯文娟","廖启智"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2535096871.webp","id":"26425063","cover_y":1000},{"directors":["文牧野"],"rate":"9.0","cover_x":1080,"star":"45","title":"我不是药神","url":"https:\/\/movie.douban.com\/subject\/26752088\/","casts":["徐峥","王传君","周一围","谭卓","章宇"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2519070834.webp","id":"26752088","cover_y":1512},{"directors":["张开宙"],"rate":"7.5","cover_x":3547,"star":"40","title":"知否知否应是绿肥红瘦","url":"https:\/\/movie.douban.com\/subject\/26928226\/","casts":["赵丽颖","冯绍峰","朱一龙","施诗","张佳宁"],"cover":"https://img1.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2537131688.webp","id":"26928226","cover_y":4729},{"directors":["吕克·贝松"],"rate":"9.4","cover_x":658,"star":"45","title":"这个杀手不太冷","url":"https:\/\/movie.douban.com\/subject\/1295644\/","casts":["让·雷诺","娜塔莉·波特曼","加里·奥德曼","丹尼·爱罗","彼得·阿佩尔"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p511118051.webp","id":"1295644","cover_y":980},{"directors":["涅提·蒂瓦里"],"rate":"9.0","cover_x":4500,"star":"45","title":"摔跤吧！爸爸","url":"https:\/\/movie.douban.com\/subject\/26387939\/","casts":["阿米尔·汗","法缇玛·萨那·纱卡","桑亚·玛荷塔","阿帕尔夏克提·库拉那","沙克希·坦沃"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2457983084.webp","id":"26387939","cover_y":6300},{"directors":["克里斯托弗·诺兰"],"rate":"9.2","cover_x":4858,"star":"45","title":"星际穿越","url":"https:\/\/movie.douban.com\/subject\/1889243\/","casts":["马修·麦康纳","安妮·海瑟薇","杰西卡·查斯坦","卡西·阿弗莱克","迈克尔·凯恩"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2206088801.webp","id":"1889243","cover_y":7200},{"directors":["林超贤"],"rate":"8.3","cover_x":1429,"star":"40","title":"红海行动","url":"https:\/\/movie.douban.com\/subject\/26861685\/","casts":["张译","黄景瑜","海清","杜江","蒋璐霞"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2514119443.webp","id":"26861685","cover_y":2000},{"directors":["是枝裕和"],"rate":"8.7","cover_x":3300,"star":"45","title":"小偷家族","url":"https:\/\/movie.douban.com\/subject\/27622447\/","casts":["中川雅也","安藤樱","松冈茉优","城桧吏","佐佐木美结"],"cover":"https://img3.doubanio.com\/view\/photo\/s_ratio_poster\/public\/p2530599636.webp","id":"27622447","cover_y":4659}]}`))
	decoder := json.NewDecoder(buf)
	movieJson := make(map[string]interface{})
	if err := decoder.Decode(&movieJson); err != nil {
		return
	}

	var urls []*url.URL
	if data, ok := movieJson["data"]; ok {
		for _, item := range data.([]interface{}) {
			if m, ok := item.(map[string]interface{}); ok {
				if rawUrl, ok := m["url"]; ok {
					u, _ := url.Parse(rawUrl.(string))
					urls = append(urls, u)
				}
			}
		}
	}

	fmt.Println(urls)
}

func TestUrlP(t *testing.T) {
	l, _ := url.Parse("https://movie.douban.com/tag/%E6%97%A5%E8%AF%AD%E6%A0%87%E7%AD%BE?start=100&type=T")

	fmt.Println(l.Query().Encode())
	query := make(url.Values)
	if l.Query().Get("start") != "" {
		query.Set("start", l.Query().Get("start"))
	}
	if l.Query().Get("type") != "" {
		query.Set("type", l.Query().Get("type"))
	}
	fmt.Println(query.Encode())

	fmt.Println(l.Scheme + "://" + l.Host + l.Path)
}

func printAll(v reflect.Value) {
	s := v
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())

		if f.Kind().String() == "struct" {
			x1 := reflect.ValueOf(f.Interface())
			fmt.Printf("type2: %s\n", x1)
			printAll(x1)
		}
	}
}

func TestCrawl(t *testing.T) {
	res := regexp.MustCompile("<a\\shref=\"([^\"]+)\"\\sclass=\"topic-link\">(.+)</a>").FindAllStringSubmatch(html200, -1)

	for _, v := range res {
		log.Println(v[1], v[2])
	}
	//log.Println(len(res))
	//log.Println(res)
}
