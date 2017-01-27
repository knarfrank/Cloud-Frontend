


function fetchScan() {
    var loaded = 0;
    if(window.location.hash) {
        url = window.location.hash.substring(1);
        $('#pageTile').html("Site Report - " + url);
        loaded = getcontent(url);

        if(loaded == 0) {

            var intervalId = window.setInterval(
            function () {

                if(loaded == 0) {
                    loaded = getcontent(url);
                    if(loaded == 1) {
                        clearInterval(intervalId);
                    }
                }
                console.log('on interval...' + loaded);
            }, 2000);
        }
    }
}

function getcontent(url) {
    jQuery.ajaxSetup({async:false});
    var loaded = 0;
    console.log("Getting content for " + url);
    $.get( "/get?url="+url, function(r) {
        if(r != "null") {

            json = $.parseJSON(r);
            if(json.length > 0) {
                console.log("Unhiding stuff");
                $('#loadingstuff').hide();
                $('#certificatesblock').show();
            }
            for (var i = 0; i < json.length; i++) {
                scan = json[i];
                $('#certificatesarea').append('<h5>Certificate #1</h5>' +
                                              '<table width="100%" class="table table-bordered table-hover" id="dataTables-example"><tbody>' +
                                              '<tr class="gradeX"><td><strong>Subject</strong></td><td>'+scan.Subject+'</td></tr>' +
                                              '<tr class="gradeX"><td><strong>Fingerprint Sha1</strong></td><td>'+scan.Fingerprint+'</td></tr>' +
                                              '<tr class="gradeX"><td><strong>Common name </strong></td><td>'+scan.CommonName+'</td></tr>' +
                                              '<tr class="gradeX"><td><strong>Alternative names</strong></td><td></td></tr>' +
                                              '<tr class="gradeX"><td><strong>Issuer</strong></td><td>'+scan.Issuer+'</td></tr>' +
                                              '<tr class="gradeX"><td><strong>Valid From</strong></td><td>'+scan.ValidFrom+'</td></tr>' +
                                              '<tr class="gradeX"><td><strong>Valid Till</strong></td><td>'+scan.ValidTill+'</td></tr>' +
                                              '<tr class="gradeX"><td><strong>Signature Algorithm</strong></td><td>'+scan.SignatureAlgorithm+'</td></tr>' +
                                              '<tr class="gradeX"><td><strong>Key</strong></td><td></td></tr></tbody></table>' );

            }


            loaded = 1;
        } else {
            loaded = 0;
        }

    });
    jQuery.ajaxSetup({async:true});
    return loaded;
}
function runScan(url) {
    $.post('/queue?url='+url, {}, function(r, d) {
		document.location = "results.html#"+url;
	});
}


function fillRecentScans() {
    $.post('/getrecent', {"limit": 5}, function(r, d) {
        json = $.parseJSON(r);
        for (var i = 0; i < json.length; i++) {
            scan = json[i];
            $('#recentscans').append('<a href="results.html#'+scan.Url+'" class="list-group-item">'+
                                     '<i class="fa fa-compress fa-fw"></i> '+scan.Url+' '+
                                     '<span class="pull-right text-muted small"><em>'+scan.Ip+'</em></span></a>');

        }
	});
}
