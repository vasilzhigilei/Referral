$(function () {
    $('#profileform').on('submit', function (e) {
        e.preventDefault();
        $.ajax({
            url: "/updateuser",
            type: "POST",
            data: $('#profileform').serialize(),
            success: function (data) {
                $('#alertresult').html("<div class=\"alert alert-success\" role=\"alert\">\n" +
                    "  Success! Profile updated.\n" +
                    "</div>");
            },
            error: function (jXHR, textStatus, errorThrown) {
                alert(errorThrown);
            }
        });
    });
});