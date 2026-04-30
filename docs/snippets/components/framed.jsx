export const Framed = ({
    image,
    background,
    link,
    alt = "Framed content",
    borderWidth = "24px",
    outerRadius = "10px",
    innerRadius = "6px",
}) => {
    // Parse CSS shorthand (1-4 values) into array of 4
    const parseShorthand = (value) => {
        const values = String(value).trim().split(/\s+/);
        switch (values.length) {
            case 1: return [values[0], values[0], values[0], values[0]];
            case 2: return [values[0], values[1], values[0], values[1]];
            case 3: return [values[0], values[1], values[2], values[1]];
            default: return values.slice(0, 4);
        }
    };

    const isZero = (v) => parseFloat(v) === 0;

    // padding: [top, right, bottom, left]
    const padding = parseShorthand(borderWidth);
    // radius: [top-left, top-right, bottom-right, bottom-left]
    const outer = parseShorthand(outerRadius);

    // Build mask to exclude background from corners where BOTH adjacent paddings are 0
    const cornerMasks = [];
    // top-left: padding[0]=top, padding[3]=left
    if (isZero(padding[0]) && isZero(padding[3])) {
        cornerMasks.push(`radial-gradient(circle at 0% 0%, transparent ${outer[0]}, black ${outer[0]})`);
    }
    // top-right: padding[0]=top, padding[1]=right
    if (isZero(padding[0]) && isZero(padding[1])) {
        cornerMasks.push(`radial-gradient(circle at 100% 0%, transparent ${outer[1]}, black ${outer[1]})`);
    }
    // bottom-right: padding[2]=bottom, padding[1]=right
    if (isZero(padding[2]) && isZero(padding[1])) {
        cornerMasks.push(`radial-gradient(circle at 100% 100%, transparent ${outer[2]}, black ${outer[2]})`);
    }
    // bottom-left: padding[2]=bottom, padding[3]=left
    if (isZero(padding[2]) && isZero(padding[3])) {
        cornerMasks.push(`radial-gradient(circle at 0% 100%, transparent ${outer[3]}, black ${outer[3]})`);
    }
    
    // Combine masks - each mask punches a hole, we need to intersect them
    const backgroundMask = cornerMasks.length > 0 ? cornerMasks.join(", ") : undefined;

    const innerImage = (
        <img
            src={image}
            alt={alt}
            noZoom={link ? true : false}
            style={{
                display: "block",
                width: "100%",
                margin: 0,
                borderRadius: 0,
            }}
        />
    );

    return (
        <div
            style={{
                display: "inline-block",
                position: "relative",
                borderRadius: outerRadius,
                overflow: "hidden",
            }}
        >
            {/* Background layer - masked to exclude corners with no padding */}
            {background && (
                <div
                    style={{
                        position: "absolute",
                        inset: 0,
                        backgroundImage: `url(${background})`,
                        backgroundSize: "cover",
                        maskImage: backgroundMask,
                        WebkitMaskImage: backgroundMask,
                        maskComposite: "intersect",
                        WebkitMaskComposite: "source-in",
                    }}
                />
            )}
            {/* Content layer with padding */}
            <div style={{ position: "relative", padding: borderWidth }}>
                <div
                    style={{
                        borderRadius: innerRadius,
                        overflow: "hidden",
                        lineHeight: 0,
                    }}
                >
                    {link ? (<a href={link}>{innerImage}</a>) : innerImage}
                </div>
            </div>
        </div>
    );
};
