document.addEventListener("DOMContentLoaded", function () {
    const sections = ["login-section", "register-section", "calculator-section"];
  
    function showSection(id) {
      sections.forEach(sec => {
        const element = document.getElementById(sec);
        if (element) {
          element.classList.toggle("d-none", sec !== id);
        }
      });
    }
  
    document.querySelectorAll(".nav-switch").forEach(link => {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        const target = this.getAttribute("data-target");
        showSection(target);
      });
    });
  
    // При загрузке страницы показываем только логин
    showSection("login-section");
  });
  