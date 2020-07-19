$(function () {
    $('profileform').on('submit', function (e) {
        e.preventDefault();
        $.ajax({
            url: $(form).attr('action') || window.location.pathname,
            type: "POST",
            data: $('profileform').serialize(),
            success: function (data) {
                console.log("lifts logged")
            },
            error: function (jXHR, textStatus, errorThrown) {
                alert(errorThrown);
            }
        });
    });
});