# Image

Load image as RGB.

## API Surface

`Image.GetAsRGB(reader io.Reader)` decodes the input image and returns an `*image.RGBA` plus the detected format name.

## Notes

Use this package when downstream image processing expects a mutable RGBA buffer instead of the original decoded image type.
