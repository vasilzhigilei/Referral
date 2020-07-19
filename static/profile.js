function submitForm(event){
    event.preventDefault()
    $.ajax({
        url: $(form).attr('action') || window.location.pathname,
        type: "POST",
        data: fetcheddata,
        success: function (data) {
            console.log("lifts logged")
        },
        error: function (jXHR, textStatus, errorThrown) {
            alert(errorThrown);
        }
    });
}