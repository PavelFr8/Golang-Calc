// Уведомления
function showAlert(message, type) {
  const alertBox = `
    <div class="alert alert-custom alert-dismissible fade show" role="alert">
      ${message}
      <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    </div>
  `;
  $('#alert-container').html(alertBox);
}

// Обработка отправки выражения
$("#calculate-form").on("submit", function(event) {
  event.preventDefault();
  const expression = $("#expression").val();

  $.ajax({
    url: "/api/v1/calculate",
    method: "POST",
    contentType: "application/json",
    data: JSON.stringify({ expression }),
    success: function(response) {
      showAlert(`Выражение отправлено успешно! ID: ${response.id}`, "success");
    },
    error: function(xhr, status, error) {
      showAlert(`Ошибка: ${xhr.responseText}`, "danger");
    }
  });
});

// Загрузка всех выражений
$("#load-expressions").on("click", function() {
  $.ajax({
    url: "/api/v1/expressions",
    method: "GET",
    success: function(response) {
      const expressions = response.expressions;
      const tbody = $("#expressions-table-body");
      tbody.empty();

      expressions.forEach(function(expr) {
        const row = `
          <tr>
            <td>${expr.id}</td>
            <td>${expr.expression}</td>
            <td>${expr.status}</td>
            <td>${expr.result || 'Не доступен'}</td>
          </tr>
        `;
        tbody.append(row);
      });
    },
    error: function(xhr, status, error) {
      showAlert(`Ошибка: ${xhr.responseText}`, "danger");
    }
  });
});

// Получение деталей выражения
$("#get-expression-form").on("submit", function(event) {
  event.preventDefault();
  const exprId = $("#expr-id").val();

  $.ajax({
    url: `/api/v1/expressions/${exprId}`,
    method: "GET",
    success: function(response) {
      const expr = response.expression;
      const details = `
        <h4>Детали выражения</h4>
        <p><strong>ID:</strong> ${expr.id}</p>
        <p><strong>Выражение:</strong> ${expr.expression}</p>
        <p><strong>Статус:</strong> ${expr.status}</p>
        <p><strong>Результат:</strong> ${expr.result || 'Не доступен'}</p>
      `;
      $("#expression-details").html(details);
    },
    error: function(xhr, status, error) {
      showAlert(`Ошибка: ${xhr.responseText}`, "danger");
    }
  });
});
