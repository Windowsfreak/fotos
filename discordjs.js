const Discord = require('discord.js');
const fetch = require('node-fetch');
const querystring = require('querystring');

const client = new Discord.Client();
const prefix = '!bot ';

const trim = (str, max) => (str.length > max ? `${str.slice(0, max - 3)}...` : str);

async function showRandomPost() {
    const body = {a: 1};
    //const response = await fetch('https://www.roofing.run/discordjs/backend.php', { method: 'post', body: JSON.stringify(body), headers: {'Content-Type': 'application/json'}});


    const response = await fetch('http://roofing.run:42888/pictures/random', { method: 'post', body: JSON.stringify(body), headers: {'Content-Type': 'application/json'}});

    const json = await response.json();

    console.log(json);
    var username = json.userName;
    var fullpath = "../html/images/Original/" + json.gallery + "/" + json.filename;
    var filename = json.filename;
    console.log('filename' + filename + "\n");

    const exampleEmbed = new Discord.MessageEmbed()
        .setColor('#0099ff')
        .setTitle('Daily Feature - Click for more shots done by @' + username)
        .setURL('https://www.roofing.run/' + username)
        .attachFiles([fullpath])
        .setImage('attachment://' + filename)
        .setFooter('Want your own shots featured? Type "!bot help" to figure out how.', 'https://i.imgur.com/wSTFkRM.png');

    // send back "Pong." to the channel the message was sent in
    //message.channel.send('Pong.');
    client.channels.cache.get('471612953284706315').send(exampleEmbed).then(result => {
        console.log('embed sent');
        client.destroy();
    });
}

client.once('ready', () => {
    console.log('Ready!');
    showRandomPost();
});

client.login('xxx');
