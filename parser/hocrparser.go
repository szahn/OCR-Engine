package parser

import (
	"bytes"
	"encoding/xml"
	pb "ocr-engine/grpc"
	"strconv"
	"strings"
)

func DescribeBlock(block *pb.Block) string {
	var builder strings.Builder
	builder.WriteString("[")
	builder.WriteString(block.Tag)
	builder.WriteString(" ")
	builder.WriteString("(")
	builder.WriteString(strconv.Itoa(int(block.X1)))
	builder.WriteString("x")
	builder.WriteString(strconv.Itoa(int(block.X2)))
	builder.WriteString("-")
	builder.WriteString(strconv.Itoa(int(block.Y1)))
	builder.WriteString("x")
	builder.WriteString(strconv.Itoa(int(block.Y2)))
	builder.WriteString(" confidence:")
	builder.WriteString(block.Confidence)
	builder.WriteString(" ")
	builder.WriteString(")")
	builder.WriteString(" '")
	builder.WriteString(block.Text)
	builder.WriteString("' ")
	builder.WriteString("]")

	return builder.String()
}

//Read HOCR spec http://kba.cloud/hocr-spec/1.2/#bbox
func Parse(hocr string) []*pb.Block {

	data := []byte(hocr)
	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)

	isOpen := false

	blocks := make([]*pb.Block, 0)

	current_block := new(pb.Block)
	current_text := new(strings.Builder)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			{
				isOpen = true
				el := xml.StartElement(se)
				current_block.Tag = el.Name.Local

				for _, a := range el.Attr {
					switch a.Name.Local {
					case "title":
						{
							props := strings.Split(a.Value, ";")
							for _, prop := range props {
								propVal := strings.TrimSpace(prop)
								propValues := strings.Split(propVal, " ")
								propHead := propValues[0]
								switch propHead {
								case "bbox":
									{
										x1, _ := strconv.Atoi(propValues[1])
										y1, _ := strconv.Atoi(propValues[2])
										x2, _ := strconv.Atoi(propValues[3])
										y2, _ := strconv.Atoi(propValues[4])

										current_block.X1 = int32(x1)
										current_block.Y1 = int32(y1)
										current_block.X2 = int32(x2)
										current_block.Y2 = int32(y2)

										current_block.HasBox = true
									}
								case "image":
									{

									}
								case "ppageno":
									{

									}
								case "x_wconf":
									{
										current_block.Confidence = propValues[1]
									}
								default:
									{

									}
								}
							}
						}
					}
				}
			}
		case xml.EndElement:
			{
				if isOpen {
					isOpen = false
					current_block.Text = current_text.String()
					blocks = append(blocks, current_block)
					current_block = new(pb.Block)
					current_text = new(strings.Builder)
				}
			}
		case xml.Attr:
			{

			}
		case xml.CharData:
			{
				if !isOpen {
					continue
				}
				text := strings.Trim(string([]byte(se)), " \r\n\t\v")
				current_text.WriteString(text)
			}
		default:
			{

			}
		}
	}

	if current_block.HasBox {
		blocks = append(blocks, current_block)
	}

	return blocks
}
