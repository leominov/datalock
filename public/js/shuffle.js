$(function() {
	const CookieName = "shuffle";
	function IsShuffleCookieExists() {
		if (typeof $.cookie(CookieName) != "undefined") {
			return true;
		}
		return false;
	}
	function Init() {
		if (IsShuffleCookieExists()) {
			$(".shuffle-button").addClass("act");
		}
	}
	$(".shuffle-button").click(function(){
		if (IsShuffleCookieExists()) {
			$.removeCookie(CookieName);
		} else {
			$.cookie(CookieName, Date.now(), {expires: 365});
		}
		$(".shuffle-button").toggleClass("act");
		location.reload(true);
	});
	Init();
});
