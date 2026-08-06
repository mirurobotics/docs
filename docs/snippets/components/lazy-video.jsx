export const LazyVideo = ({ src, alt, className, aspectRatio }) => {
    const ref = React.useRef(null);

    React.useEffect(() => {
        const video = ref.current;
        if (!video) return;

        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    video.play().catch(() => {});
                } else {
                    video.pause();
                }
            },
            { threshold: 0.25 },
        );

        observer.observe(video);
        return () => observer.disconnect();
    }, []);

    return (
        <video
            ref={ref}
            loop
            muted
            controls
            preload="none"
            playsInline
            alt={alt}
            className={["w-full", className].filter(Boolean).join(" ")}
            style={aspectRatio ? { aspectRatio } : undefined}
            src={src}
        />
    );
};
