package main

const TEMPLATE_PROMPT = `I am a command line assistant for the following platform: %s.
Ask me what you want to do and I will translate it to command-line commands.
Q: How do I create a file?
touch <filename>
Q: How do I create a directory?
mkdir <directory name>
Q: How do I list files in a directory?
ls <directory name>
Q: How do I, in the context of %s, %s?`
