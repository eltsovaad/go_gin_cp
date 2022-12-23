$(document).ready(function(){
    var active = document.getElementsByClassName("item")[0];
    active.classList.add('active');
    $('#itemslider').carousel({ interval: 0 });

    $('.carousel-showmanymoveone .item').each(function(){
        var itemToClone = $(this);

        for (i in document.getElementsByClassName("item")) {
        itemToClone = itemToClone.next();

        if (!itemToClone.length) {
        itemToClone = $(this).siblings(':first');
        }

        itemToClone.children(':first-child').clone()
        .addClass("cloneditem-"+(i))
        .appendTo($(this));
    }
    });
});