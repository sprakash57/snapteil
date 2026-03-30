import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TagFilter from "@/features/feed/TagFilter";

const allTags = ["nature", "portrait", "street"];

describe("TagFilter", () => {
  it("renders the filter trigger button", () => {
    render(
      <TagFilter allTags={allTags} selectedTags={[]} onChange={vi.fn()} />,
    );
    expect(
      screen.getByRole("button", { name: /filter by tags/i }),
    ).toBeInTheDocument();
  });

  it("does not show tag list before button is clicked", () => {
    render(
      <TagFilter allTags={allTags} selectedTags={[]} onChange={vi.fn()} />,
    );
    expect(screen.queryByText("nature")).not.toBeInTheDocument();
  });

  it("shows tag options after clicking the trigger", async () => {
    const user = userEvent.setup();
    render(
      <TagFilter allTags={allTags} selectedTags={[]} onChange={vi.fn()} />,
    );
    await user.click(screen.getByRole("button", { name: /filter by tags/i }));
    expect(screen.getByText("nature")).toBeInTheDocument();
    expect(screen.getByText("portrait")).toBeInTheDocument();
    expect(screen.getByText("street")).toBeInTheDocument();
  });

  it("calls onChange with the tag added when clicking an unselected tag", async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(
      <TagFilter allTags={allTags} selectedTags={[]} onChange={onChange} />,
    );
    await user.click(screen.getByRole("button", { name: /filter by tags/i }));
    await user.click(screen.getByText("nature"));
    expect(onChange).toHaveBeenCalledWith(["nature"]);
  });

  it("calls onChange with the tag removed when clicking an active tag", async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(
      <TagFilter
        allTags={allTags}
        selectedTags={["nature"]}
        onChange={onChange}
      />,
    );
    await user.click(screen.getByRole("button", { name: /filter by tags/i }));
    // Use getByRole to target the dropdown button (not the pill span)
    await user.click(screen.getByRole("button", { name: "nature" }));
    expect(onChange).toHaveBeenCalledWith([]);
  });

  it("shows selected tags as pills outside the dropdown", () => {
    render(
      <TagFilter
        allTags={allTags}
        selectedTags={["nature", "street"]}
        onChange={vi.fn()}
      />,
    );
    // Pills are visible without opening the dropdown
    const pills = screen.getAllByText(/^(nature|street)$/);
    expect(pills.length).toBeGreaterThanOrEqual(2);
  });

  it("removes tag when × button on its pill is clicked", async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(
      <TagFilter
        allTags={allTags}
        selectedTags={["nature"]}
        onChange={onChange}
      />,
    );
    // The × button is inside the pill
    await user.click(screen.getByRole("button", { name: "×" }));
    expect(onChange).toHaveBeenCalledWith([]);
  });

  it("closes the dropdown when clicking outside", async () => {
    const user = userEvent.setup();
    render(
      <TagFilter allTags={allTags} selectedTags={[]} onChange={vi.fn()} />,
    );
    await user.click(screen.getByRole("button", { name: /filter by tags/i }));
    expect(screen.getByText("nature")).toBeInTheDocument();

    await user.click(document.body);
    expect(screen.queryByText("nature")).not.toBeInTheDocument();
  });

  it("does not open dropdown when allTags is empty", async () => {
    const user = userEvent.setup();
    render(<TagFilter allTags={[]} selectedTags={[]} onChange={vi.fn()} />);
    await user.click(screen.getByRole("button", { name: /filter by tags/i }));
    // Nothing to show, dropdown element should not exist
    expect(screen.queryByRole("listbox")).not.toBeInTheDocument();
  });
});
