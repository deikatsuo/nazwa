// Convert nomor ke mata uang rupiah
function toRupiah(number) {
  //number = number.toString().replace(/\s+/g, '');
  if (isNaN(number) || number == '') {
    number = 0;
  }

  num = number.toString().split('').reverse().join('');
  num = num.match(/\d{1,3}/g);

  num = "Rp " + num.join(',').split('').reverse().join('');
  return num;
}

// Detect OS client
function operatingSytem() {
  var os = "Unknown OS";
  if (navigator.appVersion.indexOf("Win") != -1) os = "Windows";
  if (navigator.appVersion.indexOf("Mac") != -1) os = "Mac";
  if (navigator.appVersion.indexOf("X11") != -1) os = "Unix";
  if (navigator.appVersion.indexOf("Linux") != -1) os = "Linux";

  // Kembalikan nilai os
  return os;
}

function itemsToString(arr) {
  items = "";
  for (i = 0; i < arr.length; i++) {
    items = `${items}[${arr[i].Quantity}x ${arr[i].Product.Name}]`;
    if (arr.length > 1 && i < (arr.length - 2)) {
      items = `${items}, `;
    }
    if (arr.length > 1 && i == (arr.length - 2)) {
      items = `${items} dan `;
    }
  }

  return items;
}

function future(due) {
  now = Date.now();
  ds = due.split("/");
  td = new Date(`${ds[2]}-${ds[1]}-${ds[0]}`);

  return td > now;
}