let charMap = new Map();

window.addEventListener('DOMContentLoaded', function(){
    //set our map with the good characters
    mapGoodCharSetter();
});

function mapGoodCharSetter(){
    //Lowercase Letters
    charMap.set('a', 'a');
    charMap.set('b', 'b');
    charMap.set('c', 'c');
    charMap.set('d', 'd');
    charMap.set('e', 'e');
    charMap.set('f', 'f');
    charMap.set('g', 'g');
    charMap.set('h', 'h');
    charMap.set('i', 'i');
    charMap.set('j', 'j');
    charMap.set('k', 'k');
    charMap.set('l', 'l');
    charMap.set('m', 'm');
    charMap.set('n', 'n');
    charMap.set('o', 'o');
    charMap.set('p', 'p');
    charMap.set('q', 'q');
    charMap.set('r', 'r');
    charMap.set('s', 's');
    charMap.set('t', 't');
    charMap.set('u', 'u');
    charMap.set('v', 'v');
    charMap.set('w', 'w');
    charMap.set('x', 'x');
    charMap.set('y', 'y');
    charMap.set('z', 'z');

    //Uppercase letters
    charMap.set('a'.toUpperCase(), 'a'.toUpperCase());
    charMap.set('b'.toUpperCase(), 'b'.toUpperCase());
    charMap.set('c'.toUpperCase(), 'c'.toUpperCase());
    charMap.set('d'.toUpperCase(), 'd'.toUpperCase());
    charMap.set('e'.toUpperCase(), 'e'.toUpperCase());
    charMap.set('f'.toUpperCase(), 'f'.toUpperCase());
    charMap.set('g'.toUpperCase(), 'g'.toUpperCase());
    charMap.set('h'.toUpperCase(), 'h'.toUpperCase());
    charMap.set('i'.toUpperCase(), 'i'.toUpperCase());
    charMap.set('j'.toUpperCase(), 'j'.toUpperCase());
    charMap.set('k'.toUpperCase(), 'k'.toUpperCase());
    charMap.set('l'.toUpperCase(), 'l'.toUpperCase());
    charMap.set('m'.toUpperCase(), 'm'.toUpperCase());
    charMap.set('n'.toUpperCase(), 'n'.toUpperCase());
    charMap.set('o'.toUpperCase(), 'o'.toUpperCase());
    charMap.set('p'.toUpperCase(), 'p'.toUpperCase());
    charMap.set('q'.toUpperCase(), 'q'.toUpperCase());
    charMap.set('r'.toUpperCase(), 'r'.toUpperCase());
    charMap.set('s'.toUpperCase(), 's'.toUpperCase());
    charMap.set('t'.toUpperCase(), 't'.toUpperCase());
    charMap.set('u'.toUpperCase(), 'u'.toUpperCase());
    charMap.set('v'.toUpperCase(), 'v'.toUpperCase());
    charMap.set('w'.toUpperCase(), 'w'.toUpperCase());
    charMap.set('x'.toUpperCase(), 'x'.toUpperCase());
    charMap.set('y'.toUpperCase(), 'y'.toUpperCase());
    charMap.set('z'.toUpperCase(), 'z'.toUpperCase());

    //Numbers
    charMap.set('0', '0');
    charMap.set('1', '1');
    charMap.set('2', '2');
    charMap.set('3', '3');
    charMap.set('4', '4');
    charMap.set('5', '5');
    charMap.set('6', '6');
    charMap.set('7', '7');
    charMap.set('8', '8');
    charMap.set('9', '9');

    //Allowed Characters
    charMap.set('!', '!');
    charMap.set(')', ')');
    charMap.set('(', '(');
    charMap.set('_', '_');
    charMap.set('-', '-');
    charMap.set('>', '>');
    charMap.set('<', '<');
    charMap.set('?', '?');
    
    /* debug
    for (const [key, value] of charMap.entries()){
        console.log("Here is the charMap: " + key + " and value: " + value);
    }

    */
}

function checkInput(theInput){
    var goodCharacters = true; // A check returned if we have good characters
    for (var i = 0; i < theInput.length; i++){
        if (charMap.has(theInput[i]) != true){
            console.log("Wrong character: '" + theInput[i] + "', not allowed!");
            goodCharacters = false;
            break;
        } 
    }

    return goodCharacters;
}