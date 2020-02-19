$('document').ready(function(){
    $('#get_exec').on('click', function(){
        var key = $('#get_key').val();
        if ($.trim(key) != ''){
            var sgpost = {}
            sgpost.key = {}
            sgpost.key.key = btoa($('#get_key').val());
            $.post( ENDPOINT+"/v1/immurestproxy/item/safe/get",JSON.stringify(sgpost), function(item) {
                let table = getTable([item])
                $('#get_res').html(table);
                $('#c_get_res').addClass( "visible" ).removeClass("invisible");
                if(item.verified){
                    $('#id_getsuccess').addClass( "show" );
                    $("#id_getsuccess").fadeTo(2000, 500).slideUp(500, function(){
                        $("#id_getsuccess").slideUp(500);
                    });
                }else{
                    $('#id_getdanger').addClass( "show" );
                    $('#id_getsuccess').removeClass( "show" );
                }
            }).fail(function(error) {
                $('#c_get_res').addClass( "visible" ).removeClass("invisible");
                $('#get_res').text(error.responseJSON.message);
            });
        };
    });
    $('#set_exec').on('click', function(){
        var post = {}
        post.kv = {}
        post.kv.key = btoa($('#set_key').val());
        let val = {}
        val.ts = Date.now()
        val.val= $('#set_value').val();
        post.kv.value = btoa(JSON.stringify(val, null, 2))

        if ($.trim(post.kv.key) != '' && $.trim(post.kv.value) != ''){
            $.post( ENDPOINT+"/v1/immurestproxy/item/safe", JSON.stringify(post) , function(data) {
                //$('#c_set_res').addClass( "visible" ).removeClass("invisible");
                if(data.verified){
                    $('#id_setsuccess').addClass( "show" );
                    $("#id_setsuccess").fadeTo(2000, 500).slideUp(500, function(){
                        $("#id_setsuccess").slideUp(500);
                    });
                }else{
                    $('#id_setdanger').addClass( "show" );
                    $('#id_setsuccess').removeClass( "show" );
                }
                //$('#set_res').text(atob(post.kv.value));
            }).fail(function(data) {
                $('#set_res').text(JSON.stringify(data.responseJSON, null, 2));
            });
        };
    });
    $('#cp_exec').on('click', function(){
        index = $('#cp_index').val();
        if ($.trim(index) != ''){
            $.get( ENDPOINT+"/v1/immurestproxy/consistencyproof/"+index, function(data) {
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
            $.get( ENDPOINT+"/v1/immurestproxy/inclusionproof/"+index, function(data) {

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
            $.get( ENDPOINT+"/v1/immurestproxy/history/"+btoa(hk), function(data) {
                let table = getTable(data)
                $('#h_res').html(table);

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
            $.post( ENDPOINT+"/v1/immurestproxy/item/scan", JSON.stringify(scpost) , function(data) {

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

String.prototype.replaceAll = function(searchStr, replaceStr) {
    var str = this;

    // escape regexp special characters in search string
    searchStr = searchStr.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&');
    return str.replace(new RegExp(searchStr, 'gi'), replaceStr);
};

function getTable(data){
    $('#c_h_res').addClass( "visible" ).removeClass("invisible");
    let tbody = $('<tbody>')
    for (var i = 0; i < data.length; i++) {

        //let temp = JSON.stringify(atob(data.items[i].value), null, 2)
        //temp = temp.replaceAll("\\n", "\n").replaceAll("\\t", "\t").replaceAll("\\\"", "\"")

        let temp = JSON.parse(atob(data[i].value))
        let row = $('<tr>');
        let td0 = $('<th scope="row"></th>').text(String( data[i].index));
        let td1 = $('<td>');
        let td2 = $('<td>');
        let span1 = $('<pre>').text(temp.ts);
        let span2 = $('<pre>').text(temp.val);
        td1.append(span1);
        td2.append(span2);
        row.append(td0);
        row.append(td1);
        row.append(td2);
        tbody.append(row);
    }
    //$('#h_res').text(fdata, null, 2);
    let table = $('<table>').addClass('table table-striped');
    table.append('<thead><tr><th scope="col">#</th><th scope="col">Timestamp</th><th scope="col">Value</th></tr></thead>')
    table.append(tbody)
    return table
}
