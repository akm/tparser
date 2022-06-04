package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestStrucRecType(t *testing.T) {
	NewTypeTest(t,
		"TDateRec",
		[]rune(`
record
	Year: Integer;
	Month: (Jan, Feb, Mar, Apr, May, Jun,
			Jul, Aug, Sep, Oct, Nov, Dec);
	Day: 1..31;
end
`),
		&ast.RecType{
			FieldList: &ast.FieldList{
				FieldDecls: ast.FieldDecls{
					{
						IdentList: asttest.NewIdentList("Year"),
						Type:      &ast.OrdIdent{Ident: asttest.NewIdent("Integer")},
					},
					{
						IdentList: asttest.NewIdentList("Month"),
						Type: ast.EnumeratedType{
							{Ident: asttest.NewIdent("Jan")},
							{Ident: asttest.NewIdent("Feb")},
							{Ident: asttest.NewIdent("Mar")},
							{Ident: asttest.NewIdent("Apr")},
							{Ident: asttest.NewIdent("May")},
							{Ident: asttest.NewIdent("Jun")},
							{Ident: asttest.NewIdent("Jul")},
							{Ident: asttest.NewIdent("Aug")},
							{Ident: asttest.NewIdent("Sep")},
							{Ident: asttest.NewIdent("Oct")},
							{Ident: asttest.NewIdent("Nov")},
							{Ident: asttest.NewIdent("Dec")},
						},
					},
					{
						IdentList: asttest.NewIdentList("Day"),
						Type: &ast.SubrangeType{
							Low:  asttest.NewConstExpr(asttest.NewNumber("1")),
							High: asttest.NewConstExpr(asttest.NewNumber("31")),
						},
					},
				},
			},
		},
	).Run().RunTypeSection("TDateRec")

	NewTypeTest(t,
		"S",
		[]rune(`
record
	Name: string;
	Age: Integer;
end
`),
		&ast.RecType{
			FieldList: &ast.FieldList{
				FieldDecls: ast.FieldDecls{
					{
						IdentList: asttest.NewIdentList("Name"),
						Type:      &ast.StringType{Name: "STRING"},
					},
					{
						IdentList: asttest.NewIdentList("Age"),
						Type:      &ast.OrdIdent{Ident: asttest.NewIdent("Integer")},
					},
				},
			},
		},
	).Run().RunVarSection("S")

	NewTypeTest(t,
		"TEmployee",
		[]rune(`
record
	FirstName, LastName: string[40];
	BirthDate: TDate;
	case Salaried: Boolean of
		True: (AnnualSalary: Currency);
		False: (HourlyWage: Currency); 
end
`),
		&ast.RecType{
			FieldList: &ast.FieldList{
				FieldDecls: ast.FieldDecls{
					{
						IdentList: asttest.NewIdentList("FirstName", "LastName"),
						Type:      &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("40"))},
					},
					{
						IdentList: asttest.NewIdentList("BirthDate"),
						Type:      &ast.TypeId{Ident: asttest.NewIdent("TDate")},
					},
				},
				VariantSection: &ast.VariantSection{
					Ident:  asttest.NewIdent("Salaried"),
					TypeId: &ast.OrdIdent{Ident: asttest.NewIdent("Boolean")},
					RecVariants: ast.RecVariants{
						{
							ConstExprs: ast.ConstExprs{asttest.NewConstExpr(&ast.ValueFactor{Value: "True"})},
							FieldList: &ast.FieldList{
								FieldDecls: ast.FieldDecls{
									{
										IdentList: asttest.NewIdentList("AnnualSalary"),
										Type:      &ast.RealType{Ident: asttest.NewIdent("Currency")},
									},
								},
							},
						},
						{
							ConstExprs: ast.ConstExprs{asttest.NewConstExpr(&ast.ValueFactor{Value: "False"})},
							FieldList: &ast.FieldList{
								FieldDecls: ast.FieldDecls{
									{
										IdentList: asttest.NewIdentList("HourlyWage"),
										Type:      &ast.RealType{Ident: asttest.NewIdent("Currency")},
									},
								},
							},
						},
					},
				},
			},
		},
	).Run().RunTypeSection("TEmployee")

	NewTypeTest(t,
		"TPerson",
		[]rune(`
record
	FirstName, LastName: string[40];
	BirthDate: TDate;
	case Citizen: Boolean of
		True: (Birthplace: string[40]);
		False: (Country: string[20];
				EntryPort: string[20]; 
				EntryDate, ExitDate: TDate);
end
`),
		&ast.RecType{
			FieldList: &ast.FieldList{
				FieldDecls: ast.FieldDecls{
					{
						IdentList: asttest.NewIdentList("FirstName", "LastName"),
						Type:      &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("40"))},
					},
					{
						IdentList: asttest.NewIdentList("BirthDate"),
						Type:      &ast.TypeId{Ident: asttest.NewIdent("TDate")},
					},
				},
				VariantSection: &ast.VariantSection{
					Ident:  asttest.NewIdent("Citizen"),
					TypeId: &ast.OrdIdent{Ident: asttest.NewIdent("Boolean")},
					RecVariants: ast.RecVariants{
						{
							ConstExprs: ast.ConstExprs{asttest.NewConstExpr(&ast.ValueFactor{Value: "True"})},
							FieldList: &ast.FieldList{
								FieldDecls: ast.FieldDecls{
									{
										IdentList: asttest.NewIdentList("Birthplace"),
										Type:      &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("40"))},
									},
								},
							},
						},
						{
							ConstExprs: ast.ConstExprs{asttest.NewConstExpr(&ast.ValueFactor{Value: "False"})},
							FieldList: &ast.FieldList{
								FieldDecls: ast.FieldDecls{
									{
										IdentList: asttest.NewIdentList("Country"),
										Type:      &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("20"))},
									},
									{
										IdentList: asttest.NewIdentList("EntryPort"),
										Type:      &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("20"))},
									},
									{
										IdentList: asttest.NewIdentList("EntryDate", "ExitDate"),
										Type:      &ast.TypeId{Ident: asttest.NewIdent("TDate")},
									},
								},
							},
						},
					},
				},
			},
		},
	).Run().RunTypeSection("TPerson")

	NewTypeSectionTest(t,
		"with TShapeList",
		[]rune(`
type
	TShapeList = (Rectangle, Triangle, Circle, Ellipse, Other);
	TFigure = record
		case TShapeList of
			Rectangle: (Height, Width: Real);
			Triangle: (Side1, Side2, Angle: Real);
			Circle: (Radius: Real);
			Ellipse, Other: ();
	end;
`),
		ast.TypeSection{
			{
				Ident: asttest.NewIdent("TShapeList"),
				Type: ast.EnumeratedType{
					{Ident: asttest.NewIdent("Club")},
					{Ident: asttest.NewIdent("Diamond")},
					{Ident: asttest.NewIdent("Heart")},
					{Ident: asttest.NewIdent("Spade")},
				},
			},
			{
				Ident: asttest.NewIdent("TFigure"),
				Type: &ast.RecType{
					FieldList: &ast.FieldList{
						VariantSection: &ast.VariantSection{
							TypeId: &ast.TypeId{Ident: asttest.NewIdent("TShapeList")},
							RecVariants: ast.RecVariants{
								{
									ConstExprs: ast.ConstExprs{asttest.NewConstExpr(&ast.ValueFactor{Value: "Rectangle"})},
									FieldList: &ast.FieldList{
										FieldDecls: ast.FieldDecls{
											{
												IdentList: asttest.NewIdentList("Height", "Width"),
												Type:      &ast.RealType{Ident: asttest.NewIdent("REAL")},
											},
										},
									},
								},
								{
									ConstExprs: ast.ConstExprs{asttest.NewConstExpr(&ast.ValueFactor{Value: "Triangle"})},
									FieldList: &ast.FieldList{
										FieldDecls: ast.FieldDecls{
											{
												IdentList: asttest.NewIdentList("Side1", "Side2", "Angle"),
												Type:      &ast.RealType{Ident: asttest.NewIdent("REAL")},
											},
										},
									},
								},
								{
									ConstExprs: ast.ConstExprs{asttest.NewConstExpr(&ast.ValueFactor{Value: "Circle"})},
									FieldList: &ast.FieldList{
										FieldDecls: ast.FieldDecls{
											{
												IdentList: asttest.NewIdentList("Radius"),
												Type:      &ast.RealType{Ident: asttest.NewIdent("REAL")},
											},
										},
									},
								},
								{
									ConstExprs: ast.ConstExprs{
										asttest.NewConstExpr(&ast.ValueFactor{Value: "Ellipse"}),
										asttest.NewConstExpr(&ast.ValueFactor{Value: "Other"}),
									},
									FieldList: &ast.FieldList{
										FieldDecls: ast.FieldDecls{},
									},
								},
							},
						},
					},
				},
			},
		},
	)
}