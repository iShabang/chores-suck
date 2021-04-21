function setupShow() {
    elems = document.getElementsByClassName('hidable');
    for (e of elems) {
        e.style.display = 'none';
    }
    showElem = document.getElementById('first');
    tab = window.sessionStorage.getItem('dashTab');
    if (tab) {
        showElem = document.getElementById(tab);
    }
    showElem.style.display = 'block';
}

function display(showElem) {
    divs = document.getElementsByClassName('hidable');
    for (d of divs) {
        d.style.display = 'none';
    }
    document.getElementById(showElem).style.display = 'block';
    window.sessionStorage.setItem('dashTab',showElem);
}

setupShow()
document.getElementById('chore-btn').addEventListener('click', function() {display('first');});
document.getElementById('group-btn').addEventListener('click', function() {display('second');});
