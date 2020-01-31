$('document').ready(function(){
    $('#get_exec').on('click', function(){
        var key = $('#get_key').val();
        if ($.trim(key) != ''){
            $.get( "http://127.0.0.1:8081/v1/immurestproxy/item/"+btoa(key), function(data) {
                finaldata = {}
                finaldata.key = atob(data.key)
                finaldata.value = atob(data.value)
                $('#get_res').text(JSON.stringify(finaldata, null, 2));
            }).fail(function(error) {
                $('#get_res').text(JSON.stringify(error.responseJSON, null, 2));
            });
        };
    });
    $('#set_exec').on('click', function(){
        var post = {}
        post.key = btoa($('#set_key').val());
        post.value = btoa($('#set_value').val());
        if ($.trim(post.key) != '' && $.trim(post.value) != ''){
            $.post( "http://127.0.0.1:8081/v1/immurestproxy/item", JSON.stringify(post) , function(data) {
                $('#set_res').text(JSON.stringify(data, null, 2));
            }).fail(function(data) {
                $('#set_res').text(JSON.stringify(data.responseJSON, null, 2));
            });
        };
    });
    $('#cp_exec').on('click', function(){
        index = $('#cp_index').val();
        if ($.trim(index) != ''){
            $.get( "http://127.0.0.1:8081/v1/immurestproxy/inclusionproof/"+index, function(data) {
                $('#cp_res').text(JSON.stringify(data, null, 2));
            }).fail(function(error) {
                $('#cp_res').text(JSON.stringify(error.responseJSON, null, 2));
            });
        };
    });
    $('#ip_exec').on('click', function(){
        index = $('#ip_index').val();
        if ($.trim(index) != ''){
            $.get( "http://127.0.0.1:8081/v1/immurestproxy/inclusionproof/"+index, function(data) {
                $('#ip_res').text(JSON.stringify(data, null, 2));
            }).fail(function(error) {
                $('#ip_res').text(JSON.stringify(error.responseJSON, null, 2));
            });
        };
    });

});