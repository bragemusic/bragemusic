export default function BodyLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex justify-center w-full">
      <div className="flex flex-col p-2 pt-8 w-full sm:p-4 sm:p-8 sm:pb-8 pb-35 max-w-body text-foreground">
        {children}
      </div>
    </div>
  );
}
