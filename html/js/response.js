function responseHandler(xhr, goodFunc, badFunc, loadFunc) {
  if (xhr.readyState < 4) {
    loadFunc();
  } else if (xhr.readyState == 4) {
    if (xhr.status == 200) {
      goodFunc();
    } else {
      badFunc();
    }
  }
}
