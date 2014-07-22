(function() {

	var s = io.connect()
	s.on("data", function (e){
		console.log(e);
	})

})()

