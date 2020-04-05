function logout() {
  xhr = new XMLHttpRequest();
  xhr.onreadystatechange = (function(x) {
    return function() {
      responseHandler(x, logoutSuccess, logoutFail, logoutLoad);
    };
  })(xhr);
  xhr.open("POST", "/logout", true);
  xhr.send();
}

function logoutSuccess() {
  window.location.replace("http://localhost:8080/login.html");
}

function logoutFail() {
  var msg = "Logout failed with status " + this.status.toString();
  document.getElementById("load-text").innerHTML = msg;
}

function logoutLoad() {
  document.getElementById("load-text").innerHTML = "Loading ...";
}
