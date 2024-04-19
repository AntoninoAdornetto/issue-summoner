/*
This file provides constants and structures for determining comment syntax in source code files.
The constants define file extensions while the CommentNotation struct specifies
syntax details for various programming languages. The NewCommentNotation function retrieves the
appropriate CommentNotation based on a file extension.

The CommentNotations map contains predefined syntax for common languages such as C, Python, and Markdown,
with a default syntax for unrecognized file types. When walking a project directory, the program reads
each file extension and uses the CommentNotations map to determine the comment syntax for parsing.
*/
package issue
