function panic() {
    console.log("PANIC")
    let response = fetch('/panic')
    response.then(res => res.json()).then(d => alert(d.response))
}