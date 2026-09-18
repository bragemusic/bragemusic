import { Card, Text } from "@heroui/react";
import { SideScrollContainer } from "./SideScrollContainer";

interface RecentAlbumsProps{}

export const RecentAlbums: React.FC<RecentAlbumsProps> = () => {
    const images = [
        "https://heroui-assets.nyc3.cdn.digitaloceanspaces.com/docs/robot1.jpeg",
        "https://heroui-assets.nyc3.cdn.digitaloceanspaces.com/docs/avocado.jpeg",
        "https://heroui-assets.nyc3.cdn.digitaloceanspaces.com/docs/oranges.jpeg",
    ];

    return (
        <div className="flex flex-col w-full min-w-0">
            <Text type="h4">
            Recently Added Albums
            </Text>
            <SideScrollContainer>
                {Array.from({ length: 10 }).map((_, idx) => (
                    <Card
                        key={`scroll-shadow-lorem-cards-${idx}`}
                        className="flex flex-col gap-3 p-1 min-w-[200px] shrink-0"
                        variant="transparent"
                    >
                        <img
                            alt="Lorem Card"
                            className="object-cover rounded-xl select-none sm:w-40 sm:h-40 w-30 h-30 aspect-square shrink-0"
                            loading="lazy"
                            src={images[idx % images.length]}
                        />
                        <div className="flex flex-col flex-1 gap-1 justify-center">
                            <Card.Title className="text-sm">
                                Bridging the Future
                            </Card.Title>
                            <Card.Description className="text-xs">
                                Today, 6:30 PM
                            </Card.Description>
                        </div>
                    </Card>
                ))}
            </SideScrollContainer>
        </div>
    );
};
