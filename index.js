/*
This file is (C) Captain ALM
Under the Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License
*/
var EntryData = []
var EntryIndices = []
var SortOrderStateI = true
var SortOrderBStateI = true
var SortOrderEnabled = false
var SortValue = ""
var OrderValue = ""
function SetupJS() {
    SetupIndexArray()
    SetupJSTheme()
    SetupJSHPL()
    SetupJSHSO()
    SetupJSSOI()
}
function CreateEntry(id, name, videourl, videotype, start, end, duration) {
    EntryData[id] = {
        name: name,
        videourl: videourl,
        videotype: videotype,
        start: Date.parse(start),
        end: Date.parse(end),
        duration : parseInt(duration, 10)
    };
}
function CreateVideoPlaceholder(id) {
    var imgPH = document.createElement("img")
    imgPH.src = PlayImageURL
    imgPH.id = "play-"+id
    imgPH.alt = "Play Video"
    imgPH.title = "Play"
    imgPH.width = 360
    imgPH.style.cursor = "pointer"
    if (document.addEventListener) {
        imgPH.addEventListener("click", function () {ActivateVideo(id);})
    } else {
        imgPH.setAttribute("onclick", "ActivateVideo("+id+");")
        imgPH.onclick = function () {ActivateVideo(id);}
    }
    document.getElementById("video-" + id).appendChild(imgPH)
}
function ActivateVideo(id) {
    var holder = document.getElementById("video-" + id)
    holder.removeChild(document.getElementById("play-"+id))
    var vid = document.createElement("video")
    vid.controls = true
    vid.width = 360
    var vids = document.createElement("source")
    vids.src = EntryData[id].videourl
    vids.type = EntryData[id].videotype
    var vida = document.createElement("a")
    vida.href = EntryData[id].videourl
    vida.innerText = "The Video"
    vid.appendChild(vids)
    vid.appendChild(vida)
    holder.appendChild(vid)
    if (vid.play) {vid.play();}
}
function SetupIndexArray() {
    for (var i = 0; i < EntryData.length; i++) {
        EntryIndices[i] = i
    }
}
function SetupJSTheme() {
    var th = document.getElementById("theme")
    th.href = "#"
    if (document.addEventListener) {
        th.addEventListener("click", ToggleTheme)
    } else {
        th.setAttribute("onclick", "ToggleTheme();")
        th.onclick = ToggleTheme
    }
}
function cReplaceHistory() {
    ReplaceHistory(window.location.href)
}
function ReplaceHistory(url) {
    if (window.history) {
        if (window.history.replaceState) {
            window.history.replaceState({
                light: !!document.getElementById("so-theme"),
                order: document.getElementById("so-order").value,
                sort: document.getElementById("so-sort").value
            }, "", url);
            console.log("REPLACE")
        }
    }
}
function PushHistory(url) {
    var s = true
    if (window.history) {
        if (window.history.pushState) {
            window.history.pushState({
                light: !!document.getElementById("so-theme"),
                order: document.getElementById("so-order").value,
                sort: document.getElementById("so-sort").value
            }, "", url);
            console.log("PUSH")
            s = false
        }
    }
    if (s) {
        document.location.href = url
    }
}
function ToggleTheme() {
    var th = document.getElementById("theme")
    var thimg = document.getElementById("theme-img")
    var thsty = document.getElementById("style-theme")
    var logo = document.getElementById("logo")
    var url = document.location.href
    url = url.split("#", 1)[0].split('?', 1)[0]
    if (document.getElementById("so-theme")) {
        thimg.src = SunImageURL
        thimg.alt = "()"
        th.title = "Switch to Light Mode"
        document.getElementById("so-form").removeChild(document.getElementById("so-theme"))
        logo.href = "?"
        PushHistory(url+"?"+TheParameters+"#")
        thsty.href = CssDarkURL
    } else {
        thimg.src = MoonImageURL
        thimg.alt = "{"
        th.title = "Switch to Dark Mode"
        var thi = document.createElement("input")
        thi.name = "light"
        thi.type = "hidden"
        thi.id = "so-theme"
        document.getElementById("so-form").appendChild(thi)
        logo.href = "?light"
        if (TheParameters === "") {
            PushHistory(url+"?light#")
        } else {
            PushHistory(url+"?light&"+TheParameters+"#")
        }
        thsty.href = CssLightURL
    }
}
function SetupJSHPL(){
    if (window.history) {
        if (window.history.pushState) {
            window.addEventListener("load", cReplaceHistory)
            window.addEventListener("popstate", HandleHistoryPop)
        }
    }
}
function HandleHistoryPop(event) {
    console.log("POP")
    console.log(event.state)
    if (event.state) {
        SortOrderEnabled = false
        var isnl = !document.getElementById("so-theme")
        if ((event.state.light && isnl) || (!event.state.light && !isnl)) {ToggleTheme();}
        document.getElementById("so-order").value = event.state.order
        document.getElementById("so-sort").value = event.state.sort
        EntrySort(event.state.order, event.state.sort)
        SortOrderEnabled = true
    }
}
function SetupJSHSO() {
    var pb = document.getElementById("sort-menu-button")
    var pane = document.getElementById("so-pane")
    if (document.addEventListener) {
        document.addEventListener("click", HandleGlobalClick)
        pb.addEventListener("mouseover", HandleSortOrderBEnter)
        pb.addEventListener("mouseout", HandleSortOrderBLeave)
        pane.addEventListener("mouseover", HandleSortOrderEnter)
        pane.addEventListener("mouseout", HandleSortOrderLeave)
    } else {
        document.parentElement.setAttribute("onclick", "HandleGlobalClick();")
        pb.setAttribute("onmouseover", "HandleSortOrderBEnter();")
        pb.setAttribute("onmouseout", "HandleSortOrderBLeave();")
        pane.setAttribute("onmouseover", "HandleSortOrderEnter();")
        pane.setAttribute("onmouseout", "HandleSortOrderLeave();")
        document.parentElement.onclick = HandleGlobalClick
        pb.onmouseover = HandleSortOrderBEnter
        pb.onmouseout = HandleSortOrderBLeave
        pane.onmouseover = HandleSortOrderEnter
        pane.onmouseout = HandleSortOrderLeave
    }
}
function HandleGlobalClick() {
    if (SortOrderStateI && SortOrderBStateI) {document.getElementById("sort-menu").checked = false;}
}
function HandleSortOrderBEnter() {
    SortOrderBStateI = false
}
function HandleSortOrderBLeave(){
    SortOrderBStateI = true
}
function HandleSortOrderEnter() {
    SortOrderStateI = false
}
function HandleSortOrderLeave(){
    SortOrderStateI = true
}
function SetupJSSOI() {
    var submit = document.getElementById("so-submit")
    if (submit.parentNode) {submit.parentNode.removeChild(submit);}
    var oc = document.getElementById("so-order")
    OrderValue = oc.value
    var sc = document.getElementById("so-sort")
    SortValue = sc.value
    if (document.addEventListener) {
        oc.addEventListener("change", HandleSortOrderChange)
        sc.addEventListener("change", HandleSortOrderChange)
    } else {
        oc.setAttribute("onchange", "HandleSortOrderChange();")
        sc.setAttribute("onchange", "HandleSortOrderChange();")
        oc.onchange = HandleSortOrderChange
        sc.onchange = HandleSortOrderChange
    }
    SortOrderEnabled = true
}
function HandleSortOrderChange() {
    if (SortOrderEnabled) {EntrySort(document.getElementById("so-order").value, document.getElementById("so-sort").value);}
}
function EntrySort(o, s) {
    var ts = s.toString().toLowerCase()
    var chg = false
    if (SortValue !== s) {
        chg = true
        SortValue = s
    }
    if (chg || OrderValue !== o) {
        if (ts === "desc" || ts === "descending") {
            ts = -1
        } else {
            ts = 1
        }
        var to = o.toString().toLowerCase()
        if (to === "start") {
            if (ts < 0) {
                EntryIndices = EntryIndices.sort(SortStartD)
            } else {
                EntryIndices = EntryIndices.sort(SortStartA)
            }
        } else if (to === "end") {
            if (ts < 0) {
                EntryIndices = EntryIndices.sort(SortEndD)
            } else {
                EntryIndices = EntryIndices.sort(SortEndA)
            }
        } else if (to === "name") {
            if (ts < 0) {
                EntryIndices = EntryIndices.sort(SortNameD)
            } else {
                EntryIndices = EntryIndices.sort(SortNameA)
            }
        } else if (to === "duration") {
            if (ts < 0) {
                EntryIndices = EntryIndices.sort(SortDurationD)
            } else {
                EntryIndices = EntryIndices.sort(SortDurationA)
            }
        }
        chg = true
        OrderValue = o
    }
    if (chg) {
        TheParameters = "order="+OrderValue+"&sort="+SortValue
        var url = document.location.href
        url = url.split("#", 1)[0].split('?', 1)[0]
        if (document.getElementById("so-theme")) {
            PushHistory(url+"?light&"+TheParameters)
        } else {
            PushHistory(url+"?"+TheParameters)
        }
        for (var i = 0; i < EntryIndices.length; i++) {
            var tNode = document.getElementById("entry-"+EntryIndices[i])
            var pNode = tNode.parentNode
            tNode = pNode.removeChild(tNode)
            pNode.appendChild(tNode)
        }
    }
}
function SortStartA(a, b) {
    if (EntryData[a].start < EntryData[b].start) {
        return -1
    } else {
        return  1
    }
}
function SortStartD(a, b) {
    if (EntryData[a].start > EntryData[b].start) {
        return -1
    } else {
        return  1
    }
}
function SortEndA(a, b) {
    if (EntryData[a].end < EntryData[b].end) {
        return -1
    } else {
        return  1
    }
}
function SortEndD(a, b) {
    if (EntryData[a].end > EntryData[b].end) {
        return -1
    } else {
        return  1
    }
}
function SortNameA(a, b) {
    if (EntryData[a].name < EntryData[b].name) {
        return -1
    } else {
        return  1
    }
}
function SortNameD(a, b) {
    if (EntryData[a].name > EntryData[b].name) {
        return -1
    } else {
        return  1
    }
}
function SortDurationA(a, b) {
    if (EntryData[a].duration < EntryData[b].duration) {
        return -1
    } else {
        return  1
    }
}
function SortDurationD(a, b) {
    if (EntryData[a].duration > EntryData[b].duration) {
        return -1
    } else {
        return  1
    }
}