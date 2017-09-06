/**
 * Created by 18147 on 2017/9/1.
 */
function viewOrHide() {
    var tableObj = document.getElementById("ApiTable");
    for (var i = 1; i < tableObj.rows.length; i++) {    //遍历Table的所有Row
        $("tr#"+tableObj.rows[i].id).show();
    }

    for (var i = 1; i < tableObj.rows.length; i++) {    //遍历Table的所有Row
        tableInfo = tableObj.rows[i].cells[6].innerText;   //获取Table中单元格的内容
        if (tableInfo != $('#Api_Service_id').val()) {
            $("tr#"+tableObj.rows[i].id).hide();
        }
    }
}
function validateUri_required(field,id,alerttxt)
{
    var id_value


    id_value = id.value;

    console.log("id="+id_value);

    if(id_value != "-1") {
        with (field)
        {
            var IfRepeat=false;
            var tableObj = document.getElementById("ApiTable");
            console.log($('#flag').val())
            console.log("flag="+$('#flag').val())
            for (var i = 1; i < tableObj.rows.length; i++) {    //遍历Table的所有Row
                tableInfo = tableObj.rows[i].cells[1].innerText;   //获取Table中单元格的内容
                console.log("tableInfo="+tableInfo)
                if (tableInfo == value&&$('#flag').val()!=-1) {
                    IfRepeat = true;
                    break;
                }
            }
            if (value==null||value==""||IfRepeat) {
                alert(alerttxt);
                return false
            } else {
                return true
            }
        }
    } else {
        return true
    }
}

function validateOperation_required(field,alerttxt)
{
    with (field)
    {
        if (value==null||value=="") {
            alert(alerttxt);
            return false
        } else {
            return true
        }
    }
}

$(function() {

    $("#addTr1").on('click', function () {
        //获取table最后一行 $("#tab tr:last")
        //获取table第一行 $("#tab tr").eq(0)
        //获取table倒数第二行 $("#tab tr").eq(-2)
        var $tr=$("#MockHeaders tr:last").clone();
        $("#MockHeaders").append($tr);
    })

    $("#addTr2").on('click', function () {
        //获取table最后一行 $("#tab tr:last")
        //获取table第一行 $("#tab tr").eq(0)
        //获取table倒数第二行 $("#tab tr").eq(-2)
        var $tr=$("#mockCookies tr:last").clone();
        $("#mockCookies").append($tr);
    })

    $("#subMock").on('click', function () {
        var str = '{"contentType":"'+$("#ContentType").val()+'","value":"'+$("#MockValue").val()+'",'+'"headers":['
        var tableObj = document.getElementById("MockHeaders");
        for (var i = 1; i < tableObj.rows.length; i++) {    //遍历Table的所有Row
            tableInfo = $("#MockHeaders tr").eq(1).find('input').val();   //获取Table中单元格的内容
            console.log("tableInfo="+tableInfo)
            if (i!=tableObj.rows.length-1) {
                str = str+'{"'+ $("#MockHeaders tr").eq(i).find('input').val()+'":"'+$("#MockHeaders tr").eq(i).find('input').eq(1).val()+'"},'
            } else {
                str = str+'{"'+ $("#MockHeaders tr").eq(i).find('input').val()+'":"'+$("#MockHeaders tr").eq(i).find('input').eq(1).val()+'"}'
            }
        }
        str =str + '],"cookies":['
        var tableObj2 = document.getElementById("mockCookies");
        for (var i = 1; i < tableObj2.rows.length; i++) {    //遍历Table的所有Row
            if (i!=tableObj2.rows.length-1) {
                str = str+'"'+$("#mockCookies tr").eq(i).find('input').val()+'",'
            } else {
                str = str+'"'+$("#mockCookies tr").eq(i).find('input').val()+'"'
            }
        }
        str =str + ']}'
        $('#Api_mock').val(str);
        if ($("#clearMock").prop( "checked" )==true) {
            $('#Api_mock').val('')
        }
        $('#addUserModal').modal('hide');
        console.log("MockStr="+str)
    });





    $("#example").hover(function () {
        $(this).stop().animate({
            opacity: '1'
        }, 600);
    }, function () {
        $(this).stop().animate({
            opacity: '0.6'
        }, 1000);
    }).on('click', function () {
        $("body").append("<div id='mask'></div>");
        $("#mask").addClass("mask").fadeIn("slow");
        $("#LoginBox").fadeIn("slow");
    });
//
//按钮的透明度
    $("#loginbtn").hover(function () {
        $(this).stop().animate({
            opacity: '1'
        }, 600);
    }, function () {
        $(this).stop().animate({
            opacity: '0.8'
        }, 1000);
    });
//文本框不允许为空---按钮触发

    $("#loginbtn").on('click', function () {
        var txtName = $("#txtName").val();
        var txtPwd = $("#txtPwd").val();
        if (txtName == "" || txtName == undefined || txtName == null) {
            if (txtPwd == "" || txtPwd == undefined || txtPwd == null) {
                $(".warning").css({ display: 'block' });
            }
            else {
                $("#warn").css({ display: 'block' });
                $("#warn2").css({ display: 'none' });
            }
        }
        else {
            if (txtPwd == "" || txtPwd == undefined || txtPwd == null) {
                $("#warn").css({ display: 'none' });
                $(".warn2").css({ display: 'block' });
            }
            else {
                $(".warning").css({ display: 'none' });
            }
        }

        var params = $("#loginForm").serialize();

        $.ajax({
            url: "/login",
            type: 'POST',
            data: params,
            success: function (data, status, returndata) {
                if (returndata.getResponseHeader('status')=="true") {
                    location.reload()
                } else {
                    alert("登陆失败")
                    location.reload()
                }

            },
            error: function (returndata) {
                alert("登陆请求失败")
            }
        });

    });
//文本框不允许为空---单个文本触发
    $("#txtName").on('blur', function () {
        var txtName = $("#txtName").val();
        if (txtName == "" || txtName == undefined || txtName == null) {
            $("#warn").css({ display: 'block' });
        }
        else {
            $("#warn").css({ display: 'none' });
        }
    });
    $("#txtName").on('focus', function () {
        $("#warn").css({ display: 'none' });
    });
//
    $("#txtPwd").on('blur', function () {
        var txtName = $("#txtPwd").val();
        if (txtName == "" || txtName == undefined || txtName == null) {
            $("#warn2").css({ display: 'block' });
        }
        else {
            $("#warn2").css({ display: 'none' });
        }
    });
    $("#txtPwd").on('focus', function () {
        $("#warn2").css({ display: 'none' });
    });
    //关闭
    $(".close_btn").hover(function () {
            $(this).css({ color: 'black' }) },
        function () { $(this).css({ color: '#999' }) }).on('click', function () {
        $("#LoginBox").fadeOut("fast");
        $("#mask").css({ display: 'none' });
    });


    $("#upload").click(function () {
        console.log("name="+$("#example").text())
        if ($("#example").text() == '登陆') {
            alert("请登陆！！！！！！！！")
            return false
        }

        var formData = new FormData($("#uploadFile")[0]);

        $.ajax({
            url: "/uploadFile",
            type: 'POST',
            data: formData,
            cache: false,
            contentType: false,
            processData: false,
            success: function (returndata) {
                location.reload()
            },
            error: function (returndata) {
            }
        });

    });

    //ajax提交
    $("#ajaxBtn").click(function () {
        console.log("ajax submit")
        console.log("name="+$("#example").text())
        if ($("#example").text() == '登陆')
        {
            alert("请登陆！！！！！！！！")
            return false
        }
        with ($("#form1"))
        {
            if (validateUri_required(Api_uri,flag,"Api_uri 不能为空和重复!")==false) {
                Api_uri.focus();
                return false
            }
            if (validateOperation_required(Api_name,"Api_name 不能为空!")==false) {
                Api_name.focus();
                return false
            }
        };
        if ($("#IfStr").prop( "checked" )==true) {
            $("#Api_uri").val($("#Api_uri").val()+"$");
        }
        var params = $("#form1").serialize();
        console.log("params="+params)
        $.ajax({
            type: "POST",
            url: "/gateway/admin/add_api",
            data: params,
            success: function (msg) {
                location.reload()
            }
        });
    })
})


function modifyData(id) {
    var tdArray = document.querySelectorAll('#ApiTable #tr-'+id+' td:not(.operation)');
    console.log("tdArray="+tdArray[0])
    $("#IfStr").attr("checked", false)
    $('#Api_name').val(tdArray[0].innerText);
    $('#Api_uri').val(tdArray[1].innerText);
    $('#Api_method').val(tdArray[2].innerText);
    $('#Api_need_login').val(tdArray[3].innerText);
    $('#Api_display_name').val(tdArray[4].innerText);
    $('#Api_status').val(tdArray[5].innerText);
    $('#Api_Service_id').val(tdArray[6].innerText);
    $('#Api_Service_Provider_name').val(tdArray[7].innerText);
    $('#Api_mock').val(tdArray[8].innerText);
    $('#Api_desc').val(tdArray[9].innerText);
    $('#Api_Api_id').val(tdArray[10].innerText);
    $('#flag').val(-1)
    document.getElementById("Api_uri").focus();
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
        var parm = "Api_Api_id=" + id;//构造URL参数
        //xmlobj.open("POST", "{dede:global.cfg_templeturl/}/../include/weather.php", true); //调用weather.php
        xmlobj.open("POST", "/gateway/admin/delete_api", true); //调用weather.php
        xmlobj.setRequestHeader("cache-control","no-cache");
        xmlobj.setRequestHeader("contentType","text/html;charset=uft-8"); //指定发送的编码
        xmlobj.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;");  //设置请求头信息
        //                xmlobj.onreadystatechange = StatHandler;  //判断URL调用的状态值并处理
        xmlobj.send(parm); //设置为发送给服务器数据
        location.reload()
    }
}







var $table = $('#ApiTable'),
    $alertBtn = $('#alertBtn'),
    full_screen = false,
    window_height;

$().ready(function(){
    console.log("i'm in ready");
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
        formatShowingRows: function(pageFrom, pageTo, totalRows){
            //do nothing here, we don't want to show the text "showing x of y from..."
        },
        formatRecordsPerPage: function(pageNumber){
            return pageNumber + " rows visible";
        },
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