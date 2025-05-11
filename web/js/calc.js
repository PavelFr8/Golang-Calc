// js/calc.js

$(document).ready(function () {
    // Отправка выражения
    $('#calculate-form').on('submit', function (e) {
      e.preventDefault();
      const token = $('#jwt-token').val();
      const expression = $('#expression').val();
  
      if (!token) {
        showAlert('Введите JWT токен.', 'warning');
        return;
      }
  
      $.ajax({
        url: '/api/v1/calculate',
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + token },
        contentType: 'application/json',
        data: JSON.stringify({ expression }),
        success: function (response) {
          showAlert(`Выражение отправлено. ID: ${response.id}`, 'success');
        },
        error: function (xhr) {
          showAlert(`Ошибка: ${xhr.responseText}`);
        }
      });
    });
  
    // Загрузка всех выражений
    $('#load-expressions').on('click', function () {
      const token = $('#jwt-token').val();
  
      if (!token) {
        showAlert('Введите JWT токен.', 'warning');
        return;
      }
  
      $.ajax({
        url: '/api/v1/expressions',
        method: 'GET',
        headers: { 'Authorization': 'Bearer ' + token },
        success: function (response) {
          const expressions = response.expressions;
          const tbody = $('#expressions-table-body');
          tbody.empty();
  
          expressions.forEach(expr => {
            tbody.append(`
              <tr>
                <td>${expr.id}</td>
                <td>${expr.expression}</td>
                <td>${expr.status}</td>
                <td>${expr.result || 'Не доступен'}</td>
              </tr>
            `);
          });
        },
        error: function (xhr) {
          showAlert(`Ошибка: ${xhr.responseText}`);
        }
      });
    });
  
    // Получение деталей выражения
    $('#get-expression-form').on('submit', function (e) {
      e.preventDefault();
      const token = $('#jwt-token').val();
      const exprId = $('#expr-id').val();
  
      if (!token) {
        showAlert('Введите JWT токен.', 'warning');
        return;
      }
  
      $.ajax({
        url: `/api/v1/expressions/${exprId}`,
        method: 'GET',
        headers: { 'Authorization': 'Bearer ' + token },
        success: function (response) {
          const expr = response.expression;
          $('#expression-details').html(`
            <h4>Детали выражения</h4>
            <p><strong>ID:</strong> ${expr.id}</p>
            <p><strong>Выражение:</strong> ${expr.expression}</p>
            <p><strong>Статус:</strong> ${expr.status}</p>
            <p><strong>Результат:</strong> ${expr.result || 'Не доступен'}</p>
          `);
        },
        error: function (xhr) {
          showAlert(`Ошибка: ${xhr.responseText}`);
        }
      });
    });
  });
  