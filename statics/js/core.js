// Convert nomor ke mata uang rupiah
function toRupiah(number) {
    //number = number.toString().replace(/\s+/g, '');
    if(isNaN(number) || number == '') {
      number = 0;
    }
  
    num = number.toString().split('').reverse().join('');
    num = num.match(/\d{1,3}/g);
   
    num = "Rp " + num.join(',').split('').reverse().join('');
    return num;
  }