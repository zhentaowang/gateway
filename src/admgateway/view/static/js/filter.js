/**
 * Created by 18147 on 2017/9/2.
 */
function validate_required(field,alerttxt)
{
    with (field)
    {
        var patt = /^-?[0-9]\d*$/;
        var result = patt.test(value);
        if (value==null||value==""||!result) {
            alert(alerttxt);
            return false
        } else {
            return true
        }
    }
}

function modifyData(id) {
    var tdArray = document.querySelectorAll('#FilterTable #tr-'+id+' td:not(.operation)');
    $('#Filter_id').val(tdArray[0].innerText);
    $('#Filter_Api_id').val(tdArray[2].innerText);
    $('#Filter_name').val(tdArray[1].innerText);
    $('#Filter_seq').val(tdArray[3].innerText);
    document.getElementById("Filter_Api_id").focus();
}
var xmlobj; //定义XMLHttpRequest对象
function CreateXMLHttpRequest()
{
    if(window.ActiveXObject)
    //如果当前浏览器支持Active Xobject，则创建ActiveXObject对象
    {
        //xmlobj = new ActiveXObject("Microsoft.XMLHTTP");
        try {
            xmlobj = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (e) {
            try {
                xmlobj = new ActiveXObject("Microsoft.XMLHTTP");
            } catch (E) {
                xmlobj = false;
            }
        }
    }
    else if(window.XMLHttpRequest)
    //如果当前浏览器支持XMLHttp Request，则创建XMLHttpRequest对象
    {
        xmlobj = new XMLHttpRequest();
    }
}
function deleteData(id) {
    if(confirm('您确定删除这行数据吗？')) {
        console.log('执行删除操作');
        CreateXMLHttpRequest(); //创建对象
        var parm = "Filter_id=" + id;//构造URL参数
        //xmlobj.open("POST", "{dede:global.cfg_templeturl/}/../include/weather.php", true); //调用weather.php
        xmlobj.open("POST", "/gateway/admin/delete_filter", true); //调用weather.php
        xmlobj.setRequestHeader("cache-control","no-cache");
        xmlobj.setRequestHeader("contentType","text/html;charset=uft-8"); //指定发送的编码
        xmlobj.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;");  //设置请求头信息
//                xmlobj.onreadystatechange = StatHandler;  //判断URL调用的状态值并处理
        xmlobj.send(parm); //设置为发送给服务器数据
        location.reload()
    }
}








var $table = $('#FilterTable'),
    $alertBtn = $('#alertBtn'),
    full_screen = false,
    window_height;

$().ready(function(){

    window_height = $(window).height();
    table_height = window_height - 20;


    $table.bootstrapTable({
        toolbar: ".toolbar",

        showRefresh: false,
        search: true,
        showToggle: false,
        showColumns: true,
        pagination: true,
        striped: true,
        sortable: true,
        height: table_height,
        pageSize: 100000,
        pageList: [25,50,100,100000],
        icons: {
            refresh: 'fa fa-refresh',
            toggle: 'fa fa-th-list',
            columns: 'fa fa-columns',
            detailOpen: 'fa fa-plus-circle',
            detailClose: 'fa fa-minus-circle'
        }
    });

    $(window).resize(function () {
        $table.bootstrapTable('resetView');
    });
});