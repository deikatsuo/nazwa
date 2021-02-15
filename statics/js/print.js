function jspmWSStatus() {
    //Check JSPM WebSocket status
    if (JSPM.JSPrintManager.websocket_status == JSPM.WSStatus.Open) {
      return true;
    } else if (JSPM.JSPrintManager.websocket_status == JSPM.WSStatus.Closed) {
      console.warn('JSPrintManager (JSPM) is not installed or not running! Download JSPM Client App from https://neodynamic.com/downloads/jspm');
      return false;
    } else if (JSPM.JSPrintManager.websocket_status == JSPM.WSStatus.Blocked) {
      alert('JSPM has blocked this website!');
      return false;
    }
  }