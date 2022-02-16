Date.prototype.addHours = function(h) {
    this.setTime(this.getTime() + (h * 60 * 60 * 1000));
    return this;
}


// Set the date we're counting down to 
console.log("Making a Timer");
var timerDate = document.getElementById("timer-date").innerHTML;
console.log(timerDate);
var split = timerDate.split(" ");
var timerDateString = split[0] + "T" + split[1];

var countDownDate = new Date(timerDateString);
var timerDuration = document.getElementById("timer-duration").innerHTML;
countDownDate.addHours(timerDuration);
countDownDate = countDownDate.getTime();
// Update the count down every 1 second
var x = setInterval(function() {

    // Get today's date and time
    var now = new Date().getTime();

    // Find the distance between now and the count down date
    var distance = countDownDate - now;

    // Time calculations for days, hours, minutes and seconds

    var hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    var seconds = Math.floor((distance % (1000 * 60)) / 1000);

    // Display the result in the element with id="timer"
    document.getElementById("timer").innerHTML = hours + "h " +
        minutes + "m " + seconds + "s ";

    // If the count down is finished, write some text
    if (distance < 0) {
        clearInterval(x);
        document.getElementById("timer").innerHTML = "EXPIRED";
    }
}, 1000);