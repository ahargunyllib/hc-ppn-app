"use client";

import { Button } from "@/shared/components/ui/button";
import {
  Popover,
  PopoverPopup,
  PopoverTrigger,
} from "@/shared/components/ui/popover";
import { cn } from "@/shared/lib/utils";
import { FilterIcon, Star, X } from "lucide-react";
import { useState } from "react";

type FeedbackFiltersProps = {
  selectedRatings: number[];
  onRatingsChange: (ratings: number[]) => void;
};

export function FeedbackFilters({
  selectedRatings,
  onRatingsChange,
}: FeedbackFiltersProps) {
  const [open, setOpen] = useState(false);

  const toggleRating = (rating: number) => {
    if (selectedRatings.includes(rating)) {
      onRatingsChange(selectedRatings.filter((r) => r !== rating));
    } else {
      onRatingsChange([...selectedRatings, rating].sort());
    }
  };

  const clearFilters = () => {
    onRatingsChange([]);
    setOpen(false);
  };

  const hasActiveFilters = selectedRatings.length > 0;

  return (
    <Popover onOpenChange={setOpen} open={open}>
      <PopoverTrigger render={<Button size="sm" variant="outline" />}>
        <FilterIcon />
        Filter
        {hasActiveFilters && (
          <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1.5 font-semibold text-primary-foreground text-xs">
            {selectedRatings.length}
          </span>
        )}
      </PopoverTrigger>
      <PopoverPopup align="end" className="w-fit">
        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between">
            <h3 className="font-semibold text-sm">Filter by Rating</h3>
          </div>

          <div className="flex flex-col gap-2">
            {[5, 4, 3, 2, 1].map((rating) => {
              const isSelected = selectedRatings.includes(rating);
              return (
                <Button
                  className="justify-start"
                  key={rating}
                  onClick={() => toggleRating(rating)}
                  size="sm"
                  variant={isSelected ? "default" : "outline"}
                >
                  <Star
                    className={cn(isSelected ? "fill-current" : "fill-none")}
                  />
                  <span>
                    {rating} Star{rating !== 1 ? "s" : ""}
                  </span>
                </Button>
              );
            })}
          </div>

          {hasActiveFilters && (
            <Button onClick={clearFilters} size="xs" variant="ghost">
              <X />
              Clear Filters
            </Button>
          )}
        </div>
      </PopoverPopup>
    </Popover>
  );
}
