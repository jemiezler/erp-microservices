import { cn } from "../lib/utils";
import { Card } from "./ui/card";
export interface AppCardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
};

export const AppCard = ({ children, className, ...props }: AppCardProps) => {
  return (
    <Card className={cn('shadow-none bg-card-background rounded-2xl p-6',className)} {...props}>
      {children}
    </Card>
  );
}
