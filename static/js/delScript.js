function delItem() {
    if (confirm("Вы уверены что хотите удалить элемент?") === true) {
        let xhr = new XMLHttpRequest();
        let url = document.getElementById("delLink").getAttribute("href")

// 2. Конфигурируем его: GET-запрос на URL 'phones.json'
        xhr.open('GET', url, false);

// 3. Отсылаем запрос
        xhr.send();

// 4. Если код ответа сервера не 200, то это ошибка
        if (xhr.status !== 200) {
            // обработать ошибку
            alert( xhr.status + ': ' + xhr.statusText ); // пример вывода: 404: Not Found
        }
    }

}
