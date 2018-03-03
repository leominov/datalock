$(function() {
	const cookieName = "shuffle";
	function IsShuffleCookieExists() {
		if (typeof $.cookie(cookieName) != "undefined") {
			return true;
		}
		return false;
	}
	if (IsShuffleCookieExists()) {
		$(".shuffle-button").addClass("act");
	}
	$(".shuffle-button").click(function(){
		if (IsShuffleCookieExists()) {
			$.removeCookie(cookieName);
		} else {
			$.cookie(cookieName, Date.now())
		}
		$(".shuffle-button").toggleClass("act");
		location.reload(true);
	});
});
