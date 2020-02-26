(function() {
  const sidebar = document.body.querySelector('.sidebar');
  const sidebarToggle = document.body.querySelector('.sidebar-toggle');
  sidebarToggle.addEventListener('click', () => {
    sidebar.classList.toggle('expand');
  });
})();