$('document').ready(function(){
    $('#get_exec').on('click', function(){
        var key = $('#get_key').val();
        if ($.trim(key) != ''){
            var sgpost = {}
            sgpost.key = {}
            sgpost.key.key = btoa($('#get_key').val());

            $.post( "http://127.0.0.1:8081/v1/immurestproxy/item/safe/get",JSON.stringify(sgpost), function(data) {
                $('#get_res').text(atob(data.item.value));
                if(data.verified){
                    $('#id_getsuccess').addClass( "show" );
                    $("#id_getsuccess").fadeTo(2000, 500).slideUp(500, function(){
                        $("#id_getsuccess").slideUp(500);
                    });
                }else{
                    $('#id_getdanger').addClass( "show" );
                    $('#id_getsuccess').removeClass( "show" );
                }
            }).fail(function(error) {
                $('#get_res').text(error.responseJSON.message);
            });
        };
    });
    $('#set_exec').on('click', function(){
        var post = {}
        post.kv = {}
        post.kv.key = btoa($('#set_key').val());
        post.kv.value = btoa($('#set_value').val());
        if ($.trim(post.kv.key) != '' && $.trim(post.kv.value) != ''){
            $.post( "http://127.0.0.1:8081/v1/immurestproxy/item/safe", JSON.stringify(post) , function(data) {
                if(data.verified){
                    $('#id_setsuccess').addClass( "show" );
                    $("#id_setsuccess").fadeTo(2000, 500).slideUp(500, function(){
                        $("#id_setsuccess").slideUp(500);
                    });
                }else{
                    $('#id_setdanger').addClass( "show" );
                    $('#id_setsuccess').removeClass( "show" );
                }

                $('#set_res').text(atob(post.kv.value));
            }).fail(function(data) {
                $('#set_res').text(JSON.stringify(data.responseJSON, null, 2));
            });
        };
    });
    $('#cp_exec').on('click', function(){
        index = $('#cp_index').val();
        if ($.trim(index) != ''){
            $.get( "http://127.0.0.1:8081/v1/immurestproxy/consistencyproof/"+index, function(data) {
                fd = {}
                fd.first = data.first
                fd.second = data.second
                fd.secondRoot = stringToHex(atob(data.secondRoot))
                $('#cp_res').text(JSON.stringify(fd, null, 2));
            }).fail(function(error) {
                $('#cp_res').text(JSON.stringify(error.responseJSON, null, 2));
            });
        };
    });
    $('#ip_exec').on('click', function(){
        index = $('#ip_index').val();
        if ($.trim(index) != ''){
            $.get( "http://127.0.0.1:8081/v1/immurestproxy/inclusionproof/"+index, function(data) {

                fd = {}
                fd.at = data.at
                fd.index = data.index
                fd.root = stringToHex(atob(data.root))
                fd.leaf = stringToHex(atob(data.leaf))
                fd.path = []
                for (var i = 0; i < data.path.length; i++) {
                    ele = stringToHex(atob(data.path[i]));
                    fd.path.push(ele)
                }

                $('#ip_res').text(JSON.stringify(fd, null, 2));
            }).fail(function(error) {
                $('#ip_res').text(JSON.stringify(error.responseJSON, null, 2));
            });
        };
    });
    $('#h_exec').on('click', function(){
        hk = $('#h_key').val();
        if ($.trim(hk) != ''){
            $.get( "http://127.0.0.1:8081/v1/immurestproxy/history/"+btoa(hk), function(data) {
                fdata = [];
                for (var i = 0; i < data.items.length; i++) {
                    ele = {};
                    ele.key = atob(data.items[i].key);
                    ele.value = atob(data.items[i].value);
                    ele.index = data.items[i].index;
                    fdata.push(ele)
                }
                $('#h_res').text(JSON.stringify(fdata, null, 2));
            }).fail(function(error) {
                $('#h_res').text(JSON.stringify(error.responseJSON, null, 2));
            });
        };
    });

    $('#scan_exec').on('click', function(){
        scpost = {}
        scpost.prefix = btoa($('#scan_prefix').val());
        scpost.offset = btoa($('#scan_offset').val());
        scpost.limit = parseInt($('#scan_limit').val());
        scpost.reverse = $('#scan_reverse').is(":checked");
        if ($.trim(scpost.prefix) != ''){
            $.post( "http://127.0.0.1:8081/v1/immurestproxy/item/scan", JSON.stringify(scpost) , function(data) {

                $('#scan_res').text(JSON.stringify(serializeItems(data.items), null, 2));
            }).fail(function(data) {
                $('#scan_res').text(JSON.stringify(data.responseJSON, null, 2));
            });
        };
    });

});

function serializeItems(items) {
    fdata = [];
    for (var i = 0; i < items.length; i++) {
        ele = {};
        ele.key = atob(items[i].key);
        ele.value = atob(items[i].value);
        ele.index = items[i].index;
        fdata.push(ele)
    }
    return fdata;
}

function stringToHex (tmp) {
    var str = '',
        i = 0,
        tmp_len = tmp.length,
        c;

    for (; i < tmp_len; i += 1) {
        c = tmp.charCodeAt(i);
        str += d2h(c);
    }
    return str;
}

function d2h(d) {
    return d.toString(16);
}
