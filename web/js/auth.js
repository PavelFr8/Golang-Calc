// js/auth.js

function showAlert(message, type = "danger") {
    const alertBox = `
      <div class="alert alert-${type} alert-dismissible fade show" role="alert">
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
      </div>
    `;
    $('#alert-container').html(alertBox);
  }
  
  $(document).ready(function () {
    // Регистрация
    $('#register-form').on('submit', function (e) {
      e.preventDefault();
      const login = $('#register-username').val();
      const password = $('#register-password').val();
  
      $.ajax({
        url: '/api/v1/register',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({ login, password }),
        success: function () {
          showAlert('Регистрация успешна. Теперь войдите.', 'success');
          $('#register-form')[0].reset();
          $('#nav-login-tab').tab('show');
        },
        error: function (xhr) {
          showAlert(`Ошибка регистрации: ${xhr.responseText}`);
        }
      });
    });
  
    // Логин
    $('#login-form').on('submit', function (e) {
      e.preventDefault();
      const login = $('#login-username').val();
      const password = $('#login-password').val();
  
      $.ajax({
        url: '/api/v1/login',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({ login, password }),
        success: function (response) {
            showAlert('Успешный вход. JWT токен помещён в поле во вкладке "Калькулятор".', 'success');
            $('#jwt-token').val("");
            $('#jwt-token').val(response.token);
        },
        error: function (xhr) {
          showAlert(`Ошибка входа: ${xhr.responseText}`);
        }
      });
    });
  });
  
  